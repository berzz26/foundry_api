package profile

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	UserID                  uuid.UUID
	UserType                string
	ExperienceLevel         string
	Goals                   []string
	InterestedRoles         []string
	PreferredLocations      []string
	CompanyStagePreferences []string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}
