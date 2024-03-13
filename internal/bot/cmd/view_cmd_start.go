package cmd

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbotv2/internal/botkit"
)

func ViewCmdStart(next botkit.ViewFunc) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, `Привет, Я бот для онлайн покупки @dinasty_nature. Выберите в каталоге нужный товар, добавьте в коризну и оформите заказ:)`)); err != nil {
			return err
		}
		next(ctx, bot, update, botInfo)
		return nil
	}
}
