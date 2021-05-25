package models

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// faltan campos?
type Trabajo struct {
	ID          int       `json:"id"`
	Descripcion string    `json:"descripcion"`
	FechaInicio time.Time `json:"fechaInicio"`
	FechaFinal  time.Time `json:"fechaFinal"`
	CursoId     int       `json:"cursoId"`
	Curso       Curso     `json:"curso"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (t *Trabajo) CreateTrabajo(db *pgxpool.Pool) error {
	now := time.Now()
	return db.QueryRow(
		context.Background(),
		`INSERT INTO trabajos(descripcion, cursoId, activo,
		fechaInicio, fechaFinal, createdAt, updatedAt)
		VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		t.Descripcion, t.CursoId, t.Activo, t.FechaInicio,
		t.FechaFinal, now, now).Scan(&t.ID)
}

func (t *Trabajo) GetProfesor(db *pgxpool.Pool) error {
	return db.QueryRow(
		context.Background(),
		`SELECT descripcion, cursoId, activo, fechaInicio,
		fechaFinal, CreatedAt, UpdatedAt
		FROM trabajos
		WHERE id=$1`,
		t.ID).Scan(&t.Descripcion, &t.CursoId, &t.Activo,
		&t.FechaInicio, &t.FechaFinal, &t.CreatedAt, &t.UpdatedAt)
}

func GetTrabajos(db *pgxpool.Pool) ([]Trabajo, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT id, descripcion, cursoId, activo,
		fechaInicio, fechaFinal, CreatedAt, UpdatedAt
		FROM trabajos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trabajos := []Trabajo{}
	for rows.Next() {
		var tra Trabajo
		err := rows.Scan(
			&tra.ID, &tra.Descripcion, &tra.CursoId,
			&tra.Activo, &tra.FechaInicio, &tra.FechaFinal,
			&tra.CreatedAt, &tra.UpdatedAt)
		if err != nil {
			log.Println("Las filas obtenidas de la BD para Trabajo, no satisfacen a 'Scan'")
			return nil, err
		}
		trabajos = append(trabajos, tra)
	}

	return trabajos, nil

}

func (t *Trabajo) UpdateTrabajo(db *pgxpool.Pool) error {
	updTime := time.Now()
	_, err := db.Exec(
		context.Background(),
		`UPDATE trabajos SET descripcion=$1, cursoId=$2, activo=$3,
		fechaInicio=$4, fechaFinal=$5, updatedAt=$6
		WHERE id=$7`,
		t.Descripcion, t.CursoId, t.Activo, t.FechaInicio,
		t.FechaFinal, updTime)

	return err
}

func (t *Trabajo) DeleteTrabajo(db *pgxpool.Pool) error {
	_, err := db.Exec(
		context.Background(),
		`DELETE FROM trabajos WHERE id=$1`,
		t.ID)
	return err
}
