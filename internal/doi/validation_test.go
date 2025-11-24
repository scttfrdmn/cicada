// Copyright 2025 Scott Friedman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package doi

import (
	"strings"
	"testing"
	"time"
)

func TestDOIReadinessValidator_ValidateComplete(t *testing.T) {
	validator := NewDOIReadinessValidator()

	dataset := &Dataset{
		Title:           "Complete Dataset with All Fields",
		Authors:         []Author{
			{
				Name:        "Dr. Jane Smith",
				GivenName:   "Jane",
				FamilyName:  "Smith",
				ORCID:       "0000-0001-2345-6789",
				Affiliation: "University Lab",
			},
		},
		Publisher:       "Research Institute",
		PublicationYear: time.Now().Year(),
		ResourceType:    "Dataset",
		Description:     "This is a comprehensive description of the dataset that provides sufficient detail about the content, methodology, and purpose of the research data.",
		License:         "CC-BY-4.0",
		URL:             "https://example.com/dataset",
		Version:         "1.0",
		Language:        "en",
		Keywords:        []string{"genomics", "rna-seq", "bioinformatics"},
		Dates: []DateInfo{
			{Date: "2025-01-15", Type: "Collected"},
		},
		RelatedIdentifiers: []RelatedID{
			{Identifier: "10.1234/paper", Type: "DOI", Relation: "IsSupplementTo"},
		},
		Contributors: []Contributor{
			{Name: "Lab Tech", Type: "DataCollector"},
		},
		FundingReferences: []FundingRef{
			{FunderName: "NSF", AwardNumber: "1234567"},
		},
	}

	result := validator.Validate(dataset)

	if !result.IsReady {
		t.Errorf("Complete dataset should be ready, errors: %v", result.Errors)
	}

	if result.Score < 80 {
		t.Errorf("Complete dataset score = %.1f, want >= 80", result.Score)
	}

	if len(result.Errors) > 0 {
		t.Errorf("Complete dataset should have no errors, got: %v", result.Errors)
	}
}

func TestDOIReadinessValidator_ValidateMissingRequired(t *testing.T) {
	validator := NewDOIReadinessValidator()

	tests := []struct {
		name        string
		dataset     *Dataset
		wantError   string
		shouldFail  bool
	}{
		{
			name:       "missing title",
			dataset:    &Dataset{Authors: []Author{{Name: "Author"}}, Publisher: "Pub", PublicationYear: 2025, ResourceType: "Dataset"},
			wantError:  "title is required",
			shouldFail: true,
		},
		{
			name:       "missing authors",
			dataset:    &Dataset{Title: "Title", Publisher: "Pub", PublicationYear: 2025, ResourceType: "Dataset"},
			wantError:  "at least one creator is required",
			shouldFail: true,
		},
		{
			name:       "missing publisher",
			dataset:    &Dataset{Title: "Title", Authors: []Author{{Name: "Author"}}, PublicationYear: 2025, ResourceType: "Dataset"},
			wantError:  "publisher is required",
			shouldFail: true,
		},
		{
			name:       "missing publication year",
			dataset:    &Dataset{Title: "Title", Authors: []Author{{Name: "Author"}}, Publisher: "Pub", ResourceType: "Dataset"},
			wantError:  "publication year is required",
			shouldFail: true,
		},
		{
			name:       "missing resource type",
			dataset:    &Dataset{Title: "Title", Authors: []Author{{Name: "Author"}}, Publisher: "Pub", PublicationYear: 2025},
			wantError:  "resource type is required",
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.Validate(tt.dataset)

			if tt.shouldFail && result.IsReady {
				t.Error("Dataset with missing required field should not be ready")
			}

			hasError := false
			for _, err := range result.Errors {
				if strings.Contains(err, tt.wantError) {
					hasError = true
					break
				}
			}
			if !hasError {
				t.Errorf("Should have error containing '%s', got errors: %v", tt.wantError, result.Errors)
			}
		})
	}
}

func TestDOIReadinessValidator_ValidateUnknownCreator(t *testing.T) {
	validator := NewDOIReadinessValidator()
	validator.RequireRealAuthors = true

	dataset := &Dataset{
		Title:           "Test Dataset",
		Authors:         []Author{{Name: "Unknown Creator"}},
		Publisher:       "Publisher",
		PublicationYear: 2025,
		ResourceType:    "Dataset",
	}

	result := validator.Validate(dataset)

	if result.IsReady {
		t.Error("Dataset with Unknown Creator should not be ready when RequireRealAuthors is true")
	}

	hasUnknownError := false
	for _, err := range result.Errors {
		if strings.Contains(err, "Unknown Creator") {
			hasUnknownError = true
			break
		}
	}
	if !hasUnknownError {
		t.Error("Should have error about Unknown Creator")
	}
}

