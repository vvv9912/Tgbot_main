package callback

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vvv9912/sddb"
	"go.uber.org/zap"
	"strconv"
	"tgbotv2/internal/botkit"
	"tgbotv2/logger"
)

func ViewCallbackdeleteshoppcartposition(c botkit.CorzinaStorager) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		var Data botkit.BotCommand
		err := json.Unmarshal([]byte(update.CallbackQuery.Data), &Data)
		if err != nil {
			logger.Log.Error(logger.ErrorJsonMarshal, zap.Error(err))
			return err
		}
		article, err := strconv.Atoi(Data.Data)
		if err != nil {
			logger.Log.Error(logger.ErrorAtoi, zap.Error(err))
			return err
		}

		corz, err := c.GetShopCartByTgIdAndArticle(ctx, botInfo.TgId, article)
		if err != nil {
			logger.Log.Error("Get shoppcart", zap.Error(err))
			return err
		}
		if corz == (sddb.ShopCart{}) {
			err := fmt.Errorf("Такой позиции нет в корзине, корзина пуста")
			logger.Log.CustomWarn(err.Error(), map[string]interface{}{
				"article": article,
				"tg_id":   botInfo.TgId,
			})
			return err
		}
		if corz.Quantity == 1 {
			err := c.DeleteShopCartByTgIdAndArticle(ctx, botInfo.TgId, article)
			if err != nil {
				logger.Log.Error("Delete shoppcart", zap.Error(err))
				return err
			}
		} else {
			newQuan := corz.Quantity - 1
			err := c.UpdateShopCartByTgId(ctx, botInfo.TgId, article, newQuan)
			if err != nil {
				logger.Log.Error("Update shoppcart", zap.Error(err))
				return err
			}
		}
		//err := c.DeleteShopCartByTgId(ctx, botInfo.TgId)
		//if err != nil {
		//	fmt.Println("DeleteShopCartByTgId :", err)
		//	return err
		//}
		answ := tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID, Text: "Позиция удалена!"}
		bot.Send(answ)

		return nil
	}
}
