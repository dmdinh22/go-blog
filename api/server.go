package api

import (
	"fmt"
	"log"
	"os"

	"github.com/dmdinh22/go-blog/api/controllers"
	"github.com/dmdinh22/go-blog/api/seed"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var currentEnv = os.Getenv("ENVIRONMENT")

func Run() {
	var err error
	err = godotenv.Load()
	var dbDriver, dbUser, dbPassword, dbPort, dbHost, dbName string

	if err != nil {
		log.Fatalf("Error getting env vars", err)
	} else {
		fmt.Println("Loading env vars...")
	}

	if currentEnv != "production" {
		dbDriver = os.Getenv("DEV_DB_DRIVER")
		dbUser = os.Getenv("DEV_DB_USER")
		dbPassword = os.Getenv("DEV_DB_PASSWORD")
		dbPort = os.Getenv("DEV_DB_PORT")
		dbHost = os.Getenv("DEV_DB_HOST")
		dbName = os.Getenv("DEV_DB_NAME")
	} else {
		dbDriver = os.Getenv("DB_DRIVER")
		dbUser = os.Getenv("DB_USER")
		dbPassword = os.Getenv("DB_PASSWORD")
		dbPort = os.Getenv("DB_PORT")
		dbHost = os.Getenv("DB_HOST")
		dbName = os.Getenv("DB_NAME")
	}

	server.Initialize(dbDriver, dbUser, dbPassword, dbPort, dbHost, dbName)
	seed.Load(server.DB)
	server.Run(":8080")
}
