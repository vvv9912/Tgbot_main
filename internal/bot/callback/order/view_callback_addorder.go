package callback_order

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tgbotv2/internal/botkit"
)

func ViewCallbackAddorder() botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		dataMsgSam := botkit.BotCommand{
			Cmd: "/addSendOrderModerator",
			//Data: "",
		}
		msgSam, err := json.Marshal(dataMsgSam)
		if err != nil {
			log.Println("") //todo
			return err
		}
		dataMsgCancel := botkit.BotCommand{
			Cmd:  "/orderCancel",
			Data: "",
		}
		msgCancel, err := json.Marshal(dataMsgCancel)
		if err != nil {
			log.Println("") //todo
			return err
		}
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите способ доставки")
		var numericKeyboardInline = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Оформить заказ, о доставке с вами свяжутся(?)", string(msgSam)), //оставляю эту + добавить отмену
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Отмена", string(msgCancel)),
			),
		)
		msg.ReplyMarkup = numericKeyboardInline
		bot.Send(msg)

		return nil
	}
}
