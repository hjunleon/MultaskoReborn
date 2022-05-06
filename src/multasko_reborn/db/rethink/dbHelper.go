package db_helper

import (
	"fmt"
	"log"
	"multasko_reborn/models"
	"multasko_reborn/my_utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// type DB struct {
// 	session
// }
var sessionArray []*r.Session
var url = ""
var tables [4]string = [...]string{"resources", "toDos", "notes", "messages"}

func InitDb() {
	url = my_utils.GoDotEnvVariable("RETHINK_URL")
	fmt.Printf("url: %s\n", url)
	log.Printf("url: %s\n", url)

	session, err := r.Connect(r.ConnectOpts{
		Address: url, // endpoint without http
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	var cursor *r.Cursor
	cursor, err = r.DBCreate("multasko").Run(session) //.Exec()
	if err != nil {
		log.Println(err)
	}
	// defer cursor.Close()

	for i := 0; i < len(tables); i++ {
		cursor, err = r.TableCreate(tables[i]).Run(session) //.Exec()
		if err != nil {
			log.Println(err)
		}
	}

	defer cursor.Close()

	sessionArray = append(sessionArray, session)
}

func ListAllTables() {
	session := sessionArray[0]
	cursor, err := r.TableList().Run(session)
	if err != nil {
		log.Println(err)
	}
	var response []string
	err = cursor.All(&response)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(response)
}

func AddMessage(msg *tgbotapi.Message) bool {

	session := sessionArray[0]
	cursor, err := r.Table(tables[3]).Insert(msg).Run(session)
	if err != nil {
		log.Printf("error: %s", err)
		return false
	}
	defer cursor.Close()
	return true
}

//[K comparable, V int64 | float64]
func AddTodo(item models.Todo) bool {
	session := sessionArray[0]
	cursor, err := r.Table(tables[1]).Insert(item).Run(session)
	if err != nil {
		log.Printf("error: %s", err)
		return false
	}
	defer cursor.Close()
	return true
}

func AddResource(item models.Resource) bool {
	session := sessionArray[0]
	cursor, err := r.Table(tables[0]).Insert(item).Run(session)
	if err != nil {
		log.Printf("error: %s", err)
		return false
	}
	defer cursor.Close()
	return true
}

func AddNote(item models.Note) bool {
	session := sessionArray[0]
	cursor, err := r.Table(tables[2]).Insert(item).Run(session)
	if err != nil {
		log.Printf("error: %s", err)
		return false
	}
	defer cursor.Close()
	return true
}

func AddItem(item, itemType models.Model) bool {
	session := sessionArray[0]
	var cursor *r.Cursor
	var err error
	switch itemType {
	case models.TODO:
		cursor, err = r.Table(tables[1]).Insert(item).Run(session)
	case models.RES:
		cursor, err = r.Table(tables[0]).Insert(item).Run(session)
	case models.NOTE:
		cursor, err = r.Table(tables[2]).Insert(item).Run(session)
	default:
		log.Println("Adding nothing")
		return false
	}
	if err != nil {
		log.Printf("error: %s", err)
		return false
	}
	defer cursor.Close()
	return true
}

func TryHard() {
	InitDb()
	ListAllTables()
	test()
}
func test() {
	session := sessionArray[0]
	if session == nil {
		session, err := r.Connect(r.ConnectOpts{
			Address: url, // endpoint without http
		})
		if err != nil {
			log.Fatalln(err)
		}
		sessionArray[0] = session
	}
	res, err := r.Expr("Hello World").Run(session)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Close()
}
