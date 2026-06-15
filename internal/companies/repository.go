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

const companyFields = `id, name, slug, website, tagline, description, hiring_description, tech_stack, batch, stage, team_size, location, parent_sector, child_sector, industry, logo_url, small_logo_url,logo_source_url,small_logo_source_url, country, founded_at, linkedin_url, twitter_url, created_at, updated_at`

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
		&c.SourceLogoURL,
		&c.SourceSmallLogoURL,
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

func buildCompanyFilters(filters CompanyFilters) (string, []any, int) {
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

	if filters.Country != nil && *filters.Country != "" {
		conditions = append(conditions, fmt.Sprintf("country = $%d", argIndex))
		args = append(args, *filters.Country)
		argIndex++
	}

	if filters.MinTeamSize != nil {
		conditions = append(conditions, fmt.Sprintf("team_size >= $%d", argIndex))
		args = append(args, *filters.MinTeamSize)
		argIndex++
	}

	if filters.MaxTeamSize != nil {
		conditions = append(conditions, fmt.Sprintf("team_size <= $%d", argIndex))
		args = append(args, *filters.MaxTeamSize)
		argIndex++
	}

	if filters.Search != nil && *filters.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d OR tagline ILIKE $%d OR industry ILIKE $%d OR parent_sector ILIKE $%d OR child_sector ILIKE $%d)", argIndex, argIndex, argIndex, argIndex, argIndex, argIndex))
		args = append(args, "%"+*filters.Search+"%")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	return whereClause, args, argIndex
}

func (r *Repository) List(ctx context.Context, filters CompanyFilters) ([]Company, error) {
	whereClause, args, argIndex := buildCompanyFilters(filters)

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

	sortClause := "ORDER BY id DESC"
	if filters.Sort != nil {
		switch *filters.Sort {
		case "newest":
			sortClause = "ORDER BY id DESC"
		case "oldest":
			sortClause = "ORDER BY id ASC"
		case "team_size_desc":
			sortClause = "ORDER BY team_size DESC, id DESC"
		case "team_size_asc":
			sortClause = "ORDER BY team_size ASC, id ASC"
		case "name_asc":
			sortClause = "ORDER BY name ASC, id ASC"
		case "name_desc":
			sortClause = "ORDER BY name DESC, id DESC"
		}
	}

	query := fmt.Sprintf(`
		SELECT %s 
		FROM companies 
		%s 
		%s 
		LIMIT $%d OFFSET $%d
	`, companyFields, whereClause, sortClause, argIndex, argIndex+1)

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

func (r *Repository) Count(ctx context.Context, filters CompanyFilters) (int64, error) {
	whereClause, args, _ := buildCompanyFilters(filters)

	query := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM companies 
		%s
	`, whereClause)

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository) GetMetadata(ctx context.Context) (*CompanyMetadataResponse, error) {
	var batches []string
	var industries []string
	var stages []string

	// Query batches
	rows, err := r.db.Query(ctx, "SELECT DISTINCT batch FROM companies WHERE batch IS NOT NULL AND batch != '' ORDER BY batch DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var b string
		if err := rows.Scan(&b); err != nil {
			return nil, err
		}
		batches = append(batches, b)
	}

	// Query industries
	rowsInd, err := r.db.Query(ctx, "SELECT DISTINCT industry FROM companies WHERE industry IS NOT NULL AND industry != '' ORDER BY industry ASC")
	if err != nil {
		return nil, err
	}
	defer rowsInd.Close()
	for rowsInd.Next() {
		var ind string
		if err := rowsInd.Scan(&ind); err != nil {
			return nil, err
		}
		industries = append(industries, ind)
	}

	// Query stages
	rowsStage, err := r.db.Query(ctx, "SELECT DISTINCT stage FROM companies WHERE stage IS NOT NULL AND stage != '' ORDER BY stage ASC")
	if err != nil {
		return nil, err
	}
	defer rowsStage.Close()
	for rowsStage.Next() {
		var s string
		if err := rowsStage.Scan(&s); err != nil {
			return nil, err
		}
		stages = append(stages, s)
	}

	if batches == nil {
		batches = []string{}
	}
	if industries == nil {
		industries = []string{}
	}
	if stages == nil {
		stages = []string{}
	}

	return &CompanyMetadataResponse{
		Batches:    batches,
		Industries: industries,
		Stages:     stages,
	}, nil
}
