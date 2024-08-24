package usecase

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ogiksumanjaya/entity"
	"github.com/ogiksumanjaya/helpers"
	"github.com/ogiksumanjaya/repository"
	"log"
	"strconv"
)

type TelegramReplayUsecase struct {
	bot             *tgbotapi.BotAPI
	update          tgbotapi.Update
	telegramRepo    *repository.TelegramRepository
	accountRepo     *repository.AccountRepository
	categoryRepo    *repository.CategoryRepository
	transactionRepo *repository.TransactionRepository
}

func NewTelegramReplayUsecase(bot *tgbotapi.BotAPI, update tgbotapi.Update, telegramRepo *repository.TelegramRepository, accountRepo *repository.AccountRepository, categoryRepo *repository.CategoryRepository, transactionRepo *repository.TransactionRepository) *TelegramReplayUsecase {
	return &TelegramReplayUsecase{
		bot:             bot,
		update:          update,
		telegramRepo:    telegramRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
		transactionRepo: transactionRepo,
	}
}

func (t *TelegramReplayUsecase) StartMessageReplay() {
	chatID := t.update.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "Halo! Saya adalah bot pencatat keuangan.\n\nSilahkan gunakan perintah berikut:\n\n/masuk untuk memasukkan pemasukan dan\n/keluar untuk memasukkan pengeluaran.\n/transfer untuk melakukan transfer antar bank/account.")
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

func (t *TelegramReplayUsecase) IncreaseMessageReplayCallback(ctx context.Context) (helpers.InputValue, error) {
	chatID := t.update.Message.Chat.ID
	username := t.update.Message.From.UserName

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

	// GetBankList
	bankList, err := t.accountRepo.GetAccountList(ctx, username)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat mengambil data bank.")
		t.bot.Send(msg)
		log.Println(err)
		return helpers.InputValue{}, err
	}

	// Send message to choose bank
	msg := tgbotapi.NewMessage(chatID, "Masukan ke bank apa?")
	msg.ReplyMarkup = helpers.GetBankKeyboardButton(nil, bankList)
	t.bot.Send(msg)

	return value, nil
}

func (t *TelegramReplayUsecase) DecreaseMessageReplayCallback(ctx context.Context) (helpers.InputValue, error) {
	chatID := t.update.Message.Chat.ID
	username := t.update.Message.From.UserName

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

	// GetBankList
	bankList, err := t.accountRepo.GetAccountList(ctx, username)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat mengambil data bank.")
		t.bot.Send(msg)
		return helpers.InputValue{}, err
	}

	// Send message to choose bank
	msg := tgbotapi.NewMessage(chatID, "Keluar dari bank apa?")
	msg.ReplyMarkup = helpers.GetBankKeyboardButton(nil, bankList)
	t.bot.Send(msg)

	return value, nil
}

func (t *TelegramReplayUsecase) HandleResponseSelectedBank(ctx context.Context, dataInput helpers.InputValue) (helpers.InputValue, error) {
	chatID := t.update.CallbackQuery.Message.Chat.ID
	username := t.update.CallbackQuery.From.UserName

	// Handle response dari tombol bank
	bank := t.update.CallbackQuery.Data
	dataInput.Bank = bank

	// Get Category List
	categoryList, err := t.categoryRepo.GetCategoryList(ctx, username)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat mengambil data kategori.")
		t.bot.Send(msg)
		return helpers.InputValue{}, err
	}

	// Send message to choose category
	msg := tgbotapi.NewMessage(chatID, "Pilih kategori:")
	msg.ReplyMarkup = helpers.GetCategoryKeyboardButton(categoryList)
	t.bot.Send(msg)

	return dataInput, nil
}

