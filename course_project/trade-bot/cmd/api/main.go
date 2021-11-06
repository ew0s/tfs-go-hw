package main

import (
	"fmt"
	"os"
	"trade-bot/configs"
	"trade-bot/internal/app"
	"trade-bot/internal/pkg/handler"
	"trade-bot/internal/pkg/repository"
	"trade-bot/internal/pkg/service"

	"github.com/joho/godotenv"

	"github.com/pkg/errors"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var (
	ErrUnableToInitConfig       = errors.New("unable to init config files")
	ErrReadConfig               = errors.New("read config")
	ErrRunServer                = errors.New("run server")
	ErrUnableToConnectToDB      = errors.New("unable to connect to database")
	ErrUnableToLoadEnvVariables = errors.New("unable to load enviroment variables")
)

func main() {
	config, err := initConfig()
	if err != nil {
		log.Fatalf("%s: %s", ErrUnableToInitConfig, err)
	}

	db, err := repository.NewPostgresDB(config.Database)
	if err != nil {
		log.Fatalf("%s: %s", ErrUnableToConnectToDB, err)
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(app.Server)
	if err := srv.Run(config.Server.Port, handlers.InitRoutes()); err != nil {
		log.Fatalf("%s: %s", ErrRunServer, err)
	}
}

func initConfig() (configs.Configuration, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("configs")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal(fmt.Errorf("%s: %s", ErrReadConfig, err))
		}
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal(fmt.Errorf("%s: %s", ErrUnableToLoadEnvVariables, err))
	}

	var c configs.Configuration
	err := viper.Unmarshal(&c)
	c.Database.Password = os.Getenv("DB_PASSWORD")
	return c, err
}
