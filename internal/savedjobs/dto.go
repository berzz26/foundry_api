package savedjobs

import (
	"time"

	"github.com/berzz26/foundry_api/internal/jobs"
)

type SaveJobRequest struct {
	JobID int64 `json:"jobId" validate:"required"`
}

type SavedJobResponse struct {
	JobID   int64     `json:"jobId"`
	SavedAt time.Time `json:"savedAt"`
}

type SavedJobCard struct {
	jobs.JobCardResponse
	SavedAt time.Time `json:"savedAt"`
}

type SavedStatusResponse struct {
	Saved bool `json:"saved"`
}
