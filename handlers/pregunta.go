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

func (a *App) getPreguntaByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("GET %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de pregunta invalido")
		return
	}

	pregunta := models.Pregunta{ID: id}
	err = pregunta.GetPregunta(a.DB)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			log.Printf("GET %s code: %d ERROR: %s -- no rows", r.RequestURI,
				http.StatusNotFound, err.Error())
			respondWithError(w, http.StatusNotFound, "Pregunta no encontrado")
		default:
			log.Printf("GET %s code: %d ERROR: %s -- wtf just happened? default switch",
				r.RequestURI, http.StatusInternalServerError, err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	log.Printf("GET %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, pregunta)
	return
}

func (a *App) getPreguntasHandler(w http.ResponseWriter, r *http.Request) {
	preguntas, err := models.GetPreguntas(a.DB)
	if err != nil {
		log.Printf("GET %s code: %d ERROR: %s -- models.GetPreguntas", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("GET %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, preguntas)
	return

}

func (a *App) createPreguntaHandler(w http.ResponseWriter, r *http.Request) {
	var pregunta models.Pregunta

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&pregunta)
	if err != nil {
		log.Printf("POST %s code: %d ERROR: %s -- decoder", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()

	// hay la request debe especificamente settear el valor de pregunta.Activo,
	// debido a que por defecto se inicializa en 'false'
	err = pregunta.CreatePregunta(a.DB)
	if err != nil {
		log.Printf("POST %s code: %d ERROR: %s -- pregunta.createPregunta", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("POST %s code: %d", r.RequestURI, http.StatusCreated)
	respondWithJSON(w, http.StatusCreated, pregunta)
	return
}

func (a *App) updatePreguntaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de Pregunta invalido")
		return
	}

	var pregunta models.Pregunta
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&pregunta)
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s -- decoder", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()

	pregunta.ID = id
	err = pregunta.UpdatePregunta(a.DB)
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s pregunta.UpdatePregunta", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("PUT %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, pregunta)
	return
}

func (a *App) deletePreguntaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("DELETE %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de pregunta invalido")
		return
	}

	pregunta := models.Pregunta{ID: id}
	if err := pregunta.DeletePregunta(a.DB); err != nil {
		log.Printf("DELETE %s code: %d ERROR: %s -- pregunta.DeletePregunta", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("DELETE %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, map[string]int{"exito": 1, "id": pregunta.ID})
	return
}
