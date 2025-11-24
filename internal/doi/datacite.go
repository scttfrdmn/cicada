// internal/doi/datacite.go
package doi

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

// DOIManager handles DOI minting and management
type DOIManager struct {
	client   *DataCiteClient
	config   *DOIConfig
	metadata MetadataGenerator
}

// DOIConfig contains configuration for DOI minting
type DOIConfig struct {
	Provider     string // "datacite", "zenodo", "institution"
	RepositoryID string
	Password     string
	Prefix       string  // e.g., "10.12345"
	TestMode     bool    // Use DataCite test environment
	CostPerDOI   float64 // For billing
}

// DataCiteClient wraps DataCite API interactions
type DataCiteClient struct {
	BaseURL    string
	Username   string
	Password   string
	HTTPClient *http.Client
}

// DOI represents a Digital Object Identifier
type DOI struct {
	DOI             string                 `json:"doi"`
	URL             string                 `json:"url"`
	State           string                 `json:"state"` // draft, registered, findable
	Title           string                 `json:"title"`
	Authors         []Author               `json:"authors"`
	Publisher       string                 `json:"publisher"`
	PublicationYear int                    `json:"publication_year"`
	ResourceType    string                 `json:"resource_type"`
	License         string                 `json:"license"`
	Description     string                 `json:"description"`
	Metadata        map[string]interface{} `json:"metadata"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Downloads       int                    `json:"downloads"`
	Citations       int                    `json:"citations"`
}

// Author represents a dataset author
type Author struct {
	Name        string `json:"name"`
	GivenName   string `json:"given_name,omitempty"`
	FamilyName  string `json:"family_name,omitempty"`
	ORCID       string `json:"orcid,omitempty"`
	Affiliation string `json:"affiliation,omitempty"`
}

// DataCiteMetadata represents DataCite Metadata Schema 4.4
type DataCiteMetadata struct {
	XMLName xml.Name `xml:"resource"`
	XMLNS   string   `xml:"xmlns,attr"`

	Identifier           Identifier            `xml:"identifier"`
	Creators             []Creator             `xml:"creators>creator"`
	Titles               []Title               `xml:"titles>title"`
	Publisher            string                `xml:"publisher"`
	PublicationYear      int                   `xml:"publicationYear"`
	ResourceType         ResourceType          `xml:"resourceType"`
	Subjects             []Subject             `xml:"subjects>subject,omitempty"`
	Contributors         []DataCiteContributor `xml:"contributors>contributor,omitempty"`
	Dates                []Date                `xml:"dates>date,omitempty"`
	Language             string                `xml:"language,omitempty"`
	AlternateIdentifiers []AlternateIdentifier `xml:"alternateIdentifiers>alternateIdentifier,omitempty"`
	RelatedIdentifiers   []RelatedIdentifier   `xml:"relatedIdentifiers>relatedIdentifier,omitempty"`
	Sizes                []string              `xml:"sizes>size,omitempty"`
	Formats              []string              `xml:"formats>format,omitempty"`
	Version              string                `xml:"version,omitempty"`
	RightsList           []Rights              `xml:"rightsList>rights,omitempty"`
	Descriptions         []Description         `xml:"descriptions>description,omitempty"`
	GeoLocations         []DataCiteGeoLocation `xml:"geoLocations>geoLocation,omitempty"`
	FundingReferences    []FundingReference    `xml:"fundingReferences>fundingReference,omitempty"`
}

type Identifier struct {
	Value string `xml:",chardata"`
	Type  string `xml:"identifierType,attr"`
}

type Creator struct {
	CreatorName    string          `xml:"creatorName"`
	GivenName      string          `xml:"givenName,omitempty"`
	FamilyName     string          `xml:"familyName,omitempty"`
	NameIdentifier *NameIdentifier `xml:"nameIdentifier,omitempty"`
	Affiliation    []string        `xml:"affiliation,omitempty"`
}

type NameIdentifier struct {
	Value     string `xml:",chardata"`
	Scheme    string `xml:"nameIdentifierScheme,attr"`
	SchemeURI string `xml:"schemeURI,attr,omitempty"`
}

type Title struct {
	Value string `xml:",chardata"`
	Lang  string `xml:"xml:lang,attr,omitempty"`
	Type  string `xml:"titleType,attr,omitempty"`
}

type ResourceType struct {
	Value string `xml:",chardata"`
	Type  string `xml:"resourceTypeGeneral,attr"`
}

type Subject struct {
	Value         string `xml:",chardata"`
	SubjectScheme string `xml:"subjectScheme,attr,omitempty"`
	SchemeURI     string `xml:"schemeURI,attr,omitempty"`
}

type DataCiteContributor struct {
	ContributorName string   `xml:"contributorName"`
	ContributorType string   `xml:"contributorType,attr"`
	GivenName       string   `xml:"givenName,omitempty"`
	FamilyName      string   `xml:"familyName,omitempty"`
	Affiliation     []string `xml:"affiliation,omitempty"`
}

type Date struct {
	Value string `xml:",chardata"`
	Type  string `xml:"dateType,attr"`
}

type AlternateIdentifier struct {
	Value string `xml:",chardata"`
	Type  string `xml:"alternateIdentifierType,attr"`
}

type RelatedIdentifier struct {
	Value        string `xml:",chardata"`
	Type         string `xml:"relatedIdentifierType,attr"`
	RelationType string `xml:"relationType,attr"`
}

type Rights struct {
	Value            string `xml:",chardata"`
	RightsURI        string `xml:"rightsURI,attr,omitempty"`
	RightsIdentifier string `xml:"rightsIdentifier,attr,omitempty"`
}

type Description struct {
	Value string `xml:",chardata"`
	Type  string `xml:"descriptionType,attr"`
}

type DataCiteGeoLocation struct {
	GeoLocationPlace string `xml:"geoLocationPlace,omitempty"`
	GeoLocationPoint *Point `xml:"geoLocationPoint,omitempty"`
	GeoLocationBox   *Box   `xml:"geoLocationBox,omitempty"`
}

type Point struct {
	PointLongitude float64 `xml:"pointLongitude"`
	PointLatitude  float64 `xml:"pointLatitude"`
}

type Box struct {
	WestBoundLongitude float64 `xml:"westBoundLongitude"`
	EastBoundLongitude float64 `xml:"eastBoundLongitude"`
	SouthBoundLatitude float64 `xml:"southBoundLatitude"`
	NorthBoundLatitude float64 `xml:"northBoundLatitude"`
}

type FundingReference struct {
	FunderName       string            `xml:"funderName"`
	FunderIdentifier *FunderIdentifier `xml:"funderIdentifier,omitempty"`
	AwardNumber      string            `xml:"awardNumber,omitempty"`
	AwardTitle       string            `xml:"awardTitle,omitempty"`
}

type FunderIdentifier struct {
	Value string `xml:",chardata"`
	Type  string `xml:"funderIdentifierType,attr"`
}

// NewDOIManager creates a new DOI manager
func NewDOIManager(config *DOIConfig) *DOIManager {
	baseURL := "https://api.datacite.org"
	if config.TestMode {
		baseURL = "https://api.test.datacite.org"
	}

	client := &DataCiteClient{
		BaseURL:  baseURL,
		Username: config.RepositoryID,
		Password: config.Password,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	return &DOIManager{
		client:   client,
		config:   config,
		metadata: &StandardMetadataGenerator{},
	}
}

// MintDOI mints a new DOI for a dataset
func (m *DOIManager) MintDOI(dataset *Dataset) (*DOI, error) {
	// Generate DOI suffix
	suffix := generateDOISuffix(dataset)
	doiString := fmt.Sprintf("%s/%s", m.config.Prefix, suffix)

	// Generate DataCite metadata XML
	metadata := m.metadata.Generate(dataset)
	metadata.Identifier = Identifier{
		Value: doiString,
		Type:  "DOI",
	}

	xmlData, err := xml.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
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
				"event": "publish", // Or "register" for draft
			},
		},
	}

	_, err = m.client.CreateDOI(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create DOI: %w", err)
	}

	doi := &DOI{
		DOI:             doiString,
		URL:             dataset.URL,
		State:           "findable",
		Title:           dataset.Title,
		Authors:         dataset.Authors,
		Publisher:       dataset.Publisher,
		PublicationYear: dataset.PublicationYear,
		ResourceType:    dataset.ResourceType,
		License:         dataset.License,
		Description:     dataset.Description,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return doi, nil
}

// UpdateDOI updates an existing DOI's metadata
func (m *DOIManager) UpdateDOI(doiString string, updates map[string]interface{}) error {
	// Get current metadata
	current, err := m.client.GetDOI(doiString)
	if err != nil {
		return fmt.Errorf("failed to get current DOI: %w", err)
	}

	// Apply updates
	// TODO: Merge updates into current metadata
	_ = current // Silence unused variable warning

	// Submit update
	// TODO: Convert DOI struct to payload format
	return fmt.Errorf("UpdateDOI not yet implemented")
}

// GetDOI retrieves DOI information
func (m *DOIManager) GetDOI(doiString string) (*DOI, error) {
	return m.client.GetDOI(doiString)
}

// ListDOIs lists all DOIs for the repository
func (m *DOIManager) ListDOIs() ([]*DOI, error) {
	return m.client.ListDOIs(m.config.RepositoryID)
}

// CreateDOI creates a DOI via DataCite API
func (c *DataCiteClient) CreateDOI(payload map[string]interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/dois", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/vnd.api+json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetDOI retrieves a DOI from DataCite
func (c *DataCiteClient) GetDOI(doi string) (*DOI, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dois/%s", c.BaseURL, doi), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Accept", "application/vnd.api+json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	// TODO: Parse response into DOI struct
	return nil, nil
}

// UpdateDOI updates a DOI via DataCite
func (c *DataCiteClient) UpdateDOI(doi string, updates map[string]interface{}) error {
	jsonData, err := json.Marshal(updates)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/dois/%s", c.BaseURL, doi), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/vnd.api+json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return nil
}

// ListDOIs lists all DOIs for a repository
func (c *DataCiteClient) ListDOIs(repositoryID string) ([]*DOI, error) {
	// TODO: Implement pagination
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dois?client-id=%s", c.BaseURL, repositoryID), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Accept", "application/vnd.api+json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	// TODO: Parse response
	return nil, nil
}

// MetadataGenerator generates DataCite metadata
// Uses Dataset from provider.go
type MetadataGenerator interface {
	Generate(dataset *Dataset) *DataCiteMetadata
}

// StandardMetadataGenerator generates standard DataCite metadata
type StandardMetadataGenerator struct{}

func (g *StandardMetadataGenerator) Generate(dataset *Dataset) *DataCiteMetadata {
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

func generateDOISuffix(dataset *Dataset) string {
	// Generate a unique suffix
	// Format: cicada.{lab-name}.{year}.{sequence}
	// Example: cicada.smith-lab.2024.001
	timestamp := time.Now().Unix()
	return fmt.Sprintf("cicada.%d", timestamp)
}

func getLicenseURI(license string) string {
	licenses := map[string]string{
		"CC-BY-4.0":    "https://creativecommons.org/licenses/by/4.0/",
		"CC-BY-SA-4.0": "https://creativecommons.org/licenses/by-sa/4.0/",
		"CC0":          "https://creativecommons.org/publicdomain/zero/1.0/",
		"MIT":          "https://opensource.org/licenses/MIT",
	}

	if uri, ok := licenses[license]; ok {
		return uri
	}
	return ""
}
