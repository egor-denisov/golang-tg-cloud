package main

import (
	"log"
	"main/clients/api"
	"main/clients/telegram"
	"main/db"
	env "main/environment"
	"main/lib/e"
	"main/storage"

	"github.com/joho/godotenv"
)

func main() {
    log.Printf("App was started")
	//Uploading the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	env.SetEnvironment()
	// Connect to database
	database, err := db.New()
	if err != nil {
		log.Fatal(e.Wrap("Cannot connect to database", err))
	}
	defer database.Close()
	
	
	// Create tg bot instance
	tgClient := telegram.New(env.BOT_TOKEN, database)
	go tgClient.Listen()
	
	// Create storage for uploading and providing files
	storage := storage.New(env.BOT_TOKEN, database)
	go storage.StartUploading()

	// Create client for web version
	apiClient := api.New(database, storage)
	go apiClient.Listen()

	// Create channel for infinite work
	exit := make(chan struct{})
	<-exit
	defer close(exit)
}