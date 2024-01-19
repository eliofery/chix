package middleware

import (
	"github.com/eliofery/go-chix/pkg/chix"
)

// Example пример реализации middleware
func Example() chix.Handler {
	return func(ctx *chix.Ctx) error {
		return ctx.Next()
	}
}
