package env

import "os"

// Database environment variables
var (
	HOST     = os.Getenv("DB_HOST")
	PORT     = os.Getenv("DB_PORT")
	USER     = os.Getenv("DB_USER")
	PASSWORD = os.Getenv("DB_PASSWORD")
	DBNAME   = os.Getenv("DB_NAME")
)

// Telegram tokens
var BOT_TOKEN = os.Getenv("TG_BOT_TOKEN")

// Function for setting environment variables
func SetEnvironment() {
	HOST     = os.Getenv("DB_HOST")
	PORT     = os.Getenv("DB_PORT")
	USER     = os.Getenv("DB_USER")
	PASSWORD = os.Getenv("DB_PASSWORD")
	DBNAME   = os.Getenv("DB_NAME")
	BOT_TOKEN = os.Getenv("TG_BOT_TOKEN")
}

// Limits for api

var DOWNLOADING_LIMIT = 20971520
var UPLOADING_LIMIT = 52428800