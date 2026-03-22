package service

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/laundryos/backend/internal/domain"
	"github.com/laundryos/backend/pkg/hasher"
	"github.com/laundryos/backend/pkg/jwt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrTokenExpired       = errors.New("token expired")
)

type AuthService struct {
	db         *sqlx.DB
	jwtManager *jwt.JWTManager
}

func NewAuthService(db *sqlx.DB, jwtManager *jwt.JWTManager) *AuthService {
	return &AuthService{
		db:         db,
		jwtManager: jwtManager,
	}
}

type RegisterRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	Name       string `json:"name" binding:"required"`
	TenantName string `json:"tenant_name" binding:"required"`
	Phone      string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
	User         *domain.User   `json:"user"`
	Tenant       *domain.Tenant `json:"tenant"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	var exists bool
	err := s.db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var tenant domain.Tenant
	var address, logoURL *string
	tenantQuery := `
		INSERT INTO tenants (name, subdomain, phone, plan, is_active)
		VALUES ($1, $2, $3, 'starter', true)
		RETURNING id, name, subdomain, phone, COALESCE(address, ''), COALESCE(logo_url, ''), plan, is_active, created_at, updated_at
	`
	err = tx.QueryRowContext(ctx, tenantQuery, req.TenantName, generateSubdomain(req.TenantName), req.Phone).
		Scan(&tenant.ID, &tenant.Name, &tenant.Subdomain, &tenant.Phone, &address, &logoURL, &tenant.Plan, &tenant.IsActive, &tenant.CreatedAt, &tenant.UpdatedAt)
	if err != nil {
		return nil, err
	}
	tenant.Address = *address
	tenant.LogoURL = *logoURL

	passwordHash, err := hasher.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	var user domain.User
	userQuery := `
		INSERT INTO users (tenant_id, name, email, password_hash, role, is_active)
		VALUES ($1, $2, $3, $4, 'owner', true)
		RETURNING id, tenant_id, name, email, password_hash, role, is_active, created_at, updated_at
	`
	err = tx.QueryRowContext(ctx, userQuery, tenant.ID, req.Name, req.Email, passwordHash).
		Scan(&user.ID, &user.TenantID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, tenant.ID, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken := &domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: jwt.HashToken(tokenPair.RefreshToken),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}
	refreshQuery := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at, is_revoked)
		VALUES ($1, $2, $3, false)
	`
	_, err = tx.ExecContext(ctx, refreshQuery, refreshToken.UserID, refreshToken.TokenHash, refreshToken.ExpiresAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         &user,
		Tenant:       &tenant,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	var user domain.User
	query := `SELECT id, tenant_id, name, email, password_hash, role, is_active, created_at, updated_at FROM users WHERE email = $1`
	err := s.db.GetContext(ctx, &user, query, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !hasher.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrInvalidCredentials
	}

	var tenant domain.Tenant
	var address, logoURL string
	tenantQuery := `SELECT id, name, subdomain, phone, COALESCE(address, ''), COALESCE(logo_url, ''), plan, is_active, created_at, updated_at FROM tenants WHERE id = $1`
	err = s.db.QueryRowContext(ctx, tenantQuery, user.TenantID).Scan(
		&tenant.ID, &tenant.Name, &tenant.Subdomain, &tenant.Phone,
		&address, &logoURL, &tenant.Plan, &tenant.IsActive, &tenant.CreatedAt, &tenant.UpdatedAt,
	)
	if err != nil {
		return nil, ErrUserNotFound
	}
	tenant.Address = address
	tenant.LogoURL = logoURL

	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, tenant.ID, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken := &domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: jwt.HashToken(tokenPair.RefreshToken),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}
	refreshQuery := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at, is_revoked)
		VALUES ($1, $2, $3, false)
	`
	_, err = s.db.ExecContext(ctx, refreshQuery, refreshToken.UserID, refreshToken.TokenHash, refreshToken.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         &user,
		Tenant:       &tenant,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	tokenHash := jwt.HashToken(refreshToken)

	var storedToken domain.RefreshToken
	query := `SELECT * FROM refresh_tokens WHERE token_hash = $1 AND is_revoked = false`
	err := s.db.GetContext(ctx, &storedToken, query, tokenHash)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if time.Now().After(storedToken.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	_, err = s.db.ExecContext(ctx, `UPDATE refresh_tokens SET is_revoked = true WHERE id = $1`, storedToken.ID)
	if err != nil {
		return nil, err
	}

	var user domain.User
	userQuery := `SELECT * FROM users WHERE id = $1`
	err = s.db.GetContext(ctx, &user, userQuery, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	var tenant domain.Tenant
	tenantQuery := `SELECT * FROM tenants WHERE id = $1`
	err = s.db.GetContext(ctx, &tenant, tenantQuery, claims.TenantID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, tenant.ID, user.Role)
	if err != nil {
		return nil, err
	}

	newRefreshToken := &domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: jwt.HashToken(tokenPair.RefreshToken),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}
	refreshQuery := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at, is_revoked)
		VALUES ($1, $2, $3, false)
	`
	_, err = s.db.ExecContext(ctx, refreshQuery, newRefreshToken.UserID, newRefreshToken.TokenHash, newRefreshToken.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         &user,
		Tenant:       &tenant,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := jwt.HashToken(refreshToken)
	_, err := s.db.ExecContext(ctx, `UPDATE refresh_tokens SET is_revoked = true WHERE token_hash = $1`, tokenHash)
	return err
}

func (s *AuthService) GetCurrentUser(ctx context.Context, userID string) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE id = $1`
	err := s.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func generateSubdomain(name string) string {
	result := ""
	for _, c := range name {
		if c >= 'a' && c <= 'z' || c >= '0' && c <= '9' {
			result += string(c)
		} else if c >= 'A' && c <= 'Z' {
			result += string(c + 32)
		}
	}
	return result
}
