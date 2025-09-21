package datastore

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type PostgresDatastore struct {
	Client *pgxpool.Pool
}

func (d *PostgresDatastore) InitConnection(ctx context.Context, cfg *DBConfig) error {
	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)
	var err error
	d.Client, err = pgxpool.New(ctx, dbConnStr)
	if err != nil {
		return err
	}
	err = d.Client.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to establish connection with db: %w", err)
	}
	return nil
}
