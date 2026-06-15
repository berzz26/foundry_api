package companies

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrCompanyNotFound = errors.New("company not found")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(ctx context.Context, id int64) (*Company, error) {
	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCompanyNotFound
		}
		return nil, err
	}
	return company, nil
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*Company, error) {
	company, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCompanyNotFound
		}
		return nil, err
	}
	return company, nil
}

func (s *Service) List(ctx context.Context, filters CompanyFilters) (*CompanyListResponse, error) {
	list, err := s.repo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx, filters)
	if err != nil {
		return nil, err
	}

	cards := make([]CompanyCardResponse, len(list))
	for i, c := range list {
		cards[i] = mapToCardResponse(&c)
	}

	limit := 10
	if filters.Limit > 0 {
		limit = filters.Limit
	}
	if limit > 100 {
		limit = 100
	}

	offset := 0
	if filters.Offset > 0 {
		offset = filters.Offset
	}

	hasNext := int64(offset+limit) < total

	return &CompanyListResponse{
		Companies: cards,
		Pagination: PaginationResponse{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			HasNext: hasNext,
		},
	}, nil
}

func (s *Service) GetMetadata(ctx context.Context) (*CompanyMetadataResponse, error) {
	return s.repo.GetMetadata(ctx)
}

func mapToCardResponse(c *Company) CompanyCardResponse {
	var tagline string
	if c.Tagline != nil {
		tagline = *c.Tagline
	}
	var batch string
	if c.Batch != nil {
		batch = *c.Batch
	}
	var stage string
	if c.Stage != nil {
		stage = *c.Stage
	}
	var teamSize int
	if c.TeamSize != nil {
		teamSize = int(*c.TeamSize)
	}
	var location string
	if c.Location != nil {
		location = *c.Location
	}
	var industry string
	if c.Industry != nil {
		industry = *c.Industry
	}
	var logoURL string
	if c.LogoURL != nil {
		logoURL = *c.LogoURL
	}
	var sourceLogoURL string
	if c.SourceLogoURL != nil {
		sourceLogoURL = *c.SourceLogoURL
	}
	var smallLogoURL string
	if c.SmallLogoURL != nil {
		smallLogoURL = *c.SmallLogoURL
	}
	var sourceSmallLogoURL string
	if c.SourceSmallLogoURL != nil {
		sourceSmallLogoURL = *c.SourceSmallLogoURL
	}

	return CompanyCardResponse{
		ID:                 c.ID,
		Name:               c.Name,
		Slug:               c.Slug,
		Tagline:            tagline,
		Batch:              batch,
		Stage:              stage,
		TeamSize:           teamSize,
		Location:           location,
		Industry:           industry,
		LogoURL:            logoURL,
		SourceLogoURL:      sourceLogoURL,
		SmallLogoURL:       smallLogoURL,
		SourceSmallLogoURL: sourceSmallLogoURL,
	}
}
