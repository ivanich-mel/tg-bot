package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	db "melnik/telegram-bot/db/sqlc"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	startCmd           = "start"
	helpCmd            = "help"
	listCategoriesCmd  = "list"
	createCategoryCmd  = "create"
	addNameState       = "add_name"
	addBalanceState    = "add_balance"
	changeBalanceState = "change_balance"
)

func (b *tgBot) handleCommand(msg *tgbotapi.Message) error {
	switch msg.Command() {
	case startCmd:
		return b.handleStartCmd(msg)
	case helpCmd:
		return b.handleHelpCmd(msg)
	case createCategoryCmd:
		return b.handleCreateCategoryCmd(msg)
	case listCategoriesCmd:
		return b.handleListCategoriesCmd(msg)
	default:
		return b.sendMessage(msg.Chat.ID, "Такой команды не существует. Попробуй еще раз =)")
	}
}

func (b *tgBot) handleStartCmd(msg *tgbotapi.Message) error {
	return b.sendMessage(msg.Chat.ID, startMsg)
}

func (b *tgBot) handleHelpCmd(msg *tgbotapi.Message) error {
	return b.sendMessage(msg.Chat.ID, helpMsg)
}
func (b *tgBot) handleCreateCategoryCmd(msg *tgbotapi.Message) error {
	b.mu.Lock()
	chatID := msg.Chat.ID
	_, exists := b.userStates[chatID]
	if !exists {
		b.userStates[chatID] = &UserState{State: addNameState}
		b.mu.Unlock()
		b.sendMessage(chatID, "Введи название категории")
		return nil
	}
	return nil
}
func (b *tgBot) handleListCategoriesCmd(msg *tgbotapi.Message) error {
	message := tgbotapi.NewMessage(msg.Chat.ID, "Список Категорий: ")

	message.ReplyMarkup = b.getListCategories(100, 0)

	if _, err := b.bot.Send(message); err != nil {
		log.Printf("failed to send message: %v", err)
	}
	return nil
}

func (b *tgBot) handleMessage(msg *tgbotapi.Message) error {
	chatID := msg.Chat.ID

	b.mu.Lock()
	defer b.mu.Unlock()

	state, exists := b.userStates[chatID]
	if !exists {
		b.sendMessage(chatID, startMsg)
		return nil
	}
	switch state.State {
	case addNameState:
		state.CreateCategoryState.Name = msg.Text
		state.State = addBalanceState
		b.sendMessage(chatID, "Введи сумму, которую планиурется потратить за месяц")
	case addBalanceState:
		balance, err := strconv.ParseFloat(msg.Text, 32)
		if err != nil {
			b.sendMessage(chatID, "Пожалуйста, введите числовое значение для баланса.")
			return nil
		}
		state.CreateCategoryState.Balance = balance
		c, err := b.db.CreateCategory(context.Background(), state.CreateCategoryState)
		if err != nil {
			log.Printf("Ошибка при создании категории: %v", err)
			b.sendMessage(chatID, "Произошла ошибка при сохранении категории. Попробуйте позже.")
			return nil
		}
		delete(b.userStates, chatID)

		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Категория: '%s' успешно создана", c.Name))
		msg.ReplyMarkup = b.getListCategories(100, 0)
		if _, err := b.bot.Send(msg); err != nil {
			log.Printf("failed to send message: %v", err)
		}
	case changeBalanceState:
		subtractive, err := strconv.ParseFloat(msg.Text, 32)
		if err != nil {
			b.sendMessage(chatID, "Пожалуйста, введите числовое значение для баланса.")
			return nil
		}
		c, err := b.db.GetCategory(context.Background(), state.UpdateeCategoryState.ID)
		if err != nil {
			log.Printf("can't get category: %v", err)
			return nil
		}
		state.UpdateeCategoryState.Balance = c.Balance - subtractive
		err = b.db.UpdateCategory(context.Background(), db.UpdateCategoryParams{
			ID:      c.ID,
			Name:    c.Name,
			Balance: state.UpdateeCategoryState.Balance,
		})
		if err != nil {
			log.Printf("Ошибка при создании категории: %v", err)
			b.sendMessage(chatID, "Произошла ошибка при сохранении категории. Попробуйте позже.")
			return nil
		}
		delete(b.userStates, chatID)

		msg := tgbotapi.NewMessage(chatID, "Categories")
		msg.ReplyMarkup = b.getListCategories(100, 0)
		if _, err := b.bot.Send(msg); err != nil {
			log.Printf("failed to send message: %v", err)
		}
	}
	return nil
}

func (b *tgBot) handleCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	defer b.mu.Unlock()
	data := callbackQuery.Data
	chatID := callbackQuery.Message.Chat.ID
	b.mu.Lock()
	switch {
	case strings.HasPrefix(data, "change_balance"):
		categoryID, _ := parseCategoryID(data, "change_balance_")
		state, exists := b.userStates[chatID]
		if !exists {
			state = &UserState{}
			b.userStates[chatID] = state
		}
		state.State = changeBalanceState
		state.UpdateeCategoryState.ID = categoryID

		b.sendMessage(chatID, "Введи сумму чека")
		return nil
	}
	return nil
}
func (b *tgBot) getListCategories(limit int32, offset int32) tgbotapi.InlineKeyboardMarkup {
	categories, err := b.db.ListCategories(context.Background(), db.ListCategoriesParams{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		log.Printf("Ошибка при получении категорий: %v", err)
		return tgbotapi.InlineKeyboardMarkup{}
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, category := range categories {
		changeNameCallback := fmt.Sprintf("change_name_%d", category.ID)
		changeBalanceCallback := fmt.Sprintf("change_balance_%d", category.ID)

		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(category.Name, changeNameCallback),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%.2f", category.Balance), changeBalanceCallback),
		)
		rows = append(rows, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
func (b *tgBot) sendMessage(chatID int64, text string) (err error) {
	if chatID == 0 {
		return errors.New("chatID cannot be zero")
	}
	if text == "" {
		return errors.New("text cannot be empty string")
	}
	msg := tgbotapi.NewMessage(chatID, text)

	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("failed to send message: %v", err)
	}
	return nil
}
func parseCategoryID(data string, prefix string) (int64, error) {
	if !strings.HasPrefix(data, prefix) {
		return 0, fmt.Errorf("invalid prefix")
	}
	id, err := strconv.Atoi(strings.TrimPrefix(data, prefix))
	if err != nil {
		return 0, fmt.Errorf("failed to parse category ID: %w", err)
	}
	return int64(id), nil
}
