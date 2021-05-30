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

func TestEmptyCursoTable(t *testing.T) {
	clearTableCurso()
	clearTableUsuario()
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/cursos", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	if body != "" {
		t.Errorf("Se esperaba un array vacio. Se obtuvo %#v", body)
	}
}

func TestGetNonExistentCurso(t *testing.T) {
	clearTableCurso()
	clearTableUsuario()
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
	clearTableCurso()
	clearTableUsuario()
	ensureAuthorizedUserExists()

	var jsonStr = []byte(`
	{
		"nombre": "curso_test",
		"siglas": "CS-TS-123",
		"silabo": "silabo_test",
		"semestre": "TS-2021-II"
	}`)
	req, _ := http.NewRequest("POST", "/cursos", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["nombre"] != "curso_test" {
		t.Errorf("Expected user nombre to be 'curso_test'. Got '%v'", m["nombre"])
	}

	if m["siglas"] == "CS-TS-123" {
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
	clearTableCurso()
	clearTableUsuario()
	ensureAuthorizedUserExists()
	addCursos(1)

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/cursos/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateCurso(t *testing.T) {
	clearTableCurso()
	clearTableUsuario()
	ensureAuthorizedUserExists()

	addCursos(1)

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
	clearTableCurso()
	clearTableUsuario()
	ensureAuthorizedUserExists()

	addCursos(1)
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

const tableCursoCreationQuery = `
CREATE TABLE IF NOT EXISTS cursos
	(
		id INT PRIMARY KEY,
		nombre VARCHAR(200) NOT NULL,
		siglas VARCHAR(20) NOT NULL,
		silabo VARCHAR(200) NOT NULL,
		semestre VARCHAR(20) NOT NULL,

		activo BOOLEAN NOT NULL,
		createdAt TIMESTAMPTZ,
		updatedAt TIMESTAMPTZ
	)
`

// es posible hacer decouple de `a.DB`?
func ensureTableCursoExists() {
	_, err := a.DB.Exec(context.Background(), tableCursoCreationQuery)
	if err != nil {
		log.Printf("TEST: error creando tabla cursoss: %s", err)
	}
}

func clearTableCurso() {
	a.DB.Exec(context.Background(), "DELETE FROM cursos")
	a.DB.Exec(context.Background(), "ALTER SEQUENCE cursos_id_seq RESTART WITH 1")
}

func addCursos(count int) {
	if count < 1 {
		count = 1
	}
	now := time.Now()

	for i := 0; i < count; i++ {
		semestre := fmt.Sprintf("%.20s", "semestre_"+strconv.Itoa(i))
		a.DB.Exec(
			context.Background(),
			`INSERT INTO cursos(nombre, siglas, silabo, semestre, activo, createdAt, updatedAt)
			VALUES($1, $2, $3, $4, $5, $6, $7)`,
			"curso_test_"+strconv.Itoa(i),
			"TS-0"+strconv.Itoa(i),
			"silabo_test_"+strconv.Itoa(i),
			semestre, i%2 == 0, now, now)
	}
}
