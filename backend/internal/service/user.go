package service

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/laundryos/backend/internal/domain"
	"github.com/laundryos/backend/pkg/hasher"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserInactive       = errors.New("user is inactive")
)

type UserService struct {
	db *sqlx.DB
}

func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{db: db}
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=owner cashier"`
	Phone    string `json:"phone"`
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"omitempty,email"`
	Role     string `json:"role" binding:"omitempty,oneof=owner cashier"`
	IsActive *bool  `json:"is_active"`
	Phone    string `json:"phone"`
}

func (s *UserService) CreateUser(ctx context.Context, tenantID string, req *CreateUserRequest) (*domain.User, error) {
	var exists bool
	err := s.db.GetContext(ctx, &exists,
		`SELECT EXISTS(SELECT 1 FROM users WHERE tenant_id = $1 AND email = $2)`,
		tenantID, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	passwordHash, err := hasher.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	var user domain.User
	query := `
		INSERT INTO users (tenant_id, name, email, password_hash, role, is_active)
		VALUES ($1, $2, $3, $4, $5, true)
		RETURNING id, tenant_id, name, email, role, is_active, created_at, updated_at
	`
	err = s.db.QueryRowContext(ctx, query, tenantID, req.Name, req.Email, passwordHash, req.Role).
		Scan(&user.ID, &user.TenantID, &user.Name, &user.Email, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetUsers(ctx context.Context, tenantID string, page, limit int) ([]domain.User, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM users WHERE tenant_id = $1`
	err := s.db.GetContext(ctx, &total, countQuery, tenantID)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `
		SELECT id, tenant_id, name, email, role, is_active, created_at, updated_at
		FROM users WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	var users []domain.User
	err = s.db.SelectContext(ctx, &users, query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *UserService) GetUserByID(ctx context.Context, tenantID, userID string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, tenant_id, name, email, role, is_active, created_at, updated_at FROM users WHERE id = $1 AND tenant_id = $2`
	err := s.db.GetContext(ctx, &user, query, userID, tenantID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, tenantID, userID string, req *UpdateUserRequest) (*domain.User, error) {
	if req.Email != "" {
		var exists bool
		err := s.db.GetContext(ctx, &exists,
			`SELECT EXISTS(SELECT 1 FROM users WHERE tenant_id = $1 AND email = $2 AND id != $3)`,
			tenantID, req.Email, userID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrEmailAlreadyExists
		}
	}

	setParts := []string{}
	args := []interface{}{}
	argNum := 1

	if req.Name != "" {
		setParts = append(setParts, "name = $"+string(rune('0'+argNum)))
		args = append(args, req.Name)
		argNum++
	}
	if req.Email != "" {
		setParts = append(setParts, "email = $"+string(rune('0'+argNum)))
		args = append(args, req.Email)
		argNum++
	}
	if req.Role != "" {
		setParts = append(setParts, "role = $"+string(rune('0'+argNum)))
		args = append(args, req.Role)
		argNum++
	}
	if req.IsActive != nil {
		setParts = append(setParts, "is_active = $"+string(rune('0'+argNum)))
		args = append(args, *req.IsActive)
		argNum++
	}

	if len(setParts) == 0 {
		return s.GetUserByID(ctx, tenantID, userID)
	}

	setParts = append(setParts, "updated_at = NOW()")

	query := "UPDATE users SET "
	for i, part := range setParts {
		if i > 0 {
			query += ", "
		}
		query += part
	}
	query += " WHERE id = $" + string(rune('0'+argNum)) + " AND tenant_id = $" + string(rune('0'+argNum+1))
	args = append(args, userID, tenantID)

	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return s.GetUserByID(ctx, tenantID, userID)
}

func (s *UserService) DeleteUser(ctx context.Context, tenantID, userID string) error {
	query := `DELETE FROM users WHERE id = $1 AND tenant_id = $2`
	result, err := s.db.ExecContext(ctx, query, userID, tenantID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}
