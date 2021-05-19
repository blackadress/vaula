package models

import "time"

type Curso struct {
	ID       int    `json:"id"`
	Siglas   int    `json:"siglas"`
	Silabo   string `json:"silabo"`
	Semestre string `json:"semestre"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
