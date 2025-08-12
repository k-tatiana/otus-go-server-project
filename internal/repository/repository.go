package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"otus/go-server-project/internal/models"
)

const (
	defaultMaxConns          = 1000
	defaultMinConns          = 10
	defaultMaxConnLifetime   = time.Hour * 1
	defaultMaxConnIdleTime   = time.Minute * 30
	defaultHealthCheckPeriod = time.Minute
	defaultConnectTimeout    = time.Second * 5
)

func Config(url string) *pgxpool.Config {
	dbConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal("Failed to create a config, error", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	return dbConfig
}

type Repo struct {
	Master   *pgxpool.Pool
	Replicas []*pgxpool.Pool
}

func NewRepo(ctx context.Context, masterConfig *pgxpool.Config, slavesConfig []*pgxpool.Config, useReplicas bool) (*Repo, error) {
	masterPool, err := pgxpool.NewWithConfig(ctx, masterConfig)
	if err != nil {
		return nil, err
	}
	if err := masterPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to master database: %w", err)
	}

	var replicaPools []*pgxpool.Pool
	if useReplicas {
		for _, rCfg := range slavesConfig {
			p, err := pgxpool.NewWithConfig(ctx, rCfg)
			if err != nil {
				log.Printf("failed to connect to replica: %v", err)
				continue
			}
			if err := p.Ping(ctx); err != nil {
				log.Printf("failed to ping replica: %v", err)
				continue
			}
			log.Printf("Connected to replica: %s", rCfg.ConnConfig.Host)
			replicaPools = append(replicaPools, p)
		}
	} else {
		replicaPools = []*pgxpool.Pool{masterPool}
	}

	return &Repo{
		Master:   masterPool,
		Replicas: replicaPools,
	}, nil
}

// Login checks user credentials and returns a token if valid.
func (r *Repo) Login(ctx context.Context, login, passwordHash string) (string, error) {
	// acquire a connection from the master pool
	conn, err := r.Master.Acquire(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()
	// Start a transaction to ensure atomicity
	tx, err := conn.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err := tx.Commit(ctx); err != nil {
			tx.Rollback(ctx)
		}
	}()

	token, err := r.loginWithReturnToken(ctx, tx, login, passwordHash)
	if err != nil {
		return "", err
	}

	err = r.saveToken(ctx, tx, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

// loginWithReturnToken checks the credentials and returns the token if valid.
func (r *Repo) loginWithReturnToken(ctx context.Context, tx pgx.Tx, login, passwordHash string) (string, error) {
	var (
		token    string
		pwd_hash string
	)
	err := tx.QueryRow(
		ctx,
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
func (r *Repo) saveToken(ctx context.Context, tx pgx.Tx, token string) error {
	_, err := tx.Exec(
		ctx,
		"INSERT INTO sessions (token, expiration_time) VALUES ($1, $2)",
		token, time.Now().Add(24*time.Hour), // Example expiration time
	)
	return err
}

func (r *Repo) RegisterUser(ctx context.Context, u models.UserDTO) (string, error) {
	conn, err := r.Master.Acquire(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	row := conn.QueryRow(
		ctx,
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
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return "", errors.New("user already exists")
		}
		return "", fmt.Errorf("inserting user %w", err)
	}

	return token, nil
}

func (r *Repo) Get(ctx context.Context, id string) (models.UserDTO, error) {
	// acquire a connection from random one of the replicas
	idx := rand.Intn(len(r.Replicas))
	conn, err := r.Replicas[idx].Acquire(ctx)
	if err != nil {
		return models.UserDTO{}, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()
	// Use the connection to query the user
	var u models.UserDTO
	err = conn.QueryRow(
		ctx,
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

func (r *Repo) ValidateToken(ctx context.Context, token string) error {
	idx := rand.Intn(len(r.Replicas))
	conn, err := r.Replicas[idx].Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()
	// Check if the token exists and is valid
	var count int
	err = conn.QueryRow(
		ctx,
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

func (r *Repo) SearchUser(ctx context.Context, name, surname string) ([]models.UserDTO, error) {
	users := make([]models.UserDTO, 0)

	idx := rand.Intn(len(r.Replicas))
	conn, err := r.Replicas[idx].Acquire(ctx)
	if err != nil {
		return users, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		`SELECT name, surname, birthday, gender, interests, city, login FROM users 
		WHERE lower(name) LIKE $1 AND lower(surname) LIKE $2
		ORDER BY token`,
		strings.ToLower(name)+"%",
		strings.ToLower(surname)+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.UserDTO
		if err := rows.Scan(
			&user.Name,
			&user.Surname,
			&user.Birthday,
			&user.Gender,
			&user.Interests,
			&user.City,
			&user.Login,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}
