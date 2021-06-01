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
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`

	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (u *User) GetUser(db *pgxpool.Pool) error {
	return db.QueryRow(
		context.Background(),
		`SELECT username, password, email,
		activo, createdAt, updatedAt
		FROM usuarios
		WHERE id=$1`,
		u.ID,
	).Scan(&u.Username, &u.Password, &u.Email,
		&u.Activo, &u.CreatedAt, &u.UpdatedAt)
}

func (u *User) GetUserByUsername(db *pgxpool.Pool) error {
	return db.QueryRow(context.Background(),
		`SELECT id, password, email,
		activo, createdAt, updatedAt
		FROM usuarios
		WHERE username=$1`,
		u.Username,
	).Scan(&u.ID, &u.Password, &u.Email,
		&u.Activo, &u.CreatedAt, &u.UpdatedAt)
}

func (u *User) GetUserNoPwd(db *pgxpool.Pool) error {
	return db.QueryRow(context.Background(),
		`SELECT id, email, activo, createdAt, updatedAt
		FROM usuarios
		WHERE id=$1`,
		u.ID,
	).Scan(&u.ID, &u.Email, &u.Activo,
		&u.CreatedAt, &u.UpdatedAt)
}

func (u *User) UpdateUser(db *pgxpool.Pool) error {
	now := time.Now()
	_, err := db.Exec(context.Background(),
		`UPDATE usuarios SET username=$1, password=$2, email=$3,
		activo=$4, updatedAt=$5
		WHERE id=$6`,
		u.Username, u.Password, u.Email,
		u.Activo, now, u.ID,
	)

	return err
}

func (u *User) DeleteUser(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(),
		`DELETE FROM usuarios WHERE id=$1`,
		u.ID,
	)
	return err
}

func (u *User) CreateUser(db *pgxpool.Pool) error {
	now := time.Now()
	return db.QueryRow(context.Background(),
		`INSERT INTO usuarios(username, password, email,
		activo, createdAt, updatedAt)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id, createdAt, updatedAt`,
		u.Username, u.Password, u.Email,
		u.Activo, now, now).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func GetUsers(db *pgxpool.Pool) ([]User, error) {
	rows, err := db.Query(context.Background(),
		`SELECT id, username, email, activo, createdAt, updatedAt
		FROM usuarios`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.ID, &u.Username, &u.Email,
			&u.Activo, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			log.Printf("The rows we got from the DB can't be 'Scan'(ed) %s", err)
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

// vale verga cuando se le da un token que no es
// revisar el como se handlean los errores
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
