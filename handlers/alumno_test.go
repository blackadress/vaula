package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testing"

	"github.com/blackadress/vaula/globals"
)

func TestEmptyAlumnoTable(t *testing.T) {
	clearTableAlumno()
	clearTableUsuario()
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/alumnos", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	if body != "" {
		t.Errorf("Se esperaba un array vacio. Se obtuvo %#v", body)
	}
}

func TestGetNonExistentAlumno(t *testing.T) {
	clearTableAlumno()
	clearTableUsuario()
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/alumnos/11", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Alumno not found" {
		t.Errorf(
			"Se espera que la key 'error' sea 'Alumno not found'. Got '%s'",
			m["error"])
	}
}

func TestCreateAlumno(t *testing.T) {
	clearTableAlumno()
	clearTableUsuario()
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	var jsonStr = []byte(`
	{
		"nombres": "nom_al_test",
		"apellidos": "ap_al_test",
		"codigo": "12345678",
		"usuarioId": 1
	}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token_str)

	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["nombres"] != "nom_al_test" {
		t.Errorf("Expected user nombres to be 'nom_al_test'. Got '%v'", m["nombres"])
	}

	if m["apellidos"] == "ap_al_test" {
		t.Errorf("Expected apellidos to be 'ap_al_test'. Got '%v'", m["apellidos"])
	}

	if m["codigo"] != "12345678" {
		t.Errorf("Expected user codigo to be '12345678'. Got '%v'", m["codigo"])
	}

	if m["usuarioId"] != 1.0 {
		t.Errorf("Expected user codigo to be 'user_test@test.ts'. Got '%v'", m["codigo"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected user ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetAlumno(t *testing.T) {
	clearTableUsuario()
	addAlumnos(1)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/users/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateAlumno(t *testing.T) {
	clearTableUsuario()
	addAlumnos(1)
	addUsers(2)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/users/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	var originalAlumno map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalAlumno)

	var jsonStr = []byte(`{
		"nombres": "nom_al_test_updated",
		"apellidos": "ap_al_test_updated",
		"codigo": "11111111",
		"usuarioId": 2,
		"activo": false}`)

	req, _ = http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token_str)
	response = executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalAlumno["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalAlumno["id"], m["id"])
	}

	if m["nombres"] == originalAlumno["nombres"] {
		t.Errorf(
			"Expected the nombres to change from '%v' to '%v'. Got '%v'",
			originalAlumno["nombres"],
			m["nombres"],
			originalAlumno["nombres"],
		)
	}

	if m["apellidos"] == originalAlumno["apellidos"] {
		t.Errorf(
			"Expected the apellidos to change from '%v' to '%v'. Got '%v'",
			originalAlumno["apellidos"],
			m["apellidos"],
			originalAlumno["apellidos"],
		)
	}

	if m["codigo"] == originalAlumno["codigo"] {
		t.Errorf(
			"Expected the codigo to change from '%v', to '%v'. Got '%v'",
			originalAlumno["codigo"],
			m["codigo"],
			originalAlumno["codigo"],
		)
	}

	if m["usuarioId"] == originalAlumno["usuarioId"] {
		t.Errorf(
			"Expected the usuarioId to change from '%v', to '%v'. Got '%v'",
			originalAlumno["usuarioId"],
			m["usuarioId"],
			originalAlumno["usuarioId"],
		)
	}

	if m["activo"] == originalAlumno["activo"] {
		t.Errorf(
			"Expected the activo to change from '%v', to '%v'. Got '%v'",
			originalAlumno["activo"],
			m["activo"],
			originalAlumno["activo"],
		)
	}
}

func TestDeleteAlumno(t *testing.T) {
	clearTableUsuario()
	addAlumnos(1)
	ensureAuthorizedUserExists()
	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/users/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/users/1", nil)
	req.Header.Set("Authorization", token_str)
	response = executeRequest(req, a)
	checkResponseCode(t, http.StatusOK, response.Code)
}

const tableAlumnoCreationQuery = `
CREATE TABLE IF NOT EXISTS alumnos
	(
		id SERIAL,
		apellidos VARCHAR(200) NOT NULL,
		nombres VARCHAR(200) NOT NULL,
		codigo CHAR(8) NOT NULL,
		usuarioId INT REFERENCES usuario(id)

		activo BOOLEAN NOT NULL,
		createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updatedAt TIMESTAMPTZ
	)
`

// es posible hacer decouple de `globals.DB`?
func ensureTableAlumnoExists() {
	_, err := globals.DB.Exec(context.Background(), tableAlumnoCreationQuery)
	if err != nil {
		log.Printf("TEST: error creando tabla alumnos: %s", err)
	}
}

func clearTableAlumno() {
	globals.DB.Exec(context.Background(), "DELETE FROM alumnos")
	globals.DB.Exec(context.Background(), "ALTER SEQUENCE alumnos_id_seq RESTART WITH 1")
}

func addAlumnos(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		globals.DB.Exec(
			context.Background(),
			`INSERT INTO alumnos(valor, correcto, activo)
			VALUES($1, $2, $3)`,
			"valor_"+strconv.Itoa(i),
			i%2 == 1,
			true)
	}
}
