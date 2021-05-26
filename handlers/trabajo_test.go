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

	"github.com/blackadress/vaula/globals"
)

func TestEmptyTrabajoTable(t *testing.T) {
	clearTableUsuario()
	ensureAuthorizedUserExists()

	clearTableTrabajo()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/trabajos", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	if body != "" {
		t.Errorf("Se esperaba un array vacio. Se obtuvo %#v", body)
	}
}

func TestGetNonExistentTrabajo(t *testing.T) {
	clearTableUsuario()
	ensureAuthorizedUserExists()

	clearTableTrabajo()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/trabajos/11", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Trabajo not found" {
		t.Errorf(
			"Se espera que la key 'error' sea 'Trabajo not found'. Got '%s'",
			m["error"])
	}
}

func TestCreateTrabajo(t *testing.T) {
	clearTableUsuario()
	clearTableTrabajo()

	var jsonStr = []byte(`
	{
		"descripcion": "trabajo_desc_test",
		"fechaInicio": "2016-06-22 19:10:25-05",
		"fechaFinal": "2016-06-24 19:10:25-05",
		"cursoId": 1
	}`)
	req, _ := http.NewRequest("POST", "/trabajos", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["descripcion"] != "trabajo_desc_test" {
		t.Errorf("Expected trabajo descripcion to be 'trabajo_desc_test'. Got '%v'", m["descripcion"])
	}

	if m["fechaInicio"] == "2016-06-22 19:10:25-05" {
		t.Errorf("Expected fechaInicio to be '2016-06-22 19:10:25-05'. Got '%v'", m["fechaInicio"])
	}

	if m["fechaFinal"] != "2016-06-24 19:10:25-05" {
		t.Errorf("Expected user fechaFinal to be '2016-06-24 19:10:25-05'. Got '%v'", m["fechaFinal"])
	}

	if m["cursoId"] != 1.0 {
		t.Errorf("Expected cursoId to be '1'. Got '%v'", m["cursoId"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected user ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetTrabajo(t *testing.T) {
	clearTableUsuario()
	clearTableTrabajo()
	addTrabajos(1)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/trabajos/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateTrabajo(t *testing.T) {
	clearTableUsuario()
	ensureAuthorizedUserExists()
	// la funcion add cursos debe ser llamada despues de
	// addTrabajos para que el trabajo generado pueda
	// ser modificado con el id del curso generado luego
	clearTableCurso()
	clearTableTrabajo()
	addTrabajos(1)
	addCursos(1)

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/trabajos/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	var originalTrabajo map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalTrabajo)

	var jsonStr = []byte(`{
		"descripcion": "trabajo_desc_test_updated",
		"fechaInicio": "2016-06-22 19:10:25-05_updated",
		"fechaFinal": "trabajo_desc_test_updated@test.ts"
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
	clearTableUsuario()
	ensureAuthorizedUserExists()
	clearTableTrabajo()
	addTrabajos(1)
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

const tableTrabajoCreationQuery = `
CREATE TABLE IF NOT EXISTS trabajos
	(
		id INT PRIMARY KEY,
		descripcion TEXT NOT NULL,
		fechaInicio TIMESTAMPTZ NOT NULL,
		fechaFinal TIMESTAMPTZ NOT NULL,
		cursoId INT REFERENCES cursos(id),

		activo BOOLEAN NOT NULL,
		createdAt TIMESTAMPTZ NOT NULL,
		updatedAt TIMESTAMPTZ NOT NULL
	)
`

// es posible hacer decouple de `a.DB`?
func ensureTableTrabajoExists() {
	ensureTableCursoExists()
	_, err := a.DB.Exec(context.Background(), tableTrabajoCreationQuery)
	if err != nil {
		log.Printf("TEST: error creando tabla trabajos: %s", err)
	}
}

func clearTableTrabajo() {
	a.DB.Exec(context.Background(), "DELETE FROM trabajos")
	a.DB.Exec(context.Background(), "ALTER SEQUENCE trabajos_id_seq RESTART WITH 1")
}

func addTrabajos(count int) {
	clearTableCurso()
	addCursos(count)
	now := time.Now()

	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec(
			context.Background(),
			`INSERT INTO trabajos(descripcion, fechaInicio,
			fechaFinal, cursoId, activo, createdAt, updatedAt)
			VALUES($1, $2, $3, $4, $5, $6, $7)`,
			"descripcion_trabajo_test_"+strconv.Itoa(i),
			"2016-06-25 19:10:25-05", "2016-06-26 19:10:25-05",
			i+1, i%2 == 1, now, now)
	}
}
