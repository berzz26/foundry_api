package profile

import (
	"errors"
)

var (
	AllowedUserTypes = map[string]bool{
		"student":           true,
		"software_engineer": true,
		"founder":           true,
		"recruiter":         true,
		"investor":          true,
		"job_seeker":        true,
		"other":             true,
	}

	AllowedGoals = map[string]bool{
		"find_jobs":             true,
		"discover_startups":     true,
		"research_companies":    true,
		"find_founders":         true,
		"track_hiring_signals":  true,
		"fundraising_research": true,
		"recruiting":           true,
	}

	AllowedInterestedRoles = map[string]bool{
		"backend":   true,
		"frontend":  true,
		"fullstack": true,
		"ai_ml":     true,
		"infra":     true,
		"devops":    true,
		"security":  true,
		"product":   true,
		"design":    true,
	}

	AllowedExperienceLevels = map[string]bool{
		"student":       true,
		"0_2_years":     true,
		"3_5_years":     true,
		"5_10_years":    true,
		"10_plus_years": true,
	}

	AllowedPreferredLocations = map[string]bool{
		"remote":    true,
		"india":     true,
		"us":        true,
		"europe":    true,
		"singapore": true,
		"uae":       true,
	}

	AllowedCompanyStagePreferences = map[string]bool{
		"idea":          true,
		"pre_seed":      true,
		"seed":          true,
		"series_a":      true,
		"series_b_plus": true,
		"public":        true,
	}
)

func ValidateOnboardingRequest(req CompleteOnboardingRequest) error {
	if req.UserType == "" {
		return errors.New("user_type is required")
	}
	if !AllowedUserTypes[req.UserType] {
		return errors.New("invalid user_type value")
	}

	if req.ExperienceLevel == "" {
		return errors.New("experience_level is required")
	}
	if !AllowedExperienceLevels[req.ExperienceLevel] {
		return errors.New("invalid experience_level value")
	}

	for _, g := range req.Goals {
		if !AllowedGoals[g] {
			return errors.New("invalid goals value: " + g)
		}
	}

	for _, r := range req.InterestedRoles {
		if !AllowedInterestedRoles[r] {
			return errors.New("invalid interested_roles value: " + r)
		}
	}

	for _, l := range req.PreferredLocations {
		if !AllowedPreferredLocations[l] {
			return errors.New("invalid preferred_locations value: " + l)
		}
	}

	for _, s := range req.CompanyStagePreferences {
		if !AllowedCompanyStagePreferences[s] {
			return errors.New("invalid company_stage_preferences value: " + s)
		}
	}

	return nil
}
