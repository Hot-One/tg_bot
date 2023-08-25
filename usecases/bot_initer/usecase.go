package bot_initer

import (
	"context"
	"log"
	"telegram-bot/config"
	"telegram-bot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func New(cfg config.Config, userGetter userGetter) *UseCase {
	return &UseCase{
		cfg: cfg,
	}
}

type UseCase struct {
	cfg config.Config
}

type userGetter interface {
	Get(ctx context.Context, user models.User) (models.User, error)
}

func (uc *UseCase) Execute(ctx context.Context) (*tgbotapi.BotAPI, error) {

	bot, err := tgbotapi.NewBotAPI(uc.cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println("Error While bot.GetUpdatesChan: ", err.Error())
		return
	}

	for update := range updates {
		if update.Message == nil { // Пропустить любые не-сообщения
			continue
		}

		switch update.Message.Text {
		case "/start":

			order, err := Location(bot, update, strg.strg, Order)
			if err != nil {
				log.Println("Error while Get Location: ", err.Error())
				return
			}

			Order.Id = order.Id
			Order.Name = order.Name
			Order.Lat = order.Lat
			Order.Long = order.Long

		case "/location":
			contact, err := Contact(bot, update, strg.strg, Order)
			if err != nil {
				log.Println("Error while Get Contact: ", err.Error())
				return
			}
			Order.Phone = contact.Phone

			address, err := Address(bot, update, strg.strg, Order)
			if err != nil {
				log.Println("Error while Get Addres: ", err.Error())
				return
			}
			Order.Address = address.Address
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Напишите /start")
			bot.Send(msg)
		}
	}

	return nil
}
