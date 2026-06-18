package founders

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

const founderFields = `id, company_id, full_name, first_name, last_name, linkedin, twitter, avatar_url,avatar_source_url, avatar_thumb,avatar_thumb_source_url, avatar_medium`
const founderFieldsAll = `id, company_id, full_name, first_name, last_name, bio, linkedin, twitter, avatar_url,avatar_source_url, avatar_thumb,avatar_thumb_source_url, avatar_medium`

func scanFounder(row interface {
	Scan(dest ...any) error
}) (*Founder, error) {
	var f Founder
	err := row.Scan(
		&f.ID,
		&f.CompanyID,
		&f.FullName,
		&f.FirstName,
		&f.LastName,
		&f.Linkedin,
		&f.Twitter,
		&f.AvatarURL,
		&f.AvatarSourceURL,
		&f.AvatarThumb,
		&f.AvatarSourceThumb,
		&f.AvatarMedium,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}
func scanFounderAll(row interface {
	Scan(dest ...any) error
}) (*Founder, error) {
	var f Founder
	err := row.Scan(
		&f.ID,
		&f.CompanyID,
		&f.FullName,
		&f.FirstName,
		&f.LastName,
		&f.Bio,
		&f.Linkedin,
		&f.Twitter,
		&f.AvatarURL,
		&f.AvatarSourceURL,
		&f.AvatarThumb,
		&f.AvatarSourceThumb,
		&f.AvatarMedium,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *Repository) Create(ctx context.Context, f *Founder) (*Founder, error) {
	var idVal *int64
	if f.ID > 0 {
		idVal = &f.ID
	}

	query := fmt.Sprintf(`
		INSERT INTO founders (id, company_id, full_name, first_name, last_name, bio, linkedin, twitter, avatar_url, avatar_thumb, avatar_medium, created_at, updated_at)
		VALUES (
			CASE WHEN $1::bigint IS NOT NULL AND $1::bigint > 0 THEN $1::bigint ELSE (SELECT COALESCE(MAX(id), 0) + 1 FROM founders) END,
			$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW()
		)
		RETURNING %s
	`, founderFields)

	row := r.db.QueryRow(ctx, query,
		idVal,
		f.CompanyID,
		f.FullName,
		f.FirstName,
		f.LastName,
		f.Bio,
		f.Linkedin,
		f.Twitter,
		f.AvatarURL,
		f.AvatarThumb,
		f.AvatarMedium,
	)
	return scanFounder(row)
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Founder, error) {
	query := fmt.Sprintf(`SELECT %s FROM founders WHERE id = $1`, founderFieldsAll)
	row := r.db.QueryRow(ctx, query, id)
	return scanFounderAll(row)
}
func (r *Repository) GetByCompanyID(ctx context.Context, companyId int64) ([]Founder, error) {
	query := fmt.Sprintf(`SELECT %s FROM founders WHERE company_id = $1`, founderFieldsAll)
	rows, err := r.db.Query(ctx, query, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Founder
	for rows.Next() {
		f, err := scanFounderAll(rows)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		list = append(list, *f)
	}
	if len(list) == 0 {
		return nil, ErrFounderNotFound
	}
	return list, nil
}
func (r *Repository) List(ctx context.Context, companyID *int64, limit, offset int) ([]Founder, error) {
	var args []any
	var conditions []string
	argIndex := 1

	if companyID != nil {
		conditions = append(conditions, fmt.Sprintf("company_id = $%d", argIndex))
		args = append(args, *companyID)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT %s 
		FROM founders 
		%s 
		ORDER BY id DESC
		LIMIT $%d OFFSET $%d
	`, founderFields, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Founder
	for rows.Next() {
		f, err := scanFounder(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *f)
	}
	return list, nil
}

func (r *Repository) Update(ctx context.Context, f *Founder) (*Founder, error) {
	query := fmt.Sprintf(`
		UPDATE founders SET 
			company_id = $1, 
			full_name = $2, 
			first_name = $3, 
			last_name = $4, 
			bio = $5, 
			linkedin = $6, 
			twitter = $7, 
			avatar_url = $8, 
			avatar_thumb = $9, 
			avatar_medium = $10, 
			updated_at = NOW()
		WHERE id = $11
		RETURNING %s
	`, founderFields)

	row := r.db.QueryRow(ctx, query,
		f.CompanyID,
		f.FullName,
		f.FirstName,
		f.LastName,
		f.Bio,
		f.Linkedin,
		f.Twitter,
		f.AvatarURL,
		f.AvatarThumb,
		f.AvatarMedium,
		f.ID,
	)
	return scanFounder(row)
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM founders WHERE id = $1`, id)
	return err
}
