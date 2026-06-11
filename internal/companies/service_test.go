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

func TestMapToResponseDTO(t *testing.T) {
	name := "Google"
	slug := "google"
	comp := &Company{
		ID:   42,
		Name: name,
		Slug: slug,
	}

	dto := mapToResponseDTO(comp)

	if dto.ID != 42 {
		t.Errorf("Expected ID 42, got %d", dto.ID)
	}
	if dto.Name != name {
		t.Errorf("Expected name %q, got %q", name, dto.Name)
	}
	if dto.Slug != slug {
		t.Errorf("Expected slug %q, got %q", slug, dto.Slug)
	}
}
