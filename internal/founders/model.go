package founders

import "time"

type Founder struct {
	ID           int64      `json:"id"`
	CompanyID    *int64     `json:"companyId"`
	FullName     string     `json:"fullName"`
	FirstName    *string    `json:"firstName"`
	LastName     *string    `json:"lastName"`
	Bio          *string    `json:"bio"`
	Linkedin     *string    `json:"linkedin"`
	Twitter      *string    `json:"twitter"`
	AvatarURL    *string    `json:"avatarUrl"`
	AvatarThumb  *string    `json:"avatarThumb"`
	AvatarMedium *string    `json:"avatarMedium"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}