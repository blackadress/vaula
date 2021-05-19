package models

import "time"

type PreguntaTrabajo struct {
	ID        int     `json:"id"`
	Enunciado string  `json:"Enunciado"`
	TrabajoID int     `json:"trabajoId"`
	Trabajo   Trabajo `json:"trabajo"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
