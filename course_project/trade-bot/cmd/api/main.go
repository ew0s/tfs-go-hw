package main

import (
	"trade-bot/internal/app"
	"trade-bot/internal/pkg/handler"
	"trade-bot/internal/pkg/repository"
	"trade-bot/internal/pkg/service"

	log "github.com/sirupsen/logrus"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(app.Server)

	if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalln("server crushed")
	}
}
