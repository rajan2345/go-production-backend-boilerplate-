package errs

import "net/http"

// this is 404 error
func NewUnauthorizeError(message string, override bool) *HttpError {
	return &HttpError{
		Code:     MakeUpperCaseWithUnderscores(http.StatusText(http.StatusUnauthorized)),
		Message:  message,
		Status:   http.StatusUnauthorized,
		Override: override,
	}
}

// this is 403 error , forbidden error
func NewForbiddenError(message string, override bool) *HttpError {
	return &HttpError{
		Code:     MakeUpperCaseWithUnderscores(http.StatusText(http.StatusForbidden)),
		Message:  message,
		Status:   http.StatusForbidden,
		Override: override,
	}
}

// this is badrequest 400 error
func NewBadRequestError(message string, override bool, code *string, errors []FieldError, action *Action) *HttpError {
	formattedCode := MakeUpperCaseWithUnderscores(http.StatusText(http.StatusBadRequest))

	if code != nil {
		formattedCode = *code
	}
	return &HttpError{
		Code:     formattedCode,
		Message:  message,
		Status:   http.StatusBadRequest,
		Override: override,
		Errors:   errors,
		Action:   action,
	}
}
