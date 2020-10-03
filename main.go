package main

import (
	"fmt"
	_ "github.com/blackadress/vaula/models"
	"github.com/blackadress/vaula/webapp"
	"os"
)

func main() {
	WebApp = App{}

	WebApp.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	WebApp.Run(":8000")
	fmt.Printf("WebApp running on port 8000")
}
