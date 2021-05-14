package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/blackadress/vaula/globals"
	"github.com/blackadress/vaula/models"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

func getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
	}

	u := models.User{ID: id}
	if err := u.GetUser(globals.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetUsers(globals.DB)

	if err != nil {
		fmt.Println("error line 22")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()

	// hashing the password
	u.Password = hashAndSalt([]byte(u.Password))

	if err := u.CreateUser(globals.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, u)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var u models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()

	u.ID = id

	if err := u.UpdateUser(globals.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	u := models.User{ID: id}
	if err := u.DeleteUser(globals.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

//func refresh(w http.ResponseWriter, r *http.Request) {
//if r.Header["Token"] == nil {
//respondWithError(w, http.StatusBadRequest, "")
//return
//}
//claims := &Claims{}
//token, err := jwt.ParseWithClaims(
//r.Header["Token"][0],
//claims,
//func(token *jwt.Token) (interface{}, error) {
//if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//return nil, fmt.Errorf("There was an error while refreshing token")
//}
//return []byte(os.Getenv("SECRET_KEY")), nil
//})

//if err != nil {
//if err == jwt.ErrSignatureInvalid {
//respondWithError(w, http.StatusUnauthorized, err.Error())
//return
//}
//respondWithError(w, http.StatusBadRequest, "Token Expired")
//return
//}

//if !token.Valid {
//respondWithError(w, http.StatusUnauthorized, "Invalid token")
//return
//}

//if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
//respondWithError(w, http.StatusUnauthorized, "Too soon to request a new token")
//return
//}
//expirationTime := time.Now().Add(30 * time.Minute)
//tokenString, err := generateJWT(expirationTime, claims.UserId)
//if err != nil {
//respondWithError(w, http.StatusInternalServerError, err.Error())
//return
//}

//tkn := Token{
//UserId:      claims.UserId,
//AccessToken: tokenString,
//Expires:     expirationTime,
//}
//respondWithJSON(w, http.StatusOK, tkn)
//return
//}

func auth(w http.ResponseWriter, r *http.Request) {
	// attackers shouldn't know if a username exists on the DB
	// so we should roughly take the same amount of time
	// either if the user exists or doesn't
	var u models.User
	//josn.Unmarshal(r.Body)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		log.Printf("Formato json invalido")
		respondWithError(w, http.StatusBadRequest, "Invalid user or password")
		return
	}
	fmt.Printf("user: %s, pass: %s\n", u.Username, u.Password)
	defer r.Body.Close()

	var uFetched models.User
	uFetched.Username = u.Username

	// therefore, this piece of code can't respond without using a 'CompareHashAndPassword',
	// maybe write to Log for internal debugging purposes
	// but can't send a response just after checking
	// if the username exists on the DB
	if err := uFetched.GetUserByUsername(globals.DB); err != nil {
		log.Printf("No existe usuario en la DB")
		bcrypt.CompareHashAndPassword([]byte(uFetched.Password), []byte("thereIsNoUser"))
		respondWithError(w, http.StatusBadRequest, "Invalid user or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(uFetched.Password), []byte(u.Password)); err != nil {
		// Invalid password
		log.Printf("Password invalida")
		respondWithError(w, http.StatusBadRequest, "Invalid user or password")
		return
	} else {
		token, err := uFetched.GetJWTForUser()
		if err != nil {
			// error inesperado loggeado en la capa de modelo
			respondWithError(w, http.StatusBadRequest, "Invalid user or password")
		}

		respondWithJSON(w, http.StatusOK, token)
	}
}

func pass(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		endpoint(w, r)
	})
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			isTokenValid, err := models.ValidateToken(r.Header["Token"][0])

			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Wrong user")
			}

			if isTokenValid {
				endpoint(w, r)
			}
		} else {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		}
	})
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}

	return string(hash)
}
