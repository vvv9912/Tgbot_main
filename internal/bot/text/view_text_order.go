package text

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"tgbot_main/internal/botkit"
	"tgbot_main/internal/model"
	"tgbot_main/logger"
)

func ViewTextOrder(o botkit.OrderStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		order, err := o.GetOrdersByTgID(ctx, botInfo.TgId)
		if err != nil {
			logger.Log.AddOriginalError(err).CustomError("GetOrdersByTgID", map[string]interface{}{"tg_id": botInfo.TgId})
			return err
		}

		if len(order) == 0 {
			msg := tgbotapi.NewMessage(botInfo.TgId, "Заказов нет!")
			bot.Send(msg)
			return nil
		}

		outmsg := ""

		for i := range order {
			var Corz []model.OrderCorz
			err := json.Unmarshal([]byte(order[i].Order), &Corz)
			if err != nil {
				logger.Log.Error(logger.ErrorJsonUnmarshal, zap.Error(err))
				return err
			}

			//Проверка на корз = 0
			outmsg += fmt.Sprintf("➖➖➖➖➖➖➖\nID Заказа №%d\n", order[i].ID)
			for k := range Corz {
				outmsg += fmt.Sprintf("%d товар.\nАртикул товара: %d\nНазвание: %s\nКоличество: %d\nЦена: %0.2f\n", k+1, Corz[k].Article, Corz[k].Name, Corz[k].Quantity, Corz[k].Price)
			}
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, outmsg)

		_, err = bot.Send(msg)
		if err != nil {
			logger.Log.Error(logger.ErrorSendMessage, zap.Error(err))
			return err
		}

		return nil
	}
}
