package helpers

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ogiksumanjaya/entity"
	"strconv"
	"strings"
)

type InputValue struct {
	Username   string
	Amount     int
	Note       string
	Bank       string
	BankID     int
	Category   string
	CategoryID int
	Type       string // INCOME or EXPENSE
}

type InputTranferValue struct {
	Username string
	Amount   int
	FromBank string
	ToBank   string
}

type ReportTransaction struct {
	No            int
	Date          string
	Account       string
	Category      string
	Description   string
	IncomeAmount  float64
	ExpenseAmount float64
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

func GetBankKeyboardButton(replaceAccountName *string, bankList []entity.Account) tgbotapi.ReplyKeyboardMarkup {
	var filteredBankAccount []entity.Account

	if replaceAccountName != nil {
		for _, v := range bankList {
			if v.BankName != *replaceAccountName {
				filteredBankAccount = append(filteredBankAccount, v)
			}
		}
	} else {
		filteredBankAccount = bankList
	}

	var keyboardButtons [][]tgbotapi.KeyboardButton

	for _, account := range filteredBankAccount {
		button := tgbotapi.NewKeyboardButton(account.BankName)
		keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButtonRow(button))
	}

	keyboard := tgbotapi.NewReplyKeyboard(
		keyboardButtons...,
	)

	// Mengatur properti tambahan dari ReplyKeyboard
	keyboard.OneTimeKeyboard = true
	keyboard.ResizeKeyboard = true
	keyboard.InputFieldPlaceholder = "Pilih Bank"
	keyboard.Selective = true

	return keyboard
}

func GetCategoryKeyboardButton(categoryList []entity.Category) tgbotapi.ReplyKeyboardMarkup {
	var keyboardRows [][]tgbotapi.KeyboardButton

	for _, category := range categoryList {
		// Setiap tombol dimasukkan ke dalam baris baru untuk membuatnya vertikal
		row := tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(category.Name),
		)
		keyboardRows = append(keyboardRows, row)
	}

	// Membuat markup keyboard dengan semua baris tombol
	keyboard := tgbotapi.ReplyKeyboardMarkup{
		Keyboard:        keyboardRows,
		ResizeKeyboard:  true, // Mengatur ukuran tombol agar sesuai dengan lebar layar
		OneTimeKeyboard: true, // Menyembunyikan keyboard setelah pengguna memilih
	}

	return keyboard
}

func GetDateRangeKeyboardButton() tgbotapi.ReplyKeyboardMarkup {
	// Membuat tombol keyboard untuk memilih rentang tanggal
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Bulan Ini"),
			tgbotapi.NewKeyboardButton("Bulan Lalu"),
		),
	)

	return keyboard
}

func FormatRupiah(amount float64) string {
	// Format number to have comma as thousand separator
	formatted := fmt.Sprintf("%.2f", amount)

	// Split the formatted string to separate whole number and decimal parts
	parts := strings.Split(formatted, ".")

	// Insert commas as thousand separators
	n := len(parts[0])
	if n > 3 {
		remainder := n % 3
		if remainder > 0 {
			parts[0] = parts[0][:remainder] + "," + parts[0][remainder:]
		}
		for i := remainder + 3; i < n; i += 4 {
			parts[0] = parts[0][:i] + "," + parts[0][i:]
		}
	}

	// Combine the whole number part with the decimal part
	result := parts[0] + "." + parts[1]

	// Add the currency symbol
	return "Rp" + result
}
