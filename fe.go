package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// func main() {
// 	//ensureUserExists()
// 	//token := getTestJWT()
// 	//println("this is the token obtained", token)
// 	//getUsersT(token)

// 	fullToken := getFullJWT()
// 	fmt.Printf("%#v\n", fullToken)
// 	println("*****************************************************")

// 	newPair := refreshToken(fullToken.AccessToken)
// 	fmt.Printf("%#v\n", newPair)
// 	println("*****************************************************")

// 	newPair = refreshToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjExMCwiZXhwIjoxNjIxMzc1ODA1fQ.Bx5rWZpiQUkeEvPWaBYu2TNOic7g3DlgNB6Dj-MRe6o")
// 	fmt.Printf("%#v\n", newPair)
// 	println("*****************************************************")

// }

type Temp_jwt struct {
	UserId       int
	AccessToken  string
	RefreshToken string
}

func refreshToken(refreshToken string) Temp_jwt {
	url := "http://localhost:8000/api/refresh"
	req, _ := http.NewRequest(
		"GET",
		url,
		nil)
	req.Header.Set("Refresh", refreshToken)

	client := &http.Client{}
	resp, _ := client.Do(req)

	var jwt Temp_jwt
	body, _ := ioutil.ReadAll(resp.Body)
	println("refresh token raw", string(body))
	println("--------------------------------------")
	json.Unmarshal(body, &jwt)

	return jwt
}

func getFullJWT() Temp_jwt {
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

	return jwt
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
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
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
	resp, err := client.Do(req)
	if err != nil {
		return
	}

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
