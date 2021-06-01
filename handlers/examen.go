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

func (a *App) getExamenByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("GET %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de examen invalido")
		return
	}

	examen := models.Examen{ID: id}
	err = examen.GetExamen(a.DB)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			log.Printf("GET %s code: %d ERROR: %s -- no rows", r.RequestURI,
				http.StatusNotFound, err.Error())
			respondWithError(w, http.StatusNotFound, "Examen no encontrado")
		default:
			log.Printf("GET %s code: %d ERROR: %s -- wtf just happened? default switch",
				r.RequestURI, http.StatusInternalServerError, err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	log.Printf("GET %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, examen)
	return
}

func (a *App) getExamenesHandler(w http.ResponseWriter, r *http.Request) {
	examenes, err := models.GetExamenes(a.DB)
	if err != nil {
		log.Printf("GET %s code: %d ERROR: %s -- models.GetExamenes", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("GET %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, examenes)
	return

}

func (a *App) createExamenHandler(w http.ResponseWriter, r *http.Request) {
	var examen models.Examen

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&examen)
	if err != nil {
		log.Printf("POST %s code: %d ERROR: %s -- decoder", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()

	// hay la request debe especificamente settear el valor de examen.Activo,
	// debido a que por defecto se inicializa en 'false'
	err = examen.CreateExamen(a.DB)
	if err != nil {
		log.Printf("POST %s code: %d ERROR: %s -- examen.createExamen", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("POST %s code: %d", r.RequestURI, http.StatusCreated)
	respondWithJSON(w, http.StatusCreated, examen)
	return
}

func (a *App) updateExamenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de Examen invalido")
		return
	}

	var examen models.Examen
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&examen)
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s -- decoder", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()

	examen.ID = id
	err = examen.UpdateExamen(a.DB)
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s examen.UpdateExamen", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("PUT %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, examen)
	return
}

func (a *App) deleteExamenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("DELETE %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de examen invalido")
		return
	}

	examen := models.Examen{ID: id}
	if err := examen.DeleteExamen(a.DB); err != nil {
		log.Printf("DELETE %s code: %d ERROR: %s -- examen.DeleteExamen", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("DELETE %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, map[string]int{"exito": 1, "id": examen.ID})
	return
}
