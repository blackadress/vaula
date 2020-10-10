package models

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (u *User) GetUser(db *sql.DB) error {
	return db.QueryRow(
		`SELECT username, password, email
        FROM users
        WHERE id=$1`,
		u.ID,
	).Scan(&u.Username, &u.Password, &u.Email)
}

func (u *User) GetUserByUsername(db *sql.DB) error {
	return db.QueryRow(
		`SELECT id, password, email
        FROM users
        WHERE username=$1`,
		u.Username,
	).Scan(&u.ID, &u.Password, &u.Email)
}

func (u *User) UpdateUser(db *sql.DB) error {
	_, err := db.Exec(
		`UPDATE users SET username=$1, password=$2, email=$3
        WHERE id=$4`,
		u.Username,
		u.Password,
		u.Email,
		u.ID,
	)

	return err
}

func (u *User) DeleteUser(db *sql.DB) error {
	_, err := db.Exec(
		`DELETE FROM users WHERE id=$1`,
		u.ID,
	)
	return err
}

func (u *User) CreateUser(db *sql.DB) error {
	return db.QueryRow(
		`INSERT INTO users(username, password, email)
        VALUES($1, $2, $3)
        RETURNING id`,
		u.Username,
		u.Password,
		u.Email,
	).Scan(&u.ID)
}

func GetUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query(
		`SELECT id, username, password, email
        FROM users`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID, &u.Username, &u.Password, &u.Email,
		); err != nil {
			fmt.Println("line 35")
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
