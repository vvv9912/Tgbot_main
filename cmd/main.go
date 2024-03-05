package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vvv9912/sddb"
	"tgbotv2/internal/bot"
	"tgbotv2/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tgbotv2/internal/config"
)

func run() error {
	flagLogLevel := "info" //todo config
	if err := logger.Initialize2(flagLogLevel, config.Get().PathFileLog); err != nil {
		return err
	}
	return nil
}

func main() {

	token := config.Get().TelegramBotToken
	key := config.Get().Key
	cert := config.Get().Cert
	localhost := config.Get().Localhost
	url := config.Get().MyUrl

	log.Println("token:", token, "\n", "key:", key, "\n", "cert:", cert, "\n", "localhost:", localhost, "\n", "url:", url)

	if err := run(); err != nil {
		log.Panic(err)
		return
	}

	logger.Log.CustomError("failed to create bot", nil)
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Log.AddOriginalError(err).CustomError("failed to create bot", nil)
		return
	}
	db, err := sqlx.Connect("postgres", config.Get().DatabaseDSN)
	if err != nil {
		logger.Log.AddOriginalError(err).CustomError("failed to connect database", nil)
		return
	}
	defer db.Close()

	sProducts := sddb.NewProductsPostgresStorage(db)
	sUsers := sddb.NewStorageUser(db)
	sCorzina := sddb.NewStorageCorzina(db)
	sOrder := sddb.NewStorageOrder(db)

	mybot := bot.NewMyBot(botAPI, sUsers, sProducts, sCorzina, sOrder)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := mybot.Start(ctx, token, cert, url, key, localhost); err != nil {
		logger.Log.AddOriginalError(err).CustomError("failed to bot", nil)
		return
	}

}
