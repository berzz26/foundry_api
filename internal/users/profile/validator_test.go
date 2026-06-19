package profile

import (
	"testing"
)

func TestValidateOnboardingRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     CompleteOnboardingRequest
		wantErr bool
	}{
		{
			name: "Valid request student",
			req: CompleteOnboardingRequest{
				UserType:        "student",
				ExperienceLevel: "student",
				Goals:           []string{"find_jobs", "discover_startups"},
				InterestedRoles: []string{"backend", "ai_ml"},
				PreferredLocations: []string{"remote", "india"},
				CompanyStagePreferences: []string{"seed", "series_a"},
			},
			wantErr: false,
		},
		{
			name: "Valid request software engineer",
			req: CompleteOnboardingRequest{
				UserType:        "software_engineer",
				ExperienceLevel: "3_5_years",
				Goals:           []string{"discover_startups"},
				InterestedRoles: []string{"fullstack"},
				PreferredLocations: []string{"us"},
				CompanyStagePreferences: []string{"series_b_plus"},
			},
			wantErr: false,
		},
		{
			name: "Missing user_type",
			req: CompleteOnboardingRequest{
				ExperienceLevel: "3_5_years",
			},
			wantErr: true,
		},
		{
			name: "Missing experience_level",
			req: CompleteOnboardingRequest{
				UserType: "software_engineer",
			},
			wantErr: true,
		},
		{
			name: "Invalid user_type",
			req: CompleteOnboardingRequest{
				UserType:        "pilot",
				ExperienceLevel: "3_5_years",
			},
			wantErr: true,
		},
		{
			name: "Invalid experience_level",
			req: CompleteOnboardingRequest{
				UserType:        "software_engineer",
				ExperienceLevel: "99_years",
			},
			wantErr: true,
		},
		{
			name: "Invalid goal",
			req: CompleteOnboardingRequest{
				UserType:        "software_engineer",
				ExperienceLevel: "3_5_years",
				Goals:           []string{"world_domination"},
			},
			wantErr: true,
		},
		{
			name: "Invalid interested role",
			req: CompleteOnboardingRequest{
				UserType:        "software_engineer",
				ExperienceLevel: "3_5_years",
				InterestedRoles: []string{"manager"},
			},
			wantErr: true,
		},
		{
			name: "Invalid location",
			req: CompleteOnboardingRequest{
				UserType:        "software_engineer",
				ExperienceLevel: "3_5_years",
				PreferredLocations: []string{"mars"},
			},
			wantErr: true,
		},
		{
			name: "Invalid stage preference",
			req: CompleteOnboardingRequest{
				UserType:        "software_engineer",
				ExperienceLevel: "3_5_years",
				CompanyStagePreferences: []string{"decacorn"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOnboardingRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOnboardingRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
