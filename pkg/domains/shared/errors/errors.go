package errors

import "fmt"

type ConflictError struct {
	Message string
}

func (e ConflictError) Error() string {
	return e.Message
}

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

var (
	ErrProductAlreadyExists = ConflictError{
		Message: "product with the same name and category already exists",
	}
)

func NewValidationError(msg string, args ...any) ValidationError {
	return ValidationError{
		Message: fmt.Sprintf(msg, args...),
	}
}
