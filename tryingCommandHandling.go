package main

import (
    "fmt"
    "log"
    "os"
	"github.com/joho/godotenv"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
    u.Timeout = 30

    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil { // ignore any non-Message updates
            continue
        }

        if !update.Message.IsCommand() { // ignore any non-command Messages
            continue
        }

        // Create a new MessageConfig. We don't have text yet,
        // so we leave it empty.
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

        // Extract the command from the Message.
        switch update.Message.Command() {
        case "help":
            msg.Text = "I understand /sayhi and /status."
        case "sayhi":
            msg.Text = "Hi :)"
        case "status":
            msg.Text = "I'm ok."
        default:
            msg.Text = "I don't know that command"
        }

        if _, err := bot.Send(msg); err != nil {
            log.Panic(err)
        }
    }
}
