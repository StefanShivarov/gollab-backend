package common

import "net/http"

type ApiError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (e *ApiError) Error() string {
	return e.Message
}

func NotFound(msg string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		Message:    msg,
	}
}

func BadRequest(msg string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Message:    msg,
	}
}

func InternalServerError(msg string) *ApiError {
	return &ApiError{
		StatusCode: http.StatusInternalServerError,
		Message:    msg,
	}
}

type ErrorResponse struct {
	Error ApiError `json:"error"`
}
