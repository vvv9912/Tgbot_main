package callback_order

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"tgbotv2/internal/botkit"
	"tgbotv2/logger"
)

func ViewCallbackOrderCancel(o botkit.OrderStorager, c botkit.CorzinaStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		msg := tgbotapi.NewDeleteMessage(botInfo.TgId, update.CallbackQuery.Message.MessageID)
		bot.Send(msg)
		answ := tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID, Text: "Заказ Отменен!"}
		_, err := bot.Send(answ)
		if err != nil {
			logger.Log.Error("Delete message callback", zap.Error(err))
			return err
		}
		return nil
	}
}
