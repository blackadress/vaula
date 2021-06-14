package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/blackadress/vaula/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

func (a *App) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("GET %s code: %d ERROR: %s", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	u := models.User{ID: id}
	if err := u.GetUser(a.DB); err != nil {
		switch err {
		case pgx.ErrNoRows:
			log.Printf("GET %s code: %d ERROR: %s", r.RequestURI,
				http.StatusNotFound, err.Error())
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			log.Printf("GET %s code: %d ERROR: %s", r.RequestURI,
				http.StatusInternalServerError, err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	log.Printf("GET %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, u)
	return
}

func (a *App) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetUsers(a.DB)

	if err != nil {
		log.Printf("GET %s code: %d ERROR: %s", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("GET %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, users)
	return
}

func (a *App) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var u models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		log.Printf("POST %s code: %d ERROR: %s", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()

	// hashing the password
	u.Password = hashAndSalt([]byte(u.Password))

	if err := u.CreateUser(a.DB); err != nil {
		log.Printf("POST %s code: %d ERROR: %s", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	u.Password = "" // no regresar la password hash en la respuesta

	log.Printf("POST %s code: %d", r.RequestURI, http.StatusCreated)
	respondWithJSON(w, http.StatusCreated, u)
	return
}

func (a *App) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var u models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		log.Printf("PUT %s code: %d ERROR: %s", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()

	u.ID = id

	if err := u.UpdateUser(a.DB); err != nil {
		log.Printf("PUT %s code: %d ERROR: %s", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("PUT %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, u)
	return
}

func (a *App) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	u := models.User{ID: id}
	if err := u.DeleteUser(a.DB); err != nil {
		log.Printf("PUT %s code: %d ERROR: %s", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]int{"exito": 1, "id": u.ID})
	return
}

//
func (a *App) auth(w http.ResponseWriter, r *http.Request) {
	// attackers shouldn't know if a username exists on the DB
	// so we should roughly take the same amount of time
	// either if the user exists or doesn't
	var u models.User
	//josn.Unmarshal(r.Body)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		log.Printf("Formato json invalido")
		log.Printf("POST %s code: %d ERROR: %s", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid user or password")
		return
	}
	//fmt.Printf("user: %s, pass: %s\n", u.Username, u.Password)
	defer r.Body.Close()

	var uFetched models.User
	uFetched.Username = u.Username

	// therefore, this piece of code can't respond without using a 'CompareHashAndPassword',
	// but can't send a response just after checking
	// if the username exists on the DB
	if err := uFetched.GetUserByUsername(a.DB); err != nil {
		log.Printf("No existe usuario en la DB")
		bcrypt.CompareHashAndPassword([]byte(uFetched.Password), []byte("thereIsNoUser"))
		log.Printf("POST %s code: %d ERROR: %s", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid user or password")
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(uFetched.Password), []byte(u.Password))
	if err != nil {
		// Invalid password
		log.Printf("Password invalida")
		log.Printf("POST %s code: %d ERROR: %s", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid user or password")
		return
	} else {
		token, err := uFetched.GetJWTForUser()
		if err != nil {
			// error inesperado loggeado en la capa de modelo
			log.Printf("POST %s code: %d ERROR: %s", r.RequestURI,
				http.StatusBadRequest, err.Error())
			respondWithError(w, http.StatusBadRequest, "Invalid user or password")
			return
		}

		log.Printf("POST %s code: %d", r.RequestURI, http.StatusOK)
		respondWithJSON(w, http.StatusOK, token)
		return
	}
}

func (a *App) refresh(w http.ResponseWriter, r *http.Request) {
	if r.Header["Refresh"] == nil {
		log.Printf("POST %s code: %d ERROR: %s", r.RequestURI,
			http.StatusOK, "Cabecera no contiene 'Refresh' token")
		respondWithError(w, http.StatusBadRequest, "")
		return
	}

	isTokenValid, claims, err := models.ValidateToken(r.Header["Refresh"][0])

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Printf("GET %s code: %d ERROR: %s", r.RequestURI,
				http.StatusBadRequest, err.Error())
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		log.Printf("GET %s code: %d ERROR: %s", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Token Expired")
		return
	}

	if !isTokenValid {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Minute {
		respondWithError(w, http.StatusUnauthorized, "Demasiado pronto para pedir nuevo token")
		return
	}
	// check if userID is in DB
	u := models.User{ID: claims.UserId}
	if err := u.GetUserNoPwd(a.DB); err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	newTknPair, err := u.GetJWTForUser()
	if err != nil {
		log.Printf("%v", err.Error())
		respondWithError(w, http.StatusInternalServerError, "Error generando token")
		return
	}

	respondWithJSON(w, http.StatusOK, newTknPair)
	return
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			parseToken := func(tkn string) (string, error) {

				re := regexp.MustCompile(`Bearer\s(?P<token>.*)`)
				captured := re.FindStringSubmatch(tkn)
				if captured == nil {
					return "", fmt.Errorf("Wrong Authorization header format")
				}
				parsedTkn := captured[1]
				return parsedTkn, nil
			}
			tkn, err := parseToken(r.Header["Authorization"][0])
			if err != nil {
				log.Printf("POST %s code: %d ERROR: %s", r.RequestURI,
					http.StatusBadRequest, err.Error())
				respondWithError(w, http.StatusBadRequest, "Invalid user or password")
			}

			isTokenValid, _, err := models.ValidateToken(tkn)

			if err != nil {
				log.Printf("POST %s code: %d", r.RequestURI, http.StatusBadRequest)
				respondWithError(w, http.StatusBadRequest, "Invalid user or password")
			}

			if isTokenValid {
				endpoint(w, r)
			}
		} else {
			var s string
			for key, val := range r.Header {
				s = fmt.Sprintf("%s=\"%s\"\n", key, val)
			}
			println("no hay ['Authorization'], en los headers ", s)
			log.Printf("POST %s code: %d", r.RequestURI, http.StatusUnauthorized)
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		}
	})
}

func pass(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		endpoint(w, r)
	})
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(hash)
}
