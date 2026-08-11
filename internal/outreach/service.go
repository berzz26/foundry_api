package outreach

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Service struct {
	repo   *Repository
	sender EmailSender
}

type EmailSender interface {
	Configured() bool
	FromAddress() string
	Send(to, subject, body string) error
}

func NewService(repo *Repository, sender EmailSender) *Service {
	return &Service{repo: repo, sender: sender}
}

func (s *Service) List(ctx context.Context, filters ListFilters) (*OutreachListResponse, error) {
	seed := filters.Seed
	if seed == nil || *seed == "" {
		s := randomSeed()
		seed = &s
	}

	cards, total, err := s.repo.ListCards(ctx, filters, *seed)
	if err != nil {
		return nil, err
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

	res := make([]OutreachCardResponse, 0, len(cards))
	for _, c := range cards {
		res = append(res, mapCardToResponse(c))
	}

	return &OutreachListResponse{
		Cards: res,
		Pagination: PaginationResponse{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			HasNext: int64(offset+limit) < total,
		},
		Seed: *seed,
	}, nil
}

func randomSeed() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

func (s *Service) GetByID(ctx context.Context, outreachID int64, founderID *int64) (*OutreachCardResponse, error) {
	card, err := s.repo.GetCard(ctx, outreachID, founderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOutreachNotFound
		}
		return nil, err
	}
	res := mapCardToResponse(*card)
	return &res, nil
}

// Send delivers a message to the founder email for an outreach and records the attempt.
// When SMTP is not configured it runs in dry-run mode and records the send with
// status "dry_run" instead of delivering anything.
//
// outreachID may be 0 when the card has no outreach. In that case the founder ID
// from the request payload is used as the fallback and the subject/message come
// exclusively from the request payload.
func (s *Service) Send(ctx context.Context, outreachID int64, req SendEmailRequest, userID string) (*SendEmailResponse, error) {
	var outreach *OutreachRecord
	if outreachID > 0 {
		o, err := s.repo.GetOutreach(ctx, outreachID)
		if err != nil && !errors.Is(err, ErrOutreachNotFound) {
			return nil, err
		}
		if err == nil {
			outreach = o
		}
	}

	var recipient string
	switch {
	case req.To != nil && *req.To != "":
		recipient = *req.To
	case req.FounderID != nil:
		email, err := s.repo.GetFounderEmail(ctx, *req.FounderID)
		if err != nil {
			return nil, err
		}
		recipient = email
	case outreach != nil:
		email, err := s.repo.GetOutreachFounderEmail(ctx, outreachID)
		if err != nil {
			return nil, err
		}
		recipient = email
	default:
		return nil, ErrNoRecipient
	}

	subject := ""
	if req.Subject != nil && *req.Subject != "" {
		subject = *req.Subject
	} else if outreach != nil && outreach.Subject != nil {
		subject = *outreach.Subject
	}

	status := "sent"
	mode := "sent"
	if s.sender.Configured() {
		if err := s.sender.Send(recipient, subject, req.Message); err != nil {
			errStr := err.Error()
			rec := &SendRecord{
				OutreachID: outreachIDPtr(outreachID, outreach),
				FounderID:  req.FounderID,
				Recipient:  recipient,
				Subject:    &subject,
				Message:    req.Message,
				Status:     "failed",
				Error:      &errStr,
				UserID:     &userID,
			}
			s.repo.RecordSend(ctx, rec)
			return nil, err
		}
	} else {
		status = "dry_run"
		mode = "dry-run"
	}

	var founderID *int64
	if req.FounderID != nil {
		founderID = req.FounderID
	}

	rec := &SendRecord{
		OutreachID: outreachIDPtr(outreachID, outreach),
		FounderID:  founderID,
		Recipient:  recipient,
		Subject:    &subject,
		Message:    req.Message,
		Status:     status,
		UserID:     &userID,
	}
	sendID, err := s.repo.RecordSend(ctx, rec)
	if err != nil {
		return nil, err
	}

	return &SendEmailResponse{
		SendID:     sendID,
		OutreachID: outreachIDPtr(outreachID, outreach),
		FounderID:  founderID,
		To:         recipient,
		Subject:    subject,
		Status:     status,
		Mode:       mode,
		SentAt:     time.Now(),
	}, nil
}

