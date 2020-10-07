package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf(
			"user=%s password=%s dbname=%s sslmode=disable",
			user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.Handle("/users", getUsersHandler(*a)).Methods("GET")
	a.Router.Handle("/users", createUser(*a)).Methods("POST")
	a.Router.Handle("/users/{id:[0-9]+}", getUserByIdHandler(*a)).Methods("GET")
	a.Router.Handle("/users/{id:[0-9]+}", updateUserHandler(*a)).Methods("PUT")
	a.Router.Handle("/users/{id:[0-9]+}", deleteUser(*a)).Methods("DELETE")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
