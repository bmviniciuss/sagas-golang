package responses

import "github.com/go-playground/validator/v10"

func ValidatorErrorToFieldError(err error) []FieldError {
	var fieldErrors []FieldError
	for _, err := range err.(validator.ValidationErrors) {
		fieldErrors = append(fieldErrors, FieldError{
			Field:   err.Field(),
			Message: err.Error(),
		})
	}
	return fieldErrors
}
