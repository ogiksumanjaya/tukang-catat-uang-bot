package usecase

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ogiksumanjaya/helpers"
	"github.com/ogiksumanjaya/repository"
	"strconv"
)

type TelegramReplayUsecase struct {
	bot          *tgbotapi.BotAPI
	update       tgbotapi.Update
	telegramRepo *repository.TelegramRepository
}

func NewTelegramReplayUsecase(bot *tgbotapi.BotAPI, update tgbotapi.Update, telegramRepo *repository.TelegramRepository) *TelegramReplayUsecase {
	return &TelegramReplayUsecase{
		bot:          bot,
		update:       update,
		telegramRepo: telegramRepo,
	}
}

func (t *TelegramReplayUsecase) StartMessageReplay() {
	chatID := t.update.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "Halo! Saya adalah bot pencatat keuangan.\n\nSilahkan gunakan perintah berikut:\n\n/masuk untuk memasukkan pemasukan dan\n/keluar untuk memasukkan pengeluaran.")
	t.bot.Send(msg)
}

func (t *TelegramReplayUsecase) IncreaseMessageReplay() error {
	chatID := t.update.Message.Chat.ID
	username := t.update.Message.From.UserName

	// Check allowed username
	isAllowed := t.telegramRepo.IsAllowedUsername(username)
	if !isAllowed {
		msg := tgbotapi.NewMessage(chatID, "Maaf, kamu tidak diizinkan untuk menggunakan bot ini.")
		t.bot.Send(msg)
		return errors.New("username not allowed")
	}

	msg := tgbotapi.NewMessage(chatID, "Silahkan masukan pemasukanmu dengan contoh berikut:\n50000 gaji bulanan")
	t.bot.Send(msg)

	return nil
}

func (t *TelegramReplayUsecase) DecreaseMessageReplay() error {
	chatID := t.update.Message.Chat.ID
	username := t.update.Message.From.UserName

	// Check allowed username
	isAllowed := t.telegramRepo.IsAllowedUsername(username)
	if !isAllowed {
		msg := tgbotapi.NewMessage(chatID, "Maaf, kamu tidak diizinkan untuk menggunakan bot ini.")
		t.bot.Send(msg)
		return errors.New("username not allowed")
	}

	msg := tgbotapi.NewMessage(chatID, "Silahkan masukan pengeluaranmu dengan contoh berikut:\n50000 beli minuman")
	t.bot.Send(msg)

	return nil
}

func (t *TelegramReplayUsecase) IncreaseMessageReplayCallback() (helpers.InputValue, error) {
	chatID := t.update.Message.Chat.ID

	// Get value from text
	value, err := helpers.GetValueFromText(t.update.Message.Text)
	if err != nil {
		if err.Error() == "invalid amount" {
			msg := tgbotapi.NewMessage(chatID, "Nominal yang kamu masukan tidak valid. Silahkan coba lagi.")
			t.bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, "Format yang kamu masukan salah. Silahkan coba lagi.")
			t.bot.Send(msg)
		}
		return helpers.InputValue{}, err
	}

	// Send message to choose bank
	msg := tgbotapi.NewMessage(chatID, "Masukan ke bank apa?")
	msg.ReplyMarkup = helpers.GetBankKeyboardButton(nil)
	t.bot.Send(msg)

	return value, nil
}

func (t *TelegramReplayUsecase) DecreaseMessageReplayCallback() (helpers.InputValue, error) {
	chatID := t.update.Message.Chat.ID

	// Get value from text
	value, err := helpers.GetValueFromText(t.update.Message.Text)
	if err != nil {
		if err.Error() == "invalid amount" {
			msg := tgbotapi.NewMessage(chatID, "Nominal yang kamu masukan tidak valid. Silahkan coba lagi.")
			t.bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, "Format yang kamu masukan salah. Silahkan coba lagi.")
			t.bot.Send(msg)
		}
		return helpers.InputValue{}, err
	}

	// Send message to choose bank
	msg := tgbotapi.NewMessage(chatID, "Keluar dari bank apa?")
	msg.ReplyMarkup = helpers.GetBankKeyboardButton(nil)
	t.bot.Send(msg)

	return value, nil
}

