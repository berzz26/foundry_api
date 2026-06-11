package users

import (
	"context"
	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddUser(ctx context.Context, user *User) (*User, error) {
	user.ID = uuid.New().String()
	return s.repo.AddUser(ctx, user)
}

