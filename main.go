package main

import (
	"log"
	"main/clients/api"
	"main/db"
	env "main/environment"
	"main/lib/e"
	"main/storage"
)

func main() {
    log.Printf("App was started")
	database, err := db.New()
	if err != nil {
		log.Fatal(e.Wrap("Cannot connect to database", err))
	}
	defer database.Close()
	
	stop := make(chan bool)
	defer close(stop)

	// tgClient := telegram.New(env.BOT_TOKEN, database)
	// go tgClient.Listen()
	
	storage := storage.New(env.STORAGE_TOKEN, database)

	apiClient := api.New(database, storage)
	go apiClient.Listen()

	<-stop
}