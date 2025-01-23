package main

import (
	"database/sql"
	"log"
	"melnik/telegram-bot/config"
	db "melnik/telegram-bot/db/sqlc"
	"melnik/telegram-bot/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.Init()
	if err != nil {
		log.Fatal("error initialize config: ", err)
	}
	botApi, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Panic(err)
	}
	botApi.Debug = true

	conn, err := sql.Open(conf.DB.Driver, conf.DB.Url)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	DB := db.New(conn)
	bot := telegram.NewBot(botApi, DB)
	if err := bot.Start(); err != nil {
		log.Fatal("error start bot: ", err)
	}
}
