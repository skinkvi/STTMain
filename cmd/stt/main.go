package main

import (
	"STTMain/internal/config"
	"STTMain/internal/server"
	"STTMain/internal/service"
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Запуск
	// команда для запуска go run cmd/stt/main.go --config=./config/local.yaml
	cfg := config.MustLoad()

	// инициализация логера
	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.String("env", cfg.Env))

	// TODO: инициализация приложения (app)

	// TODO: запустить gRPC-сервер приложение

	var sttService service.SpeedTypingTestService

	grpcServer := server.NewServer(sttService)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Error("falied to start gRPC server", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Info("starting gRPC server", slog.String("address", listener.Addr().String()))
	if err := grpcServer.Start(listener); err != nil {
		log.Error("falied to start gRPC server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Gracefully stop the gRPC server when done
	grpcServer.Stop(context.Background())
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
