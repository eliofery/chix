package chix

import (
	"context"
	"net/http"
)

const requestKey key = "request"

// WithRequest добавляет запрос в контекст
func WithRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, requestKey, r)
}

// Request получает запрос из контекста
func Request(ctx context.Context) *http.Request {
	val := ctx.Value(requestKey)

	request, ok := val.(*http.Request)
	if !ok {
		return nil
	}

	return request
}
