package main

import (
	"fmt"
	"time"
)

const (
	layout = "2006-01-02"
)

func main() {
	date := "2022-02-01"
	time, _ := time.Parse(layout, date)
	fmt.Println(time)
}
