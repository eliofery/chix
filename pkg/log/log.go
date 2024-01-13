package log

import (
	"context"
	"fmt"
	"github.com/eliofery/go-chix/pkg/utils"
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// InitLog инициализирует логгер
// Пример: https://github.com/golang/go/issues/59145#issuecomment-1481920720
func InitLog() *slog.Logger {
	return slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      getLevel(),
		AddSource:  true,
		TimeFormat: "2006/01/02 15:04:05",
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.SourceKey:
				return replaceSource(a)
			default:
				return a
			}
		},
	}))
}

// getLevel возвращает уровень логгирования
func getLevel() slog.Level {
	var level slog.Level

	switch utils.GetEnv() {
	case utils.Local:
		level = slog.LevelDebug
	case utils.Prod:
		level = slog.LevelInfo
	}

	return level
}

// replaceSource путь до файла в котором был вызван лог
// Реализация с использованием ASCII
// absPath, err := filepath.Abs(source.File)
// formattedPath := fmt.Sprintf("\x1b]8;;file://%v\x1b\\%s:%v\x1b]8;;\x1b\\", absPath, relPath, source.Line)
func replaceSource(a slog.Attr) slog.Attr {
	source := a.Value.Any().(*slog.Source)

	pwd, err := os.Getwd()
	if err != nil {
		return a
	}

	relPath, err := filepath.Rel(pwd, source.File)
	if err != nil {
		return a
	}

	formattedPath := fmt.Sprintf("%s:%d", relPath, source.Line)

	return slog.Attr{
		Key:   a.Key,
		Value: slog.StringValue(formattedPath),
	}
}

// Debug выводит сообщение в лог с уровнем debug
func Debug(msg string, args ...any) {
	l := InitLog()
	if !l.Enabled(context.Background(), slog.LevelDebug) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // пропустить [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelDebug, msg, pcs[0])
	r.Add(args...)

	if err := l.Handler().Handle(context.Background(), r); err != nil {
		l.Warn("При логировании произошла ошибка", slog.String("err", err.Error()))
	}
}

// Info выводит сообщение в лог с уровнем info
func Info(msg string, args ...any) {
	l := InitLog()
	if !l.Enabled(context.Background(), slog.LevelInfo) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelInfo, msg, pcs[0])
	r.Add(args...)

	if err := l.Handler().Handle(context.Background(), r); err != nil {
		l.Warn("При логировании произошла ошибка", slog.String("err", err.Error()))
	}
}

// Warn выводит сообщение в лог с уровнем warn
func Warn(msg string, args ...any) {
	l := InitLog()
	if !l.Enabled(context.Background(), slog.LevelWarn) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelWarn, msg, pcs[0])
	r.Add(args...)

	if err := l.Handler().Handle(context.Background(), r); err != nil {
		l.Warn("При логировании произошла ошибка", slog.String("err", err.Error()))
	}
}

// Error выводит сообщение в лог с уровнем error
func Error(msg string, args ...any) {
	l := InitLog()
	if !l.Enabled(context.Background(), slog.LevelError) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelError, msg, pcs[0])
	r.Add(args...)

	if err := l.Handler().Handle(context.Background(), r); err != nil {
		l.Warn("При логировании произошла ошибка", slog.String("err", err.Error()))
	}
}
