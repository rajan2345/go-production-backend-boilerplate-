package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	nrredis "github.com/newrelic/go-agent/v3/integrations/nrredis-v9"
	"github.com/rajan2345/go-boilerplate/internal/config"
	"github.com/rajan2345/go-boilerplate/internal/database"
	"github.com/rajan2345/go-boilerplate/internal/lib/job"
	loggerPkg "github.com/rajan2345/go-boilerplate/internal/logger"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// this file will contain the core data structure which will contain our config , the logger , the new relic logger service , database instance , redis instance our
// http server and background job processing server
//This file will take all these instances of different server , different modules and services and put all of them inside a data structure . Now this data structure
// is called and passing this server instance or pointer to the server we can establish dependency injection workflow, i.e. we will pass pointer to any of the services
// where we will be needing these services
// for example if we want to take database instance , then we will pass the struct and then we will take out the database instance using that server instance

// let's create a struct first
type Server struct {
	Config        *config.Config
	Logger        *zerolog.Logger
	LoggerService *loggerPkg.LoggerService
	DB            *database.Database
	Redis         *redis.Client
	httpServer    *http.Server
	Job           *job.JobService
}

// implement the function which is going to initialize all the server instance like config , logger etc.. and put them into Server struct
func New(cfg *config.Config, logger *zerolog.Logger, loggerService *loggerPkg.LoggerService) (*Server, error) {
	db, err := database.New(cfg, logger, loggerService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Redis client with new relic integration
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Address,
	})

	// add new relic redis hooks
	if loggerService != nil && loggerService.GetApplication() != nil {
		redisClient.AddHook(nrredis.NewHook(redisClient.Options()))
	}

	// Test redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Error().Err(err).Msg("Failed to connect with Redis, Continuing without Redis")
		//Don't fail startup if redis is not working
	}

	// job service
	jobService := job.NewJobService(logger, cfg)
	jobService.InitHandlers(cfg, logger)

	if err := jobService.Start(); err != nil {
		return nil, err
	}

	server := &Server{
		Config:        cfg,
		Logger:        logger,
		LoggerService: loggerService,
		DB:            db,
		Redis:         redisClient,
		Job:           jobService,
	}
	return server, nil
}

// Now we have service , now what is needed is to function for clean start  and shutdown of this service
// first creation of httpServer and assigning it to the http.Server
func (s *Server) HttpServer(handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:         ":" + s.Config.Server.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(s.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.Config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.Config.Server.IdleTimeout) * time.Second,
	}
}

func (s *Server) Start() error {
	if s.httpServer == nil {
		return errors.New("HTTP server not initialized")
	}

	// start the logger service
	s.Logger.Info().
		Str("port", s.Config.Server.Port).
		Str("env", s.Config.Primary.Env).
		Msg("starting server")

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	if err := s.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	if s.Job != nil {
		s.Job.Stop()
	}
	return nil
}
