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
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestNewAuthenticationError(t *testing.T) {
	err := NewAuthenticationError(http.StatusUnauthorized, "invalid token")

	if err.Type != ErrorTypeAuthentication {
		t.Errorf("expected type %s, got %s", ErrorTypeAuthentication, err.Type)
	}
	if err.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, err.StatusCode)
	}
	if err.Retryable {
		t.Error("authentication errors should not be retryable")
	}
	if err.Suggestion == "" {
		t.Error("expected suggestion for authentication error")
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError(http.StatusBadRequest, "missing required field")

	if err.Type != ErrorTypeValidation {
		t.Errorf("expected type %s, got %s", ErrorTypeValidation, err.Type)
	}
	if err.Retryable {
		t.Error("validation errors should not be retryable")
	}
}

func TestNewRateLimitError(t *testing.T) {
	retryAfter := 60 * time.Second
	err := NewRateLimitError(http.StatusTooManyRequests, "rate limit exceeded", retryAfter)

	if err.Type != ErrorTypeRateLimit {
		t.Errorf("expected type %s, got %s", ErrorTypeRateLimit, err.Type)
	}
	if !err.Retryable {
		t.Error("rate limit errors should be retryable")
	}
	if err.Suggestion == "" {
		t.Error("expected suggestion with retry time")
	}
}

func TestNewNetworkError(t *testing.T) {
	underlying := errors.New("connection refused")
	err := NewNetworkError(underlying)

	if err.Type != ErrorTypeNetwork {
		t.Errorf("expected type %s, got %s", ErrorTypeNetwork, err.Type)
	}
	if !err.Retryable {
		t.Error("network errors should be retryable")
	}
	if err.Underlying != underlying {
		t.Error("underlying error not preserved")
	}
}

func TestNewAPIError_Classification(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantType   ErrorType
		wantRetry  bool
	}{
		{
			name:       "401 unauthorized",
			statusCode: http.StatusUnauthorized,
			wantType:   ErrorTypeAuthentication,
			wantRetry:  false,
		},
		{
			name:       "403 forbidden",
			statusCode: http.StatusForbidden,
			wantType:   ErrorTypeAuthentication,
			wantRetry:  false,
		},
		{
			name:       "400 bad request",
			statusCode: http.StatusBadRequest,
			wantType:   ErrorTypeValidation,
			wantRetry:  false,
		},
		{
			name:       "422 unprocessable entity",
			statusCode: http.StatusUnprocessableEntity,
			wantType:   ErrorTypeValidation,
			wantRetry:  false,
		},
		{
			name:       "429 too many requests",
			statusCode: http.StatusTooManyRequests,
			wantType:   ErrorTypeRateLimit,
			wantRetry:  true,
		},
		{
			name:       "500 internal server error",
			statusCode: http.StatusInternalServerError,
			wantType:   ErrorTypeAPI,
			wantRetry:  true,
		},
		{
			name:       "502 bad gateway",
			statusCode: http.StatusBadGateway,
			wantType:   ErrorTypeAPI,
			wantRetry:  true,
		},
		{
			name:       "503 service unavailable",
			statusCode: http.StatusServiceUnavailable,
			wantType:   ErrorTypeAPI,
			wantRetry:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAPIError(tt.statusCode, "test error")

			if err.Type != tt.wantType {
				t.Errorf("expected type %s, got %s", tt.wantType, err.Type)
			}
			if err.Retryable != tt.wantRetry {
				t.Errorf("expected retryable %v, got %v", tt.wantRetry, err.Retryable)
			}
			if err.Suggestion == "" {
				t.Error("expected suggestion for error")
			}
		})
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		Type:       ErrorTypeValidation,
		StatusCode: http.StatusBadRequest,
		Message:    "invalid input",
		Suggestion: "fix the fields",
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("error message should not be empty")
	}

	// Should contain type, status code, message, and suggestion
	expectedSubstrings := []string{"validation", "400", "invalid input", "fix the fields"}
	for _, substr := range expectedSubstrings {
		if !contains(errMsg, substr) {
			t.Errorf("error message should contain %q, got: %s", substr, errMsg)
		}
	}
}

func TestAPIError_Unwrap(t *testing.T) {
	underlying := errors.New("underlying error")
	err := &APIError{
		Type:       ErrorTypeNetwork,
		Message:    "network failure",
		Underlying: underlying,
	}

	unwrapped := err.Unwrap()
	if unwrapped != underlying {
		t.Error("unwrap should return underlying error")
	}
}

func TestWithRetry_Success(t *testing.T) {
	ctx := context.Background()
	config := &RetryConfig{
		MaxRetries:   3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
		Jitter:       false,
	}

	attempts := 0
	fn := func() error {
		attempts++
		return nil
	}

	err := WithRetry(ctx, config, fn)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}

