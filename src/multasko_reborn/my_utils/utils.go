package my_utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		//   log.Fatalf("Error loading .env file")
		fmt.Printf("Error loading .env file\n")
	}

	return os.Getenv(key)
}

//https://go.dev/src/time/format.go

var DateFormats []string = []string{
	"20060102",
	"2006/01/02",
	"2006-01-02",
	"2006-01-02T15:04:05.000Z",
	"02/01/2006",
	"02/01/06",
	"02-01-2006",
	"02/01/06_15:04", //my most used one
	"02-01-06",
	"Jan 2, 2006 at 3:04pm (MST)",
}

func ParseDates(payload string) *time.Time {
	var d time.Time
	var err error
	payload = strings.TrimSpace(payload)
	log.Printf("Parsing time string %s with length %d", payload, len(payload))
	for _, df := range DateFormats {
		// log.Printf("Format: %s\n", df)
		d, err = time.Parse(df, strings.TrimSpace(payload))
		// log.Printf("%d/%d/%d %d:%d\n", d.Day(), d.Month(), d.Year(), d.Hour(), d.Minute())
		if err == nil {
			return &d
		}
	}
	return nil
}

func GetFormattedTime(d *time.Time) string {
	return fmt.Sprintf("%02d/%02d/%02d %02d:%02d", d.Day(), d.Month(), d.Year(), d.Hour(), d.Minute())
}

func IsEmptyStr(payload string) bool {
	if len(payload) == 0 || payload == "" {
		return true
	}
	return false
}
