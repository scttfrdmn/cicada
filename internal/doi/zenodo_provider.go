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
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ZenodoProvider implements the Provider interface for Zenodo
type ZenodoProvider struct {
	client  *http.Client
	baseURL string
	token   string
	enabled bool
}

// ZenodoConfig holds Zenodo provider configuration
type ZenodoConfig struct {
	Token   string // API access token
	Sandbox bool   // Use sandbox environment
}

// ZenodoDeposition represents a Zenodo deposition (draft DOI)
type ZenodoDeposition struct {
	ID                int                    `json:"id"`
	DOI               string                 `json:"doi"`
	DOIPreview        string                 `json:"doi_url"`
	ConceptDOI        string                 `json:"conceptdoi"`
	State             string                 `json:"state"` // draft, published
	Title             string                 `json:"title"`
	Created           string                 `json:"created"`
	Modified          string                 `json:"modified"`
	Metadata          ZenodoMetadata         `json:"metadata"`
	Links             map[string]string      `json:"links"`
	Files             []ZenodoFile           `json:"files"`
}

// ZenodoMetadata represents Zenodo metadata
type ZenodoMetadata struct {
	Title              string                `json:"title"`
	UploadType         string                `json:"upload_type"` // dataset, software, etc.
	Description        string                `json:"description"`
	Creators           []ZenodoCreator       `json:"creators"`
	PublicationDate    string                `json:"publication_date,omitempty"`
	Keywords           []string              `json:"keywords,omitempty"`
	License            string                `json:"license,omitempty"`
	AccessRight        string                `json:"access_right"` // open, embargoed, restricted, closed
	Version            string                `json:"version,omitempty"`
	RelatedIdentifiers []ZenodoRelatedID     `json:"related_identifiers,omitempty"`
	Contributors       []ZenodoContributor   `json:"contributors,omitempty"`
	References         []string              `json:"references,omitempty"`
	Communities        []ZenodoCommunity     `json:"communities,omitempty"`
	Grants             []ZenodoGrant         `json:"grants,omitempty"`
}

// ZenodoCreator represents a creator/author
type ZenodoCreator struct {
	Name        string `json:"name"`
	Affiliation string `json:"affiliation,omitempty"`
	ORCID       string `json:"orcid,omitempty"`
}

// ZenodoRelatedID represents a related identifier
type ZenodoRelatedID struct {
	Identifier string `json:"identifier"`
	Relation   string `json:"relation"` // isSupplementTo, isCitedBy, etc.
}

// ZenodoContributor represents a contributor
type ZenodoContributor struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // ContactPerson, DataCollector, etc.
	Affiliation string `json:"affiliation,omitempty"`
	ORCID       string `json:"orcid,omitempty"`
}

// ZenodoCommunity represents a Zenodo community
type ZenodoCommunity struct {
	Identifier string `json:"identifier"`
}

// ZenodoGrant represents funding information
type ZenodoGrant struct {
	ID string `json:"id"` // Grant ID from Zenodo's grant registry
}

// ZenodoFile represents an uploaded file
type ZenodoFile struct {
	ID       string            `json:"id"`
	Filename string            `json:"filename"`
	Filesize int64             `json:"filesize"`
	Checksum string            `json:"checksum"`
	Links    map[string]string `json:"links"`
}

// NewZenodoProvider creates a new Zenodo provider
func NewZenodoProvider(config *ZenodoConfig) (*ZenodoProvider, error) {
	if config.Token == "" {
		return nil, fmt.Errorf("zenodo token is required")
	}

	baseURL := "https://zenodo.org/api"
	if config.Sandbox {
		baseURL = "https://sandbox.zenodo.org/api"
	}

	return &ZenodoProvider{
		client: &http.Client{
			Timeout: 60 * time.Second, // Longer timeout for file uploads
		},
		baseURL: baseURL,
		token:   config.Token,
		enabled: true,
	}, nil
}

// Name returns the provider name
func (p *ZenodoProvider) Name() string {
	return "zenodo"
}

