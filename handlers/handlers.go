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
	DB     *pgxpool.Pool
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", user, password, dbname)
	fmt.Printf("Connection string %s \n", connectionString)

	var err error
	globals.DB, err = pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		log.Fatalf("Unable to connect to database: %v", err)
		os.Exit(1)
	}
	log.Print("Si conecta con db")

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	// auth
	a.Router.HandleFunc("/api/token", auth).Methods("POST")
	a.Router.HandleFunc("/api/refresh", refresh).Methods("GET")

	// users
	a.Router.Handle("/users/{id:[0-9]+}", isAuthorized(getUserByIdHandler)).Methods("GET")
	a.Router.Handle("/users", isAuthorized(getUsersHandler)).Methods("GET")
	a.Router.Handle("/users", pass(createUserHandler)).Methods("POST")
	a.Router.Handle("/users/{id:[0-9]+}", isAuthorized(updateUserHandler)).Methods("PUT")
	a.Router.Handle("/users/{id:[0-9]+}", isAuthorized(deleteUserHandler)).Methods("DELETE")

	// alternativas
	a.Router.Handle("/alternativas/{id:[0-9]+}", isAuthorized(getAlternativaByIdHandler)).Methods("GET")
	a.Router.Handle("/alternativas", isAuthorized(getAlternativasHandler)).Methods("GET")
	a.Router.Handle("/alternativas", isAuthorized(createAlternativaHandler)).Methods("POST")
	a.Router.Handle("/alternativas/{id:[0-9]+}", isAuthorized(updateAlternativaHandler)).Methods("PUT")
	a.Router.Handle("/alternativas/{id:[0-9]+}", isAuthorized(deleteAlternativaHandler)).Methods("DELETE")

	// alumno
	a.Router.Handle("/alumnos/{id:[0-9]+}", isAuthorized(getAlumnoByIdHandler)).Methods("GET")
	a.Router.Handle("/alumnos", isAuthorized(getAlumnosHandler)).Methods("GET")
	a.Router.Handle("/alumnos", isAuthorized(createAlumnoHandler)).Methods("POST")
	a.Router.Handle("/alumnos/{id:[0-9]+}", isAuthorized(updateAlumnoHandler)).Methods("PUT")
	a.Router.Handle("/alumnos/{id:[0-9]+}", isAuthorized(deleteAlumnoHandler)).Methods("DELETE")

	// curso
	a.Router.Handle("/cursos/{id:[0-9]+}", isAuthorized(getCursoByIdHandler)).Methods("GET")
	a.Router.Handle("/cursos", isAuthorized(getCursosHandler)).Methods("GET")
	a.Router.Handle("/cursos", isAuthorized(createCursoHandler)).Methods("POST")
	a.Router.Handle("/cursos/{id:[0-9]+}", isAuthorized(updateCursoHandler)).Methods("PUT")
	a.Router.Handle("/cursos/{id:[0-9]+}", isAuthorized(deleteCursoHandler)).Methods("DELETE")

	// examen
	a.Router.Handle("/examenes{id:[0-9]+}", isAuthorized(getExamenByIdHandler)).Methods("GET")
	a.Router.Handle("/examenes", isAuthorized(getExamenesHandler)).Methods("GET")
	a.Router.Handle("/examenes", isAuthorized(createExamenHandler)).Methods("POST")
	a.Router.Handle("/examenes{id:[0-9]+}", isAuthorized(updateExamenHandler)).Methods("PUT")
	a.Router.Handle("/examenes{id:[0-9]+}", isAuthorized(deleteExamenHandler)).Methods("DELETE")

	// pregunta
	a.Router.Handle("/preguntas/{id:[0-9]+}", isAuthorized(getPreguntaByIdHandler)).Methods("GET")
	a.Router.Handle("/preguntas", isAuthorized(getPreguntasHandler)).Methods("GET")
	a.Router.Handle("/preguntas", isAuthorized(createPreguntaHandler)).Methods("POST")
	a.Router.Handle("/preguntas/{id:[0-9]+}", isAuthorized(updatePreguntaHandler)).Methods("PUT")
	a.Router.Handle("/preguntas/{id:[0-9]+}", isAuthorized(deletePreguntaHandler)).Methods("DELETE")

	// profesor
	a.Router.Handle("/profesores/{id:[0-9]+}", isAuthorized(getProfesorByIdHandler)).Methods("GET")
	a.Router.Handle("/profesores", isAuthorized(getProfesoresHandler)).Methods("GET")
	a.Router.Handle("/profesores", isAuthorized(createProfesorHandler)).Methods("POST")
	a.Router.Handle("/profesores/{id:[0-9]+}", isAuthorized(updateProfesorHandler)).Methods("PUT")
	a.Router.Handle("/profesores/{id:[0-9]+}", isAuthorized(deleteProfesorHandler)).Methods("DELETE")

	// trabajo
	a.Router.Handle("/trabajos/{id:[0-9]+}", isAuthorized(getTrabajoByIdHandler)).Methods("GET")
	a.Router.Handle("/trabajos", isAuthorized(getTrabajosHandler)).Methods("GET")
	a.Router.Handle("/trabajos", isAuthorized(createTrabajoHandler)).Methods("POST")
	a.Router.Handle("/trabajos/{id:[0-9]+}", isAuthorized(updateTrabajoHandler)).Methods("PUT")
	a.Router.Handle("/trabajos/{id:[0-9]+}", isAuthorized(deleteTrabajoHandler)).Methods("DELETE")

}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
