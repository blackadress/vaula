package models

import "time"

type Pregunta struct {
	ID        int    `json:"id"`
	Enunciado string `json:"Enunciado"`
	ExamenID  int    `json:"examenId"`
	Examen    Examen `json:"examen"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
