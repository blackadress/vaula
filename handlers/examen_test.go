package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/blackadress/vaula/utils"
)

func TestEmptyExamenTable(t *testing.T) {
	utils.ClearTableExamen(a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/examenes", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	if body != "[]" {
		t.Errorf("Se esperaba un array vacio. Se obtuvo %#v", body)
	}
}

func TestGetNonExistentExamen(t *testing.T) {
	utils.ClearTableExamen(a.DB)
	utils.ClearTableUsuario(a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/examenes/11", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Examen no encontrado" {
		t.Errorf(
			"Se espera que la key 'error' sea 'Examen no encontrado'. Got '%s'",
			m["error"])
	}
}

func TestCreateExamen(t *testing.T) {
	utils.ClearTableExamen(a.DB)
	utils.ClearTableUsuario(a.DB)
	utils.ClearTableCurso(a.DB)
	utils.AddCursos(1, a.DB)

	ensureAuthorizedUserExists()
	// si se usa 'Z' para mandar la fecha, es convertida
	// a hora UTC, quizá sea mejor almacenar todo en UTC
	// y convertir al horario de cliente en el lado del cliente?
	// "fechaInicio": "2016-06-22T19:10:25Z",
	// "fechaFinal": "2016-06-24T19:10:25Z",

	var jsonStr = []byte(`
	{
		"nombre": "examen_test",
		"fechaInicio": "2016-06-22T19:10:25-05:00",
		"fechaFinal": "2016-06-24T19:10:25-05:00",
		"cursoId": 1,
		"activo": true
	}`)
	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("POST", "/examenes", bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", token_str)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["nombre"] != "examen_test" {
		t.Errorf("Expected examen nombre to be 'examen_test'. Got '%v'", m["nombre"])
	}

	if m["fechaInicio"] != "2016-06-22T19:10:25-05:00" {
		t.Errorf("Expected examen fechaInicio to be '2016-06-22T19:10:25-05:00'. Got '%v'",
			m["fechaInicio"])
	}

	if m["fechaFinal"] != "2016-06-24T19:10:25-05:00" {
		t.Errorf("Expected examen fechaFinal to be '2016-06-24T19:10:25-05:00'. Got '%v'",
			m["fechaFinal"])
	}

	if m["cursoId"] != 1.0 {
		t.Errorf("Expected examen cursoId to be '1'. Got '%v'", m["cursoId"])
	}

	if m["activo"] != true {
		t.Errorf("Expected examen activo to be 'true'. Got '%v'", m["activo"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected examen ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetExamen(t *testing.T) {
	utils.ClearTableUsuario(a.DB)
	utils.AddExamenes(1, a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/examenes/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateExamen(t *testing.T) {
	utils.ClearTableUsuario(a.DB)
	utils.AddExamenes(1, a.DB)
	utils.AddCursos(1, a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/examenes/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	var originalExamen map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalExamen)

	var jsonStr = []byte(`{
		"nombre": "examen_test_updated",
		"fechaInicio": "2016-06-22T20:10:25+05:00",
		"fechaFinal": "2016-06-22T20:10:25+05:00",
		"cursoId": 2,
		"activo": false
	}`)

	req, _ = http.NewRequest("PUT", "/examenes/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token_str)
	response = executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalExamen["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalExamen["id"], m["id"])
	}

	if m["nombre"] == originalExamen["nombre"] {
		t.Errorf(
			"Expected the nombre to change from '%v' to '%v'. Got '%v'",
			originalExamen["nombre"],
			m["nombre"],
			originalExamen["nombre"],
		)
	}

	if m["fechaInicio"] == originalExamen["fechaInicio"] {
		t.Errorf(
			"Expected the fechaInicio to change from '%v' to '%v'. Got '%v'",
			originalExamen["fechaInicio"],
			m["fechaInicio"],
			originalExamen["fechaInicio"],
		)
	}

	if m["fechaFinal"] == originalExamen["fechaFinal"] {
		t.Errorf(
			"Expected the fechaFinal to change from '%v', to '%v'. Got '%v'",
			originalExamen["fechaFinal"],
			m["fechaFinal"],
			originalExamen["fechaFinal"],
		)
	}

	if m["cursoId"] == originalExamen["cursoId"] {
		t.Errorf(
			"Expected the cursoId to change from '%v', to '%v'. Got '%v'",
			originalExamen["cursoId"],
			m["cursoId"],
			originalExamen["cursoId"],
		)
	}

	if m["activo"] == originalExamen["activo"] {
		t.Errorf(
			"Expected the activo to change from '%v', to '%v'. Got '%v'",
			originalExamen["activo"],
			m["activo"],
			originalExamen["activo"],
		)
	}
}

func TestDeleteExamen(t *testing.T) {
	utils.ClearTableUsuario(a.DB)
	utils.AddExamenes(1, a.DB)
	ensureAuthorizedUserExists()
	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/examenes/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/examenes/1", nil)
	req.Header.Set("Authorization", token_str)
	response = executeRequest(req, a)
	checkResponseCode(t, http.StatusOK, response.Code)
}
