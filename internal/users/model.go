package users

import "time"

type User struct {
	ID                  string
	Email               string
	FirstName           string
	LastName            string
	PasswordHash        string
	ProfileImageURL     *string
	Provider            string
	ProviderID          *string
	Role                string
	OnboardingCompleted bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
