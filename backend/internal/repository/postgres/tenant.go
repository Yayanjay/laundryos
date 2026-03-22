package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/laundryos/backend/internal/domain"
)

type TenantRepository struct {
	db *sqlx.DB
}

func NewTenantRepository(db *sqlx.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		INSERT INTO tenants (name, subdomain, phone, address, plan, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		tenant.Name, tenant.Subdomain, tenant.Phone, tenant.Address, tenant.Plan, tenant.IsActive,
	).Scan(&tenant.ID, &tenant.CreatedAt, &tenant.UpdatedAt)
}

func (r *TenantRepository) GetByID(ctx context.Context, id string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	query := `SELECT * FROM tenants WHERE id = $1`
	err := r.db.GetContext(ctx, &tenant, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &tenant, err
}

func (r *TenantRepository) GetBySubdomain(ctx context.Context, subdomain string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	query := `SELECT * FROM tenants WHERE subdomain = $1`
	err := r.db.GetContext(ctx, &tenant, query, subdomain)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &tenant, err
}

func (r *TenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		UPDATE tenants SET name = $2, phone = $3, address = $4, logo_url = $5, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, tenant.ID, tenant.Name, tenant.Phone, tenant.Address, tenant.LogoURL)
	return err
}
