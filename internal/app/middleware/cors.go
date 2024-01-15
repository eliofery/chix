package middleware

import (
	"fmt"
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/eliofery/go-chix/pkg/config"
	"github.com/eliofery/go-chix/pkg/log"
	"github.com/go-chi/cors"
)

const defaultCorsMaxAge = 3600 // 1 час

// Cors настройки межсайтового взаимодействия
// Пример: https://github.com/go-chi/cors?tab=readme-ov-file#usage
func Cors(conf config.Config) chix.HandlerNext {
	log.Debug("Инициализация middleware Cors")

	url := fmt.Sprintf(
		"%s://%s:%s",
		conf.Get("http.protocol"), conf.Get("http.url"), conf.Get("http.port"),
	)

	maxAge, ok := conf.GetAny("http.cors.maxage").(int)
	if !ok {
		maxAge = defaultCorsMaxAge
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{url},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		MaxAge:           maxAge,
	})
}
