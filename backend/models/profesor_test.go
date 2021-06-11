package models

import (
	"testing"

	"github.com/blackadress/vaula/utils"
	"github.com/jackc/pgx/v4"
)

func TestCreateProfesor(t *testing.T) {
	utils.ClearTableUsuario(db)
	utils.ClearTableProfesor(db)
	utils.AddUsers(1, db)

	al := Profesor{
		Nombres:   "nom_pro prueba",
		Apellidos: "ap_prof prueba",
		UsuarioId: 1,
		Activo:    true,
	}
	err := al.CreateProfesor(db)
	if err != nil {
		t.Errorf("No se creo el profesor")
	}

	if al.ID != 1 {
		t.Errorf("Se esperaba crear un profesor con ID 1. Se obtuvo %d", al.ID)
	}
}

func TestGetProfesor(t *testing.T) {
	utils.ClearTableProfesor(db)
	utils.AddProfesores(1, db)
	al := Profesor{ID: 1}
	err := al.GetProfesor(db)

	if err != nil {
		t.Errorf("Se esperaba obtener el profesor con ID 1. Se obtuvo %v", err)
	}
}

func TestNotGetProfesor(t *testing.T) {
	utils.ClearTableProfesor(db)
	al := Profesor{ID: 1}
	err := al.GetProfesor(db)
	if err != pgx.ErrNoRows {
		t.Errorf("Se esperaba ErrNoRows, se obtuvo diferente error. ERROR %v", err)
	}
}

func TestGetProfesors(t *testing.T) {
	utils.ClearTableProfesor(db)
	utils.AddProfesores(2, db)
	profesores, err := GetProfesores(db)
	if err != nil {
		t.Errorf("Algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(profesores) != 2 {
		t.Errorf("Se esperaba obtener un array de 2 elementos. Se obtuvo: %v", profesores)
	}
}

func TestGetZeroProfesores(t *testing.T) {
	utils.ClearTableProfesor(db)

	profesores, err := GetProfesores(db)
	if err != nil {
		t.Errorf("Algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(profesores) != 0 {
		t.Errorf("Se esperaba obtener un array vacia. Se obtuvo: %v", profesores)
	}
}

func TestUpdateProfesor(t *testing.T) {
	utils.ClearTableProfesor(db)
	utils.ClearTableUsuario(db)
	utils.AddProfesores(1, db)
	utils.AddUsers(1, db)

	original_prof := Profesor{ID: 1}
	err := original_prof.GetProfesor(db)
	if err != nil {
		t.Errorf("El metodo GetProfesor fallo %s", err)
	}

	al_upd := Profesor{
		ID:        1,
		Nombres:   "nom_pro upd",
		Apellidos: "ap_prof upd",
		UsuarioId: 2,
		Activo:    false,
	}
	err = al_upd.UpdateProfesor(db)
	if err != nil {
		t.Errorf("El metodo UpdateProfesor fallo %s", err)
	}

	err = al_upd.GetProfesor(db)
	if err != nil {
		t.Errorf("El metodo GetProfesor fallo para al_upd %s", err)
	}

	if original_prof.ID != al_upd.ID {
		t.Errorf("Se esperaba que el ID no cambiara, cambio de '%d' a '%d'",
			original_prof.ID, al_upd.ID)
	}

	if original_prof.Nombres == al_upd.Nombres {
		t.Errorf("Se esperaba que los Nombres cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_prof.Nombres, al_upd.Nombres, original_prof.Nombres)
	}

	if original_prof.Apellidos == al_upd.Apellidos {
		t.Errorf("Se esperaba que los Apellidos cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_prof.Apellidos, al_upd.Apellidos, original_prof.Apellidos)
	}

	if original_prof.UsuarioId == al_upd.UsuarioId {
		t.Errorf("Se esperaba que los UsuarioId cambiaran de '%d' a '%d'. Se obtuvo %d",
			original_prof.UsuarioId, al_upd.UsuarioId, original_prof.UsuarioId)
	}

	if original_prof.Activo == al_upd.Activo {
		t.Errorf("Se esperaba que los Activo cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_prof.Activo, al_upd.Activo, original_prof.Activo)
	}

	if original_prof.CreatedAt != al_upd.CreatedAt {
		t.Errorf("Se esperaba que los CreatedAt cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_prof.CreatedAt, al_upd.CreatedAt, original_prof.CreatedAt)
	}

	if original_prof.UpdatedAt == al_upd.UpdatedAt {
		t.Errorf("Se esperaba que los UpdatedAt cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_prof.UpdatedAt, al_upd.UpdatedAt, original_prof.UpdatedAt)
	}
}

func TestDeleteProfesor(t *testing.T) {
	utils.ClearTableProfesor(db)
	utils.AddProfesores(1, db)

	al := Profesor{ID: 1}
	err := al.DeleteProfesor(db)
	if err != nil {
		t.Errorf("Ocurrio un error en el metodo DeleteProfesor")
	}
}
