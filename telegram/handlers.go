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
	startCmd          = "start"
	createCategoryCmd = "create_category"
	listBalancesCmd   = "list_balances"
	listLimitsCmd     = "list_limits"
	updateBalancesCmd = "update_balances"

	updateBalancesCallBack       = "callback_update_balances"
	listCategoriesCallBack       = "callback_list_categories"
	listLimitsCallBack           = "callback_list_limits"
	createCategoryCallBack       = "callback_create_category"
	deleteCategoryCallBack       = "callback_delete_category"
	renameCategoryCallback       = "callback_rename_category"
	changeCategoryLimitsCallback = "callback_change_limit_category"

	addNameState                = "state_add_name"
	renameCategoryState         = "state_rename_category"
	addBalanceState             = "state_add_balance"
	changeBalanceState          = "state_change_balance"
	changeCategoryNameState     = "state_change_name"
	changeCategoriesLimitsState = "state_change_limits"

	categoryOptionsAction       = "category_options_action"
	deleteCategoryAction        = "action_delete_category"
	deleteCategoryApproveAction = "yes_action_delete_approve"
	deleteCategoryDeclineAction = "no_action_delete_decline"

	updateBalancesAction        = "action_update_balance"
	updateBalancesApproveAction = "yes_action_update_balance_approve"
	updateBalancesDeclineAction = "no_action_update_balance_decline"
)

func (b *tgBot) handleCommand(msg *tgbotapi.Message) error {
	chatID := msg.Chat.ID
	switch msg.Command() {
	case startCmd:
		return b.handleStartCmd(chatID)
	case createCategoryCmd:
		return b.handleCreateCetegoryCmd(chatID)
	case listBalancesCmd:
		return b.handleListBalancesCmd(chatID)
	case listLimitsCmd:
		return b.handleListLimitsCmd(chatID)
	case updateBalancesCmd:
		return b.handleUpdateBalancesCmd(chatID)
	default:
		return b.sendMessage(msg.Chat.ID, defaultCommandMsg)
	}
}

func (b *tgBot) handleStartCmd(chatID int64) error {
	var rows [][]tgbotapi.InlineKeyboardButton

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("➕ New category", createCategoryCallBack),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📋 List balances", listCategoriesCallBack),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⚠️  List limits", listLimitsCallBack),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔄 Update balances", updateBalancesAction),
	))

	message := tgbotapi.NewMessage(chatID, startMsg)
	message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	if _, err := b.bot.Send(message); err != nil {
		log.Printf("failed to send message: %v", err)
	}
	return nil
}
func (b *tgBot) handleCreateCetegoryCmd(chatID int64) error {
	b.handleCreateCategoryCallback(chatID)
	return nil
}
func (b *tgBot) handleListBalancesCmd(chatID int64) error {
	b.handleListCategoriesCallback(chatID)
	return nil
}
func (b *tgBot) handleListLimitsCmd(chatID int64) error {
	b.handleListLimitsCallback(chatID)
	return nil
}
func (b *tgBot) handleUpdateBalancesCmd(chatID int64) error {
	b.handleUpdateBalancesAction(chatID)
	return nil
}

func (b *tgBot) handleCallback(callbackQuery *tgbotapi.CallbackQuery) error {
	data := callbackQuery.Data
	chatID := callbackQuery.Message.Chat.ID

	switch {
	case strings.HasPrefix(data, changeBalanceState):
		categoryID, _ := parseCategoryID(data, changeBalanceState+"_")
		b.handleChangeBalanceState(chatID, categoryID)
	case strings.HasPrefix(data, changeCategoryNameState):
		categoryID, _ := parseCategoryID(data, changeBalanceState+"_")
		b.handleChangeBalanceState(chatID, categoryID)
	case strings.HasPrefix(data, deleteCategoryAction):
		categoryID, _ := parseCategoryID(data, deleteCategoryAction+"_")
		b.handleDeleteCategotyAction(chatID, categoryID)
	case strings.HasPrefix(data, categoryOptionsAction):
		categoryID, _ := parseCategoryID(data, categoryOptionsAction+"_")
		b.handleCategoryOptionsAction(chatID, categoryID)
	case strings.HasPrefix(data, deleteCategoryApproveAction):
		categoryID, _ := parseCategoryID(data, deleteCategoryApproveAction+"_")
		b.handleDeleteCategoryApproveAction(chatID, categoryID)
	case strings.HasPrefix(data, renameCategoryCallback):
		categoryID, _ := parseCategoryID(data, renameCategoryCallback+"_")
		b.handleRenameCategoryCallback(chatID, categoryID)
	case strings.HasPrefix(data, changeCategoryLimitsCallback):
		categoryID, _ := parseCategoryID(data, changeCategoryLimitsCallback+"_")
		b.changeCategoryLimitsCallback(chatID, categoryID)
	case data == deleteCategoryDeclineAction:
		b.handleStartCmd(chatID)
	case data == createCategoryCallBack:
		b.handleCreateCategoryCallback(chatID)
	case data == listCategoriesCallBack:
		b.handleListCategoriesCallback(chatID)
	case data == listLimitsCallBack:
		b.handleListLimitsCallback(chatID)
	case data == deleteCategoryCallBack:
		b.handleDeleteCategoryCallback(chatID)
	case data == updateBalancesAction:
		b.handleUpdateBalancesAction(chatID)
	case data == updateBalancesApproveAction:
		b.handleUpdateBalancesApproveAction(chatID)
	case data == updateBalancesDeclineAction:
		b.handleStartCmd(chatID)
	default:
		message := tgbotapi.NewMessage(chatID, defaultCallbackMsg)
		if _, err := b.bot.Send(message); err != nil {
			log.Printf("failed to send message: %v", err)
		}
		return nil
	}

	return nil
}

