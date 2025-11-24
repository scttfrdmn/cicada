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

	"github.com/scttfrdmn/cicada/internal/metadata"
)

// WorkflowConfig configures the DOI assignment workflow
type WorkflowConfig struct {
	Publisher       string  // Default publisher
	License         string  // Default license
	LandingPageURL  string  // Base URL for landing pages
	MinQualityScore float64 // Minimum quality score (0-100)
	RequireRealAuthors bool // Reject "Unknown Creator"
	RequireDescription bool // Require non-empty description
	AutoEnrich      bool    // Automatically enrich from presets
}

// DOIWorkflow orchestrates the DOI assignment process
type DOIWorkflow struct {
	config    *WorkflowConfig
	mapper    *MetadataMapper
	validator *DOIReadinessValidator
	registry  *ProviderRegistry
}

// NewDOIWorkflow creates a new DOI workflow
func NewDOIWorkflow(config *WorkflowConfig, registry *ProviderRegistry) *DOIWorkflow {
	if config == nil {
		config = &WorkflowConfig{
			Publisher:       "Unknown Publisher",
			License:         "CC-BY-4.0",
			MinQualityScore: 60.0,
			RequireRealAuthors: true,
			RequireDescription: true,
			AutoEnrich:      false,
		}
	}

	mapper := NewMetadataMapper(config.Publisher, config.License, config.LandingPageURL)

	validator := NewDOIReadinessValidator()
	validator.MinQualityScore = config.MinQualityScore
	validator.RequireRealAuthors = config.RequireRealAuthors
	validator.RequireDescription = config.RequireDescription

	return &DOIWorkflow{
		config:    config,
		mapper:    mapper,
		validator: validator,
		registry:  registry,
	}
}

// PrepareRequest represents a DOI preparation request
type PrepareRequest struct {
	FilePath   string                 // Path to the file
	Metadata   map[string]interface{} // Extracted metadata
	Enrichment map[string]interface{} // User-provided enrichment
	PresetID   string                 // Optional instrument preset ID
}

// PrepareResult represents the result of DOI preparation
type PrepareResult struct {
	Dataset    *Dataset         // Mapped dataset
	Validation *ReadinessResult // Validation result
	Warnings   []string         // Workflow warnings
}

// Prepare prepares metadata for DOI minting
func (w *DOIWorkflow) Prepare(req *PrepareRequest) (*PrepareResult, error) {
	result := &PrepareResult{
		Warnings: []string{},
	}

	// 1. Map metadata to Dataset
	dataset, err := w.mapper.MapToDataset(req.Metadata, req.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to map metadata: %w", err)
	}

	// 2. Enrich with preset information if configured
	if w.config.AutoEnrich && req.PresetID != "" {
		presetRegistry := metadata.NewPresetRegistry()
		presetRegistry.RegisterDefaults()

		preset, err := presetRegistry.GetPreset(req.PresetID)
		if err == nil && preset != nil {
			// Enrich dataset with preset metadata
			// This is a placeholder for future preset integration
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("Preset enrichment (%s) not yet implemented", req.PresetID))
		} else if err != nil {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("Failed to load preset %s: %v", req.PresetID, err))
		}
	}

	// 3. Apply user enrichment
	if req.Enrichment != nil {
		w.mapper.EnrichDataset(dataset, req.Enrichment)
	}

	// 4. Validate readiness
	validation := w.validator.Validate(dataset)

	result.Dataset = dataset
	result.Validation = validation

	// 5. Add workflow-specific warnings
	if !validation.IsReady {
		result.Warnings = append(result.Warnings,
			"Dataset is not ready for DOI minting. See validation errors.")
	}

	return result, nil
}

// MintRequest represents a DOI minting request
type MintRequest struct {
	Dataset     *Dataset // Prepared dataset
	Provider    string   // Provider name (datacite, zenodo, etc.)
	Draft       bool     // Create as draft (not published)
	DryRun      bool     // Validate only, don't actually mint
}

// MintResult represents the result of DOI minting
type MintResult struct {
	DOI        *DOI     // Minted DOI (nil if DryRun)
	Cost       float64  // Estimated cost
	Currency   string   // Currency (USD, EUR, etc.)
	Validation *ReadinessResult // Final validation
	Warnings   []string // Minting warnings
}