func (t *TelegramReplayUsecase) HandleResponseSelectedCategory(ctx context.Context, dataInput helpers.InputValue, lastCommand string) {
	chatID := t.update.Message.Chat.ID
	username := t.update.Message.From.UserName
	dataInput.Username = username

	category := t.update.Message.Text
	dataInput.Category = category

	// Get Data Category by Name
	var dataCategory entity.Category
	dataCategory.Name = category
	dataCategory.Username = username
	categoryData, err := t.categoryRepo.GetCategoryByName(ctx, dataCategory)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat mengambil data kategori.")
		t.bot.Send(msg)
		return
	}
	dataInput.CategoryID = categoryData.ID

	// Get Data Account by Name
	var dataAccount entity.Account
	dataAccount.BankName = dataInput.Bank
	dataAccount.Username = username
	accountData, err := t.accountRepo.GetAccountByName(ctx, dataAccount)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat mengambil data bank.")
		t.bot.Send(msg)
		return
	}
	dataInput.BankID = accountData.ID

	if lastCommand == "/choose_in_category" {
		dataInput.Type = "INCOME"

		// Insert to transaction
		err = t.transactionRepo.CreateNewTransaction(ctx, dataInput)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat menyimpan data.")
			t.bot.Send(msg)
			return
		}

		// Update Balance
		newBalance := accountData.Balance + float64(dataInput.Amount)
		accountData.Balance = newBalance
		err = t.accountRepo.UpdateBalance(ctx, accountData)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat mengupdate data.")
			t.bot.Send(msg)
			return
		}

		// Kirim konfirmasi bahwa pemasukan telah tercatat
		msg := tgbotapi.NewMessage(chatID, "Pemasukan sudah tercatat.")
		// send detail pemasukan
		msg.Text = "Pemasukan: " + strconv.Itoa(dataInput.Amount) + "\n" + "Catatan: " + dataInput.Note + "\n" + "Bank: " + dataInput.Bank + "\n" + "Category: " + dataInput.Category
		t.bot.Send(msg)
	} else if lastCommand == "/choose_out_category" {
		dataInput.Type = "EXPENSE"

		// Insert to transaction
		err = t.transactionRepo.CreateNewTransaction(ctx, dataInput)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat menyimpan data.")
			t.bot.Send(msg)
			return
		}

		// Update Balance
		newBalance := accountData.Balance - float64(dataInput.Amount)
		accountData.Balance = newBalance
		err = t.accountRepo.UpdateBalance(ctx, accountData)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat mengupdate data.")
			t.bot.Send(msg)
			return
		}

		// Kirim konfirmasi bahwa pengeluaran telah tercatat
		msg := tgbotapi.NewMessage(chatID, "Pengeluaran sudah tercatat.")
		// send detail pengeluaran
		msg.Text = "Pengeluaran: " + strconv.Itoa(dataInput.Amount) + "\n" + "Catatan: " + dataInput.Note + "\n" + "Bank: " + dataInput.Bank + "\n" + "Category: " + dataInput.Category
		t.bot.Send(msg)
	}
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

func (t *TelegramReplayUsecase) TransferFromAccountButtonCallback(ctx context.Context) (helpers.InputTranferValue, error) {
	chatID := t.update.Message.Chat.ID
	username := t.update.Message.From.UserName

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

	// GetBankList
	bankList, err := t.accountRepo.GetAccountList(ctx, username)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat mengambil data bank.")
		t.bot.Send(msg)
		return helpers.InputTranferValue{}, err
	}

	// Send message to choose bank
	msg := tgbotapi.NewMessage(chatID, "Transfer dari bank apa?")
	msg.ReplyMarkup = helpers.GetBankKeyboardButton(nil, bankList)
	t.bot.Send(msg)

	return value, nil
}

