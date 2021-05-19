package models

import "time"

type Alumno struct {
	ID        int    `json:"id"`
	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos"`
	Codigo    string `json:"codigo"` // 8 characteres
	UsuarioId int    `json:"usuarioId"`
	Usuario   User   `json:"usuario"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AlumnoCurso struct {
	ID           int       `json:"id"`
	Calificacion float32   `json:"calificacion"`
	FechaInicio  time.Time `json:"fechaInicio"`
	FechaFinal   time.Time `json:"fechaFinal"`
	AlumnoId     int       `json:"alumnoId"`
	Alumno       Alumno    `json:"alumno"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AlumnoExamen struct {
	ID           int       `json:"id"`
	Calificacion float32   `json:"calificacion"`
	FechaInicio  time.Time `json:"fechaInicio"`
	FechaFinal   time.Time `json:"fechaFinal"`
	AlumnoId     int       `json:"alumnoId"`
	Alumno       Alumno    `json:"alumno"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AlumnoTrabajo struct {
	ID           int       `json:"id"`
	Calificacion float32   `json:"calificacion"`
	Uri          string    `json:"uri"`
	FechaInicio  time.Time `json:"fechaInicio"`
	FechaFinal   time.Time `json:"fechaFinal"`
	AlumnoId     int       `json:"alumnoId"`
	Alumno       Alumno    `json:"alumno"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
