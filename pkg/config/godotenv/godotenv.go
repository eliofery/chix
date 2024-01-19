package godotenv

import (
	"fmt"
	"github.com/eliofery/go-chix/pkg/config"
	"github.com/eliofery/go-chix/pkg/log"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strings"
)

const (
	defaultConfigType = "env"
)

// GoDotEnv загрузка конфигураций из переменного окружения
type GoDotEnv interface {
	config.Config
}

type goDotEnv struct {
	configName string
}

// New конструктор goDotEnv
func New(configName string) GoDotEnv {
	log.Debug("инициализация конфигурации godotenv")

	configName = fmt.Sprintf("%s.%s", configName, defaultConfigType)

	return &goDotEnv{
		configName: configName,
	}
}

// Init загрузка конфигурации
func (g *goDotEnv) Init() error {
	if err := godotenv.Load(g.configName); err != nil {
		log.Error("Не удалось загрузить переменные окружения", slog.String("err", err.Error()))
		return err
	}

	return nil
}

// Get получение конфигурации в формате string
func (g *goDotEnv) Get(key string) string {
	key = formatter(key)

	return os.Getenv(key)
}

// GetAny получение конфигурации в любом формате
func (g *goDotEnv) GetAny(key string) any {
	key = formatter(key)

	return os.Getenv(key)
}

// Formatter изменяет строку под конфигурацию
// Пример: http.port -> HTTP_PORT
func formatter(key string) string {
	key = strings.ToUpper(key)
	key = strings.ReplaceAll(key, ".", "_")

	return key
}
