// we will put all the go struct , we will read all the environment variable from .env and store the values in the go struct and we will resuse it accross the runtime
// of application without reading from environment again and again
// the reading will be like this we will read it once inside main.go and stores it in the structs and we will pass and resuse again and again

// Concept -- reflection is basically understanding structure of data at runtime and struct tags are one of the methods for validation and can be used during runtime
// in go tags are defined using backticks

//Each library can define some kind of key using which we can pass some kind of metadata for example if we are using koanf (a tag) , using this library tag we can pass the
// metadata to the struct which library will be able to parse and perform various operation

package config

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

// parent struct (should create child struct for keeping this Config struct clean)
type Config struct {
	Primary       Primary              `koanf:"primary" validate:"required"`
	Server        ServerConfig         `koanf:"server" validate:"required"`
	Database      DatabaseConfig       `koanf:"database" validate:"required"`
	Auth          AuthConfig           `koanf:"auth" validate:"required"`
	Redis         RedisConfig          `koanf:"redis" validate:"required"`
	Observability *ObservabilityConfig `koanf:"observability"`
}

// primary config struct
// for primary will need only one variable for tracking the environment
type Primary struct {
	Env string `koanf:"env" validate:"required"`
}

// server config struct
type ServerConfig struct {
	Port               string   `koanf:"port" validate:"required"`
	ReadTimeout        int      `koanf:"read_timeout" validate:"required"`
	WriteTimeout       int      `koanf:"write_timeout" validate:"reqruired"`
	IdleTimeout        int      `koanf:"idle_timeout" validate:"required"`
	CORSAllowedOrigins []string `koanf:"cors_allowed_origins" validate:"required"`
}

// database config struct
type DatabaseConfig struct {
	Host            string `koanf:"host" validate:"required"`
	Port            string `koanf:"port" validate:"required"`
	User            string `koanf:"user" validate:"required"`
	Password        string `koanf:"password" validate:"required"`
	Name            string `koanf:"name" validate:"required"`
	SSLMode         string `koanf:"ssl_mode" validate:"required"`
	MaxOpenConns    int    `koanf:"max_open_conns" validate:"required"`
	MaxIdleConns    int    `koanf:"max_idle_conns" validate:"required"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime int    `koanf:"conn_max_idle_time" validate:"required"`
}

// Auth config struct (using clerk for authentication)
type AuthConfig struct {
	SecretKey string `koanf:"secret_key" validate:"reqruired"`
}

// Redis config struct (for caching and session management)
type RedisConfig struct {
	Address string `koanf:"address" validate:"required"`
}

func LoadConfig() (*Config, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	k := koanf.New(".")

	err := k.Load(env.Provider("BOILERPLATE_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "BOILERPLATE_"))
	}), nil)

	if err != nil {
		logger.Fatal().Err(err).Msg("could not load initial environment variables")
	}

	mainConfig := &Config{}

	err = k.Unmarshal("", mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not unmarhsal config")
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(mainConfig)

	if err != nil {
		logger.Fatal().Err(err).Msg("Config validation failed")
	}

	if mainConfig.Observability == nil {
		mainConfig.Observability = DefaultObservabilityConfig()
	}

	mainConfig.Observability.ServiceName = "boilerplate"
	mainConfig.Observability.Environment = mainConfig.Primary.Env

	if err := mainConfig.Observability.Validate(); err != nil {
		logger.Fatal().Err(err).Msg("invalid Observability config")
	}

	return mainConfig, nil
}

// External Library -- validator, koanf, godotenv, zerolog.
