package chix

import (
	"errors"
	"github.com/eliofery/go-chix/pkg/log"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ru_translations "github.com/go-playground/validator/v10/translations/ru"
	"log/slog"
	"reflect"
	"strings"
)

const (
	langDefault    = "ru"
	defaultTagName = "name"
)

var (
	defaultLanguage = ru.New()
)

// CustomValidate пользовательская валидация
type CustomValidate struct {
	Tag     string
	Func    validator.Func
	Message map[string]string
}

// DefaultTranslations функции перевода
type DefaultTranslations map[string]func(v *validator.Validate, trans ut.Translator) error

// Validate валидация данных
type Validate interface {
	Validation(data any, langOption ...string) []error
	UseTagName(name string) Validate
	UseLocales(translators ...locales.Translator) Validate
	UseValidations(customValidates ...CustomValidate) Validate
	UseTranslations(dt DefaultTranslations) Validate
}
type validate struct {
	*validator.Validate
	Locales      []locales.Translator
	Validations  []CustomValidate
	Translations map[string]func(v *validator.Validate, trans ut.Translator) error

	tagName string // tagName имя тега
}

// NewValidate создание валидации
func NewValidate(v *validator.Validate) Validate {
	log.Debug("Инициализация валидации")

	var vd validate
	vd.Validate = v

	vd.UseTagName(defaultTagName)
	vd.UseLocales(defaultLanguage)
	vd.UseTranslations(DefaultTranslations{
		"ru": ru_translations.RegisterDefaultTranslations,
	})

	return &vd
}

// Validation валидация входных данных
// langOptions: https://github.com/go-playground/validator/tree/master/translations
func (v *validate) Validation(data any, langOption ...string) []error {
	langString := langDefault
	if len(langOption) > 0 {
		langString = langOption[0]
	}

	trans, lang := v.registerTranslate(langString)

	v.registerCustomValidation(trans, lang)

	// Данный код переводит сообщение об ошибке в соответствии с текущим языком
	// Сам по себе валидатор переводит только сообщение об ошибке без самих полей
	// Пример: Password обязательное поле
	// При связанных полях перевод выглядит еще не корректней
	// Пример: Password должен быть равен PasswordConfirm
	// Данную проблему можно было решить переопределив стандартные валидаторы пользовательскими
	// Но это сложно, долго и не удобно. Данная реализация автоматизирует процесс перевода
	// Так же ни кто не мешает переопределять стандартные валидаторы если это необходимо
	// Все по прежнему будет работать
	var validatorErr validator.ValidationErrors
	if err := v.Validate.Struct(data); err != nil && errors.As(err, &validatorErr) {
		var errMessages []error
		for _, validateErr := range validatorErr {
			translateError := validateErr.Translate(trans)
			field, message := v.parseErrorMessage(translateError, lang)

			associateField := v.equivalentField(data)
			param := validateErr.Param()
			transcription := associateField[param]
			eqfield, _ := v.parseErrorMessage(transcription, lang)
			if eqfield == "" {
				eqfield = transcription
			}
			messageReplace := strings.Replace(message, param, eqfield, 1)

			errMessage := errors.New(field + "" + messageReplace)
			errMessages = append(errMessages, errMessage)
		}

		return errMessages
	}

	return nil
}

// UseTagName регистрация имени поля в структуре данных
// Пример:
//
//	type User struct {
//	   Name string `json:"name" validate:"required" <имя>:"Имя"`
//	   Age  int    `json:"age" validate:"required" <имя>:"[ru:Возраст;en:Age;fr:Ajy]"`
//	}
func (v *validate) UseTagName(name string) Validate {
	v.Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		tag := field.Tag.Get(name)
		if tag == "-" {
			return ""
		}

		v.tagName = name

		return tag
	})

	return v
}

// UseLocales использование языков
func (v *validate) UseLocales(locales ...locales.Translator) Validate {
	v.Locales = locales

	return v
}

// UseValidations использование пользовательских валидаций
func (v *validate) UseValidations(customValidates ...CustomValidate) Validate {
	v.Validations = customValidates

	return v
}

// UseTranslations использование переводов
func (v *validate) UseTranslations(dt DefaultTranslations) Validate {
	v.Translations = dt

	return v
}

// registerTranslate регистрация перевода
func (v *validate) registerTranslate(lang string) (ut.Translator, string) {
	uni := ut.New(v.Locales[0], v.Locales...)

	trans, ok := uni.GetTranslator(lang)
	if !ok {
		log.Warn("Язык не поддерживается", slog.String("lang", lang))
		lang = langDefault
		trans, _ = uni.GetTranslator(lang)
	}

	if err := v.Translations[lang](v.Validate, trans); err != nil {
		log.Warn("Не удалось зарегистрировать перевод", slog.String("err", err.Error()))
	}

	return trans, lang
}

// registerCustomValidation регистрация пользовательской валидации
func (v *validate) registerCustomValidation(trans ut.Translator, lang string) {
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

			// перевод для значения имени поля зарегистрированное через UseTagName
			// поля должны иметь формат: "[ru:Возраст;en:Age;fr:Ajy]" или "Возраст"
			func(ut ut.Translator, fe validator.FieldError) string {
				translations, _ := v.parseErrorMessage(fe.Field(), lang)

				t, _ := ut.T(cv.Tag, translations)
				return t
			},
		)
		if err != nil {
			log.Warn("Не удалось перевести пользовательскую валидацию", slog.String("err", err.Error()))
		}
	}
}

// parseErrorMessage разбирает текст ошибки на поле и сообщение
// [ru:Пароль;en:Password;fr:Pass] обязательное поле -> "Пароль", "обязательное поле"
func (v *validate) parseErrorMessage(errString, lang string) (string, string) {
	separator := "]"
	closeBracketIndex := strings.Index(errString, separator)

	var field, message string
	if closeBracketIndex != -1 {
		field = v.parseTranslations(errString[1:closeBracketIndex], lang)
		message = errString[closeBracketIndex+1:]

		return field, message
	}

	return "", errString
}

// parseTranslations разбивает строку в поисках языка
// ru:Возраст;en:Age;fr:Ajy -> Возраст или Age или Ajy
func (v *validate) parseTranslations(translationString, lang string) string {
	fields := make(map[string]string)

	translations := strings.Split(translationString, ";")
	if len(translationString) == 0 {
		return translationString
	}

	for _, translation := range translations {
		parts := strings.Split(translation, ":")

		if len(parts) == 2 {
			name := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			fields[name] = value
		} else {
			fields[lang] = translationString
		}
	}

	return fields[lang]
}

// equivalentField ищет в структуре поле на которое ссылается другое поле
// Пример:
// Password        string `json:"password" name:"[ru:Пароль;en:Password;fr:Mot de passe]"`
// PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password" name:"[ru:Подтверждение пароля;en:Password confirmation;fr:Confirmation mot de passe]"`
//
//	Сформирует: map["Password":"[ru:Пароль;en:Password;fr:Mot de passe]"]
func (v *validate) equivalentField(obj any) map[string]string {
	objValue := reflect.ValueOf(obj)
	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}

	result := make(map[string]string)
	objType := objValue.Type()

	for i := 0; i < objValue.NumField(); i++ {
		field := objType.Field(i).Name

		value := objType.Field(i).Tag.Get(v.tagName)
		if value == "" {
			value = objType.Field(i).Name
		}

		result[field] = value
	}

	return result
}
