# Конфигурации микросервиса

Конфигурации поддерживаются двух видов:

- *.env (godotenv)
- *.yml (viperr)

В зависимости от того какой пакет вы инициализируете при запуске микросервиса.

```go
viperr := config.MustInit(viperr.New("local"))

godotenv := config.MustInit(godotenv.New("local"))
```

Какой бы вариант не был выбран доступ к значениям можно будет получить следующим образом:

```go
// Значение вернется как строка. 

protocol := viperr.Get("HTTP_PROTOCOL")
protocol := viperr.Get("http.protocol")

protocol := godotenv.Get("HTTP_PROTOCOL")
protocol := godotenv.Get("http.protocol")
```

Разница в использовании того или иного варианта зависит от личных предпочтений, в каком расширении файла больше нравится хранить конфигурации .env или .yml.

## Обязательные параметры конфигурации

Если какого либо из этих параметров будет отсутствовать микросервис не запустится.

```yml
# При использовании postgres
postgres:
    host: <значение>
    user: <значение>
    password: <значение>
    database: <значение>

# При использовании sqlite
sqlite:
    path: internal/<имя>.db

jwt:
    secret: <значение>
```
