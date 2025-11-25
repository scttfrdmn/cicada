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

//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/scttfrdmn/cicada/internal/doi"
)

// getZenodoTestConfig returns Zenodo test configuration from environment
func getZenodoTestConfig(t *testing.T) *doi.ZenodoConfig {
	token := os.Getenv("ZENODO_SANDBOX_TOKEN")

	if token == "" {
		t.Skip("Zenodo sandbox token not available (set ZENODO_SANDBOX_TOKEN)")
	}

	return &doi.ZenodoConfig{
		Token:   token,
		Sandbox: true,
	}
}

func TestZenodo_EndToEnd(t *testing.T) {
	config := getZenodoTestConfig(t)
	provider, err := doi.NewZenodoProvider(config)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()
	dataset := createTestDataset(fmt.Sprintf("zenodo-e2e-%d", time.Now().Unix()))

	// Test 1: Create a new deposition (mints DOI)
	t.Log("Creating new deposition...")
	newDOI, err := provider.Mint(ctx, dataset)
	if err != nil {
		t.Fatalf("failed to create deposition: %v", err)
	}

	if newDOI.DOI == "" {
		t.Fatal("created DOI is empty")
	}
	if newDOI.State != "draft" {
		t.Errorf("unexpected DOI state: %s (expected draft)", newDOI.State)
	}

	t.Logf("Successfully created deposition with DOI: %s (state: %s)", newDOI.DOI, newDOI.State)

	// Test 2: Get the deposition
	t.Log("Retrieving deposition...")
	retrievedDOI, err := provider.Get(ctx, newDOI.DOI)
	if err != nil {
		t.Fatalf("failed to get deposition: %v", err)
	}

	if retrievedDOI.DOI != newDOI.DOI {
		t.Errorf("retrieved DOI mismatch: got %s, want %s", retrievedDOI.DOI, newDOI.DOI)
	}

	t.Logf("Successfully retrieved deposition: %s", retrievedDOI.DOI)

	// Test 3: List depositions
	t.Log("Listing depositions...")
	depositions, err := provider.List(ctx)
	if err != nil {
		t.Fatalf("failed to list depositions: %v", err)
	}

	if len(depositions) == 0 {
		t.Error("expected at least one deposition in list")
	}

	found := false
	for _, d := range depositions {
		if d.DOI == newDOI.DOI {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("newly created deposition %s not found in list", newDOI.DOI)
	}

	t.Logf("Successfully listed %d depositions", len(depositions))
}

func TestZenodo_CreateDeposition(t *testing.T) {
	config := getZenodoTestConfig(t)
	provider, err := doi.NewZenodoProvider(config)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()
	dataset := createTestDataset(fmt.Sprintf("zenodo-create-%d", time.Now().Unix()))

	deposition, err := provider.Mint(ctx, dataset)
	if err != nil {
		t.Fatalf("failed to create deposition: %v", err)
	}

	if deposition.DOI == "" {
		t.Fatal("created deposition DOI is empty")
	}

	if deposition.State != "draft" {
		t.Errorf("expected draft state, got %s", deposition.State)
	}

	t.Logf("Created deposition: %s", deposition.DOI)
}

func TestZenodo_GetDeposition(t *testing.T) {
	config := getZenodoTestConfig(t)
	provider, err := doi.NewZenodoProvider(config)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()

	// First create a deposition to retrieve
	dataset := createTestDataset(fmt.Sprintf("zenodo-get-%d", time.Now().Unix()))
	newDOI, err := provider.Mint(ctx, dataset)
	if err != nil {
		t.Fatalf("failed to create deposition: %v", err)
	}

	// Now retrieve it
	retrievedDOI, err := provider.Get(ctx, newDOI.DOI)
	if err != nil {
		t.Fatalf("failed to get deposition: %v", err)
	}

	if retrievedDOI.DOI != newDOI.DOI {
		t.Errorf("DOI mismatch: got %s, want %s", retrievedDOI.DOI, newDOI.DOI)
	}

	if retrievedDOI.URL == "" {
		t.Error("retrieved deposition has no URL")
	}

	t.Logf("Retrieved deposition: %s (URL: %s)", retrievedDOI.DOI, retrievedDOI.URL)
}

func TestZenodo_AuthenticationFailure(t *testing.T) {
	// Create provider with invalid token
	config := &doi.ZenodoConfig{
		Token:   "invalid-token-12345",
		Sandbox: true,
	}

	provider, err := doi.NewZenodoProvider(config)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()
	dataset := createTestDataset("zenodo-auth-fail")

	// This should fail with authentication error
	_, err = provider.Mint(ctx, dataset)
	if err == nil {
		t.Fatal("expected authentication error, got success")
	}

	t.Logf("Got expected authentication error: %v", err)
}

func TestZenodo_ListDepositions(t *testing.T) {
	config := getZenodoTestConfig(t)
	provider, err := doi.NewZenodoProvider(config)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()

	// Create a test deposition first
	dataset := createTestDataset(fmt.Sprintf("zenodo-list-%d", time.Now().Unix()))
	newDOI, err := provider.Mint(ctx, dataset)
	if err != nil {
		t.Fatalf("failed to create deposition: %v", err)
	}

	t.Logf("Created test deposition: %s", newDOI.DOI)

	// List depositions
	depositions, err := provider.List(ctx)
	if err != nil {
		t.Fatalf("failed to list depositions: %v", err)
	}

	if len(depositions) == 0 {
		t.Fatal("expected at least one deposition")
	}

	// Verify our deposition is in the list
	found := false
	for _, d := range depositions {
		if d.DOI == newDOI.DOI {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("newly created deposition %s not found in list of %d depositions", newDOI.DOI, len(depositions))
	}

	t.Logf("Successfully listed %d depositions", len(depositions))
}

func TestZenodo_InvalidMetadata(t *testing.T) {
	config := getZenodoTestConfig(t)
	provider, err := doi.NewZenodoProvider(config)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()

	// Create invalid dataset (empty title)
	invalidDataset := &doi.Dataset{
		Title:           "", // Empty title
		Authors:         []doi.Author{},
		PublicationYear: 0,
	}

	_, err = provider.Mint(ctx, invalidDataset)
	if err == nil {
		t.Fatal("expected validation error for invalid metadata, got success")
	}

	t.Logf("Got expected validation error: %v", err)
}

func TestZenodo_GetNonExistentDOI(t *testing.T) {
	config := getZenodoTestConfig(t)
	provider, err := doi.NewZenodoProvider(config)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()

	// Try to get a DOI that doesn't exist
	nonExistentDOI := "10.5281/zenodo.99999999999"
	_, err = provider.Get(ctx, nonExistentDOI)
	if err == nil {
		t.Fatal("expected error for non-existent DOI, got success")
	}

	t.Logf("Got expected error for non-existent DOI: %v", err)
}
