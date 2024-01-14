package main

import (
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
	db.MigrateMust()
	valid := validate.New(validator.New(), ru.New(), en.New())
	_ = valid
	tokenManager := jwt.NewTokenManager(conf)
	_ = tokenManager
}
