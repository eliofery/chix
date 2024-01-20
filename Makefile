# Автоматическая сборка проекта при изменении файлов
watch:
	modd

# Сборка проекта
build:
	go build -o bin/rest cmd/rest/main.go

# Запуск базы данных
db:
	docker compose up -d

# Доступные команды
help:
	./bin/rest -help

