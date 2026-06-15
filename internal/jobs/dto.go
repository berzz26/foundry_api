package jobs

import "time"

type JobCardResponse struct {
	ID         int64                  `json:"id"`
	Title      string                 `json:"title"`
	Location   *string                `json:"location"`
	Remote     bool                   `json:"remote"`
	Salary     *SalaryResponse        `json:"salary,omitempty"`
	Equity     *EquityResponse        `json:"equity,omitempty"`
	Experience *ExperienceResponse    `json:"experience,omitempty"`
	Company    CompanyPreviewResponse `json:"company"`
	CreatedAt  *time.Time             `json:"createdAt,omitempty"`
}

type JobDetailResponse struct {
	ID               int64                  `json:"id"`
	Title            string                 `json:"title"`
	Description      *string                `json:"description,omitempty"`
	Location         *string                `json:"location"`
	Remote           bool                   `json:"remote"`
	Salary           *SalaryResponse        `json:"salary,omitempty"`
	Equity           *EquityResponse        `json:"equity,omitempty"`
	Experience       *ExperienceResponse    `json:"experience,omitempty"`
	VisaRequired     bool                   `json:"visaRequired"`
	Skills           []string               `json:"skills"`
	InterviewProcess *string                `json:"interviewProcess,omitempty"`
	TimeToHire       *int32                 `json:"timeToHire,omitempty"`
	Company          CompanyPreviewResponse `json:"company"`
	CreatedAt        *time.Time             `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time             `json:"updatedAt,omitempty"`
}

type CompanyPreviewResponse struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	LogoURL *string `json:"logoUrl,omitempty"`
	Batch   *string `json:"batch,omitempty"`
}

type SalaryResponse struct {
	Min *int32 `json:"min,omitempty"`
	Max *int32 `json:"max,omitempty"`
}

type EquityResponse struct {
	Min *float64 `json:"min,omitempty"`
	Max *float64 `json:"max,omitempty"`
}

type ExperienceResponse struct {
	MinYears *int32 `json:"minYears,omitempty"`
}

type PaginationResponse struct {
	Page    int  `json:"page"`
	Limit   int  `json:"limit"`
	Total   int  `json:"total"`
	HasNext bool `json:"hasNext"`
}

type JobListResponse struct {
	Jobs       []JobCardResponse  `json:"jobs"`
	Pagination PaginationResponse `json:"pagination"`
}

type JobRelatedResponse struct {
	Jobs []JobCardResponse `json:"jobs"`
}

type JobFeaturedResponse struct {
	Jobs []JobCardResponse `json:"jobs"`
}

type JobFilters struct {
	Page         int      `query:"page"`
	Limit        int      `query:"limit"`
	Search       *string  `query:"search"`
	Role         *string  `query:"role"`
	Location     *string  `query:"location"`
	Remote       *bool    `query:"remote"`
	SalaryMin    *int32   `query:"salaryMin"`
	SalaryMax    *int32   `query:"salaryMax"`
	VisaRequired *bool    `query:"visaRequired"`
	ExperienceMin *int32  `query:"experienceMin"`
	CompanyBatch *string  `query:"companyBatch"`
	CompanyID    *int64   `query:"companyId"`
	Sort         *string  `query:"sort"`
}
