package models

import "time"

// faltan campos?
type Trabajo struct {
	ID       int   `json:"id"`
	Curso_id int   `json:"cursoId"`
	Curso    Curso `json:"curso"`

	Activo      bool      `json:"activo"`
	FechaInicio time.Time `json:"fechaInicio"`
	FechaFinal  time.Time `json:"fechaFinal"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
