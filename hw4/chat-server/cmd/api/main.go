package main

import (
	"chat"
	"chat/memoryDB"
	"chat/pkg/handler"
	"chat/pkg/repository"
	"chat/pkg/service"
	"errors"
	log "github.com/sirupsen/logrus"
)

var (
	ErrRunServer = errors.New("cant run server")
)

func main() {
	log.SetFormatter(new(log.JSONFormatter))

	db, err := memoryDB.NewMemoryDB()
	if err != nil {
		log.Fatal(err)
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	serv := new(chat.Server)
	if err := serv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatal(ErrRunServer)
	}
}