func (t *TelegramReplayUsecase) HandleResponseTransferFromBank(ctx context.Context, dataInput helpers.InputTranferValue) helpers.InputTranferValue {
	chatID := t.update.CallbackQuery.Message.Chat.ID
	username := t.update.CallbackQuery.From.UserName

	// Handle response dari tombol bank
	bank := t.update.CallbackQuery.Data
	dataInput.FromBank = bank

	// GetBankList
	bankList, err := t.accountRepo.GetAccountList(ctx, username)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat mengambil data bank.")
		t.bot.Send(msg)
		return helpers.InputTranferValue{}
	}

	// Send message to choose bank
	msg := tgbotapi.NewMessage(chatID, "Transfer ke bank apa?")
	msg.ReplyMarkup = helpers.GetBankKeyboardButton(&bank, bankList)
	t.bot.Send(msg)

	return dataInput
}

func (t *TelegramReplayUsecase) HandleResponseTransferToBank(ctx context.Context, dataInput helpers.InputTranferValue) {
	chatID := t.update.CallbackQuery.Message.Chat.ID
	username := t.update.CallbackQuery.From.UserName

	var dataInputToTrx helpers.InputValue
	dataInputToTrx.Username = username
	dataInputToTrx.Amount = dataInput.Amount

	// Handle response dari tombol bank
	bank := t.update.CallbackQuery.Data
	dataInput.ToBank = bank

	// Get Data Category by Name
	var dataCategory entity.Category
	dataCategory.Name = "Lainnya"
	dataCategory.Username = username
	categoryData, err := t.categoryRepo.GetCategoryByName(ctx, dataCategory)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat mengambil data kategori.")
		t.bot.Send(msg)
		return
	}
	dataInputToTrx.CategoryID = categoryData.ID

	// Get Data From Account by Name
	var dataAccountFrom entity.Account
	dataAccountFrom.BankName = dataInput.FromBank
	dataAccountFrom.Username = username
	accountDataFrom, err := t.accountRepo.GetAccountByName(ctx, dataAccountFrom)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat mengambil data bank.")
		t.bot.Send(msg)
		return
	}

	// Get Data To Account by Name
	var dataAccountTo entity.Account
	dataAccountTo.BankName = dataInput.ToBank
	dataAccountTo.Username = username
	accountDataTo, err := t.accountRepo.GetAccountByName(ctx, dataAccountTo)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat mengambil data bank.")
		t.bot.Send(msg)
		return
	}

	// Insert to transaction Expense
	dataInputToTrx.Type = "EXPENSE"
	dataInputToTrx.BankID = accountDataFrom.ID
	dataInputToTrx.Note = "Transfer ke " + dataInput.ToBank
	err = t.transactionRepo.CreateNewTransaction(ctx, dataInputToTrx)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat menyimpan data.")
		t.bot.Send(msg)
		return
	}

	// Update Balance From Account
	newBalanceFrom := accountDataFrom.Balance - float64(dataInput.Amount)
	accountDataFrom.Balance = newBalanceFrom
	err = t.accountRepo.UpdateBalance(ctx, accountDataFrom)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat mengupdate data.")
		t.bot.Send(msg)
		return
	}

	// Insert to transaction Income
	dataInputToTrx.Type = "INCOME"
	dataInputToTrx.BankID = accountDataTo.ID
	dataInputToTrx.Note = "Transfer dari " + dataInput.FromBank
	err = t.transactionRepo.CreateNewTransaction(ctx, dataInputToTrx)

	// Update Balance To Account
	newBalanceTo := accountDataTo.Balance + float64(dataInput.Amount)
	accountDataTo.Balance = newBalanceTo
	err = t.accountRepo.UpdateBalance(ctx, accountDataTo)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Terjadi kesalahan saat mengupdate data.")
		t.bot.Send(msg)
		return
	}

	// Kirim konfirmasi bahwa transfer telah tercatat
	msg := tgbotapi.NewMessage(chatID, "Transfer sudah tercatat.")
	// send detail transfer
	msg.Text = "Transfer: " + strconv.Itoa(dataInput.Amount) + "\n" + "Dari Bank: " + dataInput.FromBank + "\n" + "Ke Bank: " + dataInput.ToBank
	t.bot.Send(msg)
}
