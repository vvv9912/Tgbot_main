package text

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tgbot_main/internal/botkit"
)

func ViewTextCatalog(s botkit.ProductsStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {

		catalog, err := s.GetCatalogNamesIsAvailable(ctx)
		if err != nil {
			log.Printf("[ERROR] Get catalog from ProductsStorage %v", err)
			return err

		}
		if len(catalog) == 0 {
			if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Каталог пуст!")); err != nil {
				return err
			}
			return nil

		}

		var numKeyInline tgbotapi.InlineKeyboardMarkup
		numKeyInline.InlineKeyboard = make([][]tgbotapi.InlineKeyboardButton, len(catalog))
		for i := range numKeyInline.InlineKeyboard {
			//var data string
			numKeyInline.InlineKeyboard[i] = make([]tgbotapi.InlineKeyboardButton, 1)
			numKeyInline.InlineKeyboard[i][0].Text = catalog[i] //

			cmd := botkit.BotCommand{Cmd: "/ucatalog",
				Data: catalog[i]}
			//data = fmt.Sprintf("/ucatalog\ncategory:%s", catalog[i]) //надо передавать команду +id что удалить(?) //todo передать сразу json
			data, err := json.Marshal(cmd)
			if err != nil {
				log.Println("") //todo
				//return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
				//
				//}
			}
			CallbackData := string(data)
			numKeyInline.InlineKeyboard[i][0].CallbackData = &CallbackData
		}

		msg := tgbotapi.NewMessage(update.FromChat().ID, "Выберите каталог")
		msg.ReplyMarkup = numKeyInline
		if _, err := bot.Send(msg); err != nil {
			return err
		}
		return nil
	}
}
