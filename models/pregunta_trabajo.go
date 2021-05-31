package models

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PreguntaTrabajo struct {
	ID        int     `json:"id"`
	Enunciado string  `json:"Enunciado"`
	TrabajoId int     `json:"trabajoId"`
	Trabajo   Trabajo `json:"trabajo"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (pt *PreguntaTrabajo) CreatePreguntaTrabajo(db *pgxpool.Pool) error {
	now := time.Now()
	return db.QueryRow(
		context.Background(),
		`INSERT INTO preguntasTrabajo(enunciado, trabajoId, activo, createdAt, updatedAt)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id`,
		pt.Enunciado, pt.TrabajoId, true, now, now).Scan(&pt.ID)
}

func (pt *PreguntaTrabajo) GetPreguntaTrabajo(db *pgxpool.Pool) error {
	return db.QueryRow(
		context.Background(),
		`SELECT enunciado, trabajoId, activo, createdAt, updatedAt
		FROM preguntasTrabajo
		WHERE id=$1`,
		pt.ID).Scan(&pt.Enunciado, &pt.TrabajoId,
		&pt.Activo, &pt.CreatedAt, &pt.UpdatedAt)
}

func GetPreguntasTrabajo(db *pgxpool.Pool) ([]PreguntaTrabajo, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT id, enunciado, trabajoId, activo, createdAt, updatedAt
		FROM preguntasTrabajo`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	preguntasTrabajo := []PreguntaTrabajo{}
	for rows.Next() {
		var pt PreguntaTrabajo
		err := rows.Scan(
			&pt.ID, &pt.Enunciado, &pt.TrabajoId,
			&pt.Activo, &pt.CreatedAt, &pt.UpdatedAt)
		if err != nil {
			log.Printf("Las filas obtenidas de la BD para PreguntaTrabajo, no satisfacen a 'Scan' %s",
				err)
			return nil, err
		}
		preguntasTrabajo = append(preguntasTrabajo, pt)
	}
	return preguntasTrabajo, nil
}

func (pt *PreguntaTrabajo) UpdatePreguntaTrabajo(db *pgxpool.Pool) error {
	updTime := time.Now()
	_, err := db.Exec(
		context.Background(),
		`UPDATE preguntasTrabajo SET enunciado=$1, trabajoId=$2,
		activo=$3, updatedAt=$4`,
		pt.Enunciado, pt.TrabajoId, pt.Activo, updTime)
	return err
}

func (pt *PreguntaTrabajo) DeletePreguntaTrabajo(db *pgxpool.Pool) error {
	_, err := db.Exec(
		context.Background(),
		`DELETE FROM preguntasTrabajo WHERE id=$1`,
		pt.ID)
	return err
}
