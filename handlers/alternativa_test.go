package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/blackadress/vaula/utils"
)

func TestEmptyAlternativaTable(t *testing.T) {
	utils.ClearTableAlternativa(a.DB)
	utils.ClearTableUsuario(a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/alternativas", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	if body != "[]" {
		t.Errorf("Se esperaba un array vacio. Se obtuvo %#v", body)
	}
}

func TestGetNonExistentAlternativa(t *testing.T) {
	utils.ClearTableAlternativa(a.DB)
	utils.ClearTableUsuario(a.DB)
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
	utils.ClearTableAlternativa(a.DB)
	utils.ClearTableUsuario(a.DB)
	ensureAuthorizedUserExists()

	jsonStr := []byte(`{
		"valor": "val_alt_test",
		"correcto": true,
		"activo": true
	}`)

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

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
	utils.ClearTableAlternativa(a.DB)
	utils.AddAlternativas(1, a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/alternativas/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateAlternativa(t *testing.T) {
	utils.ClearTableAlternativa(a.DB)
	utils.AddAlternativas(1, a.DB)
	utils.ClearTableUsuario(a.DB)
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
		"correcto": false,
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
	utils.ClearTableAlternativa(a.DB)
	utils.AddAlternativas(1, a.DB)
	utils.ClearTableUsuario(a.DB)
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
