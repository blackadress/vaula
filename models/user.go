package models

import (
	"database/sql"
)

type user struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func GetUsers(db *sql.DB) ([]user, error) {
	rows, err := db.Query(
		`SELECT id, username, password
        FROM users`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []user{}

	for rows.Next() {
		var u user
		if err := rows.Scan(
			&u.ID, &u.Username, &u.Password, &u.Email,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
