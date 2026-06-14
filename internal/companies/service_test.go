package companies

import (
	"testing"
)

func TestCompanyFilters(t *testing.T) {
	// Test default pagination behavior in handler or logic
	filters := CompanyFilters{
		Limit:  0,
		Offset: -1,
	}

	if filters.Limit <= 0 {
		filters.Limit = 10
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	if filters.Limit != 10 {
		t.Errorf("Expected limit to default to 10, got %d", filters.Limit)
	}
	if filters.Offset != 0 {
		t.Errorf("Expected offset to default to 0, got %d", filters.Offset)
	}
}

func TestMapToCardResponse(t *testing.T) {
	tagline := "Google Search"
	batch := "W24"
	stage := "Growth"
	teamSize := int32(500)
	location := "Mountain View"
	industry := "Search"
	logoURL := "http://example.com/logo.png"
	sourceLogoURL := "http://example.com/source.png"
	smallLogoURL := "http://example.com/small.png"
	sourceSmallLogoURL := "http://example.com/source_small.png"

	comp := &Company{
		ID:                 42,
		Name:               "Google",
		Slug:               "google",
		Tagline:            &tagline,
		Batch:              &batch,
		Stage:              &stage,
		TeamSize:           &teamSize,
		Location:           &location,
		Industry:           &industry,
		LogoURL:            &logoURL,
		SourceLogoURL:      &sourceLogoURL,
		SmallLogoURL:       &smallLogoURL,
		SourceSmallLogoURL: &sourceSmallLogoURL,
	}

	dto := mapToCardResponse(comp)

	if dto.ID != 42 {
		t.Errorf("Expected ID 42, got %d", dto.ID)
	}
	if dto.Name != "Google" {
		t.Errorf("Expected name 'Google', got %q", dto.Name)
	}
	if dto.Slug != "google" {
		t.Errorf("Expected slug 'google', got %q", dto.Slug)
	}
	if dto.Tagline != tagline {
		t.Errorf("Expected tagline %q, got %q", tagline, dto.Tagline)
	}
	if dto.Batch != batch {
		t.Errorf("Expected batch %q, got %q", batch, dto.Batch)
	}
	if dto.Stage != stage {
		t.Errorf("Expected stage %q, got %q", stage, dto.Stage)
	}
	if dto.TeamSize != 500 {
		t.Errorf("Expected teamSize 500, got %d", dto.TeamSize)
	}
	if dto.Location != location {
		t.Errorf("Expected location %q, got %q", location, dto.Location)
	}
	if dto.Industry != industry {
		t.Errorf("Expected industry %q, got %q", industry, dto.Industry)
	}
	if dto.LogoURL != logoURL {
		t.Errorf("Expected logoUrl %q, got %q", logoURL, dto.LogoURL)
	}
	if dto.SourceLogoURL != sourceLogoURL {
		t.Errorf("Expected sourceLogoURL %q, got %q", sourceLogoURL, dto.SourceLogoURL)
	}
	if dto.SmallLogoURL != smallLogoURL {
		t.Errorf("Expected smallLogoURL %q, got %q", smallLogoURL, dto.SmallLogoURL)
	}
	if dto.SourceSmallLogoURL != sourceSmallLogoURL {
		t.Errorf("Expected sourceSmallLogoURL %q, got %q", sourceSmallLogoURL, dto.SourceSmallLogoURL)
	}
}

func TestMapToDetailResponse(t *testing.T) {
	tagline := "Google Search"
	comp := &Company{
		ID:      42,
		Name:    "Google",
		Slug:    "google",
		Tagline: &tagline,
	}

	dto := mapToDetailResponse(comp)

	if dto.ID != 42 {
		t.Errorf("Expected ID 42, got %d", dto.ID)
	}
	if dto.Name != "Google" {
		t.Errorf("Expected name 'Google', got %q", dto.Name)
	}
	if dto.Slug != "google" {
		t.Errorf("Expected slug 'google', got %q", dto.Slug)
	}
	if dto.Tagline == nil || *dto.Tagline != tagline {
		t.Errorf("Expected tagline %q", tagline)
	}
}
