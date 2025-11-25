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

package sync

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func TestParseS3URI(t *testing.T) {
	tests := []struct {
		name        string
		uri         string
		wantBucket  string
		wantKey     string
		wantErr     bool
	}{
		{
			name:       "bucket only",
			uri:        "s3://my-bucket",
			wantBucket: "my-bucket",
			wantKey:    "",
			wantErr:    false,
		},
		{
			name:       "bucket with key",
			uri:        "s3://my-bucket/path/to/file.txt",
			wantBucket: "my-bucket",
			wantKey:    "path/to/file.txt",
			wantErr:    false,
		},
		{
			name:       "bucket with prefix",
			uri:        "s3://my-bucket/prefix/",
			wantBucket: "my-bucket",
			wantKey:    "prefix/",
			wantErr:    false,
		},
		{
			name:       "invalid - no s3 prefix",
			uri:        "http://my-bucket/file.txt",
			wantBucket: "",
			wantKey:    "",
			wantErr:    true,
		},
		{
			name:       "invalid - missing bucket",
			uri:        "s3://",
			wantBucket: "",
			wantKey:    "",
			wantErr:    true,
		},
		{
			name:       "invalid - empty bucket",
			uri:        "s3:///path/to/file.txt",
			wantBucket: "",
			wantKey:    "",
			wantErr:    true,
		},
		{
			name:       "bucket with nested path",
			uri:        "s3://my-bucket/a/b/c/d/file.txt",
			wantBucket: "my-bucket",
			wantKey:    "a/b/c/d/file.txt",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bucket, key, err := ParseS3URI(tt.uri)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseS3URI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if bucket != tt.wantBucket {
				t.Errorf("ParseS3URI() bucket = %v, want %v", bucket, tt.wantBucket)
			}

			if key != tt.wantKey {
				t.Errorf("ParseS3URI() key = %v, want %v", key, tt.wantKey)
			}
		})
	}
}

func TestTagsToString(t *testing.T) {
	tests := []struct {
		name string
		tags []types.Tag
		want string
	}{
		{
			name: "empty tags",
			tags: []types.Tag{},
			want: "",
		},
		{
			name: "single tag",
			tags: []types.Tag{
				{Key: aws.String("format"), Value: aws.String("CZI")},
			},
			want: "format=CZI",
		},
		{
			name: "multiple tags",
			tags: []types.Tag{
				{Key: aws.String("format"), Value: aws.String("CZI")},
				{Key: aws.String("manufacturer"), Value: aws.String("Zeiss")},
				{Key: aws.String("instrument-type"), Value: aws.String("microscopy")},
			},
			want: "format=CZI&manufacturer=Zeiss&instrument-type=microscopy",
		},
		{
			name: "tags with special characters",
			tags: []types.Tag{
				{Key: aws.String("operator"), Value: aws.String("John Doe")},
				{Key: aws.String("acquisition-date"), Value: aws.String("2025-11-23")},
			},
			want: "operator=John Doe&acquisition-date=2025-11-23",
		},
		{
			name: "tag with nil key",
			tags: []types.Tag{
				{Key: aws.String("format"), Value: aws.String("CZI")},
				{Key: nil, Value: aws.String("value")},
			},
			want: "format=CZI",
		},
		{
			name: "tag with nil value",
			tags: []types.Tag{
				{Key: aws.String("format"), Value: aws.String("CZI")},
				{Key: aws.String("empty"), Value: nil},
			},
			want: "format=CZI",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tagsToString(tt.tags)
			if got != tt.want {
				t.Errorf("tagsToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Note: WriteWithMetadata, PutObjectTagging, and GetObjectTagging require
// actual S3 client operations and are tested in integration tests.
// See internal/integration/s3_test.go for integration tests with real S3/LocalStack.
