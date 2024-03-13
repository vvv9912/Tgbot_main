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
)

func ViewCallbackWriteModer(o botkit.OrderStorager, c botkit.CorzinaStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		corz, err := c.GetShopCartDetailByTgId(ctx, botInfo.TgId)
		if err != nil {
			logger.Log.Error("Get shopcart", zap.Error(err))
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
		err = o.AddOrder(ctx, sddb.Orders{
			TgID:        botInfo.TgId,
			UserName:    botInfo.UserName,
			FirstName:   botInfo.FirstName,
			LastName:    botInfo.LastName,
			StatusOrder: sddb.StatusOrderNew,
			Pvz:         "{}",
			Order:       string(msgcorz),
			TypeDostavka: constant.D_SAMOVIVOZ,
		})
		if err != nil {
			logger.Log.Error("Add order", zap.Error(err))
			return err
		}
		//todo delete shopcart
		msg := tgbotapi.NewMessage(botInfo.TgId, "Заказ добавлен! В течение дня с вами свяжется менеджер!\nДля уточнения заказа вы можете написать @Bykova_Dina")
		bot.Send(msg)
		return nil
	}
}
