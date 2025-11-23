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
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// MaxS3Tags is the maximum number of tags allowed by S3
const MaxS3Tags = 10

// PriorityFields defines the fields to include in S3 tags, in priority order.
// S3 allows a maximum of 10 tags, so we prioritize the most important fields.
var PriorityFields = []string{
	"instrument_type",   // High priority: enables filtering by instrument type
	"format",            // High priority: file format
	"manufacturer",      // High priority: instrument manufacturer
	"instrument_model",  // Medium priority: specific instrument
	"acquisition_date",  // Medium priority: when data was acquired
	"operator",          // Low priority: who ran the instrument
	"extractor_name",    // Low priority: which extractor was used
	"schema_name",       // Low priority: metadata schema used
}

// MetadataToS3Tags converts metadata to S3 tags.
// S3 allows a maximum of 10 tags, so we prioritize the most important fields.
//
// Priority order:
//  1. instrument_type - Enables filtering by instrument type (microscopy, sequencing, etc.)
//  2. format - File format (CZI, OME-TIFF, FASTQ, etc.)
//  3. manufacturer - Instrument manufacturer (Zeiss, Illumina, etc.)
//  4. instrument_model - Specific instrument model
//  5. acquisition_date - When data was acquired
//  6. operator - Who ran the instrument
//  7. extractor_name - Which extractor was used
//  8. schema_name - Metadata schema used
//
// Returns up to 10 tags with sanitized keys and values.
func MetadataToS3Tags(metadata *Metadata) []types.Tag {
	tags := []types.Tag{}

	// Create map of available fields
	fields := make(map[string]string)

	// Add fields from Metadata struct
	if metadata.SchemaName != "" {
		fields["schema_name"] = metadata.SchemaName
	}

	// Add fields from metadata.Fields map
	for key, value := range metadata.Fields {
		if strValue, ok := value.(string); ok && strValue != "" {
			fields[key] = strValue
		} else if value != nil {
			// Convert non-string values to string
			fields[key] = fmt.Sprintf("%v", value)
		}
	}

	// Add file info fields
	if metadata.FileInfo.Format != "" {
		fields["format"] = metadata.FileInfo.Format
	}

	// Prioritize and add tags (max 10)
	for _, fieldName := range PriorityFields {
		if len(tags) >= MaxS3Tags {
			break
		}

		if value, ok := fields[fieldName]; ok && value != "" {
			tags = append(tags, types.Tag{
				Key:   aws.String(sanitizeTagKey(fieldName)),
				Value: aws.String(sanitizeTagValue(value)),
			})
		}
	}

	return tags
}

// S3TagsToMetadata converts S3 tags back to metadata fields.
// This is used when reading metadata from S3 object tags.
func S3TagsToMetadata(tags []types.Tag) map[string]interface{} {
	fields := make(map[string]interface{})

	for _, tag := range tags {
		if tag.Key != nil && tag.Value != nil {
			key := *tag.Key
			value := *tag.Value

			// Convert sanitized key back to field name
			key = strings.ReplaceAll(key, "-", "_")
			key = strings.ToLower(key)

			fields[key] = value
		}
	}

	return fields
}

// sanitizeTagKey sanitizes a metadata field name for use as an S3 tag key.
// S3 tag keys:
// - Can be up to 128 characters
// - Can contain letters, numbers, spaces, and: + - = . _ : / @
// - Are case-sensitive
func sanitizeTagKey(key string) string {
	// Replace underscores with hyphens (more common in tags)
	key = strings.ReplaceAll(key, "_", "-")

	// Remove any invalid characters
	var sanitized strings.Builder
	for _, char := range key {
		if (char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '.' || char == '_' ||
			char == ':' || char == '/' || char == '@' {
			sanitized.WriteRune(char)
		}
	}

	result := sanitized.String()

	// Ensure key is not empty and not too long
	if result == "" {
		result = "unknown"
	}
	if len(result) > 128 {
		result = result[:128]
	}

	return result
}

// sanitizeTagValue sanitizes a metadata value for use as an S3 tag value.
// S3 tag values:
// - Can be up to 256 characters
// - Can contain letters, numbers, spaces, and: + - = . _ : / @
func sanitizeTagValue(value string) string {
	// Remove any invalid characters
	var sanitized strings.Builder
	for _, char := range value {
		if (char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == ' ' || char == '-' || char == '.' || char == '_' ||
			char == ':' || char == '/' || char == '@' || char == '+' || char == '=' {
			sanitized.WriteRune(char)
		}
	}

	result := sanitized.String()

	// Trim whitespace
	result = strings.TrimSpace(result)

	// Ensure value is not empty and not too long
	if result == "" {
		result = "unknown"
	}
	if len(result) > 256 {
		result = result[:256]
	}

	return result
}

