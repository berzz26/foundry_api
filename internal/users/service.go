package users

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddUser(ctx context.Context, user *User) (*User, error) {
	// Check if user with email already exists
	existing, err := s.repo.GetByEmail(ctx, user.Email)
	if err == nil && existing != nil {
		return nil, ErrUserAlreadyExists
	} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	user.ID = uuid.New().String()
	hashed, err := s.hashPassword(user.PasswordHash)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = hashed
	return s.repo.AddUser(ctx, user)
}

func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) List(ctx context.Context, limit int, offset int) ([]User, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *Service) Update(ctx context.Context, user *User) (*User, error) {
	return s.repo.Update(ctx, user)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *Service) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}