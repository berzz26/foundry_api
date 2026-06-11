package users

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	var user User

	err := r.db.QueryRow(ctx, `
	SELECT id, email, firstName, lastName, profileImageUrl, created_at, updated_at 
	FROM users 
	WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.ProfileImageURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository)AddUser(ctx context.Context, user *User) (*User, error){
	var addedUser User
	
	err:= r.db.QueryRow(ctx, `
	INSERT INTO users (id, email, firstName, lastName, passwordHash, profileImageUrl, provider, providerId, createdAt, updatedAt) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id, email, firstName, lastName, passwordHash, profileImageUrl, provider, providerId, createdAt, updatedAt
	`, user.ID, user.Email, user.FirstName, user.LastName, user.PasswordHash, user.ProfileImageURL, user.Provider, user.ProviderID, user.CreatedAt, user.UpdatedAt).Scan(
		&addedUser.ID,
		&addedUser.Email,
		&addedUser.FirstName,
		&addedUser.LastName,
		&addedUser.PasswordHash,
		&addedUser.ProfileImageURL,
		&addedUser.Provider,
		&addedUser.ProviderID,
		&addedUser.CreatedAt,
		&addedUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &addedUser, nil
}