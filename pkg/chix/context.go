package chix

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/eliofery/go-chix/pkg/log"
	"io"
	"log/slog"
	"net/http"
)

// Map шаблон для передачи данных
type Map map[string]any

// Ctx контекст предоставляемый в обработчик
type Ctx struct {
	context.Context
	Validate
	http.ResponseWriter
	*http.Request

	status int
}

// NewCtx создание контекста
func NewCtx(ctx context.Context, validate Validate) *Ctx {
	return &Ctx{
		Context:        ctx,
		Validate:       validate,
		ResponseWriter: ResponseWriter(ctx),
		Request:        Request(ctx),

		status: http.StatusOK,
	}
}

// Status установка статуса ответа
func (ctx *Ctx) Status(status int) *Ctx {
	ctx.status = status
	return ctx
}

// ContentType установка типа контента
func (ctx *Ctx) ContentType(ct string) *Ctx {
	ctx.ResponseWriter.Header().Set("Content-Type", ct)
	return ctx
}

// Decode декодирование тела запроса
func (ctx *Ctx) Decode(data any, langOption ...string) error {
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
