package models

import "time"

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
