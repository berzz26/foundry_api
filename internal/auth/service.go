package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/berzz26/foundry_api/internal/users"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrOAuthProviderMismatch = errors.New("account linked with a different provider")
)

var jwtSecret []byte

func init() {
	secretStr := os.Getenv("JWT_SECRET")
	if secretStr == "" {
		log.Println("WARNING: JWT_SECRET environment variable is empty. Generating an ephemeral secret key. This will not scale horizontally!")
		key := make([]byte, 32)
		_, err := rand.Read(key)
		if err != nil {
			panic("failed to generate random JWT secret: " + err.Error())
		}
		jwtSecret = key
	} else {
		jwtSecret = []byte(secretStr)
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type Service struct {
	userService *users.Service
}

func NewService(userService *users.Service) *Service {
	return &Service{userService: userService}
}

func (s *Service) GenerateToken(user *users.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s *Service) Register(ctx context.Context, signupDTO *SignupDTO) (*users.User, error) {
	user := &users.User{
		Email:     signupDTO.Email,
		FirstName: signupDTO.FirstName,
		LastName:  signupDTO.LastName,
		PasswordHash: signupDTO.Password,
		Provider:  "local",
		Role:      "user",
	}

	return s.userService.AddUser(ctx, user)
}

func (s *Service) Login(ctx context.Context, email, password string) (*users.User, string, error) {
	user, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if user.Provider != "local" && user.Provider != "credentials" && user.Provider != "" {
		return nil, "", fmt.Errorf("this account uses %s authentication. Please sign in via OAuth.", user.Provider)
	}

	if !s.userService.VerifyPassword(user.PasswordHash, password) {
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *Service) GetOrCreateOAuthUser(ctx context.Context, provider, providerID, email, firstName, lastName, avatarURL string) (*users.User, string, error) {
	existingUser, err := s.userService.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		// User already exists. Verify provider matches or link it
		if existingUser.Provider != provider && existingUser.ProviderID != nil && *existingUser.ProviderID != providerID {
			return nil, "", ErrOAuthProviderMismatch
		}

		// Update provider details if they were local or empty
		needsUpdate := false
		if existingUser.Provider != provider {
			existingUser.Provider = provider
			needsUpdate = true
		}
		if existingUser.ProviderID == nil || *existingUser.ProviderID != providerID {
			existingUser.ProviderID = &providerID
			needsUpdate = true
		}
		if avatarURL != "" && (existingUser.ProfileImageURL == nil || *existingUser.ProfileImageURL == "") {
			existingUser.ProfileImageURL = &avatarURL
			needsUpdate = true
		}

		if needsUpdate {
			existingUser, err = s.userService.Update(ctx, existingUser)
			if err != nil {
				return nil, "", err
			}
		}

		token, err := s.GenerateToken(existingUser)
		if err != nil {
			return nil, "", err
		}
		return existingUser, token, nil
	}

	// Create new OAuth user
	var profileImage *string
	if avatarURL != "" {
		profileImage = &avatarURL
	}

	newUser := &users.User{
		Email:           email,
		FirstName:       firstName,
		LastName:        lastName,
		PasswordHash:    "", // Empty password for OAuth
		ProfileImageURL: profileImage,
		Provider:        provider,
		ProviderID:      &providerID,
		Role:            "user",
	}

	createdUser, err := s.userService.AddUser(ctx, newUser)
	if err != nil {
		return nil, "", err
	}

	token, err := s.GenerateToken(createdUser)
	if err != nil {
		return nil, "", err
	}

	return createdUser, token, nil
}

func (s *Service) GetProfile(ctx context.Context, userID string) (*users.User, error) {
	return s.userService.GetByID(ctx, userID)
}
