package job

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TaskWelcome = "email:welcome"
)

type WelcomeEmailPayload struct {
	To        string `json:"to"`
	FirstName string `json:"firstname"`
}

func NewWelcomeEmailTask(to, firstName string) (*asynq.Task, error) {
	payload, err := json.Marshal(WelcomeEmailPayload{
		To:        to,
		FirstName: firstName,
	})

	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskWelcome, payload,
		asynq.MaxRetry(3),
		asynq.Queue("default"),
		asynq.Timeout(30*time.Second)), nil
}

// TaksWelcome --

//    1. Definition in internal/lib/job/email_task.go:
//       The constant is defined with a specific string name.
//         L11: TaskWelcome = "email:welcome"

//    2. Creation in internal/lib/job/email_task.go:
//       Your code creates a task using this constant.
//         L29: return asynq.NewTask(TaskWelcome, payload, ...

//    3. Handling in internal/lib/job/job.go:
//       This is the most important part. The asynq task server (mux) is told to associate the TaskWelcome constant with a specific function, j.handleEmailWelcomeTask.
//         L54: mux.HandleFunc(TaskWelcome, j.handleEmailWelcomeTask)

//   So, when a worker process receives a task named "email:welcome", this line tells it to execute the j.handleEmailWelcomeTask function to perform the actual job. The constant
//   is the bridge between creating the task and executing it.
