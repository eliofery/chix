package middleware

import (
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/eliofery/go-chix/pkg/log"
)

func Example() chix.Handler {
	log.Debug("Инициализация middleware Example")

	return func(ctx *chix.Ctx) error {
		return ctx.Next()
	}
}
