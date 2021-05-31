package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)


var db *pgxpool.Pool

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("TEST: '.env' no encontrado")
		os.Exit(1)
	}
	user := os.Getenv("APP_DB_USERNAME")
	password := os.Getenv("APP_DB_PASSWORD")
	dbname := os.Getenv("APP_DB_NAME")
	connectionString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", user, password, dbname)

	db, err = pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		log.Printf("%v, no se pudo conectar a la base de datos %v", os.Stderr, err)
		os.Exit(1)
	}

	EnsureTableUsuarioExists(db)
	EnsureTableAlumnoExists(db)
	EnsureTableProfesorExists(db)
	EnsureTableExamenExists(db)
	EnsureTableCursoExists(db)
	EnsureTablePreguntaExists(db)
	EnsureTableTrabajoExists(db)
	EnsureTablePreguntaTrabajoExists(db)

	code := m.Run()

	// ClearTableAlumno(db) ya reseteado por ClearTableUsuario
	// ClearTableProfesor(db) ya reseteado por ClearTableUsuario
	ClearTableUsuario(db)
	ClearTableCurso(db)
	// ClearTableExamen(db) ya reseteado por ClearTableCurso
	// ClearTablePregunta(db) ya reseteado por ClearTableCurso
	// ClearTableTrabajo(db) ya reseteado por ClearTableCurso
	// ClearTablePreguntaTrabajo(db) ya reseteado por ClearTableCurso
	os.Exit(code)

}
