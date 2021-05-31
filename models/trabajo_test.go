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

func TestCreateTrabajo(t *testing.T) {
	ClearTableCurso(db) // limpia la tabla trabajos tambien
	AddCursos(1, db)

	loc, _ := time.LoadLocation("America/Lima")
	fechaInicio := time.Date(2021, time.June,
		20, 18, 0, 0, 0, loc)
	fechaFinal := time.Date(2021, time.June,
		22, 18, 0, 0, 0, loc)

	tra := Trabajo{
		Descripcion: "nom_trabajo_prueba",
		FechaInicio: fechaInicio,
		FechaFinal:  fechaFinal,
		CursoId:     1,
		Activo:      true,
	}
	err := tra.CreateTrabajo(db)
	if err != nil {
		t.Errorf("No se creo el trabajo %s", err)
	}

	if tra.ID != 1 {
		t.Errorf("Se esperaba crear un trabajo con ID 1. Se obtuvo %d", tra.ID)
	}
}

func TestGetTrabajo(t *testing.T) {
	ClearTableCurso(db) // limpia la tabla trabajos tambien
	AddTrabajos(1, db)

	tra := Trabajo{ID: 1}
	err := tra.GetTrabajo(db)

	if err != nil {
		t.Errorf("Se esperaba obtener el trabajo con ID 1. Se obtuvo %v", err)
	}
}

func TestNotGetTrabajo(t *testing.T) {
	ClearTableTrabajo(db)
	tra := Trabajo{ID: 1}
	err := tra.GetTrabajo(db)
	if err != pgx.ErrNoRows {
		t.Errorf("Se esperaba no ErrNoRows, se obtuvo diferente error. ERROR %v", err)
	}
}

