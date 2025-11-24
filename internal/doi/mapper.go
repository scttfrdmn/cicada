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
	"path/filepath"
	"strings"
	"time"
)

// MetadataMapper maps Cicada metadata to DOI Dataset structure
type MetadataMapper struct {
	// Default values for required fields
	DefaultPublisher string
	DefaultLicense   string
	DefaultURL       string
}

// NewMetadataMapper creates a new metadata mapper with defaults
func NewMetadataMapper(publisher, license, url string) *MetadataMapper {
	if publisher == "" {
		publisher = "Unknown Publisher"
	}
	if license == "" {
		license = "CC-BY-4.0"
	}
	return &MetadataMapper{
		DefaultPublisher: publisher,
		DefaultLicense:   license,
		DefaultURL:       url,
	}
}

// MapToDataset converts Cicada extractor metadata to DOI Dataset
func (m *MetadataMapper) MapToDataset(metadata map[string]interface{}, filename string) (*Dataset, error) {
	format, _ := metadata["format"].(string)

	// Route to format-specific mapper
	switch format {
	case "CZI":
		return m.mapCZI(metadata, filename)
	case "FASTQ":
		return m.mapFASTQ(metadata, filename)
	case "OME-TIFF":
		return m.mapOMETIFF(metadata, filename)
	default:
		return m.mapGeneric(metadata, filename)
	}
}

// mapCZI maps CZI microscopy metadata to Dataset
func (m *MetadataMapper) mapCZI(metadata map[string]interface{}, filename string) (*Dataset, error) {
	dataset := &Dataset{
		Publisher:       m.DefaultPublisher,
		PublicationYear: time.Now().Year(),
		ResourceType:    "Dataset",
		License:         m.DefaultLicense,
		Custom:          make(map[string]interface{}),
	}

	// Title from experiment name or filename
	if expName, ok := metadata["experiment_name"].(string); ok && expName != "" {
		dataset.Title = expName
	} else if sampleName, ok := metadata["sample_name"].(string); ok && sampleName != "" {
		dataset.Title = sampleName
	} else {
		dataset.Title = fmt.Sprintf("Microscopy Image: %s", filepath.Base(filename))
	}

	// Authors from experimenter/operator
	authors := make([]Author, 0)
	if experimenter, ok := metadata["experimenter"].(string); ok && experimenter != "" {
		authors = append(authors, Author{Name: experimenter})
	}
	if operator, ok := metadata["operator"].(string); ok && operator != "" && operator != metadata["experimenter"] {
		authors = append(authors, Author{Name: operator})
	}
	if len(authors) == 0 {
		authors = append(authors, Author{Name: "Unknown Creator"})
	}
	dataset.Authors = authors

	// Description from technical details
	desc := []string{}
	if manufacturer, ok := metadata["manufacturer"].(string); ok {
		desc = append(desc, fmt.Sprintf("Manufacturer: %s", manufacturer))
	}
	if model, ok := metadata["instrument_model"].(string); ok {
		desc = append(desc, fmt.Sprintf("Instrument: %s", model))
	}
	if width, ok := metadata["image_width"].(int); ok {
		if height, ok2 := metadata["image_height"].(int); ok2 {
			desc = append(desc, fmt.Sprintf("Dimensions: %d x %d pixels", width, height))
		}
	}
	if channels, ok := metadata["channel_count"].(int); ok && channels > 0 {
		desc = append(desc, fmt.Sprintf("Channels: %d", channels))
	}
	if zplanes, ok := metadata["z_planes"].(int); ok && zplanes > 1 {
		desc = append(desc, fmt.Sprintf("Z-planes: %d", zplanes))
	}
	if timepoints, ok := metadata["timepoints"].(int); ok && timepoints > 1 {
		desc = append(desc, fmt.Sprintf("Timepoints: %d", timepoints))
	}
	dataset.Description = strings.Join(desc, "; ")

	// Keywords/subjects
	keywords := []string{"microscopy", "imaging", "confocal"}
	if manufacturer, ok := metadata["manufacturer"].(string); ok {
		keywords = append(keywords, strings.ToLower(manufacturer))
	}
	if dataType, ok := metadata["data_type"].(string); ok {
		keywords = append(keywords, dataType)
	}
	dataset.Keywords = keywords

	// Dates
	if acqDate, ok := metadata["acquisition_date"].(string); ok && acqDate != "" {
		dataset.Dates = append(dataset.Dates, DateInfo{
			Date: acqDate,
			Type: "Collected",
		})
	}

	// Custom metadata for preservation
	preserveFields := []string{
		"objective_magnification", "objective_na", "zoom_factor",
		"pixel_size_x", "pixel_size_y", "pixel_size_z",
		"excitation_wavelengths", "emission_wavelengths",
		"detector_type", "pinhole_size",
	}
	for _, field := range preserveFields {
		if val, ok := metadata[field]; ok {
			dataset.Custom[field] = val
		}
	}

	return dataset, nil
}

