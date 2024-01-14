package database

import (
	"database/sql"
	"fmt"
	"github.com/eliofery/go-chix/internal"
	"github.com/eliofery/go-chix/pkg/log"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"log/slog"
)

const (
	dirMigration = "migration"
)

// Migration миграция базы данных
type Migration interface {
	Migrate() error // Migrate миграция базы данной
	MigrateMust()   // MigrateMust миграция базы данной с обработкой ошибок
}

// Migrate миграция базы данных
func (db *DB) Migrate() error {
	log.Debug("Миграция базы данных")

	goose.SetBaseFS(internal.EmbedMigration)
	defer goose.SetBaseFS(nil)

	if err := goose.SetDialect(getCurrentDialect(db.Conn)); err != nil {
		log.Error("Не удалось установить используемую БД", slog.String("err", err.Error()))
		return err
	}

	if err := goose.Up(db.Conn, dirMigration); err != nil {
		log.Error("Не удалось выполнить миграцию базы данной", slog.String("err", err.Error()))
		return err
	}

	return nil
}

// MigrateMust миграция базы данных с обработкой ошибок
func (db *DB) MigrateMust() {
	if err := db.Migrate(); err != nil {
		panic(err)
	}
}

// GetCurrentDialect получение названия используемой базы данной
func getCurrentDialect(db *sql.DB) string {
	var dialect goose.Dialect

	switch db.Driver().(type) {
	case *mysql.MySQLDriver:
		dialect = goose.DialectMySQL
	case *sqlite3.SQLiteDriver:
		dialect = goose.DialectSQLite3
	case *stdlib.Driver:
		dialect = goose.DialectPostgres
	default:
		log.Error("не удалось определить используемую БД", slog.Any("dialect", fmt.Sprintf("%T", db.Driver())))
	}

	return string(dialect)
}
