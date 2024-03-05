package callback

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vvv9912/sddb"
	"go.uber.org/zap"
	"tgbotv2/internal/botkit"
	"tgbotv2/logger"
	"time"
)

func ViewCallbackAddcorzine(c botkit.CorzinaStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		var Data botkit.BotCommand

		err := json.Unmarshal([]byte(update.CallbackQuery.Data), &Data)
		if err != nil {
			logger.Log.Error(logger.ErrorJsonUnmarshal, zap.Error(err))
			return err
		}

		var MsgAddCorzine AddCorzine

		err = json.Unmarshal([]byte(Data.Data), &MsgAddCorzine)
		if err != nil {
			logger.Log.Error(logger.ErrorJsonUnmarshal, zap.Error(err))
			return err
		}

		//Добавление в БД
		corz, err := c.GetShopCartByTgIdAndArticle(ctx, botInfo.TgId, MsgAddCorzine.Article)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = nil
				err = c.AddShopCart(ctx, sddb.ShopCart{
					TgId:      botInfo.TgId,
					Article:   MsgAddCorzine.Article,
					Quantity:  1,
					CreatedAt: time.Now().UTC().Add(3 * time.Hour),
				})
				if err != nil {
					logger.Log.Error("AddShopCart", zap.Error(err))
					return err
				}

				//send message 'add corzine' todo
				msg := tgbotapi.NewMessage(botInfo.TgId, "Добавлено в корзину!")
				_, err = bot.Send(msg)
				if err != nil {
					logger.Log.Error("AddShopCart", zap.Error(err))
					return err
				}

				answ := tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID, Text: "Добавлено в корзину!"}
				bot.Send(answ)

				return nil
			}
			logger.Log.Error("GetShopCartByTgIdAndArticle", zap.Error(err))

			return err
		}

		if corz == (sddb.ShopCart{}) { //Проверка на всякий случай //если корзина пуста
			err = c.AddShopCart(ctx, sddb.ShopCart{
				TgId:      botInfo.TgId,
				Article:   MsgAddCorzine.Article,
				Quantity:  1,
				CreatedAt: time.Now().UTC().Add(3 * time.Hour),
			})
			if err != nil {
				logger.Log.Error("AddShopCart", zap.Error(err))
				return err
			}
			msg := tgbotapi.NewMessage(botInfo.TgId, "Добавлено в корзину!")

			_, err = bot.Send(msg)
			if err != nil {
				logger.Log.Error(logger.ErrorSendMessage, zap.Error(err))
				return err
			}

			answ := tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID, Text: "Добавлено в корзину!"}
			_, err = bot.Send(answ)
			if err != nil {
				logger.Log.Error(logger.ErrorSendMessage, zap.Error(err))
				return err
			}
			return nil
		}

		corz.Quantity++
		corz.CreatedAt = time.Now().UTC().Add(3 * time.Hour)

		err = c.UpdateShopCartByTgId(ctx, botInfo.TgId, corz.Article, corz.Quantity) //todo обновлять время
		if err != nil {
			logger.Log.Error("UpdateShopCartByTgId", zap.Error(err))
			return err
		}

		msg := tgbotapi.NewMessage(botInfo.TgId, "Добавлено в корзину!")
		_, err = bot.Send(msg)
		if err != nil {
			logger.Log.Error(logger.ErrorSendMessage, zap.Error(err))
			return err
		}

		answ := tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID, Text: "Добавлено в корзину!"}

		_, err = bot.Send(answ)
		if err != nil {
			logger.Log.Error(logger.ErrorSendMessage, zap.Error(err))
			return err
		}

		//send message 'add corzine' todo
		return nil
	}
}
