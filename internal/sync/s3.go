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
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/scttfrdmn/cicada/internal/metadata"
)

// S3Backend implements Backend for AWS S3.
type S3Backend struct {
	client *s3.Client
	bucket string
}

// NewS3Backend creates a new S3 backend.
func NewS3Backend(ctx context.Context, bucket string) (*S3Backend, error) {
	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	return &S3Backend{
		client: client,
		bucket: bucket,
	}, nil
}

// List returns all files with the given prefix.
func (b *S3Backend) List(ctx context.Context, prefix string) ([]FileInfo, error) {
	var files []FileInfo

	paginator := s3.NewListObjectsV2Paginator(b.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(b.bucket),
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("list objects: %w", err)
		}

		for _, obj := range page.Contents {
			// Skip directories (keys ending with /)
			if strings.HasSuffix(*obj.Key, "/") {
				continue
			}

			files = append(files, FileInfo{
				Path:         *obj.Key,
				Size:         *obj.Size,
				ModTime:      *obj.LastModified,
				ETag:         strings.Trim(*obj.ETag, "\""), // Remove quotes
				IsDir:        false,
				StorageClass: string(obj.StorageClass),
			})
		}
	}

	return files, nil
}

// Read opens a file for reading.
func (b *S3Backend) Read(ctx context.Context, path string) (io.ReadCloser, error) {
	output, err := b.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}

	return output.Body, nil
}

// Write writes a file.
func (b *S3Backend) Write(ctx context.Context, path string, r io.Reader, size int64) error {
	_, err := b.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(path),
		Body:   r,
	})
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}

	return nil
}

// WriteWithMetadata writes a file with metadata tags.
// The metadata is converted to S3 object tags (max 10 tags, prioritized).
func (b *S3Backend) WriteWithMetadata(ctx context.Context, path string, r io.Reader, size int64, md *metadata.Metadata) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(path),
		Body:   r,
	}

	// Convert metadata to S3 tags if metadata is provided
	if md != nil {
		tags := metadata.MetadataToS3Tags(md)
		if len(tags) > 0 {
			// Convert to tagging string format
			input.Tagging = aws.String(tagsToString(tags))
		}
	}

	_, err := b.client.PutObject(ctx, input)
	if err != nil {
		return fmt.Errorf("put object with metadata: %w", err)
	}

	return nil
}

// PutObjectTagging adds or updates tags on an existing S3 object.
// The metadata is converted to S3 object tags (max 10 tags, prioritized).
func (b *S3Backend) PutObjectTagging(ctx context.Context, path string, md *metadata.Metadata) error {
	if md == nil {
		return fmt.Errorf("metadata cannot be nil")
	}

	tags := metadata.MetadataToS3Tags(md)
	if len(tags) == 0 {
		// No tags to apply
		return nil
	}

	_, err := b.client.PutObjectTagging(ctx, &s3.PutObjectTaggingInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(path),
		Tagging: &types.Tagging{
			TagSet: tags,
		},
	})
	if err != nil {
		return fmt.Errorf("put object tagging: %w", err)
	}

	return nil
}

// GetObjectTagging retrieves tags from an S3 object and converts them to metadata fields.
// Returns a map of field names to values extracted from the tags.
func (b *S3Backend) GetObjectTagging(ctx context.Context, path string) (map[string]interface{}, error) {
	output, err := b.client.GetObjectTagging(ctx, &s3.GetObjectTaggingInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, fmt.Errorf("get object tagging: %w", err)
	}

	return metadata.S3TagsToMetadata(output.TagSet), nil
}

// Delete deletes a file.
func (b *S3Backend) Delete(ctx context.Context, path string) error {
	_, err := b.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("delete object: %w", err)
	}

	return nil
}

// Stat gets file metadata.
func (b *S3Backend) Stat(ctx context.Context, path string) (*FileInfo, error) {
	output, err := b.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, fmt.Errorf("head object: %w", err)
	}

	return &FileInfo{
		Path:         path,
		Size:         *output.ContentLength,
		ModTime:      *output.LastModified,
		ETag:         strings.Trim(*output.ETag, "\""),
		IsDir:        false,
		StorageClass: string(output.StorageClass),
	}, nil
}

// Close closes the backend.
func (b *S3Backend) Close() error {
	return nil // S3 client doesn't need explicit closing
}

// ParseS3URI parses s3://bucket/key into bucket and key.
func ParseS3URI(uri string) (bucket, key string, err error) {
	if !strings.HasPrefix(uri, "s3://") {
		return "", "", fmt.Errorf("invalid S3 URI: must start with s3://")
	}

	parts := strings.SplitN(strings.TrimPrefix(uri, "s3://"), "/", 2)
	if len(parts) < 1 || parts[0] == "" {
		return "", "", fmt.Errorf("invalid S3 URI: missing bucket")
	}

	bucket = parts[0]
	if len(parts) == 2 {
		key = parts[1]
	}

	return bucket, key, nil
}

// tagsToString converts S3 tags to URL-encoded string format for PutObject.
// Format: "key1=value1&key2=value2"
func tagsToString(tags []types.Tag) string {
	var parts []string
	for _, tag := range tags {
		if tag.Key != nil && tag.Value != nil {
			parts = append(parts, fmt.Sprintf("%s=%s", *tag.Key, *tag.Value))
		}
	}
	return strings.Join(parts, "&")
}
