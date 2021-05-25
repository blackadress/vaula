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

func TestEmptyAlternativaTable(t *testing.T) {
	clearTableAlternativa()
	clearTableUsuario()
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/alternativas", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	if body != "" {
		t.Errorf("Se esperaba un array vacio. Se obtuvo %#v", body)
	}
}

func TestGetNonExistentAlternativa(t *testing.T) {
	clearTableAlternativa()
	clearTableUsuario()
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/alternativas/11", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Alternativa not found" {
		t.Errorf(
			"Se espera que la key 'error' sea 'Alternativa not found'. Got '%s'",
			m["error"])
	}
}

func TestCreateAlternativa(t *testing.T) {
	clearTableAlternativa()
	clearTableUsuario()
	ensureAuthorizedUserExists()

	var jsonStr = []byte(`
	{
		"valor": "val_alt_test",
		"correcto": true
	}`)

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["valor"] != "val_alt_test" {
		t.Errorf("Expected user 'valor' to be 'val_alt_test'. Got '%v'", m["valor"])
	}

	if m["correcto"] {
		t.Errorf("Expected 'correcto' to be 'true'. Got '%v'", m["correcto"])
	}

	if m["activo"] {
		t.Errorf("Expected 'activo' to be 'true'. Got '%v'", m["password"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected alternativa ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetAlternativa(t *testing.T) {
	clearTableAlternativa()
	addAlternativas(1)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/alternativas/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateAlternativa(t *testing.T) {
	clearTableAlternativa()
	addAlternativas(1)
	clearTableUsuario()
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/users/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	var originalAlternativa map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalAlternativa)

	var jsonStr = []byte(`{
		"valor": "alt_test_updated",
		"correcto": false}`)

	req, _ = http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token_str)
	response = executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalAlternativa["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalAlternativa["id"], m["id"])
	}

	if m["valor"] == originalAlternativa["valor"] {
		t.Errorf(
			"Expected the valor to change from '%v' to '%v'. Got '%v'",
			originalAlternativa["valor"],
			m["valor"],
			originalAlternativa["valor"],
		)
	}

	if m["correcto"] == originalAlternativa["correcto"] {
		t.Errorf(
			"Expected the correcto to change from '%v' to '%v'. Got '%v'",
			originalAlternativa["correcto"],
			m["correcto"],
			originalAlternativa["correcto"],
		)
	}
}

func TestDeleteAlternativa(t *testing.T) {
	clearTableAlternativa()
	addAlternativas(1)
	clearTableUsuario()
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

const tableAlternativaCreationQuery = `
CREATE TABLE IF NOT EXISTS alternativas
	(
		id SERIAL,
		valor VARCHAR(50) NOT NULL,
		correcto BOOLEAN NOT NULL,

		activo BOOLEAN NOT NULL,
		createdAt TIMESTAMPTZ,
		updatedAt TIMESTAMPTZ
	)
`

// es posible hacer decouple de `globals.DB`?
func ensureTableAlternativaExists() {
	_, err := globals.DB.Exec(context.Background(), tableAlternativaCreationQuery)
	if err != nil {
		log.Printf("TEST: error creando tabla alternativas: %s", err)
	}
}

func clearTableAlternativa() {
	globals.DB.Exec(context.Background(), "DELETE FROM alternativas")
	globals.DB.Exec(context.Background(), "ALTER SEQUENCE alternativas_id_seq RESTART WITH 1")
}

func addAlternativas(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		globals.DB.Exec(
			context.Background(),
			`INSERT INTO alternativas(valor, correcto, activo)
			VALUES($1, $2, $3)`,
			"valor_"+strconv.Itoa(i),
			i%2 == 1,
			true)
	}
}
