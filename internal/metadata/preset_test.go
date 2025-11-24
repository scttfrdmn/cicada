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

package metadata

import (
	"testing"
)

func TestPresetRegistry_Register(t *testing.T) {
	registry := NewPresetRegistry()
	preset := &InstrumentPreset{
		ID:           "test-preset",
		Name:         "Test Preset",
		Manufacturer: "Test Corp",
	}

	registry.Register(preset)

	retrieved, err := registry.GetPreset("test-preset")
	if err != nil {
		t.Fatalf("GetPreset() error = %v", err)
	}

	if retrieved.ID != preset.ID {
		t.Errorf("Retrieved preset ID = %v, want %v", retrieved.ID, preset.ID)
	}
}

func TestPresetRegistry_GetPreset_NotFound(t *testing.T) {
	registry := NewPresetRegistry()

	_, err := registry.GetPreset("nonexistent")
	if err == nil {
		t.Error("GetPreset() should return error for nonexistent preset")
	}
}

func TestPresetRegistry_ListPresets(t *testing.T) {
	registry := NewPresetRegistry()

	preset1 := &InstrumentPreset{ID: "preset1", Name: "Preset 1"}
	preset2 := &InstrumentPreset{ID: "preset2", Name: "Preset 2"}

	registry.Register(preset1)
	registry.Register(preset2)

	presets := registry.ListPresets()
	if len(presets) != 2 {
		t.Errorf("ListPresets() returned %d presets, want 2", len(presets))
	}
}

func TestPresetRegistry_FindPresets(t *testing.T) {
	registry := NewPresetRegistry()

	zeissPreset := &InstrumentPreset{
		ID:             "zeiss1",
		Manufacturer:   "Zeiss",
		InstrumentType: "microscopy",
	}
	illuminaPreset := &InstrumentPreset{
		ID:             "illumina1",
		Manufacturer:   "Illumina",
		InstrumentType: "sequencing",
	}

	registry.Register(zeissPreset)
	registry.Register(illuminaPreset)

	// Find by manufacturer
	zeissPresets := registry.FindPresets("Zeiss", "")
	if len(zeissPresets) != 1 {
		t.Errorf("FindPresets('Zeiss', '') returned %d presets, want 1", len(zeissPresets))
	}

	// Find by instrument type
	seqPresets := registry.FindPresets("", "sequencing")
	if len(seqPresets) != 1 {
		t.Errorf("FindPresets('', 'sequencing') returned %d presets, want 1", len(seqPresets))
	}

	// Find by both
	illuminaSeq := registry.FindPresets("Illumina", "sequencing")
	if len(illuminaSeq) != 1 {
		t.Errorf("FindPresets('Illumina', 'sequencing') returned %d presets, want 1", len(illuminaSeq))
	}

	// Find all
	allPresets := registry.FindPresets("", "")
	if len(allPresets) != 2 {
		t.Errorf("FindPresets('', '') returned %d presets, want 2", len(allPresets))
	}
}

func TestPresetRegistry_RegisterDefaults(t *testing.T) {
	registry := NewPresetRegistry()
	registry.RegisterDefaults()

	presets := registry.ListPresets()
	if len(presets) < 8 {
		t.Errorf("RegisterDefaults() registered %d presets, want at least 8", len(presets))
	}

	// Check for specific presets
	expectedIDs := []string{
		"zeiss-lsm-880",
		"zeiss-lsm-900",
		"zeiss-lsm-980",
		"illumina-novaseq",
		"illumina-miseq",
		"illumina-nextseq",
		"generic-microscopy",
		"generic-sequencing",
	}

	for _, id := range expectedIDs {
		_, err := registry.GetPreset(id)
		if err != nil {
			t.Errorf("RegisterDefaults() did not register preset %s", id)
		}
	}
}

func TestInstrumentPreset_Validate_RequiredFields(t *testing.T) {
	preset := &InstrumentPreset{
		ID:   "test",
		Name: "Test",
		RequiredFields: []FieldRequirement{
			{Name: "field1", Type: "string"},
			{Name: "field2", Type: "number"},
		},
	}

	tests := []struct {
		name     string
		metadata map[string]interface{}
		wantValid bool
		wantErrors int
	}{
		{
			name: "all required fields present",
			metadata: map[string]interface{}{
				"field1": "value1",
				"field2": 42.0,
			},
			wantValid: true,
			wantErrors: 0,
		},
		{
			name: "missing required field",
			metadata: map[string]interface{}{
				"field1": "value1",
			},
			wantValid: false,
			wantErrors: 1,
		},
		{
			name: "all required fields missing",
			metadata: map[string]interface{}{},
			wantValid: false,
			wantErrors: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := preset.Validate(tt.metadata)
			if result.IsValid != tt.wantValid {
				t.Errorf("Validate() IsValid = %v, want %v", result.IsValid, tt.wantValid)
			}
			if len(result.Errors) != tt.wantErrors {
				t.Errorf("Validate() errors count = %d, want %d", len(result.Errors), tt.wantErrors)
			}
		})
	}
}

