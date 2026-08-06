package github

import (
	"context"
	"fmt"
)

type ClientAPI interface {
	FetchUserActivity(ctx context.Context, username string) ([]Event, error)
	FetchUserRepositories(ctx context.Context, username string) ([]Repository, error)
}

type Service struct {
	client ClientAPI
}

func NewService(client ClientAPI) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) GetUserActivity(ctx context.Context, username string) ([]Event, error) {
	if username == "" {
		return nil, fmt.Errorf("username tidak boleh kosong")
	}

	events, err := s.client.FetchUserActivity(ctx, username)
	if err != nil {
		return nil, err
	}

	
	return events, nil
}
