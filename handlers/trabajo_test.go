package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/blackadress/vaula/utils"
)

func TestEmptyTrabajoTable(t *testing.T) {
	utils.ClearTableUsuario(a.DB)
	ensureAuthorizedUserExists()

	utils.ClearTableTrabajo(a.DB)

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/trabajos", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	if body != "[]" {
		t.Errorf("Se esperaba un array vacio. Se obtuvo %#v", body)
	}
}

func TestGetNonExistentTrabajo(t *testing.T) {
	utils.ClearTableUsuario(a.DB)
	utils.ClearTableTrabajo(a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/trabajos/11", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Trabajo no encontrado" {
		t.Errorf(
			"Se espera que la key 'error' sea 'Trabajo no encontrado'. Got '%s'",
			m["error"])
	}
}

func TestCreateTrabajo(t *testing.T) {
	utils.ClearTableTrabajo(a.DB)
	utils.ClearTableUsuario(a.DB)
	ensureAuthorizedUserExists()

	var jsonStr = []byte(`
	{
		"descripcion": "trabajo_desc_test",
		"fechaInicio": "2016-06-22T19:10:25-05:00",
		"fechaFinal": "2016-06-24T19:10:25-05:00",
		"cursoId": 1
	}`)
	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("POST", "/trabajos", bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", token_str)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["descripcion"] != "trabajo_desc_test" {
		t.Errorf("Expected trabajo descripcion to be 'trabajo_desc_test'. Got '%v'", m["descripcion"])
	}

	if m["fechaInicio"] != "2016-06-22T19:10:25-05:00" {
		t.Errorf("Expected fechaInicio to be '2016-06-22T19:10:25-05:00'. Got '%v'", m["fechaInicio"])
	}

	if m["fechaFinal"] != "2016-06-24T19:10:25-05:00" {
		t.Errorf("Expected user fechaFinal to be '2016-06-24T19:10:25-05:00'. Got '%v'", m["fechaFinal"])
	}

	if m["cursoId"] != 1.0 {
		t.Errorf("Expected cursoId to be '1'. Got '%v'", m["cursoId"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected user ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetTrabajo(t *testing.T) {
	utils.ClearTableUsuario(a.DB)
	utils.ClearTableTrabajo(a.DB)
	utils.AddTrabajos(1, a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/trabajos/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateTrabajo(t *testing.T) {
	utils.ClearTableUsuario(a.DB)
	utils.ClearTableCurso(a.DB)
	utils.AddTrabajos(1, a.DB)
	utils.AddCursos(1, a.DB)
	// la funcion add cursos debe ser llamada despues de
	// addTrabajos para que el trabajo generado pueda
	// ser modificado con el id del curso generado luego
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/trabajos/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	var originalTrabajo map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalTrabajo)

	var jsonStr = []byte(`{
		"descripcion": "trabajo_desc_test_updated",
		"fechaInicio": "2016-07-22T19:10:25-05:00",
		"fechaFinal": "2016-07-24T19:10:25-05:00",
		"cursoId": 2,
		"activo": false
	}`)

	req, _ = http.NewRequest("PUT", "/trabajos/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token_str)
	response = executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalTrabajo["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalTrabajo["id"], m["id"])
	}

	if m["descripcion"] == originalTrabajo["descripcion"] {
		t.Errorf(
			"Expected the descripcion to change from '%v' to '%v'. Got '%v'",
			originalTrabajo["descripcion"],
			m["descripcion"],
			originalTrabajo["descripcion"],
		)
	}

	if m["fechaInicio"] == originalTrabajo["fechaInicio"] {
		t.Errorf(
			"Expected the fechaInicio to change from '%v' to '%v'. Got '%v'",
			originalTrabajo["fechaInicio"],
			m["fechaInicio"],
			originalTrabajo["fechaInicio"],
		)
	}

	if m["fechaFinal"] == originalTrabajo["fechaFinal"] {
		t.Errorf(
			"Expected the fechaFinal to change from '%v', to '%v'. Got '%v'",
			originalTrabajo["fechaFinal"],
			m["fechaFinal"],
			originalTrabajo["fechaFinal"],
		)
	}

	if m["cursoId"] == originalTrabajo["cursoId"] {
		t.Errorf(
			"Expected the cursoId to change from '%v', to '%v'. Got '%v'",
			originalTrabajo["cursoId"],
			m["cursoId"],
			originalTrabajo["cursoId"],
		)
	}

	if m["activo"] == originalTrabajo["activo"] {
		t.Errorf(
			"Expected the activo to change from '%v', to '%v'. Got '%v'",
			originalTrabajo["activo"],
			m["activo"],
			originalTrabajo["activo"],
		)
	}
}

func TestDeleteTrabajo(t *testing.T) {
	utils.ClearTableUsuario(a.DB)
	utils.ClearTableTrabajo(a.DB)
	utils.AddTrabajos(1, a.DB)
	ensureAuthorizedUserExists()
	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/trabajos/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/trabajos/1", nil)
	req.Header.Set("Authorization", token_str)
	response = executeRequest(req, a)
	checkResponseCode(t, http.StatusOK, response.Code)
}
