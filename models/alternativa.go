package models

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Alternativa struct {
	ID       int    `json:"id"`
	Valor    string `json:"valor"`
	Correcto bool   `json:"correcto"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (a *Alternativa) CreateAlternativa(db *pgxpool.Pool) error {
	now := time.Now()
	return db.QueryRow(
		context.Background(),
		`INSERT INTO alternativas(valor, correcto, activo, createdAt, updatedAt)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id`,
		a.Valor, a.Correcto, true, now, now,
	).Scan(&a.ID)
}

func (a *Alternativa) GetAlternativa(db *pgxpool.Pool) error {
	return db.QueryRow(
		context.Background(),
		`SELECT valor, correcto, activo, createdAt, updatedAt
		`,
		a.ID,
	).Scan(&a.Valor, &a.Correcto, &a.Activo, &a.CreatedAt, &a.UpdatedAt)
}

func (a *Alternativa) GetAlternativas(db *pgxpool.Pool) ([]Alternativa, error) {
	rows, err := db.Query(context.Background(),
		`SELECT id, valor, correcto, activo, createdAt, updatedAt
		FROM alternativas`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	alternativas := []Alternativa{}

	for rows.Next() {
		var a Alternativa
		err := rows.Scan(
			&a.ID, &a.Valor, &a.Correcto, &a.Activo, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			log.Println("Las filas obtenidas de la BD no satisfacen la función 'Scan'")
			return nil, err
		}
		alternativas = append(alternativas, a)
	}

	return alternativas, nil
}

func (a *Alternativa) UpdateAlternativa(db *pgxpool.Pool) error {
	updTime := time.Now()
	_, err := db.Exec(
		context.Background(),
		`UPDATE alternativas SET valor=$1, correcto=$2, updatedAt=$3
		WHERE id=$4`,
		a.Valor, a.Correcto, updTime, a.ID,
	)

	return err
}

func (a *Alternativa) DeleteAlternativa(db *pgxpool.Pool) error {
	_, err := db.Exec(
		context.Background(),
		`DELETE FROM alternativas WHERE id=$1`,
		a.ID,
	)

	return err
}
