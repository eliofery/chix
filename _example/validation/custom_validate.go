package validation

import (
    en_translations "github.com/go-playground/validator/v10/translations/en"
    fr_translations "github.com/go-playground/validator/v10/translations/fr"
    ru_translations "github.com/go-playground/validator/v10/translations/ru"
)

// TestValidate пример реализации пользовательской валидации.
func TestValidate() chix.CustomValidate {
	// Переводы сообщения иб ошибки.
	// На место {0} подставится название поля:
    // Возраст не соответствует требования
    // Age does not meet requirements
    // Âge ne répond pas aux exigences
	messages := make(map[string]string)
	messages["ru"] = "{0} не соответствует требованиям"
	messages["en"] = "{0} does not meet requirements"
	messages["fr"] = "{0} ne répond pas aux exigences"

	return chix.CustomValidate{
		Tag: "age_foo_bar",
		Func: func(fl validator.FieldLevel) bool {
			return fl.Field().Int() >= 12 && fl.Field().Int() <= 18
		},
		Message: messages,
	}
}

// Пример использования в структуре
type User struct {
	Age int `json:"age" validate:"age_foo_bar" label:"ru:Возраст;en:Age;fr:Âge"`
}

// Пример регистрации пользовательской валидации
func main() {
    valid := chix.NewValidate(validator.New()).
        RegisterLocales(
            ru.New(),
            en.New(),
            fr.New(),
        ).
        RegisterTranslations(chix.DefaultTranslations{
            "ru": ru_translations.RegisterDefaultTranslations,
            "en": en_translations.RegisterDefaultTranslations,
            "fr": fr_translations.RegisterDefaultTranslations,
        }).
        RegisterValidations(
            TestValidate(),
        )
}
