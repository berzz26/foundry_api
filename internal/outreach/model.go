package outreach

import "time"

// OutreachRecord is a single row from the outreach table used by the send flow.
type OutreachRecord struct {
	ID          int64
	CompanyID   *int64
	CompanyName *string
	FounderName *string
	Subject     *string
	Message     *string
	Role        *string
	GeneratedAt time.Time
}

// outreachCardRow is the flattened join of outreach + founder + email + company + best-matched job.
// It is the source for a single swipe card.
type outreachCardRow struct {
	OutreachID *int64

	CompanyID                 int64
	CompanyName               string
	CompanyBatch              *string
	CompanyTagline            *string
	CompanyDescription        *string
	CompanyHiringDescription  *string
	CompanyLogoURL            *string
	CompanySourceLogoURL      *string
	CompanySmallLogoURL       *string
	CompanySmallSourceLogoURL *string
	CompanyWebsite            *string
	CompanyLocation           *string
	CompanyIndustry           *string
	CompanyStage              *string
	CompanyTeamSize           *int32
	CompanyIsHiring           bool

	FounderID        int64
	FounderCompanyID *int64
	FounderName      string
	FounderFirstName *string
	FounderLastName  *string
	FounderBio       *string
	FounderLinkedin  *string
	FounderTwitter   *string
	FounderAvatarURL *string

	EmailAddress              string
	EmailConfidence           *int32
	EmailSMTPValid            *bool
	EmailCatchAll             *bool
	EmailVerificationResponse *string

	OutreachSubject     *string
	OutreachMessage     *string
	OutreachRole        *string
	OutreachGeneratedAt *time.Time

	JobID               *int64
	JobTitle            *string
	JobDescription      *string
	JobRole             *string
	JobLocation         *string
	JobRemote           *string
	JobSalaryMin        *int32
	JobSalaryMax        *int32
	JobEquityMin        *float64
	JobEquityMax        *float64
	JobVisaRequired     *string
	JobURL              *string
	JobSkills           *string
	JobInterviewProcess *string
	JobMinExperience    *int32
	JobTimeToHire       *int32
	JobPrettyEngType    *string
	JobCreatedAt        *time.Time
}

// SendRecord is a single row written to outreach_sends. OutreachID is nil when
// the send is made for a card that has no outreach (founder-only fallback).
type SendRecord struct {
	ID         int64
	OutreachID *int64
	FounderID  *int64
	Recipient  string
	Subject    *string
	Message    string
	Status     string
	Error      *string
	UserID     *string
	SentAt     time.Time
}
