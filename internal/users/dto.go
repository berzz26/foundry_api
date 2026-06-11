package users

import "time"

type AddUserDTO struct {
	ID              string
	Email           string `json:"email" validate:"required,email"`
	FirstName       string `json:"firstName" validate:"required"`
	LastName        string `json:"lastName" validate:"required"`
	PasswordHash    string `json:"passwordHash" validate:"required"`
	ProfileImageURL string `json:"profileImageUrl" validate:"omitempty,url"`
	Provider        string `json:"provider" validate:"required"`
	ProviderID      string `json:"providerId" validate:"required"`
}

type UpdateUserDTO struct {
	Email           string `json:"email" validate:"omitempty,email"`
	FirstName       string `json:"firstName" validate:"omitempty"`
	LastName        string `json:"lastName" validate:"omitempty"`
	PasswordHash    string `json:"passwordHash" validate:"omitempty"`
	ProfileImageURL string `json:"profileImageUrl" validate:"omitempty,url"`
	Provider        string `json:"provider" validate:"omitempty"`
	ProviderID      string `json:"providerId" validate:"omitempty"`
}

type ResponseDTO struct {
	ID              string
	Email           string `json:"email" validate:"required,email"`
	FirstName       string `json:"firstName" validate:"required"`
	LastName        string `json:"lastName" validate:"required"`
	ProfileImageURL string `json:"profileImageUrl" validate:"omitempty,url"`
	Provider        string `json:"provider" validate:"required"`
	ProviderID      string `json:"providerId" validate:"required"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
