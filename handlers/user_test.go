package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/blackadress/vaula/globals"
	"github.com/joho/godotenv"
)

var a App
var BASE_URL string

func TestMain(m *testing.M) {

	if err := godotenv.Load("../.env"); err != nil {
		log.Print("TEST: no '.env' found")
	}
	BASE_URL = os.Getenv("URL")

	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	ensureTableExists()
	code := m.Run()
	//clearTable()
	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()
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

func TestGetNonExistentUser(t *testing.T) {
	clearTable()
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
	clearTable()

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
	clearTable()
	addUsers(1)
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/users/1", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
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

func TestDeleteProduct(t *testing.T) {
	clearTable()
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

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS users
	(
		id SERIAL,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		email TEXT NOT NULL,
		CONSTRAINT user_pkey PRIMARY KEY (id)
	)
`

const userInsertionQuery = `
	INSERT INTO users(username, password, email)
	VALUES('prueba', 'prueba', 'prueba@mail.com')
`

type Temp_jwt struct {
	UserId      int
	AccessToken string
	Expires     time.Time
}

func ensureTableExists() {
	if _, err := globals.DB.Exec(context.Background(), tableCreationQuery); err != nil {
		log.Printf("TEST: error creando tabla de usuarios: %s", err)
	}
}

func getTestJWT() Temp_jwt {
	userJson, err := json.Marshal(map[string]string{
		"username": "prueba", "password": "prueba"})
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("%s%s", BASE_URL, "/api/token")

	resp, err := http.Post(url, "application/json",
		bytes.NewBuffer(userJson))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var jwt Temp_jwt
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(res, &jwt)

	return jwt
	//checkResponseCode(t, http.StatusOK, response.Code)

}

func ensureAuthorizedUserExists() {
	var userJson = []byte(`
	{
		"username": "prueba",
		"password": "prueba",
		"email": "prueba@pru.eba"
	}`)
	req, _ := http.NewRequest("POST",
		"http://localhost:8000/users",
		bytes.NewBuffer(userJson))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)

	resp.Body.Close()
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))

}

func clearTable() {
	globals.DB.Exec(context.Background(), "DELETE FROM users")
	globals.DB.Exec(context.Background(), "ALTER SEQUENCE users_id_seq RESTART WITH 1")
}

func addUsers(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		globals.DB.Exec(
			context.Background(),
			`INSERT INTO users(username, password, email)
			VALUES($1, $2, $3)`,
			"user_"+strconv.Itoa(i),
			"pass"+strconv.Itoa(i),
			"em"+strconv.Itoa(i)+"@test.ts",
		)
	}
}
