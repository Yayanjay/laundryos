package service

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/laundryos/backend/internal/domain"
)

type TenantService struct {
	db *sqlx.DB
}

func NewTenantService(db *sqlx.DB) *TenantService {
	return &TenantService{db: db}
}

type UpdateTenantRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	LogoURL string `json:"logo_url"`
}

func (s *TenantService) GetTenant(ctx context.Context, tenantID string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	query := `
		SELECT id, name, subdomain, phone, COALESCE(address, ''), COALESCE(logo_url, ''), 
		       plan, is_active, created_at, updated_at 
		FROM tenants WHERE id = $1
	`
	err := s.db.QueryRowContext(ctx, query, tenantID).Scan(
		&tenant.ID, &tenant.Name, &tenant.Subdomain, &tenant.Phone,
		&tenant.Address, &tenant.LogoURL, &tenant.Plan, &tenant.IsActive,
		&tenant.CreatedAt, &tenant.UpdatedAt,
	)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &tenant, nil
}

func (s *TenantService) UpdateTenant(ctx context.Context, tenantID string, req *UpdateTenantRequest) (*domain.Tenant, error) {
	query := `
		UPDATE tenants 
		SET name = COALESCE(NULLIF($2, ''), name),
		    phone = COALESCE(NULLIF($3, ''), phone),
		    address = COALESCE(NULLIF($4, ''), address),
		    logo_url = COALESCE(NULLIF($5, ''), logo_url),
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := s.db.ExecContext(ctx, query, tenantID, req.Name, req.Phone, req.Address, req.LogoURL)
	if err != nil {
		return nil, err
	}

	return s.GetTenant(ctx, tenantID)
}

func (s *TenantService) UpdatePlan(ctx context.Context, tenantID, plan string) error {
	query := `UPDATE tenants SET plan = $2, updated_at = NOW() WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, tenantID, plan)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("tenant not found")
	}
	return nil
}

func (s *TenantService) DeactivateTenant(ctx context.Context, tenantID string) error {
	query := `UPDATE tenants SET is_active = false, updated_at = NOW() WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, tenantID)
	return err
}
