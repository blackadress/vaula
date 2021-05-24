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

func TestEmptyProfesorTable(t *testing.T) {
	clearTableProfesor()
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/profesores", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	if body != "" {
		t.Errorf("Se esperaba un array vacio. Se obtuvo %#v", body)
	}
}

func TestGetNonExistentProfesor(t *testing.T) {
	clearTableProfesor()
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/asdf/11", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Profesor not found" {
		t.Errorf(
			"Se espera que la key 'error' sea 'Profesor not found'. Got '%s'",
			m["error"])
	}
}

func TestCreateProfesor(t *testing.T) {
	clearTableUsuario()

	var jsonStr = []byte(`
	{
		"username": "user_test",
		"password": "1234",
		"email": "user_test@test.ts"
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

func TestGetProfesor(t *testing.T) {
	clearTableUsuario()
	addProfesores(1)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/users/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateProfesor(t *testing.T) {
	clearTableUsuario()
	addProfesores(1)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/users/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	var originalProfesor map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProfesor)

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

	if m["id"] != originalProfesor["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProfesor["id"], m["id"])
	}

	if m["username"] == originalProfesor["username"] {
		t.Errorf(
			"Expected the username to change from '%v' to '%v'. Got '%v'",
			originalProfesor["username"],
			m["username"],
			originalProfesor["username"],
		)
	}

	if m["password"] == originalProfesor["password"] {
		t.Errorf(
			"Expected the password to change from '%v' to '%v'. Got '%v'",
			originalProfesor["password"],
			m["password"],
			originalProfesor["password"],
		)
	}

	if m["email"] == originalProfesor["email"] {
		t.Errorf(
			"Expected the email to change from '%v', to '%v'. Got '%v'",
			originalProfesor["email"],
			m["email"],
			originalProfesor["email"],
		)
	}
}

func TestDeleteProfesor(t *testing.T) {
	clearTableUsuario()
	addProfesores(1)
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

const tableProfesorCreationQuery = `
CREATE TABLE IF NOT EXISTS profesores
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
func ensureTableProfesorExists() {
	_, err := globals.DB.Exec(context.Background(), tableProfesorCreationQuery)
	if err != nil {
		log.Printf("TEST: error creando tabla profesores: %s", err)
	}
}

func clearTableProfesor() {
	globals.DB.Exec(context.Background(), "DELETE FROM profesores")
	globals.DB.Exec(context.Background(), "ALTER SEQUENCE profesores_id_seq RESTART WITH 1")
}

func addProfesores(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		globals.DB.Exec(
			context.Background(),
			`INSERT INTO profesores(valor, correcto, activo)
			VALUES($1, $2, $3)`,
			"valor_"+strconv.Itoa(i),
			i%2 == 1,
			true)
	}
}
