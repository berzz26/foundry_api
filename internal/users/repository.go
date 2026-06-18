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

	if user.Role == "" {
		user.Role = "user"
	}

	err := r.db.QueryRow(ctx, `
	INSERT INTO users (id, email, first_name, last_name, password_hash, profile_image_url, provider, provider_id, role) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id, email, first_name, last_name, password_hash, profile_image_url, provider, provider_id, role, created_at, updated_at
	`, user.ID, user.Email, user.FirstName, user.LastName, user.PasswordHash, user.ProfileImageURL, user.Provider, user.ProviderID, user.Role).Scan(
		&addedUser.ID,
		&addedUser.Email,
		&addedUser.FirstName,
		&addedUser.LastName,
		&addedUser.PasswordHash,
		&addedUser.ProfileImageURL,
		&addedUser.Provider,
		&addedUser.ProviderID,
		&addedUser.Role,
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
	SELECT id, email, first_name, last_name, password_hash, profile_image_url, provider, provider_id, role, created_at, updated_at 
	FROM users 
	WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.ProfileImageURL,
		&user.Provider,
		&user.ProviderID,
		&user.Role,
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
	SELECT id, email, first_name, last_name, password_hash, profile_image_url, provider, provider_id, role, created_at, updated_at 
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
		&user.Role,
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
	SELECT id, email, first_name, last_name, password_hash, profile_image_url, provider, provider_id, role, created_at, updated_at 
	FROM users 
	ORDER BY created_at DESC
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
			&user.Role,
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

	if user.Role == "" {
		user.Role = "user"
	}

	err := r.db.QueryRow(ctx, `
	UPDATE users SET 
	email = $1, first_name = $2, last_name = $3, password_hash = $4, profile_image_url = $5, provider = $6, provider_id = $7, role = $8, updated_at = NOW() 
	WHERE id = $9
	RETURNING id, email, first_name, last_name, password_hash, profile_image_url, provider, provider_id, role, created_at, updated_at
	`, user.Email, user.FirstName, user.LastName, user.PasswordHash, user.ProfileImageURL, user.Provider, user.ProviderID, user.Role, user.ID).Scan(
		&updatedUser.ID,
		&updatedUser.Email,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.PasswordHash,
		&updatedUser.ProfileImageURL,
		&updatedUser.Provider,
		&updatedUser.ProviderID,
		&updatedUser.Role,
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
