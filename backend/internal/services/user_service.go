package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"interview-prep-app/internal/models"
	"interview-prep-app/internal/repositories"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// RegisterWithEmail registers a new user with email and password
func (s *UserService) RegisterWithEmail(req *models.CreateUserRequest) (*models.User, error) {
	// Check if email already exists
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		Name:         req.Name,
		Avatar:       req.Avatar,
		AuthProvider: models.AuthProviderEmail,
		ProviderID:   "", // Empty string for email users - will be handled as NULL in DB
		PasswordHash: hashedPassword,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Remove password hash from returned user
	user.PasswordHash = ""
	return user, nil
}

// LoginWithEmail authenticates a user with email and password
func (s *UserService) LoginWithEmail(email, password string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Update last login
	err = s.userRepo.UpdateLastLogin(user.ID)
	if err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login: %v\n", err)
	}

	// Remove password hash from returned user
	user.PasswordHash = ""
	return user, nil
}

// LoginWithOAuth authenticates or registers a user with OAuth
func (s *UserService) LoginWithOAuth(req *models.OAuthLoginRequest) (*models.User, error) {
	// Validate OAuth token and get user info
	userInfo, err := s.validateOAuthToken(req)
	if err != nil {
		return nil, fmt.Errorf("invalid OAuth token: %w", err)
	}

	// Try to find existing user by provider ID
	user, err := s.userRepo.GetByProviderID(req.Provider, userInfo.ProviderID)
	if err == nil {
		// User exists, update last login
		err = s.userRepo.UpdateLastLogin(user.ID)
		if err != nil {
			fmt.Printf("Failed to update last login: %v\n", err)
		}
		return user, nil
	}

	// Try to find existing user by email
	user, err = s.userRepo.GetByEmail(userInfo.Email)
	if err == nil {
		// User exists with different provider, link accounts
		// For now, we'll return an error to prevent account linking without explicit consent
		return nil, fmt.Errorf("email already exists with different provider")
	}

	// Create new user
	user = &models.User{
		Email:        userInfo.Email,
		Name:         userInfo.Name,
		Avatar:       userInfo.Avatar,
		AuthProvider: req.Provider,
		ProviderID:   userInfo.ProviderID,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id int) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Remove password hash from returned user
	user.PasswordHash = ""
	return user, nil
}

// UpdateUser updates a user's profile
func (s *UserService) UpdateUser(userID int, req *models.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	// Remove password hash from returned user
	user.PasswordHash = ""
	return user, nil
}

// GenerateRefreshToken generates a new refresh token
func (s *UserService) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// CreateRefreshToken creates and stores a refresh token
func (s *UserService) CreateRefreshToken(userID int) (string, error) {
	token, err := s.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days
	err = s.userRepo.CreateRefreshToken(userID, token, expiresAt)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateRefreshToken validates a refresh token
func (s *UserService) ValidateRefreshToken(token string) (*models.User, error) {
	refreshToken, err := s.userRepo.GetRefreshToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if refreshToken.IsRevoked {
		return nil, fmt.Errorf("refresh token revoked")
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	user, err := s.userRepo.GetByID(refreshToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Remove password hash from returned user
	user.PasswordHash = ""
	return user, nil
}

// RevokeRefreshToken revokes a refresh token
func (s *UserService) RevokeRefreshToken(token string) error {
	return s.userRepo.RevokeRefreshToken(token)
}

// hashPassword hashes a password using bcrypt
func (s *UserService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// OAuthUserInfo represents user info from OAuth provider
type OAuthUserInfo struct {
	ProviderID string
	Email      string
	Name       string
	Avatar     string
}

// validateOAuthToken validates OAuth token and returns user info
func (s *UserService) validateOAuthToken(req *models.OAuthLoginRequest) (*OAuthUserInfo, error) {
	switch req.Provider {
	case models.AuthProviderGoogle:
		return s.validateGoogleToken(req.AccessToken)
	case models.AuthProviderFacebook:
		return s.validateFacebookToken(req.AccessToken)
	case models.AuthProviderApple:
		return s.validateAppleToken(req.AccessToken)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}
}

// validateGoogleToken validates Google OAuth token
func (s *UserService) validateGoogleToken(token string) (*OAuthUserInfo, error) {
	// Google OAuth token validation
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to validate Google token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid Google token")
	}

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	err = json.NewDecoder(resp.Body).Decode(&googleUser)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Google user info: %w", err)
	}

	return &OAuthUserInfo{
		ProviderID: googleUser.ID,
		Email:      googleUser.Email,
		Name:       googleUser.Name,
		Avatar:     googleUser.Picture,
	}, nil
}

// validateFacebookToken validates Facebook OAuth token
func (s *UserService) validateFacebookToken(token string) (*OAuthUserInfo, error) {
	// Facebook OAuth token validation
	url := fmt.Sprintf("https://graph.facebook.com/me?fields=id,email,name,picture&access_token=%s", token)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to validate Facebook token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid Facebook token")
	}

	var facebookUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}

	err = json.NewDecoder(resp.Body).Decode(&facebookUser)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Facebook user info: %w", err)
	}

	return &OAuthUserInfo{
		ProviderID: facebookUser.ID,
		Email:      facebookUser.Email,
		Name:       facebookUser.Name,
		Avatar:     facebookUser.Picture.Data.URL,
	}, nil
}

// validateAppleToken validates Apple OAuth token
func (s *UserService) validateAppleToken(token string) (*OAuthUserInfo, error) {
	// Apple OAuth token validation is more complex and requires JWT verification
	// For now, we'll implement a basic validation
	// In production, you should use Apple's JWT verification

	if token == "" {
		return nil, fmt.Errorf("empty Apple token")
	}

	// Parse JWT token to extract user info
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid Apple token format")
	}

	// In a real implementation, you would:
	// 1. Verify the JWT signature using Apple's public keys
	// 2. Validate the token claims (iss, aud, exp, etc.)
	// 3. Extract user information from the token

	// For now, we'll return a placeholder implementation
	// This should be replaced with proper Apple JWT verification
	return &OAuthUserInfo{
		ProviderID: "apple_user_id",    // This should come from the JWT sub claim
		Email:      "user@example.com", // This should come from the JWT email claim
		Name:       "Apple User",       // This might not be available in Apple tokens
		Avatar:     "",                 // Apple doesn't provide avatar URLs
	}, fmt.Errorf("Apple OAuth not fully implemented - please implement JWT verification")
}

// CleanupExpiredTokens removes expired refresh tokens
func (s *UserService) CleanupExpiredTokens() error {
	return s.userRepo.CleanupExpiredRefreshTokens()
}
