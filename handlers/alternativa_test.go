package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testing"

	"github.com/blackadress/vaula/globals"
)

func TestEmptyAlternativaTable(t *testing.T) {
	clearTableAlternativa()
	ensureAuthorizedUserExists()

	token := getTestJWT()
	token_str := fmt.Sprintf("Bearer %s", token.AccessToken)

	req, _ := http.NewRequest("GET", "/alternativas", nil)
	req.Header.Set("Authorization", token_str)
	response := executeRequest(req, a)

	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	if body != "" {
		t.Errorf("Se esperaba un array vacio. Se obtuvo %#v", body)
	}
}

const tableAlternativaCreationQuery = `
CREATE TABLE IF NOT EXISTS alternativas
	(
		id SERIAL,
		valor TEXT NOT NULL,
		correcto BOOLEAN NOT NULL,

		activo BOOLEAN NOT NULL,
		createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updatedAt TIMESTAMPTZ
	)
`

// es posible hacer decouple de `globals.DB`?
func ensureTableAlternativaExists() {
	_, err := globals.DB.Exec(context.Background(), tableAlternativaCreationQuery)
	if err != nil {
		log.Printf("TEST: error creando tabla alternativas: %s", err)
	}
}

func clearTableAlternativa() {
	globals.DB.Exec(context.Background(), "DELETE FROM alternativas")
	globals.DB.Exec(context.Background(), "ALTER SEQUENCE alternativas_id_seq RESTART WITH 1")
}

func addAlternativas(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		globals.DB.Exec(
			context.Background(),
			`INSERT INTO alternativas(valor, correcto, activo)
			VALUES($1, $2, $3)`,
			"valor_"+strconv.Itoa(i),
			i%2 == 1,
			true)
	}
}