// doRequestWithRetry wraps an HTTP request with retry logic
//nolint:unused // Will be integrated into existing methods in future refactoring
func (p *ZenodoProvider) doRequestWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
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
				Message string `json:"message"`
				Status  int    `json:"status"`
				Errors  []struct {
					Field    string `json:"field"`
					Messages string `json:"message"`
				} `json:"errors"`
			}
			if json.Unmarshal(bodyBytes, &apiResp) == nil {
				if apiResp.Message != "" {
					message = apiResp.Message
				}
				if len(apiResp.Errors) > 0 {
					message = fmt.Sprintf("%s: %s", apiResp.Errors[0].Field, apiResp.Errors[0].Messages)
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
func (p *ZenodoProvider) Mint(ctx context.Context, dataset *Dataset) (*DOI, error) {
	// Step 1: Create a new deposition (draft)
	deposition, err := p.createDeposition(ctx, dataset)
	if err != nil {
		return nil, fmt.Errorf("create deposition: %w", err)
	}

	// Step 2: Upload files if any are specified
	// Note: For now, we'll skip file upload and just publish metadata
	// File upload will be implemented when we have file paths in Dataset
	// TODO: Add file upload support

	// Step 3: Publish the deposition (mints the DOI)
	published, err := p.publishDeposition(ctx, deposition.ID)
	if err != nil {
		return nil, fmt.Errorf("publish deposition: %w", err)
	}

	// Step 4: Convert to DOI struct
	doi := p.depositionToDOI(published)

	return doi, nil
}

// Update updates existing DOI metadata
func (p *ZenodoProvider) Update(ctx context.Context, doiString string, dataset *Dataset) error {
	// Zenodo requires creating a new version for updates
	// For now, return error indicating this is not supported
	// TODO: Implement new version creation workflow
	return fmt.Errorf("zenodo updates require creating a new version (not yet implemented)")
}

// Get retrieves DOI information
func (p *ZenodoProvider) Get(ctx context.Context, doiString string) (*DOI, error) {
	// Extract record ID from DOI
	// Zenodo DOIs are like: 10.5281/zenodo.123456
	// We need the numeric ID (123456)
	recordID, err := p.extractRecordID(doiString)
	if err != nil {
		return nil, fmt.Errorf("extract record ID: %w", err)
	}

	// Get deposition by ID
	reqURL := fmt.Sprintf("%s/deposit/depositions/%d", p.baseURL, recordID)

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.token))
	req.Header.Set("Content-Type", "application/json")

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

	var deposition ZenodoDeposition
	if err := json.NewDecoder(resp.Body).Decode(&deposition); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	doi := p.depositionToDOI(&deposition)
	return doi, nil
}

// List returns all DOIs for this provider
func (p *ZenodoProvider) List(ctx context.Context) ([]*DOI, error) {
	dois := []*DOI{}
	page := 1
	pageSize := 100

	for {
		reqURL := fmt.Sprintf("%s/deposit/depositions?page=%d&size=%d", p.baseURL, page, pageSize)

		req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.token))
		req.Header.Set("Content-Type", "application/json")

		resp, err := p.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("execute request: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		}

		var depositions []ZenodoDeposition
		if err := json.NewDecoder(resp.Body).Decode(&depositions); err != nil {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("decode response: %w", err)
		}
		_ = resp.Body.Close()

		if len(depositions) == 0 {
			break
		}

		for _, dep := range depositions {
			doi := p.depositionToDOI(&dep)
			dois = append(dois, doi)
		}

		// Check if there are more pages
		if len(depositions) < pageSize {
			break
		}

		page++
	}

	return dois, nil
}

// Validate checks if dataset metadata is valid for Zenodo
func (p *ZenodoProvider) Validate(dataset *Dataset) error {
	// Required fields for Zenodo
	if dataset.Title == "" {
		return fmt.Errorf("title is required")
	}
	if len(dataset.Authors) == 0 {
		return fmt.Errorf("at least one author is required")
	}
	if dataset.Description == "" {
		return fmt.Errorf("description is required")
	}
	// Zenodo doesn't require publisher, it sets itself as publisher
	// Publication year is optional, defaults to current year

	return nil
}

// EstimateCost returns estimated cost for minting (Zenodo is free)
func (p *ZenodoProvider) EstimateCost(dataset *Dataset) (float64, string, error) {
	// Zenodo is completely free
	return 0.0, "USD", nil
}

// IsEnabled returns true if provider is configured and enabled
func (p *ZenodoProvider) IsEnabled() bool {
	return p.enabled
}

