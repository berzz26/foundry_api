package founders

import "time"

type CreateFounderDTO struct {
	ID                *int64  `json:"id"`
	CompanyID         *int64  `json:"companyId" validate:"required"`
	FullName          string  `json:"fullName" validate:"required"`
	FirstName         *string `json:"firstName"`
	LastName          *string `json:"lastName"`
	Bio               *string `json:"bio"`
	Linkedin          *string `json:"linkedin" validate:"omitempty,url"`
	Twitter           *string `json:"twitter" validate:"omitempty,url"`
	AvatarURL         *string `json:"avatarUrl" validate:"omitempty,url"`
	AvatarSourceURL   *string `json:"avatarSourceUrl" validate:"omitempty,url"`
	AvatarThumb       *string `json:"avatarThumb" validate:"omitempty,url"`
	AvatarSourceThumb *string `json:"avatarSourceThumb" validate:"omitempty,url"`
	AvatarMedium      *string `json:"avatarMedium" validate:"omitempty,url"`
}

type UpdateFounderDTO struct {
	CompanyID    *int64  `json:"companyId"`
	FullName     *string `json:"fullName"`
	FirstName    *string `json:"firstName"`
	LastName     *string `json:"lastName"`
	Bio          *string `json:"bio"`
	Linkedin     *string `json:"linkedin" validate:"omitempty,url"`
	Twitter      *string `json:"twitter" validate:"omitempty,url"`
	AvatarURL    *string `json:"avatarUrl" validate:"omitempty,url"`
	AvatarThumb  *string `json:"avatarThumb" validate:"omitempty,url"`
	AvatarMedium *string `json:"avatarMedium" validate:"omitempty,url"`
}

type FounderResponseDTO struct {
	ID                int64     `json:"id"`
	CompanyID         *int64    `json:"companyId"`
	FullName          string    `json:"fullName"`
	FirstName         *string   `json:"firstName"`
	LastName          *string   `json:"lastName"`
	Bio               *string   `json:"bio"`
	Linkedin          *string   `json:"linkedin"`
	Twitter           *string   `json:"twitter"`
	AvatarURL         *string   `json:"avatarUrl"`
	AvatarSourceURL   *string   `json:"avatarSourceUrl"`
	AvatarThumb       *string   `json:"avatarThumb"`
	AvatarSourceThumb *string   `json:"avatarSourceThumb"`
	AvatarMedium      *string   `json:"avatarMedium"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type PaginationResponse struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasNext bool  `json:"hasNext"`
}

type FounderListResponse struct {
	Founders   []FounderResponseDTO `json:"founders"`
	Pagination PaginationResponse    `json:"pagination"`
}

