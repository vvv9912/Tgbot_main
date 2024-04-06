package cmd

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbot_main/internal/botkit"
)

func ViewCmdAdminButton() botkit.ViewFunc {
	var numericKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Каталог"),
			tgbotapi.NewKeyboardButton("Мои заказы"),
			tgbotapi.NewKeyboardButton("Корзина"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("FAQ"),
			tgbotapi.NewKeyboardButton("HELP"),
			//tgbotapi.NewKeyboardButton("6"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Настройки админа"),
			//tgbotapi.NewKeyboardButton("Мои заказы"),
			//tgbotapi.NewKeyboardButton("Корзина"),
		),
	)
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Меню добавлено")
		msg.ReplyMarkup = numericKeyboard
		if _, err := bot.Send(msg); err != nil {
			return err
		}
		return nil
	}
}
