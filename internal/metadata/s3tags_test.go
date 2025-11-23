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

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func TestMetadataToS3Tags(t *testing.T) {
	tests := []struct {
		name     string
		metadata *Metadata
		wantLen  int
		wantTags map[string]string // Expected tag key-value pairs
	}{
		{
			name: "empty metadata",
			metadata: &Metadata{
				Fields: map[string]interface{}{},
			},
			wantLen:  0,
			wantTags: map[string]string{},
		},
		{
			name: "metadata with priority fields",
			metadata: &Metadata{
				SchemaName: "zeiss_czi",
				Fields: map[string]interface{}{
					"instrument_type":  "microscopy",
					"manufacturer":     "Zeiss",
					"instrument_model": "LSM 980",
					"acquisition_date": "2025-11-23",
				},
				FileInfo: FileInfo{
					Format: "CZI",
				},
			},
			wantLen: 6, // 5 field tags + schema_name
			wantTags: map[string]string{
				"instrument-type":  "microscopy",
				"format":           "CZI",
				"manufacturer":     "Zeiss",
				"instrument-model": "LSM 980",
				"acquisition-date": "2025-11-23",
				"schema-name":      "zeiss_czi",
			},
		},
		{
			name: "metadata with all priority fields",
			metadata: &Metadata{
				SchemaName: "zeiss_czi",
				Fields: map[string]interface{}{
					"instrument_type":  "microscopy",
					"manufacturer":     "Zeiss",
					"instrument_model": "LSM 980",
					"acquisition_date": "2025-11-23",
					"operator":         "John Doe",
					"extractor_name":   "zeiss_czi",
					"extra_field_1":    "value1", // Not in PriorityFields, won't be included
					"extra_field_2":    "value2", // Not in PriorityFields, won't be included
					"extra_field_3":    "value3", // Not in PriorityFields, won't be included
				},
				FileInfo: FileInfo{
					Format: "CZI",
				},
			},
			wantLen: 8, // All 8 fields from PriorityFields
			wantTags: map[string]string{
				"instrument-type":  "microscopy",
				"format":           "CZI",
				"manufacturer":     "Zeiss",
				"instrument-model": "LSM 980",
				"acquisition-date": "2025-11-23",
				"operator":         "John Doe",
				"extractor-name":   "zeiss_czi",
				"schema-name":      "zeiss_czi",
			},
		},
		{
			name: "metadata with non-priority fields only included if in PriorityFields",
			metadata: &Metadata{
				Fields: map[string]interface{}{
					"instrument_type": "microscopy",
					"pixel_count":     1024,        // Not in PriorityFields
					"is_valid":        true,        // Not in PriorityFields
				},
				FileInfo: FileInfo{
					Format: "CZI",
				},
			},
			wantLen: 2, // Only priority fields: instrument_type and format
			wantTags: map[string]string{
				"instrument-type": "microscopy",
				"format":          "CZI",
			},
		},
		{
			name: "metadata with empty values",
			metadata: &Metadata{
				Fields: map[string]interface{}{
					"instrument_type": "microscopy",
					"manufacturer":    "",
					"operator":        nil,
				},
				FileInfo: FileInfo{
					Format: "CZI",
				},
			},
			wantLen: 2,
			wantTags: map[string]string{
				"instrument-type": "microscopy",
				"format":          "CZI",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags := MetadataToS3Tags(tt.metadata)

			if len(tags) != tt.wantLen {
				t.Errorf("MetadataToS3Tags() returned %d tags, want %d", len(tags), tt.wantLen)
			}

			// Verify expected tags are present
			tagMap := make(map[string]string)
			for _, tag := range tags {
				if tag.Key != nil && tag.Value != nil {
					tagMap[*tag.Key] = *tag.Value
				}
			}

			for key, expectedValue := range tt.wantTags {
				if gotValue, ok := tagMap[key]; !ok {
					t.Errorf("MetadataToS3Tags() missing expected tag %q", key)
				} else if gotValue != expectedValue {
					t.Errorf("MetadataToS3Tags() tag %q = %q, want %q", key, gotValue, expectedValue)
				}
			}
		})
	}
}

