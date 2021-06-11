package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/blackadress/vaula/utils"
)

func TestEmptyCursoTable(t *testing.T) {
	utils.ClearTableCurso(a.DB)
	utils.ClearTableUsuario(a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/cursos", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	if body != "[]" {
		t.Errorf("Se esperaba un array vacio. Se obtuvo %#v", body)
	}
}

func TestGetNonExistentCurso(t *testing.T) {
	utils.ClearTableCurso(a.DB)
	utils.ClearTableUsuario(a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/cursos/11", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Curso no encontrado" {
		t.Errorf(
			"Se espera que la key 'error' sea 'Curso no encontrado'. Got '%s'",
			m["error"])
	}
}

func TestCreateCurso(t *testing.T) {
	utils.ClearTableCurso(a.DB)
	utils.ClearTableUsuario(a.DB)
	ensureAuthorizedUserExists()

	var jsonStr = []byte(`
	{
		"nombre": "curso_test",
		"siglas": "CS-TS-123",
		"silabo": "silabo_test",
		"semestre": "TS-2021-II"
	}`)

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("POST", "/cursos", bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", token_str)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["nombre"] != "curso_test" {
		t.Errorf("Expected user nombre to be 'curso_test'. Got '%v'", m["nombre"])
	}

	if m["siglas"] != "CS-TS-123" {
		t.Errorf("Expected siglas to be 'CS-TS-123'. Got '%v'", m["siglas"])
	}

	if m["silabo"] != "silabo_test" {
		t.Errorf("Expected silabo to be 'silabo_test'. Got '%v'", m["silabo"])
	}

	if m["semestre"] != "TS-2021-II" {
		t.Errorf("Expected semestre to be 'TS-2021-II'. Got '%v'", m["semestre"])
	}

	if m["activo"] == true {
		t.Errorf("Expected activo to be 'true'. Got '%v'", m["activo"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected user ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetCurso(t *testing.T) {
	utils.ClearTableCurso(a.DB)
	utils.ClearTableUsuario(a.DB)
	utils.AddCursos(1, a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/cursos/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateCurso(t *testing.T) {
	utils.ClearTableCurso(a.DB)
	utils.ClearTableUsuario(a.DB)
	utils.AddCursos(1, a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/cursos/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	var originalCurso map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalCurso)

	var jsonStr = []byte(`{
		"nombre": "curso_test_updated",
		"siglas": "SIG-UPD",
		"silabo": "silabo_test_upd",
		"semestre": "UP-3030",
		"activo": false}`)

	req, _ = http.NewRequest("PUT", "/cursos/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token_str)
	response = executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalCurso["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalCurso["id"], m["id"])
	}

	if m["nombre"] == originalCurso["nombre"] {
		t.Errorf(
			"Expected the nombre to change from '%v' to '%v'. Got '%v'",
			originalCurso["nombre"],
			m["nombre"],
			originalCurso["nombre"],
		)
	}

	if m["siglas"] == originalCurso["siglas"] {
		t.Errorf(
			"Expected the siglas to change from '%v' to '%v'. Got '%v'",
			originalCurso["siglas"],
			m["siglas"],
			originalCurso["siglas"],
		)
	}

	if m["silabo"] == originalCurso["silabo"] {
		t.Errorf(
			"Expected the silabo to change from '%v', to '%v'. Got '%v'",
			originalCurso["silabo"],
			m["silabo"],
			originalCurso["silabo"],
		)
	}

	if m["semestre"] == originalCurso["semestre"] {
		t.Errorf(
			"Expected the semestre to change from '%v', to '%v'. Got '%v'",
			originalCurso["semestre"],
			m["semestre"],
			originalCurso["semestre"],
		)
	}

	if m["activo"] == originalCurso["activo"] {
		t.Errorf(
			"Expected the activo to change from '%v', to '%v'. Got '%v'",
			originalCurso["activo"],
			m["activo"],
			originalCurso["activo"],
		)
	}
}

func TestDeleteCurso(t *testing.T) {
	utils.ClearTableCurso(a.DB)
	utils.ClearTableUsuario(a.DB)
	utils.AddCursos(1, a.DB)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/cursos/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/cursos/1", nil)
	req.Header.Set("Authorization", token_str)
	response = executeRequest(req, a)
	checkResponseCode(t, http.StatusOK, response.Code)
}
