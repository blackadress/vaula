package globals

import (
	"github.com/blackadress/vaula/webapp"
)

var WebApp = webapp.App{}

func Init() {
	globals.WebApp.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	globals.WebApp.Run(":8000")
	fmt.Printf("WebApp running on port 8000")
}
