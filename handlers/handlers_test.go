package handlers

import (
	"log"
	"os"
	"testing"

	"github.com/blackadress/vaula/models"
	"github.com/blackadress/vaula/utils"
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
	utils.EnsureTableUsuarioExists(a.DB)
	// utils.EnsureTableAlternativaExists(a.DB) TODO
	utils.EnsureTableAlumnoExists(a.DB)
	utils.EnsureTableCursoExists(a.DB)
	utils.EnsureTableExamenExists(a.DB)
	utils.EnsureTablePreguntaExists(a.DB)
	utils.EnsureTableProfesorExists(a.DB)
	utils.EnsureTableTrabajoExists(a.DB)

	code := m.Run()

	// limpiar las tablas de la BD
	// ClearTableAlumno(db) ya reseteado por ClearTableUsuario
	// ClearTableProfesor(db) ya reseteado por ClearTableUsuario
	utils.ClearTableUsuario(a.DB)
	utils.ClearTableCurso(a.DB)
	utils.ClearTableAlternativa(a.DB)
	// ClearTableExamen(db) ya reseteado por ClearTableCurso
	// ClearTablePregunta(db) ya reseteado por ClearTableCurso
	// ClearTableTrabajo(db) ya reseteado por ClearTableCurso
	// ClearTablePreguntaTrabajo(db) ya reseteado por ClearTableCurso
	os.Exit(code)
}

func getTestJWT() models.JWToken {
	usuario := models.User{Username: "prueba", Password: "prueba"}

	err := usuario.GetUserByUsername(a.DB)
	if err != nil {
		log.Fatalf("Error no se encuentra el usuario prueba para autorización, %s", err)
	}

	token, err := usuario.GetJWTForUser()
	if err != nil {
		log.Fatalf("Error inesperado en GetJWTForUser, %s", err)
	}

	return token
}

func ensureAuthorizedUserExists() {
	user := models.User{Username: "prueba", Password: "prueba", Email: "prueba@pru.eba", Activo: true}
	user.Password = hashAndSalt([]byte(user.Password))

	err := user.CreateUser(a.DB)
	if err != nil {
		log.Fatalf("Error en el metodo CreateUser, %s", err)
	}
}
