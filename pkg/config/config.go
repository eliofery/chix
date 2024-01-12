package config

import (
	"github.com/eliofery/go-chix/pkg/log"
)

// Config конфигурация
type Config interface {
	Init() error           // Init загрузка конфигурации
	Get(key string) string // Get получение конфигурации в формате string
	GetAny(key string) any // GetAny получение конфигурации в любом формате
}

// Load загрузка конфигурации
// Пример: config.Load(viperr.New())
func Load(config Config) (Config, error) {
	log.Debug("Загрузка конфигурации")

	if err := config.Init(); err != nil {
		return nil, err
	}

	return config, nil
}

// MustInit инициализация конфигурации с обработкой ошибок
func MustInit(config Config) Config {
	conf, err := Load(config)
	if err != nil {
		panic(err)
	}

	return conf
}
