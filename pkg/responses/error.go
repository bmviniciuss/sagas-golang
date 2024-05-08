package responses

import "net/http"

type Error struct {
	Status int         `json:"-"`
	ID     string      `json:"id"`
	Err    ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewInternalServerErrorResponse(id string) Error {
	return Error{
		ID:     id,
		Status: http.StatusInternalServerError,
		Err: ErrorDetail{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		},
	}
}

func NewNotFoundErrorResponse(id string) Error {
	return Error{
		ID:     id,
		Status: http.StatusNotFound,
		Err: ErrorDetail{
			Code:    http.StatusNotFound,
			Message: "Not Found",
		},
	}
}

func NewBadRequestErrorResponse(id string, details []FieldError) Error {
	return Error{
		ID:     id,
		Status: http.StatusBadRequest,
		Err: ErrorDetail{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Details: details,
		},
	}
}
