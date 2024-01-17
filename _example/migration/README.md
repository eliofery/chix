# База данных

База данных поддерживаются двух видов:

- postgres
- sqlite

В зависимости от того какой пакет вы инициализируете при запуске микросервиса.

```go
postgres := database.MustConnect(postgres.New(conf))

sqlite := database.MustConnect(sqlite.New(conf))
```

## Миграции

Миграции хранятся в каталоге **internal/migration**. Для миграций используется пакет **goose**.

[goose - миграция базы данных](https://github.com/pressly/goose)

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
go get github.com/pressly/goose/v3
```

Миграция базы данных происходит автоматически при запуске микросервиса. Важно иметь в каталоге с миграциями хотя бы один sql файл с миграцией иначе приложение не запустится. 

### Пример описания файла миграции:

Его так же можно использовать в качестве заглушки чтобы микросервис не паниковал если миграции не нужны.

```sql
-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
```

### Создание миграции

```bash
# Создание миграции
goose -dir ./internal/migration create <имя миграции> sql

# Переименовывает миграции с формата даты создания в порядковый номер создания
# 20250104093011_<имя миграции>.sql -> 00001_<имя миграции>.sql
goose -dir ./internal/migration fix
```

### Статус миграций

```bash
# Вариант 1 (длинный)
goose -dir internal/migration postgres "postgresql://<пользователь>:<пароль>@<хост>:<порт>/<БД>?sslmode=disable" status

# Вариант 2 (короткий)
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING=postgresql://<пользователь>:<пароль>@<хост>:<порт>/<БД>?sslmode=disable

goose -dir internal/migration status
```

### Выполнить миграцию

```bash
goose -dir internal/migration up
```

### Откат миграции

```bash
goose -dir internal/migration down
```
