package models

import (
	"testing"

	"github.com/blackadress/vaula/utils"
	"github.com/jackc/pgx/v4"
)

func TestCreateAlumno(t *testing.T) {
	utils.ClearTableUsuario(db)
	utils.ClearTableAlumno(db)
	utils.AddUsers(1, db)

	al := Alumno{
		Nombres:   "nom_al prueba",
		Apellidos: "ap_al prueba",
		Codigo:    "11111111",
		UsuarioId: 1,
		Activo:    true,
	}
	err := al.CreateAlumno(db)
	if err != nil {
		t.Errorf("No se creo el alumno")
	}

	if al.ID != 1 {
		t.Errorf("Se esperaba crear un alumno con ID 1. Se obtuvo %d", al.ID)
	}
}

func TestGetAlumno(t *testing.T) {
	utils.ClearTableAlumno(db)
	utils.AddAlumnos(1, db)
	al := Alumno{ID: 1}
	err := al.GetAlumno(db)

	if err != nil {
		t.Errorf("Se esperaba obtener el alumno con ID 1. Se obtuvo %v", err)
	}
}

func TestNotGetAlumno(t *testing.T) {
	utils.ClearTableAlumno(db)
	al := Alumno{ID: 1}
	err := al.GetAlumno(db)
	if err != pgx.ErrNoRows {
		t.Errorf("Se esperaba error ErrNoRows, se obtuvo diferente error. ERROR %v", err)
	}
}

func TestGetAlumnos(t *testing.T) {
	utils.ClearTableAlumno(db)
	utils.AddAlumnos(2, db)
	alumnos, err := GetAlumnos(db)
	if err != nil {
		t.Errorf("Metodo alumno.GetAlumnos no funciona %s", err)
	}

	if len(alumnos) != 2 {
		t.Errorf("Se esperaba obtener un array de 2 elementos. Se obtuvo: %v", alumnos)
	}
}

func TestGetZeroAlumnos(t *testing.T) {
	utils.ClearTableAlumno(db)

	alumnos, err := GetAlumnos(db)
	if err != nil {
		t.Errorf("Metodo alumnos.GetAlumnos no funciona %s", err)
	}

	if len(alumnos) != 0 {
		t.Errorf("Se esperaba obtener un array vacia. Se obtuvo: %v", alumnos)
	}
}

func TestUpdateAlumno(t *testing.T) {
	utils.ClearTableAlumno(db)
	utils.ClearTableUsuario(db)
	utils.AddAlumnos(1, db)
	utils.AddUsers(1, db)

	original_al := Alumno{ID: 1}
	err := original_al.GetAlumno(db)
	if err != nil {
		t.Errorf("El metodo GetAlumno fallo %s", err)
	}

	al_upd := Alumno{
		ID:        1,
		Nombres:   "nom_al upd",
		Apellidos: "ap_al upd",
		Codigo:    "22222222",
		UsuarioId: 2,
		Activo:    false,
	}
	err = al_upd.UpdateAlumno(db)
	if err != nil {
		t.Errorf("El metodo UpdateAlumno fallo %s", err)
	}

	err = al_upd.GetAlumno(db)
	if err != nil {
		t.Errorf("El metodo GetAlumno fallo para al_upd %s", err)
	}

	if original_al.ID != al_upd.ID {
		t.Errorf("Se esperaba que el ID no cambiara, cambio de '%d' a '%d'",
			original_al.ID, al_upd.ID)
	}

	if original_al.Nombres == al_upd.Nombres {
		t.Errorf("Se esperaba que los Nombres cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_al.Nombres, al_upd.Nombres, original_al.Nombres)
	}

	if original_al.Apellidos == al_upd.Apellidos {
		t.Errorf("Se esperaba que los Apellidos cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_al.Apellidos, al_upd.Apellidos, original_al.Apellidos)
	}

	if original_al.Codigo == al_upd.Codigo {
		t.Errorf("Se esperaba que los Codigo cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_al.Codigo, al_upd.Codigo, original_al.Codigo)
	}

	if original_al.UsuarioId == al_upd.UsuarioId {
		t.Errorf("Se esperaba que los UsuarioId cambiaran de '%d' a '%d'. Se obtuvo %d",
			original_al.UsuarioId, al_upd.UsuarioId, original_al.UsuarioId)
	}

	if original_al.Activo == al_upd.Activo {
		t.Errorf("Se esperaba que los Activo cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_al.Activo, al_upd.Activo, original_al.Activo)
	}

	if original_al.CreatedAt != al_upd.CreatedAt {
		t.Errorf("Se esperaba que los CreatedAt cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_al.CreatedAt, al_upd.CreatedAt, original_al.CreatedAt)
	}

	if original_al.UpdatedAt == al_upd.UpdatedAt {
		t.Errorf("Se esperaba que los UpdatedAt cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_al.UpdatedAt, al_upd.UpdatedAt, original_al.UpdatedAt)
	}
}

func TestDeleteAlumno(t *testing.T) {
	utils.ClearTableAlumno(db)
	utils.AddAlumnos(1, db)

	al := Alumno{ID: 1}
	err := al.DeleteAlumno(db)
	if err != nil {
		t.Errorf("Ocurrio un error en el metodo DeleteAlumno")
	}
}
