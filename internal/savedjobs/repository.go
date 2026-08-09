package savedjobs

import (
	"context"
	"errors"
	"fmt"

	"github.com/berzz26/foundry_api/internal/jobs"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

var ErrSavedJobNotFound = errors.New("saved job not found")

const fkViolation = "23503"

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == fkViolation
}

func (r *Repository) Save(ctx context.Context, userID string, jobID int64) (*SavedJob, error) {
	var sj SavedJob
	err := r.db.QueryRow(ctx, `
		INSERT INTO saved_jobs (user_id, job_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, job_id) DO UPDATE SET user_id = EXCLUDED.user_id
		RETURNING id, user_id, job_id, created_at
	`, userID, jobID).Scan(&sj.ID, &sj.UserID, &sj.JobID, &sj.CreatedAt)
	if err != nil {
		if isForeignKeyViolation(err) {
			return nil, ErrSavedJobNotFound
		}
		return nil, err
	}
	return &sj, nil
}

func (r *Repository) Delete(ctx context.Context, userID string, jobID int64) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM saved_jobs WHERE user_id = $1 AND job_id = $2
	`, userID, jobID)
	return err
}

func (r *Repository) IsSaved(ctx context.Context, userID string, jobID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM saved_jobs WHERE user_id = $1 AND job_id = $2)
	`, userID, jobID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// SavedJobWithCompany is a saved job joined with its job and company preview.
type SavedJobWithCompany struct {
	SavedJob
	Job               jobs.JobWithCompany
	CompanyName       *string
	CompanyLogo       *string
	CompanyLogoSource *string
	CompanyBatch      *string
}

const savedJobFields = `
	sj.id, sj.user_id, sj.job_id, sj.created_at
`

const jobFields = `
	j.id, j.company_id, j.title, j.description, j.job_type, j.role, j.location,
	j.remote, j.salary_min, j.salary_max, j.equity_min, j.equity_max,
	j.visa_required, j.job_url, j.created_at, j.updated_at, j.state,
	j.skills::text, j.show_path, j.interview_process, j.time_to_hire, j.visa, j.min_experience,
	j.pretty_eng_type
`

const companyJoinFields = `c.name, c.logo_url, c.logo_source_url, c.batch`

func (r *Repository) ListByUser(ctx context.Context, userID string) ([]SavedJobWithCompany, error) {
	query := fmt.Sprintf(`
		SELECT %s, %s, %s
		FROM saved_jobs sj
		JOIN jobs j ON j.id = sj.job_id
		LEFT JOIN companies c ON j.company_id = c.id
		WHERE sj.user_id = $1
		ORDER BY sj.created_at DESC
	`, savedJobFields, jobFields, companyJoinFields)

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []SavedJobWithCompany
	for rows.Next() {
		var s SavedJobWithCompany
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.JobID, &s.CreatedAt,
			&s.Job.ID, &s.Job.CompanyID, &s.Job.Title, &s.Job.Description, &s.Job.JobType, &s.Job.Role, &s.Job.Location,
			&s.Job.Remote, &s.Job.SalaryMin, &s.Job.SalaryMax, &s.Job.EquityMin, &s.Job.EquityMax,
			&s.Job.VisaRequired, &s.Job.JobURL, &s.Job.CreatedAt, &s.Job.UpdatedAt, &s.Job.State,
			&s.Job.Skills, &s.Job.ShowPath, &s.Job.InterviewProcess, &s.Job.TimeToHire, &s.Job.Visa, &s.Job.MinExperience,
			&s.Job.PrettyEngType,
			&s.CompanyName, &s.CompanyLogo, &s.CompanyLogoSource, &s.CompanyBatch,
		); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}
