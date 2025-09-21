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

func (m *Service) Search(ctx context.Context, request datastore.ServiceSearchRequest) ( datastore.ServiceSearchResponse, error) {
	searchResponse, err := m.Repo.SearchService(ctx, request)
	return searchResponse, err
}

func NewService(repo datastore.ServiceRepo) *Service {
	return &Service{Repo: repo}
}
