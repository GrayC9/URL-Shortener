package storage

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type Storage interface {
	SaveURL(shortCode, originalURL string) error
	GetURL(shortCode string) (string, error)
	GetShortCode(originalURL string) (string, error)
}

type MariaDBStorage struct {
	db *sql.DB
}

func NewMariaDBStorage(dsn string) (*MariaDBStorage, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &MariaDBStorage{db: db}, nil
}

func (s *MariaDBStorage) SaveURL(shortCode, originalURL string) error {
	_, err := s.db.Exec("INSERT INTO urls (short_code, original_url) VALUES (?, ?)", shortCode, originalURL)
	return err
}

func (s *MariaDBStorage) GetURL(shortCode string) (string, error) {
	var originalURL string
	err := s.db.QueryRow("SELECT original_url FROM urls WHERE short_code = ?", shortCode).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("url not found")
		}
		return "", err
	}
	return originalURL, nil
}

func (s *MariaDBStorage) GetShortCode(originalURL string) (string, error) {
	var shortCode string
	err := s.db.QueryRow("SELECT short_code FROM urls WHERE original_url = ?", originalURL).Scan(&shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("short url not found")
		}
		return "", err
	}
	return shortCode, nil
}
