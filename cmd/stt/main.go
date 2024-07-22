package main

import (
	"STTMain/internal/config"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Запуск
	cfg := config.MustLoad()

	// инициализация логера
	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.String("env", cfg.Env))

	// TODO: инициализация приложения (app)

	// TODO: запустить gRPC-сервер приложение
}

// для логгера
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

// func setupPrettySlog() *slog.Logger {
// 	opts := slogpretty.PrettyHandlerOptions{
// 		SlogOpts: &slog.HandlerOptions{
// 			Level: slog.LevelDebug,
// 		},
// 	}

// 	handler := opts.NewPrettyHandler(os.Stdout)

// 	return slog.New(handler)
// }