// Mint mints a DOI for the dataset
func (w *DOIWorkflow) Mint(ctx context.Context, req *MintRequest) (*MintResult, error) {
	result := &MintResult{
		Warnings: []string{},
	}

	// 1. Final validation
	validation := w.validator.Validate(req.Dataset)
	result.Validation = validation

	if !validation.IsReady && !req.DryRun {
		return nil, fmt.Errorf("dataset is not ready for DOI minting: %d errors", len(validation.Errors))
	}

	// 2. Get provider
	var provider Provider
	if req.Provider != "" {
		if p, ok := w.registry.Get(req.Provider); ok {
			provider = p
		} else {
			return nil, fmt.Errorf("provider not found: %s", req.Provider)
		}
	} else {
		provider = w.registry.GetActive()
		if provider == nil {
			return nil, fmt.Errorf("no active DOI provider configured")
		}
	}

	if !provider.IsEnabled() {
		return nil, fmt.Errorf("provider %s is not enabled", provider.Name())
	}

	// 3. Estimate cost
	cost, currency, err := provider.EstimateCost(req.Dataset)
	if err != nil {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("Failed to estimate cost: %v", err))
	}
	result.Cost = cost
	result.Currency = currency

	// 4. Dry run check
	if req.DryRun {
		result.Warnings = append(result.Warnings, "Dry run - DOI not actually minted")
		return result, nil
	}

	// 5. Validate with provider
	if err := provider.Validate(req.Dataset); err != nil {
		return nil, fmt.Errorf("provider validation failed: %w", err)
	}

	// 6. Mint DOI
	doi, err := provider.Mint(ctx, req.Dataset)
	if err != nil {
		return nil, fmt.Errorf("failed to mint DOI: %w", err)
	}

	result.DOI = doi

	// 7. Add success info
	if req.Draft {
		result.Warnings = append(result.Warnings, "DOI created as draft - not yet published")
	}

	return result, nil
}

// PrepareAndMint is a convenience method that combines Prepare and Mint
func (w *DOIWorkflow) PrepareAndMint(ctx context.Context, prepareReq *PrepareRequest, mintReq *MintRequest) (*MintResult, error) {
	// Prepare
	prepResult, err := w.Prepare(prepareReq)
	if err != nil {
		return nil, fmt.Errorf("prepare failed: %w", err)
	}

	if !prepResult.Validation.IsReady {
		return nil, fmt.Errorf("dataset is not ready: %d errors", len(prepResult.Validation.Errors))
	}

	// Set dataset in mint request
	mintReq.Dataset = prepResult.Dataset

	// Mint
	mintResult, err := w.Mint(ctx, mintReq)
	if err != nil {
		return nil, fmt.Errorf("mint failed: %w", err)
	}

	// Merge warnings
	mintResult.Warnings = append(prepResult.Warnings, mintResult.Warnings...)

	return mintResult, nil
}

// ValidateMetadata validates metadata without minting
func (w *DOIWorkflow) ValidateMetadata(metadata map[string]interface{}, filename string) (*PrepareResult, error) {
	dataset, err := w.mapper.MapToDataset(metadata, filename)
	if err != nil {
		return nil, err
	}

	validation := w.validator.Validate(dataset)

	return &PrepareResult{
		Dataset:    dataset,
		Validation: validation,
		Warnings:   []string{},
	}, nil
}

// GetRecommendations returns recommendations for improving metadata
func (w *DOIWorkflow) GetRecommendations(validation *ReadinessResult) []string {
	return w.validator.GetRecommendations(validation)
}

// UpdateDataset updates a Dataset with new information
func (w *DOIWorkflow) UpdateDataset(dataset *Dataset, updates map[string]interface{}) {
	w.mapper.EnrichDataset(dataset, updates)
}

// EstimateCost estimates the cost of minting a DOI with a specific provider
func (w *DOIWorkflow) EstimateCost(dataset *Dataset, providerName string) (float64, string, error) {
	provider, ok := w.registry.Get(providerName)
	if !ok {
		return 0, "", fmt.Errorf("provider not found: %s", providerName)
	}

	return provider.EstimateCost(dataset)
}

// ListProviders returns all registered providers
func (w *DOIWorkflow) ListProviders() []Provider {
	return w.registry.List()
}

// GetActiveProvider returns the currently active provider
func (w *DOIWorkflow) GetActiveProvider() Provider {
	return w.registry.GetActive()
}

// SetActiveProvider sets the active provider by name
func (w *DOIWorkflow) SetActiveProvider(name string) error {
	return w.registry.SetActive(name)
}

// PreviewDataCiteXML generates DataCite XML for preview
func (w *DOIWorkflow) PreviewDataCiteXML(dataset *Dataset) (string, error) {
	generator := &StandardMetadataGenerator{}
	dcMetadata := generator.Generate(dataset)

	// Set temporary DOI for preview
	dcMetadata.Identifier = Identifier{
		Value: "10.XXXXX/preview",
		Type:  "DOI",
	}

	xmlBytes, err := dcMetadata.ToXML()
	if err != nil {
		return "", fmt.Errorf("failed to generate XML: %w", err)
	}

	return string(xmlBytes), nil
}

// ToXML converts DataCiteMetadata to XML
func (m *DataCiteMetadata) ToXML() ([]byte, error) {
	// This method should be in datacite.go, but adding here for workflow completeness
	// In practice, use xml.MarshalIndent
	return nil, fmt.Errorf("ToXML not implemented - use xml.MarshalIndent")
}
