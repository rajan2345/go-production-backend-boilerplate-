// this file will contain service provided and managed by the library asynq,
package job

import (
	"github.com/hibiken/asynq"
	zerolog "github.com/jackc/pgx-zerolog"
)

// this is just for sending email
// define a struct which will contain the instance of the library
type JobService struct {
	Client *asynq.Client
	server *asynq.Server
	logger *zerolog.Logger
}

// next we will create a function which will actually create the instance of the library and will return the instance of the
// jobservice
