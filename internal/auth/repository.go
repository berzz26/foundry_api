package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshToken struct {
	ID          int
	UserID      string
	TokenHash   string
	ExpiresAt   time.Time
	CreatedAt   time.Time
}

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, tokenHash, expiresAt)
	return err
}

func (r *Repository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	var rt RefreshToken
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`, tokenHash).Scan(&rt.ID, &rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE id = $1`, id)
	return err
}
