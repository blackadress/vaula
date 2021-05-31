package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/blackadress/vaula/utils"
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
	// utils.ClearTableAlternativa(a.DB)
	// utils.ClearTableUsuario(a.DB)
	// utils.ClearTableAlumno(a.DB)
	// utils.ClearTableCurso(a.DB)
	// utils.ClearTableExamen(a.DB)
	// utils.ClearTablePregunta(a.DB)
	// utils.ClearTableProfesor(a.DB)
	// utils.ClearTableTrabajo(a.DB)
	os.Exit(code)
}

type Temp_jwt struct {
	UserId       int
	AccessToken  string
	RefreshToken string
}

func getTestJWT() Temp_jwt {
	userJson, err := json.Marshal(map[string]string{
		"username": "prueba", "password": "prueba"})
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("%s%s", BASE_URL, "/api/token")

	resp, err := http.Post(url, "application/json",
		bytes.NewBuffer(userJson))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var jwt Temp_jwt
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(res, &jwt)

	return jwt
	//checkResponseCode(t, http.StatusOK, response.Code)

}

func ensureAuthorizedUserExists() {
	var userJson = []byte(`
	{
		"username": "prueba",
		"password": "prueba",
		"email": "prueba@pru.eba"
	}`)
	req, _ := http.NewRequest("POST",
		"http://localhost:8000/users",
		bytes.NewBuffer(userJson))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Algo inesperado paso %s, probablemente el servidor no este activo", err)
	}

	resp.Body.Close()
}
