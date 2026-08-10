package savedjobs

import (
	"context"
	"sort"

	"github.com/berzz26/foundry_api/internal/companies"
	"github.com/berzz26/foundry_api/internal/jobs"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SaveJob(ctx context.Context, userID string, jobID int64) (*SaveResponse, error) {
	si, err := s.repo.SaveJob(ctx, userID, jobID)
	if err != nil {
		return nil, err
	}
	return &SaveResponse{
		JobID:     si.JobID,
		CompanyID: &si.CompanyID,
		SavedAt:   si.CreatedAt,
	}, nil
}

func (s *Service) SaveCompany(ctx context.Context, userID string, companyID int64) (*SaveResponse, error) {
	si, err := s.repo.SaveCompany(ctx, userID, companyID)
	if err != nil {
		return nil, err
	}
	return &SaveResponse{
		CompanyID: &si.CompanyID,
		SavedAt:   si.CreatedAt,
	}, nil
}

func (s *Service) DeleteJob(ctx context.Context, userID string, jobID int64) error {
	return s.repo.DeleteJob(ctx, userID, jobID)
}

func (s *Service) DeleteCompany(ctx context.Context, userID string, companyID int64) error {
	return s.repo.DeleteCompany(ctx, userID, companyID)
}

func (s *Service) IsJobSaved(ctx context.Context, userID string, jobID int64) (bool, error) {
	return s.repo.IsJobSaved(ctx, userID, jobID)
}

func (s *Service) IsCompanySaved(ctx context.Context, userID string, companyID int64) (bool, error) {
	return s.repo.IsCompanySaved(ctx, userID, companyID)
}

func (s *Service) ListByUser(ctx context.Context, userID string) ([]SavedItemCard, error) {
	items, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	cards := make([]SavedItemCard, 0, len(items))
	for _, item := range items {
		if item.Job != nil {
			card := jobs.ToJobCardResponse(*item.Job)
			cards = append(cards, SavedItemCard{
				Type:    "job",
				Job:     &card,
				SavedAt: item.CreatedAt,
			})
		} else if item.Company != nil {
			card := companies.ToCardResponse(item.Company)
			cards = append(cards, SavedItemCard{
				Type:    "company",
				Company: &card,
				SavedAt: item.CreatedAt,
			})
		}
	}

	sort.SliceStable(cards, func(i, j int) bool {
		return cards[i].SavedAt.After(cards[j].SavedAt)
	})

	return cards, nil
}
