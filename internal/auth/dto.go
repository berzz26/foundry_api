package auth

import "github.com/berzz26/foundry_api/internal/users"

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignupDTO struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Password  string `json:"password" validate:"required,min=8"`
}

type AuthResponseDTO struct {
	Token        string            `json:"token"`
	RefreshToken string            `json:"refreshToken"`
	User         users.ResponseDTO `json:"user"`
}

type RefreshTokenRequestDTO struct {
	RefreshToken string `json:"refreshToken"`
}
