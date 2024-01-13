package validate

import (
	"errors"
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/eliofery/go-chix/pkg/log"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ru_translations "github.com/go-playground/validator/v10/translations/ru"
	"log/slog"
)

const (
	langDefault = "ru"
)

type Validate struct {
	*validator.Validate                      // Validate Валидатор
	Translators         []locales.Translator // Translators переводы
}

// New создание валидации
func New(v *validator.Validate, translators ...locales.Translator) chix.Validate {
	log.Debug("Инициализация валидации")

	return &Validate{
		Validate:    v,
		Translators: translators,
	}
}

// Validation валидация входных данных
// langOptions: https://github.com/go-playground/validator/tree/master/translations
func (v *Validate) Validation(data any, langOptions ...string) []error {
	var (
		validatorErr validator.ValidationErrors
		errMessages  []error
	)

	langString := langDefault
	if len(langOptions) > 0 {
		langString = langOptions[0]
	}

	lang := v.setLang(langString, v.Translators...)
	if err := v.Validate.Struct(data); err != nil && errors.As(err, &validatorErr) {
		for _, validateErr := range validatorErr {
			errMessage := errors.New(validateErr.Translate(lang))
			errMessages = append(errMessages, errMessage)
		}
	}

	return errMessages
}

// RegisterValidations регистрация пользовательской валидации
func (v *Validate) RegisterValidations(customValidates ...chix.CustomValidate) {
	for _, cv := range customValidates {
		err := v.Validate.RegisterValidation(cv.Name, cv.Func)
		if err != nil {
			log.Warn("Не удалось зарегистрировать валидацию", slog.String("err", err.Error()))
		}
	}
}

// setLang перевод ошибок валидации
func (v *Validate) setLang(lang string, translators ...locales.Translator) ut.Translator {
	if len(translators) == 0 {
		translators = append(translators, ru.New())
	}

	uni := ut.New(translators[0], translators...)

	trans, ok := uni.GetTranslator(lang)
	if !ok {
		log.Warn("Язык не поддерживается", slog.String("lang", lang))
		trans, _ = uni.GetTranslator(langDefault)
	}

	if err := ru_translations.RegisterDefaultTranslations(v.Validate, trans); err != nil {
		log.Warn("Не удалось зарегистрировать перевод", slog.String("err", err.Error()))
	}

	return trans
}
