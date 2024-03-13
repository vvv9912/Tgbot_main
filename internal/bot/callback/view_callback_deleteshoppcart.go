package callback

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbot_main/internal/botkit"
	"tgbot_main/logger"
)

func ViewCallbackdeleteshoppcart(c botkit.CorzinaStorager) botkit.ViewFunc {
	//view_callback_deleteshoppcartposition.go
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {

		err := c.DeleteShopCartByTgId(ctx, botInfo.TgId)
		if err != nil {
			logger.Log.AddOriginalError(err).CustomError("Delete Shop Cart By TgId", map[string]interface{}{
				"tg_id": botInfo.TgId,
			})
			return err
		}
		answ := tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID, Text: "Корзина удалена!"}
		bot.Send(answ)
		return nil
	}
}