func TestDOIReadinessValidator_ValidateInvalidYear(t *testing.T) {
	validator := NewDOIReadinessValidator()

	tests := []struct {
		name string
		year int
		want bool // true if should be valid
	}{
		{"too old", 1800, false},
		{"valid old", 1900, true},
		{"current year", time.Now().Year(), true},
		{"next year", time.Now().Year() + 1, true},
		{"future", time.Now().Year() + 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := &Dataset{
				Title:           "Test",
				Authors:         []Author{{Name: "Author", GivenName: "A", FamilyName: "Author"}},
				Publisher:       "Publisher",
				PublicationYear: tt.year,
				ResourceType:    "Dataset",
			}

			result := validator.Validate(dataset)

			hasYearError := false
			for _, err := range result.Errors {
				if strings.Contains(err, "year") || strings.Contains(err, "range") {
					hasYearError = true
					break
				}
			}

			if tt.want && hasYearError {
				t.Errorf("Year %d should be valid but got error", tt.year)
			}
			if !tt.want && !hasYearError {
				t.Errorf("Year %d should be invalid but got no error", tt.year)
			}
		})
	}
}

func TestDOIReadinessValidator_QualityScore(t *testing.T) {
	validator := NewDOIReadinessValidator()

	tests := []struct {
		name      string
		dataset   *Dataset
		minScore  float64
		maxScore  float64
	}{
		{
			name: "minimal valid",
			dataset: &Dataset{
				Title:           "Minimal Dataset",
				Authors:         []Author{{Name: "Author", GivenName: "A", FamilyName: "Author"}},
				Publisher:       "Publisher",
				PublicationYear: 2025,
				ResourceType:    "Dataset",
			},
			minScore: 50,
			maxScore: 65,
		},
		{
			name: "with description",
			dataset: &Dataset{
				Title:           "Dataset with Description",
				Authors:         []Author{{Name: "Author", GivenName: "A", FamilyName: "Author"}},
				Publisher:       "Publisher",
				PublicationYear: 2025,
				ResourceType:    "Dataset",
				Description:     "This is a detailed description of the dataset with sufficient length to meet quality standards.",
			},
			minScore: 60,
			maxScore: 75,
		},
		{
			name: "with rich metadata",
			dataset: &Dataset{
				Title:           "Rich Metadata Dataset",
				Authors:         []Author{{Name: "Dr. Smith", GivenName: "John", FamilyName: "Smith", ORCID: "0000-0001-2345-6789", Affiliation: "University"}},
				Publisher:       "Research Institute",
				PublicationYear: 2025,
				ResourceType:    "Dataset",
				Description:     "Comprehensive description with methodological details and context about the research.",
				License:         "CC-BY-4.0",
				URL:             "https://example.com",
				Keywords:        []string{"keyword1", "keyword2", "keyword3", "keyword4", "keyword5"},
				Version:         "1.0",
			},
			minScore: 85,
			maxScore: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.Validate(tt.dataset)

			if result.Score < tt.minScore || result.Score > tt.maxScore {
				t.Errorf("Score = %.1f, want between %.1f and %.1f", result.Score, tt.minScore, tt.maxScore)
			}
		})
	}
}

func TestDOIReadinessValidator_Warnings(t *testing.T) {
	validator := NewDOIReadinessValidator()
	validator.RequireDescription = false
	validator.RequireLicense = false

	dataset := &Dataset{
		Title:           "Dataset",
		Authors:         []Author{{Name: "Author", GivenName: "A", FamilyName: "Author"}},
		Publisher:       "Unknown Publisher",
		PublicationYear: 2025,
		ResourceType:    "Dataset",
		// Missing description, license, URL, keywords
	}

	result := validator.Validate(dataset)

	if len(result.Warnings) == 0 {
		t.Error("Should have warnings for missing optional fields")
	}

	// Should warn about Unknown Publisher
	hasPublisherWarning := false
	for _, warning := range result.Warnings {
		if strings.Contains(warning, "Unknown Publisher") {
			hasPublisherWarning = true
			break
		}
	}
	if !hasPublisherWarning {
		t.Error("Should warn about Unknown Publisher")
	}
}

