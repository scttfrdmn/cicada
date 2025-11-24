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
)

func TestMetadataMapper_MapCZI(t *testing.T) {
	mapper := NewMetadataMapper("Test Lab", "CC-BY-4.0", "https://example.com")

	metadata := map[string]interface{}{
		"format":                  "CZI",
		"experiment_name":         "Cell Culture Imaging",
		"sample_name":             "HeLa Cells",
		"experimenter":            "Dr. Jane Smith",
		"operator":                "Lab Tech",
		"manufacturer":            "Zeiss",
		"instrument_model":        "LSM 880",
		"image_width":             2048,
		"image_height":            2048,
		"channel_count":           3,
		"z_planes":                10,
		"timepoints":              5,
		"acquisition_date":        "2025-01-15T10:30:00Z",
		"objective_magnification": 63.0,
		"objective_na":            1.4,
		"pixel_size_x":            0.1,
		"pixel_size_y":            0.1,
		"data_type":               "image",
	}

	dataset, err := mapper.MapToDataset(metadata, "/data/experiment.czi")
	if err != nil {
		t.Fatalf("MapToDataset() error = %v", err)
	}

	// Check required fields
	if dataset.Title != "Cell Culture Imaging" {
		t.Errorf("Title = %v, want Cell Culture Imaging", dataset.Title)
	}

	if len(dataset.Authors) != 2 {
		t.Errorf("len(Authors) = %d, want 2", len(dataset.Authors))
	} else {
		if dataset.Authors[0].Name != "Dr. Jane Smith" {
			t.Errorf("Authors[0].Name = %v, want Dr. Jane Smith", dataset.Authors[0].Name)
		}
		if dataset.Authors[1].Name != "Lab Tech" {
			t.Errorf("Authors[1].Name = %v, want Lab Tech", dataset.Authors[1].Name)
		}
	}

	if dataset.Publisher != "Test Lab" {
		t.Errorf("Publisher = %v, want Test Lab", dataset.Publisher)
	}

	if dataset.ResourceType != "Dataset" {
		t.Errorf("ResourceType = %v, want Dataset", dataset.ResourceType)
	}

	// Check description includes technical details
	if !strings.Contains(dataset.Description, "Zeiss") {
		t.Errorf("Description should contain manufacturer")
	}
	if !strings.Contains(dataset.Description, "2048 x 2048") {
		t.Errorf("Description should contain dimensions")
	}
	if !strings.Contains(dataset.Description, "Channels: 3") {
		t.Errorf("Description should contain channel count")
	}

	// Check keywords
	if len(dataset.Keywords) == 0 {
		t.Error("Keywords should not be empty")
	}
	hasKeyword := func(keyword string) bool {
		for _, kw := range dataset.Keywords {
			if kw == keyword {
				return true
			}
		}
		return false
	}
	if !hasKeyword("microscopy") {
		t.Error("Keywords should include 'microscopy'")
	}

	// Check dates
	if len(dataset.Dates) != 1 {
		t.Errorf("len(Dates) = %d, want 1", len(dataset.Dates))
	} else {
		if dataset.Dates[0].Type != "Collected" {
			t.Errorf("Dates[0].Type = %v, want Collected", dataset.Dates[0].Type)
		}
	}

	// Check custom metadata preservation
	if dataset.Custom["objective_magnification"] != 63.0 {
		t.Error("Custom metadata should preserve objective_magnification")
	}
}

