package chix

import (
	"fmt"
	"github.com/eliofery/go-chix/pkg/config"
	"github.com/eliofery/go-chix/pkg/database"
	"github.com/eliofery/go-chix/pkg/log"
	"github.com/fatih/color"
	"log/slog"
	"os"
	"strings"
	"time"
)

const (
	version         = "1.0.0"
	protocolDefault = "http"
	urlDefault      = "127.0.0.1"
	portDefault     = "3000"
)

// Тип key ключ контекста
type key string

// App структура фреймворка Chix
// Основой фреймворка является Сhi роутер
// Chix это обертка вокруг chi роутера, с расширением функционала
// Теперь роуты выглядят как во всех современных фреймворка
// func (ctx *router.Ctx) error {}
type App struct {
	config      config.Config
	db          *database.DB
	middlewares []Handler
	routes      []func(router *Router)
	validate    Validate
}

func NewApp(db *database.DB, cfg config.Config) *App {
	log.Debug("Инициализация chix")

	return &App{
		config: cfg,
		db:     db,
	}
}

// UseExtends использование расширений
// На данный момент поддерживается только валидатор
func (a *App) UseExtends(extends ...any) *App {
	log.Debug("Регистрация расширений")

	for _, extend := range extends {
		switch extend.(type) {
		case Validate:
			a.validate = extend.(Validate)
		default:
			log.Warn("Неизвестное расширение", slog.String("extend", fmt.Sprintf("%T", extend)))
		}
	}

	return a
}

// UseMiddlewares использование промежуточное программное обеспечение
func (a *App) UseMiddlewares(middlewares ...Handler) *App {
	log.Debug("Регистрация middlewares")

	a.middlewares = append(a.middlewares, middlewares...)

	return a
}

// UseRoutes использование маршрутов
func (a *App) UseRoutes(routes ...func(router *Router)) *App {
	log.Debug("Регистрация маршрутов")

	a.routes = append(a.routes, routes...)

	return a
}

// MustRun запуск приложения с обработкой ошибок
func (a *App) MustRun() {
	defer func() {
		if err := a.db.Conn.Close(); err != nil {
			log.Error("Не удалось закрыть соединение с базой данных", slog.String("err", err.Error()))
		}
	}()

	if err := a.db.Migrate(); err != nil {
		log.Error("Не удалось выполнить миграцию базы данных", slog.String("err", err.Error()))
	}

	server := NewRouter(a.validate)
	a.registerMiddlewares(server, a.middlewares)
	a.registerRoutes(server, a.routes)
	a.printLogo(server.GetStatistic())

	log.Info("Запуск сервера")
	if err := a.listen(server); err != nil {
		panic(err)
	}
	log.Info("Остановка сервера")
}

// RegisterMiddlewares регистрация промежуточного программного обеспечения
func (a *App) registerMiddlewares(r *Router, middlewares []Handler) {
	r.Use(middlewares...)
}

// RegisterRoutes регистрация маршрутов
func (a *App) registerRoutes(r *Router, routes []func(router *Router)) {
	for _, route := range routes {
		route(r)
	}
}

// MustListen регистрация маршрутов
func (a *App) listen(r *Router) error {
	return r.Listen(fmt.Sprintf("%s:%s", a.getUrl(), a.getPort()))
}

// printLogo печать логотипа в терминал
func (a *App) printLogo(statistic map[string]int) {
	c := color.New(color.FgHiBlue).Add(color.Bold, color.Italic)
	h := color.New(color.FgHiCyan).Add(color.Bold)
	i := color.New(color.FgHiYellow).Add(color.Bold, color.Italic)
	x := color.New(color.FgHiGreen).Add(color.Bold)

	logo := `
+---------------------------------+
|         Chix ver. ` + version + `         |
+---------------------------------+
| СССССС  HH  HH  IIIIII  XX   XX |
| ССС     HH  HH    II     XX XX  |
| СС      HHHHHH    II      XXX   |
| ССС     HH  HH    II     XX XX  |
| СССССС  HH  HH  IIIIII  XX   XX |
+---------------------------------+
`

	lines := strings.Split(logo, "\n")
	for key, line := range lines {
		if key == 0 {
			continue
		}

		for _, char := range line {
			switch char {
			case 'С':
				_, _ = c.Print("@")
			case 'H':
				_, _ = h.Print("@")
			case 'I':
				_, _ = i.Print("@")
			case 'X':
				_, _ = x.Print("@")
			default:
				if key == 2 {
					_, _ = color.New(color.FgHiWhite).Add(color.Bold).Print(string(char))
					continue
				}
				fmt.Print(string(char))
			}
			time.Sleep(time.Millisecond * 2)
		}

		if key != len(lines)-1 {
			time.Sleep(time.Millisecond * 15)
			fmt.Println()
		}
	}

	color.HiBlue(fmt.Sprintf("| Количество middlewares: %d\n", statistic["middlewares"]))
	time.Sleep(time.Millisecond * 100)
	color.HiCyan(fmt.Sprintf("| Количество маршрутов: %d\n", statistic["routes"]))
	time.Sleep(time.Millisecond * 100)
	color.HiYellow(fmt.Sprintf("| PID: %d\n", os.Getpid()))
	time.Sleep(time.Millisecond * 100)
	color.HiGreen(fmt.Sprintf("| Сервер: %s://%s:%s\n", a.getProtocol(), a.getUrl(), a.getPort()))
	fmt.Println(`+---------------------------------+`)
}

func (a *App) getProtocol() string {
	protocol := a.config.Get("http.protocol")
	if protocol == "" {
		protocol = protocolDefault
	}

	return protocol
}

func (a *App) getUrl() string {
	url := a.config.Get("http.url")
	if url == "" {
		url = urlDefault
	}

	return url
}

func (a *App) getPort() string {
	port := a.config.Get("http.port")
	if port == "" {
		port = portDefault
	}

	return port
}
