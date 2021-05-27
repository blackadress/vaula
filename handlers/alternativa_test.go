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
	"time"
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

	// checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	if body != "[]" {
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
	if m["error"] != "Alternativa no encontrada" {
		t.Errorf(
			"Se espera que la key 'error' sea 'Alternativa no encontrada'. Got '%s'",
			m["error"])
	}
}

func TestCreateAlternativa(t *testing.T) {
	clearTableAlternativa()
	clearTableUsuario()
	ensureAuthorizedUserExists()

	jsonStr := []byte(`{
		"valor": "val_alt_test",
		"correcto": true,
		"activo": true
	}`)

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)
	fmt.Printf("------------------------------\n")
	fmt.Printf("token: %s\n", token_str)

	req, _ := http.NewRequest("POST", "/alternativas", bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", token_str)

	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["valor"] != "val_alt_test" {
		t.Errorf("Expected user 'valor' to be 'val_alt_test'. Got '%v'", m["valor"])
	}

	if m["correcto"] == "true" {
		t.Errorf("Expected 'correcto' to be 'true'. Got '%#v'", m["correcto"])
	}

	if m["activo"] == "true" {
		t.Errorf("Expected 'activo' to be 'true'. Got '%#v'", m["activo"])
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

	req, _ := http.NewRequest("GET", "/alternativas/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	var originalAlternativa map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalAlternativa)

	var jsonStr = []byte(`{
		"valor": "alt_test_updated",
		"correcto": true,
		"activo": false
	}`)

	req, _ = http.NewRequest("PUT", "/alternativas/1", bytes.NewBuffer(jsonStr))
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

	if m["activo"] == originalAlternativa["activo"] {
		t.Errorf(
			"Expected the activo to change from '%v' to '%v'. Got '%v'",
			originalAlternativa["activo"],
			m["activo"],
			originalAlternativa["activo"],
		)
	}

	if m["updatedAt"] == originalAlternativa["updatedAt"] {
		t.Errorf(
			"Expected the updatedAt to change from '%v' to '%v'. Got '%v'",
			originalAlternativa["updatedAt"],
			m["updatedAt"],
			originalAlternativa["updatedAt"],
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

	req, _ := http.NewRequest("GET", "/alternativas/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/alternativas/1", nil)
	req.Header.Set("Authorization", token_str)
	response = executeRequest(req, a)
	checkResponseCode(t, http.StatusOK, response.Code)
}

const tableAlternativaCreationQuery = `
CREATE TABLE IF NOT EXISTS alternativas
	(
		id SERIAL PRIMARY KEY,
		valor VARCHAR(50) NOT NULL,
		correcto BOOLEAN NOT NULL,

		activo BOOLEAN NOT NULL,
		createdAt TIMESTAMPTZ,
		updatedAt TIMESTAMPTZ
	)
`

// es posible hacer decouple de `a.DB`?
func ensureTableAlternativaExists() {
	_, err := a.DB.Exec(context.Background(), tableAlternativaCreationQuery)
	if err != nil {
		log.Printf("TEST: error creando tabla alternativas: %s", err)
	}
}

func clearTableAlternativa() {
	a.DB.Exec(context.Background(), "DELETE FROM alternativas")
	a.DB.Exec(context.Background(), "ALTER SEQUENCE alternativas_id_seq RESTART WITH 1")
}

func addAlternativas(count int) {
	now := time.Now()
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec(
			context.Background(),
			`INSERT INTO alternativas(valor, correcto, activo, createdAt, updatedAt)
			VALUES($1, $2, $3, $4, $5)`,
			"valor_"+strconv.Itoa(i), i%2 == 1, i%2 == 0, now, now)
	}
}
