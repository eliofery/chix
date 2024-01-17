package chix

import (
	"context"
	"errors"
	"github.com/eliofery/go-chix/pkg/log"
	"github.com/eliofery/go-chix/pkg/utils"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Handler обработчик, контроллер
type Handler func(ctx *Ctx) error

// Router обертка над chi роутером
type Router struct {
	*chi.Mux
	Validate Validate

	statistic map[string]int
}

// NewRouter создание роутера
func NewRouter(validate Validate) *Router {
	return &Router{
		Mux:      chi.NewRouter(),
		Validate: validate,

		statistic: make(map[string]int),
	}
}

// handleCtx запускает обработчик роутера
func (rt *Router) handler(handler Handler, w http.ResponseWriter, r *http.Request) {
	ctx := NewCtx(w, r, rt.Validate)

	if err := handler(ctx); err != nil {
		err = ctx.JSON(Map{
			"success": false,
			"message": utils.FirstToUpper(err.Error()),
		})
		if err != nil {
			log.Error("Не удалось обработать запрос", slog.String("handler", err.Error()))
			http.Error(ctx.ResponseWriter, "Не предвиденная ошибка", http.StatusInternalServerError)
		}
	}
}

// Get запрос на получение данных
func (rt *Router) Get(path string, handler Handler) {
	rt.statistic["routes"]++
	rt.Mux.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handler(handler, w, r)
	})
}

// Post запрос на сохранение данных
func (rt *Router) Post(path string, handler Handler) {
	rt.statistic["routes"]++
	rt.Mux.Post(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handler(handler, w, r)
	})
}

// Put запрос на обновление всех данных
func (rt *Router) Put(path string, handler Handler) {
	rt.statistic["routes"]++
	rt.Mux.Put(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handler(handler, w, r)
	})
}

// Patch запрос на обновление конкретных данных
func (rt *Router) Patch(path string, handler Handler) {
	rt.statistic["routes"]++
	rt.Mux.Patch(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handler(handler, w, r)
	})
}

// Delete запрос на удаление данных
func (rt *Router) Delete(path string, handler Handler) {
	rt.statistic["routes"]++
	rt.Mux.Delete(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handler(handler, w, r)
	})
}

// NotFound обрабатывает 404 ошибку
func (rt *Router) NotFound(handler Handler) {
	rt.statistic["routes"]++
	rt.Mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		rt.handler(handler, w, r)
	})
}

// MethodNotAllowed обрабатывает 405 ошибку
func (rt *Router) MethodNotAllowed(handler Handler) {
	rt.statistic["routes"]++
	rt.Mux.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		rt.handler(handler, w, r)
	})
}

// Use добавляет промежуточное программное обеспечение
func (rt *Router) Use(handlers ...Handler) {
	for _, handler := range handlers {
		rt.statistic["middlewares"]++

		currentHandler := handler
		rt.Mux.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				ctx = WithNextHandler(ctx, next)

				rt.handler(currentHandler, w, r.WithContext(ctx))
			})
		})
	}
}

// Group группирует роутеры
func (rt *Router) Group(fn func(r *Router)) *Router {
	im := rt.With()

	if fn != nil {
		fn(&im)
	}

	return &im
}

// With добавляет встроенное промежуточное программное обеспечение для обработчика конечной точки
func (rt *Router) With(middlewares ...func(http.Handler) http.Handler) Router {
	return Router{
		Mux:      rt.Mux.With(middlewares...).(*chi.Mux),
		Validate: rt.Validate,

		statistic: rt.statistic,
	}
}

// Route создает вложенность роутеров
func (rt *Router) Route(pattern string, fn func(r *Router)) *Router {
	subRouter := &Router{
		Mux:      chi.NewRouter(),
		Validate: rt.Validate,

		statistic: rt.statistic,
	}

	fn(subRouter)
	rt.Mount(pattern, subRouter)

	return subRouter
}

// Mount добавляет вложенность роутеров
func (rt *Router) Mount(pattern string, router *Router) {
	rt.Mux.Mount(pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.Mux.ServeHTTP(w, r)
	}))
}

// ServeHTTP возвращает весь пул роутеров
func (rt *Router) ServeHTTP() http.HandlerFunc {
	return rt.Mux.ServeHTTP
}

// Listen запускает сервер
// Реализация: https://github.com/go-chi/chi/blob/master/_examples/graceful/main.go
func (rt *Router) Listen(addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: rt.ServeHTTP(),
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	ch := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error("Не удалось запустить сервер", slog.String("err", err.Error()))
				ch <- ctx.Err()
			}
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		panic(err)
	case <-ctx.Done():
		timeoutCtx, done := context.WithTimeout(context.Background(), time.Second*10)
		defer done()

		go func() {
			<-timeoutCtx.Done()
			if errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) {
				log.Error("Время корректного завершения работы истекло. Принудительный выход", slog.String("err", timeoutCtx.Err().Error()))
			}
		}()

		if err := server.Shutdown(timeoutCtx); err != nil {
			log.Error("Не удалось остановить сервер", slog.String("err", err.Error()))
		}
	}

	return nil
}

// GetStatistic возвращает статистику использования роутеров
func (rt *Router) GetStatistic() map[string]int {
	return rt.statistic
}
