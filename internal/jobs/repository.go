package jobs

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

// Joined representation of a job with its company preview
type JobWithCompany struct {
	Job
	CompanyName      *string
	CompanyLogo      *string
	CompanyLogoSource *string
	CompanyBatch     *string
}

const jobFields = `
	j.id, j.company_id, j.title, j.description, j.job_type, j.role, j.location, 
	j.remote, j.salary_min, j.salary_max, j.equity_min, j.equity_max, 
	j.visa_required, j.job_url, j.created_at, j.updated_at, j.state, 
	j.skills::text, j.show_path, j.interview_process, j.time_to_hire, j.visa, j.min_experience
`

const companyJoinFields = `c.name, c.logo_url, c.logo_source_url, c.batch`

func scanJobWithCompany(row interface{ Scan(dest ...any) error }) (*JobWithCompany, error) {
	var j JobWithCompany
	err := row.Scan(
		&j.ID, &j.CompanyID, &j.Title, &j.Description, &j.JobType, &j.Role, &j.Location,
		&j.Remote, &j.SalaryMin, &j.SalaryMax, &j.EquityMin, &j.EquityMax,
		&j.VisaRequired, &j.JobURL, &j.CreatedAt, &j.UpdatedAt, &j.State,
		&j.Skills, &j.ShowPath, &j.InterviewProcess, &j.TimeToHire, &j.Visa, &j.MinExperience,
		&j.CompanyName, &j.CompanyLogo, &j.CompanyLogoSource, &j.CompanyBatch,
	)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*JobWithCompany, error) {
	query := fmt.Sprintf(`
		SELECT %s, %s 
		FROM jobs j
		LEFT JOIN companies c ON j.company_id = c.id
		WHERE j.id = $1
	`, jobFields, companyJoinFields)

	row := r.db.QueryRow(ctx, query, id)
	return scanJobWithCompany(row)
}

func (r *Repository) List(ctx context.Context, filters JobFilters) ([]JobWithCompany, int, error) {
	var args []any
	var conditions []string
	argIndex := 1

	if filters.Search != nil && *filters.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(j.title ILIKE $%d OR j.description ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+*filters.Search+"%")
		argIndex++
	}

	if filters.Role != nil && *filters.Role != "" {
		conditions = append(conditions, fmt.Sprintf("j.role ILIKE $%d", argIndex))
		args = append(args, "%"+*filters.Role+"%")
		argIndex++
	}

	if filters.Location != nil && *filters.Location != "" {
		conditions = append(conditions, fmt.Sprintf("j.location ILIKE $%d", argIndex))
		args = append(args, "%"+*filters.Location+"%")
		argIndex++
	}

	if filters.Remote != nil && *filters.Remote {
		conditions = append(conditions, "(j.remote = 'yes' OR j.remote = 'only')")
	}

	if filters.SalaryMin != nil && *filters.SalaryMin > 0 {
		conditions = append(conditions, fmt.Sprintf("j.salary_max >= $%d", argIndex))
		args = append(args, *filters.SalaryMin)
		argIndex++
	}

	if filters.SalaryMax != nil && *filters.SalaryMax > 0 {
		conditions = append(conditions, fmt.Sprintf("j.salary_min <= $%d", argIndex))
		args = append(args, *filters.SalaryMax)
		argIndex++
	}

	if filters.VisaRequired != nil && *filters.VisaRequired {
		conditions = append(conditions, "(j.visa_required = 'yes' OR j.visa_required = 'possible')")
	}

	if filters.ExperienceMin != nil {
		conditions = append(conditions, fmt.Sprintf("j.min_experience <= $%d", argIndex))
		args = append(args, *filters.ExperienceMin)
		argIndex++
	}

	if filters.CompanyBatch != nil && *filters.CompanyBatch != "" {
		conditions = append(conditions, fmt.Sprintf("c.batch = $%d", argIndex))
		args = append(args, *filters.CompanyBatch)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`SELECT COUNT(j.id) FROM jobs j LEFT JOIN companies c ON j.company_id = c.id %s`, whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Pagination
	page := filters.Page
	if page < 1 {
		page = 1
	}
	limit := filters.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	// Order
	orderClause := "ORDER BY j.created_at DESC"
	if filters.Sort != nil {
		switch *filters.Sort {
		case "salary_desc":
			orderClause = "ORDER BY j.salary_max DESC NULLS LAST"
		case "salary_asc":
			orderClause = "ORDER BY j.salary_min ASC NULLS LAST"
		case "recent":
			orderClause = "ORDER BY j.created_at DESC NULLS LAST"
		}
	}

	// Data query
	query := fmt.Sprintf(`
		SELECT %s, %s 
		FROM jobs j
		LEFT JOIN companies c ON j.company_id = c.id
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, jobFields, companyJoinFields, whereClause, orderClause, argIndex, argIndex+1)

	dataArgs := append(args, limit, offset)
	rows, err := r.db.Query(ctx, query, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []JobWithCompany
	for rows.Next() {
		j, err := scanJobWithCompany(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, *j)
	}

	return list, total, nil
}

func (r *Repository) GetRelated(ctx context.Context, jobID int64, role string, companyID *int64, limit int) ([]JobWithCompany, error) {
	if limit > 20 {
		limit = 10
	}

	query := fmt.Sprintf(`
		SELECT %s, %s 
		FROM jobs j
		LEFT JOIN companies c ON j.company_id = c.id
		WHERE j.id != $1 AND (j.role = $2 OR (j.company_id = $3 AND $3 IS NOT NULL))
		ORDER BY j.created_at DESC
		LIMIT $4
	`, jobFields, companyJoinFields)

	rows, err := r.db.Query(ctx, query, jobID, role, companyID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []JobWithCompany
	for rows.Next() {
		j, err := scanJobWithCompany(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *j)
	}
	return list, nil
}

func (r *Repository) GetFeatured(ctx context.Context, limit int) ([]JobWithCompany, error) {
	if limit > 20 {
		limit = 10
	}
	// Fetch high salary or recent ones
	query := fmt.Sprintf(`
		SELECT %s, %s 
		FROM jobs j
		LEFT JOIN companies c ON j.company_id = c.id
		ORDER BY j.salary_max DESC NULLS LAST, j.created_at DESC
		LIMIT $1
	`, jobFields, companyJoinFields)

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []JobWithCompany
	for rows.Next() {
		j, err := scanJobWithCompany(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *j)
	}
	return list, nil
}
