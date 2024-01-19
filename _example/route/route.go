package route

import (
	"github.com/eliofery/go-chix/internal/app/controller"
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/eliofery/go-chix/pkg/log"
)

const (
	apiV1 = "/v1"
)

// Route маршрутизатор
type Route interface {
	ErrorRoute(*chix.Router)
	AuthRoute(*chix.Router)
}

type route struct {
	Handler *controller.ServiceController
}

// NewRouter создание маршрутов
func NewRouter(handler controller.ServiceController) Route {
	log.Debug("Инициализация маршрутов")

	return &route{Handler: &handler}
}