func TestInstrumentPreset_Validate_OptionalFields(t *testing.T) {
	preset := &InstrumentPreset{
		ID:   "test",
		Name: "Test",
		OptionalFields: []FieldRequirement{
			{Name: "optional1", Type: "string"},
			{Name: "optional2", Type: "number"},
		},
	}

	metadata := map[string]interface{}{
		"optional1": "value1",
		// optional2 missing
	}

	result := preset.Validate(metadata)

	// Should be valid (optional fields don't invalidate)
	if !result.IsValid {
		t.Errorf("Validate() IsValid = false, want true for missing optional fields")
	}

	// Should have warnings
	if len(result.Warnings) == 0 {
		t.Error("Validate() should have warnings for missing optional fields")
	}
}

func TestValidateFieldValue_StringType(t *testing.T) {
	tests := []struct {
		name    string
		req     FieldRequirement
		value   interface{}
		wantErr bool
	}{
		{
			name:    "valid string",
			req:     FieldRequirement{Name: "field", Type: "string"},
			value:   "test",
			wantErr: false,
		},
		{
			name:    "invalid type",
			req:     FieldRequirement{Name: "field", Type: "string"},
			value:   42,
			wantErr: true,
		},
		{
			name:    "string with pattern match",
			req:     FieldRequirement{Name: "field", Type: "string", Pattern: "^[A-Z]+$"},
			value:   "ABC",
			wantErr: false,
		},
		{
			name:    "string with pattern no match",
			req:     FieldRequirement{Name: "field", Type: "string", Pattern: "^[A-Z]+$"},
			value:   "abc",
			wantErr: true,
		},
		{
			name:    "string with enum match",
			req:     FieldRequirement{Name: "field", Type: "string", Enum: []string{"A", "B", "C"}},
			value:   "B",
			wantErr: false,
		},
		{
			name:    "string with enum no match",
			req:     FieldRequirement{Name: "field", Type: "string", Enum: []string{"A", "B", "C"}},
			value:   "D",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFieldValue(tt.req, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFieldValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFieldValue_NumberType(t *testing.T) {
	minVal := 0.0
	maxVal := 100.0

	tests := []struct {
		name    string
		req     FieldRequirement
		value   interface{}
		wantErr bool
	}{
		{
			name:    "valid float64",
			req:     FieldRequirement{Name: "field", Type: "number"},
			value:   42.5,
			wantErr: false,
		},
		{
			name:    "valid int",
			req:     FieldRequirement{Name: "field", Type: "number"},
			value:   42,
			wantErr: false,
		},
		{
			name:    "invalid type",
			req:     FieldRequirement{Name: "field", Type: "number"},
			value:   "not a number",
			wantErr: true,
		},
		{
			name:    "number within range",
			req:     FieldRequirement{Name: "field", Type: "number", MinValue: &minVal, MaxValue: &maxVal},
			value:   50.0,
			wantErr: false,
		},
		{
			name:    "number below min",
			req:     FieldRequirement{Name: "field", Type: "number", MinValue: &minVal},
			value:   -10.0,
			wantErr: true,
		},
		{
			name:    "number above max",
			req:     FieldRequirement{Name: "field", Type: "number", MaxValue: &maxVal},
			value:   150.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFieldValue(tt.req, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFieldValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFieldValue_BooleanType(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"valid true", true, false},
		{"valid false", false, false},
		{"invalid string", "true", true},
		{"invalid number", 1, true},
	}

	req := FieldRequirement{Name: "field", Type: "boolean"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFieldValue(req, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFieldValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFieldValue_ArrayType(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"valid []interface{}", []interface{}{"a", "b"}, false},
		{"valid []string", []string{"a", "b"}, false},
		{"valid []int", []int{1, 2}, false},
		{"invalid string", "not an array", true},
		{"invalid number", 42, true},
	}

	req := FieldRequirement{Name: "field", Type: "array"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFieldValue(req, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFieldValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFieldValue_ObjectType(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"valid object", map[string]interface{}{"key": "value"}, false},
		{"invalid string", "not an object", true},
		{"invalid array", []string{"a", "b"}, true},
	}

	req := FieldRequirement{Name: "field", Type: "object"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFieldValue(req, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFieldValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPresetValidationResult_QualityScore(t *testing.T) {
	tests := []struct {
		name       string
		result     *PresetValidationResult
		wantScore  float64
	}{
		{
			name: "all fields present",
			result: &PresetValidationResult{
				Present: []string{"field1", "field2", "field3"},
				Missing: []string{},
				Errors:  []string{},
			},
			wantScore: 100.0,
		},
		{
			name: "half fields present",
			result: &PresetValidationResult{
				Present: []string{"field1"},
				Missing: []string{"field2"},
				Errors:  []string{},
			},
			wantScore: 50.0,
		},
		{
			name: "no fields",
			result: &PresetValidationResult{
				Present: []string{},
				Missing: []string{},
				Errors:  []string{},
			},
			wantScore: 0.0,
		},
		{
			name: "with errors penalty",
			result: &PresetValidationResult{
				Present: []string{"field1", "field2"},
				Missing: []string{},
				Errors:  []string{"error1"},
			},
			wantScore: 90.0, // 100 - (1 error * 10)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := tt.result.QualityScore()
			if score != tt.wantScore {
				t.Errorf("QualityScore() = %v, want %v", score, tt.wantScore)
			}
		})
	}
}

func TestInstrumentPreset_GenerateTemplate(t *testing.T) {
	preset := &InstrumentPreset{
		ID:   "test",
		Name: "Test",
		RequiredFields: []FieldRequirement{
			{Name: "field1", Type: "string", Description: "Field 1", Example: "example1"},
			{Name: "field2", Type: "number", Description: "Field 2", Example: 42},
		},
		OptionalFields: []FieldRequirement{
			{Name: "optional1", Type: "string", Example: "opt1"},
		},
	}

	template := preset.GenerateTemplate()

	// Check required fields with examples
	if template["field1"] != "example1" {
		t.Errorf("Template field1 = %v, want 'example1'", template["field1"])
	}

	if template["field2"] != 42 {
		t.Errorf("Template field2 = %v, want 42", template["field2"])
	}

	// Check optional fields
	if template["optional1"] != "opt1" {
		t.Errorf("Template optional1 = %v, want 'opt1'", template["optional1"])
	}
}

func TestZeissLSM880Preset(t *testing.T) {
	preset := zeissLSM880Preset()

	if preset.ID != "zeiss-lsm-880" {
		t.Errorf("Preset ID = %v, want 'zeiss-lsm-880'", preset.ID)
	}

	if preset.Manufacturer != "Zeiss" {
		t.Errorf("Manufacturer = %v, want 'Zeiss'", preset.Manufacturer)
	}

	if preset.InstrumentType != "microscopy" {
		t.Errorf("InstrumentType = %v, want 'microscopy'", preset.InstrumentType)
	}

	// Test validation with valid CZI metadata
	metadata := map[string]interface{}{
		"format":           "CZI",
		"manufacturer":     "Zeiss",
		"instrument_model": "LSM 880",
		"image_width":      1024,
		"image_height":     1024,
	}

	result := preset.Validate(metadata)
	if !result.IsValid {
		t.Errorf("Valid CZI metadata failed validation: %v", result.Errors)
	}

	// Test validation with missing required field
	incompleteMetadata := map[string]interface{}{
		"format":       "CZI",
		"manufacturer": "Zeiss",
	}

	result2 := preset.Validate(incompleteMetadata)
	if result2.IsValid {
		t.Error("Incomplete metadata should fail validation")
	}
}

func TestIlluminaNovaSeqPreset(t *testing.T) {
	preset := illuminaNovaSeqPreset()

	if preset.ID != "illumina-novaseq" {
		t.Errorf("Preset ID = %v, want 'illumina-novaseq'", preset.ID)
	}

	if preset.Manufacturer != "Illumina" {
		t.Errorf("Manufacturer = %v, want 'Illumina'", preset.Manufacturer)
	}

	if preset.InstrumentType != "sequencing" {
		t.Errorf("InstrumentType = %v, want 'sequencing'", preset.InstrumentType)
	}

	// Test validation with valid FASTQ metadata
	metadata := map[string]interface{}{
		"format":           "FASTQ",
		"instrument_type":  "sequencing",
		"total_reads":      10000,
		"mean_read_length": 150.0,
	}

	result := preset.Validate(metadata)
	if !result.IsValid {
		t.Errorf("Valid FASTQ metadata failed validation: %v", result.Errors)
	}
}

func TestGenericMicroscopyPreset(t *testing.T) {
	preset := genericMicroscopyPreset()

	if preset.ID != "generic-microscopy" {
		t.Errorf("Preset ID = %v, want 'generic-microscopy'", preset.ID)
	}

	// Test validation with minimal microscopy metadata
	metadata := map[string]interface{}{
		"format":          "TIFF",
		"instrument_type": "microscopy",
		"image_width":     512,
		"image_height":    512,
	}

	result := preset.Validate(metadata)
	if !result.IsValid {
		t.Errorf("Valid microscopy metadata failed validation: %v", result.Errors)
	}
}

func TestGenericSequencingPreset(t *testing.T) {
	preset := genericSequencingPreset()

	if preset.ID != "generic-sequencing" {
		t.Errorf("Preset ID = %v, want 'generic-sequencing'", preset.ID)
	}

	// Test validation with minimal sequencing metadata
	metadata := map[string]interface{}{
		"format":          "FASTQ",
		"instrument_type": "sequencing",
		"total_reads":     5000,
	}

	result := preset.Validate(metadata)
	if !result.IsValid {
		t.Errorf("Valid sequencing metadata failed validation: %v", result.Errors)
	}
}
