package utils

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

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

// ALUMNOS
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