func TestGetTrabajos(t *testing.T) {
	ClearTableTrabajo(db)
	AddTrabajos(2, db)
	trabajos, err := GetTrabajos(db)
	if err != nil {
		t.Errorf("algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(trabajos) != 2 {
		t.Errorf("Se esperaba obtener un array de 2 elementos. Se obtuvo: %v", trabajos)
	}

}

func TestGetZeroTrabajos(t *testing.T) {
	ClearTableTrabajo(db)

	trabajos, err := GetTrabajos(db)
	if err != nil {
		t.Errorf("Algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(trabajos) != 0 {
		t.Errorf("Se esperaba obtener un array vacia. Se obtuvo: %v", trabajos)
	}
}

func TestUpdateTrabajo(t *testing.T) {
	ClearTableCurso(db)
	AddTrabajos(1, db)
	AddCursos(1, db)

	original_tra := Trabajo{ID: 1}
	err := original_tra.GetTrabajo(db)
	if err != nil {
		t.Errorf("El metodo GetTrabajo fallo %s", err)
	}

	loc, _ := time.LoadLocation("America/Lima")
	fechaInicio := time.Date(2023, time.June,
		20, 18, 0, 0, 0, loc)
	fechaFinal := time.Date(2023, time.June,
		22, 18, 0, 0, 0, loc)

	tra_upd := Trabajo{
		ID:          1,
		Descripcion: "nom_trabajo_prueba_upd",
		FechaInicio: fechaInicio,
		FechaFinal:  fechaFinal,
		CursoId:     2,
		Activo:      false,
	}
	err = tra_upd.UpdateTrabajo(db)
	if err != nil {
		t.Errorf("El metodo UpdateTrabajo fallo %s", err)
	}

	err = tra_upd.GetTrabajo(db)
	if err != nil {
		t.Errorf("El metodo GetTrabajo fallo para tra_upd %s", err)
	}

	if original_tra.ID != tra_upd.ID {
		t.Errorf("Se esperaba que el ID no cambiara, cambio de '%d' a '%d'",
			original_tra.ID, tra_upd.ID)
	}

	if original_tra.Descripcion == tra_upd.Descripcion {
		t.Errorf("Se esperaba que los Descripcion cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_tra.Descripcion, tra_upd.Descripcion, original_tra.Descripcion)
	}

	if original_tra.FechaInicio == tra_upd.FechaInicio {
		t.Errorf("Se esperaba que los FechaInicio cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_tra.FechaInicio, tra_upd.FechaInicio, original_tra.FechaInicio)
	}

	if original_tra.FechaFinal == tra_upd.FechaFinal {
		t.Errorf("Se esperaba que los FechaFinal cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_tra.FechaFinal, tra_upd.FechaFinal, original_tra.FechaFinal)
	}

	if original_tra.CursoId == tra_upd.CursoId {
		t.Errorf("Se esperaba que los CursoId cambiaran de '%d' a '%d'. Se obtuvo %d",
			original_tra.CursoId, tra_upd.CursoId, original_tra.CursoId)
	}

	if original_tra.Activo == tra_upd.Activo {
		t.Errorf("Se esperaba que los Activo cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_tra.Activo, tra_upd.Activo, original_tra.Activo)
	}

	if original_tra.CreatedAt != tra_upd.CreatedAt {
		t.Errorf("Se esperaba que CreatedAt no cambiara de '%v' a '%v'. Se obtuvo %v",
			original_tra.CreatedAt, tra_upd.CreatedAt, original_tra.CreatedAt)
	}

	if original_tra.UpdatedAt == tra_upd.UpdatedAt {
		t.Errorf("Se esperaba que UpdatedAt cambiara de '%v' a '%v'. Se obtuvo %v",
			original_tra.UpdatedAt, tra_upd.UpdatedAt, original_tra.UpdatedAt)
	}
}

func TestDeleteTrabajo(t *testing.T) {
	ClearTableTrabajo(db)
	AddTrabajos(1, db)

	tra := Trabajo{ID: 1}
	err := tra.DeleteTrabajo(db)
	if err != nil {
		t.Errorf("Ocurrio un error en el metodo DeleteTrabajo")
	}
}

const tableTrabajoCreationQuery = `
CREATE TABLE IF NOT EXISTS trabajos
	(
		id SERIAL PRIMARY KEY,
		descripcion TEXT NOT NULL,
		fechaInicio TIMESTAMPTZ NOT NULL,
		fechaFinal TIMESTAMPTZ NOT NULL,
		cursoId INT REFERENCES cursos(id),

		activo BOOLEAN NOT NULL,
		createdAt TIMESTAMPTZ NOT NULL,
		updatedAt TIMESTAMPTZ NOT NULL
	)
`

func EnsureTableTrabajoExists(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), tableTrabajoCreationQuery)
	if err != nil {
		log.Printf("TEST: error creando tabla trabajos: %s", err)
	}
}

func ClearTableTrabajo(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), "DELETE FROM trabajos")
	if err != nil {
		log.Printf("Error deleteando contenidos de la tabla Trabajo %s", err)
	}
	_, err = db.Exec(context.Background(), "ALTER SEQUENCE trabajos_id_seq RESTART WITH 1")
	if err != nil {
		log.Printf("Error reseteando secuencia de trabajo_id %s", err)
	}

}

func AddTrabajos(count int, db *pgxpool.Pool) {
	AddCursos(count, db)
	if count < 1 {
		count = 1
	}
	now := time.Now()
	loc, _ := time.LoadLocation("America/Lima")
	fechaInicio := time.Date(2022, time.June,
		20, 18, 0, 0, 0, loc)
	fechaFinal := time.Date(2022, time.June,
		22, 18, 0, 0, 0, loc)

	for i := 0; i < count; i++ {
		_, err := db.Exec(
			context.Background(),
			`INSERT INTO trabajos(descripcion, fechaInicio, fechaFinal,
				cursoId, activo, createdAt, updatedAt)
			VALUES($1, $2, $3, $4, $5, $6, $7)`,
			"trabajo_test_"+strconv.Itoa(i),
			fechaInicio, fechaFinal,
			i+1, i%2 == 0, now, now)

		if err != nil {
			log.Printf("Error adding trabajos %s", err)
		}
	}
}
