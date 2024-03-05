package text

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbotv2/internal/botkit"
)

func ViewTextFaq() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		var numericKeyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButtonWebApp("Приложенька ", tgbotapi.WebAppInfo{URL: "https://google.com"}),
			),
		)
		msg := tgbotapi.NewMessage(botInfo.TgId, "лял")
		msg.ReplyMarkup = numericKeyboard
		bot.Send(msg) //todo

		return nil
	}
}
