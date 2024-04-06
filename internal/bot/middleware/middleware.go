package middleware

import (
	"context"
	"database/sql"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vvv9912/sddb"
	"go.uber.org/zap"
	"log"
	"tgbot_main/internal/bot/constant"
	"tgbot_main/logger"

	"tgbot_main/internal/botkit"
)

type Middleware struct {
	UserStorage botkit.UsersStorager
}

func NewMiddleware(userStorage botkit.UsersStorager) *Middleware {
	return &Middleware{UserStorage: userStorage}
}

// сюда кэш можно
func (m *Middleware) MwAdminOnly(next botkit.ViewFunc) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {

		next(ctx, bot, update, botInfo)
		//проверка на админа
		return nil
	}
}

func (m *Middleware) MwUsersOnly(next botkit.ViewFunc) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, botInfo botkit.BotInfo) error {

		uStatus, uState, err := m.UserStorage.GetStatusUserByTgID(ctx, update.FromChat().ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = nil
			} else {
				logger.Log.Error("Get status user tg id", zap.Error(err))
				return err
			}
		}
		botInfo.IdStatus = uStatus
		botInfo.IdState = uState

		switch uStatus {
		case constant.UNoUser:
			err = m.UserStorage.AddUser(ctx, sddb.Users{
				TgID:       botInfo.TgId,
				StatusUser: constant.UUser,
				StateUser:  constant.E_STATE_NOTHING,
			})
			if err != nil {
				log.Println(err)
				return err
			}
			botInfo.IdState = constant.UUser
			botInfo.IdStatus = constant.E_STATE_NOTHING
		case constant.UAdmin:
			//
		case constant.UUser:
			//
		}

		err = next(ctx, bot, update, botInfo)

		if err != nil {
			return err
		}
		return nil
	}
}
