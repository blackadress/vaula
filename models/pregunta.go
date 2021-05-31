package models

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Pregunta struct {
	ID        int    `json:"id"`
	Enunciado string `json:"Enunciado"`
	ExamenId  int    `json:"examenId"`
	Examen    Examen `json:"examen"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (p *Pregunta) CreatePregunta(db *pgxpool.Pool) error {
	now := time.Now()
	return db.QueryRow(
		context.Background(),
		`INSERT INTO preguntas(enunciado, examenId, activo, createdAt, updatedAt)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id`,
		p.Enunciado, p.ExamenId, p.Activo, now, now).Scan(&p.ID)
}

func (p *Pregunta) GetPregunta(db *pgxpool.Pool) error {
	return db.QueryRow(
		context.Background(),
		`SELECT enunciado, examenId, activo, createdAt, updatedAt
		FROM preguntas
		WHERE id=$1`,
		p.ID).Scan(&p.Enunciado, &p.ExamenId, &p.Activo,
		&p.CreatedAt, &p.UpdatedAt)
}

func GetPreguntas(db *pgxpool.Pool) ([]Pregunta, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT id, enunciado, examenId, activo, createdAt, updatedAt
		FROM preguntas`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	preguntas := []Pregunta{}
	for rows.Next() {
		var p Pregunta
		err := rows.Scan(
			&p.ID, &p.Enunciado, &p.ExamenId,
			&p.Activo, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			log.Printf("Las filas obtenidas de la BD para Pregunta, no satisfacen a 'Scan' %s",
				err)
			return nil, err
		}
		preguntas = append(preguntas, p)
	}
	return preguntas, nil
}

func (p *Pregunta) UpdatePregunta(db *pgxpool.Pool) error {
	updTime := time.Now()
	_, err := db.Exec(
		context.Background(),
		`UPDATE preguntas SET enunciado=$1, examenId=$2,
		activo=$3, updatedAt=$4
		WHERE id=$5`,
		p.Enunciado, p.ExamenId, p.Activo, updTime, p.ID)

	return err
}

func (p *Pregunta) DeletePregunta(db *pgxpool.Pool) error {
	_, err := db.Exec(
		context.Background(),
		`DELETE FROM preguntas WHERE id=$1`,
		p.ID)

	return err
}
