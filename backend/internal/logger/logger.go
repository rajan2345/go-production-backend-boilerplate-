package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/zerologWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rajan2345/go-boilerplate/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// logging consists of two aspects -- one is monitoring and other is observability
// for observability we are using the service called NewRelic
// NewRelic offers different tools for monitoring and observability so that we can properly instrument different parts of the application
// e.g. Parts  -- Database , Requests etc..

// first thing create a NewRelic service, here we are storing an instance of NewRelic Application in a struct called LoggerService
type LoggerService struct {
	nrApp *newrelic.Application
}

// initialization of the instance that is logger service with newrelic integration
func NewLoggerService(cfg *config.ObservabilityConfig) *LoggerService {
	service := &LoggerService{}

	if cfg.NewRelic.LicenseKey == "" {
		fmt.Println("New Relic License key not provided , skipping initialization")
		return service
	}

	// for initialization of newrelic service , parameters are needed to be passed
	// for this a new slice is being created and parameters are being appended to it
	var ConfigOption []newrelic.ConfigOption
	ConfigOption = append(ConfigOption,
		newrelic.ConfigAppName(cfg.ServiceName),
		newrelic.ConfigLicense(cfg.NewRelic.LicenseKey),
		newrelic.ConfigAppLogForwardingEnabled(cfg.NewRelic.AppLogForwardingEnabled),
		newrelic.ConfigDistributedTracerEnabled(cfg.NewRelic.DistributedTracingEnabled),
	)

	//add debug logic only if explicitly enabled
	if cfg.NewRelic.DebugLogging {
		ConfigOption = append(ConfigOption, newrelic.ConfigDebugLogger(os.Stdout))
	}

	// now let's initialize the newrelic , as we have all the parameter needed for initialization
	app, err := newrelic.NewApplication(ConfigOption...)
	if err != nil {
		fmt.Printf("Failed to initialize new relic application: %v\n", err)
		return service
	}
	service.nrApp = app
	fmt.Printf("New Relic app initialized for app: %v\n", cfg.ServiceName)
	return service
}

// for graceful shutdown we will be creating a reciever function
func (ls *LoggerService) Shutdown() {
	if ls.nrApp != nil {
		ls.nrApp.Shutdown(10 * time.Second)
	}
}

// now for accessing the nrApp instance we will be creating a getter function
func (ls *LoggerService) GetApplication() *newrelic.Application {
	return ls.nrApp
}

// creating a new logger
func NewLogger(level string, isProd bool) zerolog.Logger {
	return NewLoggerWithService(
		&config.ObservabilityConfig{
			Logging: config.LoggingConfig{
				Level: level,
			},
			Environment: func() string {
				if isProd {
					return "production"
				} else {
					return "development"
				}
			}(),
		}, nil)
}

// creates a logger with full config
func NewLoggerWithConfig(cfg *config.ObservabilityConfig) zerolog.Logger {
	return NewLoggerWithService(cfg, nil)
}

// NewLoggerWithService creates a zerolog logger with NewRelic integration
func NewLoggerWithService(cfg *config.ObservabilityConfig, loggerService *LoggerService) zerolog.Logger {
	var logLevel zerolog.Level
	level := cfg.GetLoggingLevel()

	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	// don't set global level let each logger have its own level
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var writer io.Writer

	//setup base writer
	var baseWriter io.Writer

	if cfg.IsProduction() && cfg.Logging.Format == "json" {

		//in production write to std out
		baseWriter = os.Stdout

		// wrap with newrelic zerologwriter for log forwarding in production
		if loggerService != nil && loggerService.nrApp != nil {
			nWriter := zerologWriter.New(baseWriter, loggerService.nrApp)
			writer = nWriter
		} else {
			writer = baseWriter
		}
	} else {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05"}
		writer = consoleWriter
	}

	//from here on newRelic log forwarding is handled by the zerologWriter configuration

	logger := zerolog.New(writer).
		Level(logLevel).
		With().
		Timestamp().
		Str("service", cfg.ServiceName).
		Str("environment", cfg.Environment).
		Logger()

	// if development include stack traces in error logs
	if !cfg.IsProduction() {
		logger = logger.With().Stack().Logger()
	}

	return logger

}

// withtraceContext add newrelic transaction context to the logger
func WithTraceContext(logger zerolog.Logger, txn *newrelic.Transaction) zerolog.Logger {
	if txn == nil {
		return logger
	}

	//get trace metadata from transaction
	traceMetadata := txn.GetTraceMetadata()

	return logger.With().
		Str("trace_id", traceMetadata.TraceID).
		Str("span_id", traceMetadata.SpanID).
		Logger()
}

func FormatSQLWithArgs(sql string, args []any) string {
	result := sql
	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i+1)
		value := fmt.Sprintf("'%v'", arg)
		result = strings.Replace(result, placeholder, value, 1)
	}
	return result
}

// a database logger NewPgxLogger
func NewPgxLogger(level zerolog.Level) zerolog.Logger {
	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
		FormatFieldValue: func(i any) string {
			switch v := i.(type) {
			case string:
				if len(v) > 200 {
					return v[:200] + "..." // Truncate long SQL queries
				}
				return v
			case []byte:
				var obj interface{}
				if err := json.Unmarshal(v, &obj); err == nil {
					pretty, _ := json.MarshalIndent(obj, ",", "  ")
					return "\n" + string(pretty)
				}
				return string(v)
			default:
				return fmt.Sprintf("%v", v)
			}
		},
	}
	return zerolog.New(writer).
		Level(level).
		With().
		Timestamp().
		Str("component", "database").
		Logger()
}

// getting log level
func GetPgxTraceLogLevel(level zerolog.Level) int {
	switch level {
	case zerolog.DebugLevel:
		return 6
	case zerolog.InfoLevel:
		return 5
	case zerolog.WarnLevel:
		return 4
	case zerolog.ErrorLevel:
		return 3
	case zerolog.FatalLevel:
		return 2
	case zerolog.PanicLevel:
		return 1
	default:
		return 0
	}
}