// createDeposition creates a new Zenodo deposition
func (p *ZenodoProvider) createDeposition(ctx context.Context, dataset *Dataset) (*ZenodoDeposition, error) {
	// Convert Dataset to Zenodo metadata
	metadata := p.datasetToZenodoMetadata(dataset)

	// Create deposition request
	payload := map[string]interface{}{
		"metadata": metadata,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	reqURL := fmt.Sprintf("%s/deposit/depositions", p.baseURL)

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.token))
	req.Header.Set("Content-Type", "application/json")

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

	var deposition ZenodoDeposition
	if err := json.Unmarshal(bodyBytes, &deposition); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &deposition, nil
}

// publishDeposition publishes a Zenodo deposition (mints DOI)
func (p *ZenodoProvider) publishDeposition(ctx context.Context, depositionID int) (*ZenodoDeposition, error) {
	reqURL := fmt.Sprintf("%s/deposit/depositions/%d/actions/publish", p.baseURL, depositionID)

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var deposition ZenodoDeposition
	if err := json.Unmarshal(bodyBytes, &deposition); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &deposition, nil
}

// uploadFile uploads a file to a Zenodo deposition
//nolint:unused // Reserved for future file upload functionality
func (p *ZenodoProvider) uploadFile(ctx context.Context, depositionID int, filePath string) (*ZenodoFile, error) {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("copy file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close writer: %w", err)
	}

	// Upload file
	reqURL := fmt.Sprintf("%s/deposit/depositions/%d/files", p.baseURL, depositionID)

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.token))
	req.Header.Set("Content-Type", writer.FormDataContentType())

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

	var zenodoFile ZenodoFile
	if err := json.Unmarshal(bodyBytes, &zenodoFile); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &zenodoFile, nil
}

// datasetToZenodoMetadata converts Dataset to Zenodo metadata
func (p *ZenodoProvider) datasetToZenodoMetadata(dataset *Dataset) ZenodoMetadata {
	metadata := ZenodoMetadata{
		Title:       dataset.Title,
		UploadType:  "dataset", // Default to dataset
		Description: dataset.Description,
		AccessRight: "open",    // Default to open access
	}

	// Creators/authors
	for _, author := range dataset.Authors {
		creator := ZenodoCreator{
			Name:        author.Name,
			Affiliation: author.Affiliation,
			ORCID:       author.ORCID,
		}
		metadata.Creators = append(metadata.Creators, creator)
	}

	// Publication date
	if dataset.PublicationYear > 0 {
		metadata.PublicationDate = fmt.Sprintf("%d-01-01", dataset.PublicationYear)
	} else {
		metadata.PublicationDate = time.Now().Format("2006-01-02")
	}

	// Keywords
	metadata.Keywords = dataset.Keywords

	// License
	if dataset.License != "" {
		metadata.License = mapLicenseToZenodo(dataset.License)
	}

	// Version
	metadata.Version = dataset.Version

	// Related identifiers
	for _, relatedID := range dataset.RelatedIdentifiers {
		zenodoRelated := ZenodoRelatedID{
			Identifier: relatedID.Identifier,
			Relation:   mapRelationToZenodo(relatedID.Relation),
		}
		metadata.RelatedIdentifiers = append(metadata.RelatedIdentifiers, zenodoRelated)
	}

	// Contributors
	for _, contrib := range dataset.Contributors {
		zenodoContrib := ZenodoContributor{
			Name:        contrib.Name,
			Type:        mapContributorTypeToZenodo(contrib.Type),
			ORCID:       contrib.ORCID,
		}
		if len(contrib.Affiliations) > 0 {
			zenodoContrib.Affiliation = contrib.Affiliations[0]
		}
		metadata.Contributors = append(metadata.Contributors, zenodoContrib)
	}

	return metadata
}

// depositionToDOI converts Zenodo deposition to DOI struct
func (p *ZenodoProvider) depositionToDOI(deposition *ZenodoDeposition) *DOI {
	doi := &DOI{
		DOI:             deposition.DOI,
		URL:             deposition.Links["html"],
		State:           deposition.State,
		Title:           deposition.Metadata.Title,
		Publisher:       "Zenodo",
		ResourceType:    deposition.Metadata.UploadType,
		License:         deposition.Metadata.License,
		Description:     deposition.Metadata.Description,
		Metadata:        map[string]interface{}{"zenodo": deposition},
	}

	// Parse publication date
	if deposition.Metadata.PublicationDate != "" {
		if t, err := time.Parse("2006-01-02", deposition.Metadata.PublicationDate); err == nil {
			doi.PublicationYear = t.Year()
		}
	}

	// Parse created/updated
	if t, err := time.Parse(time.RFC3339, deposition.Created); err == nil {
		doi.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, deposition.Modified); err == nil {
		doi.UpdatedAt = t
	}

	// Convert creators to authors
	for _, creator := range deposition.Metadata.Creators {
		author := Author{
			Name:        creator.Name,
			Affiliation: creator.Affiliation,
			ORCID:       creator.ORCID,
		}
		doi.Authors = append(doi.Authors, author)
	}

	return doi
}