func TestWithRetry_EventualSuccess(t *testing.T) {
	ctx := context.Background()
	config := &RetryConfig{
		MaxRetries:   3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
		Jitter:       false,
	}

	attempts := 0
	fn := func() error {
		attempts++
		if attempts < 3 {
			return NewAPIError(http.StatusInternalServerError, "temporary failure")
		}
		return nil
	}

	err := WithRetry(ctx, config, fn)
	if err != nil {
		t.Errorf("expected no error after retries, got %v", err)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestWithRetry_NonRetryableError(t *testing.T) {
	ctx := context.Background()
	config := &RetryConfig{
		MaxRetries:   3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
		Jitter:       false,
	}

	attempts := 0
	fn := func() error {
		attempts++
		return NewAuthenticationError(http.StatusUnauthorized, "invalid credentials")
	}

	err := WithRetry(ctx, config, fn)
	if err == nil {
		t.Error("expected error")
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt (no retries), got %d", attempts)
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Error("error should be APIError")
	}
	if apiErr.Type != ErrorTypeAuthentication {
		t.Errorf("expected authentication error, got %s", apiErr.Type)
	}
}

func TestWithRetry_MaxRetriesExceeded(t *testing.T) {
	ctx := context.Background()
	config := &RetryConfig{
		MaxRetries:   2,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
		Jitter:       false,
	}

	attempts := 0
	fn := func() error {
		attempts++
		return NewAPIError(http.StatusInternalServerError, "persistent failure")
	}

	err := WithRetry(ctx, config, fn)
	if err == nil {
		t.Error("expected error after max retries")
	}
	if attempts != 3 { // Initial + 2 retries
		t.Errorf("expected 3 attempts, got %d", attempts)
	}

	errMsg := err.Error()
	if !contains(errMsg, "max retries exceeded") {
		t.Errorf("error should indicate max retries exceeded, got: %s", errMsg)
	}
}

func TestWithRetry_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	config := &RetryConfig{
		MaxRetries:   5,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     1 * time.Second,
		Multiplier:   2.0,
		Jitter:       false,
	}

	attempts := 0
	fn := func() error {
		attempts++
		if attempts == 2 {
			cancel() // Cancel context on second attempt
		}
		return NewAPIError(http.StatusInternalServerError, "failure")
	}

	err := WithRetry(ctx, config, fn)
	if err == nil {
		t.Error("expected context cancellation error")
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestWithRetry_DefaultConfig(t *testing.T) {
	ctx := context.Background()

	attempts := 0
	fn := func() error {
		attempts++
		if attempts == 1 {
			return NewAPIError(http.StatusInternalServerError, "failure")
		}
		return nil
	}

	err := WithRetry(ctx, nil, fn) // Use default config
	if err != nil {
		t.Errorf("expected no error with default config, got %v", err)
	}
}

func TestCalculateBackoff(t *testing.T) {
	config := &RetryConfig{
		InitialDelay: 1 * time.Second,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
	}

	tests := []struct {
		attempt  int
		expected time.Duration
	}{
		{0, 1 * time.Second},
		{1, 2 * time.Second},
		{2, 4 * time.Second},
		{3, 8 * time.Second},
		{4, 16 * time.Second},
		{5, 30 * time.Second}, // Capped at MaxDelay
		{10, 30 * time.Second}, // Still capped
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("attempt_%d", tt.attempt), func(t *testing.T) {
			delay := CalculateBackoff(tt.attempt, config)
			if delay != tt.expected {
				t.Errorf("attempt %d: expected %v, got %v", tt.attempt, tt.expected, delay)
			}
		})
	}
}

func TestCalculateBackoff_DefaultConfig(t *testing.T) {
	delay := CalculateBackoff(0, nil)
	if delay <= 0 {
		t.Error("delay should be positive with default config")
	}
}

func TestWithRetry_ExponentialBackoff(t *testing.T) {
	ctx := context.Background()
	config := &RetryConfig{
		MaxRetries:   3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
		Jitter:       false,
	}

	attempts := 0
	startTime := time.Now()
	fn := func() error {
		attempts++
		if attempts < 3 {
			return NewAPIError(http.StatusInternalServerError, "temporary failure")
		}
		return nil
	}

	err := WithRetry(ctx, config, fn)
	elapsed := time.Since(startTime)

	if err != nil {
		t.Errorf("expected success after retries, got %v", err)
	}

	// Should have waited at least: 10ms (first retry) + 20ms (second retry) = 30ms
	minExpectedDelay := 30 * time.Millisecond
	if elapsed < minExpectedDelay {
		t.Errorf("expected at least %v delay, got %v", minExpectedDelay, elapsed)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
