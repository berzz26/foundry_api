package jobs

import (
	"time"
)

type Job struct {
	ID               int64      `json:"id"`
	CompanyID        *int64     `json:"companyId"`
	Title            string     `json:"title"`
	Description      *string    `json:"description"`
	JobType          *string    `json:"jobType"`
	Role             *string    `json:"role"`
	Location         *string    `json:"location"`
	Remote           *string    `json:"remote"`
	SalaryMin        *int32     `json:"salaryMin"`
	SalaryMax        *int32     `json:"salaryMax"`
	EquityMin        *float64   `json:"equityMin"`
	EquityMax        *float64   `json:"equityMax"`
	VisaRequired     *string    `json:"visaRequired"`
	JobURL           *string    `json:"jobUrl"`
	CreatedAt        *time.Time `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt"`
	State            *string    `json:"state"`
	Skills           *string    `json:"skills"` // pgx scanning jsonb to string/[]byte or just string, or custom struct
	ShowPath         *string    `json:"showPath"`
	InterviewProcess *string    `json:"interviewProcess"`
	TimeToHire       *int32     `json:"timeToHire"`
	Visa             *string    `json:"visa"`
	MinExperience    *int32     `json:"minExperience"`
}