func (b *tgBot) handleCreateCategoryCallback(chatID int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, exists := b.userStates[chatID]
	if !exists {
		b.userStates[chatID] = &UserState{State: addNameState}
		b.sendMessage(chatID, createCategoryMsg)
		return nil
	}
	return nil
}
func (b *tgBot) handleRenameCategoryCallback(chatID int64, categoryID int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	state, exists := b.userStates[chatID]
	if !exists {
		state = &UserState{}
		b.userStates[chatID] = state
	}
	state.State = renameCategoryState
	state.UpdateeCategoryState.ID = categoryID

	b.sendMessage(chatID, renameCategoryMsg)
	return nil
}
func (b *tgBot) handleListCategoriesCallback(chatID int64) error {
	message := tgbotapi.NewMessage(chatID, listCategoriesBalanceMsg)
	message.ReplyMarkup = b.listCategoriesBalanceKeyboard(100, 0)
	if _, err := b.bot.Send(message); err != nil {
		log.Printf("failed to send message: %v", err)
	}
	return nil
}
func (b *tgBot) handleListLimitsCallback(chatID int64) error {
	message := tgbotapi.NewMessage(chatID, listCategoriesPermanentBalanceMsg)
	message.ReplyMarkup = b.listCategoriesLimitsKeyboard(100, 0)
	if _, err := b.bot.Send(message); err != nil {
		log.Printf("failed to send message: %v", err)
	}
	return nil
}
func (b *tgBot) handleDeleteCategoryCallback(chatID int64) error {
	message := tgbotapi.NewMessage(chatID, listCategoriesDeleteMsg)
	message.ReplyMarkup = b.listCategoriesDeleteKeyboard(100, 0)
	if _, err := b.bot.Send(message); err != nil {
		log.Printf("failed to send message: %v", err)
	}
	return nil
}
func (b *tgBot) handleDeleteCategoryApproveAction(chatID int64, categoryID int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	c, err := b.db.GetCategory(context.Background(), categoryID)
	if err != nil {
		log.Printf("Error retrieving category: %v", err)
		return err
	}
	err = b.db.DeleteCategory(context.Background(), c.ID)
	if err != nil {
		log.Printf("Error deleting category: %v", err)
		return err
	}
	message := tgbotapi.NewMessage(chatID, fmt.Sprintf(deleteCategorySuccessMsg, c.Name))
	message.ReplyMarkup = b.listCategoriesBalanceKeyboard(100, 0)
	if _, err := b.bot.Send(message); err != nil {
		log.Printf("failed to send message: %v", err)
	}
	return nil
}
func (b *tgBot) handleUpdateBalancesApproveAction(chatID int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	categories, err := b.db.ListCategories(context.Background(), db.ListCategoriesParams{
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		log.Printf("Error retrieving categories: %v", err)
		return err
	}

	for _, c := range categories {
		err := b.db.UpdateCategory(context.Background(), db.UpdateCategoryParams{
			ID:               c.ID,
			Name:             c.Name,
			Balance:          c.PermanentBalance,
			PermanentBalance: c.PermanentBalance,
		})
		if err != nil {
			log.Printf("Error updating category %s: %v", c.Name, err)
		}
	}

	message := tgbotapi.NewMessage(chatID, updateBalancesSuccessMsg)
	message.ReplyMarkup = b.listCategoriesBalanceKeyboard(100, 0)

	if _, err := b.bot.Send(message); err != nil {
		log.Printf("Error sending message: %v", err)
	}

	return nil
}
func (b *tgBot) handleChangeBalanceState(chatID int64, categoryID int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	state, exists := b.userStates[chatID]
	if !exists {
		state = &UserState{}
		b.userStates[chatID] = state
	}
	state.State = changeBalanceState
	state.UpdateeCategoryState.ID = categoryID

	b.sendMessage(chatID, enterReceiptMsg)
	return nil
}
func (b *tgBot) changeCategoryLimitsCallback(chatID int64, categoryID int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	state, exists := b.userStates[chatID]
	if !exists {
		state = &UserState{}
		b.userStates[chatID] = state
	}
	state.State = changeCategoriesLimitsState
	state.UpdateeCategoryState.ID = categoryID

	b.sendMessage(chatID, enterLimitMsg)
	return nil
}
func (b *tgBot) handleCategoryOptionsAction(chatID int64, categoryID int64) error {
	var rows [][]tgbotapi.InlineKeyboardButton

	b.mu.Lock()
	defer b.mu.Unlock()

	c, err := b.db.GetCategory(context.Background(), categoryID)
	if err != nil {
		log.Printf("failed to get category: %v", err)
		return nil
	}
	if c.ID == 0 {
		message := tgbotapi.NewMessage(chatID, categoryNotFoundMsg)
		if _, err := b.bot.Send(message); err != nil {
			log.Printf("failed to send message: %v", err)
		}
	}
	row := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✏️ Rename", fmt.Sprintf("%s_%d", renameCategoryCallback, categoryID)),
		tgbotapi.NewInlineKeyboardButtonData("❌ Delete", fmt.Sprintf("%s_%d", deleteCategoryAction, categoryID)),
	)

	rows = append(rows, row)
	message := tgbotapi.NewMessage(chatID, fmt.Sprintf(categorySelectedMsg, c.Name))
	message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	if _, err := b.bot.Send(message); err != nil {
		log.Printf("failed to send message: %v", err)
	}
	return nil
}

