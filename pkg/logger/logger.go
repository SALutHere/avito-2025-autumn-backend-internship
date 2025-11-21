package logger

import (
	"log/slog"
	"os"
	"sync"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var (
	log  *slog.Logger
	once sync.Once
)

func Setup(env string) *slog.Logger {
	once.Do(func() {
		switch env {
		case envLocal:
			log = slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}),
			)
		case envDev:
			log = slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}),
			)
		case envProd:
			log = slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
					Level: slog.LevelInfo,
				}),
			)
		default:
			log = slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}),
			)
		}
	})

	return log
}

func L() *slog.Logger {
	return log
}
