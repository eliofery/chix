package chix

import (
	"context"
	"net/http"
)

const nextKey key = "next"

func WithNextHandler(ctx context.Context, next http.Handler) context.Context {
	return context.WithValue(ctx, nextKey, next)
}

func NextHandler(ctx context.Context) http.Handler {
	val := ctx.Value(nextKey)

	next, ok := val.(http.Handler)
	if !ok {
		return nil
	}

	return next
}
