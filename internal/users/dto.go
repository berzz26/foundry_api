package users

import "time"

type AddUserDTO struct {
	Email           string  `json:"email" validate:"required,email"`
	FirstName       string  `json:"firstName" validate:"required"`
	LastName        string  `json:"lastName" validate:"required"`
	Password        string  `json:"password" validate:"required,min=8"`
	ProfileImageURL *string `json:"profileImageUrl" validate:"omitempty,url"`
	Provider        string  `json:"provider" validate:"required"` // e.g., "local", "google", "github"
	ProviderID      *string `json:"providerId" validate:"omitempty"`
}

type UpdateUserDTO struct {
	Email           *string `json:"email" validate:"omitempty,email"`
	FirstName       *string `json:"firstName" validate:"omitempty"`
	LastName        *string `json:"lastName" validate:"omitempty"`
	Password        *string `json:"password" validate:"omitempty,min=8"`
	ProfileImageURL *string `json:"profileImageUrl" validate:"omitempty,url"`
	Provider        *string `json:"provider" validate:"omitempty"`
	ProviderID      *string `json:"providerId" validate:"omitempty"`
}

type ResponseDTO struct {
	ID              string    `json:"id"`
	Email           string    `json:"email"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	ProfileImageURL *string   `json:"profileImageUrl,omitempty"`
	Provider        string    `json:"provider"`
	ProviderID      *string   `json:"providerId,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
