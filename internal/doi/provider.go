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
	"context"
	"fmt"
)

// Provider interface for pluggable DOI minting services
type Provider interface {
	// Name returns the provider name (datacite, zenodo, disabled, etc.)
	Name() string

	// Mint creates a new DOI for a dataset
	Mint(ctx context.Context, dataset *Dataset) (*DOI, error)

	// Update updates existing DOI metadata
	Update(ctx context.Context, doi string, dataset *Dataset) error

	// Get retrieves DOI information
	Get(ctx context.Context, doi string) (*DOI, error)

	// List returns all DOIs for this provider
	List(ctx context.Context) ([]*DOI, error)

	// Validate checks if dataset metadata is valid for this provider
	Validate(dataset *Dataset) error

	// EstimateCost returns estimated cost for minting (some providers charge)
	// Returns: (cost, currency, error)
	EstimateCost(dataset *Dataset) (float64, string, error)

	// IsEnabled returns true if provider is configured and enabled
	IsEnabled() bool
}

// Dataset represents a dataset to be published with a DOI
type Dataset struct {
	// Required fields
	Title           string   `json:"title"`
	Authors         []Author `json:"authors"`
	Publisher       string   `json:"publisher"`
	PublicationYear int      `json:"publication_year"`
	ResourceType    string   `json:"resource_type"` // Dataset, Software, Image, etc.

	// Optional fields
	Description     string            `json:"description,omitempty"`
	License         string            `json:"license,omitempty"`         // e.g., "CC-BY-4.0"
	URL             string            `json:"url,omitempty"`             // Landing page URL
	Version         string            `json:"version,omitempty"`
	Language        string            `json:"language,omitempty"`        // ISO 639-1 code
	Keywords        []string          `json:"keywords,omitempty"`
	RelatedIdentifiers []RelatedID    `json:"related_identifiers,omitempty"`
	Contributors    []Contributor     `json:"contributors,omitempty"`
	FundingReferences []FundingRef    `json:"funding_references,omitempty"`
	Sizes           []string          `json:"sizes,omitempty"`           // Size information
	Formats         []string          `json:"formats,omitempty"`         // Format information (MIME types, etc.)

	// Dates
	Dates           []DateInfo        `json:"dates,omitempty"`

	// GeoLocation (for datasets with spatial coverage)
	GeoLocations    []GeoLocation     `json:"geo_locations,omitempty"`

	// Custom metadata
	Custom          map[string]interface{} `json:"custom,omitempty"`

	// S3 location (if applicable)
	S3Bucket        string            `json:"s3_bucket,omitempty"`
	S3Prefix        string            `json:"s3_prefix,omitempty"`
}

// RelatedID represents a related identifier (paper, dataset, etc.)
type RelatedID struct {
	Identifier   string `json:"identifier"`     // DOI, URL, etc.
	Type         string `json:"type"`           // DOI, URL, ARK, etc.
	Relation     string `json:"relation"`       // IsSupplementTo, IsCitedBy, etc.
	ResourceType string `json:"resource_type,omitempty"`
}

// Contributor represents a dataset contributor (not author)
type Contributor struct {
	Name         string   `json:"name"`
	GivenName    string   `json:"given_name,omitempty"`
	FamilyName   string   `json:"family_name,omitempty"`
	Type         string   `json:"type"`           // ContactPerson, DataCollector, etc.
	Affiliations []string `json:"affiliations,omitempty"`
	ORCID        string   `json:"orcid,omitempty"`
}

// FundingRef represents funding information
type FundingRef struct {
	FunderName       string `json:"funder_name"`
	FunderIdentifier string `json:"funder_identifier,omitempty"` // CrossRef Funder ID
	AwardNumber      string `json:"award_number,omitempty"`
	AwardTitle       string `json:"award_title,omitempty"`
}

// DateInfo represents a date with type
type DateInfo struct {
	Date string `json:"date"` // ISO 8601 format
	Type string `json:"type"` // Collected, Valid, Accepted, etc.
}

// GeoLocation represents geographic coverage
type GeoLocation struct {
	Place      string      `json:"place,omitempty"`
	Point      *GeoPoint   `json:"point,omitempty"`
	Box        *GeoBox     `json:"box,omitempty"`
	Polygon    []GeoPoint  `json:"polygon,omitempty"`
}

// GeoPoint represents a geographic point
type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// GeoBox represents a geographic bounding box
type GeoBox struct {
	WestLongitude float64 `json:"west_longitude"`
	EastLongitude float64 `json:"east_longitude"`
	SouthLatitude float64 `json:"south_latitude"`
	NorthLatitude float64 `json:"north_latitude"`
}

// ProviderRegistry manages DOI providers
type ProviderRegistry struct {
	providers map[string]Provider
	active    Provider  // Currently active provider
}

// NewProviderRegistry creates a new provider registry
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providers: make(map[string]Provider),
	}
}

// Register registers a provider
func (r *ProviderRegistry) Register(provider Provider) {
	r.providers[provider.Name()] = provider
}

// SetActive sets the active provider by name
func (r *ProviderRegistry) SetActive(name string) error {
	provider, ok := r.providers[name]
	if !ok {
		return fmt.Errorf("provider not found: %s", name)
	}

	if !provider.IsEnabled() {
		return fmt.Errorf("provider is not enabled: %s", name)
	}

	r.active = provider
	return nil
}

// GetActive returns the currently active provider
func (r *ProviderRegistry) GetActive() Provider {
	return r.active
}

// Get returns a provider by name
func (r *ProviderRegistry) Get(name string) (Provider, bool) {
	provider, ok := r.providers[name]
	return provider, ok
}

// List returns all registered providers
func (r *ProviderRegistry) List() []Provider {
	providers := make([]Provider, 0, len(r.providers))
	for _, provider := range r.providers {
		providers = append(providers, provider)
	}
	return providers
}

// DisabledProvider is a no-op provider for when DOI minting is disabled
type DisabledProvider struct{}

// NewDisabledProvider creates a new disabled provider
func NewDisabledProvider() *DisabledProvider {
	return &DisabledProvider{}
}

func (p *DisabledProvider) Name() string {
	return "disabled"
}

func (p *DisabledProvider) Mint(ctx context.Context, dataset *Dataset) (*DOI, error) {
	return nil, fmt.Errorf("DOI minting is disabled")
}

func (p *DisabledProvider) Update(ctx context.Context, doi string, dataset *Dataset) error {
	return fmt.Errorf("DOI minting is disabled")
}

func (p *DisabledProvider) Get(ctx context.Context, doi string) (*DOI, error) {
	return nil, fmt.Errorf("DOI minting is disabled")
}

func (p *DisabledProvider) List(ctx context.Context) ([]*DOI, error) {
	return nil, fmt.Errorf("DOI minting is disabled")
}

func (p *DisabledProvider) Validate(dataset *Dataset) error {
	return fmt.Errorf("DOI minting is disabled")
}

func (p *DisabledProvider) EstimateCost(dataset *Dataset) (float64, string, error) {
	return 0, "USD", nil
}

func (p *DisabledProvider) IsEnabled() bool {
	return true  // Disabled provider is always "enabled" (available)
}

// ConfigureFromFile loads provider configuration from file
func (r *ProviderRegistry) ConfigureFromFile(path string) error {
	// TODO: Implement configuration loading from YAML/JSON
	return fmt.Errorf("not implemented")
}
