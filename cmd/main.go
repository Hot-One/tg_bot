package main

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"telegram-bot/config"
	"telegram-bot/models"
	"telegram-bot/storage"
	"telegram-bot/storage/postgres"
)

type Store struct {
	cfg  *config.Config
	strg storage.StorageI
}

func NewStore(cfg *config.Config, strg storage.StorageI) *Store {
	return &Store{
		cfg:  cfg,
		strg: strg,
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

	// Here Laod Configs
	config := config.Load()

	// Here We Connected To Database
	pgconn, err := postgres.NewConnectionPostgres(&config)
	if err != nil {
		panic("postgres no connection: " + err.Error())
	}

	// Here We Connected Store
	strg := NewStore(&config, pgconn)

	bot, err := tgbotapi.NewBotAPI(config.BotToken)
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
	var Order *models.Order
	for update := range updates {
		if update.Message == nil { // Пропустить любые не-сообщения
			continue
		}
		if update.Message.Text == "/start" {
			order, err := Location(bot, update, strg.strg, Order)
			if err != nil {
				log.Println("Error while Get Location: ", err.Error())
				return
			}
			Order.Id = order.Id
			Order.Name = order.Name
			Order.Lat = order.Lat
			Order.Long = order.Long
		} else if update.Message.Location != nil {
			contact, err := Contact(bot, update, strg.strg, Order)
			if err != nil {
				log.Println("Error while Get Contact: ", err.Error())
				return
			}
			Order.Phone = contact.Phone
		} else if update.Message != nil {
			fmt.Println("AAAAAAAAAAAAAAAAAAAA: ", Order)
			address, err := Address(bot, update, strg.strg, Order)
			if err != nil {
				log.Println("Error while Get Addres: ", err.Error())
				return
			}
			Order.Address = address.Address
		}
	}
}
