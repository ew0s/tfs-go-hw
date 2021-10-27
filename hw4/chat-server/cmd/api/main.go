package main

import (
	"chat"
	"chat/memoryDB"
	"chat/pkg/handler"
	"chat/pkg/repository"
	"chat/pkg/service"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

var (
	ErrRunServer      = errors.New("cant run server")
	ErrShutdownServer = errors.New("error happened while shutting down")
)

func main() {
	db, err := memoryDB.NewMemoryDB()
	if err != nil {
		log.Fatal(err)
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	serv := new(chat.Server)
	go func() {
		if err := serv.Run("8000", handlers.InitRoutes()); err != nil && err != http.ErrServerClosed {
			log.Fatal(fmt.Errorf("%v: %v", ErrRunServer, err))
		}
	}()

	log.Info("Chat server started")
	log.Info("Signal captured: ", <-quit)
	log.Info("Chat server shutting down")

	if err := serv.Shutdown(context.Background()); err != nil {
		log.Errorf("%s: %s", ErrShutdownServer, err)
	}

	log.Info("Chat server successfully stopped")
}
