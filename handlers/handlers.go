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
	//a.Router.HandleFunc("/users", getUsersHandler).Methods("GET")
	// auth
	a.Router.HandleFunc("/api/token", auth).Methods("GET")
	a.Router.HandleFunc("/users", getUsersHandler).Methods("GET")
	a.Router.HandleFunc("/users", createUser).Methods("POST")
	a.Router.HandleFunc("/users/{id:[0-9]+}", getUserByIdHandler).Methods("GET")
	a.Router.HandleFunc("/users/{id:[0-9]+}", updateUserHandler).Methods("PUT")
	a.Router.HandleFunc("/users/{id:[0-9]+}", deleteUser).Methods("DELETE")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
