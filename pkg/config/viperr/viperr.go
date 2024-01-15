package viperr

import (
	"errors"
	"github.com/eliofery/go-chix/pkg/config"
	"github.com/eliofery/go-chix/pkg/log"
	"github.com/spf13/viper"
	"log/slog"
	"strings"
)

const (
	defaultConfigType = "yml"
	defaultConfigPath = "internal/config"
)

// Viperr загрузка конфигураций из yml файлов
type Viperr interface {
	config.Config

	AddConfigType(configType string) *viperr
	AddConfigPath(configPath ...string) *viperr
	GetAny(key string) any
}

type viperr struct {
	configName  string
	configType  string
	configPaths []string
}

// New конструктор viperr
func New(configName string, configPaths ...string) Viperr {
	log.Debug("Инициализация конфигурации Viperr")

	paths := []string{defaultConfigPath}
	if len(configPaths) > 0 {
		paths = configPaths
	}

	return &viperr{
		configName:  configName,
		configType:  defaultConfigType,
		configPaths: paths,
	}
}

// AddConfigType добавление типа конфигурации
func (v *viperr) AddConfigType(configType string) *viperr {
	v.configType = configType

	return v
}

// AddConfigPath добавление путей конфигурации
func (v *viperr) AddConfigPath(paths ...string) *viperr {
	v.configPaths = append(v.configPaths, paths...)

	return v
}

// Init загрузка конфигурации
func (v *viperr) Init() error {
	viper.SetConfigName(v.configName)
	viper.SetConfigType(v.configType)
	for _, configPath := range v.configPaths {
		viper.AddConfigPath(configPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError

		if errors.As(err, &configFileNotFoundError) {
			log.Error("Конфигурационный файл не найден", slog.String("err", err.Error()))
			return err
		}

		log.Error("Не удалось прочитать конфигурационный файл", slog.String("err", err.Error()))
		return err
	}

	return nil
}

// Get получение конфигурации в формате string
func (v *viperr) Get(key string) string {
	key = formatter(key)

	return viper.GetString(key)
}

// GetAny получение конфигурации в любом формате
func (v *viperr) GetAny(key string) any {
	key = formatter(key)

	return viper.Get(key)
}

// Formatter изменяет строку под конфигурацию
// Пример: HTTP_PORT -> http.port
func formatter(key string) string {
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "_", ".")

	return key
}
