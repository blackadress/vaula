package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

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
	a.DB, err = pgxpool.Connect(context.Background(), connectionString)
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
	a.Router.HandleFunc("/api/token", a.auth).Methods("POST")
	a.Router.HandleFunc("/api/refresh", a.refresh).Methods("GET")

	// users
	a.Router.Handle("/users/{id:[0-9]+}", isAuthorized(a.getUserByIdHandler)).Methods("GET")
	a.Router.Handle("/users", isAuthorized(a.getUsersHandler)).Methods("GET")
	a.Router.Handle("/users", pass(a.createUserHandler)).Methods("POST")
	a.Router.Handle("/users/{id:[0-9]+}", isAuthorized(a.updateUserHandler)).Methods("PUT")
	a.Router.Handle("/users/{id:[0-9]+}", isAuthorized(a.deleteUserHandler)).Methods("DELETE")

	// alternativas
	a.Router.Handle("/alternativas/{id:[0-9]+}", isAuthorized(a.getAlternativaByIdHandler)).Methods("GET")
	a.Router.Handle("/alternativas", isAuthorized(a.getAlternativasHandler)).Methods("GET")
	a.Router.Handle("/alternativas", isAuthorized(a.createAlternativaHandler)).Methods("POST")
	a.Router.Handle("/alternativas/{id:[0-9]+}", isAuthorized(a.updateAlternativaHandler)).Methods("PUT")
	a.Router.Handle("/alternativas/{id:[0-9]+}", isAuthorized(a.deleteAlternativaHandler)).Methods("DELETE")

	// alumno
	a.Router.Handle("/alumnos/{id:[0-9]+}", isAuthorized(a.getAlumnoByIdHandler)).Methods("GET")
	a.Router.Handle("/alumnos", isAuthorized(a.getAlumnosHandler)).Methods("GET")
	a.Router.Handle("/alumnos", isAuthorized(a.createAlumnoHandler)).Methods("POST")
	a.Router.Handle("/alumnos/{id:[0-9]+}", isAuthorized(a.updateAlumnoHandler)).Methods("PUT")
	a.Router.Handle("/alumnos/{id:[0-9]+}", isAuthorized(a.deleteAlumnoHandler)).Methods("DELETE")

	// curso
	a.Router.Handle("/cursos/{id:[0-9]+}", isAuthorized(a.getCursoByIdHandler)).Methods("GET")
	a.Router.Handle("/cursos", isAuthorized(a.getCursosHandler)).Methods("GET")
	a.Router.Handle("/cursos", isAuthorized(a.createCursoHandler)).Methods("POST")
	a.Router.Handle("/cursos/{id:[0-9]+}", isAuthorized(a.updateCursoHandler)).Methods("PUT")
	a.Router.Handle("/cursos/{id:[0-9]+}", isAuthorized(a.deleteCursoHandler)).Methods("DELETE")

	// examen
	a.Router.Handle("/examenes/{id:[0-9]+}", isAuthorized(a.getExamenByIdHandler)).Methods("GET")
	a.Router.Handle("/examenes", isAuthorized(a.getExamenesHandler)).Methods("GET")
	a.Router.Handle("/examenes", isAuthorized(a.createExamenHandler)).Methods("POST")
	a.Router.Handle("/examenes/{id:[0-9]+}", isAuthorized(a.updateExamenHandler)).Methods("PUT")
	a.Router.Handle("/examenes/{id:[0-9]+}", isAuthorized(a.deleteExamenHandler)).Methods("DELETE")

	// pregunta
	a.Router.Handle("/preguntas/{id:[0-9]+}", isAuthorized(a.getPreguntaByIdHandler)).Methods("GET")
	a.Router.Handle("/preguntas", isAuthorized(a.getPreguntasHandler)).Methods("GET")
	a.Router.Handle("/preguntas", isAuthorized(a.createPreguntaHandler)).Methods("POST")
	a.Router.Handle("/preguntas/{id:[0-9]+}", isAuthorized(a.updatePreguntaHandler)).Methods("PUT")
	a.Router.Handle("/preguntas/{id:[0-9]+}", isAuthorized(a.deletePreguntaHandler)).Methods("DELETE")

	// profesor
	a.Router.Handle("/profesores/{id:[0-9]+}", isAuthorized(a.getProfesorByIdHandler)).Methods("GET")
	a.Router.Handle("/profesores", isAuthorized(a.getProfesoresHandler)).Methods("GET")
	a.Router.Handle("/profesores", isAuthorized(a.createProfesorHandler)).Methods("POST")
	a.Router.Handle("/profesores/{id:[0-9]+}", isAuthorized(a.updateProfesorHandler)).Methods("PUT")
	a.Router.Handle("/profesores/{id:[0-9]+}", isAuthorized(a.deleteProfesorHandler)).Methods("DELETE")

	// trabajo
	a.Router.Handle("/trabajos/{id:[0-9]+}", isAuthorized(a.getTrabajoByIdHandler)).Methods("GET")
	a.Router.Handle("/trabajos", isAuthorized(a.getTrabajosHandler)).Methods("GET")
	a.Router.Handle("/trabajos", isAuthorized(a.createTrabajoHandler)).Methods("POST")
	a.Router.Handle("/trabajos/{id:[0-9]+}", isAuthorized(a.updateTrabajoHandler)).Methods("PUT")
	a.Router.Handle("/trabajos/{id:[0-9]+}", isAuthorized(a.deleteTrabajoHandler)).Methods("DELETE")

}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
