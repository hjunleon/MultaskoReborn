package main

import (
	"fmt"
    "log"
    "os"
	"github.com/joho/godotenv"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
    tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("1"),
        tgbotapi.NewKeyboardButton("2"),
        tgbotapi.NewKeyboardButton("3"),
    ),
    tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("4"),
        tgbotapi.NewKeyboardButton("5"),
        tgbotapi.NewKeyboardButton("6"),
    ),
)
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")
  
	if err != nil {
	//   log.Fatalf("Error loading .env file")
      fmt.Printf("Error loading .env file\n")
	}
  
	return os.Getenv(key)
  }
func main() {
    TELE_API := goDotEnvVariable("TELEGRAM_APITOKEN")
    
    fmt.Printf("API TOKEN: %s\n", TELE_API)
    bot, err := tgbotapi.NewBotAPI(TELE_API)
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true

    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil { // ignore non-Message updates
            continue
        }

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

        switch update.Message.Text {
        case "open":
            msg.ReplyMarkup = numericKeyboard
        case "close":
            msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
        }

        if _, err := bot.Send(msg); err != nil {
            log.Panic(err)
        }
    }
}
