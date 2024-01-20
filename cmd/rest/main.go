package main

import (
	"flag"
	"github.com/eliofery/go-chix/internal/app/controller"
	"github.com/eliofery/go-chix/internal/app/middleware"
	"github.com/eliofery/go-chix/internal/app/repository"
	"github.com/eliofery/go-chix/internal/app/route"
	"github.com/eliofery/go-chix/internal/app/service"
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/eliofery/go-chix/pkg/config"
	"github.com/eliofery/go-chix/pkg/config/viperr"
	"github.com/eliofery/go-chix/pkg/database"
	"github.com/eliofery/go-chix/pkg/database/postgres"
	"github.com/eliofery/go-chix/pkg/jwt"
	"github.com/eliofery/go-chix/pkg/log"
	"github.com/eliofery/go-chix/pkg/utils"
	"github.com/go-playground/validator/v10"
)

// TODO: изменить логику получение окружения (utils.GetEnv())
func init() {
	flag.String("env", "local", "Имя окружения должно соответствовать имени конфигурационного файла")
}

func main() {
	flag.Parse()
	utils.PrintEnv(log.InitLog())

	conf := config.MustInit(viperr.New(utils.GetEnv()))
	db := database.MustConnect(postgres.New(conf))
	tokenManager := jwt.NewTokenManager(conf)
	validate := chix.NewValidate(validator.New())

	dao := repository.NewDAO(db.Conn)
	handler := controller.NewServiceController(
		service.NewAuthService(dao, tokenManager),
	)
	routes := route.NewRouter(handler)

	chix.NewApp(db, conf).
		UseExtends(validate).
		UseMiddlewares(
			middleware.Cors(conf),
			middleware.SetUserIdFromToken(dao, tokenManager),
		).
		UseRoutes(
			routes.ErrorRoute,
			routes.AuthRoute,
		).
		MustRun()
}
