package callback_order

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vvv9912/sddb"
	"go.uber.org/zap"
	"tgbotv2/internal/bot/constant"
	"tgbotv2/internal/botkit"
	"tgbotv2/internal/model"
	"tgbotv2/logger"
	"time"
)

func ViewCallbackAddOrderSamovivoz(o botkit.OrderStorager, c botkit.CorzinaStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		corz, err := c.GetShopCartDetailByTgId(ctx, botInfo.TgId)
		if err != nil {
			logger.Log.AddOriginalError(err).CustomError("GetShopCartDetailByTgId", map[string]interface{}{"tg_id": botInfo.TgId})
			return err
		}

		if len(corz) == 0 {
			msg := tgbotapi.NewMessage(botInfo.TgId, "Корзина пуста!")
			bot.Send(msg)
			return nil
		}

		corzForOrder := make([]model.OrderCorz, len(corz))

		for i := range corz {
			corzForOrder[i].Article = corz[i].Article
			corzForOrder[i].Quantity = corz[i].Quantity
			corzForOrder[i].Price = corz[i].Price
			corzForOrder[i].ID = corz[i].ID
			corzForOrder[i].Name = corz[i].Name
		}

		msgcorz, err := json.Marshal(corzForOrder)
		if err != nil {
			logger.Log.Error(logger.ErrorJsonMarshal, zap.Error(err))
			return err
		}

		order := sddb.Orders{

			TgID:        botInfo.TgId,
			UserName:    botInfo.UserName,
			FirstName:   botInfo.FirstName,
			LastName:    botInfo.LastName,
			StatusOrder: 0,
			Pvz:         "{}",
			Order:       string(msgcorz),

			CreatedAt:    time.Now().UTC().Add(3 * time.Hour),
			TypeDostavka: constant.D_SAMOVIVOZ,
		}

		err = o.AddOrder(ctx, order)
		if err != nil {
			logger.Log.AddOriginalError(err).CustomError("AddOrder", map[string]interface{}{"order": order})
			return err
		}

		msg := tgbotapi.NewMessage(botInfo.TgId, "Заказ добавлен!")
		bot.Send(msg)

		answ := tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID, Text: "Заказ добавлен!"}

		_, err = bot.Send(answ)
		if err != nil {
			logger.Log.Error(logger.ErrorSendMessage, zap.Error(err))
			return err
		}

		return nil
	}
}
