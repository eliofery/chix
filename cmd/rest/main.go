package main

import (
	"github.com/eliofery/go-chix/internal/app/repository"
	"github.com/eliofery/go-chix/internal/app/service"
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/eliofery/go-chix/pkg/config"
	"github.com/eliofery/go-chix/pkg/config/viperr"
	"github.com/eliofery/go-chix/pkg/database"
	"github.com/eliofery/go-chix/pkg/database/postgres"
	"github.com/eliofery/go-chix/pkg/jwt"
	"github.com/eliofery/go-chix/pkg/log"
	"github.com/eliofery/go-chix/pkg/utils"
	"github.com/eliofery/go-chix/pkg/validate"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ru"
	"github.com/go-playground/validator/v10"
	"log/slog"
)

func main() {
	log.Info("Используемое окружение", slog.String("env", utils.GetEnv()))

	conf := config.MustInit(viperr.New(utils.GetEnv()))
	db := database.MustConnect(postgres.New(conf))
	valid := validate.New(validator.New(), ru.New(), en.New())
	tokenManager := jwt.NewTokenManager(conf)
	_ = tokenManager

	dao := repository.NewDAO(db.Conn)
	users, err := service.NewUserService(dao).GetUsers()
	if err != nil {
		log.Error("Не удалось получить пользователей", slog.String("err", err.Error()))
	}
	log.Info("Список пользователей", slog.Any("users", users))

	chix.NewApp(db, conf).
		UseExtends(valid).
		UseMiddlewares().
		UseRoutes().
		MustRun()
}
