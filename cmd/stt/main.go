package main

import (
	"STTMain/internal/config"
	"STTMain/internal/storage"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	// Запуск
	// команда для запуска go run cmd/stt/main.go --config=./config/local.yaml
	cfg := config.MustLoad()

	// инициализация логера
	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.String("env", cfg.Env))

	// Инициализация хранилища
	postgresStorage, err := storage.NewPostgresStorage(cfg.Storage.Postgres.URL)
	if err != nil {
		log.Error("failed to initialize Postgres storage", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer postgresStorage.Close()

	// Инициализация сервиса

	// Настройка HTTP сервера
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/api/text", func(w http.ResponseWriter, r *http.Request) {
		words, err := postgresStorage.GetRandomWords(25)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"text": words[0].Text})
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("failed to upgrade WebSocket connection", slog.String("error", err.Error()))
			return
		}
		defer conn.Close()

		// Отправка текста для ввода
		words, err := postgresStorage.GetRandomWords(25)
		if err != nil {
			log.Error("failed to get random words", slog.String("error", err.Error()))
			return
		}
		textToType := ""
		for _, word := range words {
			textToType += word.Text + " "
		}
		conn.WriteJSON(map[string]string{"type": "text", "text": textToType})

		// Обработка ввода пользователя
		startTime := time.Now()
		var inputText string
		var correctIndices []int
		var incorrectIndices []int
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Error("failed to read WebSocket message", slog.String("error", err.Error()))
				break
			}
			var data struct {
				Type string `json:"type"`
				Text string `json:"text"`
			}
			if err := json.Unmarshal(msg, &data); err != nil {
				log.Error("failed to unmarshal WebSocket message", slog.String("error", err.Error()))
				continue
			}
			if data.Type == "input" {
				inputText = data.Text
				correctIndices, incorrectIndices = highlightErrors(textToType, inputText)
				conn.WriteJSON(map[string]interface{}{
					"type":      "highlight",
					"correct":   correctIndices,
					"incorrect": incorrectIndices,
				})
			}
			if len(inputText) >= len(textToType) {
				break
			}
		}

		// Подсчет ошибок и скорости
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		errors := len(incorrectIndices)
		speed := float64(len(inputText)) / duration.Minutes()

		// Отправка результатов
		conn.WriteJSON(map[string]interface{}{
			"type":   "result",
			"errors": errors,
			"speed":  speed,
		})
	})

	log.Info("starting HTTP server", slog.String("address", ":8080"))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Error("failed to start HTTP server", slog.String("error", err.Error()))
		os.Exit(1)
	}
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

func highlightErrors(textToType, inputText string) ([]int, []int) {
	var correctIndices []int
	var incorrectIndices []int

	for i, char := range inputText {
		if i < len(textToType) && char == rune(textToType[i]) {
			correctIndices = append(correctIndices, i)
		} else {
			incorrectIndices = append(incorrectIndices, i)
		}
	}

	return correctIndices, incorrectIndices
}
