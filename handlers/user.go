package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/blackadress/vaula/globals"
	"github.com/blackadress/vaula/models"

	jwt "github.com/dgrijalva/jwt-go"
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

type Claims struct {
	UserId int `json:"userId"`
	jwt.StandardClaims
}

type Token struct {
	UserId      int       `json:"userId"`
	AccessToken string    `json:"accessToken"`
	Expires     time.Time `json:"expires"`
}

func refresh(w http.ResponseWriter, r *http.Request) {
	if r.Header["Token"] == nil {
		respondWithError(w, http.StatusBadRequest, "")
		return
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(r.Header["Token"][0], claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error while refreshing token")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusBadRequest, "Token Expired")
		return
	}

	if !token.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		respondWithError(w, http.StatusUnauthorized, "Too soon to request a new token")
		return
	}
	expirationTime := time.Now().Add(30 * time.Minute)
	tokenString, err := generateJWT(expirationTime, claims.UserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tkn := Token{
		UserId:      claims.UserId,
		AccessToken: tokenString,
		Expires:     expirationTime,
	}
	respondWithJSON(w, http.StatusOK, tkn)
	return
}

func auth(w http.ResponseWriter, r *http.Request) {
	var u models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user or password")
		return
	}
	defer r.Body.Close()

	var uFetched models.User
	uFetched.Username = u.Username

	if err := uFetched.GetUserByUsername(globals.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(uFetched.Password), []byte(u.Password)); err != nil {
		// Invalid password
		respondWithError(w, http.StatusBadRequest, "Invalid user or password")
		return
	} else {
		expirationTime := time.Now().Add(time.Minute * 30)
		validToken, err := generateJWT(expirationTime, uFetched.ID)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		token := Token{
			UserId:      uFetched.ID,
			AccessToken: validToken,
			Expires:     expirationTime,
		}
		respondWithJSON(w, http.StatusOK, token)
	}
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}

	return string(hash)
}

func generateJWT(expirationTime time.Time, userId int) (string, error) {
	claims := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// get this from env
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		fmt.Errorf("Something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretKey := []byte(os.Getenv("SECRET_KEY"))
		if r.Header["Token"] != nil {
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(r.Header["Token"][0], claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error while validatin token")
				}
				return secretKey, nil
			})

			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Wrong user")
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		}
	})
}