func TestDOIReadinessValidator_OptionalFields(t *testing.T) {
	validator := NewDOIReadinessValidator()
	validator.RequireDescription = false

	// Minimal but valid
	minimal := &Dataset{
		Title:           "Minimal",
		Authors:         []Author{{Name: "Author", GivenName: "A", FamilyName: "Author"}},
		Publisher:       "Publisher",
		PublicationYear: 2025,
		ResourceType:    "Dataset",
	}

	minResult := validator.Validate(minimal)

	// With optional fields
	enhanced := &Dataset{
		Title:           "Enhanced",
		Authors:         []Author{{Name: "Author", GivenName: "A", FamilyName: "Author"}},
		Publisher:       "Publisher",
		PublicationYear: 2025,
		ResourceType:    "Dataset",
		Description:     "Detailed description of the dataset with comprehensive information.",
		License:         "CC-BY-4.0",
		Keywords:        []string{"keyword1", "keyword2", "keyword3"},
		Version:         "1.0",
	}

	enhResult := validator.Validate(enhanced)

	// Enhanced should score higher
	if enhResult.Score <= minResult.Score {
		t.Errorf("Enhanced dataset score (%.1f) should be higher than minimal (%.1f)",
			enhResult.Score, minResult.Score)
	}
}

func TestIsValidResourceType(t *testing.T) {
	tests := []struct {
		resourceType string
		want         bool
	}{
		{"Dataset", true},
		{"Software", true},
		{"Image", true},
		{"Text", true},
		{"Other", true},
		{"InvalidType", false},
		{"dataset", false}, // Case sensitive
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.resourceType, func(t *testing.T) {
			got := isValidResourceType(tt.resourceType)
			if got != tt.want {
				t.Errorf("isValidResourceType(%q) = %v, want %v", tt.resourceType, got, tt.want)
			}
		})
	}
}

func TestGetQualityLevel(t *testing.T) {
	tests := []struct {
		score float64
		want  string
	}{
		{100, "Excellent"},
		{85, "Excellent"},
		{80, "Excellent"},
		{75, "Good"},
		{60, "Good"},
		{50, "Moderate"},
		{40, "Moderate"},
		{30, "Poor"},
		{20, "Poor"},
		{10, "Very Poor"},
		{0, "Very Poor"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := GetQualityLevel(tt.score)
			if got != tt.want {
				t.Errorf("GetQualityLevel(%.1f) = %v, want %v", tt.score, got, tt.want)
			}
		})
	}
}

func TestDOIReadinessValidator_GetRecommendations(t *testing.T) {
	validator := NewDOIReadinessValidator()

	tests := []struct {
		name    string
		result  *ReadinessResult
		wantKey string // Key phrase that should appear in recommendations
	}{
		{
			name: "with errors",
			result: &ReadinessResult{
				Score:  30,
				Errors: []string{"title is required"},
			},
			wantKey: "Fix all errors",
		},
		{
			name: "low quality",
			result: &ReadinessResult{
				Score:  35,
				Errors: []string{},
			},
			wantKey: "quality is low",
		},
		{
			name: "moderate quality",
			result: &ReadinessResult{
				Score:  55,
				Errors: []string{},
			},
			wantKey: "quality is moderate",
		},
		{
			name: "good quality",
			result: &ReadinessResult{
				Score:  75,
				Errors: []string{},
			},
			wantKey: "quality is good",
		},
		{
			name: "excellent quality",
			result: &ReadinessResult{
				Score:  90,
				Errors: []string{},
			},
			wantKey: "excellent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recommendations := validator.GetRecommendations(tt.result)
			if len(recommendations) == 0 {
				t.Error("GetRecommendations() should return at least one recommendation")
			}

			// Check that key phrase appears in recommendations
			found := false
			for _, rec := range recommendations {
				if strings.Contains(rec, tt.wantKey) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Recommendations should contain '%s', got: %v", tt.wantKey, recommendations)
			}
		})
	}
}

func TestDOIReadinessValidator_MinQualityThreshold(t *testing.T) {
	validator := NewDOIReadinessValidator()
	validator.MinQualityScore = 70.0

	dataset := &Dataset{
		Title:           "Dataset",
		Authors:         []Author{{Name: "Author", GivenName: "A", FamilyName: "Author"}},
		Publisher:       "Publisher",
		PublicationYear: 2025,
		ResourceType:    "Dataset",
	}

	result := validator.Validate(dataset)

	// Should not be ready even though all required fields present
	if result.IsReady {
		t.Error("Dataset below quality threshold should not be ready")
	}

	// Should have warning about quality score
	hasQualityWarning := false
	for _, warning := range result.Warnings {
		if strings.Contains(warning, "Quality score") && strings.Contains(warning, "threshold") {
			hasQualityWarning = true
			break
		}
	}
	if !hasQualityWarning {
		t.Error("Should have warning about quality score below threshold")
	}
}
