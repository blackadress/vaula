package models

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Curso struct {
	ID       int    `json:"id"`
	Nombre   int    `json:"nombre"`
	Siglas   int    `json:"siglas"`
	Silabo   string `json:"silabo"`
	Semestre string `json:"semestre"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (c *Curso) CreateCurso(db *pgxpool.Pool) error {
	now := time.Now()
	return db.QueryRow(
		context.Background(),
		`INSERT INTO cursos(siglas, nombre, silabo, semestre, activo, createdAt, updatedAt)
		VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		c.Siglas, c.Nombre, c.Silabo, c.Semestre,
		c.Activo, now, now,
	).Scan(&c.ID)
}

func (c *Curso) GetCurso(db *pgxpool.Pool) error {
	return db.QueryRow(
		context.Background(),
		`SELECT siglas, nombre, silabo, semestre, activo, createdAt, updatedAt
		FROM cursos
		WHERE id=$1`,
		c.ID,
	).Scan(&c.Siglas, &c.Nombre, &c.Silabo, &c.Semestre, &c.Activo, &c.CreatedAt, &c.UpdatedAt)
}

func getCursos(db *pgxpool.Pool) ([]Curso, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT siglas, nombre, silabo, semestre, activo, createdAt, updatedAt
		FROM cursos`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cursos := []Curso{}

	for rows.Next() {
		var c Curso
		err := rows.Scan(
			&c.ID, &c.Siglas, &c.Nombre, &c.Silabo,
			&c.Semestre, &c.Activo, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			log.Println("Las filas obtenidas de la BD para Curso, no satisfacen a 'Scan'")
			return nil, err
		}
		cursos = append(cursos, c)
	}
	return cursos, nil
}

func (c *Curso) UpdateCurso(db *pgxpool.Pool) error {
	updTime := time.Now()
	_, err := db.Exec(
		context.Background(),
		`UPDATE cursos SET siglas=$1, nombre=$2, silabo=$3, 
		semestre=$4, activo=$5, updatedA=$6
		WHERE id=$7`,
		c.Siglas, c.Nombre, c.Silabo, c.Semestre, c.Activo, updTime, c.ID,
	)

	return err
}

func (c *Curso) DeleteCurso(db *pgxpool.Pool) error {
	_, err := db.Exec(
		context.Background(),
		`DELETE FROM cursos WHERE id=$1`,
		c.ID,
	)
	return err
}
