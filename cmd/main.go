package main

import (
	"context"
	"fmt"
	"log"

	"telegram-bot/config"
	"telegram-bot/models"
	"telegram-bot/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	cfg *config.Config
}

func NewStore(cfg *config.Config) *Store {

	return &Store{
		cfg: cfg,
	}
}

func Location(bot *tgbotapi.BotAPI, update tgbotapi.Update, strg storage.StorageI, order *models.Order) (*models.Order, error) {
	if update.Message.Text == "/start" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome!")
		bot.Send(msg)
		// Create a custom keyboard with a "Send Location" button
		location := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButtonLocation("Send Location"),
			),
		)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome! Please send your location.")
		msg.ReplyMarkup = location
		bot.Send(msg)
		id, err := strg.Order().Create(context.Background(), &models.OrderCreate{
			Name:  update.Message.Chat.FirstName,
			Phone: update.Message.Chat.FirstName,
		})
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		order.Id = id
	} else if update.Message.Location != nil {
		latitude := fmt.Sprintf("%f", update.Message.Location.Latitude)
		longitude := fmt.Sprintf("%f", update.Message.Location.Latitude)
		_, err := strg.Order().Update(context.Background(), &models.OrderUpdate{
			Id:   order.Id,
			Name: update.Message.Chat.FirstName,
			Lat:  latitude,
			Long: longitude,
		})
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		order.Name = update.Message.Chat.FirstName
		order.Lat = latitude
		order.Long = longitude
	}
	return order, nil
}

func Contact(bot *tgbotapi.BotAPI, update tgbotapi.Update, strg storage.StorageI, order *models.Order) (*models.Order, error) {
	contact := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact("Send Contact"),
		),
	)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome! Please send your Contact.")
	msg.ReplyMarkup = contact
	bot.Send(msg)
	if update.Message.Contact != nil {
		contact := update.Message.Contact.PhoneNumber
		_, err := strg.Order().Update(context.Background(), &models.OrderUpdate{
			Id:    order.Id,
			Name:  order.Name,
			Phone: contact,
			Lat:   order.Lat,
			Long:  order.Long,
		})
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		order.Phone = contact
	}
	return order, nil
}

func Address(bot *tgbotapi.BotAPI, update tgbotapi.Update, strg storage.StorageI, order *models.Order) (*models.Order, error) {
	address := tgbotapi.NewMessage(update.Message.Chat.ID, "Напишите польностью свой Адрес Например:\n Улица: Себзар, Дом:00, Квартира:00, Этаж:00 !")
	bot.Send(address)
	if update.Message != nil {
		address := update.Message.Text
		_, err := strg.Order().Update(context.Background(), &models.OrderUpdate{
			Id:      order.Id,
			Name:    order.Name,
			Phone:   order.Phone,
			Lat:     order.Lat,
			Long:    order.Long,
			Address: address,
		})
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		order.Address = address
	}
	return order, nil
}

func main() {
	cfg := config.Load()

	conn, err := sqlx.ConnectContext(context.Background(), "postgres", fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresDatabase,
		cfg.PostgresPassword,
		cfg.PostgresPort,
	))
	if err != nil {
		panic(err)
	}

	conn.SetMaxIdleConns(cfg.PostgresMaxConnection)

}
