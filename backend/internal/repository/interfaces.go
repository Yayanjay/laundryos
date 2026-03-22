package repository

import (
	"context"

	"github.com/laundryos/backend/internal/domain"
)

type TenantRepository interface {
	Create(ctx context.Context, tenant *domain.Tenant) error
	GetByID(ctx context.Context, id string) (*domain.Tenant, error)
	GetBySubdomain(ctx context.Context, subdomain string) (*domain.Tenant, error)
	Update(ctx context.Context, tenant *domain.Tenant) error
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, tenantID, email string) (*domain.User, error)
	GetByTenant(ctx context.Context, tenantID string) ([]domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *domain.RefreshToken) error
	GetByHash(ctx context.Context, hash string) (*domain.RefreshToken, error)
	Revoke(ctx context.Context, id string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error
}
