package storage

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Word struct {
	gorm.Model
	Text string `gorm:"not null"`
}

type PostgresStorage struct {
	DB *gorm.DB
}

func NewPostgresStorage(dataSourceName string) (*PostgresStorage, error) {
	const op = "internal.storage.NewPostgresStorage"

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dataSourceName,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.AutoMigrate(&Word{})
	if err != nil {
		log.Fatalf("миграция не прошла: %s: %v", op, err)
	}

	return &PostgresStorage{DB: db}, nil
}

func (s *PostgresStorage) Close() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (s *PostgresStorage) GetRandomWords(count int) ([]Word, error) {
	var words []Word
	if err := s.DB.Order("RANDOM()").Limit(count).Find(&words).Error; err != nil {
		return nil, err
	}
	return words, nil
}
