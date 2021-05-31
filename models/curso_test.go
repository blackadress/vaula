package models

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestCreateCurso(t *testing.T) {
	ClearTableCurso(db)

	al := Curso{
		Nombre:   "nom_curso_prueba",
		Siglas:   "sig_cur_pru",
		Silabo:   "sil_cur_pru",
		Semestre: "se_cur_pru",
		Activo:   true,
	}
	err := al.CreateCurso(db)
	if err != nil {
		t.Errorf("No se creo el curso %s", err)
	}

	if al.ID != 1 {
		t.Errorf("Se esperaba crear un curso con ID 1. Se obtuvo %d", al.ID)
	}
}

func TestGetCurso(t *testing.T) {
	ClearTableCurso(db)
	AddCursos(1, db)
	al := Curso{ID: 1}
	err := al.GetCurso(db)

	if err != nil {
		t.Errorf("Se esperaba obtener el curso con ID 1. Se obtuvo %v", err)
	}
}

func TestNotGetCurso(t *testing.T) {
	ClearTableCurso(db)
	al := Curso{ID: 1}
	err := al.GetCurso(db)
	if err != pgx.ErrNoRows {
		t.Errorf("Se esperaba no obtener ningun curso, se obtuvo diferente error. ERROR %v", err)
	}
}

func TestGetCursos(t *testing.T) {
	ClearTableCurso(db)
	AddCursos(2, db)
	cursos, err := GetCursos(db)
	if err != nil {
		t.Errorf("Algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(cursos) != 2 {
		t.Errorf("Se esperaba obtener un array de 2 elementos. Se obtuvo: %v", cursos)
	}

}

func TestGetZeroCursos(t *testing.T) {
	ClearTableCurso(db)

	cursos, err := GetCursos(db)
	if err != nil {
		t.Errorf("Algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(cursos) != 0 {
		t.Errorf("Se esperaba obtener un array vacia. Se obtuvo: %v", cursos)
	}
}

func TestUpdateCurso(t *testing.T) {
	ClearTableCurso(db)
	AddCursos(1, db)

	original_al := Curso{ID: 1}
	err := original_al.GetCurso(db)
	if err != nil {
		t.Errorf("El metodo GetCurso fallo %s", err)
	}

	al_upd := Curso{
		ID:        1,
		Nombre:   "nom_curso_prueba_upd",
		Siglas: "sig_cur_pru_upd",
		Silabo:    "sil_cur_pru_upd",
		Semestre: "se_cur_pru_upd",
		Activo:    false,
	}
	err = al_upd.UpdateCurso(db)
	if err != nil {
		t.Errorf("El metodo UpdateCurso fallo %s", err)
	}

	err = al_upd.GetCurso(db)
	if err != nil {
		t.Errorf("El metodo GetCurso fallo para al_upd %s", err)
	}

	if original_al.ID != al_upd.ID {
		t.Errorf("Se esperaba que el ID no cambiara, cambio de '%d' a '%d'",
			original_al.ID, al_upd.ID)
	}

	if original_al.Nombre == al_upd.Nombre {
		t.Errorf("Se esperaba que los Nombre cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_al.Nombre, al_upd.Nombre, original_al.Nombre)
	}

	if original_al.Siglas == al_upd.Siglas {
		t.Errorf("Se esperaba que los Siglas cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_al.Siglas, al_upd.Siglas, original_al.Siglas)
	}

	if original_al.Silabo == al_upd.Silabo {
		t.Errorf("Se esperaba que los Silabo cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_al.Silabo, al_upd.Silabo, original_al.Silabo)
	}

	if original_al.Semestre == al_upd.Semestre {
		t.Errorf("Se esperaba que los Semestre cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_al.Semestre, al_upd.Semestre, original_al.Semestre)
	}

	if original_al.Activo == al_upd.Activo {
		t.Errorf("Se esperaba que los Activo cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_al.Activo, al_upd.Activo, original_al.Activo)
	}

	if original_al.CreatedAt != al_upd.CreatedAt {
		t.Errorf("Se esperaba que los CreatedAt cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_al.CreatedAt, al_upd.CreatedAt, original_al.CreatedAt)
	}

	if original_al.UpdatedAt == al_upd.UpdatedAt {
		t.Errorf("Se esperaba que los UpdatedAt cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_al.UpdatedAt, al_upd.UpdatedAt, original_al.UpdatedAt)
	}
}

func TestDeleteCurso(t *testing.T) {
	ClearTableCurso(db)
	AddCursos(1, db)

	al := Curso{ID: 1}
	err := al.DeleteCurso(db)
	if err != nil {
		t.Errorf("Ocurrio un error en el metodo DeleteCurso")
	}
}

const tableCursoCreationQuery = `
CREATE TABLE IF NOT EXISTS cursos
	(
		id SERIAL PRIMARY KEY,
		nombre VARCHAR(200) NOT NULL,
		siglas VARCHAR(20) NOT NULL,
		silabo VARCHAR(200) NOT NULL,
		semestre VARCHAR(20) NOT NULL,

		activo BOOLEAN NOT NULL,
		createdAt TIMESTAMPTZ NOT NULL,
		updatedAt TIMESTAMPTZ NOT NULL
	)
`

func EnsureTableCursoExists(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), tableCursoCreationQuery)
	if err != nil {
		log.Printf("TEST: error creando tabla cursos: %s", err)
	}
}

func ClearTableCurso(db *pgxpool.Pool) {
	ClearTablePregunta(db)
	ClearTableTrabajo(db)
	ClearTableExamen(db)
	_, err := db.Exec(context.Background(), "DELETE FROM cursos")
	if err != nil {
		log.Printf("Error deleteando tabla %s", err)
	}
	_, err = db.Exec(context.Background(), "ALTER SEQUENCE cursos_id_seq RESTART WITH 1")
	if err != nil {
		log.Printf("Error alterando secuencia de curso_id %s", err)
	}

}

func AddCursos(count int, db *pgxpool.Pool) {
	if count < 1 {
		count = 1
	}
	now := time.Now()

	for i := 0; i < count; i++ {
		semestre := fmt.Sprintf("%.20s", "semestre_"+strconv.Itoa(i))
		db.Exec(
			context.Background(),
			`INSERT INTO cursos(nombre, siglas, silabo, semestre, activo, createdAt, updatedAt)
			VALUES($1, $2, $3, $4, $5, $6, $7)`,
			"curso_test_"+strconv.Itoa(i),
			"TS-0"+strconv.Itoa(i),
			"silabo_test_"+strconv.Itoa(i),
			semestre, i%2 == 0, now, now)
	}
}
