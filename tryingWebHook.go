package main

import (
	"log"
    "fmt"
	"net/http"
	"github.com/joho/godotenv"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	bot, err := tgbotapi.NewBotAPI(TELE_API)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	wh, _ := tgbotapi.NewWebhookWithCert("https://www.example.com:8443/"+bot.Token, "cert.pem")

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)

	for update := range updates {
		log.Printf("%+v\n", update)
	}
}