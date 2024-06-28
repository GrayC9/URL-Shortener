package storage

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type Storage struct {
	urls  map[string]string
	mutex sync.RWMutex
	logg  *logrus.Logger
}

func NewStorage(log *logrus.Logger) *Storage {
	return &Storage{
		urls: make(map[string]string),
		logg: log,
	}
}

func (s *Storage) Save(short, original string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.urls[short] = original
}

func (s *Storage) Get(short string) string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if url, ok := s.urls[short]; ok {
		return url
	}
	s.logg.Errorln("No such shortcode in storage")
	return ""
}
