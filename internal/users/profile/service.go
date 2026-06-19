package profile

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

func (s *Service) CompleteOnboarding(
	ctx context.Context,
	userID uuid.UUID,
	req CompleteOnboardingRequest,
) error {
	// 1. Validate payload
	if err := ValidateOnboardingRequest(req); err != nil {
		return err
	}

	// 2. Map request DTO to model
	p := &UserProfile{
		UserID:                  userID,
		UserType:                req.UserType,
		ExperienceLevel:         req.ExperienceLevel,
		Goals:                   req.Goals,
		InterestedRoles:         req.InterestedRoles,
		PreferredLocations:      req.PreferredLocations,
		CompanyStagePreferences: req.CompanyStagePreferences,
	}

	// 3. Execute transactional write
	return s.repo.CompleteOnboardingTx(ctx, userID, p)
}

func (s *Service) GetProfile(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
	return s.repo.GetProfile(ctx, userID)
}
