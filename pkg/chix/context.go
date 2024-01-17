package chix

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/eliofery/go-chix/pkg/log"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

const IssuerKey key = "issuer"

// Map шаблон для передачи данных
type Map map[string]any

// Ctx контекст предоставляемый в обработчик
type Ctx struct {
	context.Context
	http.ResponseWriter
	*http.Request
	Validate
	NextHandler http.Handler

	status int
}

// NewCtx создание контекста
func NewCtx(
	w http.ResponseWriter,
	r *http.Request,
	validate Validate,
) *Ctx {
	return &Ctx{
		Context:        r.Context(),
		ResponseWriter: w,
		Request:        r,
		Validate:       validate,
		NextHandler:    NextHandler(r.Context()),

		status: http.StatusOK,
	}
}

// Status установка статуса ответа
func (ctx *Ctx) Status(status int) *Ctx {
	ctx.status = status
	return ctx
}

func (ctx *Ctx) Header(key, value string) {
	ctx.ResponseWriter.Header().Set(key, value)
}

// ContentType установка типа контента
func (ctx *Ctx) ContentType(ct string) *Ctx {
	ctx.Header("Content-Type", ct)
	return ctx
}

// Decode декодирование тела запроса
func (ctx *Ctx) Decode(data any, langOption ...string) error {
	if len(langOption) == 0 {
		langOption = append(langOption, ctx.getActualLang())
	}

	err := json.NewDecoder(ctx.Request.Body).Decode(data)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("пустое тело запроса")
		}

		log.Error("Не удалось сформировать объект", slog.String("err", err.Error()))
		return err
	}

	if ctx.Validate != nil {
		if errMessages := ctx.Validate.Validation(data, langOption...); errMessages != nil {
			ctx.Status(http.StatusBadRequest)
			return errors.Join(errMessages...)
		}
	}

	return nil
}

// JSON формирование json ответа
func (ctx *Ctx) JSON(data Map) error {
	ctx.ContentType("application/json")
	ctx.ResponseWriter.WriteHeader(ctx.status)

	encoder := json.NewEncoder(ctx.ResponseWriter)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		log.Error("Не удалось сформировать json", slog.String("err", err.Error()))
		return err
	}

	return nil
}

// Next обработка следующего обработчика
func (ctx *Ctx) Next() error {
	ctx.NextHandler.ServeHTTP(ctx.ResponseWriter, ctx.Request)

	return nil
}

// Locals добавление/получение данных в контексте
func (ctx *Ctx) Locals(key any, value ...any) any {
	if len(value) == 0 {
		return ctx.Request.Context().Value(key)
	}

	saveCtx := context.WithValue(ctx, key, value[0])
	ctx.Request = ctx.WithContext(saveCtx)

	return value[0]
}

// Get получить содержимое заголовка
func (ctx *Ctx) Get(key string, defaultValue ...string) string {
	header := ctx.Request.Header.Get(key)

	if len(header) == 0 {
		log.Debug("Заголовок не найден", slog.Any("header", key))

		if len(defaultValue) == 0 {
			return ""
		}

		return defaultValue[0]
	}

	return header
}

// langDetected получение языка
func (ctx *Ctx) getActualLang() string {
	langHeader := ctx.Get("Accept-Language")

	if len(langHeader) == 0 {
		return ""
	}

	for _, group := range strings.Split(langHeader, ";") {
		for _, lang := range strings.Split(group, ",") {
			if len(lang) == 2 {
				return lang
			}
		}
	}

	return ""
}

// GetUserIdFromToken получение идентификатора пользователя
func (ctx *Ctx) GetUserIdFromToken() int {
	return ctx.Locals(IssuerKey).(int)
}
