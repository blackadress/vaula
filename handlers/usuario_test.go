package handlers

// los test necesitan que la aplicacion este funcionando para
// poder probar los JWT

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

func TestEmptyUsuarioTable(t *testing.T) {
	clearTableUsuario()
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/users", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	if body != `[{"id":1,"username":"prueba","password":"","email":"prueba@pru.eba"}]` {
		t.Errorf("Expected an array with one element. Got %#v", body)
	}
}

func TestGetNonExistentUsuario(t *testing.T) {
	clearTableUsuario()
	ensureAuthorizedUserExists()
	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/users/11", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "User not found" {
		t.Errorf(
			"Expected the 'error' key of the response to be set to 'User not found'. Got '%s'",
			m["error"])
	}
}

func TestCreateUser(t *testing.T) {
	clearTableUsuario()

	var jsonStr = []byte(`
	{
		"username": "user_test",
		"password": "1234",
		"email": "user_test@test.ts"
	}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, a)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["username"] != "user_test" {
		t.Errorf("Expected user username to be 'user_test'. Got '%v'", m["username"])
	}

	if m["password"] == "1234" {
		t.Errorf("Expected password to have been hashed, it is still '%v'", m["password"])
	}

	if m["email"] != "user_test@test.ts" {
		t.Errorf("Expected user email to be 'user_test@test.ts'. Got '%v'", m["email"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected user ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetUser(t *testing.T) {
	clearTableUsuario()
	addUsers(1)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/users/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateUser(t *testing.T) {
	clearTableUsuario()
	addUsers(1)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/users/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)
	var originalUser map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalUser)

	var jsonStr = []byte(`{
		"username": "user_test_updated",
		"password": "1234_updated",
		"email": "user_test_updated@test.ts"}`)

	req, _ = http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token_str)
	response = executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalUser["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalUser["id"], m["id"])
	}

	if m["username"] == originalUser["username"] {
		t.Errorf(
			"Expected the username to change from '%v' to '%v'. Got '%v'",
			originalUser["username"],
			m["username"],
			originalUser["username"],
		)
	}

	if m["password"] == originalUser["password"] {
		t.Errorf(
			"Expected the password to change from '%v' to '%v'. Got '%v'",
			originalUser["password"],
			m["password"],
			originalUser["password"],
		)
	}

	if m["email"] == originalUser["email"] {
		t.Errorf(
			"Expected the email to change from '%v', to '%v'. Got '%v'",
			originalUser["email"],
			m["email"],
			originalUser["email"],
		)
	}
}

func TestDeleteUser(t *testing.T) {
	clearTableUsuario()
	addUsers(1)
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

func TestRefreshToken(t *testing.T) {
	clearTableUsuario()
	addUsers(1)
	ensureAuthorizedUserExists()
	//token := getTestJWT()

}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS usuarios
	(
		id INT PRIMARY KEY NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		email TEXT NOT NULL,

		activo BOOLEAN,
		createdAt TIMESTAMPTZ,
		updatedAt TIMESTAMPTZ
	)
`

const userInsertionQuery = `
	INSERT INTO users(username, password, email)
	VALUES('prueba', 'prueba', 'prueba@mail.com')
`

func ensureTableUsuarioExists() {
	if _, err := a.DB.Exec(context.Background(), tableCreationQuery); err != nil {
		log.Printf("TEST: error creando tabla de usuarios: %s", err)
	}
}

func clearTableUsuario() {
	a.DB.Exec(context.Background(), "DELETE FROM users")
	a.DB.Exec(context.Background(), "ALTER SEQUENCE users_id_seq RESTART WITH 1")
}

func addUsers(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec(
			context.Background(),
			`INSERT INTO users(username, password, email)
			VALUES($1, $2, $3)`,
			"user_"+strconv.Itoa(i),
			"pass"+strconv.Itoa(i),
			"em"+strconv.Itoa(i)+"@test.ts",
		)
	}
}
