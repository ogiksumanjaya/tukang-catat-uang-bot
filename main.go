package main

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/ogiksumanjaya/config"
	"github.com/ogiksumanjaya/helpers"
	"github.com/ogiksumanjaya/repository"
	"github.com/ogiksumanjaya/usecase"
	"log"
	"time"
)

func main() {
	var cfg config.Config

	connStr := cfg.GetDatabaseConfig()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bot, updates, err := cfg.InitConfigTelegram()
	if err != nil {
		log.Panic(err)
	}

	var lastCommand string
	var dataInput helpers.InputValue
	var dataInputTransfer helpers.InputTranferValue

	userAllowed := cfg.GetTelegram().UsernameAllowed
	tgRepo := repository.NewTelegramRepository(userAllowed)
	accountRepo := repository.NewAccountRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	for update := range updates {
		// Create a new context for each update with a 15-second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		tgReplayUC := usecase.NewTelegramReplayUsecase(bot, update, tgRepo, accountRepo, categoryRepo, transactionRepo)
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				tgReplayUC.StartMessageReplay()
			case "/masuk":
				err := tgReplayUC.IncreaseMessageReplay()
				if err != nil {
					continue
				}
				lastCommand = "/masuk"
			case "/keluar":
				err := tgReplayUC.DecreaseMessageReplay()
				if err != nil {
					continue
				}
				lastCommand = "/keluar"
			case "/transfer":
				err := tgReplayUC.StartTransferToAccount()
				if err != nil {
					continue
				}
				lastCommand = "/transfer"
			default:
				if lastCommand == "/masuk" {
					value, err := tgReplayUC.IncreaseMessageReplayCallback(ctx)
					if err != nil {
						continue
					}
					dataInput = value
					lastCommand = "/choose_in_bank"
				} else if lastCommand == "/keluar" {
					value, err := tgReplayUC.DecreaseMessageReplayCallback(ctx)
					if err != nil {
						continue
					}
					dataInput = value
					lastCommand = "/choose_out_bank"
				} else if lastCommand == "/transfer" {
					value, err := tgReplayUC.TransferFromAccountButtonCallback(ctx)
					if err != nil {
						continue
					}
					dataInputTransfer = value
					lastCommand = "/choose_transfer_from_bank"
				} else if lastCommand == "/choose_in_category" {
					tgReplayUC.HandleResponseSelectedCategory(ctx, dataInput, lastCommand)
					lastCommand = ""
				} else if lastCommand == "/choose_out_category" {
					tgReplayUC.HandleResponseSelectedCategory(ctx, dataInput, lastCommand)
					lastCommand = ""
				}
			}
		} else if update.CallbackQuery != nil {
			if lastCommand == "/choose_in_bank" {
				value, err := tgReplayUC.HandleResponseSelectedBank(ctx, dataInput)
				if err != nil {
					continue
				}
				dataInput = value

				lastCommand = "/choose_in_category"
			} else if lastCommand == "/choose_out_bank" {
				value, err := tgReplayUC.HandleResponseSelectedBank(ctx, dataInput)
				if err != nil {
					continue
				}
				dataInput = value

				lastCommand = "/choose_out_category"
			} else if lastCommand == "/choose_transfer_from_bank" {
				value := tgReplayUC.HandleResponseTransferFromBank(ctx, dataInputTransfer)
				dataInputTransfer = value

				// Clear last command
				lastCommand = "/choose_transfer_to_bank"
			} else if lastCommand == "/choose_transfer_to_bank" {
				tgReplayUC.HandleResponseTransferToBank(ctx, dataInputTransfer)

				// Clear last command
				lastCommand = ""
			}
		}
	}
}
