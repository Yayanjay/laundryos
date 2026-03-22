package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/laundryos/backend/internal/domain"
)

type RefreshTokenRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository(db *sqlx.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at, is_revoked)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(ctx, query,
		token.UserID, token.TokenHash, token.ExpiresAt, token.IsRevoked,
	).Scan(&token.ID, &token.CreatedAt)
}

func (r *RefreshTokenRepository) GetByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	var token domain.RefreshToken
	query := `SELECT * FROM refresh_tokens WHERE token_hash = $1 AND is_revoked = false`
	err := r.db.GetContext(ctx, &token, query, hash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &token, err
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, id string) error {
	query := `UPDATE refresh_tokens SET is_revoked = true WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *RefreshTokenRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	query := `UPDATE refresh_tokens SET is_revoked = true WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
