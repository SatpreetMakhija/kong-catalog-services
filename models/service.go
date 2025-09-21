package models

import (
	"context"

	"kong.com/catalog/datastore"
)

type Service struct{
	Repo datastore.ServiceRepo	
}

func (m *Service) FetchServiceById(ctx context.Context, id string) (datastore.Service, error) {
	service, err := m.Repo.FindServiceById(ctx, id)
	return service, err
}

func NewService(repo datastore.ServiceRepo) *Service {
	return &Service{Repo: repo}
}
