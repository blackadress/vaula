package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	token := getTestJWT()
	println("this is the token obtained", token)
	getUsersT(token)
	//ensureUserExists()

	//headerEg := "Bearer woeirji389#$!@asd"
	//re := regexp.MustCompile(`Bearer\b* `)
	//fmt.Printf("%q\n", re.Split(headerEg, 2))
}

type Temp_jwt struct {
	UserId      int
	AccessToken string
	Expires     time.Time
}

func getUsersT(tkn string) {
	bearerToken := fmt.Sprintf("Bearer %s", tkn)
	println("this is the Bearer token: ", bearerToken)
	url := "http://localhost:8000/users"
	req, _ := http.NewRequest(
		"GET",
		url, nil)
	req.Header.Set("Authorization", bearerToken)

	client := &http.Client{}
	resp, _ := client.Do(req)

	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func getTestJWT() string {
	userJson, err := json.Marshal(map[string]string{
		"username": "prueba",
		"password": "prueba",
	})
	if err != nil {
		println(err)
	}

	req, _ := http.NewRequest(
		"POST",
		"http://localhost:8000/api/token",
		bytes.NewBuffer(userJson))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)

	var jwt Temp_jwt
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &jwt)

	return jwt.AccessToken
}

func ensureUserExists() {
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

	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}

/*
	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req.Header.Set("Authorization", token_str)
*/
