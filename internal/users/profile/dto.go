package profile

type CompleteOnboardingRequest struct {
	UserType                string   `json:"user_type"`
	ExperienceLevel         string   `json:"experience_level"`
	Goals                   []string `json:"goals"`
	InterestedRoles         []string `json:"interested_roles"`
	PreferredLocations      []string `json:"preferred_locations"`
	CompanyStagePreferences []string `json:"company_stage_preferences"`
}

type UserProfileResponse struct {
	UserType                string   `json:"user_type"`
	ExperienceLevel         string   `json:"experience_level"`
	Goals                   []string `json:"goals"`
	InterestedRoles         []string `json:"interested_roles"`
	PreferredLocations      []string `json:"preferred_locations"`
	CompanyStagePreferences []string `json:"company_stage_preferences"`
}
