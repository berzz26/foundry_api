package companies

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

const companyFields = `id, name, slug, website, tagline, description, hiring_description, tech_stack, batch, stage, team_size, location, parent_sector, child_sector, industry, logo_url, small_logo_url, country, founded_at, linkedin_url, twitter_url, created_at, updated_at`

func scanCompany(row interface {
	Scan(dest ...any) error
}) (*Company, error) {
	var c Company
	err := row.Scan(
		&c.ID,
		&c.Name,
		&c.Slug,
		&c.Website,
		&c.Tagline,
		&c.Description,
		&c.HiringDescription,
		&c.TechStack,
		&c.Batch,
		&c.Stage,
		&c.TeamSize,
		&c.Location,
		&c.ParentSector,
		&c.ChildSector,
		&c.Industry,
		&c.LogoURL,
		&c.SmallLogoURL,
		&c.Country,
		&c.FoundedAt,
		&c.LinkedinURL,
		&c.TwitterURL,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Company, error) {
	query := fmt.Sprintf(`SELECT %s FROM companies WHERE id = $1`, companyFields)
	row := r.db.QueryRow(ctx, query, id)
	return scanCompany(row)
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (*Company, error) {
	query := fmt.Sprintf(`SELECT %s FROM companies WHERE slug = $1`, companyFields)
	row := r.db.QueryRow(ctx, query, slug)
	return scanCompany(row)
}

func (r *Repository) List(ctx context.Context, filters CompanyFilters) ([]Company, error) {
	var args []any
	var conditions []string
	argIndex := 1

	if filters.Industry != nil && *filters.Industry != "" {
		conditions = append(conditions, fmt.Sprintf("industry ILIKE $%d", argIndex))
		args = append(args, "%"+*filters.Industry+"%")
		argIndex++
	}

	if filters.Batch != nil && *filters.Batch != "" {
		conditions = append(conditions, fmt.Sprintf("batch = $%d", argIndex))
		args = append(args, *filters.Batch)
		argIndex++
	}

	if filters.Stage != nil && *filters.Stage != "" {
		conditions = append(conditions, fmt.Sprintf("stage = $%d", argIndex))
		args = append(args, *filters.Stage)
		argIndex++
	}

	if filters.Location != nil && *filters.Location != "" {
		conditions = append(conditions, fmt.Sprintf("location ILIKE $%d", argIndex))
		args = append(args, "%"+*filters.Location+"%")
		argIndex++
	}

	if filters.Search != nil && *filters.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d OR tagline ILIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, "%"+*filters.Search+"%")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	limit := 10
	if filters.Limit > 0 {
		limit = filters.Limit
	}
	if limit > 100 {
		limit = 100
	}

	offset := 0
	if filters.Offset > 0 {
		offset = filters.Offset
	}

	query := fmt.Sprintf(`
		SELECT %s 
		FROM companies 
		%s 
		ORDER BY id DESC 
		LIMIT $%d OFFSET $%d
	`, companyFields, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Company
	for rows.Next() {
		c, err := scanCompany(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *c)
	}

	return list, nil
}
