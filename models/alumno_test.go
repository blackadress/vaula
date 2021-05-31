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

func TestCreateAlumno(t *testing.T) {
	ClearTableUsuario(db)
	ClearTableAlumno(db)
	AddUsers(1, db)

	al := Alumno{
		Nombres:   "nom_al prueba",
		Apellidos: "ap_al prueba",
		Codigo:    "11111111",
		UsuarioId: 1,
		Activo:    true,
	}
	err := al.CreateAlumno(db)
	if err != nil {
		t.Errorf("No se creo el alumno")
	}

	if al.ID != 1 {
		t.Errorf("Se esperaba crear un alumno con ID 1. Se obtuvo %d", al.ID)
	}
}

func TestGetAlumno(t *testing.T) {
	ClearTableAlumno(db)
	AddAlumnos(1, db)
	al := Alumno{ID: 1}
	err := al.GetAlumno(db)

	if err != nil {
		t.Errorf("Se esperaba obtener el alumno con ID 1. Se obtuvo %v", err)
	}
}

func TestNotGetAlumno(t *testing.T) {
	ClearTableAlumno(db)
	al := Alumno{ID: 1}
	err := al.GetAlumno(db)
	if err != pgx.ErrNoRows {
		t.Errorf("Se esperaba no obtener ningun alumno, se obtuvo diferente error. ERROR %v", err)
	}
}

func TestGetAlumnos(t *testing.T) {
	ClearTableAlumno(db)
	AddAlumnos(2, db)
	alumnos, err := GetAlumnos(db)
	if err != nil {
		t.Errorf("Algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(alumnos) != 2 {
		t.Errorf("Se esperaba obtener un array de 2 elementos. Se obtuvo: %v", alumnos)
	}
}

func TestGetZeroAlumnos(t *testing.T) {
	ClearTableAlumno(db)

	alumnos, err := GetAlumnos(db)
	if err != nil {
		t.Errorf("Algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(alumnos) != 0 {
		t.Errorf("Se esperaba obtener un array vacia. Se obtuvo: %v", alumnos)
	}
}

func TestUpdateAlumno(t *testing.T) {
	ClearTableAlumno(db)
	ClearTableUsuario(db)
	AddAlumnos(1, db)
	AddUsers(1, db)

	original_al := Alumno{ID: 1}
	err := original_al.GetAlumno(db)
	if err != nil {
		t.Errorf("El metodo GetAlumno fallo %s", err)
	}

	al_upd := Alumno{
		ID:        1,
		Nombres:   "nom_al upd",
		Apellidos: "ap_al upd",
		Codigo:    "22222222",
		UsuarioId: 2,
		Activo:    false,
	}
	err = al_upd.UpdateAlumno(db)
	if err != nil {
		t.Errorf("El metodo UpdateAlumno fallo %s", err)
	}

	err = al_upd.GetAlumno(db)
	if err != nil {
		t.Errorf("El metodo GetAlumno fallo para al_upd %s", err)
	}

	if original_al.ID != al_upd.ID {
		t.Errorf("Se esperaba que el ID no cambiara, cambio de '%d' a '%d'",
			original_al.ID, al_upd.ID)
	}

	if original_al.Nombres == al_upd.Nombres {
		t.Errorf("Se esperaba que los Nombres cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_al.Nombres, al_upd.Nombres, original_al.Nombres)
	}

	if original_al.Apellidos == al_upd.Apellidos {
		t.Errorf("Se esperaba que los Apellidos cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_al.Apellidos, al_upd.Apellidos, original_al.Apellidos)
	}

	if original_al.Codigo == al_upd.Codigo {
		t.Errorf("Se esperaba que los Codigo cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_al.Codigo, al_upd.Codigo, original_al.Codigo)
	}

	if original_al.UsuarioId == al_upd.UsuarioId {
		t.Errorf("Se esperaba que los UsuarioId cambiaran de '%d' a '%d'. Se obtuvo %d",
			original_al.UsuarioId, al_upd.UsuarioId, original_al.UsuarioId)
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

func TestDeleteAlumno(t *testing.T) {
	ClearTableAlumno(db)
	AddAlumnos(1, db)

	al := Alumno{ID: 1}
	err := al.DeleteAlumno(db)
	if err != nil {
		t.Errorf("Ocurrio un error en el metodo DeleteAlumno")
	}
}

const tableAlumnoCreationQuery = `
CREATE TABLE IF NOT EXISTS alumnos
	(
		id SERIAL PRIMARY KEY,
		apellidos VARCHAR(200) NOT NULL,
		nombres VARCHAR(200) NOT NULL,
		codigo CHAR(8) NOT NULL,
		usuarioId INT REFERENCES usuarios(id),

		activo BOOLEAN NOT NULL,
		createdAt TIMESTAMPTZ,
		updatedAt TIMESTAMPTZ
	)
`

func EnsureTableAlumnoExists(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), tableAlumnoCreationQuery)
	if err != nil {
		log.Printf("TEST: error creando tabla alumnos: %s", err)
	}
}

func ClearTableAlumno(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), "DELETE FROM alumnos")
	if err != nil {
		log.Printf("Error deleteando contenidos de la tabla alumno %s", err)
	}
	_, err = db.Exec(context.Background(), "ALTER SEQUENCE alumnos_id_seq RESTART WITH 1")
	if err != nil {
		log.Printf("Error reseteando secuencia de alumno_id %s", err)
	}
}

func AddAlumnos(count int, db *pgxpool.Pool) {
	ClearTableUsuario(db)
	AddUsers(count, db)
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		now := time.Now()
		codigo := fmt.Sprintf("%.8s", strconv.Itoa(i)+"00000000")

		db.Exec(
			context.Background(),
			`INSERT INTO alumnos(apellidos, nombres, codigo,
			usuarioId, activo, createdAt, updatedAt)
			VALUES($1, $2, $3, $4, $5, $6, $7)`,
			"ap_test_"+strconv.Itoa(i),
			"nom_test_"+strconv.Itoa(i),
			codigo, i+1, i%2 == 0, now, now)
	}
}

// USUARIOS TEST

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS usuarios
	(
		id INT PRIMARY KEY NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		email TEXT NOT NULL,

		activo BOOLEAN,
		createdAt TIMESTAMPTZ NOT NULL,
		updatedAt TIMESTAMPTZ NOT NULL
	)
`

func EnsureTableUsuarioExists(db *pgxpool.Pool) {
	if _, err := db.Exec(context.Background(), tableCreationQuery); err != nil {
		log.Printf("TEST: error creando tabla de usuarios: %s", err)
	}
}

func ClearTableUsuario(db *pgxpool.Pool) {
	ClearTableAlumno(db)
	ClearTableProfesor(db)
	_, err := db.Exec(context.Background(), "DELETE FROM usuarios")
	if err != nil {
		log.Printf("Error deleteando contenidos de la tabla usuarios %s", err)
	}
	_, err = db.Exec(context.Background(), "ALTER SEQUENCE usuarios_id_seq RESTART WITH 1")
	if err != nil {
		log.Printf("Error alterando secuencia de usuario_id %s", err)
	}
}

func AddUsers(count int, db *pgxpool.Pool) {
	now := time.Now()
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		db.Exec(
			context.Background(),
			`INSERT INTO usuarios(username, password, email, activo, createdAt, updatedAt)
			VALUES($1, $2, $3, $4, $5, $6)`,
			"user_"+strconv.Itoa(i),
			"pass"+strconv.Itoa(i),
			"em"+strconv.Itoa(i)+"@test.ts",
			i%2 == 0, now, now,
		)
	}
}
