package route

import (
	"github.com/eliofery/go-chix/internal/app/controller"
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/eliofery/go-chix/pkg/log"
)

const (
	apiV1 = "/api/v1"
)

// Route маршрутизатор
type Route interface {
	ErrorRoute(router *chix.Router) // ErrorRoute маршруты для обработки ошибок
	UserRoute(router *chix.Router)  // UserRoute маршруты для пользователей
}

type route struct {
	Handler *controller.ServiceController
}

// NewRouter создание маршрутов
func NewRouter(handler controller.ServiceController) Route {
	log.Debug("Инициализация маршрутов")

	return &route{Handler: &handler}
}
