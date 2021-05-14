package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/blackadress/vaula/globals"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

type App struct {
	Router *mux.Router
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", user, password, dbname)
	fmt.Printf("Connection string %s \n", connectionString)

	var err error
	globals.DB, err = pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		log.Print(err)
		os.Exit(1)
	}
	log.Print("Si conecta con db")

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	// auth
	a.Router.HandleFunc("/api/token", auth).Methods("POST")
	//a.Router.HandleFunc("/api/refresh", refresh).Methods("GET")

	// users
	a.Router.Handle("/users", isAuthorized(getUsersHandler)).Methods("GET")
	a.Router.Handle("/users", pass(createUser)).Methods("POST")
	a.Router.Handle("/users/{id:[0-9]+}", isAuthorized(getUserByIdHandler)).Methods("GET")
	a.Router.Handle("/users/{id:[0-9]+}", isAuthorized(updateUserHandler)).Methods("PUT")
	a.Router.Handle("/users/{id:[0-9]+}", isAuthorized(deleteUser)).Methods("DELETE")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
