package users

import "time"

type User struct {
	ID              string
	Email           string
	FirstName       string
	LastName        string
	PasswordHash    string
	ProfileImageURL *string
	Provider        string
	ProviderID      *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
