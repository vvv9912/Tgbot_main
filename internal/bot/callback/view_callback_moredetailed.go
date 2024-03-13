package callback

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"tgbot_main/internal/botkit"
	"tgbot_main/logger"
)

func ViewCallbackMoredetailed(p botkit.ProductsStorager, c botkit.CorzinaStorager) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {

		var Data botkit.BotCommand
		err := json.Unmarshal([]byte(update.CallbackQuery.Data), &Data)
		if err != nil {
			logger.Log.Error(logger.ErrorJsonUnmarshal, zap.Error(err))
			return err
		}

		var MsgAddCorzine AddCorzine
		err = json.Unmarshal([]byte(Data.Data), &MsgAddCorzine)
		if err != nil {
			logger.Log.Error(logger.ErrorJsonUnmarshal, zap.Error(err))
			return err
		}
		product, err := p.GetProductByArticle(ctx, MsgAddCorzine.Article)

		//Обновлять предыдущее сообщение (добавлять подробности)
		if product.Article == 0 {
			return err
		}

		text := fmt.Sprintf("Артикул: %d\nНазвание: %s\n%s\nЦена: %0.2fрублей\n", product.Article, product.Name, product.Description, product.Price)
		ms1 := tgbotapi.NewEditMessageText(botInfo.TgId, update.CallbackQuery.Message.MessageID, text)

		dataAddCorz := AddCorzine{
			Article: product.Article,
		}
		msgAddCorz, err := json.Marshal(dataAddCorz)
		if err != nil {
			logger.Log.Error(logger.ErrorJsonMarshal, zap.Error(err))
			return err
		}
		dataMsg := botkit.BotCommand{
			Cmd:  "/addCorzine",
			Data: string(msgAddCorz),
		}
		sss, err := json.Marshal(dataMsg)
		if err != nil {
			logger.Log.Error(logger.ErrorJsonMarshal, zap.Error(err))
			return err
		}

		var numericKeyboardInline = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Добавить в корзину", string(sss)),
			),
		)
		ms1.ReplyMarkup = &numericKeyboardInline
		_, err = bot.Send(ms1)
		if err != nil {
			logger.Log.Error(logger.ErrorSendMessage, zap.Error(err))
			return err
		}

		return nil
	}
}
