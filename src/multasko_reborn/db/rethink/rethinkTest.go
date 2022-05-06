// package rethinkdb_test
package db_helper


import (
	"fmt"
	"log"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)


func Example() {
	url := "127.0.0.1:49154"
	session, err := r.Connect(r.ConnectOpts{
		Address: url, // endpoint without http
	})
	if err != nil {
		log.Fatalln(err)
	}

	res, err := r.Expr("Hello World").Run(session)
	if err != nil {
		log.Fatalln(err)
	}

	var response string
	err = res.One(&response)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(response)

	// Output:
	// Hello World
}

// func main() {
// 	Example()
// }