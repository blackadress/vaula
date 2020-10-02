package main

import (
	"fmt"
	"github.com/blackadress/vaula/handlers"
	_ "github.com/blackadress/vaula/models"
	"os"
)

var app App

func main() {
	fmt.Printf("Hello from main\n")
	handlers.HelloFromUserHandlers()
	app = App{}

	app.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	app.Run(":8000")
}
