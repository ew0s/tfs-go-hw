package telegramBot

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"trade-bot/pkg/client/models"
	"trade-bot/pkg/client/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

var (
	ErrCouldNotSendMessage            = errors.New("could not send message")
	ErrExitFromSignUpInput            = errors.New("exited from sign up input")
	ErrExitFromSignInInput            = errors.New("exited from sign in input")
	ErrExitFromSendOrderInput         = errors.New("exited from send order input")
	ErrUnableToReadFromUpdatesChannel = errors.New("unable to read from updates channel")
	ErrUserAlreadyLoggedIn            = errors.New("user already logged in")
)

const (
	startCommand             = "/start"
	helpCommand              = "/help"
	signUpCommand            = "/sign_up"
	exitFromSignUpCommand    = "/exit_from_sign_up"
	signInCommand            = "/sign_in"
	exitFromSignInCommand    = "/exit_from_sign_in"
	sendOrderCommand         = "/send_order"
	exitFromSendOrderCommand = "/exit_from_send_order"
	logoutCommand            = "/logout"
)

type BotMan struct {
	bot              *tgbotapi.BotAPI
	tradeBotServices *service.Service
	usersJWT         map[string]string
}

func NewBotMan(bot *tgbotapi.BotAPI, tradeBotServices *service.Service) *BotMan {
	return &BotMan{bot: bot, tradeBotServices: tradeBotServices, usersJWT: map[string]string{}}
}

func (b *BotMan) ServeTelegram() {
	updates := b.bot.ListenForWebhook("/")

	for update := range updates {
		if update.Message != nil {
			if !update.Message.IsCommand() {
				continue
			}

			log.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)

			chatID := update.Message.Chat.ID

			switch update.Message.Text {
			case startCommand:
				message := tgbotapi.NewMessage(chatID, startMessage)
				message.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				b.sendMessage(chatID, message)

			case helpCommand:
				message := tgbotapi.NewMessage(chatID, helpMessage)
				message.ReplyToMessageID = update.Message.MessageID
				b.sendMessage(chatID, message)

			case signUpCommand:
				if _, ok := b.usersJWT[update.Message.From.UserName]; ok {
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", signUpErrMessage, ErrUserAlreadyLoggedIn))
					b.sendMessage(chatID, errMessage)
				}

				message := tgbotapi.NewMessage(chatID, signUpMessage)
				message.ReplyToMessageID = update.Message.MessageID
				b.sendMessage(chatID, message)

				if err := b.executeSignUp(updates); err != nil {
					log.Warn(err)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", signUpErrMessage, err.Error()))
					b.sendMessage(chatID, errMessage)
				} else {
					successMessage := tgbotapi.NewMessage(chatID, signUpSuccessMessage)
					b.sendMessage(chatID, successMessage)
				}

			case signInCommand:
				message := tgbotapi.NewMessage(chatID, signInMessage)
				message.ReplyToMessageID = update.Message.MessageID
				b.sendMessage(chatID, message)

				token, err := b.executeSignIn(updates)
				if err != nil {
					log.Warn(err)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", signInErrMessage, err.Error()))
					b.sendMessage(chatID, errMessage)
					continue
				}

				b.usersJWT[update.Message.From.UserName] = token
				successMessage := tgbotapi.NewMessage(chatID, signInSuccessMessgae)
				b.sendMessage(chatID, successMessage)

			case logoutCommand:
				token, err := b.userIdentity(update.Message.From.UserName)
				if err != nil {
					log.Warn(err)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", logoutErrMessage, err.Error()))
					b.sendMessage(chatID, errMessage)
					continue
				}

				_, err = b.tradeBotServices.Authorization.Logout(models.LogoutInput{JWTToken: token})
				if err != nil {
					log.Warn(err)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", logoutErrMessage, err.Error()))
					b.sendMessage(chatID, errMessage)
					continue
				}

				delete(b.usersJWT, update.Message.From.UserName)
				successMessage := tgbotapi.NewMessage(chatID, logoutSuccessMessgae)
				b.sendMessage(chatID, successMessage)

			case sendOrderCommand:
				token, err := b.userIdentity(update.Message.From.UserName)
				if err != nil {
					log.Warn(err)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", sendOrderErrMessage, err.Error()))
					b.sendMessage(chatID, errMessage)
					continue
				}

				message := tgbotapi.NewMessage(chatID, sendOrderMessage)
				b.sendMessage(chatID, message)

				if err := b.executeSendOrder(updates, token); err != nil {
					log.Warn(err)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", sendOrderErrMessage, err.Error()))
					b.sendMessage(chatID, errMessage)
					continue
				}

				successMessage := tgbotapi.NewMessage(chatID, sendOrderSuccessMessage)
				b.sendMessage(chatID, successMessage)

			default:
				message := tgbotapi.NewMessage(chatID, invalidCommandMessage)
				b.sendMessage(chatID, message)
			}
		}
	}
}

