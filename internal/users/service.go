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
func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *Service) GetByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.GetByEmail(ctx, email)
}
func (s *Service) List(ctx context.Context, limit int, offset int) ([]User, error) {
	return s.repo.List(ctx, limit, offset)
}
func (s *Service) Update(ctx context.Context, user *User) (*User, error) {
	return s.repo.Update(ctx, user)
}
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}