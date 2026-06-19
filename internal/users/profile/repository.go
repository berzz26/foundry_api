package profile

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// CompleteOnboardingTx runs database queries in a transaction
func (r *Repository) CompleteOnboardingTx(ctx context.Context, userID uuid.UUID, profile *UserProfile) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Create or Update user_profile (Upsert)
	queryProfile := `
		INSERT INTO user_profiles (
			user_id, user_type, experience_level, goals, interested_roles, preferred_locations, company_stage_preferences, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			user_type = EXCLUDED.user_type,
			experience_level = EXCLUDED.experience_level,
			goals = EXCLUDED.goals,
			interested_roles = EXCLUDED.interested_roles,
			preferred_locations = EXCLUDED.preferred_locations,
			company_stage_preferences = EXCLUDED.company_stage_preferences,
			updated_at = NOW()
	`
	_, err = tx.Exec(ctx, queryProfile,
		profile.UserID,
		profile.UserType,
		profile.ExperienceLevel,
		profile.Goals,
		profile.InterestedRoles,
		profile.PreferredLocations,
		profile.CompanyStagePreferences,
	)
	if err != nil {
		return err
	}

	// 2. Set users.onboarding_completed = true
	queryUser := `
		UPDATE users
		SET onboarding_completed = TRUE, updated_at = NOW()
		WHERE id = $1
	`
	res, err := tx.Exec(ctx, queryUser, userID.String())
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return tx.Commit(ctx)
}

func (r *Repository) GetProfile(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
	query := `
		SELECT user_id, user_type, experience_level, goals, interested_roles, preferred_locations, company_stage_preferences, created_at, updated_at
		FROM user_profiles
		WHERE user_id = $1
	`
	var p UserProfile
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&p.UserID,
		&p.UserType,
		&p.ExperienceLevel,
		&p.Goals,
		&p.InterestedRoles,
		&p.PreferredLocations,
		&p.CompanyStagePreferences,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}
	return &p, nil
}