// mapFASTQ maps FASTQ sequencing metadata to Dataset
func (m *MetadataMapper) mapFASTQ(metadata map[string]interface{}, filename string) (*Dataset, error) {
	dataset := &Dataset{
		Publisher:       m.DefaultPublisher,
		PublicationYear: time.Now().Year(),
		ResourceType:    "Dataset",
		License:         m.DefaultLicense,
		Custom:          make(map[string]interface{}),
	}

	// Title from filename
	baseName := filepath.Base(filename)
	baseName = strings.TrimSuffix(baseName, ".gz")
	baseName = strings.TrimSuffix(baseName, ".fastq")
	baseName = strings.TrimSuffix(baseName, ".fq")

	// Clean up common naming patterns
	baseName = strings.ReplaceAll(baseName, "_", " ")
	baseName = strings.ReplaceAll(baseName, "-", " ")

	dataset.Title = fmt.Sprintf("Sequencing Data: %s", baseName)

	// Authors - require user configuration for sequencing data
	dataset.Authors = []Author{{Name: "Unknown Creator"}}

	// Description from quality metrics
	desc := []string{}
	if totalReads, ok := metadata["total_reads"].(int); ok {
		desc = append(desc, fmt.Sprintf("Total reads: %d", totalReads))
	}
	if totalBases, ok := metadata["total_bases"].(int64); ok {
		desc = append(desc, fmt.Sprintf("Total bases: %d", totalBases))
	}
	if meanLength, ok := metadata["mean_read_length"].(float64); ok {
		desc = append(desc, fmt.Sprintf("Mean read length: %.1f bp", meanLength))
	}
	if gcContent, ok := metadata["gc_content_percent"].(float64); ok {
		desc = append(desc, fmt.Sprintf("GC content: %.1f%%", gcContent))
	}
	if meanQuality, ok := metadata["mean_quality_score"].(float64); ok {
		desc = append(desc, fmt.Sprintf("Mean quality: %.1f", meanQuality))
	}

	// Add pairing information
	if isPaired, ok := metadata["is_paired_end"].(bool); ok && isPaired {
		if readPair, ok := metadata["read_pair"].(string); ok {
			desc = append(desc, fmt.Sprintf("Paired-end sequencing (%s)", readPair))
		}
	}

	dataset.Description = strings.Join(desc, "; ")

	// Keywords
	keywords := []string{"genomics", "sequencing", "high-throughput sequencing", "nucleotide sequence"}
	if dataType, ok := metadata["data_type"].(string); ok {
		keywords = append(keywords, dataType)
	}
	dataset.Keywords = keywords

	// Related identifiers for paired files
	if isPaired, ok := metadata["is_paired_end"].(bool); ok && isPaired {
		if readPair, ok := metadata["read_pair"].(string); ok {
			// Create placeholder for the paired file
			// This will need to be filled in by the user or workflow
			otherPair := "R2"
			relType := "IsPartOf"
			if readPair == "R2" || readPair == "2" {
				otherPair = "R1"
			}

			dataset.Custom["paired_end"] = map[string]interface{}{
				"this_pair":  readPair,
				"other_pair": otherPair,
				"relation":   relType,
			}
		}
	}

	// Custom metadata
	preserveFields := []string{
		"total_reads", "total_bases", "mean_read_length", "min_read_length", "max_read_length",
		"gc_content_percent", "mean_quality_score", "min_quality_score", "max_quality_score",
		"compression",
	}
	for _, field := range preserveFields {
		if val, ok := metadata[field]; ok {
			dataset.Custom[field] = val
		}
	}

	return dataset, nil
}

