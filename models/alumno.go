package models

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Alumno struct {
	ID        int    `json:"id"`
	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos"`
	Codigo    string `json:"codigo"` // 8 characteres
	UsuarioId int    `json:"usuarioId"`
	Usuario   User   `json:"usuario"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (a *Alumno) CreateAlumno(db *pgxpool.Pool) error {
	now := time.Now()
	return db.QueryRow(
		context.Background(),
		`INSERT INTO alumnos(nombres, apellidos, codigo, usuarioId, createdAt, updatedAt)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		a.Nombres, a.Apellidos, a.Codigo, a.UsuarioId, now, now,
	).Scan(&a.ID)
}

func (a *Alumno) GetAlumno(db *pgxpool.Pool) error {
	return db.QueryRow(
		context.Background(),
		`SELECT nombres, apellidos, codigo, usuarioId, createdAt, updatedAt
		FROM alumnos
		WHERE id=$1`,
		a.ID,
	).Scan(&a.Nombres, &a.Apellidos, &a.Codigo, &a.UsuarioId, &a.CreatedAt, &a.UpdatedAt)
}

func (a *Alumno) GetAlumnos(db *pgxpool.Pool) ([]Alumno, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT nombres, apellidos, codigo, usuarioId, createdAt, updatedAt
		FROM alumnos`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	alumnos := []Alumno{}

	for rows.Next() {
		var a Alumno
		err := rows.Scan(
			&a.ID, &a.Nombres, &a.Apellidos, &a.Codigo,
			&a.UsuarioId, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			log.Println("Las filas obtenidas de la BD para Alumno, no satisfacen a 'Scan'")
			return nil, err
		}
		alumnos = append(alumnos, a)
	}
	return alumnos, nil
}

func (a *Alumno) UpdateAlternativa(db *pgxpool.Pool) error {
	now := time.Now()
	_, err := db.Exec(
		context.Background(),
		`UPDATE alumnos SET nombres=$1, apellidos=$2, codigo=$3, usuarioId=$4, updatedAt=$5
		WHERE id=$6`,
		a.Nombres, a.Apellidos, a.Codigo, a.UsuarioId, now, a.ID,
	)

	return err
}

func (a *Alumno) DeleteAlumno(db *pgxpool.Pool) error {
	_, err := db.Exec(
		context.Background(),
		`DELETE FROM alumnos WHERE id=$1`,
		a.ID,
	)

	return err
}

type AlumnoCurso struct {
	ID           int       `json:"id"`
	Calificacion float32   `json:"calificacion"`
	FechaInicio  time.Time `json:"fechaInicio"`
	FechaFinal   time.Time `json:"fechaFinal"`
	AlumnoId     int       `json:"alumnoId"`
	Alumno       Alumno    `json:"alumno"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (ac *AlumnoCurso) CreateAlumnoCurso(db *pgxpool.Pool) error {
	now := time.Now()
	return db.QueryRow(
		context.Background(),
		`INSERT INTO alumnoCurso(calificacion, fechaInicio, 
		fechaFinal, alumnoId, activo, createdAt, updatedAt)
		VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		ac.Calificacion, ac.FechaInicio, ac.FechaFinal,
		ac.AlumnoId, now, now,
	).Scan(&ac.Calificacion, &ac.FechaInicio, &ac.FechaFinal,
		&ac.AlumnoId, &ac.CreatedAt, &ac.UpdatedAt)
}

func (ac *AlumnoCurso) GetAlumnoCurso(db *pgxpool.Pool) error {
	return db.QueryRow(
		context.Background(),
		`SELECT calificacion, fechaInicio, fechaFinal, alumnoId, createdAt, updatedAt
		FROM alumnoCurso
		WHERE id=$1`,
		ac.ID,
	).Scan(&ac.Calificacion, &ac.FechaInicio, &ac.FechaFinal, &ac.AlumnoId, &ac.CreatedAt, &ac.UpdatedAt)
}

func (ac *AlumnoCurso) GetAlumnoCursos(db *pgxpool.Pool) ([]AlumnoCurso, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT id, calificacion, fechaInicio, fechaFinal,
		alumnoId, createdAt, updatedAt
		FROM alumnoCurso`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	alumnoCursos := []AlumnoCurso{}

	for rows.Next() {
		var ac AlumnoCurso
		err := rows.Scan(
			&ac.ID, &ac.Calificacion, &ac.FechaInicio, &ac.FechaFinal,
			&ac.AlumnoId, &ac.CreatedAt, &ac.UpdatedAt)
		if err != nil {
			log.Println("Las filas obtenidas de la BD para Alumno Curso, no satisfacen a 'Scan'")
			return nil, err
		}
		alumnoCursos = append(alumnoCursos, ac)
	}

	return alumnoCursos, nil
}

func (ac *AlumnoCurso) UpdateAlumnoCurso(db *pgxpool.Pool) error {
	updTime := time.Now()
	_, err := db.Exec(
		context.Background(),
		`UPDATE alumnoCurso SET calificacion=$1, fechaInicio=$2,
		fechaFinal=$3, alumnoId=$4, updatedAt=$5
		WHERE id=$1`,
		ac.Calificacion, ac.FechaInicio, ac.FechaFinal, ac.AlumnoId, updTime,
	)
	return err
}

func (ac *AlumnoCurso) DeleteAlumnoCurso(db *pgxpool.Pool) error {
	_, err := db.Exec(
		context.Background(),
		`DELETE FROM alumnoCurso WHERE id=$1`,
		ac.ID,
	)

	return err
}

type AlumnoExamen struct {
	ID           int       `json:"id"`
	Calificacion float32   `json:"calificacion"`
	FechaInicio  time.Time `json:"fechaInicio"`
	FechaFinal   time.Time `json:"fechaFinal"`
	AlumnoId     int       `json:"alumnoId"`
	Alumno       Alumno    `json:"alumno"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (ae *AlumnoExamen) CreateAlumnoExamen(db *pgxpool.Pool) error {
	now := time.Now()
	return db.QueryRow(
		context.Background(),
		`INSERT INTO alumnoExamen(calificacion, fechaInicio, 
		fechaFinal, alumnoId, activo, createdAt, updatedAt)
		VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		ae.Calificacion, ae.FechaInicio, ae.FechaFinal,
		ae.AlumnoId, now, now,
	).Scan(&ae.Calificacion, &ae.FechaInicio, &ae.FechaFinal,
		&ae.AlumnoId, &ae.CreatedAt, &ae.UpdatedAt)
}

func (ae *AlumnoExamen) GetAlumnoExamen(db *pgxpool.Pool) error {
	return db.QueryRow(
		context.Background(),
		`SELECT calificacion, fechaInicio, fechaFinal, alumnoId, createdAt, updatedAt
		FROM alumnoExamen
		WHERE id=$1`,
		ae.ID,
	).Scan(&ae.Calificacion, &ae.FechaInicio, &ae.FechaFinal, &ae.AlumnoId, &ae.CreatedAt, &ae.UpdatedAt)
}

func (ae *AlumnoExamen) GetAlumnoExamenes(db *pgxpool.Pool) ([]AlumnoExamen, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT id, calificacion, fechaInicio, fechaFinal,
		alumnoId, createdAt, updatedAt
		FROM alumnoExamen`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	alumnoExamens := []AlumnoExamen{}

	for rows.Next() {
		var ae AlumnoExamen
		err := rows.Scan(
			&ae.ID, &ae.Calificacion, &ae.FechaInicio, &ae.FechaFinal,
			&ae.AlumnoId, &ae.CreatedAt, &ae.UpdatedAt)
		if err != nil {
			log.Println("Las filas obtenidas de la BD para Alumno Curso, no satisfacen a 'Scan'")
			return nil, err
		}
		alumnoExamens = append(alumnoExamens, ae)
	}

	return alumnoExamens, nil
}

func (ae *AlumnoExamen) UpdateAlumnoExamen(db *pgxpool.Pool) error {
	updTime := time.Now()
	_, err := db.Exec(
		context.Background(),
		`UPDATE alumnoExamen SET calificacion=$1, fechaInicio=$2,
		fechaFinal=$3, alumnoId=$4, updatedAt=$5
		WHERE id=$1`,
		ae.Calificacion, ae.FechaInicio, ae.FechaFinal, ae.AlumnoId, updTime,
	)
	return err
}

func (ae *AlumnoExamen) DeleteAlumnoExamen(db *pgxpool.Pool) error {
	_, err := db.Exec(
		context.Background(),
		`DELETE FROM alumnoExamen WHERE id=$1`,
		ae.ID,
	)

	return err
}

type AlumnoTrabajo struct {
	ID           int       `json:"id"`
	Calificacion float32   `json:"calificacion"`
	Uri          string    `json:"uri"`
	FechaInicio  time.Time `json:"fechaInicio"`
	FechaFinal   time.Time `json:"fechaFinal"`
	AlumnoId     int       `json:"alumnoId"`
	Alumno       Alumno    `json:"alumno"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (at *AlumnoTrabajo) CreateAlumnoTrabajo(db *pgxpool.Pool) error {
	now := time.Now()
	return db.QueryRow(
		context.Background(),
		`INSERT INTO alumnoTrabajo(calificacion, uri, fechaInicio, 
		fechaFinal, alumnoId, activo, createdAt, updatedAt)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`,
		at.Calificacion, at.Uri, at.FechaInicio,
		at.FechaFinal, at.AlumnoId, now, now,
	).Scan(&at.Calificacion, &at.FechaInicio, &at.FechaFinal,
		&at.AlumnoId, &at.CreatedAt, &at.UpdatedAt)
}

func (at *AlumnoTrabajo) GetAlumnoTrabajo(db *pgxpool.Pool) error {
	return db.QueryRow(
		context.Background(),
		`SELECT calificacion, uri, fechaInicio, fechaFinal, alumnoId, createdAt, updatedAt
		FROM alumnoTrabajo
		WHERE id=$1`,
		at.ID,
	).Scan(&at.Calificacion, &at.Uri, &at.FechaInicio,
		&at.FechaFinal, &at.AlumnoId, &at.CreatedAt, &at.UpdatedAt)
}

func (at *AlumnoTrabajo) GetAlumnoTrabajoes(db *pgxpool.Pool) ([]AlumnoTrabajo, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT id, calificacion, uri, fechaInicio,
		fechaFinal, alumnoId, createdAt, updatedAt
		FROM alumnoTrabajo`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	alumnoTrabajos := []AlumnoTrabajo{}

	for rows.Next() {
		var at AlumnoTrabajo
		err := rows.Scan(
			&at.ID, &at.Calificacion, &at.Uri, &at.FechaInicio,
			&at.FechaFinal, &at.AlumnoId, &at.CreatedAt, &at.UpdatedAt)
		if err != nil {
			log.Println("Las filas obtenidas de la BD para Alumno Curso, no satisfacen a 'Scan'")
			return nil, err
		}
		alumnoTrabajos = append(alumnoTrabajos, at)
	}

	return alumnoTrabajos, nil
}

func (at *AlumnoTrabajo) UpdateAlumnoTrabajo(db *pgxpool.Pool) error {
	updTime := time.Now()
	_, err := db.Exec(
		context.Background(),
		`UPDATE alumnoTrabajo SET calificacion=$1, uri=$2, fechaInicio=$3,
		fechaFinal=$4, alumnoId=$5, updatedAt=$6
		WHERE id=$1`,
		at.Calificacion, at.Uri, at.FechaInicio, at.FechaFinal, at.AlumnoId, updTime,
	)
	return err
}

func (at *AlumnoTrabajo) DeleteAlumnoTrabajo(db *pgxpool.Pool) error {
	_, err := db.Exec(
		context.Background(),
		`DELETE FROM alumnoTrabajo WHERE id=$1`,
		at.ID,
	)

	return err
}