func TestMetadataMapper_MapFASTQ(t *testing.T) {
	mapper := NewMetadataMapper("Genomics Lab", "CC0", "")

	metadata := map[string]interface{}{
		"format":               "FASTQ",
		"total_reads":          1000000,
		"total_bases":          int64(150000000),
		"mean_read_length":     150.0,
		"gc_content_percent":   45.2,
		"mean_quality_score":   35.5,
		"is_paired_end":        true,
		"read_pair":            "R1",
		"compression":          "gzip",
		"data_type":            "nucleotide_sequence",
	}

	dataset, err := mapper.MapToDataset(metadata, "/data/sample_R1.fastq.gz")
	if err != nil {
		t.Fatalf("MapToDataset() error = %v", err)
	}

	// Check title generation
	if !strings.Contains(dataset.Title, "Sequencing Data") {
		t.Errorf("Title should contain 'Sequencing Data', got %v", dataset.Title)
	}

	// Check default author
	if len(dataset.Authors) != 1 || dataset.Authors[0].Name != "Unknown Creator" {
		t.Error("FASTQ should default to Unknown Creator")
	}

	// Check description includes quality metrics
	if !strings.Contains(dataset.Description, "Total reads: 1000000") {
		t.Error("Description should contain total reads")
	}
	if !strings.Contains(dataset.Description, "GC content: 45.2%") {
		t.Error("Description should contain GC content")
	}
	if !strings.Contains(dataset.Description, "Paired-end sequencing (R1)") {
		t.Error("Description should mention paired-end")
	}

	// Check keywords
	hasKeyword := func(keyword string) bool {
		for _, kw := range dataset.Keywords {
			if kw == keyword {
				return true
			}
		}
		return false
	}
	if !hasKeyword("genomics") {
		t.Error("Keywords should include 'genomics'")
	}
	if !hasKeyword("sequencing") {
		t.Error("Keywords should include 'sequencing'")
	}

	// Check custom metadata for paired info
	if dataset.Custom["paired_end"] == nil {
		t.Error("Custom metadata should include paired_end info")
	}
}

func TestMetadataMapper_MapOMETIFF(t *testing.T) {
	mapper := NewMetadataMapper("Imaging Center", "CC-BY-4.0", "")

	metadata := map[string]interface{}{
		"format":                "OME-TIFF",
		"extraction_note":       "OME-TIFF format recognized",
		"implementation_status": "framework - requires TIFF library",
	}

	dataset, err := mapper.MapToDataset(metadata, "/data/image.ome.tiff")
	if err != nil {
		t.Fatalf("MapToDataset() error = %v", err)
	}

	// Check title
	if !strings.Contains(dataset.Title, "OME-TIFF Image") {
		t.Errorf("Title should contain 'OME-TIFF Image', got %v", dataset.Title)
	}

	// Check keywords
	hasKeyword := func(keyword string) bool {
		for _, kw := range dataset.Keywords {
			if kw == keyword {
				return true
			}
		}
		return false
	}
	if !hasKeyword("OME-TIFF") {
		t.Error("Keywords should include 'OME-TIFF'")
	}

	// Check implementation note preserved
	if dataset.Custom["implementation_status"] == nil {
		t.Error("Custom metadata should preserve implementation_status")
	}
}

func TestMetadataMapper_MapGeneric(t *testing.T) {
	mapper := NewMetadataMapper("Research Lab", "MIT", "")

	metadata := map[string]interface{}{
		"format":          "Unknown",
		"data_type":       "tabular",
		"instrument_type": "sensor",
		"custom_field":    "custom_value",
	}

	dataset, err := mapper.MapToDataset(metadata, "/data/experiment.dat")
	if err != nil {
		t.Fatalf("MapToDataset() error = %v", err)
	}

	// Check generic mapping
	if dataset.Title == "" {
		t.Error("Title should not be empty")
	}

	// Check custom fields preserved
	if dataset.Custom["custom_field"] != "custom_value" {
		t.Error("Custom fields should be preserved in generic mapping")
	}

	// Check keywords from metadata
	hasKeyword := func(keyword string) bool {
		for _, kw := range dataset.Keywords {
			if kw == keyword {
				return true
			}
		}
		return false
	}
	if !hasKeyword("tabular") {
		t.Error("Keywords should include data_type")
	}
	if !hasKeyword("sensor") {
		t.Error("Keywords should include instrument_type")
	}
}

