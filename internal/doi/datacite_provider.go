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
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DataCiteProvider implements the Provider interface for DataCite
type DataCiteProvider struct {
	client      *http.Client
	baseURL     string
	repositoryID string
	password    string
	prefix      string // DOI prefix (extracted from repository ID or configured)
	enabled     bool
}

// DataCiteConfig holds DataCite provider configuration
type DataCiteConfig struct {
	RepositoryID string // e.g., "10.5072/FK2" or "CLIENT.MEMBER"
	Password     string
	Prefix       string // Optional: DOI prefix override
	Sandbox      bool   // Use sandbox environment
}

// NewDataCiteProvider creates a new DataCite provider
func NewDataCiteProvider(config *DataCiteConfig) (*DataCiteProvider, error) {
	if config.RepositoryID == "" {
		return nil, fmt.Errorf("DataCite repository ID is required")
	}
	if config.Password == "" {
		return nil, fmt.Errorf("DataCite password is required")
	}

	baseURL := "https://api.datacite.org"
	if config.Sandbox {
		baseURL = "https://api.test.datacite.org"
	}

	// Extract prefix from repository ID if not provided
	prefix := config.Prefix
	if prefix == "" {
		// Repository ID might be like "10.5072/FK2" or "CLIENT.MEMBER"
		// For DOI prefix, we want just the "10.5072" part
		parts := strings.Split(config.RepositoryID, "/")
		if len(parts) > 0 {
			prefix = parts[0]
		} else {
			// If no slash, use the repository ID directly
			prefix = config.RepositoryID
		}
	}

	return &DataCiteProvider{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:      baseURL,
		repositoryID: config.RepositoryID,
		password:     config.Password,
		prefix:       prefix,
		enabled:      true,
	}, nil
}

// Name returns the provider name
func (p *DataCiteProvider) Name() string {
	return "datacite"
}

// doRequestWithRetry wraps an HTTP request with retry logic
//nolint:unused // Will be integrated into existing methods in future refactoring
func (p *DataCiteProvider) doRequestWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var lastErr error

	retryFunc := func() error {
		// Execute the request
		var err error
		resp, err = p.client.Do(req)
		if err != nil {
			// Network error - retryable
			return NewNetworkError(err)
		}

		// Check status code and create appropriate error
		if resp.StatusCode >= 400 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			message := string(bodyBytes)

			// Parse API error response if JSON
			var apiResp struct {
				Errors []struct {
					Title  string `json:"title"`
					Detail string `json:"detail"`
				} `json:"errors"`
			}
			if json.Unmarshal(bodyBytes, &apiResp) == nil && len(apiResp.Errors) > 0 {
				message = apiResp.Errors[0].Detail
				if message == "" {
					message = apiResp.Errors[0].Title
				}
			}

			return NewAPIError(resp.StatusCode, message)
		}

		return nil
	}

	lastErr = WithRetry(ctx, DefaultRetryConfig(), retryFunc)
	if lastErr != nil {
		return nil, lastErr
	}

	return resp, nil
}

// Mint creates a new DOI for a dataset
func (p *DataCiteProvider) Mint(ctx context.Context, dataset *Dataset) (*DOI, error) {
	// Generate DOI suffix
	suffix := generateDOISuffix(dataset)
	doiString := fmt.Sprintf("%s/%s", p.prefix, suffix)

	// Generate DataCite metadata XML
	metadata := generateDataCiteMetadata(dataset)
	metadata.Identifier = Identifier{
		Value: doiString,
		Type:  "DOI",
	}

	xmlData, err := xml.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal metadata: %w", err)
	}

	// Add XML header
	xmlString := xml.Header + string(xmlData)

	// Create DOI via DataCite API
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "dois",
			"attributes": map[string]interface{}{
				"doi":   doiString,
				"url":   dataset.URL,
				"xml":   xmlString,
				"event": "publish", // Creates findable DOI
			},
		},
	}

	response, err := p.createDOI(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("create DOI via API: %w", err)
	}

	// Parse response to extract DOI information
	doi, err := p.parseDataCiteResponse(response)
	if err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return doi, nil
}

