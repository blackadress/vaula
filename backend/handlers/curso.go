package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/blackadress/vaula/models"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

func (a *App) getCursoByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("GET %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de curso invalido")
		return
	}

	curso := models.Curso{ID: id}
	err = curso.GetCurso(a.DB)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			log.Printf("GET %s code: %d ERROR: %s -- no rows", r.RequestURI,
				http.StatusNotFound, err.Error())
			respondWithError(w, http.StatusNotFound, "Curso no encontrado")
		default:
			log.Printf("GET %s code: %d ERROR: %s -- wtf just happened? default switch",
				r.RequestURI, http.StatusInternalServerError, err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	log.Printf("GET %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, curso)
	return
}

func (a *App) getCursosHandler(w http.ResponseWriter, r *http.Request) {
	cursos, err := models.GetCursos(a.DB)
	if err != nil {
		log.Printf("GET %s code: %d ERROR: %s -- models.GetCursos", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("GET %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, cursos)
	return

}

func (a *App) createCursoHandler(w http.ResponseWriter, r *http.Request) {
	var curso models.Curso

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&curso)
	if err != nil {
		log.Printf("POST %s code: %d ERROR: %s -- decoder", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()

	// hay la request debe especificamente settear el valor de curso.Activo,
	// debido a que por defecto se inicializa en 'false'
	err = curso.CreateCurso(a.DB)
	if err != nil {
		log.Printf("POST %s code: %d ERROR: %s -- curso.createCurso", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("POST %s code: %d", r.RequestURI, http.StatusCreated)
	respondWithJSON(w, http.StatusCreated, curso)
	return
}

func (a *App) updateCursoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de Curso invalido")
		return
	}

	var curso models.Curso
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&curso)
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s -- decoder", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()

	curso.ID = id
	err = curso.UpdateCurso(a.DB)
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s curso.UpdateCurso", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("PUT %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, curso)
	return
}

func (a *App) deleteCursoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("DELETE %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de curso invalido")
		return
	}

	curso := models.Curso{ID: id}
	if err := curso.DeleteCurso(a.DB); err != nil {
		log.Printf("DELETE %s code: %d ERROR: %s -- curso.DeleteCurso", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("DELETE %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, map[string]int{"exito": 1, "id": curso.ID})
	return
}
