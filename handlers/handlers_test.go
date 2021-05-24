package handlers

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var a App
var BASE_URL string

func TestMain(m *testing.M) {

	if err := godotenv.Load("../.env"); err != nil {
		log.Print("TEST: '.env' no encontrado")
	}
	BASE_URL = os.Getenv("URL")

	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	// asegurarse de que todas las tablas existen
	ensureTableUsuarioExists()
	ensureTableAlternativaExists()

	code := m.Run()

	// limpiar las tablas de la BD
	clearTableUsuario()
	clearTableAlternativa()
	os.Exit(code)
}
