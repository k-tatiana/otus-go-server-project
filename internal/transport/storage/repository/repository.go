package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMaxConns          = 100
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
	Replicas *pgxpool.Pool
}

func NewRepo(ctx context.Context, masterConfig *pgxpool.Config, slavesConfig *pgxpool.Config, useReplicas bool) (*Repo, error) {
	masterPool, err := pgxpool.NewWithConfig(ctx, masterConfig)
	if err != nil {
		return nil, err
	}
	if err := masterPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to master database: %w", err)
	}

	var replicaPool *pgxpool.Pool
	if useReplicas {
		p, err := pgxpool.NewWithConfig(ctx, slavesConfig)
		if err != nil {
			log.Fatalf("failed to connect to replica: %v", err)
		}
		if err := p.Ping(ctx); err != nil {
			log.Fatalf("failed to ping replica: %v", err)
		}
		replicaPool = p
		log.Printf("Connected to replica: %s", slavesConfig.ConnConfig.Host)
	} else {
		replicaPool = masterPool
	}

	return &Repo{
		Master:   masterPool,
		Replicas: replicaPool,
	}, nil
}
