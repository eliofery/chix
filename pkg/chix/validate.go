package chix

import "github.com/go-playground/validator/v10"

// Validate валидация данных
type Validate interface {
	Validation(data any, langOptions ...string) []error    // Валидация входных данных
	RegisterValidations(customValidates ...CustomValidate) // Регистрация пользовательских валидации
}

// CustomValidate пользовательская валидация
type CustomValidate struct {
	Name string
	Func validator.Func
}