// outreachIDPtr returns the outreach ID when an outreach record was resolved,
// and nil otherwise (founder-only send with no outreach).
func outreachIDPtr(outreachID int64, outreach *OutreachRecord) *int64 {
	if outreach != nil {
		return &outreachID
	}
	return nil
}

func mapCardToResponse(c outreachCardRow) OutreachCardResponse {
	company := CompanyCard{
		ID:                c.CompanyID,
		Name:              c.CompanyName,
		Batch:             c.CompanyBatch,
		Tagline:           c.CompanyTagline,
		Description:       c.CompanyDescription,
		HiringDescription: c.CompanyHiringDescription,
		Website:           c.CompanyWebsite,
		Location:          c.CompanyLocation,
		Industry:          c.CompanyIndustry,
		Stage:             c.CompanyStage,
		TeamSize:          c.CompanyTeamSize,
		IsHiring:          c.CompanyIsHiring,
	}
	if c.CompanySourceLogoURL != nil && *c.CompanySourceLogoURL != "" {
		company.LogoURL = c.CompanySourceLogoURL
	} else {
		company.LogoURL = c.CompanyLogoURL
	}
	if c.CompanySmallSourceLogoURL != nil && *c.CompanySmallSourceLogoURL != "" {
		company.SmallLogoURL = c.CompanySmallSourceLogoURL
	} else {
		company.SmallLogoURL = c.CompanySmallLogoURL
	}

	founder := FounderCard{
		ID:        c.FounderID,
		CompanyID: c.FounderCompanyID,
		FullName:  c.FounderName,
		FirstName: c.FounderFirstName,
		LastName:  c.FounderLastName,
		Bio:       c.FounderBio,
		Linkedin:  c.FounderLinkedin,
		Twitter:   c.FounderTwitter,
		AvatarURL: c.FounderAvatarURL,
	}

	email := EmailCard{
		Address:              c.EmailAddress,
		Confidence:           c.EmailConfidence,
		SMTPValid:            c.EmailSMTPValid,
		CatchAll:             c.EmailCatchAll,
		VerificationResponse: c.EmailVerificationResponse,
	}

	var job *JobCard
	if c.JobID != nil {
		skills := []string{}
		if c.JobSkills != nil && *c.JobSkills != "" && *c.JobSkills != "[]" {
			json.Unmarshal([]byte(*c.JobSkills), &skills)
		}
		job = &JobCard{
			ID:               *c.JobID,
			Title:            deref(c.JobTitle, ""),
			Role:             c.JobRole,
			Description:      c.JobDescription,
			Location:         c.JobLocation,
			Remote:           c.JobRemote,
			SalaryMin:        c.JobSalaryMin,
			SalaryMax:        c.JobSalaryMax,
			EquityMin:        c.JobEquityMin,
			EquityMax:        c.JobEquityMax,
			VisaRequired:     c.JobVisaRequired,
			JobURL:           c.JobURL,
			Skills:           skills,
			InterviewProcess: c.JobInterviewProcess,
			MinExperience:    c.JobMinExperience,
			TimeToHire:       c.JobTimeToHire,
			PrettyEngType:    c.JobPrettyEngType,
			CreatedAt:        c.JobCreatedAt,
		}
	}

	return OutreachCardResponse{
		OutreachID: c.OutreachID,
		FounderID:  c.FounderID,
		Company:    company,
		Founder:    founder,
		Email:      email,
		Job:        job,
		Outreach: OutreachInfo{
			Subject:     c.OutreachSubject,
			Message:     c.OutreachMessage,
			Role:        c.OutreachRole,
			GeneratedAt: c.OutreachGeneratedAt,
		},
		Contactable: c.EmailAddress != "",
	}
}

func deref[T any](p *T, fallback T) T {
	if p == nil {
		return fallback
	}
	return *p
}
