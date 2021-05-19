package main

import (
	"fmt"
	"log"
	"os"

	"github.com/blackadress/vaula/handlers"
	"github.com/joho/godotenv"
)

var app = handlers.App{}

func main() {

	app.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	app.Run(":8000")
	fmt.Printf("WebApp running on port 8000")
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No '.env' found")
	}
}