func (b *tgBot) handleDeleteCategotyAction(chatID int64, categoryID int64) error {
	var rows [][]tgbotapi.InlineKeyboardButton

	b.mu.Lock()
	defer b.mu.Unlock()
	row := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Yes", fmt.Sprintf("%s_%d", deleteCategoryApproveAction, categoryID)),
		tgbotapi.NewInlineKeyboardButtonData("No", deleteCategoryDeclineAction),
	)
	rows = append(rows, row)
	message := tgbotapi.NewMessage(chatID, deleteCategoryMsg)
	message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	if _, err := b.bot.Send(message); err != nil {
		log.Printf("failed to send message: %v", err)
	}
	return nil
}
func (b *tgBot) handleUpdateBalancesAction(chatID int64) error {
	var rows [][]tgbotapi.InlineKeyboardButton

	b.mu.Lock()
	defer b.mu.Unlock()
	row := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Yes", updateBalancesApproveAction),
		tgbotapi.NewInlineKeyboardButtonData("No", updateBalancesDeclineAction),
	)
	rows = append(rows, row)
	message := tgbotapi.NewMessage(chatID, updateBalancesMsg)
	message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
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
		b.sendMessage(chatID, enterAllowanceBalanceMsg)
	case addBalanceState:
		balance, err := strconv.ParseFloat(msg.Text, 32)
		if err != nil {
			b.sendMessage(chatID, incorrectNumberBalanceMsg)
			return nil
		}
		state.CreateCategoryState.Balance = balance
		c, err := b.db.CreateCategory(context.Background(), db.CreateCategoryParams{
			Name:             state.CreateCategoryState.Name,
			Balance:          state.CreateCategoryState.Balance,
			PermanentBalance: state.CreateCategoryState.Balance,
		})
		if err != nil {
			log.Printf("failed to create category: %v", err)
			b.sendMessage(chatID, unknownErrorMsg)
			return nil
		}
		delete(b.userStates, chatID)

		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(categoryCreatedMsg, c.Name))
		msg.ReplyMarkup = b.listCategoriesBalanceKeyboard(100, 0)
		if _, err := b.bot.Send(msg); err != nil {
			log.Printf("failed to send message: %v", err)
		}
	case changeBalanceState:
		subtractive, err := strconv.ParseFloat(msg.Text, 32)
		if err != nil {
			b.sendMessage(chatID, incorrectNumberBalanceMsg)
			return nil
		}
		c, err := b.db.GetCategory(context.Background(), state.UpdateeCategoryState.ID)
		if err != nil {
			log.Printf("failed to get category: %v", err)
			return nil
		}
		state.UpdateeCategoryState.Balance = c.Balance - subtractive
		err = b.db.UpdateCategory(context.Background(), db.UpdateCategoryParams{
			ID:               c.ID,
			Name:             c.Name,
			Balance:          state.UpdateeCategoryState.Balance,
			PermanentBalance: c.PermanentBalance,
		})
		if err != nil {
			log.Printf("failed to create category: %v", err)
			b.sendMessage(chatID, unknownErrorMsg)
			return nil
		}
		delete(b.userStates, chatID)

		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(updateBalanceSuccessMsg, c.Name))
		msg.ReplyMarkup = b.listCategoriesBalanceKeyboard(100, 0)
		if _, err := b.bot.Send(msg); err != nil {
			log.Printf("failed to send message: %v", err)
		}
	case changeCategoriesLimitsState:
		newLimit, err := strconv.ParseFloat(msg.Text, 32)
		if err != nil {
			b.sendMessage(chatID, incorrectNumberBalanceMsg)
			return nil
		}
		c, err := b.db.GetCategory(context.Background(), state.UpdateeCategoryState.ID)
		if err != nil {
			log.Printf("failed to get category: %v", err)
			return nil
		}
		err = b.db.UpdateCategory(context.Background(), db.UpdateCategoryParams{
			ID:               c.ID,
			Name:             c.Name,
			Balance:          c.Balance,
			PermanentBalance: newLimit,
		})
		if err != nil {
			log.Printf("failed to update category: %v", err)
			b.sendMessage(chatID, unknownErrorMsg)
			return nil
		}
		delete(b.userStates, chatID)

		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(updateLimitSuccessMsg, c.Name))
		msg.ReplyMarkup = b.listCategoriesLimitsKeyboard(100, 0)
		if _, err := b.bot.Send(msg); err != nil {
			log.Printf("failed to send message: %v", err)
		}
	case renameCategoryState:
		categoryName := msg.Text
		fmt.Println(categoryName)
		c, err := b.db.GetCategory(context.Background(), state.UpdateeCategoryState.ID)
		if err != nil {
			log.Printf("failed to get category: %v", err)
			return nil
		}
		err = b.db.UpdateCategory(context.Background(), db.UpdateCategoryParams{
			ID:               c.ID,
			Name:             categoryName,
			Balance:          c.Balance,
			PermanentBalance: c.PermanentBalance,
		})
		if err != nil {
			log.Printf("failed to create category: %v", err)
			b.sendMessage(chatID, unknownErrorMsg)
			return nil
		}
		delete(b.userStates, chatID)

		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(categoryRenamedMsg, c.Name, categoryName))
		msg.ReplyMarkup = b.listCategoriesBalanceKeyboard(100, 0)
		if _, err := b.bot.Send(msg); err != nil {
			log.Printf("failed to send message: %v", err)
		}
	}
	return nil
}

