package companies

import "time"

type Company struct {
	ID                int64      `json:"id"`
	Name              string     `json:"name"`
	Slug              string     `json:"slug"`
	Website           *string    `json:"website"`
	Tagline           *string    `json:"tagline"`
	Description       *string    `json:"description"`
	HiringDescription *string    `json:"hiringDescription"`
	TechStack         *string    `json:"techStack"`
	Batch             *string    `json:"batch"`
	Stage             *string    `json:"stage"`
	TeamSize          *int32     `json:"teamSize"`
	Location          *string    `json:"location"`
	ParentSector      *string    `json:"parentSector"`
	ChildSector       *string    `json:"childSector"`
	Industry          *string    `json:"industry"`
	LogoURL           *string    `json:"logoUrl"`
	SmallLogoURL      *string    `json:"smallLogoUrl"`
	Country           *string    `json:"country"`
	FoundedAt         *time.Time `json:"foundedAt"`
	LinkedinURL       *string    `json:"linkedinUrl"`
	TwitterURL        *string    `json:"twitterUrl"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}