// mapOMETIFF maps OME-TIFF microscopy metadata to Dataset
func (m *MetadataMapper) mapOMETIFF(metadata map[string]interface{}, filename string) (*Dataset, error) {
	dataset := &Dataset{
		Publisher:       m.DefaultPublisher,
		PublicationYear: time.Now().Year(),
		ResourceType:    "Dataset",
		License:         m.DefaultLicense,
		Custom:          make(map[string]interface{}),
	}

	// Title from filename (OME-TIFF extraction not yet complete)
	dataset.Title = fmt.Sprintf("OME-TIFF Image: %s", filepath.Base(filename))

	// Authors - require user configuration
	dataset.Authors = []Author{{Name: "Unknown Creator"}}

	// Description
	dataset.Description = "OME-TIFF microscopy image. Full metadata extraction requires TIFF library support."

	// Keywords
	dataset.Keywords = []string{"microscopy", "biological imaging", "OME-TIFF"}

	// Note about implementation status
	if note, ok := metadata["extraction_note"].(string); ok {
		dataset.Custom["extraction_note"] = note
	}
	if status, ok := metadata["implementation_status"].(string); ok {
		dataset.Custom["implementation_status"] = status
	}

	return dataset, nil
}

// mapGeneric provides generic mapping for unknown formats
func (m *MetadataMapper) mapGeneric(metadata map[string]interface{}, filename string) (*Dataset, error) {
	dataset := &Dataset{
		Publisher:       m.DefaultPublisher,
		PublicationYear: time.Now().Year(),
		ResourceType:    "Dataset",
		License:         m.DefaultLicense,
		Custom:          make(map[string]interface{}),
	}

	// Title
	format, _ := metadata["format"].(string)
	if format != "" {
		dataset.Title = fmt.Sprintf("%s Data: %s", format, filepath.Base(filename))
	} else {
		dataset.Title = filepath.Base(filename)
	}

	// Authors
	dataset.Authors = []Author{{Name: "Unknown Creator"}}

	// Description
	if format != "" {
		dataset.Description = fmt.Sprintf("Scientific data in %s format", format)
	} else {
		dataset.Description = "Scientific dataset"
	}

	// Extract any available information
	if dataType, ok := metadata["data_type"].(string); ok {
		dataset.Keywords = append(dataset.Keywords, dataType)
	}
	if instrumentType, ok := metadata["instrument_type"].(string); ok {
		dataset.Keywords = append(dataset.Keywords, instrumentType)
	}

	// Preserve all metadata in custom fields
	for key, val := range metadata {
		// Skip standard fields already processed
		if key != "format" && key != "file_name" && key != "extractor_name" && key != "schema_name" {
			dataset.Custom[key] = val
		}
	}

	return dataset, nil
}

