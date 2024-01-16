package middleware

import (
	"errors"
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/eliofery/go-chix/pkg/log"
	"net/http"
)

func Example2() chix.HandlerCtx {
	log.Debug("Инициализация middleware Example2")

	return func(ctx *chix.Ctx) error {
		ctx.Status(http.StatusForbidden)
		return errors.New("доступ запрещен")
	}
}
