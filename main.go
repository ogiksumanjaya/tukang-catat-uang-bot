package main

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ogiksumanjaya/tukang-catat-uang-bot/config"
)

func main() {
	var cfg config.Config
	if _, ok := os.LookupEnv("SERVER_REST_PORT"); !ok {
		cfg.Setup(".env")
	}

	bot, err := tgbotapi.NewBotAPI(cfg.GetTelegram().Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	var lastCommand string
	var dataInput InputValue

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Halo! Saya adalah bot pencatat keuangan. Silahkan gunakan perintah /masuk untuk memasukkan pemasukan dan /keluar untuk memasukkan pengeluaran.")
				bot.Send(msg)
			case "/masuk":
				isAllowed := IsAllowedUsername(cfg.GetTelegram().UsernameAllowed, update.Message.From.UserName)

				if !isAllowed {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Maaf, kamu tidak diizinkan untuk menggunakan bot ini.")
					bot.Send(msg)
					continue
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Silahkan masukan pemasukanmu dengan contoh berikut:\n50000 gaji bulanan")
				bot.Send(msg)
				lastCommand = "/masuk"
			case "/keluar":
				isAllowed := IsAllowedUsername(cfg.GetTelegram().UsernameAllowed, update.Message.From.UserName)

				if !isAllowed {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Maaf, kamu tidak diizinkan untuk menggunakan bot ini.")
					bot.Send(msg)
					continue
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Silahkan masukan pengeluaranmu dengan contoh berikut:\n50000 beli minuman")
				bot.Send(msg)
				lastCommand = "/keluar"
			default:
				if lastCommand == "/masuk" {
					// Get value from text
					value, err := GetValueFromText(update.Message.Text)
					if err != nil {
						if err.Error() == "invalid amount" {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Nominal yang kamu masukan tidak valid. Silahkan coba lagi.")
							bot.Send(msg)
						} else {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Format yang kamu masukan salah. Silahkan coba lagi.")
							bot.Send(msg)
						}
						continue
					}

					dataInput = value

					// Kirim pesan untuk memilih bank
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Masukan ke bank apa?")
					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("CASH", "CASH"),
							tgbotapi.NewInlineKeyboardButtonData("BCA", "BCA"),
							tgbotapi.NewInlineKeyboardButtonData("CIMB", "CIMB"),
							tgbotapi.NewInlineKeyboardButtonData("JENIUS", "JENIUS"),
						),
					)
					msg.ReplyMarkup = keyboard
					bot.Send(msg)

					lastCommand = "/choose_in_bank"
				} else if lastCommand == "/keluar" {
					// Get value from text
					value, err := GetValueFromText(update.Message.Text)
					if err != nil {
						if err.Error() == "invalid amount" {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Nominal yang kamu masukan tidak valid. Silahkan coba lagi.")
							bot.Send(msg)
						} else {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Format yang kamu masukan salah. Silahkan coba lagi.")
							bot.Send(msg)
						}
						continue
					}

					dataInput = value

					// Kirim pesan untuk memilih bank
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Keluarkan dari bank apa?")
					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("CASH", "CASH"),
							tgbotapi.NewInlineKeyboardButtonData("BCA", "BCA"),
							tgbotapi.NewInlineKeyboardButtonData("CIMB", "CIMB"),
							tgbotapi.NewInlineKeyboardButtonData("JENIUS", "JENIUS"),
						),
					)
					msg.ReplyMarkup = keyboard
					bot.Send(msg)

					lastCommand = "/choose_out_bank"
				}
			}
		} else if update.CallbackQuery != nil {
			if lastCommand == "/choose_in_bank" {
				// Handle response dari tombol bank
				bank := update.CallbackQuery.Data
				dataInput.Bank = bank

				// Kirim konfirmasi bahwa pemasukan telah tercatat
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Pemasukan sudah tercatat.")
				// send detail pemasukan
				msg.Text = "Pemasukan: " + strconv.Itoa(dataInput.Amount) + "\n" + "Catatan: " + dataInput.Note + "\n" + "Bank: " + dataInput.Bank
				bot.Send(msg)

				// Clear last command
				lastCommand = ""
			} else if lastCommand == "/choose_out_bank" {
				// Handle response dari tombol bank
				bank := update.CallbackQuery.Data
				dataInput.Bank = bank

				// Kirim konfirmasi bahwa pengeluaran telah tercatat
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Pengeluaran sudah tercatat.")
				// send detail pengeluaran
				msg.Text = "Pengeluaran: " + strconv.Itoa(dataInput.Amount) + "\n" + "Catatan: " + dataInput.Note + "\n" + "Bank: " + dataInput.Bank
				bot.Send(msg)

				// Clear last command
				lastCommand = ""
			}
		}
	}
}

func IsAllowedUsername(stringAllowed, username string) bool {
	allowedUsername := strings.Split(stringAllowed, ",")
	for _, allowed := range allowedUsername {
		if username == allowed {
			return true
		}
	}
	return false
}

type InputValue struct {
	Amount int
	Note   string
	Bank   string
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
