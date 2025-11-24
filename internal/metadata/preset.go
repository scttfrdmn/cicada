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

// Package metadata provides metadata extraction for scientific instrument files.
//
// # Instrument Presets
//
// This file implements an instrument preset system for validating and managing
// metadata from scientific instruments. Presets define expected metadata fields,
// validation rules, and requirements for specific instrument types.
//
// ## Overview
//
// The preset system helps ensure metadata quality and FAIR compliance by:
//   - Defining expected metadata fields for each instrument
//   - Validating extracted metadata against instrument specifications
//   - Providing templates for metadata collection
//   - Ensuring consistency across datasets from the same instrument
//
// ## Preset Structure
//
// Each preset includes:
//   - Instrument identification (ID, name, manufacturer, model)
//   - Required metadata fields (must be present)
//   - Optional metadata fields (should be present when available)
//   - Field validation rules (type, format, range constraints)
//   - Supported file formats
//   - Documentation and references
//
// ## Use Cases
//
// **Metadata Validation**:
//
//	preset := registry.GetPreset("zeiss-lsm-880")
//	result := preset.Validate(metadata)
//	if !result.IsValid {
//	    // Handle missing or invalid fields
//	}
//
// **Template Generation**:
//
//	preset := registry.GetPreset("illumina-novaseq")
//	template := preset.GenerateTemplate()
//	// Use template for data entry forms
//
// **Quality Assessment**:
//
//	result := preset.Validate(metadata)
//	score := result.QualityScore() // 0-100
//	// Assess metadata completeness
//
// ## Preset Registry
//
// The PresetRegistry manages multiple presets and provides:
//   - Preset lookup by ID, manufacturer, or model
//   - Listing available presets
//   - Filtering presets by instrument type
//   - Loading presets from YAML files
//
// ## FAIR Principles
//
// Presets support FAIR data principles:
//   - **Findable**: Ensures key identifiers are present
//   - **Accessible**: Validates access-related metadata
//   - **Interoperable**: Enforces standard field names and formats
//   - **Reusable**: Ensures provenance and usage information
//
// ## Implementation Notes
//
// This system:
//   - Uses Go structs for type safety
//   - Supports YAML file storage for easy editing
//   - Provides flexible validation rules
//   - Allows custom field validators
//   - Includes built-in presets for common instruments
//
package metadata

import (
	"fmt"
	"regexp"
	"strings"
)

// InstrumentPreset defines expected metadata for a specific instrument.
type InstrumentPreset struct {
	// Identification
	ID           string   `yaml:"id" json:"id"`
	Name         string   `yaml:"name" json:"name"`
	Manufacturer string   `yaml:"manufacturer" json:"manufacturer"`
	Models       []string `yaml:"models" json:"models"`
	Description  string   `yaml:"description" json:"description"`

	// Instrument classification
	InstrumentType string   `yaml:"instrument_type" json:"instrument_type"` // microscopy, sequencing, mass_spec, etc.
	DataTypes      []string `yaml:"data_types" json:"data_types"`           // image, nucleotide_sequence, mass_spectrum, etc.

	// Supported file formats
	FileFormats []string `yaml:"file_formats" json:"file_formats"` // .czi, .fastq, etc.

	// Metadata requirements
	RequiredFields []FieldRequirement `yaml:"required_fields" json:"required_fields"`
	OptionalFields []FieldRequirement `yaml:"optional_fields" json:"optional_fields"`

	// References and documentation
	Documentation string   `yaml:"documentation,omitempty" json:"documentation,omitempty"`
	References    []string `yaml:"references,omitempty" json:"references,omitempty"`
}

