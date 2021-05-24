package models

import (
	"context"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Email       string    `json:"email"`
	Activo      bool      `json:"activo"`
	FechaInicio time.Time `json:"fechaInicio"`
	FechaFinal  time.Time `json:"fechaFinal"`
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

func (u *User) GetUserNoPwd(db *pgxpool.Pool) error {
	return db.QueryRow(context.Background(),
		`SELECT id, email
		FROM users
		WHERE id=$1`,
		u.ID,
	).Scan(&u.ID, &u.Email)
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
		`SELECT id, username, email
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
			&u.ID, &u.Username, &u.Email,
		); err != nil {
			log.Println("The rows we got from the DB can't be 'Scan'(ed)")
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
	UserId       int    `json:"userId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (u *User) GetJWTForUser() (JWToken, error) {
	var token JWToken
	expirationTime := time.Now().Add(time.Minute * 30)
	// TODO test refresh time token
	// expirationTime := time.Now().Add(time.Second * 3)

	validToken, err := generateJWT(expirationTime, u.ID)
	if err != nil {
		log.Printf("Error inesperado generando access token")
		return token, err
	}

	expirationTime = time.Now().Add(time.Hour * 24 * 7)
	refreshToken, err := generateJWT(expirationTime, u.ID)
	if err != nil {
		log.Printf("Error inesperado generando refresh token")
		return token, err
	}

	token = JWToken{
		UserId:       u.ID,
		AccessToken:  validToken,
		RefreshToken: refreshToken,
	}
	return token, err
}

func ValidateToken(tkn string) (bool, Claims, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tkn,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Unexpected signing method: %v", token.Header["alg"])
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKey, nil
		})

	if err != nil {
		return token.Valid, *claims, err
	}

	return token.Valid, *claims, err
}

// mas tarde se puede generar el access token y refresh separados
// para eso en lugar de userId, podria aceptar un objeto Claims
// y revisar si tiene otros campos fuera de userId para generar
// el token
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
