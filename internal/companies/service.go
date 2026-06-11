package companies

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrCompanyNotFound = errors.New("company not found")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(ctx context.Context, id int64) (*Company, error) {
	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCompanyNotFound
		}
		return nil, err
	}
	return company, nil
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*Company, error) {
	company, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCompanyNotFound
		}
		return nil, err
	}
	return company, nil
}

func (s *Service) List(ctx context.Context, filters CompanyFilters) ([]Company, error) {
	return s.repo.List(ctx, filters)
}
