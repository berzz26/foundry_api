package jobs

import (
	"context"
	"encoding/json"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func mapToJobCardResponse(j JobWithCompany) JobCardResponse {
	remote := false
	if j.Remote != nil && (*j.Remote == "yes" || *j.Remote == "only") {
		remote = true
	}

	var salary *SalaryResponse
	if j.SalaryMin != nil || j.SalaryMax != nil {
		salary = &SalaryResponse{
			Min: j.SalaryMin,
			Max: j.SalaryMax,
		}
	}

	var equity *EquityResponse
	if j.EquityMin != nil || j.EquityMax != nil {
		equity = &EquityResponse{
			Min: j.EquityMin,
			Max: j.EquityMax,
		}
	}

	var experience *ExperienceResponse
	if j.MinExperience != nil {
		experience = &ExperienceResponse{
			MinYears: j.MinExperience,
		}
	}

	company := CompanyPreviewResponse{}
	if j.CompanyID != nil {
		company.ID = *j.CompanyID
	}
	if j.CompanyName != nil {
		company.Name = *j.CompanyName
	}
	if j.CompanyLogoSource != nil {
		company.LogoURL = j.CompanyLogoSource
	} else if j.CompanyLogo != nil {
		company.LogoURL = j.CompanyLogo
	}
	if j.CompanyBatch != nil {
		company.Batch = j.CompanyBatch
	}

	return JobCardResponse{
		ID:               j.ID,
		Title:            j.Title,
		Location:         j.Location,
		Remote:           remote,
		Salary:           salary,
		Equity:           equity,
		Experience:       experience,
		Company:          company,
		InterviewProcess: j.InterviewProcess,
		PrettyEngType:    j.PrettyEngType,
		CreatedAt:        j.CreatedAt,
	}
}

func mapToJobDetailResponse(j JobWithCompany) JobDetailResponse {
	remote := false
	if j.Remote != nil && (*j.Remote == "yes" || *j.Remote == "only") {
		remote = true
	}

	visaRequired := false
	if j.VisaRequired != nil && (*j.VisaRequired == "yes" || *j.VisaRequired == "possible") {
		visaRequired = true
	}

	var salary *SalaryResponse
	if j.SalaryMin != nil || j.SalaryMax != nil {
		salary = &SalaryResponse{
			Min: j.SalaryMin,
			Max: j.SalaryMax,
		}
	}

	var equity *EquityResponse
	if j.EquityMin != nil || j.EquityMax != nil {
		equity = &EquityResponse{
			Min: j.EquityMin,
			Max: j.EquityMax,
		}
	}

	var experience *ExperienceResponse
	if j.MinExperience != nil {
		experience = &ExperienceResponse{
			MinYears: j.MinExperience,
		}
	}

	company := CompanyPreviewResponse{}
	if j.CompanyID != nil {
		company.ID = *j.CompanyID
	}
	if j.CompanyName != nil {
		company.Name = *j.CompanyName
	}
	if j.CompanyLogoSource != nil {
		company.LogoURL = j.CompanyLogoSource
	} else if j.CompanyLogo != nil {
		company.LogoURL = j.CompanyLogo
	}
	if j.CompanyBatch != nil {
		company.Batch = j.CompanyBatch
	}

	var skills []string
	if j.Skills != nil && *j.Skills != "" && *j.Skills != "[]" {
		json.Unmarshal([]byte(*j.Skills), &skills)
	}
	if skills == nil {
		skills = []string{}
	}

	return JobDetailResponse{
		ID:               j.ID,
		Title:            j.Title,
		Description:      j.Description,
		Location:         j.Location,
		Remote:           remote,
		Salary:           salary,
		Equity:           equity,
		Experience:       experience,
		VisaRequired:     visaRequired,
		Skills:           skills,
		InterviewProcess: j.InterviewProcess,
		PrettyEngType:    j.PrettyEngType,
		TimeToHire:       j.TimeToHire,
		Company:          company,
		CreatedAt:        j.CreatedAt,
		UpdatedAt:        j.UpdatedAt,
	}
}

func (s *Service) List(ctx context.Context, filters JobFilters) (*JobListResponse, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 20
	}

	jobs, total, err := s.repo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	cards := make([]JobCardResponse, 0, len(jobs))
	for _, j := range jobs {
		cards = append(cards, mapToJobCardResponse(j))
	}

	hasNext := (filters.Page * filters.Limit) < total

	return &JobListResponse{
		Jobs: cards,
		Pagination: PaginationResponse{
			Page:    filters.Page,
			Limit:   filters.Limit,
			Total:   total,
			HasNext: hasNext,
		},
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*JobDetailResponse, error) {
	job, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := mapToJobDetailResponse(*job)
	return &res, nil
}

func (s *Service) GetRelated(ctx context.Context, id int64) (*JobRelatedResponse, error) {
	job, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	role := ""
	if job.Role != nil {
		role = *job.Role
	}

	jobs, err := s.repo.GetRelated(ctx, job.ID, role, job.CompanyID, 10)
	if err != nil {
		return nil, err
	}

	cards := make([]JobCardResponse, 0, len(jobs))
	for _, j := range jobs {
		cards = append(cards, mapToJobCardResponse(j))
	}

	return &JobRelatedResponse{Jobs: cards}, nil
}

func (s *Service) GetFeatured(ctx context.Context) (*JobFeaturedResponse, error) {
	jobs, err := s.repo.GetFeatured(ctx, 10)
	if err != nil {
		return nil, err
	}

	cards := make([]JobCardResponse, 0, len(jobs))
	for _, j := range jobs {
		cards = append(cards, mapToJobCardResponse(j))
	}

	return &JobFeaturedResponse{Jobs: cards}, nil
}

func (s *Service) GetRandomJobs(ctx context.Context, limit int) (*JobListResponse, error) {
	jobs, err := s.repo.GetRandomJobs(ctx, limit)
	if err != nil {
		return nil, err
	}

	cards := make([]JobCardResponse, 0, len(jobs))
	for _, j := range jobs {
		cards = append(cards, mapToJobCardResponse(j))
	}

	return &JobListResponse{
		Jobs: cards,
		Pagination: PaginationResponse{
			Page:    1,
			Limit:   limit,
			Total:   limit,
			HasNext: false,
		},
	}, nil
}
