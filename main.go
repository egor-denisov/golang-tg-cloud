package main

import (
	"log"
	"main/clients/telegram"
	"main/db"
	"main/lib/e"
)


func main() {
    log.Printf("App was started")
	database, err := db.New()
	if err != nil {
		log.Fatal(e.Wrap("Cannot connect to database", err))
	}
	defer database.Close()

	tgClient := telegram.New(TOKEN, database)
	tgClient.Listen()
}