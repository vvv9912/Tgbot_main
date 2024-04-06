package bot

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbot_main/internal/bot/callback"
	callback3 "tgbot_main/internal/bot/callback/order"
	"tgbot_main/internal/bot/cmd"
	"tgbot_main/internal/bot/middleware"
	"tgbot_main/internal/bot/text"
	"tgbot_main/internal/botkit"
	"tgbot_main/logger"
)

type MyBot struct {
	botApi *tgbotapi.BotAPI
	api    *botkit.Bot
}

func NewMyBot(botApi *tgbotapi.BotAPI, sUsers botkit.UsersStorager, sProducts botkit.ProductsStorager, sCorzina botkit.CorzinaStorager, sOrder botkit.OrderStorager) *MyBot {
	mybot := botkit.New(botApi)

	mw := middleware.NewMiddleware(sUsers)

	mybot.RegisterCmdView("start", mw.MwUsersOnly(cmd.ViewCmdStart(cmd.ViewCmdButton())))
	mybot.RegisterCmdView("button", cmd.ViewCmdButton()) //check user
	mybot.RegisterCmdView("adminbutton", cmd.ViewCmdAdminButton())
	//
	mybot.RegisterTextView("Каталог", mw.MwUsersOnly(text.ViewTextCatalog(sProducts)))
	mybot.RegisterTextView("Корзина", mw.MwUsersOnly(text.ViewTextCorzine(sCorzina)))
	mybot.RegisterTextView("Мои заказы", mw.MwUsersOnly(text.ViewTextOrder(sOrder)))
	//
	mybot.RegisterCallbackView("/ucatalog", callback.ViewCallbackUcatalog(sProducts))                              //todo mw check user
	mybot.RegisterCallbackView("/addCorzine", callback.ViewCallbackAddcorzine(sCorzina))                           //todo mw check user
	mybot.RegisterCallbackView("/moredetailed", callback.ViewCallbackMoredetailed(sProducts, sCorzina))            //todo mw check user
	mybot.RegisterCallbackView("/deleteshoppcart", callback.ViewCallbackdeleteshoppcart(sCorzina))                 //todo mw check user
	mybot.RegisterCallbackView("/deleteshoppcartposition", callback.ViewCallbackdeleteshoppcartposition(sCorzina)) //todo mw check user

	mybot.RegisterCallbackView("/addorder", callback3.ViewCallbackAddorder())                                   //todo mw check user
	mybot.RegisterCallbackView("/addOrderSamovivoz", callback3.ViewCallbackAddOrderSamovivoz(sOrder, sCorzina)) //todo mw check user
	mybot.RegisterCallbackView("/addSendOrderModerator", callback3.ViewCallbackWriteModer(sOrder, sCorzina))    //todo mw check user
	mybot.RegisterCallbackView("/orderCancel", callback3.ViewCallbackOrderCancel(sOrder, sCorzina))             //todo mw check user
	//
	mybot.RegisterTextView("FAQ", text.ViewTextFaq())

	return &MyBot{botApi: botApi, api: mybot}
}

func (b *MyBot) Start(ctx context.Context, token, cert, url, key, localhost string) error {
	webhook, err := b.api.Webhook(token, cert, url)
	if err != nil {
		logger.Log.AddOriginalError(err).CustomError("failed to create webhook", nil)
		return err
	}

	resp, err := b.botApi.Request(webhook)
	if err != nil {
		logger.Log.AddOriginalError(err).CustomError("failed to request webhook", map[string]interface{}{"resp": resp})
		return err
	}

	info, err := b.botApi.GetWebhookInfo()
	if err != nil {
		logger.Log.AddOriginalError(err).CustomError("failed to get webhook info", nil)
		return err
	}
	if info.LastErrorDate != 0 {
		logger.Log.CustomError("Telegram callback failed", map[string]interface{}{"info": info})
	}

	if err := b.api.Run(ctx, token, cert, key, localhost); err != nil {
		if !errors.Is(err, context.Canceled) {
			logger.Log.AddOriginalError(err).CustomError("failed to start bot", nil)
			return err
		}
		logger.Log.CustomInfo("bot stopped", nil)
	}
	return nil

}
