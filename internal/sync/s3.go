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