func (b *tgBot) listCategoriesBalanceKeyboard(limit int32, offset int32) tgbotapi.InlineKeyboardMarkup {
	categories, err := b.db.ListCategories(context.Background(), db.ListCategoriesParams{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		log.Printf("failed to get categories: %v", err)
		return tgbotapi.InlineKeyboardMarkup{}
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	var total float64 = 0
	for _, category := range categories {
		categoryOptionsAction := fmt.Sprintf("%s_%d", categoryOptionsAction, category.ID)
		changeBalanceCallback := fmt.Sprintf("%s_%d", changeBalanceState, category.ID)
		balance := fmt.Sprintf("%.2f", category.Balance)
		total += category.Balance
		if category.Balance <= 0 {
			balance = fmt.Sprintf("%s %s", balance, "❗️❗️❗️")
		}
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(category.Name, categoryOptionsAction),
			tgbotapi.NewInlineKeyboardButtonData(balance, changeBalanceCallback),
		)

		rows = append(rows, row)
	}
	row := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📊 Total:", "total"),
		tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%.2f", total), "total"),
	)
	rows = append(rows, row)
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
func (b *tgBot) listCategoriesDeleteKeyboard(limit int32, offset int32) tgbotapi.InlineKeyboardMarkup {
	categories, err := b.db.ListCategories(context.Background(), db.ListCategoriesParams{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		log.Printf("failed to get categories: %v", err)
		return tgbotapi.InlineKeyboardMarkup{}
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, category := range categories {
		deleteCategoryCallbackData := fmt.Sprintf("%s_%d", deleteCategoryAction, category.ID)

		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(category.Name, deleteCategoryCallbackData),
		)

		rows = append(rows, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func (b *tgBot) listCategoriesLimitsKeyboard(limit int32, offset int32) tgbotapi.InlineKeyboardMarkup {
	categories, err := b.db.ListCategories(context.Background(), db.ListCategoriesParams{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		log.Printf("failed to get categories: %v", err)
		return tgbotapi.InlineKeyboardMarkup{}
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, category := range categories {
		changeNameCallback := fmt.Sprintf("%s_%d", categoryOptionsAction, category.ID)
		changePermanentBalanceCallback := fmt.Sprintf("%s_%d", changeCategoryLimitsCallback, category.ID)

		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(category.Name, changeNameCallback),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%.0f", category.PermanentBalance), changePermanentBalanceCallback),
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
