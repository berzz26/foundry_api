package outreach

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrOutreachNotFound = errors.New("outreach not found")
	ErrNoRecipient      = errors.New("no recipient email available for this outreach")
	ErrNoFounderEmail   = errors.New("founder has no verified email")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// deckCardFields selects the flattened swipe card. The LATERAL join picks the
// best job for the outreach: an exact title match first, then a fuzzy match,
// then the company's most recent job.
const deckCardFields = `
	o.id AS outreach_id,

	c.id AS company_id,
	c.name AS company_name,
	c.batch AS company_batch,
	c.tagline AS company_tagline,
	c.description AS company_description,
	c.hiring_description AS company_hiring_description,
	c.logo_url AS company_logo_url,
	c.logo_source_url AS company_source_logo_url,
	c.small_logo_url AS company_small_logo_url,
	c.small_logo_source_url AS company_small_source_logo_url,
	c.website AS company_website,
	c.location AS company_location,
	c.industry AS company_industry,
	c.stage AS company_stage,
	c.team_size AS company_team_size,
	c.is_hiring AS company_is_hiring,

	f.id AS founder_id,
	f.company_id AS founder_company_id,
	COALESCE(f.full_name, o.founder_name) AS founder_name,
	f.first_name AS founder_first_name,
	f.last_name AS founder_last_name,
	f.bio AS founder_bio,
	f.linkedin AS founder_linkedin,
	f.twitter AS founder_twitter,
	f.avatar_url AS founder_avatar_url,

	fe.email AS email_address,
	fe.confidence AS email_confidence,
	fe.smtp_valid AS email_smtp_valid,
	fe.catch_all AS email_catch_all,
	fe.verification_response AS email_verification_response,

	o.subject AS outreach_subject,
	o.message AS outreach_message,
	o.role AS outreach_role,
	o.generated_at AS outreach_generated_at,

	j.id AS job_id,
	j.title AS job_title,
	j.description AS job_description,
	j.role AS job_role,
	j.location AS job_location,
	j.remote AS job_remote,
	j.salary_min AS job_salary_min,
	j.salary_max AS job_salary_max,
	j.equity_min AS job_equity_min,
	j.equity_max AS job_equity_max,
	j.visa_required AS job_visa_required,
	j.job_url AS job_url,
	j.skills::text AS job_skills,
	j.interview_process AS job_interview_process,
	j.min_experience AS job_min_experience,
	j.time_to_hire AS job_time_to_hire,
	j.pretty_eng_type AS job_pretty_eng_type,
	j.created_at AS job_created_at
`

const deckFromClause = `
	FROM outreach o
	JOIN founders f ON f.company_id = o.company_id
	JOIN founder_emails2 fe ON fe.founder_id = f.id
	LEFT JOIN companies c ON c.id = o.company_id
	LEFT JOIN LATERAL (
		SELECT j.*
		FROM jobs j
		WHERE j.company_id = o.company_id
		ORDER BY
			CASE
				WHEN o.role IS NOT NULL AND TRIM(o.role) <> '' AND j.title = o.role THEN 0
				WHEN o.role IS NOT NULL AND TRIM(o.role) <> '' AND (j.title ILIKE '%' || o.role || '%' OR o.role ILIKE '%' || j.title || '%') THEN 1
				ELSE 2
			END,
			j.created_at DESC
		LIMIT 1
	) j ON TRUE
`

const deckFromEmailsClause = `
	FROM founder_emails2 fe
	JOIN founders f ON f.id = fe.founder_id
	LEFT JOIN companies c ON c.id = f.company_id
	LEFT JOIN LATERAL (
		SELECT o.*
		FROM outreach o
		WHERE o.company_id = f.company_id
		ORDER BY o.generated_at DESC, o.id DESC
		LIMIT 1
	) o ON TRUE
	LEFT JOIN LATERAL (
		SELECT j.*
		FROM jobs j
		WHERE j.company_id = f.company_id
		ORDER BY
			CASE
				WHEN o.role IS NOT NULL AND TRIM(o.role) <> '' AND j.title = o.role THEN 0
				WHEN o.role IS NOT NULL AND TRIM(o.role) <> '' AND (j.title ILIKE '%' || o.role || '%' OR o.role ILIKE '%' || j.title || '%') THEN 1
				ELSE 2
			END,
			j.created_at DESC
		LIMIT 1
	) j ON TRUE
`

const deckWhereClause = `
	WHERE fe.email IS NOT NULL AND TRIM(fe.email) <> ''
`

