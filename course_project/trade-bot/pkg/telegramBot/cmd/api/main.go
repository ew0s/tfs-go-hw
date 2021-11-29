package main

import (
	"fmt"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"trade-bot/configs"
	"trade-bot/pkg/client/app"
	"trade-bot/pkg/client/service"
	"trade-bot/pkg/telegramBot"
)

var (
	ErrReadConfig                = errors.New("read config")
	ErrUnableToLoadEnvVariables  = errors.New("unable to load enviroment variable")
	ErrUnableToCreateTelegramBot = errors.New("unable to create telegram bot")
	ErrUnableToCreateClient      = errors.New("unable to create client")
	ErrSetupBot                  = errors.New("setup bot")
)

var (
	telegramAPIToken string
	webhookURL       string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(fmt.Errorf("%s: %s", ErrUnableToLoadEnvVariables, err))
	}

	telegramAPIToken = os.Getenv("TELEGRAM_APITOKEN")
	webhookURL = os.Getenv("WEBHOOK_URL")
}

func setBot() (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(telegramAPIToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrSetupBot, err)
	}

	log.Infof("Authorized on account: %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookURL))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrSetupBot, err)
	}

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("server stopped: ", err)
		}
	}()

	return bot, nil
}

func main() {
	config, err := initConfig()
	if err != nil {
		log.Fatalf("init config: %s", err)
	}

	bot, err := setBot()
	if err != nil {
		log.Fatalf("%s: %s", ErrUnableToCreateTelegramBot, err)
	}

	client, err := app.NewClient(config.Client)
	if err != nil {
		log.Fatalf("%s: %s", ErrUnableToCreateClient, err)
	}
	s := service.NewService(client)

	botman := telegramBot.NewBotMan(bot, s)
	botman.ServeTelegram()
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

	var c configs.Configuration
	err := viper.Unmarshal(&c)
	return c, err
}
