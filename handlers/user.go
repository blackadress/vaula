package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/blackadress/vaula/models"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

func HelloFromUserHandlers() {
	fmt.Printf("Hello from handlers\n")
}

func getUserByIdHandler(app App) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			id, err := strconv.Atoi(vars["id"])
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid user ID")
			}

			u := models.User{ID: id}
			if err := u.GetUser(app.DB); err != nil {
				switch err {
				case sql.ErrNoRows:
					respondWithError(w, http.StatusNotFound, "User not found")
				default:
					respondWithError(w, http.StatusInternalServerError, err.Error())
				}
				return
			}
			respondWithJSON(w, http.StatusOK, u)
		})
}

func getUsersHandler(app App) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			users, err := models.GetUsers(app.DB)

			if err != nil {
				fmt.Println("error line 22")
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			respondWithJSON(w, http.StatusOK, users)
		})
}

func createUser(app App) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var u models.User
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&u); err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid payload")
				return
			}
			defer r.Body.Close()

			// hashing the password
			u.Password = hashAndSalt([]byte(u.Password))

			if err := u.CreateUser(app.DB); err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			respondWithJSON(w, http.StatusCreated, u)
		})
}

func updateUserHandler(app App) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
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

			if err := u.UpdateUser(app.DB); err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			respondWithJSON(w, http.StatusOK, u)
		})
}

func deleteUser(app App) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			id, err := strconv.Atoi(vars["id"])
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid user ID")
				return
			}

			u := models.User{ID: id}
			if err := u.DeleteUser(app.DB); err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
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

//func asd(app App) http.Handler {
//return http.HandlerFunc(
//func(w http.ResponseWriter, r *http.Request) {
//asd
//})
//}
