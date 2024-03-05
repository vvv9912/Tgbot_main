package text

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbotv2/internal/botkit"
)

func ViewTextCorzine(c botkit.CorzinaStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {
		corz, err := c.GetShopCartDetailByTgId(ctx, botInfo.TgId)
		if err != nil {
			return err
		}
		if len(corz) == 0 {
			msg := tgbotapi.NewMessage(botInfo.TgId, "Корзина пуста!")
			bot.Send(msg)
			return nil
		}
		outmsg := ""
		var numKeyInline tgbotapi.InlineKeyboardMarkup
		numKeyInline.InlineKeyboard = make([][]tgbotapi.InlineKeyboardButton, 2)

		var numKeyInline2 tgbotapi.InlineKeyboardMarkup
		numKeyInline2.InlineKeyboard = make([][]tgbotapi.InlineKeyboardButton, len(corz))
		for i := range corz {
			outmsg += fmt.Sprintf("%d товар.\nАртикул товара:%d\nНазвание:%s\n Количество:%d\n", i+1, corz[i].Article, corz[i].Name, corz[i].Quantity)
			for k := range numKeyInline.InlineKeyboard {

				switch k {
				case 0:
					var structData botkit.BotCommand
					structData.Cmd = "/deleteshoppcart"
					data, err := json.Marshal(structData)
					stringData := string(data)
					if err != nil {
						return err
					}
					numKeyInline.InlineKeyboard[k] = make([]tgbotapi.InlineKeyboardButton, 1) //В поле создаем еще поле
					numKeyInline.InlineKeyboard[k][0].Text = "Очистить корзину"
					//data = "/deleteshoppcart" //надо передавать команду +id что удалить(?) //todo
					numKeyInline.InlineKeyboard[k][0].CallbackData = &stringData
				case 1:
					var structData botkit.BotCommand
					structData.Cmd = "/addorder"
					structData.Data = fmt.Sprintf("%d", corz[i].Article)
					data, err := json.Marshal(structData)
					if err != nil {
						return err
					}
					stringData := string(data)
					numKeyInline.InlineKeyboard[k] = make([]tgbotapi.InlineKeyboardButton, 1) //В поле создаем еще поле
					numKeyInline.InlineKeyboard[k][0].Text = "Оформить заказ"
					//data = "/placeanorder" //надо передавать команду +id что удалить(?) //todo
					numKeyInline.InlineKeyboard[k][0].CallbackData = &stringData
				}
			}

		}
		msg1 := tgbotapi.NewMessage(update.Message.Chat.ID, outmsg)
		msg1.ReplyMarkup = numKeyInline
		msgsend, err := bot.Send(msg1)
		if err != nil {
			return err
		}
		_ = msgsend // зачем
		for k := range numKeyInline2.InlineKeyboard {
			var structData botkit.BotCommand
			structData.Cmd = "/deleteshoppcartposition"
			structData.Data = fmt.Sprintf("%d", corz[k].Article)
			data, err := json.Marshal(structData)
			if err != nil {
				return err
			}
			stringData := string(data)
			numKeyInline2.InlineKeyboard[k] = make([]tgbotapi.InlineKeyboardButton, 1)
			numKeyInline2.InlineKeyboard[k][0].Text = corz[k].Name //
			//data = fmt.Sprintf("/deleteshopposition\narticle:%d\nmsg:%d", corz[i].Article) //зачем то передавал от bot.send -> сюда |msg.id
			numKeyInline2.InlineKeyboard[k][0].CallbackData = &stringData

		}
		msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, "Удалить из корзины позиции:")
		msg2.ReplyMarkup = numKeyInline2
		bot.Send(msg2)
		return nil
	}
}
