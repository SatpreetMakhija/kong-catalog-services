package datastore

import (
	"context"
	"fmt"

)

type ServiceRepo interface {
	FindServiceById(context.Context, string) (Service, error)
	SearchService(context.Context, ServiceSearchRequest) (ServiceSearchResponse, error)
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

func (r *PostgresServiceRepo) SearchService(ctx context.Context, searchRequest ServiceSearchRequest) (ServiceSearchResponse, error) {

	query := `
with base as ( 
	select id, name, description, version from services
	where 
	($1::text is null or name ilike $1)
	and ($2::text is null or version = $2)
)
	select id, name, description, version from base
	limit $3 offset $4
	`
	offset := (searchRequest.Page - 1) * searchRequest.PageSize

	var versionParameter any
	var nameParameter any

	if searchRequest.Version == "" {
		versionParameter = nil
	} else {
		versionParameter = searchRequest.Version
	}

	if searchRequest.Name == "" {
		nameParameter = nil 
	} else {
		nameParameter = "%" + searchRequest.Name + "%"
	}

	rows, err := r.ds.Client.Query(ctx, query, nameParameter, versionParameter, &searchRequest.PageSize, &offset)

	if err != nil {
		return ServiceSearchResponse{}, err
	}

	services := make([]Service, 0, searchRequest.PageSize)

	for rows.Next() {
		var s Service
		err := rows.Scan(&s.Id, &s.Name, &s.Description, &s.Version)
		if err != nil {
			return ServiceSearchResponse{}, fmt.Errorf("failed to scan result: %w", err)
		}
		services = append(services, s)
	}
	return ServiceSearchResponse{Items: services, Page: searchRequest.Page, PageSize: searchRequest.PageSize}, nil
}
