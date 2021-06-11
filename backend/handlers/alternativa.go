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

func (a *App) getAlternativaByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("GET %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de alternativa invalido")
		return
	}

	alt := models.Alternativa{ID: id}
	err = alt.GetAlternativa(a.DB)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			log.Printf("GET %s code: %d ERROR: %s -- no rows", r.RequestURI,
				http.StatusNotFound, err.Error())
			respondWithError(w, http.StatusNotFound, "Alternativa no encontrada")
		default:
			log.Printf("GET %s code: %d ERROR: %s -- wtf just happened? default switch", r.RequestURI,
				http.StatusInternalServerError, err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	log.Printf("GET %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, alt)
	return
}

func (a *App) getAlternativasHandler(w http.ResponseWriter, r *http.Request) {
	alternativas, err := models.GetAlternativas(a.DB)
	if err != nil {
		log.Printf("GET %s code: %d ERROR: %s -- no se obtuvieron alternativas", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("GET %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, alternativas)
	return

}

func (a *App) createAlternativaHandler(w http.ResponseWriter, r *http.Request) {
	var alt models.Alternativa
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&alt)
	if err != nil {
		log.Printf("POST %s code: %d ERROR: %s -- decoder", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()

	// hay la request debe especificamente settear el valor de alt.Activo,
	// debido a que por defecto se inicializa en 'false'
	err = alt.CreateAlternativa(a.DB)
	if err != nil {
		log.Printf("POST %s code: %d ERROR: %s -- models.createAlternativa", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("POST %s code: %d", r.RequestURI, http.StatusCreated)
	respondWithJSON(w, http.StatusCreated, alt)
	return
}

func (a *App) updateAlternativaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de Alternativa invalido")
		return
	}

	var alt models.Alternativa
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&alt)
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s -- decoder", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()

	alt.ID = id
	err = alt.UpdateAlternativa(a.DB)
	if err != nil {
		log.Printf("PUT %s code: %d ERROR: %s -- alternativa.UpdateAlternativa", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("PUT %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, alt)
	return
}

func (a *App) deleteAlternativaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("DELETE %s code: %d ERROR: %s -- strconv", r.RequestURI,
			http.StatusBadRequest, err.Error())
		respondWithError(w, http.StatusBadRequest, "ID de alternativa invalido")
		return
	}

	alt := models.Alternativa{ID: id}
	if err := alt.DeleteAlternativa(a.DB); err != nil {
		log.Printf("DELETE %s code: %d ERROR: %s -- alternativa.DeleteAlternativa", r.RequestURI,
			http.StatusInternalServerError, err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("DELETE %s code: %d", r.RequestURI, http.StatusOK)
	respondWithJSON(w, http.StatusOK, map[string]int{"exito": 1, "id": alt.ID})
	return
}
