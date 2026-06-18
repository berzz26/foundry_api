package founders

import (
	"testing"
)

func TestMapToResponseDTO(t *testing.T) {
	fullName := "Jane Doe"
	firstName := "Jane"
	lastName := "Doe"
	bio := "Serial entrepreneur"
	linkedin := "https://linkedin.com/in/janedoe"
	twitter := "https://twitter.com/janedoe"
	avatar := "https://example.com/avatar.jpg"
	avatarThumb := "https://example.com/avatar_thumb.jpg"
	avatarMedium := "https://example.com/avatar_medium.jpg"
	companyID := int64(12738)

	f := &Founder{
		ID:           42,
		CompanyID:    &companyID,
		FullName:     fullName,
		FirstName:    &firstName,
		LastName:     &lastName,
		Bio:          &bio,
		Linkedin:     &linkedin,
		Twitter:      &twitter,
		AvatarURL:    &avatar,
		AvatarThumb:  &avatarThumb,
		AvatarMedium: &avatarMedium,
	}

	dto := mapToResponseDTO(f)

	if dto.ID != 42 {
		t.Errorf("Expected ID 42, got %d", dto.ID)
	}
	if dto.FullName != fullName {
		t.Errorf("Expected full name %q, got %q", fullName, dto.FullName)
	}
	if dto.FirstName == nil || *dto.FirstName != firstName {
		t.Errorf("Expected first name %q", firstName)
	}
	if dto.LastName == nil || *dto.LastName != lastName {
		t.Errorf("Expected last name %q", lastName)
	}
	if dto.Bio == nil || *dto.Bio != bio {
		t.Errorf("Expected bio %q", bio)
	}
	if dto.Linkedin == nil || *dto.Linkedin != linkedin {
		t.Errorf("Expected linkedin URL %q", linkedin)
	}
	if dto.Twitter == nil || *dto.Twitter != twitter {
		t.Errorf("Expected twitter URL %q", twitter)
	}
	if dto.AvatarURL == nil || *dto.AvatarURL != avatar {
		t.Errorf("Expected avatar URL %q", avatar)
	}
	if dto.AvatarThumb == nil || *dto.AvatarThumb != avatarThumb {
		t.Errorf("Expected avatar thumb URL %q", avatarThumb)
	}
	if dto.AvatarMedium == nil || *dto.AvatarMedium != avatarMedium {
		t.Errorf("Expected avatar medium URL %q", avatarMedium)
	}
}

func TestDTOValidation(t *testing.T) {
	t.Run("CreateFounderDTO validation invalid url", func(t *testing.T) {
		companyID := int64(12738)
		linkedin := "not-a-url"
		dto := CreateFounderDTO{
			CompanyID: &companyID,
			FullName:  "John Smith",
			Linkedin:  &linkedin,
		}

		err := validate.Struct(dto)
		if err == nil {
			t.Error("Expected validation error for invalid linkedin url, got nil")
		}
	})

	t.Run("CreateFounderDTO validation valid", func(t *testing.T) {
		companyID := int64(12738)
		linkedin := "https://linkedin.com/in/johnsmith"
		dto := CreateFounderDTO{
			CompanyID: &companyID,
			FullName:  "John Smith",
			Linkedin:  &linkedin,
		}

		err := validate.Struct(dto)
		if err != nil {
			t.Errorf("Expected no validation error, got %v", err)
		}
	})
}
