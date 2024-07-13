package storage

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Storage interface {
	SaveURL(shortCode, originalURL string) error
	GetURL(shortCode string) (string, error)
	GetShortCode(originalURL string) (string, error)
	IncrementClickCount(shortCode string) error
	UpdateLastAccessed(shortCode string) error
}

type MariaDBStorage struct {
	db *sql.DB
}

func NewMariaDBStorage(dsn string) (*MariaDBStorage, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &MariaDBStorage{db: db}, nil
}

func (m *MariaDBStorage) SaveURL(shortCode, originalURL string) error {
	_, err := m.db.Exec("INSERT INTO urls (short_code, original_url) VALUES (?, ?)", shortCode, originalURL)
	return err
}

func (m *MariaDBStorage) GetURL(shortCode string) (string, error) {
	var originalURL string
	err := m.db.QueryRow("SELECT original_url FROM urls WHERE short_code = ?", shortCode).Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

func (m *MariaDBStorage) GetShortCode(originalURL string) (string, error) {
	var shortCode string
	err := m.db.QueryRow("SELECT short_code FROM urls WHERE original_url = ?", originalURL).Scan(&shortCode)
	if err != nil {
		return "", err
	}
	return shortCode, nil
}

func (m *MariaDBStorage) IncrementClickCount(shortCode string) error {
	_, err := m.db.Exec("UPDATE urls SET click_count = click_count + 1 WHERE short_code = ?", shortCode)
	return err
}

func (m *MariaDBStorage) UpdateLastAccessed(shortCode string) error {
	_, err := m.db.Exec("UPDATE urls SET last_accessed_at = ? WHERE short_code = ?", time.Now())
	return err
}