// FieldRequirement defines requirements for a metadata field.
type FieldRequirement struct {
	Name        string      `yaml:"name" json:"name"`
	Description string      `yaml:"description" json:"description"`
	Type        string      `yaml:"type" json:"type"`                   // string, number, boolean, array, object
	Format      string      `yaml:"format,omitempty" json:"format,omitempty"` // date-time, url, email, etc.
	Pattern     string      `yaml:"pattern,omitempty" json:"pattern,omitempty"` // regex pattern
	MinValue    *float64    `yaml:"min_value,omitempty" json:"min_value,omitempty"`
	MaxValue    *float64    `yaml:"max_value,omitempty" json:"max_value,omitempty"`
	Enum        []string    `yaml:"enum,omitempty" json:"enum,omitempty"`       // allowed values
	Example     interface{} `yaml:"example,omitempty" json:"example,omitempty"` // example value
}

// PresetRegistry manages instrument presets.
type PresetRegistry struct {
	presets map[string]*InstrumentPreset
}

// NewPresetRegistry creates a new preset registry.
func NewPresetRegistry() *PresetRegistry {
	return &PresetRegistry{
		presets: make(map[string]*InstrumentPreset),
	}
}

// Register registers a preset in the registry.
func (r *PresetRegistry) Register(preset *InstrumentPreset) {
	r.presets[preset.ID] = preset
}

// GetPreset retrieves a preset by ID.
func (r *PresetRegistry) GetPreset(id string) (*InstrumentPreset, error) {
	preset, exists := r.presets[id]
	if !exists {
		return nil, fmt.Errorf("preset not found: %s", id)
	}
	return preset, nil
}

// ListPresets returns all registered presets.
func (r *PresetRegistry) ListPresets() []*InstrumentPreset {
	presets := make([]*InstrumentPreset, 0, len(r.presets))
	for _, preset := range r.presets {
		presets = append(presets, preset)
	}
	return presets
}

// FindPresets finds presets matching criteria.
func (r *PresetRegistry) FindPresets(manufacturer, instrumentType string) []*InstrumentPreset {
	var matches []*InstrumentPreset
	for _, preset := range r.presets {
		manufacturerMatch := manufacturer == "" || strings.EqualFold(preset.Manufacturer, manufacturer)
		typeMatch := instrumentType == "" || strings.EqualFold(preset.InstrumentType, instrumentType)

		if manufacturerMatch && typeMatch {
			matches = append(matches, preset)
		}
	}
	return matches
}

// RegisterDefaults registers built-in presets for common instruments.
func (r *PresetRegistry) RegisterDefaults() {
	// Register Zeiss LSM presets
	r.Register(zeissLSM880Preset())
	r.Register(zeissLSM900Preset())
	r.Register(zeissLSM980Preset())

	// Register Illumina sequencing presets
	r.Register(illuminaNovaSeqPreset())
	r.Register(illuminaMiSeqPreset())
	r.Register(illuminaNextSeqPreset())

	// Register generic presets
	r.Register(genericMicroscopyPreset())
	r.Register(genericSequencingPreset())
}

