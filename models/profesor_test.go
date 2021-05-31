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

func TestCreateProfesor(t *testing.T) {
	ClearTableUsuario(db)
	ClearTableProfesor(db)
	AddUsers(1, db)

	al := Profesor{
		Nombres:   "nom_pro prueba",
		Apellidos: "ap_prof prueba",
		UsuarioId: 1,
		Activo:    true,
	}
	err := al.CreateProfesor(db)
	if err != nil {
		t.Errorf("No se creo el profesor")
	}

	if al.ID != 1 {
		t.Errorf("Se esperaba crear un profesor con ID 1. Se obtuvo %d", al.ID)
	}
}

func TestGetProfesor(t *testing.T) {
	ClearTableProfesor(db)
	AddProfesores(1, db)
	al := Profesor{ID: 1}
	err := al.GetProfesor(db)

	if err != nil {
		t.Errorf("Se esperaba obtener el profesor con ID 1. Se obtuvo %v", err)
	}
}

func TestNotGetProfesor(t *testing.T) {
	ClearTableProfesor(db)
	al := Profesor{ID: 1}
	err := al.GetProfesor(db)
	if err != pgx.ErrNoRows {
		t.Errorf("Se esperaba ErrNoRows, se obtuvo diferente error. ERROR %v", err)
	}
}

func TestGetProfesors(t *testing.T) {
	ClearTableProfesor(db)
	AddProfesores(2, db)
	profesores, err := GetProfesores(db)
	if err != nil {
		t.Errorf("Algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(profesores) != 2 {
		t.Errorf("Se esperaba obtener un array de 2 elementos. Se obtuvo: %v", profesores)
	}
}

func TestGetZeroProfesores(t *testing.T) {
	ClearTableProfesor(db)

	profesores, err := GetProfesores(db)
	if err != nil {
		t.Errorf("Algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(profesores) != 0 {
		t.Errorf("Se esperaba obtener un array vacia. Se obtuvo: %v", profesores)
	}
}

func TestUpdateProfesor(t *testing.T) {
	ClearTableProfesor(db)
	ClearTableUsuario(db)
	AddProfesores(1, db)
	AddUsers(1, db)

	original_prof := Profesor{ID: 1}
	err := original_prof.GetProfesor(db)
	if err != nil {
		t.Errorf("El metodo GetProfesor fallo %s", err)
	}

	al_upd := Profesor{
		ID:        1,
		Nombres:   "nom_pro upd",
		Apellidos: "ap_prof upd",
		UsuarioId: 2,
		Activo:    false,
	}
	err = al_upd.UpdateProfesor(db)
	if err != nil {
		t.Errorf("El metodo UpdateProfesor fallo %s", err)
	}

	err = al_upd.GetProfesor(db)
	if err != nil {
		t.Errorf("El metodo GetProfesor fallo para al_upd %s", err)
	}

	if original_prof.ID != al_upd.ID {
		t.Errorf("Se esperaba que el ID no cambiara, cambio de '%d' a '%d'",
			original_prof.ID, al_upd.ID)
	}

	if original_prof.Nombres == al_upd.Nombres {
		t.Errorf("Se esperaba que los Nombres cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_prof.Nombres, al_upd.Nombres, original_prof.Nombres)
	}

	if original_prof.Apellidos == al_upd.Apellidos {
		t.Errorf("Se esperaba que los Apellidos cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_prof.Apellidos, al_upd.Apellidos, original_prof.Apellidos)
	}

	if original_prof.UsuarioId == al_upd.UsuarioId {
		t.Errorf("Se esperaba que los UsuarioId cambiaran de '%d' a '%d'. Se obtuvo %d",
			original_prof.UsuarioId, al_upd.UsuarioId, original_prof.UsuarioId)
	}

	if original_prof.Activo == al_upd.Activo {
		t.Errorf("Se esperaba que los Activo cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_prof.Activo, al_upd.Activo, original_prof.Activo)
	}

	if original_prof.CreatedAt != al_upd.CreatedAt {
		t.Errorf("Se esperaba que los CreatedAt cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_prof.CreatedAt, al_upd.CreatedAt, original_prof.CreatedAt)
	}

	if original_prof.UpdatedAt == al_upd.UpdatedAt {
		t.Errorf("Se esperaba que los UpdatedAt cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_prof.UpdatedAt, al_upd.UpdatedAt, original_prof.UpdatedAt)
	}
}

func TestDeleteProfesor(t *testing.T) {
	ClearTableProfesor(db)
	AddProfesores(1, db)

	al := Profesor{ID: 1}
	err := al.DeleteProfesor(db)
	if err != nil {
		t.Errorf("Ocurrio un error en el metodo DeleteProfesor")
	}
}

const tableProfesorCreationQuery = `
CREATE TABLE IF NOT EXISTS profesores
	(
		id SERIAL PRIMARY KEY,
		apellidos VARCHAR(200) NOT NULL,
		nombres VARCHAR(200) NOT NULL,
		usuarioId INT REFERENCES usuarios(id),

		activo BOOLEAN NOT NULL,
		createdAt TIMESTAMPTZ NOT NULL,
		updatedAt TIMESTAMPTZ NOT NULL
	)
`

func EnsureTableProfesorExists(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), tableProfesorCreationQuery)
	if err != nil {
		log.Printf("TEST: error creando tabla profesores: %s", err)
	}
}

func ClearTableProfesor(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), "DELETE FROM profesores")
	if err != nil {
		log.Printf("Error deleteando contenidos de la tabla profesores %s", err)
	}
	_, err = db.Exec(context.Background(), "ALTER SEQUENCE profesores_id_seq RESTART WITH 1")
	if err != nil {
		log.Printf("Error reseteando secuencia de profesor_id %s", err)
	}
}

func AddProfesores(count int, db *pgxpool.Pool) {
	ClearTableUsuario(db)
	AddUsers(count, db)
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		now := time.Now()

		db.Exec(
			context.Background(),
			`INSERT INTO profesores(apellidos, nombres,
			usuarioId, activo, createdAt, updatedAt)
			VALUES($1, $2, $3, $4, $5, $6)`,
			"ap_test_"+strconv.Itoa(i),
			"nom_test_"+strconv.Itoa(i),
			i+1, i%2 == 0, now, now)
	}
}
