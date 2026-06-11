package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Service struct {
	DB *pgxpool.Pool
}

func New(databaseURL string) (*Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return &Service{
		DB: db,
	}, nil
}

func (s *Service) Close() {
	s.DB.Close()
}
