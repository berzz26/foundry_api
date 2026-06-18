package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/berzz26/foundry_api/internal/users"
	"github.com/berzz26/foundry_api/pkg/config"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type Handler struct {
	service *Service
	cfg     *config.Config
}

func NewHandler(service *Service, cfg *config.Config) *Handler {
	return &Handler{service: service, cfg: cfg}
}

func (h *Handler) Signup(c *fiber.Ctx) error {
	dto := new(SignupDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := validate.Struct(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	createdUser, err := h.service.Register(ctx, dto)
	if err != nil {
		if errors.Is(err, users.ErrUserAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "User already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(mapToResponseDTO(createdUser))
}

func (h *Handler) Login(c *fiber.Ctx) error {
	dto := new(LoginDTO)
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := validate.Struct(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	user, token, err := h.service.Login(ctx, dto.Email, dto.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	h.setAuthCookie(c, token)

	return c.JSON(AuthResponseDTO{
		Token: token,
		User:  mapToResponseDTO(user),
	})
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	c.ClearCookie("__Secure-token", "token")
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GoogleLogin(c *fiber.Ctx) error {
	if h.cfg.GoogleClientID == "" || h.cfg.GoogleSecret == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Google OAuth is not configured on this server.",
		})
	}

	redirectURI := fmt.Sprintf("%s/api/v1/auth/google/callback", h.cfg.HostURL)
	googleAuthURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email+profile&state=google-auth-state",
		url.QueryEscape(h.cfg.GoogleClientID),
		url.QueryEscape(redirectURI),
	)

	return c.Redirect(googleAuthURL, fiber.StatusTemporaryRedirect)
}

func (h *Handler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Authorization code not provided by Google",
		})
	}

	redirectTarget := os.Getenv("REDIRECT_URL")
	if redirectTarget == "" {
		redirectTarget = "http://localhost:3001"
	}

	redirectURI := fmt.Sprintf("%s/api/v1/auth/google/callback", h.cfg.HostURL)

	// Exchange code for token
	resp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"code":          {code},
		"client_id":     {h.cfg.GoogleClientID},
		"client_secret": {h.cfg.GoogleSecret},
		"redirect_uri":  {redirectURI},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to exchange authorization code with Google",
		})
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse Google token response",
		})
	}

	// Fetch user info
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	respInfo, err := http.DefaultClient.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user profile from Google",
		})
	}
	defer respInfo.Body.Close()

	var userInfo struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		GivenName string `json:"given_name"`
		FamilyName string `json:"family_name"`
		Picture   string `json:"picture"`
	}
	if err := json.NewDecoder(respInfo.Body).Decode(&userInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse Google user profile",
		})
	}

	if userInfo.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Google profile does not provide an email address",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	_, token, err := h.service.GetOrCreateOAuthUser(
		ctx,
		"google",
		userInfo.ID,
		userInfo.Email,
		userInfo.GivenName,
		userInfo.FamilyName,
		userInfo.Picture,
	)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	h.setAuthCookie(c, token)

	return c.Redirect(redirectTarget, fiber.StatusTemporaryRedirect)
}

func (h *Handler) GithubLogin(c *fiber.Ctx) error {
	if h.cfg.GithubClientID == "" || h.cfg.GithubSecret == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "GitHub OAuth is not configured on this server.",
		})
	}

	redirectURI := fmt.Sprintf("%s/api/v1/auth/github/callback", h.cfg.HostURL)
	githubAuthURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user+user:email&state=github-auth-state",
		url.QueryEscape(h.cfg.GithubClientID),
		url.QueryEscape(redirectURI),
	)

	return c.Redirect(githubAuthURL, fiber.StatusTemporaryRedirect)
}

func (h *Handler) GithubCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Authorization code not provided by GitHub",
		})
	}

	redirectTarget := os.Getenv("REDIRECT_URL")
	if redirectTarget == "" {
		redirectTarget = "http://localhost:3001"
	}

	redirectURI := fmt.Sprintf("%s/api/v1/auth/github/callback", h.cfg.HostURL)

	// Exchange code for token
	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(url.Values{
		"code":          {code},
		"client_id":     {h.cfg.GithubClientID},
		"client_secret": {h.cfg.GithubSecret},
		"redirect_uri":  {redirectURI},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to exchange authorization code with GitHub",
		})
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse GitHub token response",
		})
	}

	// Fetch user profile
	reqInfo, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	reqInfo.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	reqInfo.Header.Set("Accept", "application/json")
	respInfo, err := http.DefaultClient.Do(reqInfo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user profile from GitHub",
		})
	}
	defer respInfo.Body.Close()

	var githubInfo struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(respInfo.Body).Decode(&githubInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse GitHub user profile",
		})
	}

	// If email is empty/private, fetch emails list
	if githubInfo.Email == "" {
		reqEmails, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
		reqEmails.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
		reqEmails.Header.Set("Accept", "application/json")
		respEmails, err := http.DefaultClient.Do(reqEmails)
		if err == nil {
			defer respEmails.Body.Close()
			var emails []struct {
				Email    string `json:"email"`
				Primary  bool   `json:"primary"`
				Verified bool   `json:"verified"`
			}
			if err := json.NewDecoder(respEmails.Body).Decode(&emails); err == nil {
				for _, e := range emails {
					if e.Primary {
						githubInfo.Email = e.Email
						break
					}
				}
				if githubInfo.Email == "" && len(emails) > 0 {
					githubInfo.Email = emails[0].Email
				}
			}
		}
	}

	if githubInfo.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "GitHub profile does not provide an email address",
		})
	}

	// Parse first & last name from GitHub Name
	firstName := githubInfo.Login
	lastName := ""
	if githubInfo.Name != "" {
		parts := strings.SplitN(githubInfo.Name, " ", 2)
		firstName = parts[0]
		if len(parts) > 1 {
			lastName = parts[1]
		}
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	_, token, err := h.service.GetOrCreateOAuthUser(
		ctx,
		"github",
		strconv.Itoa(githubInfo.ID),
		githubInfo.Email,
		firstName,
		lastName,
		githubInfo.AvatarURL,
	)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	h.setAuthCookie(c, token)

	return c.Redirect(redirectTarget, fiber.StatusTemporaryRedirect)
}

func (h *Handler) setAuthCookie(c *fiber.Ctx, token string) {
	secure := false
	if os.Getenv("APP_ENV") == "production" {
		secure = true
	}

	c.Cookie(&fiber.Cookie{
		Name:     "__Secure-token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Lax",
		Path:     "/",
	})
	
	// Local dev fallback cookie (non __Secure prefixed name)
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})
}

func mapToResponseDTO(u *users.User) users.ResponseDTO {
	return users.ResponseDTO{
		ID:              u.ID,
		Email:           u.Email,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		ProfileImageURL: u.ProfileImageURL,
		Provider:        u.Provider,
		ProviderID:      u.ProviderID,
		Role:            u.Role,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

func (h *Handler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized. Please sign in.",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	user, err := h.service.GetProfile(ctx, userID.(string))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(mapToResponseDTO(user))
}
