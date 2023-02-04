package rest

import (
	"encoding/json"
	"errors"
	"net/http"
)

// restError represents an error
// which encodes the HTTP status code too
type restError struct {
	Code int        `json:"code"`
	Err  errorBlock `json:"error"`
}

type errorBlock struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func newError(status int, message string) *restError {
	return &restError{
		Code: status,
		Err: errorBlock{
			Status:  http.StatusText(status),
			Message: message,
		},
	}
}

var _ error = (*restError)(nil)

func writeError(w http.ResponseWriter, err error) {
	e := asHTTPError(err)

	bytes, _ := json.MarshalIndent(e, "", "  ")

	w.WriteHeader(e.Code)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = w.Write(bytes)
	_, _ = w.Write([]byte("\n"))
}

func (e *restError) Error() string {
	return e.Err.Message
}

func asHTTPError(err error) *restError {
	if err == nil {
		return nil
	}

	var httpErr *restError
	if errors.As(err, &httpErr) {
		return httpErr
	}

	return newError(http.StatusInternalServerError, err.Error())
}

func NotFound(message string) error {
	return newError(http.StatusNotFound, message)
}

func MethodNotAllowed(message string) error {
	return newError(http.StatusMethodNotAllowed, message)
}

func BadRequest(message string) error {
	return newError(http.StatusBadRequest, message)
}

func Unauthorized(message string) error {
	return newError(http.StatusUnauthorized, message)
}

func Forbidden(message string) error {
	return newError(http.StatusForbidden, message)
}

func InternalServerError(message string) error {
	return newError(http.StatusInternalServerError, message)
}
