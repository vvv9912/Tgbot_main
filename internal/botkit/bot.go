package botkit

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"tgbot_main/logger"
	"time"
)

//логика работы с ботом
import (
	"log"
	"runtime/debug"
)

type BotCommand struct {
	Cmd  string `json:"cmd,omitempty"`
	Data string `json:"data,omitempty"` //в дата упоквано в зависимости от сообщения модель
}
type Bot struct {
	api           *tgbotapi.BotAPI
	cmdViews      map[string]ViewFunc // комманды тг бота
	textViews     map[string]ViewFunc
	callbackViews map[string]ViewFunc
}
type BotInfo struct {
	TgId      int64
	IdStatus  int
	IdState   int
	UserName  string
	FirstName string
	LastName  string
}

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo BotInfo) error

func New(api *tgbotapi.BotAPI) *Bot {
	return &Bot{api: api}
}

func (b *Bot) RegisterTextView(cmd string, view ViewFunc) {
	if b.textViews == nil {
		b.textViews = make(map[string]ViewFunc)
	}
	b.textViews[cmd] = view
}
func (b *Bot) RegisterCallbackView(cmd string, view ViewFunc) {
	if b.callbackViews == nil {
		b.callbackViews = make(map[string]ViewFunc)
	}
	b.callbackViews[cmd] = view
}
func (b *Bot) RegisterCmdView(cmd string, view ViewFunc) {
	if b.cmdViews == nil {
		b.cmdViews = make(map[string]ViewFunc)
	}
	b.cmdViews[cmd] = view
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() { //перехват паники
		if p := recover(); p != nil {
			log.Printf("[ERROR] panic recovered: %v\n%s", p, string(debug.Stack()))
		}
	}()

	var view ViewFunc
	if update.Message == nil {
		if update.CallbackQuery != nil {

			var Data BotCommand

			err := json.Unmarshal([]byte(update.CallbackQuery.Data), &Data)
			if err != nil {
				log.Printf("[ERROR] Json преобразование callback %v", err)
				return
			}
			clbck := Data.Cmd
			callbackView, ok := b.callbackViews[clbck]
			if !ok {
				return
			}
			view = callbackView
		} else if update.InlineQuery != nil {
			//тут  чтото
		} else {

			return
		}
	} else {
		if update.Message.IsCommand() {
			cmd := update.Message.Command()

			cmdView, ok := b.cmdViews[cmd]
			if !ok {
				return
			}

			view = cmdView

		} else if update.Message.Document != nil {

			cmd := update.Message.Caption
			cmdView, ok := b.cmdViews[cmd]
			if !ok {
				return
			}

			view = cmdView
		} else {
			//Если текстовая команда
			text := update.Message.Text
			if text == "" {
				return
			}
			textView, ok := b.textViews[text]
			if !ok {
				return
			}
			view = textView

		}
	}
	var botInfo BotInfo
	botInfo.TgId = update.FromChat().ID
	botInfo.UserName = update.FromChat().UserName
	botInfo.FirstName = update.FromChat().FirstName
	botInfo.LastName = update.FromChat().LastName
	if err := view(ctx, b.api, update, botInfo); err != nil {
		log.Printf("[ERROR] failed to handle update: %v", err)

		b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintln(err)))
		if _, err := b.api.Send(
			tgbotapi.NewMessage(update.Message.Chat.ID, "internal error"),
		); err != nil {
			log.Printf("[ERROR] failed to send message: %v", err)
		}
	}
}

func (b *Bot) Run(ctx context.Context, token string, cert string, key string, localhost string) error {

	listen := "/bot" + token
	log.Println("listen:", listen)
	updates := b.api.ListenForWebhook(listen)
	log.Println("serverStated")
	go func() {
		err := http.ListenAndServe(localhost, nil)
		//err := http.ListenAndServeTLS(localhost, cert, key, nil)
		if err != nil {
			log.Println("ListenAndServeTLS:", err)
			panic(err)
			return
		}
	}()

	for {
		select {
		case update := <-updates:
			go func(update tgbotapi.Update) {
				updateCtx, updateCancel := context.WithTimeout(ctx, 60*time.Second)
				b.handleUpdate(updateCtx, update)
				//select c таймаутом. То что запрос долго обрабатывается
				defer updateCancel()
			}(update)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) Webhook(botToken string, certPath string, url string) (*tgbotapi.WebhookConfig, error) {

	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		logger.Log.Error("Read certPath", zap.Error(err))
		return nil, err
	}
	RequestFileData := tgbotapi.FileBytes{
		Name:  "cert",
		Bytes: cert,
	}

	urlbot := fmt.Sprintf("%s:8443/bot%s", url, botToken)
	log.Println(urlbot)

	webcookconf, err := tgbotapi.NewWebhookWithCert(urlbot, RequestFileData)
	if err != nil {
		logger.Log.Error("NewWebhookWithCert", zap.Error(err))
		return nil, err
	}

	return &webcookconf, nil
}
