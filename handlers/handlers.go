package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/blackadress/vaula/globals"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf(
			"user=%s password=%s dbname=%s sslmode=disable",
			user, password, dbname)

	var err error
	globals.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	// auth
	a.Router.HandleFunc("/api/token", auth).Methods("POST")
	a.Router.HandleFunc("/api/refresh", refresh).Methods("GET")

	// users
	a.Router.Handle("/users", isAuthorized(getUsersHandler)).Methods("GET")
	a.Router.Handle("/users", isAuthorized(createUser)).Methods("POST")
	a.Router.Handle("/users/{id:[0-9]+}", isAuthorized(getUserByIdHandler)).Methods("GET")
	a.Router.Handle("/users/{id:[0-9]+}", isAuthorized(updateUserHandler)).Methods("PUT")
	a.Router.Handle("/users/{id:[0-9]+}", isAuthorized(deleteUser)).Methods("DELETE")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
