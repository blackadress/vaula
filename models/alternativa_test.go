package models

import (
	"testing"

	"github.com/blackadress/vaula/utils"
	"github.com/jackc/pgx/v4"
)

func TestCreateAlternativa(t *testing.T) {
	utils.ClearTableAlternativa(db)

	al := Alternativa{
		Valor:    "valor_alt_prueba",
		Correcto: true,
		Activo:   true,
	}
	err := al.CreateAlternativa(db)
	if err != nil {
		t.Errorf("No se creo el curso %s", err)
	}

	if al.ID != 1 {
		t.Errorf("Se esperaba crear un curso con ID 1. Se obtuvo %d", al.ID)
	}
}

func TestGetAlternativa(t *testing.T) {
	utils.ClearTableAlternativa(db)
	utils.AddAlternativas(1, db)
	al := Alternativa{ID: 1}
	err := al.GetAlternativa(db)

	if err != nil {
		t.Errorf("Se esperaba obtener el alternativa con ID 1. Se obtuvo %v", err)
	}
}

func TestNotGetAlternativa(t *testing.T) {
	utils.ClearTableAlternativa(db)
	al := Alternativa{ID: 1}
	err := al.GetAlternativa(db)
	if err != pgx.ErrNoRows {
		t.Errorf("Se esperaba no obtener ningun curso, se obtuvo diferente error. ERROR %v", err)
	}
}

func TestGetAlternativas(t *testing.T) {
	utils.ClearTableAlternativa(db)
	utils.AddAlternativas(2, db)
	cursos, err := GetAlternativas(db)
	if err != nil {
		t.Errorf("Algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(cursos) != 2 {
		t.Errorf("Se esperaba obtener un array de 2 elementos. Se obtuvo: %v", cursos)
	}

}

func TestGetZeroAlternativas(t *testing.T) {
	utils.ClearTableAlternativa(db)

	cursos, err := GetAlternativas(db)
	if err != nil {
		t.Errorf("Algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(cursos) != 0 {
		t.Errorf("Se esperaba obtener un array vacia. Se obtuvo: %v", cursos)
	}
}

func TestUpdateAlternativa(t *testing.T) {
	utils.ClearTableAlternativa(db)
	utils.AddAlternativas(1, db)

	original_al := Alternativa{ID: 1}
	err := original_al.GetAlternativa(db)
	if err != nil {
		t.Errorf("El metodo GetAlternativa fallo %s", err)
	}

	al_upd := Alternativa{
		ID:       1,
		Valor:    "valor_alt_prueba_upd",
		Correcto: false,
		Activo:   false,
	}
	err = al_upd.UpdateAlternativa(db)
	if err != nil {
		t.Errorf("El metodo UpdateAlternativa fallo %s", err)
	}

	err = al_upd.GetAlternativa(db)
	if err != nil {
		t.Errorf("El metodo GetAlternativa fallo para al_upd %s", err)
	}

	if original_al.ID != al_upd.ID {
		t.Errorf("Se esperaba que el ID no cambiara, cambio de '%d' a '%d'",
			original_al.ID, al_upd.ID)
	}

	if original_al.Valor == al_upd.Valor {
		t.Errorf("Se esperaba que Valor cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_al.Valor, al_upd.Valor, original_al.Valor)
	}

	if original_al.Correcto == al_upd.Correcto {
		t.Errorf("Se esperaba que Correcto cambiaran de '%t' a '%t'. Se obtuvo %t",
			original_al.Correcto, al_upd.Correcto, original_al.Correcto)
	}

	if original_al.Activo == al_upd.Activo {
		t.Errorf("Se esperaba que Activo cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_al.Activo, al_upd.Activo, original_al.Activo)
	}

	if original_al.CreatedAt != al_upd.CreatedAt {
		t.Errorf("Se esperaba que CreatedAt no cambiara de '%v' a '%v'. Se obtuvo %v",
			original_al.CreatedAt, al_upd.CreatedAt, original_al.CreatedAt)
	}

	if original_al.UpdatedAt == al_upd.UpdatedAt {
		t.Errorf("Se esperaba que UpdatedAt cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_al.UpdatedAt, al_upd.UpdatedAt, original_al.UpdatedAt)
	}
}

func TestDeleteAlternativa(t *testing.T) {
	utils.ClearTableAlternativa(db)
	utils.AddAlternativas(1, db)

	al := Alternativa{ID: 1}
	err := al.DeleteAlternativa(db)
	if err != nil {
		t.Errorf("Ocurrio un error en el metodo DeleteAlternativa")
	}
}
