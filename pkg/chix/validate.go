package chix

import (
	"errors"
	"github.com/eliofery/go-chix/pkg/log"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"reflect"
	"strings"
)

const (
	langDefault    = "ru"
	defaultTagName = "name"
)

// CustomValidate пользовательская валидация
type CustomValidate struct {
	Tag     string
	Func    validator.Func
	Message map[string]string
}

// LanguageProcessor функции перевода
type LanguageProcessor map[string]func(v *validator.Validate, trans ut.Translator) error

// Validate валидация данных
type Validate interface {
	Validation(data any, langOptions ...string) []error             // Валидация входных данных
	RegisterTagName(name ...string) Validate                        // Регистрация имени поля в структуре данных
	RegisterLanguages(translators ...locales.Translator) Validate   // Регистрация переводов
	RegisterValidations(customValidates ...CustomValidate) Validate // Регистрация пользовательских валидации
	RegisterLanguagesProcess(lp LanguageProcessor) Validate         // Регистрация функций перевода
}

type validate struct {
	*validator.Validate                                                                   // Validate Валидатор
	Translators         []locales.Translator                                              // Translators переводы
	Validations         []CustomValidate                                                  // Validations пользовательские валидации
	LanguageProcessors  map[string]func(v *validator.Validate, trans ut.Translator) error // LanguageProcessors функции перевода
}

// NewValidate создание валидации
func NewValidate(v *validator.Validate) Validate {
	log.Debug("Инициализация валидации")

	return &validate{
		Validate: v,
	}
}

// Validation валидация входных данных
// langOptions: https://github.com/go-playground/validator/tree/master/translations
func (v *validate) Validation(data any, langOptions ...string) []error {
	trans, lang := v.registerTranslate(langOptions...)

	v.registerAndTranslateCustomValidation(trans, lang)

	var validatorErr validator.ValidationErrors
	if err := v.Validate.Struct(data); err != nil && errors.As(err, &validatorErr) {
		var errMessages []error
		for _, validateErr := range validatorErr {
			errMessage := errors.New(validateErr.Translate(trans))
			errMessages = append(errMessages, errMessage)
		}

		return errMessages
	}

	return nil
}

// RegisterTagName регистрация имени поля в структуре данных
// Пример:
//
//	type User struct {
//	   Name string `json:"name" validate:"required" <имя>:"Имя"`
//	   Age  int    `json:"age" validate:"required" <имя>:"ru:Возраст;en:Age;fr:Ajy"`
//	}
func (v *validate) RegisterTagName(name ...string) Validate {
	if len(name) == 0 {
		name[0] = defaultTagName
	}

	v.Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		tag := field.Tag.Get(name[0])
		if tag == "-" {
			return ""
		}
		return tag
	})

	return v
}

// RegisterLanguages регистрация пользовательских переводов
func (v *validate) RegisterLanguages(translators ...locales.Translator) Validate {
	v.Translators = translators

	return v
}

// RegisterValidations регистрация пользовательской валидации
func (v *validate) RegisterValidations(customValidates ...CustomValidate) Validate {
	v.Validations = customValidates

	return v
}

func (v *validate) RegisterLanguagesProcess(lp LanguageProcessor) Validate {
	v.LanguageProcessors = lp

	return v
}

// registerTranslate регистрация перевода
func (v *validate) registerTranslate(langOptions ...string) (ut.Translator, string) {
	lang := langDefault
	if len(langOptions) > 0 {
		lang = langOptions[0]
	}

	if len(v.Translators) == 0 {
		v.Translators = append(v.Translators, ru.New())
	}

	uni := ut.New(v.Translators[0], v.Translators...)

	trans, ok := uni.GetTranslator(lang)
	if !ok {
		log.Warn("Язык не поддерживается", slog.String("lang", lang))
		lang = langDefault
		trans, _ = uni.GetTranslator(lang)
	}

	if err := v.LanguageProcessors[lang](v.Validate, trans); err != nil {
		log.Warn("Не удалось зарегистрировать перевод", slog.String("err", err.Error()))
	}

	return trans, lang
}

// registerAndTranslateCustomValidation регистрация и перевод пользовательской валидации
func (v *validate) registerAndTranslateCustomValidation(trans ut.Translator, lang string) {
	for _, cv := range v.Validations {
		if err := v.Validate.RegisterValidation(cv.Tag, cv.Func); err != nil || cv.Message[lang] == "" {
			log.Warn("Не удалось зарегистрировать пользовательскую валидацию", slog.String("err", err.Error()))

			return
		}

		err := v.Validate.RegisterTranslation(
			cv.Tag,
			trans,

			// перевод для сообщения
			func(ut ut.Translator) error {
				return ut.Add(cv.Tag, cv.Message[lang], true)
			},

			// перевод для значения имени поля зарегистрированное через RegisterTagName
			// поля должны иметь формат: "ru:Возраст;en:Age;fr:Ajy" или "Возраст"
			func(ut ut.Translator, fe validator.FieldError) string {
				fields := make(map[string]string)

				langs := strings.Split(fe.Field(), ";")
				for _, l := range langs {
					parts := strings.Split(l, ":")

					if len(parts) == 2 {
						name := strings.TrimSpace(parts[0])
						value := strings.TrimSpace(parts[1])

						fields[name] = value
					} else {
						fields[lang] = fe.Field()
					}
				}

				t, _ := ut.T(cv.Tag, fields[lang])
				return t
			},
		)
		if err != nil {
			log.Warn("Не удалось зарегистрировать пользовательскую валидацию", slog.String("err", err.Error()))
		}
	}
}
