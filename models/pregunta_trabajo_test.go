package models

// import (
// 	"context"
// 	"log"
// 	"strconv"
// 	"testing"
// 	"time"

// 	"github.com/jackc/pgx/v4"
// 	"github.com/jackc/pgx/v4/pgxpool"
// )

// func TestCreatePreguntaTrabajo(t *testing.T) {
// 	ClearTableTrabajo(db) // limpia la tabla preguntas tambien

// 	exa := PreguntaTrabajo{
// 		Enunciado: "enun_preg_prueba",
// 		TrabajoId:  1,
// 		Activo:    true,
// 	}
// 	err := exa.CreatePreguntaTrabajo(db)
// 	if err != nil {
// 		t.Errorf("No se creo el examen %s", err)
// 	}

// 	if exa.ID != 1 {
// 		t.Errorf("Se esperaba crear una pregunta con ID 1. Se obtuvo %d", exa.ID)
// 	}
// }

// func TestGetPreguntaTrabajo(t *testing.T) {
// 	ClearTableTrabajo(db) // limpia la tabla preguntas
// 	AddPreguntaTrabajos(1, db)

// 	exa := PreguntaTrabajo{ID: 1}
// 	err := exa.GetPreguntaTrabajo(db)

// 	if err != nil {
// 		t.Errorf("Se esperaba obtener la pregunta con ID 1. Se obtuvo %v", err)
// 	}
// }

// func TestNotGetPreguntaTrabajo(t *testing.T) {
// 	ClearTablePreguntaTrabajo(db)
// 	exa := PreguntaTrabajo{ID: 1}
// 	err := exa.GetPreguntaTrabajo(db)
// 	if err != pgx.ErrNoRows {
// 		t.Errorf("Se esperaba obtener ErrNoRows, se obtuvo diferente error. ERROR %v", err)
// 	}
// }

// func TestGetPreguntaTrabajos(t *testing.T) {
// 	ClearTablePreguntaTrabajo(db)
// 	AddPreguntaTrabajos(2, db)
// 	preguntas, err := GetPreguntaTrabajos(db)
// 	if err != nil {
// 		t.Errorf("algo salio mal con la comunicacion con la DB %s", err)
// 	}

// 	if len(preguntas) != 2 {
// 		t.Errorf("Se esperaba obtener un array de 2 elementos. Se obtuvo: %v", preguntas)
// 	}

// }

// func TestGetZeroPreguntaTrabajos(t *testing.T) {
// 	ClearTablePreguntaTrabajo(db)

// 	preguntas, err := GetPreguntaTrabajos(db)
// 	if err != nil {
// 		t.Errorf("Algo salio mal con la comunicacion con la DB %s", err)
// 	}

// 	if len(preguntas) != 0 {
// 		t.Errorf("Se esperaba obtener un array vacia. Se obtuvo: %v", preguntas)
// 	}
// }

// func TestUpdatePreguntaTrabajo(t *testing.T) {
// 	ClearTableTrabajo(db) // limpia la tabla preguntas tambien
// 	AddPreguntaTrabajos(1, db)

// 	original_ex := PreguntaTrabajo{ID: 1}
// 	err := original_ex.GetPreguntaTrabajo(db)
// 	if err != nil {
// 		t.Errorf("El metodo GetPreguntaTrabajo fallo %s", err)
// 	}

// 	exa_upd := PreguntaTrabajo{
// 		ID:         1,
// 		Enunciado:  "enun_preg_prueba_upd",
// 		TrabajoId:   2,
// 		Activo:     false,
// 	}
// 	err = exa_upd.UpdatePreguntaTrabajo(db)
// 	if err != nil {
// 		t.Errorf("El metodo UpdatePreguntaTrabajo fallo %s", err)
// 	}

// 	err = exa_upd.GetPreguntaTrabajo(db)
// 	if err != nil {
// 		t.Errorf("El metodo GetPreguntaTrabajo fallo para exa_upd %s", err)
// 	}

// 	if original_ex.ID != exa_upd.ID {
// 		t.Errorf("Se esperaba que el ID no cambiara, cambio de '%d' a '%d'",
// 			original_ex.ID, exa_upd.ID)
// 	}