func (t *TelegramReplayUsecase) HandleResponseFromBank(dataInput helpers.InputValue) {
	chatID := t.update.CallbackQuery.Message.Chat.ID

	// Handle response dari tombol bank
	bank := t.update.CallbackQuery.Data
	dataInput.Bank = bank

	// Kirim konfirmasi bahwa pemasukan telah tercatat
	msg := tgbotapi.NewMessage(chatID, "Pemasukan sudah tercatat.")
	// send detail pemasukan
	msg.Text = "Pemasukan: " + strconv.Itoa(dataInput.Amount) + "\n" + "Catatan: " + dataInput.Note + "\n" + "Bank: " + dataInput.Bank
	t.bot.Send(msg)

}

func (t *TelegramReplayUsecase) HandleResponseToBank(dataInput helpers.InputValue) {
	chatID := t.update.CallbackQuery.Message.Chat.ID

	// Handle response dari tombol bank
	bank := t.update.CallbackQuery.Data
	dataInput.Bank = bank

	// Kirim konfirmasi bahwa pengeluaran telah tercatat
	msg := tgbotapi.NewMessage(chatID, "Pengeluaran sudah tercatat.")
	// send detail pengeluaran
	msg.Text = "Pengeluaran: " + strconv.Itoa(dataInput.Amount) + "\n" + "Catatan: " + dataInput.Note + "\n" + "Bank: " + dataInput.Bank
	t.bot.Send(msg)
}

func (t *TelegramReplayUsecase) StartTransferToAccount() error {
	chatID := t.update.Message.Chat.ID
	username := t.update.Message.From.UserName

	// Check allowed username
	isAllowed := t.telegramRepo.IsAllowedUsername(username)
	if !isAllowed {
		msg := tgbotapi.NewMessage(chatID, "Maaf, kamu tidak diizinkan untuk menggunakan bot ini.")
		t.bot.Send(msg)
		return errors.New("username not allowed")
	}

	msg := tgbotapi.NewMessage(chatID, "Silahkan masukan nominal yang akan transfer")
	t.bot.Send(msg)

	return nil
}

func (t *TelegramReplayUsecase) TransferFromAccountButtonCallback() (helpers.InputTranferValue, error) {
	chatID := t.update.Message.Chat.ID

	// Get value from text
	value, err := helpers.ConvertNominalToInteger(t.update.Message.Text)
	if err != nil {
		if err.Error() == "invalid amount" {
			msg := tgbotapi.NewMessage(chatID, "Nominal yang kamu masukan tidak valid. Silahkan coba lagi.")
			t.bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, "Format yang kamu masukan salah. Silahkan coba lagi.")
			t.bot.Send(msg)
		}
		return helpers.InputTranferValue{}, err
	}

	// Send message to choose bank
	msg := tgbotapi.NewMessage(chatID, "Transfer dari bank apa?")
	msg.ReplyMarkup = helpers.GetBankKeyboardButton(nil)
	t.bot.Send(msg)

	return value, nil
}

func (t *TelegramReplayUsecase) HandleResponseTransferFromBank(dataInput helpers.InputTranferValue) helpers.InputTranferValue {
	chatID := t.update.CallbackQuery.Message.Chat.ID

	// Handle response dari tombol bank
	bank := t.update.CallbackQuery.Data
	dataInput.FromBank = bank

	// Send message to choose bank
	msg := tgbotapi.NewMessage(chatID, "Transfer ke bank apa?")
	msg.ReplyMarkup = helpers.GetBankKeyboardButton(&bank)
	t.bot.Send(msg)

	return dataInput
}

func (t *TelegramReplayUsecase) HandleResponseTransferToBank(dataInput helpers.InputTranferValue) {
	chatID := t.update.CallbackQuery.Message.Chat.ID

	// Handle response dari tombol bank
	bank := t.update.CallbackQuery.Data
	dataInput.ToBank = bank

	// Kirim konfirmasi bahwa transfer telah tercatat
	msg := tgbotapi.NewMessage(chatID, "Transfer sudah tercatat.")
	// send detail transfer
	msg.Text = "Transfer: " + strconv.Itoa(dataInput.Amount) + "\n" + "Dari Bank: " + dataInput.FromBank + "\n" + "Ke Bank: " + dataInput.ToBank
	t.bot.Send(msg)
}
