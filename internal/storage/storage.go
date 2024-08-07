package storage

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type URL struct {
	OriginalURL    string    `json:"original_url"`
	ShortCode      string    `json:"short_code"`
	LastAccessedAt time.Time `json:"last_accessed_at"`
}

type Storage interface {
	SaveURL(shortCode, originalURL string) error
	GetURL(shortCode string) (string, error)
	GetShortCode(originalURL string) (string, error)
	IncrementClickCount(shortCode string) error
	UpdateLastAccessed(shortCode string) error
	GetPopularURLs(limit int) ([]URL, error)
	CreateUser(string, string) error
	EnterUser(string, string) (int, error)
}

type MariaDBStorage struct {
	db *sql.DB
}

func NewMariaDBStorage(config string) (*MariaDBStorage, error) {
	db, err := sql.Open("mysql", config)
	if err != nil {
		return nil, err
	}
	if err := RetryPing(db); err != nil {
		return nil, err
	}

	return &MariaDBStorage{db: db}, nil
}

func RetryPing(db *sql.DB) error {
	var err error
	for i := 0; i < 5; i++ {
		if err = db.Ping(); err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	return err
}

func (m *MariaDBStorage) SaveURL(shortCode, originalURL string) error {
	_, err := m.db.Exec("INSERT INTO urls (short_code, original_url) VALUES (?, ?)", shortCode, originalURL)
	if err != nil {
		return err
	}
	return nil
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
	_, err := m.db.Exec("UPDATE urls SET last_accessed_at = ? WHERE short_code = ?", time.Now(), shortCode)
	return err
}

func (m *MariaDBStorage) CreateUser(name, hash_password string) error {
	query := "INSERT INTO users (username, hash_password, created_at) VALUES (?, ?, ?)"
	tm := time.Now()
	_, err := m.db.Exec(query, name, hash_password, tm.Format(time.DateTime))
	if err != nil {
		return err
	}
	return nil
}

func (m *MariaDBStorage) GetPopularURLs(limit int) ([]URL, error) {
	rows, err := m.db.Query("SELECT original_url, short_code FROM urls ORDER BY click_count DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []URL
	for rows.Next() {
		var url URL
		if err := rows.Scan(&url.OriginalURL, &url.ShortCode); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func (m *MariaDBStorage) EnterUser(name, pass string) (int, error) {
	query := "SELECT id, hash_password FROM users WHERE username = ?"
	var id int
	var hash string
	err := m.db.QueryRow(query, name).Scan(&id, &hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("no such user")
		}
		return 0, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)); err != nil {
		return 0, err
	}
	return id, nil
}
