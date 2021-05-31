// package main

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/blackadress/vaula/models"
// 	"github.com/jackc/pgx/v4/pgxpool"
// 	"github.com/joho/godotenv"
// )

// func init() {
// 	if err := godotenv.Load(); err != nil {
// 		log.Print("No '.env' found")
// 	}
// }

// func main() {
// 	ensureUserExists()
// 	token := getTestJWT()
// 	fmt.Printf("this is the token obtained '%s'\n", token)
// 	getUsersT(token)

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

// test handlers
// func main() {
// 	// ensureUserExists()
// 	token := getTestJWT()
// 	fmt.Printf("this is the token obtained '%s'\n", token)
// 	// getUsersT(token)
// 	createAltT(token)
// }

// test models
// func main() {
// 	alt := models.Alternativa{Valor: "alt_prueba", Correcto: true}
// 	user := os.Getenv("APP_DB_USERNAME")
// 	password := os.Getenv("APP_DB_PASSWORD")
// 	dbname := os.Getenv("APP_DB_NAME")
// 	// x := 5

// 	connectionString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", user, password, dbname)
// 	db, err := pgxpool.Connect(context.Background(), connectionString)
// 	if err != nil {
// 		println("no se conecto con la base de datos ", err)
// 	} else {
// 		// row := test.QueryRow(context.Background(), `select * from usuarios where id=$1`, x)
// 		// alt.CreateAlternativa(db)
// 		// fmt.Printf("%#v", alt)

// 		alt = models.Alternativa{ID: 2}
// 		alt.GetAlternativa(db)
// 		fmt.Printf("%#v", alt)
// 		db.Close()
// 	}
// }

// type Temp_jwt struct {
// 	UserId       int
// 	AccessToken  string
// 	RefreshToken string
// }

// func refreshToken(refreshToken string) Temp_jwt {
// 	url := "http://localhost:8000/api/refresh"
// 	req, _ := http.NewRequest(
// 		"GET",
// 		url,
// 		nil)
// 		req.Header.Set("Refresh", refreshToken)

// 		client := &http.Client{}
// 		resp, _ := client.Do(req)

// 		var jwt Temp_jwt
// 		body, _ := ioutil.ReadAll(resp.Body)
// 		println("refresh token raw", string(body))
// 		println("--------------------------------------")
// 		json.Unmarshal(body, &jwt)

// 		return jwt
// 	}

// 	func getFullJWT() Temp_jwt {
// 		userJson, err := json.Marshal(map[string]string{
// 			"username": "prueba",
// 			"password": "prueba",
// 		})
// 		if err != nil {
// 			println(err)
// 		}

// 		req, _ := http.NewRequest(
// 			"POST",
// 			"http://localhost:8000/api/token",
// 			bytes.NewBuffer(userJson))
// 			req.Header.Set("Content-Type", "application/json")

// 			client := &http.Client{}
// 			resp, _ := client.Do(req)

// 			var jwt Temp_jwt
// 			body, _ := ioutil.ReadAll(resp.Body)
// 			json.Unmarshal(body, &jwt)

// 			return jwt
// 		}

// 		func getUsersT(tkn string) {
// 			bearerToken := fmt.Sprintf("Bearer %s", tkn)
// 			println("this is the Bearer token: ", bearerToken)
// 			url := "http://localhost:8000/users"
// 			req, _ := http.NewRequest(
// 				"GET",
// 				url, nil)
// 				req.Header.Set("Authorization", bearerToken)

// 				client := &http.Client{}
// 				resp, err := client.Do(req)
// 				if err != nil {
// 					return
// 				}

// 				defer resp.Body.Close()
// 				fmt.Println("response Status:", resp.Status)
// 				fmt.Println("response Headers:", resp.Header)
// 				//body, _ := ioutil.ReadAll(resp.Body)
// 				//fmt.Println("response Body:", string(body))
// 			}

// 			func createAltT(tkn string) {
// 				bearerToken := fmt.Sprintf("Bearer %s", tkn)
// 				url := "http://localhost:8000/alternativas"
// 				jsonStr := []byte(`
// 				{
// 					"valor": "alt_prueba",
// 					"correcto": true,
// 					"activo": true
// 				}`)

// 				req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
// 				req.Header.Set("Authorization", bearerToken)

// 				client := &http.Client{}
// 				resp, err := client.Do(req)
// 				if err != nil {
// 					return
// 				}

// 				defer resp.Body.Close()
// 				fmt.Println("response Status:", resp.Status)
// 				fmt.Println("response Headers:", resp.Header)
// 			}

// 			func getTestJWT() string {
// 				userJson, err := json.Marshal(map[string]string{
// 					"username": "prueba",
// 					"password": "prueba",
// 				})
// 				if err != nil {
// 					println(err)
// 				}

// 				req, _ := http.NewRequest(
// 					"POST",
// 					"http://localhost:8000/api/token",
// 					bytes.NewBuffer(userJson))
// 					req.Header.Set("Content-Type", "application/json")

// 					client := &http.Client{}
// 					resp, _ := client.Do(req)

// 					var jwt Temp_jwt
// 					body, _ := ioutil.ReadAll(resp.Body)
// 					json.Unmarshal(body, &jwt)

// 					return jwt.AccessToken
// 				}

// 				func ensureUserExists() {
// 					var userJson = []byte(`
// 					{
// 						"username": "prueba",
// 						"password": "prueba",
// 						"email": "prueba@pru.eba"
// 					}`)
// 					req, _ := http.NewRequest("POST",
// 					"http://localhost:8000/users",
// 					bytes.NewBuffer(userJson))
// 					req.Header.Set("Content-Type", "application/json")

// 					client := &http.Client{}
// 					resp, err := client.Do(req)
// 					if err != nil {
// 						return
// 					}

// 					defer resp.Body.Close()
// 					fmt.Println("response Status:", resp.Status)
// 					fmt.Println("response Headers:", resp.Header)
// 					body, _ := ioutil.ReadAll(resp.Body)
// 					fmt.Println("response Body:", string(body))

// 				}

// 				/*
// 				token := getTestJWT()
// 				token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

// 				req.Header.Set("Authorization", token_str)
// 				*/