// PresetValidationResult contains results from validating metadata against a preset.
type PresetValidationResult struct {
	IsValid  bool     `json:"is_valid"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
	Missing  []string `json:"missing,omitempty"`
	Present  []string `json:"present,omitempty"`
}

// Validate validates metadata against the preset requirements.
func (p *InstrumentPreset) Validate(metadata map[string]interface{}) *PresetValidationResult {
	result := &PresetValidationResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
		Missing:  []string{},
		Present:  []string{},
	}

	// Check required fields
	for _, req := range p.RequiredFields {
		value, exists := metadata[req.Name]
		if !exists {
			result.IsValid = false
			result.Errors = append(result.Errors, fmt.Sprintf("missing required field: %s", req.Name))
			result.Missing = append(result.Missing, req.Name)
			continue
		}

		result.Present = append(result.Present, req.Name)

		// Validate field value
		if err := validateFieldValue(req, value); err != nil {
			result.IsValid = false
			result.Errors = append(result.Errors, fmt.Sprintf("invalid value for %s: %v", req.Name, err))
		}
	}

	// Check optional fields (warnings only)
	for _, req := range p.OptionalFields {
		value, exists := metadata[req.Name]
		if !exists {
			result.Warnings = append(result.Warnings, fmt.Sprintf("missing optional field: %s", req.Name))
			result.Missing = append(result.Missing, req.Name)
			continue
		}

		result.Present = append(result.Present, req.Name)

		// Validate field value
		if err := validateFieldValue(req, value); err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("invalid value for %s: %v", req.Name, err))
		}
	}

	return result
}

// validateFieldValue validates a field value against requirements.
func validateFieldValue(req FieldRequirement, value interface{}) error {
	// Type validation
	switch req.Type {
	case "string":
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", value)
		}

		// Pattern validation
		if req.Pattern != "" {
			matched, err := regexp.MatchString(req.Pattern, str)
			if err != nil {
				return fmt.Errorf("invalid pattern: %v", err)
			}
			if !matched {
				return fmt.Errorf("does not match pattern %s", req.Pattern)
			}
		}

		// Enum validation
		if len(req.Enum) > 0 {
			found := false
			for _, allowed := range req.Enum {
				if str == allowed {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("must be one of: %v", req.Enum)
			}
		}

	case "number":
		var num float64
		switch v := value.(type) {
		case float64:
			num = v
		case float32:
			num = float64(v)
		case int:
			num = float64(v)
		case int64:
			num = float64(v)
		default:
			return fmt.Errorf("expected number, got %T", value)
		}

		// Range validation
		if req.MinValue != nil && num < *req.MinValue {
			return fmt.Errorf("must be >= %f", *req.MinValue)
		}
		if req.MaxValue != nil && num > *req.MaxValue {
			return fmt.Errorf("must be <= %f", *req.MaxValue)
		}

	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean, got %T", value)
		}

	case "array":
		// Check if it's a slice
		switch value.(type) {
		case []interface{}, []string, []int, []float64:
			// Valid array types
		default:
			return fmt.Errorf("expected array, got %T", value)
		}

	case "object":
		if _, ok := value.(map[string]interface{}); !ok {
			return fmt.Errorf("expected object, got %T", value)
		}
	}

	return nil
}

// QualityScore calculates a quality score (0-100) based on validation results.
func (v *PresetValidationResult) QualityScore() float64 {
	if len(v.Present) == 0 && len(v.Missing) == 0 {
		return 0
	}

	total := len(v.Present) + len(v.Missing)
	score := float64(len(v.Present)) / float64(total) * 100

	// Penalize for errors (required field issues)
	errorPenalty := float64(len(v.Errors)) * 10
	score -= errorPenalty

	if score < 0 {
		score = 0
	}

	return score
}

// GenerateTemplate generates a metadata template from the preset.
func (p *InstrumentPreset) GenerateTemplate() map[string]interface{} {
	template := make(map[string]interface{})

	// Add required fields
	for _, req := range p.RequiredFields {
		if req.Example != nil {
			template[req.Name] = req.Example
		} else {
			template[req.Name] = fmt.Sprintf("<%s>", req.Description)
		}
	}

	// Add optional fields with nil/example values
	for _, req := range p.OptionalFields {
		if req.Example != nil {
			template[req.Name] = req.Example
		}
	}

	return template
}

// Built-in preset definitions

func zeissLSM880Preset() *InstrumentPreset {
	return &InstrumentPreset{
		ID:             "zeiss-lsm-880",
		Name:           "Zeiss LSM 880",
		Manufacturer:   "Zeiss",
		Models:         []string{"LSM 880"},
		Description:    "Zeiss LSM 880 confocal laser scanning microscope",
		InstrumentType: "microscopy",
		DataTypes:      []string{"image"},
		FileFormats:    []string{".czi"},
		RequiredFields: []FieldRequirement{
			{Name: "format", Type: "string", Description: "File format", Enum: []string{"CZI"}, Example: "CZI"},
			{Name: "manufacturer", Type: "string", Description: "Instrument manufacturer", Enum: []string{"Zeiss"}, Example: "Zeiss"},
			{Name: "instrument_model", Type: "string", Description: "Microscope model", Example: "LSM 880"},
			{Name: "image_width", Type: "number", Description: "Image width in pixels", MinValue: ptr(1.0), Example: 1024},
			{Name: "image_height", Type: "number", Description: "Image height in pixels", MinValue: ptr(1.0), Example: 1024},
		},
		OptionalFields: []FieldRequirement{
			{Name: "acquisition_date", Type: "string", Description: "When data was acquired", Format: "date-time"},
			{Name: "operator", Type: "string", Description: "Operator name"},
			{Name: "objective_magnification", Type: "number", Description: "Objective magnification", Example: 40.0},
			{Name: "objective_na", Type: "number", Description: "Objective numerical aperture", MinValue: ptr(0.0), MaxValue: ptr(2.0), Example: 1.3},
			{Name: "pixel_size_x_um", Type: "number", Description: "Pixel size X in micrometers", MinValue: ptr(0.0)},
			{Name: "pixel_size_y_um", Type: "number", Description: "Pixel size Y in micrometers", MinValue: ptr(0.0)},
			{Name: "num_channels", Type: "number", Description: "Number of channels", MinValue: ptr(1.0)},
			{Name: "channels", Type: "array", Description: "Channel information"},
		},
		Documentation: "Zeiss LSM 880 is a confocal laser scanning microscope for high-resolution fluorescence imaging",
		References: []string{
			"https://www.zeiss.com/microscopy/us/products/confocal-microscopes/lsm-880.html",
		},
	}
}

func zeissLSM900Preset() *InstrumentPreset {
	preset := zeissLSM880Preset()
	preset.ID = "zeiss-lsm-900"
	preset.Name = "Zeiss LSM 900"
	preset.Models = []string{"LSM 900"}
	preset.Description = "Zeiss LSM 900 confocal laser scanning microscope with Airyscan 2"
	preset.References = []string{
		"https://www.zeiss.com/microscopy/us/products/confocal-microscopes/lsm-900.html",
	}
	return preset
}

func zeissLSM980Preset() *InstrumentPreset {
	preset := zeissLSM880Preset()
	preset.ID = "zeiss-lsm-980"
	preset.Name = "Zeiss LSM 980"
	preset.Models = []string{"LSM 980"}
	preset.Description = "Zeiss LSM 980 confocal laser scanning microscope"
	preset.References = []string{
		"https://www.zeiss.com/microscopy/us/products/confocal-microscopes/lsm-980.html",
	}
	return preset
}

func illuminaNovaSeqPreset() *InstrumentPreset {
	return &InstrumentPreset{
		ID:             "illumina-novaseq",
		Name:           "Illumina NovaSeq",
		Manufacturer:   "Illumina",
		Models:         []string{"NovaSeq 6000", "NovaSeq X", "NovaSeq X Plus"},
		Description:    "Illumina NovaSeq series high-throughput sequencers",
		InstrumentType: "sequencing",
		DataTypes:      []string{"nucleotide_sequence"},
		FileFormats:    []string{".fastq", ".fq", ".fastq.gz", ".fq.gz"},
		RequiredFields: []FieldRequirement{
			{Name: "format", Type: "string", Description: "File format", Enum: []string{"FASTQ"}, Example: "FASTQ"},
			{Name: "instrument_type", Type: "string", Description: "Instrument type", Enum: []string{"sequencing"}, Example: "sequencing"},
			{Name: "total_reads", Type: "number", Description: "Total number of reads", MinValue: ptr(0.0)},
			{Name: "mean_read_length", Type: "number", Description: "Mean read length in bases", MinValue: ptr(1.0)},
		},
		OptionalFields: []FieldRequirement{
			{Name: "is_paired_end", Type: "boolean", Description: "Paired-end sequencing", Example: true},
			{Name: "read_pair", Type: "string", Description: "Read pair identifier", Enum: []string{"R1", "R2", "1", "2"}},
			{Name: "mean_quality_score", Type: "number", Description: "Mean Phred quality score", MinValue: ptr(0.0), MaxValue: ptr(100.0)},
			{Name: "gc_content_percent", Type: "number", Description: "GC content percentage", MinValue: ptr(0.0), MaxValue: ptr(100.0)},
		},
		Documentation: "Illumina NovaSeq is a high-throughput sequencing platform for large-scale genomic studies",
		References: []string{
			"https://www.illumina.com/systems/sequencing-platforms/novaseq.html",
		},
	}
}

func illuminaMiSeqPreset() *InstrumentPreset {
	preset := illuminaNovaSeqPreset()
	preset.ID = "illumina-miseq"
	preset.Name = "Illumina MiSeq"
	preset.Models = []string{"MiSeq"}
	preset.Description = "Illumina MiSeq benchtop sequencer"
	preset.References = []string{
		"https://www.illumina.com/systems/sequencing-platforms/miseq.html",
	}
	return preset
}

func illuminaNextSeqPreset() *InstrumentPreset {
	preset := illuminaNovaSeqPreset()
	preset.ID = "illumina-nextseq"
	preset.Name = "Illumina NextSeq"
	preset.Models = []string{"NextSeq 500", "NextSeq 550", "NextSeq 1000", "NextSeq 2000"}
	preset.Description = "Illumina NextSeq series mid-throughput sequencers"
	preset.References = []string{
		"https://www.illumina.com/systems/sequencing-platforms/nextseq.html",
	}
	return preset
}

func genericMicroscopyPreset() *InstrumentPreset {
	return &InstrumentPreset{
		ID:             "generic-microscopy",
		Name:           "Generic Microscopy",
		Manufacturer:   "Various",
		Models:         []string{},
		Description:    "Generic preset for microscopy instruments",
		InstrumentType: "microscopy",
		DataTypes:      []string{"image"},
		FileFormats:    []string{".tif", ".tiff", ".czi", ".nd2", ".lif", ".ome.tif", ".ome.tiff"},
		RequiredFields: []FieldRequirement{
			{Name: "format", Type: "string", Description: "File format"},
			{Name: "instrument_type", Type: "string", Description: "Instrument type", Enum: []string{"microscopy"}},
			{Name: "image_width", Type: "number", Description: "Image width in pixels", MinValue: ptr(1.0)},
			{Name: "image_height", Type: "number", Description: "Image height in pixels", MinValue: ptr(1.0)},
		},
		OptionalFields: []FieldRequirement{
			{Name: "manufacturer", Type: "string", Description: "Instrument manufacturer"},
			{Name: "instrument_model", Type: "string", Description: "Instrument model"},
			{Name: "acquisition_date", Type: "string", Description: "Acquisition date"},
			{Name: "operator", Type: "string", Description: "Operator name"},
			{Name: "objective_magnification", Type: "number", Description: "Objective magnification"},
		},
		Documentation: "Generic preset for microscopy data validation",
	}
}

func genericSequencingPreset() *InstrumentPreset {
	return &InstrumentPreset{
		ID:             "generic-sequencing",
		Name:           "Generic Sequencing",
		Manufacturer:   "Various",
		Models:         []string{},
		Description:    "Generic preset for sequencing instruments",
		InstrumentType: "sequencing",
		DataTypes:      []string{"nucleotide_sequence"},
		FileFormats:    []string{".fastq", ".fq", ".fastq.gz", ".fq.gz"},
		RequiredFields: []FieldRequirement{
			{Name: "format", Type: "string", Description: "File format", Enum: []string{"FASTQ"}},
			{Name: "instrument_type", Type: "string", Description: "Instrument type", Enum: []string{"sequencing"}},
			{Name: "total_reads", Type: "number", Description: "Total reads", MinValue: ptr(0.0)},
		},
		OptionalFields: []FieldRequirement{
			{Name: "mean_read_length", Type: "number", Description: "Mean read length"},
			{Name: "is_paired_end", Type: "boolean", Description: "Paired-end sequencing"},
			{Name: "mean_quality_score", Type: "number", Description: "Mean quality score"},
			{Name: "gc_content_percent", Type: "number", Description: "GC content percentage"},
		},
		Documentation: "Generic preset for sequencing data validation",
	}
}

// Helper function to create float64 pointers
func ptr(f float64) *float64 {
	return &f
}