func scanOutreachCard(row interface {
	Scan(dest ...any) error
}) (*outreachCardRow, error) {
	var c outreachCardRow
	err := row.Scan(
		&c.OutreachID,

		&c.CompanyID,
		&c.CompanyName,
		&c.CompanyBatch,
		&c.CompanyTagline,
		&c.CompanyDescription,
		&c.CompanyHiringDescription,
		&c.CompanyLogoURL,
		&c.CompanySourceLogoURL,
		&c.CompanySmallLogoURL,
		&c.CompanySmallSourceLogoURL,
		&c.CompanyWebsite,
		&c.CompanyLocation,
		&c.CompanyIndustry,
		&c.CompanyStage,
		&c.CompanyTeamSize,
		&c.CompanyIsHiring,

		&c.FounderID,
		&c.FounderCompanyID,
		&c.FounderName,
		&c.FounderFirstName,
		&c.FounderLastName,
		&c.FounderBio,
		&c.FounderLinkedin,
		&c.FounderTwitter,
		&c.FounderAvatarURL,

		&c.EmailAddress,
		&c.EmailConfidence,
		&c.EmailSMTPValid,
		&c.EmailCatchAll,
		&c.EmailVerificationResponse,

		&c.OutreachSubject,
		&c.OutreachMessage,
		&c.OutreachRole,
		&c.OutreachGeneratedAt,

		&c.JobID,
		&c.JobTitle,
		&c.JobDescription,
		&c.JobRole,
		&c.JobLocation,
		&c.JobRemote,
		&c.JobSalaryMin,
		&c.JobSalaryMax,
		&c.JobEquityMin,
		&c.JobEquityMax,
		&c.JobVisaRequired,
		&c.JobURL,
		&c.JobSkills,
		&c.JobInterviewProcess,
		&c.JobMinExperience,
		&c.JobTimeToHire,
		&c.JobPrettyEngType,
		&c.JobCreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) ListCards(ctx context.Context, filters ListFilters, seed string) ([]outreachCardRow, int64, error) {
	var args []any
	var conditions []string
	argIndex := 1

	if filters.CompanyID != nil {
		conditions = append(conditions, fmt.Sprintf("f.company_id = $%d", argIndex))
		args = append(args, *filters.CompanyID)
		argIndex++
	}
	if filters.Search != nil && *filters.Search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(c.name ILIKE $%d OR o.company_name ILIKE $%d OR COALESCE(f.full_name, o.founder_name) ILIKE $%d OR j.title ILIKE $%d)",
			argIndex, argIndex, argIndex, argIndex))
		args = append(args, "%"+*filters.Search+"%")
		argIndex++
	}
	if filters.HasJob != nil {
		if *filters.HasJob {
			conditions = append(conditions, "j.id IS NOT NULL")
		} else {
			conditions = append(conditions, "j.id IS NULL")
		}
	}

	where := deckWhereClause
	if len(conditions) > 0 {
		where += " AND " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) %s %s", deckFromEmailsClause, where)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limit := filters.Limit
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset := filters.Offset
	if offset < 0 {
		offset = 0
	}

	dataArgs := append(args, seed, limit, offset)
	query := fmt.Sprintf(`
		SELECT %s
		%s
		%s
		ORDER BY
			CASE WHEN c.batch ~ '[0-9]+$' THEN SUBSTRING(c.batch FROM '[0-9]+$')::int ELSE NULL END DESC NULLS LAST,
			md5($%d::text || ':' || c.id::text),
			c.id
		LIMIT $%d OFFSET $%d
	`, deckCardFields, deckFromEmailsClause, where, argIndex, argIndex+1, argIndex+2)

	rows, err := r.db.Query(ctx, query, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var cards []outreachCardRow
	for rows.Next() {
		c, err := scanOutreachCard(rows)
		if err != nil {
			return nil, 0, err
		}
		cards = append(cards, *c)
	}

	return cards, total, nil
}

func (r *Repository) GetCard(ctx context.Context, outreachID int64, founderID *int64) (*outreachCardRow, error) {
	conditions := "AND o.id = $1"
	args := []any{outreachID}
	if founderID != nil {
		conditions += " AND f.id = $2"
		args = append(args, *founderID)
	}

	query := fmt.Sprintf(`
		SELECT %s
		%s
		%s
		%s
		ORDER BY o.generated_at DESC, o.id DESC
		LIMIT 1
	`, deckCardFields, deckFromClause, deckWhereClause, conditions)

	row := r.db.QueryRow(ctx, query, args...)
	return scanOutreachCard(row)
}

func (r *Repository) GetOutreach(ctx context.Context, id int64) (*OutreachRecord, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, company_id, company_name, founder_name, subject, message, role, generated_at
		FROM outreach
		WHERE id = $1
	`, id)

	var o OutreachRecord
	err := row.Scan(
		&o.ID,
		&o.CompanyID,
		&o.CompanyName,
		&o.FounderName,
		&o.Subject,
		&o.Message,
		&o.Role,
		&o.GeneratedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOutreachNotFound
		}
		return nil, err
	}
	return &o, nil
}

func (r *Repository) GetFounderEmail(ctx context.Context, founderID int64) (string, error) {
	var email string
	err := r.db.QueryRow(ctx, `
		SELECT email FROM founder_emails2
		WHERE founder_id = $1 AND email IS NOT NULL AND TRIM(email) <> ''
		LIMIT 1
	`, founderID).Scan(&email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNoFounderEmail
		}
		return "", err
	}
	return email, nil
}

// GetOutreachFounderEmail returns the first verified founder email for an outreach's company.
func (r *Repository) GetOutreachFounderEmail(ctx context.Context, outreachID int64) (string, error) {
	var email string
	err := r.db.QueryRow(ctx, `
		SELECT fe.email
		FROM outreach o
		JOIN founders f ON f.company_id = o.company_id
		JOIN founder_emails2 fe ON fe.founder_id = f.id
		WHERE o.id = $1 AND fe.email IS NOT NULL AND TRIM(fe.email) <> ''
		ORDER BY f.id
		LIMIT 1
	`, outreachID).Scan(&email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNoRecipient
		}
		return "", err
	}
	return email, nil
}

func (r *Repository) RecordSend(ctx context.Context, rec *SendRecord) (int64, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO outreach_sends (
			outreach_id, founder_id, recipient_email, subject, message, status, error, user_id, sent_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, rec.OutreachID, rec.FounderID, rec.Recipient, rec.Subject, rec.Message, rec.Status, rec.Error, rec.UserID, time.Now())

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
