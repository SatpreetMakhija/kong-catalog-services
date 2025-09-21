package datastore

import (
	"context"
	"fmt"
)

type ServiceRepo interface {
	FindServiceById(context.Context, string) (Service, error)
}

type PostgresServiceRepo struct {
	ds *PostgresDatastore
}

func NewPostgresServiceRepo(ds *PostgresDatastore) *PostgresServiceRepo {
	return &PostgresServiceRepo{ds: ds}
}

func (r *PostgresServiceRepo) FindServiceById(ctx context.Context, id string) (Service, error) {
	query := `select id, name, description, version from services where id = $1`
	row := r.ds.Client.QueryRow(ctx, query, id)
	var service Service
	err := row.Scan(&service.Id, &service.Name, &service.Description, &service.Version)
	if err != nil {
		return Service{}, fmt.Errorf("failed to scan row while fetching service by id: %w", err)
	}
	return service, nil
}
