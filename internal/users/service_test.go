package users

import (
	"testing"
)

func TestVerifyPassword(t *testing.T) {
	s := NewService(nil)
	password := "supersecret123"

	hashed, err := s.hashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error when hashing password, got %v", err)
	}

	if !s.VerifyPassword(hashed, password) {
		t.Errorf("Expected password verification to succeed")
	}

	if s.VerifyPassword(hashed, "wrongpassword") {
		t.Errorf("Expected password verification to fail for wrong password")
	}
}

func TestDTOValidation(t *testing.T) {
	t.Run("AddUserDTO validation", func(t *testing.T) {
		dto := AddUserDTO{
			Email:     "invalid-email",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "short",
			Provider:  "local",
		}

		err := validate.Struct(dto)
		if err == nil {
			t.Error("Expected validation error for invalid email and short password, got nil")
		}
	})

	t.Run("AddUserDTO validation valid", func(t *testing.T) {
		dto := AddUserDTO{
			Email:     "john.doe@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "strongpassword123",
			Provider:  "local",
		}

		err := validate.Struct(dto)
		if err != nil {
			t.Errorf("Expected no validation error, got %v", err)
		}
	})
}
