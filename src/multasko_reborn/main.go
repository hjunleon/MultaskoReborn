package main

import (
	"fmt"
	"log"
	db_helper "multasko_reborn/db/rethink"
	models "multasko_reborn/models"
	my_utils "multasko_reborn/my_utils"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// from "multasko_reborn/models" import

type BotState int64

const (
	IDLE     BotState = 0
	TODO     BotState = 1
	NOTE     BotState = 2
	RESOURCE BotState = 3
)

// Response send message back to author
// func (handler EchoHandler) HandleMessage(m tgbotapi.Message) error {
// 	_, err := m.QuickSend(m.Message().Text)
// 	return err
// }

var bot *tgbotapi.BotAPI
var TeleBotState BotState = IDLE
var CurTodo models.Todo = models.Todo{}

func handleToDo() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		log.Printf("Update Id %d", update.UpdateID)
	}

}

func main() {
	db_helper.TryHard()

	// return
	TELE_API := my_utils.GoDotEnvVariable("TELEGRAM_APITOKEN")

	fmt.Printf("API TOKEN: %s\n", TELE_API)
	var err error
	bot, err = tgbotapi.NewBotAPI(TELE_API)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		log.Printf("Update Id %d", update.UpdateID)
		if update.Message == nil {
			continue
		}

		// perhaps save the message to DB too
		writeStatus := db_helper.AddMessage(update.Message)
		if writeStatus {
			log.Println("Message Written successfullly!")
		} else {
			log.Println("Message Writing failed!")
		}
		if update.Message.IsCommand() { // ignore any non-Message updates
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch TeleBotState {
			case IDLE:
				// Extract the command from the Message.
				switch update.Message.Command() {
				case "help":
					msg.Text = "I understand /note, /todo and /resource. If not stated will be treated as note"
				case "note":
					content := update.Message.CommandArguments()
					if my_utils.IsEmptyStr(content) {
						msg.Text = fmt.Sprintf("Your note can't be empty :(")
						break
					}
					curNote := models.Note{
						Timestamp: update.Message.Time(),
						Info:      content,
					}
					writeStatus := db_helper.AddNote(curNote)
					if writeStatus {
						log.Println("Written successfullly!")
					} else {
						log.Println("Writing failed!")
					}
					msg.Text = fmt.Sprintf("Added note %s :)", content)
				case "todo":
					// check for command Args
					arg := update.Message.CommandArguments()
					args := strings.Split(arg, "/d")
					fmt.Printf("Arg: %s, length: %d\n", arg, len(args))
					if len(args) == 1 {
						if my_utils.IsEmptyStr(args[0]) {
							msg.Text = "Invalid todo :("
							break
						}
						msg.Text = "Any deadline?"
						CurTodo.Info = args[0]
						TeleBotState = TODO
						break
					}
					if len(args) == 0 {
						msg.Text = "Invalid todo :("
						break
					}
					// args = strings.TrimSpace(args)
					// check if there's date
					d := my_utils.ParseDates(args[1])
					CurTodo.Info = args[0]
					if d == nil {
						msg.Text = "Any deadline?"
						TeleBotState = TODO
						break
					}
					CurTodo.Deadline = *d

					timeStr := fmt.Sprintf("%d/%d/%d %d:%d", d.Day(), d.Month(), d.Year(), d.Hour(), d.Minute())
					msg.Text = fmt.Sprintf("Added todo %s with deadline %s", args[0], timeStr)
					CurTodo.Timestamp = update.Message.Time()
					writeStatus := db_helper.AddTodo(CurTodo)
					if writeStatus {
						log.Println("Written successfullly!")
					} else {
						log.Println("Writing failed!")
					}
				case "resource":
					content := update.Message.CommandArguments()
					if my_utils.IsEmptyStr(content) {
						msg.Text = fmt.Sprintf("Your resource can't be empty :(")
						break
					}
					curResource := models.Resource{
						Timestamp: update.Message.Time(),
						Info:      content,
					}
					writeStatus := db_helper.AddResource(curResource)
					if writeStatus {
						log.Println("Written successfullly!")
					} else {
						log.Println("Writing failed!")
					}
					msg.Text = fmt.Sprintf("Added resource %s :D", content)
				default:
					msg.Text = "I don't know that command"
				}

			case TODO:
				log.Printf("TODO STATE")
				arg := update.Message.Text
				d := my_utils.ParseDates(arg)
				if d == nil {
					msg.Text = "Invalid deadline"
					TeleBotState = IDLE
					break
				}
				CurTodo.Deadline = *d
				timeStr := my_utils.GetFormattedTime(d)
				msg.Text = fmt.Sprintf("Added todo %s with deadline %s", CurTodo.Info, timeStr)
				CurTodo.Timestamp = update.Message.Time()
				writeStatus := db_helper.AddTodo(CurTodo)
				if writeStatus {
					log.Println("Written successfullly!")
				} else {
					log.Println("Writing failed!")
				}
			default:
				msg.Text = "No idea how you got here :\\"
			}
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		} else {
			log.Println("Just a message not command")

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch TeleBotState {
			case IDLE:
				content := update.Message.CommandArguments()
				if my_utils.IsEmptyStr(content) {
					msg.Text = fmt.Sprintf("Your note can't be empty :(")
					break
				}
				curNote := models.Note{
					Timestamp: update.Message.Time(),
					Info:      content,
				}
				writeStatus := db_helper.AddNote(curNote)
				if writeStatus {
					log.Println("Written successfullly!")
				} else {
					log.Println("Writing failed!")
				}
				msg.Text = fmt.Sprintf("Added note %s :)", content)
			case TODO:
				log.Printf("TODO STATE")
				arg := update.Message.Text
				d := my_utils.ParseDates(arg)
				if d == nil {
					msg.Text = "Invalid deadline"
					TeleBotState = IDLE
					break
				}
				timeStr := my_utils.GetFormattedTime(d)
				CurTodo.Deadline = *d
				CurTodo.Timestamp = update.Message.Time()
				msg.Text = fmt.Sprintf("Added todo %s with deadline %s", CurTodo.Info, timeStr)
				writeStatus := db_helper.AddTodo(CurTodo)
				if writeStatus {
					log.Println("Written successfullly!")
				} else {
					log.Println("Writing failed!")
				}
				TeleBotState = IDLE
			}
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
	}
}
