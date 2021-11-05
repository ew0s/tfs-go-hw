package main

import (
	"trade-bot/configs"
	"trade-bot/internal/app"
	"trade-bot/internal/pkg/handler"
	"trade-bot/internal/pkg/repository"
	"trade-bot/internal/pkg/service"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := initConfig()
	if err != nil {
		log.Fatalf("unable to init config files: %s\n", err)
	}

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(app.Server)

	if err := srv.Run(config.Server.Port, handlers.InitRoutes()); err != nil {
		log.Fatalln("server crushed")
	}
}

func initConfig() (configs.Configuration, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("configs")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	var c configs.Configuration

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return configs.Configuration{}, err
		}
	}

	err := viper.Unmarshal(&c)
	return c, err
}
