package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("excludespace", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		return !strings.Contains(value, " ")
	})
}
