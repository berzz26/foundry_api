package outreach

import "time"

// OutreachCardResponse is one swipe card: company + founder + contact email + best-matched job + pre-written DM.
type OutreachCardResponse struct {
	OutreachID  int64        `json:"outreachId"`
	FounderID   int64        `json:"founderId"`
	Company     CompanyCard  `json:"company"`
	Founder     FounderCard  `json:"founder"`
	Email       EmailCard    `json:"email"`
	Job         *JobCard     `json:"job,omitempty"`
	Outreach    OutreachInfo `json:"outreach"`
	Contactable bool         `json:"contactable"`
}

type CompanyCard struct {
	ID                int64   `json:"id"`
	Name              string  `json:"name"`
	Batch             *string `json:"batch,omitempty"`
	Tagline           *string `json:"tagline,omitempty"`
	Description       *string `json:"description,omitempty"`
	HiringDescription *string `json:"hiringDescription,omitempty"`
	LogoURL           *string `json:"logoUrl,omitempty"`
	SmallLogoURL      *string `json:"smallLogoUrl,omitempty"`
	Website           *string `json:"website,omitempty"`
	Location          *string `json:"location,omitempty"`
	Industry          *string `json:"industry,omitempty"`
	Stage             *string `json:"stage,omitempty"`
	TeamSize          *int32  `json:"teamSize,omitempty"`
	IsHiring          bool    `json:"isHiring"`
}

type FounderCard struct {
	ID        int64   `json:"id"`
	CompanyID *int64  `json:"companyId,omitempty"`
	FullName  string  `json:"fullName"`
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	Linkedin  *string `json:"linkedin,omitempty"`
	Twitter   *string `json:"twitter,omitempty"`
	AvatarURL *string `json:"avatarUrl,omitempty"`
}

type EmailCard struct {
	Address              string  `json:"address"`
	Confidence           *int32  `json:"confidence,omitempty"`
	SMTPValid            *bool   `json:"smtpValid,omitempty"`
	CatchAll             *bool   `json:"catchAll,omitempty"`
	VerificationResponse *string `json:"verificationResponse,omitempty"`
}

type JobCard struct {
	ID               int64      `json:"id"`
	Title            string     `json:"title"`
	Role             *string    `json:"role,omitempty"`
	Description      *string    `json:"description,omitempty"`
	Location         *string    `json:"location,omitempty"`
	Remote           *string    `json:"remote,omitempty"`
	SalaryMin        *int32     `json:"salaryMin,omitempty"`
	SalaryMax        *int32     `json:"salaryMax,omitempty"`
	EquityMin        *float64   `json:"equityMin,omitempty"`
	EquityMax        *float64   `json:"equityMax,omitempty"`
	VisaRequired     *string    `json:"visaRequired,omitempty"`
	JobURL           *string    `json:"jobUrl,omitempty"`
	Skills           []string   `json:"skills"`
	InterviewProcess *string    `json:"interviewProcess,omitempty"`
	MinExperience    *int32     `json:"minExperience,omitempty"`
	TimeToHire       *int32     `json:"timeToHire,omitempty"`
	PrettyEngType    *string    `json:"prettyEngType,omitempty"`
	CreatedAt        *time.Time `json:"createdAt,omitempty"`
}

type OutreachInfo struct {
	Subject     *string    `json:"subject,omitempty"`
	Message     *string    `json:"message,omitempty"`
	Role        *string    `json:"role,omitempty"`
	GeneratedAt *time.Time `json:"generatedAt,omitempty"`
}

type PaginationResponse struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasNext bool  `json:"hasNext"`
}

type OutreachListResponse struct {
	Cards      []OutreachCardResponse `json:"cards"`
	Pagination PaginationResponse     `json:"pagination"`
}

type ListFilters struct {
	Limit     int     `query:"limit"`
	Offset    int     `query:"offset"`
	CompanyID *int64  `query:"companyId"`
	Search    *string `query:"search"`
	HasJob    *bool   `query:"hasJob"`
}

// SendEmailRequest is the body of POST /outreach/:id/send.
type SendEmailRequest struct {
	Subject   *string `json:"subject"`
	Message   string  `json:"message" validate:"required"`
	To        *string `json:"to" validate:"omitempty,email"`
	FounderID *int64  `json:"founderId"`
}

type SendEmailResponse struct {
	SendID     int64     `json:"sendId"`
	OutreachID int64     `json:"outreachId"`
	FounderID  *int64    `json:"founderId,omitempty"`
	To         string    `json:"to"`
	Subject    string    `json:"subject"`
	Status     string    `json:"status"`
	Mode       string    `json:"mode"`
	SentAt     time.Time `json:"sentAt"`
}
