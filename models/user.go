package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (u *User) GetUser(db *pgxpool.Pool) error {
	return db.QueryRow(
		context.Background(),
		`SELECT username, password, email
        FROM users
        WHERE id=$1`,
		u.ID,
	).Scan(&u.Username, &u.Password, &u.Email)
}

func (u *User) GetUserByUsername(db *pgxpool.Pool) error {
	return db.QueryRow(context.Background(),
		`SELECT id, password, email
        FROM users
        WHERE username=$1`,
		u.Username,
	).Scan(&u.ID, &u.Password, &u.Email)
}

func (u *User) UpdateUser(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(),
		`UPDATE users SET username=$1, password=$2, email=$3
        WHERE id=$4`,
		u.Username,
		u.Password,
		u.Email,
		u.ID,
	)

	return err
}

func (u *User) DeleteUser(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(),
		`DELETE FROM users WHERE id=$1`,
		u.ID,
	)
	return err
}

func (u *User) CreateUser(db *pgxpool.Pool) error {
	return db.QueryRow(context.Background(),
		`INSERT INTO users(username, password, email)
        VALUES($1, $2, $3)
        RETURNING id`,
		u.Username,
		u.Password,
		u.Email,
	).Scan(&u.ID)
}

func GetUsers(db *pgxpool.Pool) ([]User, error) {
	rows, err := db.Query(context.Background(),
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

type Claims struct {
	UserId int `json:"userId"`
	jwt.StandardClaims
}

type JWToken struct {
	UserId      int       `json:"userId"`
	AccessToken string    `json:"accessToken"`
	Expires     time.Time `json:"expires"`
}

func (u *User) GetJWTForUser() (JWToken, error) {
	var token JWToken
	expirationTime := time.Now().Add(time.Minute * 30)

	validToken, err := generateJWT(expirationTime, u.ID)
	if err != nil {
		log.Printf("Error inesperado en la capa modelo de jwt")
		return token, err
	}

	token = JWToken{
		UserId:      u.ID,
		AccessToken: validToken,
		Expires:     expirationTime,
	}
	return token, err
}

func ValidateToken(bearerToken string) (bool, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		bearerToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				err_msg := "There was an error while validating token"
				log.Printf("%s", err_msg)
				return nil, fmt.Errorf("%s", err_msg)
			}
			return secretKey, nil
		})

	return token.Valid, err
}

func generateJWT(expirationTime time.Time, userId int) (string, error) {
	claims := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// get this from env
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		log.Printf("Something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}
