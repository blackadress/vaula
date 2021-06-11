package models

import (
	"testing"

	"github.com/blackadress/vaula/utils"
	"github.com/jackc/pgx/v4"
)

func TestCreateUser(t *testing.T) {
	utils.ClearTableUsuario(db)

	user := User{
		Username:   "username_prueba",
		Password: "user_pass_prueba",
		Email:    "user@test.ts",
		Activo:    true,
	}
	err := user.CreateUser(db)
	if err != nil {
		t.Errorf("No se creo el usuario")
	}

	if user.ID != 1 {
		t.Errorf("Se esperaba crear un usuario con ID 1. Se obtuvo %d", user.ID)
	}
}

func TestGetUser(t *testing.T) {
	utils.ClearTableUsuario(db)
	utils.AddUsers(1, db)
	user := User{ID: 1}
	err := user.GetUser(db)

	if err != nil {
		t.Errorf("Se esperaba obtener el usuario con ID 1. Se obtuvo %v", err)
	}
}

func TestNotGetUser(t *testing.T) {
	utils.ClearTableUsuario(db)
	user := User{ID: 1}
	err := user.GetUser(db)
	if err != pgx.ErrNoRows {
		t.Errorf("Se esperaba error ErrNoRows, se obtuvo diferente error. ERROR %v", err)
	}
}

func TestGetUsers(t *testing.T) {
	utils.ClearTableUsuario(db)
	utils.AddUsers(2, db)
	users, err := GetUsers(db)
	if err != nil {
		t.Errorf("Metodo user.GetUsers no funciona %s", err)
	}

	if len(users) != 2 {
		t.Errorf("Se esperaba obtener un array de 2 elementos. Se obtuvo: %v", users)
	}
}

func TestGetZeroUsers(t *testing.T) {
	utils.ClearTableUsuario(db)

	users, err := GetUsers(db)
	if err != nil {
		t.Errorf("Medodo user.GetUsers no funciona %s", err)
	}

	if len(users) != 0 {
		t.Errorf("Se esperaba obtener un array vacia. Se obtuvo: %v", users)
	}
}

func TestUpdateUser(t *testing.T) {
	utils.ClearTableUsuario(db)
	utils.ClearTableUsuario(db)
	utils.AddUsers(1, db)
	utils.AddUsers(1, db)

	original_user := User{ID: 1}
	err := original_user.GetUser(db)
	if err != nil {
		t.Errorf("El metodo GetUser fallo %s", err)
	}

	user_upd := User{
		ID:        1,
		Username:   "username_prueba_upd",
		Password: "user_pass_prueba_upd",
		Email:    "user@test.ts_upd",
		Activo:    false,
	}
	err = user_upd.UpdateUser(db)
	if err != nil {
		t.Errorf("El metodo UpdateUser fallo %s", err)
	}

	err = user_upd.GetUser(db)
	if err != nil {
		t.Errorf("El metodo GetUser fallo para user_upd %s", err)
	}

	if original_user.ID != user_upd.ID {
		t.Errorf("Se esperaba que el ID no cambiara, cambio de '%d' a '%d'",
			original_user.ID, user_upd.ID)
	}

	if original_user.Username == user_upd.Username {
		t.Errorf("Se esperaba que los Username cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_user.Username, user_upd.Username, original_user.Username)
	}

	if original_user.Password == user_upd.Password {
		t.Errorf("Se esperaba que los Password cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_user.Password, user_upd.Password, original_user.Password)
	}

	if original_user.Email == user_upd.Email {
		t.Errorf("Se esperaba que los Email cambiaran de '%s' a '%s'. Se obtuvo %s",
			original_user.Email, user_upd.Email, original_user.Email)
	}

	if original_user.Activo == user_upd.Activo {
		t.Errorf("Se esperaba que los Activo cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_user.Activo, user_upd.Activo, original_user.Activo)
	}

	if original_user.CreatedAt != user_upd.CreatedAt {
		t.Errorf("Se esperaba que los CreatedAt cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_user.CreatedAt, user_upd.CreatedAt, original_user.CreatedAt)
	}

	if original_user.UpdatedAt == user_upd.UpdatedAt {
		t.Errorf("Se esperaba que los UpdatedAt cambiaran de '%v' a '%v'. Se obtuvo %v",
			original_user.UpdatedAt, user_upd.UpdatedAt, original_user.UpdatedAt)
	}
}

func TestDeleteUser(t *testing.T) {
	utils.ClearTableUsuario(db)
	utils.AddUsers(1, db)

	user := User{ID: 1}
	err := user.DeleteUser(db)
	if err != nil {
		t.Errorf("Ocurrio un error en el metodo DeleteUser")
	}
}
