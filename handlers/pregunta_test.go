package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/blackadress/vaula/utils"
)

func TestEmptyPreguntaTable(t *testing.T) {
	utils.ClearTablePregunta(a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/preguntas", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	if body != "[]" {
		t.Errorf("Se esperaba un array vacio. Se obtuvo %#v", body)
	}
}

func TestGetNonExistentPregunta(t *testing.T) {
	utils.ClearTablePregunta(a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/preguntas/11", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Pregunta no encontrado" {
		t.Errorf(
			"Se espera que la key 'error' sea 'Pregunta no encontrado'. Got '%s'",
			m["error"])
	}
}

func TestCreatePregunta(t *testing.T) {
	utils.ClearTableUsuario(a.DB)
	utils.ClearTableExamen(a.DB)
	utils.AddExamenes(1, a.DB)
	ensureAuthorizedUserExists()

	var jsonStr = []byte(`
	{
		"enunciado": "enunciado_test",
		"examenId": 1
	}`)

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("POST", "/preguntas", bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", token_str)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["enunciado"] != "enunciado_test" {
		t.Errorf("Expected pregunta enunciado to be 'enunciado_test'. Got '%v'", m["enunciado"])
	}

	if m["examenId"] != 1.0 {
		t.Errorf("Expected examenId to be '1'. Got '%v'", m["examenId"])
	}

	if m["activo"] == true {
		t.Errorf("Expected pregunta activo to be 'true'. Got '%v'", m["activo"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected pregunta ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetPregunta(t *testing.T) {
	utils.ClearTableUsuario(a.DB)
	utils.AddPreguntas(1, a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/preguntas/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdatePregunta(t *testing.T) {
	utils.ClearTableUsuario(a.DB)
	utils.AddPreguntas(1, a.DB)
	utils.AddExamenes(1, a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/preguntas/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	var originalPregunta map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalPregunta)

	var jsonStr = []byte(`{
		"enunciado": "enunciado_test_updated",
		"examenId": 2,
		"activo": true
	}`)

	req, _ = http.NewRequest("PUT", "/preguntas/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token_str)
	response = executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalPregunta["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalPregunta["id"], m["id"])
	}

	if m["enunciado"] == originalPregunta["enunciado"] {
		t.Errorf(
			"Expected the enunciado to change from '%v' to '%v'. Got '%v'",
			originalPregunta["enunciado"],
			m["enunciado"],
			originalPregunta["enunciado"],
		)
	}

	if m["examenId"] == originalPregunta["examenId"] {
		t.Errorf(
			"Expected the examenId to change from '%v' to '%v'. Got '%v'",
			originalPregunta["examenId"],
			m["examenId"],
			originalPregunta["examenId"],
		)
	}

	if m["activo"] == originalPregunta["activo"] {
		t.Errorf(
			"Expected the activo to change from '%v', to '%v'. Got '%v'",
			originalPregunta["activo"],
			m["activo"],
			originalPregunta["activo"],
		)
	}
}

func TestDeletePregunta(t *testing.T) {
	utils.ClearTableUsuario(a.DB)
	utils.AddPreguntas(1, a.DB)
	ensureAuthorizedUserExists()
	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/preguntas/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/preguntas/1", nil)
	req.Header.Set("Authorization", token_str)
	response = executeRequest(req, a)
	checkResponseCode(t, http.StatusOK, response.Code)
}
