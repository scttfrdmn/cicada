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
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"
)

// ErrorType represents the category of error
type ErrorType string

const (
	// ErrorTypeAPI represents a general API error
	ErrorTypeAPI ErrorType = "api"
	// ErrorTypeAuthentication represents authentication failures
	ErrorTypeAuthentication ErrorType = "authentication"
	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation ErrorType = "validation"
	// ErrorTypeRateLimit represents rate limiting errors
	ErrorTypeRateLimit ErrorType = "rate_limit"
	// ErrorTypeNetwork represents network connectivity errors
	ErrorTypeNetwork ErrorType = "network"
)

// APIError represents an error from a DOI provider API
type APIError struct {
	Type       ErrorType
	StatusCode int
	Message    string
	Suggestion string
	Underlying error
	Retryable  bool
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Suggestion != "" {
		return fmt.Sprintf("%s (HTTP %d): %s. Suggestion: %s", e.Type, e.StatusCode, e.Message, e.Suggestion)
	}
	return fmt.Sprintf("%s (HTTP %d): %s", e.Type, e.StatusCode, e.Message)
}

// Unwrap implements error unwrapping
func (e *APIError) Unwrap() error {
	return e.Underlying
}

// IsRetryable returns whether this error should trigger a retry
func (e *APIError) IsRetryable() bool {
	return e.Retryable
}

// NewAuthenticationError creates an authentication error
func NewAuthenticationError(statusCode int, message string) *APIError {
	return &APIError{
		Type:       ErrorTypeAuthentication,
		StatusCode: statusCode,
		Message:    message,
		Suggestion: "check your API credentials and ensure they are valid",
		Retryable:  false,
	}
}

// NewValidationError creates a validation error
func NewValidationError(statusCode int, message string) *APIError {
	return &APIError{
		Type:       ErrorTypeValidation,
		StatusCode: statusCode,
		Message:    message,
		Suggestion: "review the validation errors and fix the required fields",
		Retryable:  false,
	}
}

// NewRateLimitError creates a rate limit error
func NewRateLimitError(statusCode int, message string, retryAfter time.Duration) *APIError {
	suggestion := fmt.Sprintf("rate limit exceeded, wait %v before retrying", retryAfter)
	return &APIError{
		Type:       ErrorTypeRateLimit,
		StatusCode: statusCode,
		Message:    message,
		Suggestion: suggestion,
		Retryable:  true,
	}
}

// NewNetworkError creates a network error
func NewNetworkError(err error) *APIError {
	return &APIError{
		Type:       ErrorTypeNetwork,
		StatusCode: 0,
		Message:    err.Error(),
		Suggestion: "check your network connection and try again",
		Underlying: err,
		Retryable:  true,
	}
}

// NewAPIError creates a general API error from an HTTP response
func NewAPIError(statusCode int, message string) *APIError {
	errorType := ErrorTypeAPI
	retryable := false
	suggestion := ""

	// Classify based on status code
	switch {
	case statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden:
		errorType = ErrorTypeAuthentication
		suggestion = "check your API credentials and ensure they are valid"
	case statusCode == http.StatusBadRequest || statusCode == http.StatusUnprocessableEntity:
		errorType = ErrorTypeValidation
		suggestion = "review the request data and ensure all required fields are valid"
	case statusCode == http.StatusTooManyRequests:
		errorType = ErrorTypeRateLimit
		suggestion = "rate limit exceeded, wait before retrying"
		retryable = true
	case statusCode >= 500 && statusCode < 600:
		// Server errors are retryable
		suggestion = "the server encountered an error, will retry automatically"
		retryable = true
	}

	return &APIError{
		Type:       errorType,
		StatusCode: statusCode,
		Message:    message,
		Suggestion: suggestion,
		Retryable:  retryable,
	}
}

// RetryConfig configures retry behavior
type RetryConfig struct {
	// MaxRetries is the maximum number of retry attempts
	MaxRetries int
	// InitialDelay is the initial delay before the first retry
	InitialDelay time.Duration
	// MaxDelay is the maximum delay between retries
	MaxDelay time.Duration
	// Multiplier is the backoff multiplier
	Multiplier float64
	// Jitter adds randomness to retry delays to avoid thundering herd
	Jitter bool
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:   3,
		InitialDelay: 1 * time.Second,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
		Jitter:       true,
	}
}

// RetryFunc is a function that can be retried
type RetryFunc func() error

// WithRetry executes a function with retry logic and exponential backoff
func WithRetry(ctx context.Context, config *RetryConfig, fn RetryFunc) error {
	if config == nil {
		config = DefaultRetryConfig()
	}

	var lastErr error
	delay := config.InitialDelay

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// Execute the function
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if error is retryable
		var apiErr *APIError
		if !isRetryableError(err, &apiErr) {
			return err
		}

		// Don't retry on last attempt
		if attempt == config.MaxRetries {
			break
		}

		// Calculate delay with exponential backoff
		if attempt > 0 {
			delay = time.Duration(float64(delay) * config.Multiplier)
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}
		}

		// Add jitter if enabled
		actualDelay := delay
		if config.Jitter {
			// Add random jitter of ±25%
			jitter := time.Duration(rand.Float64() * float64(delay) * 0.5)
			if rand.Float64() < 0.5 {
				actualDelay = delay - jitter
			} else {
				actualDelay = delay + jitter
			}
		}

		// Wait before retry, respecting context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(actualDelay):
			// Continue to next attempt
		}
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// isRetryableError checks if an error should trigger a retry
func isRetryableError(err error, apiErr **APIError) bool {
	// Check if it's an APIError with Retryable flag
	if apiError, ok := err.(*APIError); ok {
		*apiErr = apiError
		return apiError.Retryable
	}

	// Check for context errors (don't retry)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return false
	}

	// Network errors are generally retryable
	// This is a simplified check - in production you'd want more sophisticated detection
	return true
}

// CalculateBackoff calculates the backoff delay for a given attempt
func CalculateBackoff(attempt int, config *RetryConfig) time.Duration {
	if config == nil {
		config = DefaultRetryConfig()
	}

	delay := float64(config.InitialDelay) * math.Pow(config.Multiplier, float64(attempt))
	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}

	return time.Duration(delay)
}
