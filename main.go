package main

import (
	"log"
	"main/clients/api"
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
	
	stop := make(chan bool)
	defer close(stop)

	// tgClient := telegram.New(env.TOKEN, database)
	// go tgClient.Listen()
	
	apiClient := api.New(database)
	go apiClient.Listen()

	<-stop
}