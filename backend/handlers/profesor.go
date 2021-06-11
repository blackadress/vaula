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

func (a *App) getProfesorByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("GET %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de profesor invalido")
		return
	}

	profesor := models.Profesor{ID: id}
	err = profesor.GetProfesor(a.DB)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			log.Printf("GET %s code: %d ERROR: %s -- no rows", r.RequestURI,
				http.StatusNotFound, err.Error())
			respondWithError(w, http.StatusNotFound, "Profesor no encontrado")
		default:
			log.Printf("GET %s code: %d ERROR: %s -- wtf just happened? default switch",
				r.RequestURI, http.StatusInternalServerError, err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	log.Printf("GET %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, profesor)
	return
}

func (a *App) getProfesoresHandler(w http.ResponseWriter, r *http.Request) {
	profesores, err := models.GetProfesores(a.DB)
	if err != nil {
		log.Printf("GET %s code: %d ERROR: %s -- models.GetProfesores", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("GET %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, profesores)
	return

}

func (a *App) createProfesorHandler(w http.ResponseWriter, r *http.Request) {
	var profesor models.Profesor

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&profesor)
	if err != nil {
		log.Printf("POST %s code: %d ERROR: %s -- decoder", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()

	// hay la request debe especificamente settear el valor de profesor.Activo,
	// debido a que por defecto se inicializa en 'false'
	err = profesor.CreateProfesor(a.DB)
	if err != nil {
		log.Printf("POST %s code: %d ERROR: %s -- profesor.createProfesor", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("POST %s code: %d", r.RequestURI, http.StatusCreated)
	respondWithJSON(w, http.StatusCreated, profesor)
	return
}

func (a *App) updateProfesorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de Profesor invalido")
		return
	}

	var profesor models.Profesor
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&profesor)
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s -- decoder", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()

	profesor.ID = id
	err = profesor.UpdateProfesor(a.DB)
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s profesor.UpdateProfesor", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("PUT %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, profesor)
	return
}

func (a *App) deleteProfesorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("DELETE %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de profesor invalido")
		return
	}

	profesor := models.Profesor{ID: id}
	if err := profesor.DeleteProfesor(a.DB); err != nil {
		log.Printf("DELETE %s code: %d ERROR: %s -- profesor.DeleteProfesor", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("DELETE %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, map[string]int{"exito": 1, "id": profesor.ID})
	return
}
