package savedjobs

import (
	"time"

	"github.com/berzz26/foundry_api/internal/companies"
	"github.com/berzz26/foundry_api/internal/jobs"
)

type SaveJobRequest struct {
	JobID int64 `json:"jobId" validate:"required"`
}

type SaveCompanyRequest struct {
	CompanyID int64 `json:"companyId" validate:"required"`
}

type SaveResponse struct {
	JobID     *int64    `json:"jobId,omitempty"`
	CompanyID *int64    `json:"companyId,omitempty"`
	SavedAt   time.Time `json:"savedAt"`
}

type SavedItemCard struct {
	Type    string                         `json:"type"` // "job" | "company"
	Job     *jobs.JobCardResponse          `json:"job,omitempty"`
	Company *companies.CompanyCardResponse `json:"company,omitempty"`
	SavedAt time.Time                      `json:"savedAt"`
}

type SavedStatusResponse struct {
	Saved bool `json:"saved"`
}