// EnrichDataset enriches a Dataset with user-provided information
func (m *MetadataMapper) EnrichDataset(dataset *Dataset, enrichment map[string]interface{}) {
	// Override or add authors
	if authorsData, ok := enrichment["authors"].([]interface{}); ok {
		authors := make([]Author, 0, len(authorsData))
		for _, authorData := range authorsData {
			if authorMap, ok := authorData.(map[string]interface{}); ok {
				author := Author{}
				if name, ok := authorMap["name"].(string); ok {
					author.Name = name
				}
				if givenName, ok := authorMap["given_name"].(string); ok {
					author.GivenName = givenName
				}
				if familyName, ok := authorMap["family_name"].(string); ok {
					author.FamilyName = familyName
				}
				if orcid, ok := authorMap["orcid"].(string); ok {
					author.ORCID = orcid
				}
				if affiliation, ok := authorMap["affiliation"].(string); ok {
					author.Affiliation = affiliation
				}
				if author.Name != "" {
					authors = append(authors, author)
				}
			}
		}
		if len(authors) > 0 {
			dataset.Authors = authors
		}
	}

	// Override or add other fields
	if title, ok := enrichment["title"].(string); ok && title != "" {
		dataset.Title = title
	}
	if description, ok := enrichment["description"].(string); ok && description != "" {
		dataset.Description = description
	}
	if publisher, ok := enrichment["publisher"].(string); ok && publisher != "" {
		dataset.Publisher = publisher
	}
	if year, ok := enrichment["publication_year"].(int); ok && year > 0 {
		dataset.PublicationYear = year
	}
	if license, ok := enrichment["license"].(string); ok && license != "" {
		dataset.License = license
	}
	if url, ok := enrichment["url"].(string); ok && url != "" {
		dataset.URL = url
	}
	if keywords, ok := enrichment["keywords"].([]string); ok {
		dataset.Keywords = append(dataset.Keywords, keywords...)
	}

	// Add funding references
	if fundingData, ok := enrichment["funding"].([]interface{}); ok {
		for _, fundData := range fundingData {
			if fundMap, ok := fundData.(map[string]interface{}); ok {
				funding := FundingRef{}
				if funderName, ok := fundMap["funder_name"].(string); ok {
					funding.FunderName = funderName
				}
				if awardNumber, ok := fundMap["award_number"].(string); ok {
					funding.AwardNumber = awardNumber
				}
				if awardTitle, ok := fundMap["award_title"].(string); ok {
					funding.AwardTitle = awardTitle
				}
				if funding.FunderName != "" {
					dataset.FundingReferences = append(dataset.FundingReferences, funding)
				}
			}
		}
	}
}

// GenerateTitle generates a descriptive title from metadata
func GenerateTitle(metadata map[string]interface{}, filename string) string {
	format, _ := metadata["format"].(string)

	// Try to find meaningful names
	if expName, ok := metadata["experiment_name"].(string); ok && expName != "" {
		return expName
	}
	if sampleName, ok := metadata["sample_name"].(string); ok && sampleName != "" {
		return sampleName
	}

	// Use format and filename
	baseName := filepath.Base(filename)
	if format != "" {
		return fmt.Sprintf("%s Data: %s", format, baseName)
	}

	return baseName
}

// ExtractKeywords extracts relevant keywords from metadata
func ExtractKeywords(metadata map[string]interface{}) []string {
	keywords := make([]string, 0)
	seen := make(map[string]bool)

	// Add format-specific keywords
	format, _ := metadata["format"].(string)
	switch format {
	case "CZI":
		keywords = append(keywords, "microscopy", "imaging", "confocal")
	case "FASTQ":
		keywords = append(keywords, "genomics", "sequencing", "nucleotide sequence")
	case "OME-TIFF":
		keywords = append(keywords, "microscopy", "biological imaging", "OME-TIFF")
	}

	// Add from metadata fields
	if dataType, ok := metadata["data_type"].(string); ok && dataType != "" {
		if !seen[dataType] {
			keywords = append(keywords, dataType)
			seen[dataType] = true
		}
	}
	if instrumentType, ok := metadata["instrument_type"].(string); ok && instrumentType != "" {
		if !seen[instrumentType] {
			keywords = append(keywords, instrumentType)
			seen[instrumentType] = true
		}
	}
	if manufacturer, ok := metadata["manufacturer"].(string); ok && manufacturer != "" {
		mfr := strings.ToLower(manufacturer)
		if !seen[mfr] {
			keywords = append(keywords, mfr)
			seen[mfr] = true
		}
	}

	return keywords
}
