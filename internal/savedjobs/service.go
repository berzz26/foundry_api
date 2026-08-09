package savedjobs

import (
	"context"

	"github.com/berzz26/foundry_api/internal/jobs"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Save(ctx context.Context, userID string, jobID int64) (*SavedJobResponse, error) {
	sj, err := s.repo.Save(ctx, userID, jobID)
	if err != nil {
		return nil, err
	}
	return &SavedJobResponse{
		JobID:   sj.JobID,
		SavedAt: sj.CreatedAt,
	}, nil
}

func (s *Service) Delete(ctx context.Context, userID string, jobID int64) error {
	return s.repo.Delete(ctx, userID, jobID)
}

func (s *Service) IsSaved(ctx context.Context, userID string, jobID int64) (bool, error) {
	return s.repo.IsSaved(ctx, userID, jobID)
}

func (s *Service) ListByUser(ctx context.Context, userID string) ([]SavedJobCard, error) {
	items, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	cards := make([]SavedJobCard, 0, len(items))
	for _, item := range items {
		card := jobs.ToJobCardResponse(item.Job)
		cards = append(cards, SavedJobCard{
			JobCardResponse: card,
			SavedAt:         item.CreatedAt,
		})
	}
	return cards, nil
}
