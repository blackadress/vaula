package models

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Examen struct {
	ID          int       `json:"id"`
	Nombre      string    `json:"nombre"`
	FechaInicio time.Time `json:"fechaInicio"`
	FechaFinal  time.Time `json:"fechaFinal"`
	CursoId     int       `json:"cursoId"`
	Curso       Curso     `json:"curso"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (e *Examen) CreateExamen(db *pgxpool.Pool) error {
	now := time.Now()
	return db.QueryRow(
		context.Background(),
		`INSERT INTO examenes(nombre, fechaInicio, fechaFinal,
		cursoId, activo, createdAt, updatedAt)
		VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		e.Nombre, e.FechaInicio, e.FechaFinal, e.CursoId,
		e.Activo, now, now).Scan(&e.ID)
}

func (e *Examen) GetExamen(db *pgxpool.Pool) error {
	return db.QueryRow(
		context.Background(),
		`SELECT nombre, fechaInicio, fechaFinal,
		cursoId, activo, createdAt, updatedAt
		FROM examenes
		WHERE id=$1`,
		e.ID,
	).Scan(&e.Nombre, &e.FechaInicio, &e.FechaFinal,
		&e.CursoId, &e.Activo, &e.CreatedAt, &e.UpdatedAt)
}

func GetExamenes(db *pgxpool.Pool) ([]Examen, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT id, nombre, fechaInicio, fechaFinal, cursoId,
		activo, createdAt, updatedAt
		FROM examenes`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	examenes := []Examen{}

	for rows.Next() {
		var e Examen
		err := rows.Scan(
			&e.ID, &e.Nombre, &e.FechaInicio, &e.FechaFinal, &e.CursoId,
			&e.Activo, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			log.Printf("Las filas obtenidas de la BD para Examen, no satisfacen a 'Scan' %s",
				err)
			return nil, err
		}
		examenes = append(examenes, e)
	}
	return examenes, nil
}

func (e *Examen) UpdateExamen(db *pgxpool.Pool) error {
	updTime := time.Now()
	_, err := db.Exec(
		context.Background(),
		`UPDATE examenes SET nombre=$1, fechaInicio=$2, fechaFinal=$3,
		cursoId=$4, activo=$5, updatedAt=$6`,
		e.Nombre, e.FechaInicio, e.FechaFinal, e.CursoId, e.Activo, updTime)
	return err
}

func (e *Examen) DeleteExamen(db *pgxpool.Pool) error {
	_, err := db.Exec(
		context.Background(),
		`DELETE FROM examenes WHERE id=$1`,
		e.ID)
	return err
}
