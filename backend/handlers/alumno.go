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

func (a *App) getAlumnoByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("GET %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de alumno invalido")
		return
	}

	alum := models.Alumno{ID: id}
	err = alum.GetAlumno(a.DB)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			log.Printf("GET %s code: %d ERROR: %s -- no rows", r.RequestURI,
				http.StatusNotFound, err.Error())
			respondWithError(w, http.StatusNotFound, "Alumno no encontrado")
		default:
			log.Printf("GET %s code: %d ERROR: %s -- wtf just happened? default switch",
				r.RequestURI, http.StatusInternalServerError, err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	log.Printf("GET %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, alum)
	return
}

func (a *App) getAlumnosHandler(w http.ResponseWriter, r *http.Request) {
	alumnos, err := models.GetAlumnos(a.DB)
	if err != nil {
		log.Printf("GET %s code: %d ERROR: %s -- models.GetAlumnos", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("GET %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, alumnos)
	return

}

func (a *App) createAlumnoHandler(w http.ResponseWriter, r *http.Request) {
	var alum models.Alumno

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&alum)
	if err != nil {
		log.Printf("POST %s code: %d ERROR: %s -- decoder", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()

	// hay la request debe especificamente settear el valor de alum.Activo,
	// debido a que por defecto se inicializa en 'false'
	err = alum.CreateAlumno(a.DB)
	if err != nil {
		log.Printf("POST %s code: %d ERROR: %s -- alumno.createAlumno", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("POST %s code: %d", r.RequestURI, http.StatusCreated)
	respondWithJSON(w, http.StatusCreated, alum)
	return
}

func (a *App) updateAlumnoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de Alumno invalido")
		return
	}

	var alum models.Alumno
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&alum)
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s -- decoder", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()

	alum.ID = id
	err = alum.UpdateAlumno(a.DB)
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s alumno.UpdateAlumno", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("PUT %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, alum)
	return
}

func (a *App) deleteAlumnoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("DELETE %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de alumno invalido")
		return
	}

	alum := models.Alumno{ID: id}
	if err := alum.DeleteAlumno(a.DB); err != nil {
		log.Printf("DELETE %s code: %d ERROR: %s -- alumno.DeleteAlumno", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("DELETE %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, map[string]int{"exito": 1, "id": alum.ID})
	return
}
