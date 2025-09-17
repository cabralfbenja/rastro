package app_errors

type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) *NotFoundError {
	if message == "" {
		message = "Resource not found"
	}
	return &NotFoundError{Message: message}
}