func TestS3TagsToMetadata(t *testing.T) {
	tests := []struct {
		name      string
		tags      []types.Tag
		wantFields map[string]interface{}
	}{
		{
			name:      "empty tags",
			tags:      []types.Tag{},
			wantFields: map[string]interface{}{},
		},
		{
			name: "tags with hyphens converted to underscores",
			tags: []types.Tag{
				{Key: strPtr("instrument-type"), Value: strPtr("microscopy")},
				{Key: strPtr("format"), Value: strPtr("CZI")},
				{Key: strPtr("acquisition-date"), Value: strPtr("2025-11-23")},
			},
			wantFields: map[string]interface{}{
				"instrument_type":  "microscopy",
				"format":           "CZI",
				"acquisition_date": "2025-11-23",
			},
		},
		{
			name: "tags with uppercase converted to lowercase",
			tags: []types.Tag{
				{Key: strPtr("Format"), Value: strPtr("CZI")},
				{Key: strPtr("Manufacturer"), Value: strPtr("Zeiss")},
			},
			wantFields: map[string]interface{}{
				"format":       "CZI",
				"manufacturer": "Zeiss",
			},
		},
		{
			name: "tags with nil key or value",
			tags: []types.Tag{
				{Key: strPtr("format"), Value: strPtr("CZI")},
				{Key: nil, Value: strPtr("value")},
				{Key: strPtr("key"), Value: nil},
			},
			wantFields: map[string]interface{}{
				"format": "CZI",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := S3TagsToMetadata(tt.tags)

			if len(fields) != len(tt.wantFields) {
				t.Errorf("S3TagsToMetadata() returned %d fields, want %d", len(fields), len(tt.wantFields))
			}

			for key, expectedValue := range tt.wantFields {
				if gotValue, ok := fields[key]; !ok {
					t.Errorf("S3TagsToMetadata() missing expected field %q", key)
				} else if gotValue != expectedValue {
					t.Errorf("S3TagsToMetadata() field %q = %v, want %v", key, gotValue, expectedValue)
				}
			}
		})
	}
}

func TestSanitizeTagKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{
			name: "underscores to hyphens",
			key:  "instrument_type",
			want: "instrument-type",
		},
		{
			name: "valid characters",
			key:  "acquisition-date",
			want: "acquisition-date",
		},
		{
			name: "alphanumeric with special chars",
			key:  "key_with.dots:and/slashes@test",
			want: "key-with.dots:and/slashes@test",
		},
		{
			name: "invalid characters removed",
			key:  "key$with%invalid&chars!",
			want: "keywithinvalidchars",
		},
		{
			name: "empty string",
			key:  "",
			want: "unknown",
		},
		{
			name: "too long",
			key:  "this_is_a_very_long_key_name_that_exceeds_the_maximum_length_allowed_by_s3_which_is_128_characters_and_we_need_to_make_sure_it_gets_truncated_properly",
			want: "this-is-a-very-long-key-name-that-exceeds-the-maximum-length-allowed-by-s3-which-is-128-characters-and-we-need-to-make-sure-it-g",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeTagKey(tt.key)
			if got != tt.want {
				t.Errorf("sanitizeTagKey() = %q, want %q", got, tt.want)
			}
			if len(got) > 128 {
				t.Errorf("sanitizeTagKey() returned key with length %d, exceeds max 128", len(got))
			}
		})
	}
}

func TestSanitizeTagValue(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "valid value",
			value: "Zeiss LSM 980",
			want:  "Zeiss LSM 980",
		},
		{
			name:  "alphanumeric with special chars",
			value: "value-with.dots:and/slashes@test+equals=sign",
			want:  "value-with.dots:and/slashes@test+equals=sign",
		},
		{
			name:  "invalid characters removed",
			value: "value$with%invalid&chars!",
			want:  "valuewithinvalidchars",
		},
		{
			name:  "whitespace trimmed",
			value: "  value with spaces  ",
			want:  "value with spaces",
		},
		{
			name:  "empty string",
			value: "",
			want:  "unknown",
		},
		{
			name:  "too long",
			value: "this_is_a_very_long_value_that_exceeds_the_maximum_length_allowed_by_s3_which_is_256_characters_and_we_need_to_make_sure_it_gets_truncated_properly_so_lets_add_more_text_here_to_make_it_longer_and_longer_until_we_reach_the_limit_and_beyond_to_test_the_truncation",
			want:  "this_is_a_very_long_value_that_exceeds_the_maximum_length_allowed_by_s3_which_is_256_characters_and_we_need_to_make_sure_it_gets_truncated_properly_so_lets_add_more_text_here_to_make_it_longer_and_longer_until_we_reach_the_limit_and_beyond_to_test_the_trun",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeTagValue(tt.value)
			if got != tt.want {
				t.Errorf("sanitizeTagValue() = %q, want %q", got, tt.want)
			}
			if len(got) > 256 {
				t.Errorf("sanitizeTagValue() returned value with length %d, exceeds max 256", len(got))
			}
		})
	}
}

// Helper function for creating string pointers
func strPtr(s string) *string {
	return &s
}
