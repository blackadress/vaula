package utils

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// USUARIOS
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

// CURSO
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
	ClearTablePreguntaTrabajo(db)
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

// EXAMEN
const tableExamenCreationQuery = `
CREATE TABLE IF NOT EXISTS examenes
	(
		id SERIAL PRIMARY KEY,
		nombre VARCHAR(200) NOT NULL,
		fechaInicio TIMESTAMPTZ NOT NULL,
		fechaFinal TIMESTAMPTZ NOT NULL,
		cursoId INT REFERENCES cursos(id),

		activo BOOLEAN NOT NULL,
		createdAt TIMESTAMPTZ NOT NULL,
		updatedAt TIMESTAMPTZ NOT NULL
	)
`

func EnsureTableExamenExists(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), tableExamenCreationQuery)
	if err != nil {
		log.Printf("TEST: error creando tabla examenes: %s", err)
	}
}

func ClearTableExamen(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), "DELETE FROM examenes")
	if err != nil {
		log.Printf("Error deleteando contenidos de la tabla Examen %s", err)
	}
	_, err = db.Exec(context.Background(), "ALTER SEQUENCE examenes_id_seq RESTART WITH 1")
	if err != nil {
		log.Printf("Error reseteando secuencia de examen_id %s", err)
	}

}

func AddExamenes(count int, db *pgxpool.Pool) {
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
			`INSERT INTO examenes(nombre, fechaInicio, fechaFinal,
				cursoId, activo, createdAt, updatedAt)
			VALUES($1, $2, $3, $4, $5, $6, $7)`,
			"examen_test_"+strconv.Itoa(i),
			fechaInicio, fechaFinal,
			i+1, i%2 == 0, now, now)

		if err != nil {
			log.Printf("Error adding examenes %s", err)
		}
	}
}

// PREGUNTA
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

// PREGUNTA TRABAJO
const tablePreguntaTrabajoCreationQuery = `
CREATE TABLE IF NOT EXISTS preguntasTrabajo
	(
		id SERIAL PRIMARY KEY,
		enunciado TEXT NOT NULL,
		trabajoId INT REFERENCES cursos(id),

		activo BOOLEAN NOT NULL,
		createdAt TIMESTAMPTZ NOT NULL,
		updatedAt TIMESTAMPTZ NOT NULL
	)
`

func EnsureTablePreguntaTrabajoExists(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), tablePreguntaTrabajoCreationQuery)
	if err != nil {
		log.Printf("TEST: error creando tabla preguntasTrabajo: %s", err)
	}
}

func ClearTablePreguntaTrabajo(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), "DELETE FROM preguntasTrabajo")
	if err != nil {
		log.Printf("Error deleteando contenidos de la tabla PreguntasTrabajo %s", err)
	}
	_, err = db.Exec(context.Background(), "ALTER SEQUENCE preguntasTrabajo_id_seq RESTART WITH 1")
	if err != nil {
		log.Printf("Error reseteando secuencia de pregunta_id %s", err)
	}

}

func AddPreguntaTrabajos(count int, db *pgxpool.Pool) {
	AddTrabajos(count, db)
	if count < 1 {
		count = 1
	}
	now := time.Now()

	for i := 0; i < count; i++ {
		_, err := db.Exec(
			context.Background(),
			`INSERT INTO preguntasTrabajo(enunciado, trabajoId, 
				activo, createdAt, updatedAt)
			VALUES($1, $2, $3, $4, $5)`,
			"preg_enun_test_"+strconv.Itoa(i),
			i+1, i%2 == 0, now, now)

		if err != nil {
			log.Printf("Error adding preguntasTrabajo %s", err)
		}
	}
}

// PROFESOR
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

// TRABAJO
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

// ALTERNATIVAS
const tableAlternativaCreationQuery = `
CREATE TABLE IF NOT EXISTS alternativas
	(
		id SERIAL PRIMARY KEY,
		valor VARCHAR(100) NOT NULL,
		correcto BOOLEAN NOT NULL,

		activo BOOLEAN NOT NULL,
		createdAt TIMESTAMPTZ NOT NULL,
		updatedAt TIMESTAMPTZ NOT NULL
	)
`

func EnsureTableAlternativaExists(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), tableAlternativaCreationQuery)
	if err != nil {
		log.Printf("TEST: error creando tabla alternativas: %s", err)
	}
}

func ClearTableAlternativa(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), "DELETE FROM alternativas")
	if err != nil {
		log.Printf("Error deleteando contenidos de la tabla Alternativa %s", err)
	}
	_, err = db.Exec(context.Background(), "ALTER SEQUENCE alternativas_id_seq RESTART WITH 1")
	if err != nil {
		log.Printf("Error reseteando secuencia de alternativa_id %s", err)
	}

}

func AddAlternativas(count int, db *pgxpool.Pool) {
	AddCursos(count, db)
	if count < 1 {
		count = 1
	}
	now := time.Now()

	for i := 0; i < count; i++ {
		_, err := db.Exec(
			context.Background(),
			`INSERT INTO alternativas(valor, correcto, activo, createdAt, updatedAt)
			VALUES($1, $2, $3, $4, $5)`,
			"alternativa_valor_"+strconv.Itoa(i),
			i%2 == 0, i%2 == 0, now, now)

		if err != nil {
			log.Printf("Error adding trabajos %s", err)
		}
	}
}
