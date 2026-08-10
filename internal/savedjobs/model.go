package savedjobs

import "time"

type SavedItem struct {
	ID        int64
	UserID    string
	JobID     *int64
	CompanyID int64
	CreatedAt time.Time
}
