package repository

import (
	"errors"
	"fmt"
	"otus/go-server-project/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

type Repo struct {
	db *pgx.ConnPool
}

func NewRepo(db *pgx.ConnPool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) OpenTx() (*pgx.Tx, error) {
	return r.db.Begin()
}

func (r *Repo) CommitOrRollback(tx *pgx.Tx) {
	if p := recover(); p != nil {
		tx.Rollback()
		panic(p)
	} else if err := tx.Commit(); err != nil {
		tx.Rollback()
	}
}

// Login checks user credentials and returns a token if valid.
func (r *Repo) Login(login, passwordHash string) (string, error) {
	tx, err := r.OpenTx()
	if err != nil {
		return "", err
	}
	defer r.CommitOrRollback(tx)

	token, err := r.loginWithReturnToken(tx, login, passwordHash)
	if err != nil {
		return "", err
	}

	err = r.saveToken(tx, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

// loginWithReturnToken checks the credentials and returns the token if valid.
func (r *Repo) loginWithReturnToken(tx *pgx.Tx, login, passwordHash string) (string, error) {
	var (
		token    string
		pwd_hash string
	)
	err := tx.QueryRow(
		"SELECT password_hash FROM users WHERE token=$1",
		login,
	).Scan(&pwd_hash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", models.ErrNoSuchUser
		}
		return "", err
	}
	if pwd_hash != passwordHash {
		return "", models.ErrInvalidCredentials
	}
	token = uuid.New().String() // Generate a new token

	return token, nil
}

// saveToken saves the token for the user (dummy implementation).
func (r *Repo) saveToken(tx *pgx.Tx, token string) error {
	_, err := tx.Exec(
		"INSERT INTO sessions (token, expiration_time) VALUES ($1, $2)",
		token, time.Now().Add(24*time.Hour), // Example expiration time
	)
	return err
}

func (r *Repo) RegisterUser(u models.UserDTO) (string, error) {
	row := r.db.QueryRow(
		`INSERT INTO users (
			name,
			surname,
			birthday,
			gender,
			interests,
			city,
			login,
			password_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING token`,
		u.Name, u.Surname, u.Birthday, u.Gender, u.Interests, u.City, u.Login, u.PasswordHash,
	)
	var token string
	if err := row.Scan(&token); err != nil {
		// Check for unique constraint violation (duplicate user)
		// Assuming 'login' is unique in the users table
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			return "", errors.New("user already exists")
		}
		return "", fmt.Errorf("inserting user %w", err)
	}

	return token, nil
}

func (r *Repo) Get(id string) (models.UserDTO, error) {
	var u models.UserDTO
	err := r.db.QueryRow(
		"SELECT name, surname, birthday, gender, interests, city, login FROM users WHERE token=$1",
		id,
	).Scan(&u.Name, &u.Surname, &u.Birthday, &u.Gender, &u.Interests, &u.City, &u.Login)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.UserDTO{}, errors.New("user not found")
		}
		return models.UserDTO{}, fmt.Errorf("failed to get user: %w", err)
	}
	return u, nil
}

func (r *Repo) ValidateToken(token string) error {
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM sessions WHERE token=$1 AND expiration_time > NOW()",
		token,
	).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}
	if count == 0 {
		return models.ErrUnauthorized
	}
	return nil
}
