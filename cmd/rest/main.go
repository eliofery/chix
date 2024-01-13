package main

import (
	"github.com/eliofery/go-chix/pkg/config"
	"github.com/eliofery/go-chix/pkg/config/viperr"
	"github.com/eliofery/go-chix/pkg/database"
	"github.com/eliofery/go-chix/pkg/database/sqlite"
	"github.com/eliofery/go-chix/pkg/log"
	"github.com/eliofery/go-chix/pkg/utils"
	"log/slog"
)

func main() {
	log.Info("Используемое окружение", slog.String("env", utils.GetEnv()))

	conf := config.MustInit(viperr.New(utils.GetEnv()))
	db := database.MustConnect(sqlite.New(conf))
	_ = db
}
