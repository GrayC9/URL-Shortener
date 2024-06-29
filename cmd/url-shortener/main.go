package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/GrayC9/URL-Shortener/internal/handlers"
	"github.com/GrayC9/URL-Shortener/internal/server"
	"github.com/GrayC9/URL-Shortener/internal/service"
	"github.com/GrayC9/URL-Shortener/internal/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	logg := logrus.New()
	storage := storage.NewStorage(logg)
	service := service.NewService(storage)
	router := handlers.NewRouter(service, logg)
	srv := server.New(logg)

	go func() {
		if err := srv.Run(router); err != nil {
			logg.Fatalln(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	go func() {
		var wg sync.WaitGroup

		wg.Add(1)

		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logg.Errorln(err)
			return
		}

		wg.Wait()
	}()

	logg.Infoln("Gracefully...")
}
