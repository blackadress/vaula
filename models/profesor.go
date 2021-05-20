package models

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Profesor struct {
	ID        int    `json:"id"`
	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos"`
	UsuarioId int    `json:"usuarioId"`
	Usuario   User   `json:"usuario"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (p *Profesor) CreateProfesor(db *pgxpool.Pool) error {
	now := time.Now()
	return db.QueryRow(
		context.Background(),
		`INSERT INTO profesores(nombres, apellidos, usuarioId,
		activo, createdAt, updatedAt)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		p.Nombres, p.Apellidos, p.UsuarioId,
		p.Activo, now, now).Scan(&p.ID)
}

func (p *Profesor) GetProfesor(db *pgxpool.Pool) error {
	return db.QueryRow(
		context.Background(),
		`SELECT nombres, apellidos, usuarioId,
		activo, createdAt, updatedAt
		FROM profesores
		WHERE id=$1`,
		p.ID).Scan(&p.Nombres, &p.Apellidos, &p.UsuarioId,
		&p.Activo, &p.CreatedAt, &p.UpdatedAt)
}

func GetProfesores(db *pgxpool.Pool) ([]Profesor, error) {
	rows, err := db.Query(
		context.Background(),
		`SELECT id, nombres, apellidos, usuarioId,
		activo, createdAt, updatedAt
		FROM profesores`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	profesores := []Profesor{}
	for rows.Next() {
		var prof Profesor
		err := rows.Scan(
			&prof.ID, &prof.Nombres, &prof.Apellidos,
			&prof.UsuarioId, &prof.Activo, &prof.CreatedAt, &prof.UpdatedAt)
		if err != nil {
			log.Println("Las filas obtenidas de la BD para Profesor, no satisfacen a 'Scan'")
			return nil, err
		}
		profesores = append(profesores, prof)
	}

	return profesores, nil
}

func (p *Profesor) UpdateProfesor(db *pgxpool.Pool) error {
	updTime := time.Now()
	_, err := db.Exec(
		context.Background(),
		`UPDATE profesores SET nombres=$1, apellidos=$2,
		usuarioId=$3, activo=$4, updatedAt=$5
		WHERE id=$6`,
		p.Nombres, p.Apellidos,
		p.UsuarioId, p.Activo, updTime, p.ID)

	return err
}

func (p *Profesor) DeleteProfesor(db *pgxpool.Pool) error {
	_, err := db.Exec(
		context.Background(),
		`DELETE FROM profesores WHERE id=$1`,
		p.ID)

	return err
}

