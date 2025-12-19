package errs

import "net/http"

type Status struct {
	status  int
	message string
}

func (e *Status) Error() string {
	return e.message
}

func (e *Status) Status() int {
	return e.status
}

func NewStatusError(message string, status int) *Status {
	return &Status{
		message: message,
		status: status,
	}
}

var StatusInternalServerError = NewStatusError(
	"Internal Server Error",
	http.StatusInternalServerError,
)

var StatusNotFound = NewStatusError(
	"Requested Resource Wasn't Found",
	http.StatusNotFound,
)

var StatusTimeout = NewStatusError(
	"Operation timeout",
	http.StatusRequestTimeout,
)
