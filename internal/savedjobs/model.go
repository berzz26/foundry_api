package savedjobs

import "time"

type SavedJob struct {
	ID        int64
	UserID    string
	JobID     int64
	CreatedAt time.Time
}
