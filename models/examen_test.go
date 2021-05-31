package models

import (
	"testing"
	"time"

	"github.com/blackadress/vaula/utils"
	"github.com/jackc/pgx/v4"
)

func TestCreateExamen(t *testing.T) {
	utils.ClearTableCurso(db) // limpia la tabla examenes tambien
	utils.AddCursos(1, db)

	loc, _ := time.LoadLocation("America/Lima")
	fechaInicio := time.Date(2021, time.June,
		20, 18, 0, 0, 0, loc)
	fechaFinal := time.Date(2021, time.June,
		22, 18, 0, 0, 0, loc)

	exa := Examen{
		Nombre:      "nom_examen_prueba",
		FechaInicio: fechaInicio,
		FechaFinal:  fechaFinal,
		CursoId:     1,
		Activo:      true,
	}
	err := exa.CreateExamen(db)
	if err != nil {
		t.Errorf("No se creo el examen %s", err)
	}

	if exa.ID != 1 {
		t.Errorf("Se esperaba crear un examen con ID 1. Se obtuvo %d", exa.ID)
	}
}

func TestGetExamen(t *testing.T) {
	utils.ClearTableCurso(db) // limpia la tabla examenes
	utils.AddExamenes(1, db)

	exa := Examen{ID: 1}
	err := exa.GetExamen(db)

	if err != nil {
		t.Errorf("Se esperaba obtener el examen con ID 1. Se obtuvo %v", err)
	}
}

func TestNotGetExamen(t *testing.T) {
	utils.ClearTableExamen(db)
	exa := Examen{ID: 1}
	err := exa.GetExamen(db)
	if err != pgx.ErrNoRows {
		t.Errorf("Se esperaba no ErrNoRows, se obtuvo diferente error. ERROR %v", err)
	}
}

func TestGetExamens(t *testing.T) {
	utils.ClearTableExamen(db)
	utils.AddExamenes(2, db)
	examenes, err := GetExamenes(db)
	if err != nil {
		t.Errorf("algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(examenes) != 2 {
		t.Errorf("Se esperaba obtener un array de 2 elementos. Se obtuvo: %v", examenes)
	}

}

func TestGetZeroExamens(t *testing.T) {
	utils.ClearTableExamen(db)

	examenes, err := GetExamenes(db)
	if err != nil {
		t.Errorf("Algo salio mal con la comunicacion con la DB %s", err)
	}

	if len(examenes) != 0 {
		t.Errorf("Se esperaba obtener un array vacia. Se obtuvo: %v", examenes)
	}
}

func TestUpdateExamen(t *testing.T) {
	utils.ClearTableCurso(db)
	utils.AddExamenes(1, db)
	utils.AddCursos(1, db)

	original_ex := Examen{ID: 1}
	err := original_ex.GetExamen(db)
	if err != nil {
		t.Errorf("El metodo GetExamen fallo %s", err)
	}

	loc, _ := time.LoadLocation("America/Lima")
	fechaInicio := time.Date(2023, time.June,
		20, 18, 0, 0, 0, loc)
	fechaFinal := time.Date(2023, time.June,
		22, 18, 0, 0, 0, loc)

	exa_upd := Examen{
		ID:          1,
		Nombre:      "nom_examen_prueba_upd",
		FechaInicio: fechaInicio,
		FechaFinal:  fechaFinal,
		CursoId:     2,
		Activo:      false,
	}
	err = exa_upd.UpdateExamen(db)
	if err != nil {
		t.Errorf("El metodo UpdateExamen fallo %s", err)
	}

	err = exa_upd.GetExamen(db)
	if err != nil {
		t.Errorf("El metodo GetExamen fallo para exa_upd %s", err)
	}

	if original_ex.ID != exa_upd.ID {
		t.Errorf("Se esperaba que el ID no cambiara, cambio de '%d' a '%d'",
			original_ex.ID, exa_upd.ID)
	}

	if original_ex.Nombre == exa_upd.Nombre {
		t.Errorf("Se esperaba que los Nombre cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_ex.Nombre, exa_upd.Nombre, original_ex.Nombre)
	}

	if original_ex.FechaInicio == exa_upd.FechaInicio {
		t.Errorf("Se esperaba que los FechaInicio cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_ex.FechaInicio, exa_upd.FechaInicio, original_ex.FechaInicio)
	}

	if original_ex.FechaFinal == exa_upd.FechaFinal {
		t.Errorf("Se esperaba que los FechaFinal cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_ex.FechaFinal, exa_upd.FechaFinal, original_ex.FechaFinal)
	}

	if original_ex.CursoId == exa_upd.CursoId {
		t.Errorf("Se esperaba que los CursoId cambiaran de '%d' a '%d'. Se obtuvo %d",
			original_ex.CursoId, exa_upd.CursoId, original_ex.CursoId)
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

func TestDeleteExamen(t *testing.T) {
	utils.ClearTableExamen(db)
	utils.AddExamenes(1, db)

	exa := Examen{ID: 1}
	err := exa.DeleteExamen(db)
	if err != nil {
		t.Errorf("Ocurrio un error en el metodo DeleteExamen")
	}
}
