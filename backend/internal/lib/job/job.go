// this file will contain service provided and managed by the library asynq,
package job

import (
	"github.com/hibiken/asynq"
	"github.com/rajan2345/go-boilerplate/internal/config"
	"github.com/rs/zerolog"
)

// this is just for sending email
// define a struct which will contain the instance of the library
type JobService struct {
	Client *asynq.Client
	Server *asynq.Server
	logger *zerolog.Logger
}

// next we will create a function which will actually create the instance of the library and will return the instance of the
// jobservice
// Dependency Injection : - will be using a functionality as a dependeny , not like object defined in a place , used everywhere
// client is the the thing through which we will insert the jobs , server is the module which is going to execute the jobs .
func NewJobService(logger *zerolog.Logger, cfg *config.Config) *JobService {
	redisAddr := cfg.Redis.Address

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6, // higher priority for important mails
				"default":  3, // lower that higher priority , use cases less
				"low":      1, // lowest of priority
			},
		},
	)

	return &JobService{
		Client: client,
		Server: server,
		logger: logger,
	}
}

// graceful shutdown -- one of the pillars of the production grade code, and starting of the server .
func (j *JobService) Start() error {
	// Register task handler
	// creation of server -> registration of handler-> create the handler or handler exists -> logging

	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskWelcome, j.handleWelcomeEmailTask)

	j.logger.Info().Msg("Starting backgroud job server")
	if err := j.Server.Start(mux); err != nil {
		return err
	}

	return nil
}

func (j *JobService) Stop() {
	j.logger.Info().Msg("Backgroud server is stopping")
	j.Server.Shutdown()
	j.Client.Close()
}

// creation of handler and task types
// strigify json -- in go it is json types
