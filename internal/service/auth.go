package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"

	"github.com/momokii/go-cli-notes/internal/model"
	"github.com/momokii/go-cli-notes/internal/repository"
	"github.com/momokii/go-cli-notes/internal/util"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo       repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	hasher         *util.PasswordHasher
	jwtManager     *util.JWTManager
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	hasher *util.PasswordHasher,
	jwtManager *util.JWTManager,
) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		refreshTokenRepo: refreshTokenRepo,
		hasher:         hasher,
		jwtManager:     jwtManager,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error) {
	// Validate request
	if err := util.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("%w: %s", model.ErrValidation, err)
	}

	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("check email exists: %w", err)
	}
	if exists {
		return nil, model.ErrAPIEmailExists
	}

	// Check if username already exists
	exists, err = s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("check username exists: %w", err)
	}
	if exists {
		return nil, model.ErrAPIUsernameExists
	}

	// Hash password
	hash, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	// Create user
	user := &model.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hash,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// Log activity
	// TODO: Log registration activity

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	// Validate request
	if err := util.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("%w: %s", model.ErrValidation, err)
	}

	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, model.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("find user: %w", err)
	}

	// Verify password
	valid, err := s.hasher.Verify(req.Password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("verify password: %w", err)
	}
	if !valid {
		return nil, model.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, model.ErrUnauthorized
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	// Store refresh token hash
	_ = s.hashToken(refreshToken)
	// TODO: Store refresh token in database

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log but don't fail the request
		fmt.Printf("warning: update last login failed: %v\n", err)
	}

	// Log activity
	// TODO: Log login activity

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.jwtManager.GetAccessExpiration(),
		User:         user,
	}, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthResponse, error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", model.ErrInvalidToken, err)
	}

	// Check if refresh token is revoked
	_ = s.hashToken(refreshToken)
	// TODO: Check if token is revoked in database

	// Get user
	user, err := s.userRepo.FindByID(ctx, uuid.MustParse(claims.UserID))
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, model.ErrUnauthorized
	}

	// Generate new tokens
	newAccessToken, err := s.jwtManager.GenerateAccessToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	// TODO: Revoke old refresh token and store new one

	return &model.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    s.jwtManager.GetAccessExpiration(),
		User:         user,
	}, nil
}

// Logout revokes a refresh token
func (s *AuthService) Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	_ = s.hashToken(refreshToken)
	// TODO: Revoke refresh token in database

	// Log activity
	// TODO: Log logout activity

	return nil
}

// hashToken creates a SHA256 hash of a token for storage
func (s *AuthService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
