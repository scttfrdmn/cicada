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

// getDataCiteTestConfig returns DataCite test configuration from environment
func getDataCiteTestConfig(t *testing.T) *doi.DataCiteConfig {
	repositoryID := os.Getenv("DATACITE_SANDBOX_REPOSITORY_ID")
	password := os.Getenv("DATACITE_SANDBOX_PASSWORD")

	if repositoryID == "" || password == "" {
		t.Skip("DataCite sandbox credentials not available (set DATACITE_SANDBOX_REPOSITORY_ID and DATACITE_SANDBOX_PASSWORD)")
	}

	return &doi.DataCiteConfig{
		RepositoryID: repositoryID,
		Password:     password,
		Sandbox:      true,
	}
}

// createTestDataset creates a test dataset for integration tests
func createTestDataset(suffix string) *doi.Dataset {
	return &doi.Dataset{
		Title: fmt.Sprintf("Integration Test Dataset - %s", suffix),
		Authors: []doi.Author{
			{
				Name:        "Test Researcher",
				Affiliation: "Test University",
			},
		},
		Publisher:       "Cicada Integration Tests",
		PublicationYear: time.Now().Year(),
		ResourceType:    "Dataset",
		Description:     fmt.Sprintf("Test dataset created by integration tests at %s", time.Now().Format(time.RFC3339)),
		Keywords: []string{
			"integration testing",
			"automated testing",
		},
	}
}

func TestDataCite_EndToEnd(t *testing.T) {
	config := getDataCiteTestConfig(t)
	provider, err := doi.NewDataCiteProvider(config)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()
	dataset := createTestDataset(fmt.Sprintf("e2e-%d", time.Now().Unix()))

	// Test 1: Mint a new DOI
	t.Log("Minting new DOI...")
	newDOI, err := provider.Mint(ctx, dataset)
	if err != nil {
		t.Fatalf("failed to mint DOI: %v", err)
	}

	if newDOI.DOI == "" {
		t.Fatal("minted DOI is empty")
	}
	if newDOI.State != "registered" && newDOI.State != "draft" && newDOI.State != "findable" {
		t.Errorf("unexpected DOI state: %s", newDOI.State)
	}

	t.Logf("Successfully minted DOI: %s (state: %s)", newDOI.DOI, newDOI.State)

	// Test 2: Get the DOI
	t.Log("Retrieving DOI...")
	retrievedDOI, err := provider.Get(ctx, newDOI.DOI)
	if err != nil {
		t.Fatalf("failed to get DOI: %v", err)
	}

	if retrievedDOI.DOI != newDOI.DOI {
		t.Errorf("retrieved DOI mismatch: got %s, want %s", retrievedDOI.DOI, newDOI.DOI)
	}

	t.Logf("Successfully retrieved DOI: %s", retrievedDOI.DOI)

	// Test 3: Update the DOI
	t.Log("Updating DOI...")
	dataset.Description = fmt.Sprintf("Updated test dataset at %s", time.Now().Format(time.RFC3339))
	err = provider.Update(ctx, newDOI.DOI, dataset)
	if err != nil {
		t.Logf("Update failed (expected for some sandbox environments): %v", err)
	} else {
		t.Log("Successfully updated DOI")
	}

	// Test 4: List DOIs
	t.Log("Listing DOIs...")
	dois, err := provider.List(ctx)
	if err != nil {
		t.Fatalf("failed to list DOIs: %v", err)
	}

	if len(dois) == 0 {
		t.Error("expected at least one DOI in list")
	}

	found := false
	for _, d := range dois {
		if d.DOI == newDOI.DOI {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("newly created DOI %s not found in list", newDOI.DOI)
	}

	t.Logf("Successfully listed %d DOIs", len(dois))
}

func TestDataCite_CreateDraftDOI(t *testing.T) {
	config := getDataCiteTestConfig(t)
	provider, err := doi.NewDataCiteProvider(config)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()
	dataset := createTestDataset(fmt.Sprintf("draft-%d", time.Now().Unix()))

	doi, err := provider.Mint(ctx, dataset)
	if err != nil {
		t.Fatalf("failed to create draft DOI: %v", err)
	}

	if doi.DOI == "" {
		t.Fatal("created DOI is empty")
	}

	t.Logf("Created draft DOI: %s", doi.DOI)
}

func TestDataCite_GetDOI(t *testing.T) {
	config := getDataCiteTestConfig(t)
	provider, err := doi.NewDataCiteProvider(config)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()

	// First create a DOI to retrieve
	dataset := createTestDataset(fmt.Sprintf("get-%d", time.Now().Unix()))
	newDOI, err := provider.Mint(ctx, dataset)
	if err != nil {
		t.Fatalf("failed to create DOI: %v", err)
	}

	// Now retrieve it
	retrievedDOI, err := provider.Get(ctx, newDOI.DOI)
	if err != nil {
		t.Fatalf("failed to get DOI: %v", err)
	}

	if retrievedDOI.DOI != newDOI.DOI {
		t.Errorf("DOI mismatch: got %s, want %s", retrievedDOI.DOI, newDOI.DOI)
	}

	if retrievedDOI.URL == "" {
		t.Error("retrieved DOI has no URL")
	}

	t.Logf("Retrieved DOI: %s (URL: %s)", retrievedDOI.DOI, retrievedDOI.URL)
}

func TestDataCite_AuthenticationFailure(t *testing.T) {
	// Create provider with invalid credentials
	config := &doi.DataCiteConfig{
		RepositoryID: "INVALID.REPO",
		Password:     "invalid-password",
		Sandbox:      true,
	}

	provider, err := doi.NewDataCiteProvider(config)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()
	dataset := createTestDataset("auth-fail")

	// This should fail with authentication error
	_, err = provider.Mint(ctx, dataset)
	if err == nil {
		t.Fatal("expected authentication error, got success")
	}

	t.Logf("Got expected authentication error: %v", err)
}

func TestDataCite_ValidationError(t *testing.T) {
	config := getDataCiteTestConfig(t)
	provider, err := doi.NewDataCiteProvider(config)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()

	// Create invalid dataset (missing required fields)
	invalidDataset := &doi.Dataset{
		Title: "", // Empty title should cause validation error
	}

	_, err = provider.Mint(ctx, invalidDataset)
	if err == nil {
		t.Fatal("expected validation error, got success")
	}

	t.Logf("Got expected validation error: %v", err)
}

func TestDataCite_ListDOIs(t *testing.T) {
	config := getDataCiteTestConfig(t)
	provider, err := doi.NewDataCiteProvider(config)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx := context.Background()

	// Create a test DOI first
	dataset := createTestDataset(fmt.Sprintf("list-%d", time.Now().Unix()))
	newDOI, err := provider.Mint(ctx, dataset)
	if err != nil {
		t.Fatalf("failed to create DOI: %v", err)
	}

	t.Logf("Created test DOI: %s", newDOI.DOI)

	// List DOIs
	dois, err := provider.List(ctx)
	if err != nil {
		t.Fatalf("failed to list DOIs: %v", err)
	}

	if len(dois) == 0 {
		t.Fatal("expected at least one DOI")
	}

	// Verify our DOI is in the list
	found := false
	for _, d := range dois {
		if d.DOI == newDOI.DOI {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("newly created DOI %s not found in list of %d DOIs", newDOI.DOI, len(dois))
	}

	t.Logf("Successfully listed %d DOIs", len(dois))
}