// 	if original_ex.Enunciado == exa_upd.Enunciado {
// 		t.Errorf("Se esperaba que los Enunciado cambiaran de '%s' a '%s'. Se obtuvo %s",
// 			original_ex.Enunciado, exa_upd.Enunciado, original_ex.Enunciado)
// 	}

// 	if original_ex.TrabajoId == exa_upd.TrabajoId {
// 		t.Errorf("Se esperaba que los TrabajoId cambiaran de '%d' a '%d'. Se obtuvo %d",
// 			original_ex.TrabajoId, exa_upd.TrabajoId, original_ex.TrabajoId)
// 	}


// 	if original_ex.Activo == exa_upd.Activo {
// 		t.Errorf("Se esperaba que los Activo cambiaran de '%v' a '%v'. Se obtuvo %v",
// 			original_ex.Activo, exa_upd.Activo, original_ex.Activo)
// 	}

// 	if original_ex.CreatedAt != exa_upd.CreatedAt {
// 		t.Errorf("Se esperaba que CreatedAt no cambiara de '%v' a '%v'. Se obtuvo %v",
// 			original_ex.CreatedAt, exa_upd.CreatedAt, original_ex.CreatedAt)
// 	}

// 	if original_ex.UpdatedAt == exa_upd.UpdatedAt {
// 		t.Errorf("Se esperaba que UpdatedAt cambiara de '%v' a '%v'. Se obtuvo %v",
// 			original_ex.UpdatedAt, exa_upd.UpdatedAt, original_ex.UpdatedAt)
// 	}
// }

// func TestDeletePreguntaTrabajo(t *testing.T) {
// 	ClearTablePreguntaTrabajo(db)
// 	AddPreguntaTrabajos(1, db)

// 	exa := PreguntaTrabajo{ID: 1}
// 	err := exa.DeletePreguntaTrabajo(db)
// 	if err != nil {
// 		t.Errorf("Ocurrio un error en el metodo DeletePreguntaTrabajo")
// 	}
// }

// const tablePreguntaTrabajoCreationQuery = `
// CREATE TABLE IF NOT EXISTS preguntas
// 	(
// 		id SERIAL PRIMARY KEY,
// 		enunciado TEXT NOT NULL,
// 		examenId INT REFERENCES cursos(id),

// 		activo BOOLEAN NOT NULL,
// 		createdAt TIMESTAMPTZ NOT NULL,
// 		updatedAt TIMESTAMPTZ NOT NULL
// 	)
// `

// func EnsureTablePreguntaTrabajoExists(db *pgxpool.Pool) {
// 	_, err := db.Exec(context.Background(), tablePreguntaTrabajoCreationQuery)
// 	if err != nil {
// 		log.Printf("TEST: error creando tabla pregunta: %s", err)
// 	}
// }

// func ClearTablePreguntaTrabajo(db *pgxpool.Pool) {
// 	_, err := db.Exec(context.Background(), "DELETE FROM preguntas")
// 	if err != nil {
// 		log.Printf("Error deleteando contenidos de la tabla PreguntaTrabajo %s", err)
// 	}
// 	_, err = db.Exec(context.Background(), "ALTER SEQUENCE preguntas_id_seq RESTART WITH 1")
// 	if err != nil {
// 		log.Printf("Error reseteando secuencia de pregunta_id %s", err)
// 	}

// }

// func AddPreguntaTrabajos(count int, db *pgxpool.Pool) {
// 	AddTrabajoes(count, db)
// 	if count < 1 {
// 		count = 1
// 	}
// 	now := time.Now()

// 	for i := 0; i < count; i++ {
// 		_, err := db.Exec(
// 			context.Background(),
// 			`INSERT INTO preguntas(enunciado, examenId, 
// 				activo, createdAt, updatedAt)
// 			VALUES($1, $2, $3, $4, $5)`,
// 			"preg_enun_test_"+strconv.Itoa(i),
// 			i+1, i%2 == 0, now, now)

// 		if err != nil {
// 			log.Printf("Error adding preguntas %s", err)
// 		}
// 	}
// }
