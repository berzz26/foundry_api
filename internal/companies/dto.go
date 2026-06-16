package companies

import "time"

type CompanyCardResponse struct {
	ID                 int64  `json:"id"`
	Name               string `json:"name"`
	Slug               string `json:"slug"`
	Tagline            string `json:"tagline"`
	Batch              string `json:"batch"`
	OpenRoles          int    `json:"openRoles"`
	Stage              string `json:"stage"`
	TeamSize           int    `json:"teamSize"`
	Location           string `json:"location"`
	Industry           string `json:"industry"`
	LogoURL            string `json:"logoUrl"`
	SourceLogoURL      string `json:"sourceLogoUrl"`
	SmallLogoURL       string `json:"smallLogoUrl"`
	SourceSmallLogoURL string `json:"sourceSmallLogoUrl"`
}

type CompanyDetailResponse struct {
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
	LogoURL            *string    `json:"logoUrl,omitempty"`
	SourceLogoURL      *string    `json:"sourceLogoUrl,omitempty"`
	SmallLogoURL       *string    `json:"smallLogoUrl,omitempty"`
	SourceSmallLogoURL *string    `json:"sourceSmallLogoUrl,omitempty"`
	Country            *string    `json:"country,omitempty"`
	FoundedAt          *time.Time `json:"foundedAt,omitempty"`
	LinkedinURL        *string    `json:"linkedinUrl,omitempty"`
	TwitterURL         *string    `json:"twitterUrl,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
}

type PaginationResponse struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasNext bool  `json:"hasNext"`
}

type CompanyListResponse struct {
	Companies  []CompanyCardResponse `json:"companies"`
	Pagination PaginationResponse    `json:"pagination"`
}

type CompanyMetadataResponse struct {
	Batches        []string `json:"batches"`
	Industries     []string `json:"industries"`
	Stages         []string `json:"stages"`
	TotalCompanies int64    `json:"totalCompanies"`
}

type CompanyFilters struct {
	Industry    *string `query:"industry"`
	Batch       *string `query:"batch"`
	Stage       *string `query:"stage"`
	Location    *string `query:"location"`
	Search      *string `query:"search"`
	MinTeamSize *int    `query:"minTeamSize"`
	MaxTeamSize *int    `query:"maxTeamSize"`
	Country     *string `query:"country"`
	Sort        *string `query:"sort"`
	Limit       int     `query:"limit"`
	Offset      int     `query:"offset"`
}