// extractRecordID extracts the numeric record ID from a Zenodo DOI
func (p *ZenodoProvider) extractRecordID(doiString string) (int, error) {
	// Zenodo DOIs are like: 10.5281/zenodo.123456
	parts := strings.Split(doiString, ".")
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid Zenodo DOI format: %s", doiString)
	}

	// Last part should be the numeric ID
	var recordID int
	if _, err := fmt.Sscanf(parts[len(parts)-1], "%d", &recordID); err != nil {
		return 0, fmt.Errorf("parse record ID: %w", err)
	}

	return recordID, nil
}

// mapLicenseToZenodo maps common license names to Zenodo license IDs
func mapLicenseToZenodo(license string) string {
	licenseMap := map[string]string{
		"CC-BY-4.0":       "cc-by-4.0",
		"CC-BY-SA-4.0":    "cc-by-sa-4.0",
		"CC-BY-NC-4.0":    "cc-by-nc-4.0",
		"CC-BY-NC-SA-4.0": "cc-by-nc-sa-4.0",
		"CC-BY-ND-4.0":    "cc-by-nd-4.0",
		"CC-BY-NC-ND-4.0": "cc-by-nc-nd-4.0",
		"CC0":             "cc-zero",
		"MIT":             "mit",
		"GPL-3.0":         "gpl-3.0",
		"Apache-2.0":      "apache-2.0",
	}

	if zenodoLicense, ok := licenseMap[license]; ok {
		return zenodoLicense
	}

	// Return as-is if not in map
	return strings.ToLower(license)
}

// mapRelationToZenodo maps DataCite relation types to Zenodo relation types
func mapRelationToZenodo(relation string) string {
	relationMap := map[string]string{
		"IsSupplementTo": "isSupplementTo",
		"IsCitedBy":      "isCitedBy",
		"Cites":          "cites",
		"IsVersionOf":    "isVersionOf",
		"HasVersion":     "hasVersion",
		"IsPartOf":       "isPartOf",
		"HasPart":        "hasPart",
		"IsReferencedBy": "isReferencedBy",
		"References":     "references",
		"IsDocumentedBy": "isDocumentedBy",
		"Documents":      "documents",
		"IsCompiledBy":   "isCompiledBy",
		"Compiles":       "compiles",
		"IsVariantFormOf": "isVariantFormOf",
		"IsOriginalFormOf": "isOriginalFormOf",
		"IsIdenticalTo":  "isIdenticalTo",
		"IsAlternateIdentifier": "isAlternateIdentifier",
	}

	if zenodoRelation, ok := relationMap[relation]; ok {
		return zenodoRelation
	}

	// Return as-is if not in map
	return relation
}

// mapContributorTypeToZenodo maps DataCite contributor types to Zenodo types
func mapContributorTypeToZenodo(contribType string) string {
	typeMap := map[string]string{
		"ContactPerson":     "ContactPerson",
		"DataCollector":     "DataCollector",
		"DataCurator":       "DataCurator",
		"DataManager":       "DataManager",
		"Distributor":       "Distributor",
		"Editor":            "Editor",
		"HostingInstitution": "HostingInstitution",
		"Producer":          "Producer",
		"ProjectLeader":     "ProjectLeader",
		"ProjectManager":    "ProjectManager",
		"ProjectMember":     "ProjectMember",
		"RegistrationAgency": "RegistrationAgency",
		"RegistrationAuthority": "RegistrationAuthority",
		"RelatedPerson":     "RelatedPerson",
		"Researcher":        "Researcher",
		"ResearchGroup":     "ResearchGroup",
		"RightsHolder":      "RightsHolder",
		"Sponsor":           "Sponsor",
		"Supervisor":        "Supervisor",
		"WorkPackageLeader": "WorkPackageLeader",
		"Other":             "Other",
	}

	if zenodoType, ok := typeMap[contribType]; ok {
		return zenodoType
	}

	return "Other"
}
