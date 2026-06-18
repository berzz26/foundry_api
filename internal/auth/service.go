package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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
	authRepo    *Repository
}

func NewService(userService *users.Service, authRepo *Repository) *Service {
	return &Service{userService: userService, authRepo: authRepo}
}

func (s *Service) GenerateToken(user *users.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s *Service) GenerateRefreshToken(ctx context.Context, userID string) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	//never store the raw refresh token in the db
	refreshToken := hex.EncodeToString(b)
	
	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])
	//refresh tokens expire after a week, user logs out
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	err = s.authRepo.CreateRefreshToken(ctx, userID, tokenHash, expiresAt)
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func (s *Service) RefreshSession(ctx context.Context, refreshToken string) (*users.User, string, string, error) {
	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])
	
	rt, err := s.authRepo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, "", "", errors.New("invalid refresh token")
	}

	if rt.RevokedAt != nil {
		return nil, "", "", errors.New("refresh token has been revoked")
	}
	
	_ = s.authRepo.DeleteRefreshToken(ctx, rt.ID)
	
	if time.Now().After(rt.ExpiresAt) {
		return nil, "", "", errors.New("refresh token expired")
	}
	
	user, err := s.userService.GetByID(ctx, rt.UserID)
	if err != nil {
		return nil, "", "", errors.New("user not found")
	}
	
	newAccessToken, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", "", err
	}
	
	newRefreshToken, err := s.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, "", "", err
	}
	
	return user, newAccessToken, newRefreshToken, nil
}

func (s *Service) RevokeSession(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])
	
	return s.authRepo.RevokeRefreshToken(ctx, tokenHash)
}

func (s *Service) Register(ctx context.Context, signupDTO *SignupDTO) (*users.User, string, string, error) {
	user := &users.User{
		Email:     signupDTO.Email,
		FirstName: signupDTO.FirstName,
		LastName:  signupDTO.LastName,
		PasswordHash: signupDTO.Password,
		Provider:  "local",
		Role:      "user",
	}

	createdUser, err := s.userService.AddUser(ctx, user)
	if err != nil {
		return nil, "", "", err
	}

	token, err := s.GenerateToken(createdUser)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := s.GenerateRefreshToken(ctx, createdUser.ID)
	if err != nil {
		return nil, "", "", err
	}

	return createdUser, token, refreshToken, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*users.User, string, string, error) {
	user, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", "", ErrInvalidCredentials
	}

	if user.Provider != "local" && user.Provider != "credentials" && user.Provider != "" {
		return nil, "", "", fmt.Errorf("this account uses %s authentication. Please sign in via OAuth.", user.Provider)
	}

	if !s.userService.VerifyPassword(user.PasswordHash, password) {
		return nil, "", "", ErrInvalidCredentials
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := s.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, "", "", err
	}

	return user, token, refreshToken, nil
}

func (s *Service) GetOrCreateOAuthUser(ctx context.Context, provider, providerID, email, firstName, lastName, avatarURL string) (*users.User, string, string, error) {
	existingUser, err := s.userService.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		// User already exists. Verify provider matches or link it
		if existingUser.Provider != provider && existingUser.ProviderID != nil && *existingUser.ProviderID != providerID {
			return nil, "", "", ErrOAuthProviderMismatch
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
				return nil, "", "", err
			}
		}

		token, err := s.GenerateToken(existingUser)
		if err != nil {
			return nil, "", "", err
		}
		refreshToken, err := s.GenerateRefreshToken(ctx, existingUser.ID)
		if err != nil {
			return nil, "", "", err
		}
		return existingUser, token, refreshToken, nil
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
		return nil, "", "", err
	}

	token, err := s.GenerateToken(createdUser)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := s.GenerateRefreshToken(ctx, createdUser.ID)
	if err != nil {
		return nil, "", "", err
	}

	return createdUser, token, refreshToken, nil
}

func (s *Service) GetProfile(ctx context.Context, userID string) (*users.User, error) {
	return s.userService.GetByID(ctx, userID)
}
