package route

import (
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/eliofery/go-chix/pkg/log"
)

// ErrorRoute маршруты для обработки ошибок
func (r *route) ErrorRoute(router *chix.Router) {
	log.Debug("Инициализация ErrorRoute")

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
