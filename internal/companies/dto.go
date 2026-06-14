package companies

import "time"

type CompanyResponseDTO struct {
	ID                 int64      `json:"id"`
	Name               string     `json:"name"`
	Slug               string     `json:"slug"`
	Website            *string    `json:"website,omitempty"`
	Tagline            *string    `json:"tagline,omitempty"`
	Description        *string    `json:"description,omitempty"`
	HiringDescription  *string    `json:"hiringDescription,omitempty"`
	TechStack          *string    `json:"techStack,omitempty"`
	Batch              *string    `json:"batch,omitempty"`
	Stage              *string    `json:"stage,omitempty"`
	TeamSize           *int32     `json:"teamSize,omitempty"`
	Location           *string    `json:"location,omitempty"`
	ParentSector       *string    `json:"parentSector,omitempty"`
	ChildSector        *string    `json:"childSector,omitempty"`
	Industry           *string    `json:"industry,omitempty"`
	LogoURL            *string    `json:"sourceLogoUrl,omitempty"`
	SourceLogoURL      *string    `json:"logoUrl,omitempty"`
	SmallLogoURL       *string    `json:"smallLogoUrl,omitempty"`
	SourceSmallLogoURL *string    `json:"sourceSmallLogoUrl,omitempty"`
	Country            *string    `json:"country,omitempty"`
	FoundedAt          *time.Time `json:"foundedAt,omitempty"`
	LinkedinURL        *string    `json:"linkedinUrl,omitempty"`
	TwitterURL         *string    `json:"twitterUrl,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
}

type CompanyFilters struct {
	Industry *string `query:"industry"`
	Batch    *string `query:"batch"`
	Stage    *string `query:"stage"`
	Location *string `query:"location"`
	Search   *string `query:"search"`
	Limit    int     `query:"limit"`
	Offset   int     `query:"offset"`
}
