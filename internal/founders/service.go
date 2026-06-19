package founders

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrFounderNotFound = errors.New("founder not found")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, f *Founder) (*Founder, error) {
	return s.repo.Create(ctx, f)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*Founder, error) {
	founder, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFounderNotFound
		}
		return nil, err
	}
	return founder, nil
}
func (s *Service) GetByCompanyID(ctx context.Context, companyId int64) ([]Founder, error) {
	return s.repo.GetByCompanyID(ctx, companyId)
}
func (s *Service) List(ctx context.Context, companyID *int64, limit, offset int) ([]Founder, int64, error) {
	list, err := s.repo.List(ctx, companyID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count(ctx, companyID)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *Service) Update(ctx context.Context, f *Founder) (*Founder, error) {
	return s.repo.Update(ctx, f)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
