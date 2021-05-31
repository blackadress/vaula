package models

import (
	"context"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestCreatePregunta(t *testing.T) {
	ClearTableExamen(db) // limpia la tabla preguntas tambien

	exa := Pregunta{
		Enunciado: "enun_preg_prueba",
		ExamenId:  1,
		Activo:    true,
	}
	err := exa.CreatePregunta(db)
	if err != nil {
		t.Errorf("No se creo el examen %s", err)
	}

	if exa.ID != 1 {
		t.Errorf("Se esperaba crear una pregunta con ID 1. Se obtuvo %d", exa.ID)
	}
}

func TestGetPregunta(t *testing.T) {
	ClearTableExamen(db) // limpia la tabla preguntas
	AddPreguntas(1, db)

	exa := Pregunta{ID: 1}
	err := exa.GetPregunta(db)

	if err != nil {
		t.Errorf("Se esperaba obtener la pregunta con ID 1. Se obtuvo %v", err)
	}
}

func TestNotGetPregunta(t *testing.T) {
	ClearTablePregunta(db)
	exa := Pregunta{ID: 1}
	err := exa.GetPregunta(db)
	if err != pgx.ErrNoRows {
		t.Errorf("Se esperaba obtener ErrNoRows, se obtuvo diferente error. ERROR %v", err)
	}
}

func TestGetPreguntas(t *testing.T) {
	ClearTablePregunta(db)
	AddPreguntas(2, db)
	preguntas, err := GetPreguntas(db)
	if err != nil {
		t.Errorf("algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(preguntas) != 2 {
		t.Errorf("Se esperaba obtener un array de 2 elementos. Se obtuvo: %v", preguntas)
	}

}

func TestGetZeroPreguntas(t *testing.T) {
	ClearTablePregunta(db)

	preguntas, err := GetPreguntas(db)
	if err != nil {
		t.Errorf("Algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(preguntas) != 0 {
		t.Errorf("Se esperaba obtener un array vacia. Se obtuvo: %v", preguntas)
	}
}

func TestUpdatePregunta(t *testing.T) {
	ClearTableExamen(db) // limpia la tabla preguntas tambien
	AddPreguntas(1, db)

	original_ex := Pregunta{ID: 1}
	err := original_ex.GetPregunta(db)
	if err != nil {
		t.Errorf("El metodo GetPregunta fallo %s", err)
	}

	exa_upd := Pregunta{
		ID:         1,
		Enunciado:  "enun_preg_prueba_upd",
		ExamenId:   2,
		Activo:     false,
	}
	err = exa_upd.UpdatePregunta(db)
	if err != nil {
		t.Errorf("El metodo UpdatePregunta fallo %s", err)
	}

	err = exa_upd.GetPregunta(db)
	if err != nil {
		t.Errorf("El metodo GetPregunta fallo para exa_upd %s", err)
	}

	if original_ex.ID != exa_upd.ID {
		t.Errorf("Se esperaba que el ID no cambiara, cambio de '%d' a '%d'",
			original_ex.ID, exa_upd.ID)
	}

	if original_ex.Enunciado == exa_upd.Enunciado {
		t.Errorf("Se esperaba que los Enunciado cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_ex.Enunciado, exa_upd.Enunciado, original_ex.Enunciado)
	}

	if original_ex.ExamenId == exa_upd.ExamenId {
		t.Errorf("Se esperaba que los ExamenId cambiaran de '%d' a '%d'. Se obtuvo %d",
			original_ex.ExamenId, exa_upd.ExamenId, original_ex.ExamenId)
	}


	if original_ex.Activo == exa_upd.Activo {
		t.Errorf("Se esperaba que los Activo cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_ex.Activo, exa_upd.Activo, original_ex.Activo)
	}

	if original_ex.CreatedAt != exa_upd.CreatedAt {
		t.Errorf("Se esperaba que CreatedAt no cambiara de '%v' a '%v'. Se obtuvo %v",
			original_ex.CreatedAt, exa_upd.CreatedAt, original_ex.CreatedAt)
	}

	if original_ex.UpdatedAt == exa_upd.UpdatedAt {
		t.Errorf("Se esperaba que UpdatedAt cambiara de '%v' a '%v'. Se obtuvo %v",
			original_ex.UpdatedAt, exa_upd.UpdatedAt, original_ex.UpdatedAt)
	}
}

func TestDeletePregunta(t *testing.T) {
	ClearTablePregunta(db)
	AddPreguntas(1, db)

	exa := Pregunta{ID: 1}
	err := exa.DeletePregunta(db)
	if err != nil {
		t.Errorf("Ocurrio un error en el metodo DeletePregunta")
	}
}

const tablePreguntaCreationQuery = `
CREATE TABLE IF NOT EXISTS preguntas
	(
		id SERIAL PRIMARY KEY,
		enunciado TEXT NOT NULL,
		examenId INT REFERENCES cursos(id),

		activo BOOLEAN NOT NULL,
		createdAt TIMESTAMPTZ NOT NULL,
		updatedAt TIMESTAMPTZ NOT NULL
	)
`

func EnsureTablePreguntaExists(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), tablePreguntaCreationQuery)
	if err != nil {
		log.Printf("TEST: error creando tabla pregunta: %s", err)
	}
}

func ClearTablePregunta(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), "DELETE FROM preguntas")
	if err != nil {
		log.Printf("Error deleteando contenidos de la tabla Pregunta %s", err)
	}
	_, err = db.Exec(context.Background(), "ALTER SEQUENCE preguntas_id_seq RESTART WITH 1")
	if err != nil {
		log.Printf("Error reseteando secuencia de pregunta_id %s", err)
	}

}

func AddPreguntas(count int, db *pgxpool.Pool) {
	AddExamenes(count, db)
	if count < 1 {
		count = 1
	}
	now := time.Now()

	for i := 0; i < count; i++ {
		_, err := db.Exec(
			context.Background(),
			`INSERT INTO preguntas(enunciado, examenId, 
				activo, createdAt, updatedAt)
			VALUES($1, $2, $3, $4, $5)`,
			"preg_enun_test_"+strconv.Itoa(i),
			i+1, i%2 == 0, now, now)

		if err != nil {
			log.Printf("Error adding preguntas %s", err)
		}
	}
}
