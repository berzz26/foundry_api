package founders

import "time"

type CreateFounderDTO struct {
	ID           *int64  `json:"id"`
	CompanyID    *int64  `json:"companyId" validate:"required"`
	FullName     string  `json:"fullName" validate:"required"`
	FirstName    *string `json:"firstName"`
	LastName     *string `json:"lastName"`
	Bio          *string `json:"bio"`
	Linkedin     *string `json:"linkedin" validate:"omitempty,url"`
	Twitter      *string `json:"twitter" validate:"omitempty,url"`
	AvatarURL    *string `json:"avatarUrl" validate:"omitempty,url"`
	AvatarThumb  *string `json:"avatarThumb" validate:"omitempty,url"`
	AvatarMedium *string `json:"avatarMedium" validate:"omitempty,url"`
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
