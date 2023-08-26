package bot_initer

import (
	"context"
	"fmt"
	"log"
	"telegram-bot/config"
	"telegram-bot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func New(cfg config.Config, userProvider userProvider, orderProvider orderProvider) *UseCase {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	return &UseCase{
		bot:          bot,
		cfg:          cfg,
		userProvider: userProvider,
	}
}

type UseCase struct {
	cfg           config.Config
	bot           *tgbotapi.BotAPI
	userProvider  userProvider
	orderProvider orderProvider
}

type userProvider interface {
	Get(ctx context.Context, user models.User) ([]models.User, error)
	Create(ctx context.Context, user models.User) (models.User, error)
	Update(ctx context.Context, user models.User) (models.User, error)
}

type orderProvider interface {
	Get(ctx context.Context, order models.Order) ([]models.Order, error)
	Create(ctx context.Context, order models.Order) (models.Order, error)
	Update(ctx context.Context, order models.Order) (models.Order, error)
}

type locationSaver interface {
	Save(ctx context.Context, user models.User, long, lat string) (models.User, error)
}

func (uc *UseCase) Execute(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := uc.bot.GetUpdatesChan(u)
	if err != nil {
		log.Println("Error While bot.GetUpdatesChan: ", err.Error())
		return err
	}

	for update := range updates {
		if update.Message == nil { // Пропустить любые не-сообщения
			continue
		}

		switch update.Message.Text {
		case "/start":
			var user models.User
			users, err := uc.userProvider.Get(ctx, models.User{
				TelegramUserID: update.Message.From.ID,
			})
			if err != nil {
				log.Println(err.Error())
				return err
			}
			leed := update.Message.From

			if len(users) == 0 {
				user, err = uc.userProvider.Create(ctx, models.User{
					PhoneNumber:    "",
					FirstName:      leed.FirstName,
					LastName:       leed.LastName,
					TelegramUserID: leed.ID,
					Status:         models.StatusUnverified,
				})
				if err != nil {
					log.Println(err.Error())
					return err
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Привет, %s %s", user.FirstName, user.LastName))
				uc.bot.Send(msg)
			}

			if user.PhoneNumber == "" {
				contactKeyboard := tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButtonContact("Отправить контакт"),
					),
				)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я бот для оказания услуг, чтобы начать отправьте свой контакт")
				msg.ReplyMarkup = contactKeyboard

				uc.bot.Send(msg)

				if update.Message.Contact == nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Напишите /start")
					uc.bot.Send(msg)
					continue
				}

				contact := update.Message.Contact

				users, err = uc.userProvider.Get(ctx, models.User{
					TelegramUserID: update.Message.From.ID,
				})
				if err != nil {
					log.Println(err.Error())
					return err
				}
				user := users[0]

				user.PhoneNumber = contact.PhoneNumber

				user, err = uc.userProvider.Update(ctx, user)
				if err != nil {
					log.Println(err.Error())
					return err
				}
			}

		case "/order":
			users, err := uc.userProvider.Get(ctx, models.User{
				TelegramUserID: update.Message.From.ID,
			})
			if err != nil {
				log.Println(err.Error())
				return err
			}
			if len(users) == 0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Напишите /start")
				uc.bot.Send(msg)
				continue
			}

			user := users[0]

			orders, err := uc.orderProvider.Get(ctx, models.Order{
				UserID: user.ID,
			})
			if err != nil {
				log.Println(err.Error())
				return err
			}

			fmt.Println(orders)

		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Напишите /start")
			uc.bot.Send(msg)
		}
	}
	return nil
}

func (uc *UseCase) Location(user models.User, update tgbotapi.Update) error {
	locationKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonLocation("Send Location"),
		),
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Отлично, %s, теперь отправьте свою локацию", user.FirstName))
	msg.ReplyMarkup = locationKeyboard
	uc.bot.Send(msg)

	if update.Message.Location == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Напишите /start")
		uc.bot.Send(msg)
	}

	return nil
}

// func Contact(update tgbotapi.Update) error {
// 	contact := tgbotapi.NewReplyKeyboard(
// 		tgbotapi.NewKeyboardButtonRow(
// 			tgbotapi.NewKeyboardButtonContact("Send Contact"),
// 		),
// 	)
// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome! Please send your Contact.")
// 	msg.ReplyMarkup = contact
// 	bot.Send(msg)
// 	if update.Message.Contact != nil {
// 		contact := update.Message.Contact.PhoneNumber
// 		_, err := strg.Order().Update(context.Background(), &models.OrderUpdate{
// 			Id:    order.Id,
// 			Name:  order.Name,
// 			Phone: contact,
// 			Lat:   order.Lat,
// 			Long:  order.Long,
// 		})
// 		if err != nil {
// 			log.Println(err.Error())
// 			return nil, err
// 		}
// 		order.Phone = contact
// 	}
// 	return order, nil
// }

// func Address(update tgbotapi.Update) error {
// 	address := tgbotapi.NewMessage(update.Message.Chat.ID, "Напишите польностью свой Адрес Например:\n Улица: Себзар, Дом:00, Квартира:00, Этаж:00 !")
// 	bot.Send(address)
// 	if update.Message != nil {
// 		address := update.Message.Text
// 		_, err := strg.Order().Update(context.Background(), &models.OrderUpdate{
// 			Id:      order.Id,
// 			Name:    order.Name,
// 			Phone:   order.Phone,
// 			Lat:     order.Lat,
// 			Long:    order.Long,
// 			Address: address,
// 		})
// 		if err != nil {
// 			log.Println(err.Error())
// 			return nil, err
// 		}
// 		order.Address = address
// 	}
// 	return order, nil
// }
