// need to create a custom error type , solely for sending to the client (client e.g. postman or frontend)
package errs

import "strings"

// for form based error , like tag verification can be done through field , error through slice late(Error []ErrorField , slice of errors ) multiple type of error can be
// displayed
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// next is action type , we use it when the user session in expired
// Actions -- requests or operations that fail to be processed correctly by the server-side application, resulting in an error condition
type ActionType string

const (
	ActionTypeRedirect ActionType = "redirect"
)

type Action struct {
	Type    ActionType `json:"type"`
	Message string     `json:"message"`
	Value   string     `json:"value"`
}

type HttpError struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Status   int    `json:"status"`
	Override bool   `json:"override"` // It's about replacing generic error behavior with tailored, application-specific error management.

	// field level error
	Errors []FieldError `json:"errors"`
	Action *Action      `json:"action"`
}

// reciever method which will be called
func (e *HttpError) Error() string {
	return e.Message
}

func (e *HttpError) Is(target error) bool {
	_, ok := target.(*HttpError)
	return ok
}

func (e *HttpError) WithMessage(message string) *HttpError {
	return &HttpError{
		Code:     e.Code,
		Message:  message,
		Status:   e.Status,
		Override: e.Override,
		Errors:   e.Errors,
		Action:   e.Action,
	}
}

func MakeUpperCaseWithUnderscores(str string) string {
	return strings.ToUpper(strings.ReplaceAll(str, " ", "_"))
}