func TestMetadataMapper_EnrichDataset(t *testing.T) {
	mapper := NewMetadataMapper("", "", "")

	dataset := &Dataset{
		Title:   "Original Title",
		Authors: []Author{{Name: "Unknown Creator"}},
	}

	enrichment := map[string]interface{}{
		"title":       "Enhanced Title",
		"description": "Detailed description",
		"authors": []interface{}{
			map[string]interface{}{
				"name":        "Dr. John Doe",
				"given_name":  "John",
				"family_name": "Doe",
				"orcid":       "0000-0001-2345-6789",
				"affiliation": "University Lab",
			},
			map[string]interface{}{
				"name":       "Dr. Jane Smith",
				"given_name": "Jane",
				"family_name": "Smith",
			},
		},
		"keywords": []string{"additional", "keywords"},
		"funding": []interface{}{
			map[string]interface{}{
				"funder_name":  "NSF",
				"award_number": "1234567",
				"award_title":  "Research Grant",
			},
		},
	}

	mapper.EnrichDataset(dataset, enrichment)

	// Check enriched title
	if dataset.Title != "Enhanced Title" {
		t.Errorf("Title should be enriched, got %v", dataset.Title)
	}

	// Check enriched description
	if dataset.Description != "Detailed description" {
		t.Errorf("Description should be enriched, got %v", dataset.Description)
	}

	// Check enriched authors
	if len(dataset.Authors) != 2 {
		t.Errorf("Should have 2 authors, got %d", len(dataset.Authors))
	} else {
		if dataset.Authors[0].Name != "Dr. John Doe" {
			t.Errorf("Authors[0].Name = %v, want Dr. John Doe", dataset.Authors[0].Name)
		}
		if dataset.Authors[0].ORCID != "0000-0001-2345-6789" {
			t.Errorf("Authors[0].ORCID should be set")
		}
	}

	// Check additional keywords
	if len(dataset.Keywords) < 2 {
		t.Error("Keywords should be added")
	}

	// Check funding
	if len(dataset.FundingReferences) != 1 {
		t.Errorf("Should have 1 funding reference, got %d", len(dataset.FundingReferences))
	} else {
		if dataset.FundingReferences[0].FunderName != "NSF" {
			t.Error("Funding funder_name should be set")
		}
	}
}

func TestGenerateTitle(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]interface{}
		filename string
		want     string
	}{
		{
			name:     "from experiment_name",
			metadata: map[string]interface{}{"experiment_name": "Test Experiment"},
			filename: "/data/file.czi",
			want:     "Test Experiment",
		},
		{
			name:     "from sample_name",
			metadata: map[string]interface{}{"sample_name": "Sample A"},
			filename: "/data/file.fastq",
			want:     "Sample A",
		},
		{
			name:     "from format and filename",
			metadata: map[string]interface{}{"format": "CZI"},
			filename: "/data/image.czi",
			want:     "CZI Data: image.czi",
		},
		{
			name:     "from filename only",
			metadata: map[string]interface{}{},
			filename: "/data/experiment.dat",
			want:     "experiment.dat",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateTitle(tt.metadata, tt.filename)
			if got != tt.want {
				t.Errorf("GenerateTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractKeywords(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]interface{}
		want     []string
	}{
		{
			name:     "CZI format",
			metadata: map[string]interface{}{"format": "CZI", "data_type": "image", "manufacturer": "Zeiss"},
			want:     []string{"microscopy", "imaging", "confocal", "image", "zeiss"},
		},
		{
			name:     "FASTQ format",
			metadata: map[string]interface{}{"format": "FASTQ", "data_type": "nucleotide_sequence"},
			want:     []string{"genomics", "sequencing", "nucleotide sequence", "nucleotide_sequence"},
		},
		{
			name:     "no duplicates",
			metadata: map[string]interface{}{"format": "CZI", "data_type": "image"},
			want:     []string{"microscopy", "imaging", "confocal", "image"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractKeywords(tt.metadata)
			if len(got) == 0 {
				t.Error("ExtractKeywords() should not return empty slice")
			}
			// Check that expected keywords are present
			for _, wantKW := range tt.want {
				found := false
				for _, gotKW := range got {
					if gotKW == wantKW {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("ExtractKeywords() missing expected keyword: %v", wantKW)
				}
			}
		})
	}
}
