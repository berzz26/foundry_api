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

func (r *Repository) AddUser(ctx context.Context, user *User) (*User, error) {
	var addedUser User

	err := r.db.QueryRow(ctx, `
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
func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	err := r.db.QueryRow(ctx, `
	SELECT id, email, firstName, lastName, passwordHash, profileImageUrl, provider, providerId, createdAt, updatedAt 
	FROM users 
	WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.ProfileImageURL,
		&user.Provider,
		&user.ProviderID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) List(ctx context.Context, limit int, offset int) ([]User, error) {
	var users []User

	rows, err := r.db.Query(ctx, `
	SELECT id, email, firstName, lastName, passwordHash, profileImageUrl, provider, providerId, createdAt, updatedAt 
	FROM users 
	LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.PasswordHash,
			&user.ProfileImageURL,
			&user.Provider,
			&user.ProviderID,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
func (r *Repository) Update(ctx context.Context, user *User) (*User, error) {
	var updatedUser User

	err := r.db.QueryRow(ctx, `
	UPDATE users SET 
	email = $1, firstName = $2, lastName = $3, passwordHash = $4, profileImageUrl = $5, provider = $6, providerId = $7, createdAt = $8, updatedAt = $9 
	WHERE id = $10
	RETURNING id, email, firstName, lastName, passwordHash, profileImageUrl, provider, providerId, createdAt, updatedAt
	`, user.Email, user.FirstName, user.LastName, user.PasswordHash, user.ProfileImageURL, user.Provider, user.ProviderID, user.CreatedAt, user.UpdatedAt, user.ID).Scan(
		&updatedUser.ID,
		&updatedUser.Email,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.PasswordHash,
		&updatedUser.ProfileImageURL,
		&updatedUser.Provider,
		&updatedUser.ProviderID,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
	DELETE FROM users WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	return nil
}