// Update updates existing DOI metadata
func (p *DataCiteProvider) Update(ctx context.Context, doiString string, dataset *Dataset) error {
	// Generate updated metadata
	metadata := generateDataCiteMetadata(dataset)
	metadata.Identifier = Identifier{
		Value: doiString,
		Type:  "DOI",
	}

	xmlData, err := xml.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	xmlString := xml.Header + string(xmlData)

	// Update DOI via DataCite API
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "dois",
			"attributes": map[string]interface{}{
				"doi": doiString,
				"url": dataset.URL,
				"xml": xmlString,
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	// URL-encode the DOI
	encodedDOI := url.PathEscape(doiString)
	reqURL := fmt.Sprintf("%s/dois/%s", p.baseURL, encodedDOI)

	req, err := http.NewRequestWithContext(ctx, "PUT", reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.SetBasicAuth(p.repositoryID, p.password)
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Accept", "application/vnd.api+json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// Get retrieves DOI information
func (p *DataCiteProvider) Get(ctx context.Context, doiString string) (*DOI, error) {
	// URL-encode the DOI
	encodedDOI := url.PathEscape(doiString)
	reqURL := fmt.Sprintf("%s/dois/%s", p.baseURL, encodedDOI)

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.SetBasicAuth(p.repositoryID, p.password)
	req.Header.Set("Accept", "application/vnd.api+json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("DOI not found: %s", doiString)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	doi, err := p.parseDataCiteResponse(response)
	if err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return doi, nil
}

// List returns all DOIs for this provider
func (p *DataCiteProvider) List(ctx context.Context) ([]*DOI, error) {
	dois := []*DOI{}
	page := 1
	pageSize := 100

	for {
		reqURL := fmt.Sprintf("%s/dois?client-id=%s&page[size]=%d&page[number]=%d",
			p.baseURL, p.repositoryID, pageSize, page)

		req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}

		req.SetBasicAuth(p.repositoryID, p.password)
		req.Header.Set("Accept", "application/vnd.api+json")

		resp, err := p.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("execute request: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("read response body: %w", err)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			return nil, fmt.Errorf("unmarshal response: %w", err)
		}

		// Parse data array
		data, ok := response["data"].([]interface{})
		if !ok {
			break
		}

		if len(data) == 0 {
			break
		}

		for _, item := range data {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			doi, err := p.parseDataCiteItem(itemMap)
			if err != nil {
				// Log error but continue
				continue
			}

			dois = append(dois, doi)
		}

		// Check if there are more pages
		if len(data) < pageSize {
			break
		}

		page++
	}

	return dois, nil
}

// Validate checks if dataset metadata is valid for DataCite
func (p *DataCiteProvider) Validate(dataset *Dataset) error {
	// Required fields for DataCite
	if dataset.Title == "" {
		return fmt.Errorf("title is required")
	}
	if len(dataset.Authors) == 0 {
		return fmt.Errorf("at least one author is required")
	}
	if dataset.Publisher == "" {
		return fmt.Errorf("publisher is required")
	}
	if dataset.PublicationYear == 0 {
		return fmt.Errorf("publication year is required")
	}
	if dataset.ResourceType == "" {
		return fmt.Errorf("resource type is required")
	}
	if dataset.URL == "" {
		return fmt.Errorf("URL (landing page) is required")
	}

	// Validate publication year is reasonable
	currentYear := time.Now().Year()
	if dataset.PublicationYear < 1900 || dataset.PublicationYear > currentYear+1 {
		return fmt.Errorf("publication year must be between 1900 and %d", currentYear+1)
	}

	return nil
}

// EstimateCost returns estimated cost for minting (DataCite has institutional fees)
func (p *DataCiteProvider) EstimateCost(dataset *Dataset) (float64, string, error) {
	// DataCite costs are typically institutional membership fees, not per-DOI
	// Return 0 since the cost is covered by institution
	return 0.0, "USD", nil
}

// IsEnabled returns true if provider is configured and enabled
func (p *DataCiteProvider) IsEnabled() bool {
	return p.enabled
}

// createDOI creates a DOI via DataCite API
func (p *DataCiteProvider) createDOI(ctx context.Context, payload map[string]interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/dois", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.SetBasicAuth(p.repositoryID, p.password)
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Accept", "application/vnd.api+json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return result, nil
}

// parseDataCiteResponse parses DataCite API response into DOI struct
func (p *DataCiteProvider) parseDataCiteResponse(response map[string]interface{}) (*DOI, error) {
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format: missing data")
	}

	return p.parseDataCiteItem(data)
}

// parseDataCiteItem parses a single DataCite item into DOI struct
func (p *DataCiteProvider) parseDataCiteItem(data map[string]interface{}) (*DOI, error) {
	// Extract attributes
	attributes, ok := data["attributes"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format: missing attributes")
	}

	doi := &DOI{
		DOI:      getStringField(attributes, "doi"),
		URL:      getStringField(attributes, "url"),
		State:    getStringField(attributes, "state"),
		Metadata: attributes,
	}

	// Parse titles
	if titles, ok := attributes["titles"].([]interface{}); ok && len(titles) > 0 {
		if titleMap, ok := titles[0].(map[string]interface{}); ok {
			doi.Title = getStringField(titleMap, "title")
		}
	}

	// Parse creators/authors
	if creators, ok := attributes["creators"].([]interface{}); ok {
		for _, creator := range creators {
			if creatorMap, ok := creator.(map[string]interface{}); ok {
				author := Author{
					Name:       getStringField(creatorMap, "name"),
					GivenName:  getStringField(creatorMap, "givenName"),
					FamilyName: getStringField(creatorMap, "familyName"),
				}

				// Parse name identifiers for ORCID
				if nameIdentifiers, ok := creatorMap["nameIdentifiers"].([]interface{}); ok {
					for _, ni := range nameIdentifiers {
						if niMap, ok := ni.(map[string]interface{}); ok {
							if getStringField(niMap, "nameIdentifierScheme") == "ORCID" {
								author.ORCID = getStringField(niMap, "nameIdentifier")
							}
						}
					}
				}

				// Parse affiliation
				if affiliations, ok := creatorMap["affiliation"].([]interface{}); ok && len(affiliations) > 0 {
					if affMap, ok := affiliations[0].(map[string]interface{}); ok {
						author.Affiliation = getStringField(affMap, "name")
					}
				}

				doi.Authors = append(doi.Authors, author)
			}
		}
	}

	// Parse other fields
	doi.Publisher = getStringField(attributes, "publisher")
	if pubYear, ok := attributes["publicationYear"].(float64); ok {
		doi.PublicationYear = int(pubYear)
	}

	// Parse resource type
	if types, ok := attributes["types"].(map[string]interface{}); ok {
		doi.ResourceType = getStringField(types, "resourceTypeGeneral")
	}

	// Parse descriptions
	if descriptions, ok := attributes["descriptions"].([]interface{}); ok && len(descriptions) > 0 {
		if descMap, ok := descriptions[0].(map[string]interface{}); ok {
			doi.Description = getStringField(descMap, "description")
		}
	}

	// Parse rights/license
	if rightsList, ok := attributes["rightsList"].([]interface{}); ok && len(rightsList) > 0 {
		if rightsMap, ok := rightsList[0].(map[string]interface{}); ok {
			doi.License = getStringField(rightsMap, "rights")
		}
	}

	// Parse dates
	if created, ok := attributes["created"].(string); ok {
		if t, err := time.Parse(time.RFC3339, created); err == nil {
			doi.CreatedAt = t
		}
	}
	if updated, ok := attributes["updated"].(string); ok {
		if t, err := time.Parse(time.RFC3339, updated); err == nil {
			doi.UpdatedAt = t
		}
	}

	return doi, nil
}

// getStringField safely extracts a string field from a map
func getStringField(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

// generateDataCiteMetadata generates DataCite metadata from Dataset
func generateDataCiteMetadata(dataset *Dataset) *DataCiteMetadata {
	metadata := &DataCiteMetadata{
		XMLNS: "http://datacite.org/schema/kernel-4",
		Identifier: Identifier{
			Type: "DOI",
		},
		Publisher:       dataset.Publisher,
		PublicationYear: dataset.PublicationYear,
		ResourceType: ResourceType{
			Value: dataset.ResourceType,
			Type:  "Dataset", // resourceTypeGeneral
		},
	}

	// Titles
	metadata.Titles = []Title{
		{Value: dataset.Title},
	}

	// Creators (authors)
	for _, author := range dataset.Authors {
		creator := Creator{
			CreatorName: author.Name,
			GivenName:   author.GivenName,
			FamilyName:  author.FamilyName,
		}

		if author.ORCID != "" {
			creator.NameIdentifier = &NameIdentifier{
				Value:     author.ORCID,
				Scheme:    "ORCID",
				SchemeURI: "https://orcid.org",
			}
		}

		if author.Affiliation != "" {
			creator.Affiliation = []string{author.Affiliation}
		}

		metadata.Creators = append(metadata.Creators, creator)
	}

	// Subjects (keywords)
	for _, keyword := range dataset.Keywords {
		metadata.Subjects = append(metadata.Subjects, Subject{
			Value: keyword,
		})
	}

	// Descriptions
	if dataset.Description != "" {
		metadata.Descriptions = []Description{
			{
				Value: dataset.Description,
				Type:  "Abstract",
			},
		}
	}

	// Rights (license)
	if dataset.License != "" {
		metadata.RightsList = []Rights{
			{
				Value:     dataset.License,
				RightsURI: getLicenseURI(dataset.License),
			},
		}
	}

	// Sizes and formats
	if len(dataset.Sizes) > 0 {
		metadata.Sizes = dataset.Sizes
	}
	if len(dataset.Formats) > 0 {
		metadata.Formats = dataset.Formats
	}

	// Version
	if dataset.Version != "" {
		metadata.Version = dataset.Version
	}

	// Related identifiers
	for _, relatedID := range dataset.RelatedIdentifiers {
		metadata.RelatedIdentifiers = append(metadata.RelatedIdentifiers, RelatedIdentifier{
			Value:        relatedID.Identifier,
			Type:         relatedID.Type,
			RelationType: relatedID.Relation,
		})
	}

	// Funding
	for _, funding := range dataset.FundingReferences {
		fundingRef := FundingReference{
			FunderName:  funding.FunderName,
			AwardNumber: funding.AwardNumber,
			AwardTitle:  funding.AwardTitle,
		}
		if funding.FunderIdentifier != "" {
			fundingRef.FunderIdentifier = &FunderIdentifier{
				Value: funding.FunderIdentifier,
				Type:  "Crossref Funder ID",
			}
		}
		metadata.FundingReferences = append(metadata.FundingReferences, fundingRef)
	}

	// GeoLocations
	for _, geoLoc := range dataset.GeoLocations {
		dcGeoLoc := DataCiteGeoLocation{
			GeoLocationPlace: geoLoc.Place,
		}
		if geoLoc.Point != nil {
			dcGeoLoc.GeoLocationPoint = &Point{
				PointLongitude: geoLoc.Point.Longitude,
				PointLatitude:  geoLoc.Point.Latitude,
			}
		}
		if geoLoc.Box != nil {
			dcGeoLoc.GeoLocationBox = &Box{
				WestBoundLongitude: geoLoc.Box.WestLongitude,
				EastBoundLongitude: geoLoc.Box.EastLongitude,
				SouthBoundLatitude: geoLoc.Box.SouthLatitude,
				NorthBoundLatitude: geoLoc.Box.NorthLatitude,
			}
		}
		metadata.GeoLocations = append(metadata.GeoLocations, dcGeoLoc)
	}

	// Dates
	for _, dateInfo := range dataset.Dates {
		metadata.Dates = append(metadata.Dates, Date{
			Value: dateInfo.Date,
			Type:  dateInfo.Type,
		})
	}

	// Language
	if dataset.Language != "" {
		metadata.Language = dataset.Language
	}

	// Contributors
	for _, contrib := range dataset.Contributors {
		dcContrib := DataCiteContributor{
			ContributorName: contrib.Name,
			ContributorType: contrib.Type,
			GivenName:       contrib.GivenName,
			FamilyName:      contrib.FamilyName,
		}
		if len(contrib.Affiliations) > 0 {
			dcContrib.Affiliation = contrib.Affiliations
		}
		metadata.Contributors = append(metadata.Contributors, dcContrib)
	}

	return metadata
}
