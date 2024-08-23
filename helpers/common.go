package helpers

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type InputValue struct {
	Amount int
	Note   string
	Bank   string
}

type InputTranferValue struct {
	Amount   int
	FromBank string
	ToBank   string
}

func GetValueFromText(text string) (InputValue, error) {
	var value InputValue

	parts := strings.SplitN(text, " ", 2)
	if len(parts) < 2 {
		return value, errors.New("invalid format")
	}

	// Konversi nominal ke integer
	amount, err := strconv.Atoi(parts[0])
	if err != nil {
		return value, errors.New("invalid amount")
	}

	value.Amount = amount
	value.Note = parts[1]

	return value, nil
}

func ConvertNominalToInteger(nominal string) (InputTranferValue, error) {
	var valueTf InputTranferValue

	amount, err := strconv.Atoi(nominal)
	if err != nil {
		return InputTranferValue{}, errors.New("invalid amount")
	}

	valueTf.Amount = amount

	return valueTf, nil
}

func GetBankKeyboardButton(replaceAccount *string) tgbotapi.InlineKeyboardMarkup {
	bankAccount := []string{"CASH", "BCA", "CIMB", "JENIUS"}

	var filteredBankAccount []string

	if replaceAccount != nil {
		for _, v := range bankAccount {
			if v != *replaceAccount {
				filteredBankAccount = append(filteredBankAccount, v)
			}
		}
	} else {
		filteredBankAccount = bankAccount
	}

	var keyboardButtons []tgbotapi.InlineKeyboardButton

	for _, account := range filteredBankAccount {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(account, account))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(keyboardButtons...),
	)
	return keyboard
}
