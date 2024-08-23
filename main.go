package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/ogiksumanjaya/config"
	"github.com/ogiksumanjaya/helpers"
	"github.com/ogiksumanjaya/repository"
	"github.com/ogiksumanjaya/usecase"
	"log"
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

	for update := range updates {
		tgReplayUC := usecase.NewTelegramReplayUsecase(bot, update, tgRepo)
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
					value, err := tgReplayUC.IncreaseMessageReplayCallback()
					if err != nil {
						continue
					}
					dataInput = value
					lastCommand = "/choose_in_bank"
				} else if lastCommand == "/keluar" {
					value, err := tgReplayUC.DecreaseMessageReplayCallback()
					if err != nil {
						continue
					}
					dataInput = value
					lastCommand = "/choose_out_bank"
				} else if lastCommand == "/transfer" {
					value, err := tgReplayUC.TransferFromAccountButtonCallback()
					if err != nil {
						continue
					}
					dataInputTransfer = value
					lastCommand = "/choose_transfer_from_bank"
				}
			}
		} else if update.CallbackQuery != nil {
			if lastCommand == "/choose_in_bank" {
				// Handle response dari tombol bank
				tgReplayUC.HandleResponseFromBank(dataInput)

				// Clear last command
				lastCommand = ""
			} else if lastCommand == "/choose_out_bank" {
				tgReplayUC.HandleResponseToBank(dataInput)

				// Clear last command
				lastCommand = ""
			} else if lastCommand == "/choose_transfer_from_bank" {
				value := tgReplayUC.HandleResponseTransferFromBank(dataInputTransfer)
				dataInputTransfer = value

				// Clear last command
				lastCommand = "/choose_transfer_to_bank"
			} else if lastCommand == "/choose_transfer_to_bank" {
				tgReplayUC.HandleResponseTransferToBank(dataInputTransfer)

				// Clear last command
				lastCommand = ""
			}
		}
	}
}
