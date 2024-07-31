package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Word struct {
	gorm.Model
	Text string `gorm:"not null"`
}

func main() {
	// Подключение к базе данных
	dsn := "postgres://postgres:password@localhost:5433/STTDB?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Миграция модели
	err = db.AutoMigrate(&Word{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Чтение слов из файла и вставка их в базу данных
	file, err := os.Open("words.txt")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		err = db.Create(&Word{Text: word}).Error
		if err != nil {
			log.Printf("failed to insert word: %v", err)
		} else {
			fmt.Printf("Inserted word: %s\n", word)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
}
