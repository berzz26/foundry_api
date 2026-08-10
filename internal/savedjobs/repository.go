package savedjobs

import (
	"context"
	"errors"
	"fmt"

	"github.com/berzz26/foundry_api/internal/companies"
	"github.com/berzz26/foundry_api/internal/jobs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

var (
	ErrSavedJobNotFound     = errors.New("saved job not found")
	ErrSavedCompanyNotFound = errors.New("saved company not found")
)

const (
	fkViolation      = "23503"
	notNullViolation = "23502"
)

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == fkViolation
}

func isNotNullViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == notNullViolation
}

func (r *Repository) SaveJob(ctx context.Context, userID string, jobID int64) (*SavedItem, error) {
	var si SavedItem
	err := r.db.QueryRow(ctx, `
		INSERT INTO saved_jobs (user_id, job_id, company_id)
		SELECT $1, $2, company_id FROM jobs WHERE id = $2
		ON CONFLICT (user_id, job_id) WHERE job_id IS NOT NULL
		DO UPDATE SET company_id = EXCLUDED.company_id
		RETURNING id, user_id, job_id, company_id, created_at
	`, userID, jobID).Scan(&si.ID, &si.UserID, &si.JobID, &si.CompanyID, &si.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || isForeignKeyViolation(err) || isNotNullViolation(err) {
			return nil, ErrSavedJobNotFound
		}
		return nil, err
	}
	return &si, nil
}

func (r *Repository) SaveCompany(ctx context.Context, userID string, companyID int64) (*SavedItem, error) {
	var si SavedItem
	err := r.db.QueryRow(ctx, `
		INSERT INTO saved_jobs (user_id, company_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, company_id) WHERE job_id IS NULL
		DO UPDATE SET user_id = EXCLUDED.user_id
		RETURNING id, user_id, job_id, company_id, created_at
	`, userID, companyID).Scan(&si.ID, &si.UserID, &si.JobID, &si.CompanyID, &si.CreatedAt)
	if err != nil {
		if isForeignKeyViolation(err) {
			return nil, ErrSavedCompanyNotFound
		}
		return nil, err
	}
	return &si, nil
}

func (r *Repository) DeleteJob(ctx context.Context, userID string, jobID int64) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM saved_jobs WHERE user_id = $1 AND job_id = $2
	`, userID, jobID)
	return err
}

func (r *Repository) DeleteCompany(ctx context.Context, userID string, companyID int64) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM saved_jobs WHERE user_id = $1 AND company_id = $2 AND job_id IS NULL
	`, userID, companyID)
	return err
}

func (r *Repository) IsJobSaved(ctx context.Context, userID string, jobID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM saved_jobs WHERE user_id = $1 AND job_id = $2)
	`, userID, jobID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repository) IsCompanySaved(ctx context.Context, userID string, companyID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM saved_jobs WHERE user_id = $1 AND company_id = $2 AND job_id IS NULL)
	`, userID, companyID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// SavedItemWithDetails is a saved job or company joined with its details.
type SavedItemWithDetails struct {
	SavedItem
	Job     *jobs.JobWithCompany
	Company *companies.Company
}

const savedItemFields = `
	sj.id, sj.user_id, sj.job_id, sj.company_id, sj.created_at
`

const jobFields = `
	j.id, j.company_id, j.title, j.description, j.job_type, j.role, j.location,
	j.remote, j.salary_min, j.salary_max, j.equity_min, j.equity_max,
	j.visa_required, j.job_url, j.created_at, j.updated_at, j.state,
	j.skills::text, j.show_path, j.interview_process, j.time_to_hire, j.visa, j.min_experience,
	j.pretty_eng_type
`

const companyJoinFields = `c.name, c.logo_url, c.logo_source_url, c.batch`

const companyCardFields = `
	co.id, co.name, co.slug, co.tagline, co.batch, co.stage, co.team_size,
	co.location, co.industry, co.logo_url, co.logo_source_url, co.small_logo_url,
	co.small_logo_source_url,
	COALESCE((SELECT COUNT(*) FROM jobs WHERE company_id = co.id), 0) AS open_roles
`

func scanSavedItem(row interface{ Scan(dest ...any) error }) (*SavedItem, error) {
	var si SavedItem
	err := row.Scan(&si.ID, &si.UserID, &si.JobID, &si.CompanyID, &si.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &si, nil
}

func (r *Repository) ListByUser(ctx context.Context, userID string) ([]SavedItemWithDetails, error) {
	jobQuery := fmt.Sprintf(`
		SELECT %s, %s, %s
		FROM saved_jobs sj
		JOIN jobs j ON j.id = sj.job_id
		LEFT JOIN companies c ON j.company_id = c.id
		WHERE sj.user_id = $1
		ORDER BY sj.created_at DESC
	`, savedItemFields, jobFields, companyJoinFields)

	rows, err := r.db.Query(ctx, jobQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []SavedItemWithDetails
	for rows.Next() {
		var s SavedItemWithDetails
		var j jobs.JobWithCompany
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.JobID, &s.CompanyID, &s.CreatedAt,
			&j.ID, &j.CompanyID, &j.Title, &j.Description, &j.JobType, &j.Role, &j.Location,
			&j.Remote, &j.SalaryMin, &j.SalaryMax, &j.EquityMin, &j.EquityMax,
			&j.VisaRequired, &j.JobURL, &j.CreatedAt, &j.UpdatedAt, &j.State,
			&j.Skills, &j.ShowPath, &j.InterviewProcess, &j.TimeToHire, &j.Visa, &j.MinExperience,
			&j.PrettyEngType,
			&j.CompanyName, &j.CompanyLogo, &j.CompanyLogoSource, &j.CompanyBatch,
		); err != nil {
			return nil, err
		}
		s.Job = &j
		list = append(list, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	companyQuery := fmt.Sprintf(`
		SELECT %s, %s
		FROM saved_jobs sj
		JOIN companies co ON co.id = sj.company_id
		WHERE sj.user_id = $1 AND sj.job_id IS NULL
		ORDER BY sj.created_at DESC
	`, savedItemFields, companyCardFields)

	rows2, err := r.db.Query(ctx, companyQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	for rows2.Next() {
		var s SavedItemWithDetails
		var co companies.Company
		if err := rows2.Scan(
			&s.ID, &s.UserID, &s.JobID, &s.CompanyID, &s.CreatedAt,
			&co.ID, &co.Name, &co.Slug, &co.Tagline, &co.Batch, &co.Stage, &co.TeamSize,
			&co.Location, &co.Industry, &co.LogoURL, &co.SourceLogoURL, &co.SmallLogoURL,
			&co.SourceSmallLogoURL, &co.OpenRoles,
		); err != nil {
			return nil, err
		}
		s.Company = &co
		list = append(list, s)
	}
	if err := rows2.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
