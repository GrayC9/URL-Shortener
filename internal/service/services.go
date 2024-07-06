package service

import "github.com/GrayC9/URL-Shortener/internal/storage"

type DBManage interface {
	Save(short, original string)
	Get(short string) string
}

type Service struct {
	DBManage
}

func NewService(storage *storage.MaridDB) *Service {
	return &Service{
		DBManage: storage,
	}
}
