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
	"fmt"
	"strings"
	"time"
)

// DOIReadinessValidator validates Dataset completeness for DOI minting
type DOIReadinessValidator struct {
	RequireRealAuthors bool // Reject "Unknown Creator"
	RequireDescription bool // Require non-empty description
	RequireLicense     bool // Require license information
	MinQualityScore    float64 // Minimum quality score (0-100)
}

// NewDOIReadinessValidator creates a new validator with default settings
func NewDOIReadinessValidator() *DOIReadinessValidator {
	return &DOIReadinessValidator{
		RequireRealAuthors: true,
		RequireDescription: true,
		RequireLicense:     false,
		MinQualityScore:    60.0,
	}
}

// ReadinessResult represents validation results
type ReadinessResult struct {
	IsReady  bool     `json:"is_ready"`
	Score    float64  `json:"score"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
	Missing  []string `json:"missing,omitempty"`
	Present  []string `json:"present,omitempty"`
}

// Validate validates a Dataset for DOI readiness
func (v *DOIReadinessValidator) Validate(dataset *Dataset) *ReadinessResult {
	result := &ReadinessResult{
		IsReady:  true,
		Errors:   []string{},
		Warnings: []string{},
		Missing:  []string{},
		Present:  []string{},
	}

	// Required fields validation
	v.validateRequired(dataset, result)

	// Optional fields assessment
	v.assessOptional(dataset, result)

	// Calculate quality score
	result.Score = v.calculateQualityScore(dataset, result)

	// Check if meets minimum requirements
	if len(result.Errors) > 0 {
		result.IsReady = false
	}
	if result.Score < v.MinQualityScore {
		result.IsReady = false
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("Quality score %.1f is below minimum threshold %.1f",
				result.Score, v.MinQualityScore))
	}

	return result
}

// validateRequired checks all required fields
func (v *DOIReadinessValidator) validateRequired(dataset *Dataset, result *ReadinessResult) {
	// 1. Identifier - Generated at mint time, not checked here
	result.Present = append(result.Present, "identifier (auto-generated)")

	// 2. Creators (Authors)
	if len(dataset.Authors) == 0 {
		result.Errors = append(result.Errors, "at least one creator is required")
		result.Missing = append(result.Missing, "creators")
	} else {
		hasValid := false
		for _, author := range dataset.Authors {
			if author.Name == "" {
				result.Errors = append(result.Errors, "creator name cannot be empty")
				continue
			}
			if v.RequireRealAuthors && (author.Name == "Unknown Creator" || author.Name == "Unknown") {
				result.Errors = append(result.Errors,
					"creator must be specified (currently set to 'Unknown Creator')")
				result.Missing = append(result.Missing, "real creator names")
				continue
			}
			hasValid = true
		}
		if hasValid {
			result.Present = append(result.Present, "creators")
		}
	}

	// 3. Titles
	if dataset.Title == "" {
		result.Errors = append(result.Errors, "title is required")
		result.Missing = append(result.Missing, "title")
	} else {
		result.Present = append(result.Present, "title")
		// Check for placeholder titles
		if strings.Contains(strings.ToLower(dataset.Title), "untitled") ||
			strings.Contains(strings.ToLower(dataset.Title), "unnamed") {
			result.Warnings = append(result.Warnings,
				"title appears to be a placeholder, consider providing a descriptive title")
		}
	}

	// 4. Publisher
	if dataset.Publisher == "" {
		result.Errors = append(result.Errors, "publisher is required")
		result.Missing = append(result.Missing, "publisher")
	} else {
		result.Present = append(result.Present, "publisher")
		if dataset.Publisher == "Unknown Publisher" {
			result.Warnings = append(result.Warnings,
				"publisher should be specified (currently set to 'Unknown Publisher')")
		}
	}

	// 5. PublicationYear
	currentYear := time.Now().Year()
	if dataset.PublicationYear == 0 {
		result.Errors = append(result.Errors, "publication year is required")
		result.Missing = append(result.Missing, "publication_year")
	} else if dataset.PublicationYear < 1900 || dataset.PublicationYear > currentYear+1 {
		result.Errors = append(result.Errors,
			fmt.Sprintf("publication year %d is outside valid range (1900-%d)",
				dataset.PublicationYear, currentYear+1))
	} else {
		result.Present = append(result.Present, "publication_year")
	}

	// 6. ResourceType
	if dataset.ResourceType == "" {
		result.Errors = append(result.Errors, "resource type is required")
		result.Missing = append(result.Missing, "resource_type")
	} else {
		result.Present = append(result.Present, "resource_type")
		// Validate against DataCite controlled vocabulary
		if !isValidResourceType(dataset.ResourceType) {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("resource type '%s' may not be in DataCite controlled vocabulary",
					dataset.ResourceType))
		}
	}
}

// assessOptional assesses optional but recommended fields
func (v *DOIReadinessValidator) assessOptional(dataset *Dataset, result *ReadinessResult) {
	// Description
	if dataset.Description == "" {
		result.Missing = append(result.Missing, "description")
		if v.RequireDescription {
			result.Errors = append(result.Errors, "description is required")
		} else {
			result.Warnings = append(result.Warnings, "description is recommended for better discoverability")
		}
	} else {
		result.Present = append(result.Present, "description")
		if len(dataset.Description) < 50 {
			result.Warnings = append(result.Warnings,
				"description is very short, consider adding more detail")
		}
	}

	// License/Rights
	if dataset.License == "" {
		result.Missing = append(result.Missing, "license")
		if v.RequireLicense {
			result.Errors = append(result.Errors, "license is required")
		} else {
			result.Warnings = append(result.Warnings, "license/rights information is recommended")
		}
	} else {
		result.Present = append(result.Present, "license")
	}

	// URL (landing page)
	if dataset.URL == "" {
		result.Missing = append(result.Missing, "url")
		result.Warnings = append(result.Warnings, "landing page URL is recommended")
	} else {
		result.Present = append(result.Present, "url")
		if !strings.HasPrefix(dataset.URL, "http://") && !strings.HasPrefix(dataset.URL, "https://") {
			result.Warnings = append(result.Warnings, "URL should start with http:// or https://")
		}
	}

	// Keywords/Subjects
	if len(dataset.Keywords) == 0 {
		result.Missing = append(result.Missing, "keywords")
		result.Warnings = append(result.Warnings, "keywords/subjects are recommended for discoverability")
	} else {
		result.Present = append(result.Present, "keywords")
	}

	// Version
	if dataset.Version == "" {
		result.Missing = append(result.Missing, "version")
	} else {
		result.Present = append(result.Present, "version")
	}

	// Language
	if dataset.Language == "" {
		result.Missing = append(result.Missing, "language")
	} else {
		result.Present = append(result.Present, "language")
	}

	// Dates
	if len(dataset.Dates) == 0 {
		result.Missing = append(result.Missing, "dates")
	} else {
		result.Present = append(result.Present, "dates")
	}

	// Related identifiers
	if len(dataset.RelatedIdentifiers) == 0 {
		result.Missing = append(result.Missing, "related_identifiers")
	} else {
		result.Present = append(result.Present, "related_identifiers")
	}

	// Contributors
	if len(dataset.Contributors) == 0 {
		result.Missing = append(result.Missing, "contributors")
	} else {
		result.Present = append(result.Present, "contributors")
	}

	// Funding
	if len(dataset.FundingReferences) == 0 {
		result.Missing = append(result.Missing, "funding_references")
	} else {
		result.Present = append(result.Present, "funding_references")
	}

	// GeoLocation
	if len(dataset.GeoLocations) == 0 {
		result.Missing = append(result.Missing, "geo_locations")
	} else {
		result.Present = append(result.Present, "geo_locations")
	}

	// Check author details
	hasORCID := false
	hasAffiliation := false
	for _, author := range dataset.Authors {
		if author.ORCID != "" {
			hasORCID = true
		}
		if author.Affiliation != "" {
			hasAffiliation = true
		}
	}
	if !hasORCID {
		result.Warnings = append(result.Warnings, "author ORCIDs are recommended for attribution")
	}
	if !hasAffiliation {
		result.Warnings = append(result.Warnings, "author affiliations are recommended")
	}
}

// calculateQualityScore calculates a quality score (0-100)
func (v *DOIReadinessValidator) calculateQualityScore(dataset *Dataset, result *ReadinessResult) float64 {
	score := 0.0

	// Required fields: 60 points total (all-or-nothing per category)
	requiredScore := 0.0
	hasRequiredCreators := len(dataset.Authors) > 0 && !v.hasUnknownCreator(dataset)
	hasRequiredTitle := dataset.Title != ""
	hasRequiredPublisher := dataset.Publisher != "" && dataset.Publisher != "Unknown Publisher"
	hasRequiredYear := dataset.PublicationYear >= 1900 && dataset.PublicationYear <= time.Now().Year()+1
	hasRequiredType := dataset.ResourceType != ""

	if hasRequiredCreators {
		requiredScore += 15.0
	}
	if hasRequiredTitle {
		requiredScore += 10.0
	}
	if hasRequiredPublisher {
		requiredScore += 10.0
	}
	if hasRequiredYear {
		requiredScore += 10.0
	}
	if hasRequiredType {
		requiredScore += 15.0
	}
	score += requiredScore

	// Recommended fields: 40 points total
	if dataset.Description != "" && len(dataset.Description) >= 50 {
		score += 10.0
	} else if dataset.Description != "" {
		score += 5.0
	}

	if dataset.License != "" {
		score += 5.0
	}

	if dataset.URL != "" && (strings.HasPrefix(dataset.URL, "http://") || strings.HasPrefix(dataset.URL, "https://")) {
		score += 5.0
	}

	if len(dataset.Keywords) > 0 {
		score += 5.0
		if len(dataset.Keywords) >= 5 {
			score += 2.0
		}
	}

	if len(dataset.RelatedIdentifiers) > 0 {
		score += 3.0
	}

	if len(dataset.Dates) > 0 {
		score += 2.0
	}

	if len(dataset.Contributors) > 0 {
		score += 2.0
	}

	if len(dataset.FundingReferences) > 0 {
		score += 3.0
	}

	if dataset.Version != "" {
		score += 2.0
	}

	if dataset.Language != "" {
		score += 1.0
	}

	// Author quality bonuses
	hasORCID := false
	hasAffiliation := false
	hasFullNames := true
	for _, author := range dataset.Authors {
		if author.ORCID != "" {
			hasORCID = true
		}
		if author.Affiliation != "" {
			hasAffiliation = true
		}
		if author.GivenName == "" || author.FamilyName == "" {
			hasFullNames = false
		}
	}
	if hasORCID {
		score += 3.0
	}
	if hasAffiliation {
		score += 2.0
	}
	if hasFullNames {
		score += 2.0
	}

	// Penalty for errors
	errorPenalty := float64(len(result.Errors)) * 10.0
	score -= errorPenalty
	if score < 0 {
		score = 0
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score
}

// hasUnknownCreator checks if dataset has placeholder creator
func (v *DOIReadinessValidator) hasUnknownCreator(dataset *Dataset) bool {
	for _, author := range dataset.Authors {
		if author.Name == "Unknown Creator" || author.Name == "Unknown" {
			return true
		}
	}
	return false
}

// isValidResourceType checks if resource type is in DataCite vocabulary
func isValidResourceType(rt string) bool {
	validTypes := map[string]bool{
		"Audiovisual":            true,
		"Book":                   true,
		"BookChapter":            true,
		"Collection":             true,
		"ComputationalNotebook":  true,
		"ConferencePaper":        true,
		"ConferenceProceeding":   true,
		"DataPaper":              true,
		"Dataset":                true,
		"Dissertation":           true,
		"Event":                  true,
		"Image":                  true,
		"Instrument":             true,
		"InteractiveResource":    true,
		"Journal":                true,
		"JournalArticle":         true,
		"Model":                  true,
		"OutputManagementPlan":   true,
		"PeerReview":             true,
		"PhysicalObject":         true,
		"Preprint":               true,
		"Report":                 true,
		"Software":               true,
		"Sound":                  true,
		"Standard":               true,
		"StudyRegistration":      true,
		"Text":                   true,
		"Workflow":               true,
		"Other":                  true,
	}
	return validTypes[rt]
}

// GetRecommendations returns specific recommendations for improving metadata
func (v *DOIReadinessValidator) GetRecommendations(result *ReadinessResult) []string {
	recommendations := []string{}

	// Critical issues first
	if len(result.Errors) > 0 {
		recommendations = append(recommendations,
			"Fix all errors before minting DOI:")
		for _, err := range result.Errors {
			recommendations = append(recommendations, "  - "+err)
		}
	}

	// Score-based recommendations
	if result.Score < 40 {
		recommendations = append(recommendations,
			"Metadata quality is low. Consider adding:",
			"  - Complete author information (names, ORCIDs, affiliations)",
			"  - Detailed description (at least 100 words)",
			"  - Keywords for discoverability",
			"  - License information",
			"  - Related identifiers (publications, datasets)")
	} else if result.Score < 60 {
		recommendations = append(recommendations,
			"Metadata quality is moderate. To improve:",
			"  - Add missing recommended fields",
			"  - Enhance description with methods and context",
			"  - Include author ORCIDs for proper attribution")
	} else if result.Score < 80 {
		recommendations = append(recommendations,
			"Metadata quality is good. Optional improvements:",
			"  - Add funding information",
			"  - Include temporal/spatial context if applicable",
			"  - Link to related publications")
	} else {
		recommendations = append(recommendations,
			"Metadata quality is excellent! Ready for DOI minting.")
	}

	return recommendations
}

// GetQualityLevel returns a human-readable quality level
func GetQualityLevel(score float64) string {
	switch {
	case score >= 80:
		return "Excellent"
	case score >= 60:
		return "Good"
	case score >= 40:
		return "Moderate"
	case score >= 20:
		return "Poor"
	default:
		return "Very Poor"
	}
}
