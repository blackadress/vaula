package handlers

import (
	"fmt"
	"net/http"

	"github.com/blackadress/vaula/main"
	"github.com/blackadress/vaula/models"
)

func HelloFromUserHandlers() {
	fmt.Printf("Hello from handlers\n")
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	products, err := models.GetUsers(main.WebApp.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}
