package models

import (
	"testing"

	"github.com/blackadress/vaula/utils"
	"github.com/jackc/pgx/v4"
)

func TestCreatePreguntaTrabajo(t *testing.T) {
	utils.ClearTableTrabajo(db) // limpia la tabla preguntaTrabajo tambien

	exa := PreguntaTrabajo{
		Enunciado: "enun_preg_trab_prueba",
		TrabajoId: 1,
		Activo:    true,
	}
	err := exa.CreatePreguntaTrabajo(db)
	if err != nil {
		t.Errorf("No se creo la pregunta_trabajo %s", err)
	}

	if exa.ID != 1 {
		t.Errorf("Se esperaba crear una pregunta_trabajo con ID 1. Se obtuvo %d", exa.ID)
	}
}

func TestGetPreguntaTrabajo(t *testing.T) {
	utils.ClearTableTrabajo(db) // limpia la tabla preguntasTrabajo
	utils.AddPreguntaTrabajos(1, db)

	exa := PreguntaTrabajo{ID: 1}
	err := exa.GetPreguntaTrabajo(db)

	if err != nil {
		t.Errorf("Se esperaba obtener la pregunta_trabajo con ID 1. Se obtuvo %v", err)
	}
}

func TestNotGetPreguntaTrabajo(t *testing.T) {
	utils.ClearTablePreguntaTrabajo(db)
	exa := PreguntaTrabajo{ID: 1}
	err := exa.GetPreguntaTrabajo(db)
	if err != pgx.ErrNoRows {
		t.Errorf("Se esperaba obtener ErrNoRows, se obtuvo diferente error. ERROR %v", err)
	}
}

func TestGetPreguntaTrabajos(t *testing.T) {
	utils.ClearTablePreguntaTrabajo(db)
	utils.AddPreguntaTrabajos(2, db)
	preguntasTrabajo, err := GetPreguntasTrabajo(db)
	if err != nil {
		t.Errorf("algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(preguntasTrabajo) != 2 {
		t.Errorf("Se esperaba obtener un array de 2 elementos. Se obtuvo: %v", preguntasTrabajo)
	}
}

func TestGetZeroPreguntaTrabajos(t *testing.T) {
	utils.ClearTablePreguntaTrabajo(db)

	preguntasTrabajo, err := GetPreguntasTrabajo(db)
	if err != nil {
		t.Errorf("Algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(preguntasTrabajo) != 0 {
		t.Errorf("Se esperaba obtener un array vacia. Se obtuvo: %v", preguntasTrabajo)
	}
}

func TestUpdatePreguntaTrabajo(t *testing.T) {
	utils.ClearTableTrabajo(db) // limpia la tabla preguntasTrabajo tambien
	utils.AddPreguntaTrabajos(1, db)

	original_ex := PreguntaTrabajo{ID: 1}
	err := original_ex.GetPreguntaTrabajo(db)
	if err != nil {
		t.Errorf("El metodo GetPreguntaTrabajo fallo %s", err)
	}

	exa_upd := PreguntaTrabajo{
		ID:        1,
		Enunciado: "enun_preg_trab_prueba_upd",
		TrabajoId: 2,
		Activo:    false,
	}
	err = exa_upd.UpdatePreguntaTrabajo(db)
	if err != nil {
		t.Errorf("El metodo UpdatePreguntaTrabajo fallo %s", err)
	}

	err = exa_upd.GetPreguntaTrabajo(db)
	if err != nil {
		t.Errorf("El metodo GetPreguntaTrabajo fallo para exa_upd %s", err)
	}

	if original_ex.ID != exa_upd.ID {
		t.Errorf("Se esperaba que el ID no cambiara, cambio de '%d' a '%d'",
			original_ex.ID, exa_upd.ID)
	}

	if original_ex.Enunciado == exa_upd.Enunciado {
		t.Errorf("Se esperaba que los Enunciado cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_ex.Enunciado, exa_upd.Enunciado, original_ex.Enunciado)
	}

	if original_ex.TrabajoId == exa_upd.TrabajoId {
		t.Errorf("Se esperaba que los TrabajoId cambiaran de '%d' a '%d'. Se obtuvo %d",
			original_ex.TrabajoId, exa_upd.TrabajoId, original_ex.TrabajoId)
	}

	if original_ex.Activo == exa_upd.Activo {
		t.Errorf("Se esperaba que los Activo cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_ex.Activo, exa_upd.Activo, original_ex.Activo)
	}

	if original_ex.CreatedAt != exa_upd.CreatedAt {
		t.Errorf("Se esperaba que CreatedAt no cambiara de '%v' a '%v'. Se obtuvo %v",
			original_ex.CreatedAt, exa_upd.CreatedAt, original_ex.CreatedAt)
	}

	if original_ex.UpdatedAt == exa_upd.UpdatedAt {
		t.Errorf("Se esperaba que UpdatedAt cambiara de '%v' a '%v'. Se obtuvo %v",
			original_ex.UpdatedAt, exa_upd.UpdatedAt, original_ex.UpdatedAt)
	}
}

func TestDeletePreguntaTrabajo(t *testing.T) {
	utils.ClearTablePreguntaTrabajo(db)
	utils.AddPreguntaTrabajos(1, db)

	exa := PreguntaTrabajo{ID: 1}
	err := exa.DeletePreguntaTrabajo(db)
	if err != nil {
		t.Errorf("Ocurrio un error en el metodo DeletePreguntaTrabajo")
	}
}