func (b *BotMan) userIdentity(username string) (string, error) {
	val, ok := b.usersJWT[username]
	if !ok {
		return "", fmt.Errorf("user not logged in")
	}
	return val, nil
}

func (b *BotMan) sendMessage(chatID int64, message tgbotapi.MessageConfig) {
	if _, err := b.bot.Send(message); err != nil {
		log.Warnf("%s: [chatID] - %d", ErrCouldNotSendMessage, chatID)
	}
}

func (b *BotMan) executeSendOrder(updates tgbotapi.UpdatesChannel, token string) error {
	input, err := b.getSendOrderInput(updates)
	if err != nil {
		return err
	}
	input.JwtToken = token

	_, err = b.tradeBotServices.OrdersManagee.SendOrder(input)
	return err
}

func (b *BotMan) getSendOrderInput(updates tgbotapi.UpdatesChannel) (models.SendOrderInput, error) {
	for update := range updates {
		if update.Message == nil {
			return models.SendOrderInput{}, nil
		}

		switch update.Message.Text {
		case exitFromSendOrderCommand:
			return models.SendOrderInput{}, ErrExitFromSendOrderInput
		default:
			inputValues := strings.FieldsFunc(update.Message.Text, split)
			if len(inputValues) != 3 {
				return models.SendOrderInput{}, fmt.Errorf("invalid count of arguments")
			}
			if inputValues[1] != "buy" && inputValues[1] != "sell" {
				return models.SendOrderInput{}, fmt.Errorf("invalid send order Side argument")
			}
			amount, err := strconv.Atoi(inputValues[2])
			if err != nil {
				return models.SendOrderInput{}, fmt.Errorf("invalid send order Size argument")
			}
			return models.SendOrderInput{
				OrderType: "mkt",
				Symbol:    inputValues[0],
				Side:      inputValues[1],
				Size:      amount,
			}, nil
		}
	}

	return models.SendOrderInput{}, ErrUnableToReadFromUpdatesChannel
}

func (b *BotMan) executeSignIn(updates tgbotapi.UpdatesChannel) (string, error) {
	input, err := b.getSignInInput(updates)
	if err != nil {
		return "", err
	}

	resp, err := b.tradeBotServices.Authorization.SignIn(input)
	if err != nil {
		return "", err
	}
	return resp.AcessToken, nil
}

func (b *BotMan) getSignInInput(updates tgbotapi.UpdatesChannel) (models.SignInInput, error) {
	for update := range updates {
		if update.Message == nil {
			return models.SignInInput{}, nil
		}

		switch update.Message.Text {
		case exitFromSignInCommand:
			return models.SignInInput{}, ErrExitFromSignInInput
		default:
			inputValues := strings.FieldsFunc(update.Message.Text, split)
			if len(inputValues) != 2 {
				return models.SignInInput{}, fmt.Errorf("invalid count of arguments")
			}
			return models.SignInInput{
				Username: inputValues[0],
				Password: inputValues[1],
			}, nil
		}
	}

	return models.SignInInput{}, ErrUnableToReadFromUpdatesChannel
}

func (b *BotMan) executeSignUp(updates tgbotapi.UpdatesChannel) error {
	input, err := b.getSignUpInput(updates)
	if err != nil {
		return err
	}

	_, err = b.tradeBotServices.Authorization.SignUp(input)
	return err
}

func (b *BotMan) getSignUpInput(updates tgbotapi.UpdatesChannel) (models.SignUpInput, error) {
	for update := range updates {
		if update.Message == nil {
			return models.SignUpInput{}, nil
		}

		switch update.Message.Text {
		case exitFromSignUpCommand:
			return models.SignUpInput{}, ErrExitFromSignUpInput
		default:
			inputValues := strings.FieldsFunc(update.Message.Text, split)
			if len(inputValues) != 5 {
				return models.SignUpInput{}, fmt.Errorf("invalid count of arguments")
			}
			return models.SignUpInput{
				Name:          inputValues[0],
				Username:      inputValues[1],
				Password:      inputValues[2],
				PublicAPIKey:  inputValues[3],
				PrivateAPIKey: inputValues[4],
			}, nil
		}
	}

	return models.SignUpInput{}, ErrUnableToReadFromUpdatesChannel
}

func split(r rune) bool {
	return r == ' ' || r == '\n'
}
