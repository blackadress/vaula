package models

import "time"

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
