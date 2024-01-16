package main

import (
	"github.com/eliofery/go-chix/internal/app/controller"
	"github.com/eliofery/go-chix/internal/app/middleware"
	"github.com/eliofery/go-chix/internal/app/repository"
	"github.com/eliofery/go-chix/internal/app/service"
	"github.com/eliofery/go-chix/internal/route"
	"github.com/eliofery/go-chix/internal/validation"
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/eliofery/go-chix/pkg/config"
	"github.com/eliofery/go-chix/pkg/config/viperr"
	"github.com/eliofery/go-chix/pkg/database"
	"github.com/eliofery/go-chix/pkg/database/postgres"
	"github.com/eliofery/go-chix/pkg/jwt"
	"github.com/eliofery/go-chix/pkg/log"
	"github.com/eliofery/go-chix/pkg/utils"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/fr"
	"github.com/go-playground/locales/ru"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	fr_translations "github.com/go-playground/validator/v10/translations/fr"
	ru_translations "github.com/go-playground/validator/v10/translations/ru"
)

func main() {
	utils.PrintEnv(log.InitLog())

	conf := config.MustInit(viperr.New(utils.GetEnv()))
	db := database.MustConnect(postgres.New(conf))
	valid := chix.NewValidate(validator.New()).
		RegisterTagName("label").
		RegisterLocales(
			ru.New(),
			en.New(),
			fr.New(),
		).
		RegisterTranslations(chix.DefaultTranslations{
			"ru": ru_translations.RegisterDefaultTranslations,
			"en": en_translations.RegisterDefaultTranslations,
			"fr": fr_translations.RegisterDefaultTranslations,
		}).
		RegisterValidations(
			validation.TestValidate(),
		)
	tokenManager := jwt.NewTokenManager(conf)

	dao := repository.NewDAO(db.Conn)
	handler := controller.NewServiceController(
		service.NewAuthService(dao, tokenManager),
		service.NewUserService(dao),
	)
	routes := route.NewRouter(handler)

	chix.NewApp(db, conf).
		UseExtends(valid).
		UseMiddlewares(
			middleware.Cors(conf),
		).
		UseRoutes(
			routes.ErrorRoute,
			routes.UserRoute,
		).
		MustRun()
}
