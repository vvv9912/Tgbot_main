package callback

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"tgbot_main/internal/botkit"
	"tgbot_main/logger"
)

func ViewCallbackUcatalog(s botkit.ProductsStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		var Data botkit.BotCommand
		err := json.Unmarshal([]byte(update.CallbackQuery.Data), &Data)
		if err != nil {
			logger.Log.Error(logger.ErrorJsonUnmarshal, zap.Error(err))
			return err
		}
		//Тут запрос в БД и поиск артикулей по названию каталога
		sProducts, err := s.GetProductsByCatalogIsAvailable(ctx, Data.Data)
		if err != nil {
			logger.Log.Error("Get products by catalog", zap.Error(err))
			return err
		}
		if len(sProducts) == 0 {
			msg := tgbotapi.NewMessage(update.FromChat().ID, "Товары отсутствуют!")
			if _, err := bot.Send(msg); err != nil {
				logger.Log.Error(logger.ErrorSendMessage, zap.Error(err))
				return err
			}
		}
		//
		for i := range sProducts {

			text := fmt.Sprintf("Артикул: %d\nНазвание: %s\nПодходит для: \nЦена: %0.2f рублей\n", sProducts[i].Article, sProducts[i].Name, sProducts[i].Price)

			if len(sProducts[i].PhotoUrl) != 0 {

				//todo photo

				var mediaGroup []interface{}

				for k := range sProducts[i].PhotoUrl {
					photoFileBytes := tgbotapi.FileBytes{
						Name:  "photo1.jpg",
						Bytes: sProducts[i].PhotoUrl[k],
					}
					photo1 := tgbotapi.NewInputMediaPhoto(photoFileBytes)
					if k == 0 {
						photo1.Caption = sProducts[i].Name //описание строго под одной фото
					}

					mediaGroup = append(mediaGroup, photo1)
				}

				msg := tgbotapi.NewMediaGroup(update.CallbackQuery.From.ID, mediaGroup)
				bot.Send(msg)

				var dataAddCorz AddCorzine
				dataAddCorz.Article = sProducts[i].Article
				msgAddCorz, err := json.Marshal(dataAddCorz)
				if err != nil {
					logger.Log.Error(logger.ErrorJsonUnmarshal, zap.Error(err))
					return err
				}

				data := botkit.BotCommand{Cmd: "/moredetailed",
					Data: string(msgAddCorz)}
				podrobnee, err := json.Marshal(data)
				dataMsg := botkit.BotCommand{
					Cmd:  "/addCorzine",
					Data: string(msgAddCorz), //не над
				}

				sss, err := json.Marshal(dataMsg)
				if err != nil {
					logger.Log.Error(logger.ErrorJsonUnmarshal, zap.Error(err))
				}

				var numericKeyboardInline = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Подробнее", string(podrobnee)),
						tgbotapi.NewInlineKeyboardButtonData("Добавить в корзину", string(sss)),
					),
				)
				msg2 := tgbotapi.NewMessage(update.CallbackQuery.From.ID, text)
				msg2.ReplyMarkup = numericKeyboardInline

				bot.Send(msg2) //todo

			} else { //если нет фото
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)

				dataAddCorz := AddCorzine{
					Article: sProducts[i].Article,
				}
				msgAddCorz, err := json.Marshal(dataAddCorz)
				if err != nil {
					logger.Log.Error(logger.ErrorJsonMarshal, zap.Error(err))
					return err
				}

				dataMsg := botkit.BotCommand{
					Cmd:  "/addCorzine",
					Data: string(msgAddCorz),
				}

				data := botkit.BotCommand{Cmd: "/moredetailed",
					Data: string(msgAddCorz)}

				podrobnee, err := json.Marshal(data)
				if err != nil {
					logger.Log.Error(logger.ErrorJsonMarshal, zap.Error(err))
					return err
				}

				sss, err := json.Marshal(dataMsg)
				if err != nil {
					logger.Log.Error(logger.ErrorJsonMarshal, zap.Error(err))
					return err
				}

				var numericKeyboardInline = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Подробнее", string(podrobnee)),
						tgbotapi.NewInlineKeyboardButtonData("Добавить в корзину", string(sss)),
					),
				)

				msg.ReplyMarkup = numericKeyboardInline
				_, err = bot.Send(msg) //todo
				if err != nil {
					logger.Log.Error(logger.ErrorSendMessage, zap.Error(err))
					return err
				}

			}
		}

		bot.Send(tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID, Text: "Добавлено!"})

		return nil
	}
}
