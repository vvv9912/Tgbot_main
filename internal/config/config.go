package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfighcl"
	"log"
	"sync"
)

type Config struct {
	DatabaseDSN      string `hcl:"database_dsn" env:"DATABASE_DSN" `
	TelegramBotToken string `hcl:"telegram_bot_token,omitempty" env:"TELEGRAM_BOT_TOKEN"`
	MyUrl            string `hcl:"my_url,omitempty" env:"MY_URL"`
	Cert             string `hcl:"cert,omitempty" env:"CERT"`
	Key              string `hcl:"key,omitempty" env:"KEY"`
	Localhost        string `hcl:"localhost,omitempty" env:"LOCALHOST"`
	PathFileLog      string `hcl:"PathFileLog,omitempty" env:"PathFileLog"`
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	once.Do(func() {
		loader := aconfig.LoaderFor(&cfg, aconfig.Config{
			EnvPrefix: "NFB",
			Files:     []string{"./config.hcl", "./config.local.hcl"},
			FileDecoders: map[string]aconfig.FileDecoder{
				".hcl": aconfighcl.New(),
			},
		})

		if err := loader.Load(); err != nil {
			log.Printf("[ERROR] failed to load config: %v", err)
		}
	})

	return cfg
}
