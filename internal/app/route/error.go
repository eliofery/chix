package route

import (
	"github.com/eliofery/go-chix/pkg/chix"
)

// ErrorRoute маршруты для обработки ошибок
func (r *route) ErrorRoute(router *chix.Router) {
	router.NotFound(func(ctx *chix.Ctx) error {
		return ctx.JSON(chix.Map{
			"success": false,
			"message": "Страница не найдена",
		})
	})

	router.MethodNotAllowed(func(ctx *chix.Ctx) error {
		return ctx.JSON(chix.Map{
			"success": false,
			"message": "Метод не поддерживается",
		})
	})
}
