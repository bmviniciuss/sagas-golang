package validator

import (
	"reflect"
	"strings"

	goval "github.com/go-playground/validator/v10"
)

func New() *goval.Validate {
	validate := goval.New(goval.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	return validate
}
