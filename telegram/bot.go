package telegram

import (
	db "melnik/telegram-bot/db/sqlc"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type tgBot struct {
	bot        *tgbotapi.BotAPI
	userStates map[int64]*UserState
	db         *db.Queries
	mu         sync.Mutex
}

type UserState struct {
	State                string
	CreateCategoryState  db.CreateCategoryParams
	UpdateeCategoryState db.UpdateCategoryParams
}

func NewBot(bot *tgbotapi.BotAPI, db *db.Queries) *tgBot {
	return &tgBot{
		bot:        bot,
		db:         db,
		userStates: make(map[int64]*UserState),
	}
}

func (b *tgBot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			b.handleCallback(update.CallbackQuery)
		}
		if update.Message != nil {
			if update.Message.IsCommand() {
				b.handleCommand(update.Message)
			} else {
				b.handleMessage(update.Message)
			}
		}
	}

	return nil
}
