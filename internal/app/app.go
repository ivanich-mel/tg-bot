// Package app wires application dependencies and runtime lifecycle.
package app

import (
	"context"
	"database/sql"
	"fmt"

	"melnik/telegram-bot/config"
	db "melnik/telegram-bot/db/sqlc"
	telegram "melnik/telegram-bot/pkg/telegramm"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

func Run(ctx context.Context) error {
	conf, err := config.Init()
	if err != nil {
		return fmt.Errorf("error initialize config: %w", err)
	}
	botApi, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		return fmt.Errorf("could not initiate a new bot: %w", err)
	}
	botApi.Debug = true

	conn, err := sql.Open(conf.DB.Driver, conf.DB.Url)
	if err != nil {
		return fmt.Errorf("cannot connect to db: %w", err)
	}
	defer conn.Close()
	if err := conn.PingContext(ctx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}
	DB := db.New(conn)
	bot := telegram.NewBot(botApi, DB)
	if err := bot.Start(ctx); err != nil {
		return fmt.Errorf("error start bot: %w", err)
	}
	return nil
}
