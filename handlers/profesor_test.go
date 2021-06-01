package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/blackadress/vaula/utils"
)

func TestEmptyProfesorTable(t *testing.T) {
	utils.ClearTableProfesor(a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/profesores", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	if body != "[]" {
		t.Errorf("Se esperaba un array vacio. Se obtuvo %#v", body)
	}
}

func TestGetNonExistentProfesor(t *testing.T) {
	utils.ClearTableProfesor(a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/profesores/11", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Profesor no encontrado" {
		t.Errorf(
			"Se espera que la key 'error' sea 'Profesor no encontrado'. Got '%s'",
			m["error"])
	}
}

func TestCreateProfesor(t *testing.T) {
	utils.ClearTableUsuario(a.DB)
	utils.ClearTableProfesor(a.DB)
	utils.AddUsers(1, a.DB)
	ensureAuthorizedUserExists()

	var jsonStr = []byte(`
	{
		"nombres": "profesor_test",
		"apellidos": "ap_profesor_test",
		"usuarioId": 1,
		"activo": true
	}`)
	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)
	req, _ := http.NewRequest("POST", "/profesores", bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", token_str)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["nombres"] != "profesor_test" {
		t.Errorf("Expected profesor nombres to be 'profesor_test'. Got '%v'", m["nombres"])
	}

	if m["apellidos"] != "ap_profesor_test" {
		t.Errorf("Expected profesor apellidos to be 'ap_profesor_test'. Got '%v'", m["apellidos"])
	}

	if m["usuarioId"] != 1.0 {
		t.Errorf("Expected profesor usuarioId to be '1'. Got '%v'", m["usuarioId"])
	}

	if m["activo"] == "true" {
		t.Errorf("Expected profesor activo to be true. Got '%v'", m["activo"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected profesor ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetProfesor(t *testing.T) {
	utils.ClearTableUsuario(a.DB)

	utils.ClearTableProfesor(a.DB)
	utils.AddProfesores(1, a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/profesores/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateProfesor(t *testing.T) {
	utils.ClearTableUsuario(a.DB)
	utils.AddUsers(1, a.DB)
	utils.ClearTableProfesor(a.DB)
	utils.AddProfesores(1, a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/profesores/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	var originalProfesor map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProfesor)

	var jsonStr = []byte(`{
		"nombres": "profesor_test_updated",
		"apellidos": "ap_profesor_test_updated",
		"usuarioId": 2,
		"activo": false
	}`)

	req, _ = http.NewRequest("PUT", "/profesores/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token_str)
	response = executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalProfesor["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProfesor["id"], m["id"])
	}

	if m["nombres"] == originalProfesor["nombres"] {
		t.Errorf(
			"Expected the nombres to change from '%v' to '%v'. Got '%v'",
			originalProfesor["nombres"],
			m["nombres"],
			originalProfesor["nombres"],
		)
	}

	if m["apellidos"] == originalProfesor["apellidos"] {
		t.Errorf(
			"Expected the apellidos to change from '%v' to '%v'. Got '%v'",
			originalProfesor["apellidos"],
			m["apellidos"],
			originalProfesor["apellidos"],
		)
	}

	if m["usuarioId"] == originalProfesor["usuarioId"] {
		t.Errorf(
			"Expected the usuarioId to change from '%v', to '%v'. Got '%v'",
			originalProfesor["usuarioId"],
			m["usuarioId"],
			originalProfesor["usuarioId"],
		)
	}

	if m["activo"] == originalProfesor["activo"] {
		t.Errorf(
			"Expected the activo to change from '%v', to '%v'. Got '%v'",
			originalProfesor["activo"],
			m["activo"],
			originalProfesor["activo"],
		)
	}
}

func TestDeleteProfesor(t *testing.T) {
	utils.ClearTableUsuario(a.DB)
	utils.ClearTableProfesor(a.DB)
	utils.AddProfesores(1, a.DB)

	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/profesores/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/profesores/1", nil)
	req.Header.Set("Authorization", token_str)
	response = executeRequest(req, a)
	checkResponseCode(t, http.StatusOK, response.Code)
}
