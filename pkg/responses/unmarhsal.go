package responses

import (
	"encoding/json"
	"errors"
	"io"
)

func ParseErrorToResponse(reqID string, err error) Error {
	if err == nil {
		return NewInternalServerErrorResponse(reqID)
	}

	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		return NewBadRequestErrorResponse(reqID, []FieldError{
			{
				Field:   "body",
				Message: "Invalid JSON",
			},
		})
	}
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return NewBadRequestErrorResponse(reqID, []FieldError{
			{
				Field:   "body",
				Message: "Invalid JSON",
			},
		})
	}
	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return NewBadRequestErrorResponse(reqID, []FieldError{
			{
				Field:   unmarshalTypeError.Field,
				Message: "Invalid value",
			},
		})
	}

	if errors.Is(err, io.EOF) {
		return NewBadRequestErrorResponse(reqID, []FieldError{
			{
				Field:   "body",
				Message: "Empty body",
			},
		})
	}

	return NewInternalServerErrorResponse(reqID)
}
