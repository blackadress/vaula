package models

import (
	"database/sql"
	"time"
)

type Token struct {
	ID          int       `json:"id"`
	AccessToken string    `json:"access_token"`
	AccessUuid  string    `json:"access_uuid"`
	AtExpires   time.Time `json:"at_expires"`
}

func (t *Token) GetToken(db *sql.DB) error {
	return db.QueryRow(
		`SELECT access_token, access_uuid, at_expires
        FROM tokens
        WHERE id=$1`,
		t.ID,
	).Scan(&t.AccessToken, &t.AccessUuid, &t.AtExpires)
}

func (t *Token) UpdateToken(db *sql.DB) error {
	_, err := db.Exec(
		`UPDATE tokens SET access_token=$1, access_uuid=$2, at_expires=$3
        WHERE id=$4`,
		t.AccessToken,
		t.AccessUuid,
		t.AtExpires,
		t.ID,
	)

	return err
}

func (t *Token) CreateToken(db *sql.DB) error {
	return db.QueryRow(
		`INSERT INTO tokens(access_token, access_uuid, at_expires)
        VALUES($1, $2, $3)
        RETURNING id`,
		t.AccessToken,
		t.AccessUuid,
		t.AtExpires,
	).Scan(&t.ID)
}
