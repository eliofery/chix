package validation

import (
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/go-playground/validator/v10"
)

// TestValidate пример реализации пользовательской валидации
func TestValidate() chix.CustomValidate {
	messages := make(map[string]string)
	messages["ru"] = "{0} не соответствует требованиям"
	messages["en"] = "{0} does not meet requirements"
	messages["fr"] = "{0} ne répond pas aux exigences"

	return chix.CustomValidate{
		Tag: "example",
		Func: func(fl validator.FieldLevel) bool {
			return fl.Field().Int() >= 12 && fl.Field().Int() <= 18
		},
		Message: messages,
	}
}
