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
	clearTableUsuario()

	var jsonStr = []byte(`
	{
		"nombres": "nom_al_test",
		"apellidos": "ap_al_test",
		"codigo": "user_test@test.ts"
	}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["username"] != "user_test" {
		t.Errorf("Expected user username to be 'user_test'. Got '%v'", m["username"])
	}

	if m["password"] == "1234" {
		t.Errorf("Expected password to have been hashed, it is still '%v'", m["password"])
	}

	if m["email"] != "user_test@test.ts" {
		t.Errorf("Expected user email to be 'user_test@test.ts'. Got '%v'", m["email"])
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
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/users/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	var originalAlumno map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalAlumno)

	var jsonStr = []byte(`{
		"username": "user_test_updated",
		"password": "1234_updated",
		"email": "user_test_updated@test.ts"}`)

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

	if m["username"] == originalAlumno["username"] {
		t.Errorf(
			"Expected the username to change from '%v' to '%v'. Got '%v'",
			originalAlumno["username"],
			m["username"],
			originalAlumno["username"],
		)
	}

	if m["password"] == originalAlumno["password"] {
		t.Errorf(
			"Expected the password to change from '%v' to '%v'. Got '%v'",
			originalAlumno["password"],
			m["password"],
			originalAlumno["password"],
		)
	}

	if m["email"] == originalAlumno["email"] {
		t.Errorf(
			"Expected the email to change from '%v', to '%v'. Got '%v'",
			originalAlumno["email"],
			m["email"],
			originalAlumno["email"],
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
		valor TEXT NOT NULL,
		correcto BOOLEAN NOT NULL,

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
