// internal/metadata/schema.go
package metadata

import (
	"fmt"
	"regexp"
	"time"
)

// Schema represents a metadata schema definition
type Schema struct {
	Version          string                 `yaml:"schema_version" json:"schema_version"`
	Name             string                 `yaml:"name" json:"name"`
	Description      string                 `yaml:"description" json:"description"`
	Domain           string                 `yaml:"domain" json:"domain"`
	OntologyBase     string                 `yaml:"ontology_base,omitempty" json:"ontology_base,omitempty"`
	Extends          []string               `yaml:"extends,omitempty" json:"extends,omitempty"`
	RequiredFields   []string               `yaml:"required_fields" json:"required_fields"`
	Fields           map[string]FieldSchema `yaml:",inline" json:"fields"`
	ValidationRules  []ValidationRule       `yaml:"validation,omitempty" json:"validation,omitempty"`
	OntologyMappings map[string]string      `yaml:"ontology_mappings,omitempty" json:"ontology_mappings,omitempty"`
	FileFormats      FileFormats            `yaml:"file_formats,omitempty" json:"file_formats,omitempty"`
	Facets           []FacetConfig          `yaml:"facets,omitempty" json:"facets,omitempty"`
}

// FieldSchema defines the schema for a single field
type FieldSchema struct {
	Type        string                 `yaml:"type" json:"type"` // string, number, integer, boolean, array, object, date, datetime
	Required    bool                   `yaml:"required,omitempty" json:"required,omitempty"`
	RequiredIf  string                 `yaml:"required_if,omitempty" json:"required_if,omitempty"`
	Default     interface{}            `yaml:"default,omitempty" json:"default,omitempty"`
	Description string                 `yaml:"description,omitempty" json:"description,omitempty"`
	Examples    []interface{}          `yaml:"examples,omitempty" json:"examples,omitempty"`
	Vocabulary  []string               `yaml:"vocabulary,omitempty" json:"vocabulary,omitempty"`
	Pattern     string                 `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	Units       string                 `yaml:"units,omitempty" json:"units,omitempty"`
	Range       *Range                 `yaml:"range,omitempty" json:"range,omitempty"`
	Fields      map[string]FieldSchema `yaml:"fields,omitempty" json:"fields,omitempty"` // For object types
	Items       *FieldSchema           `yaml:"items,omitempty" json:"items,omitempty"`   // For array types
	MinItems    int                    `yaml:"min_items,omitempty" json:"min_items,omitempty"`
	MaxItems    int                    `yaml:"max_items,omitempty" json:"max_items,omitempty"`
	Auto        string                 `yaml:"auto,omitempty" json:"auto,omitempty"` // Auto-generated values
	Ontology    string                 `yaml:"ontology,omitempty" json:"ontology,omitempty"`
}

// Range defines min/max constraints
type Range struct {
	Min interface{} `yaml:"min,omitempty" json:"min,omitempty"`
	Max interface{} `yaml:"max,omitempty" json:"max,omitempty"`
}

// ValidationRule defines a custom validation rule
type ValidationRule struct {
	Rule    string `yaml:"rule" json:"rule"`
	Message string `yaml:"message" json:"message"`
}

// FileFormats defines expected file formats
type FileFormats struct {
	Primary   []string `yaml:"primary,omitempty" json:"primary,omitempty"`
	Processed []string `yaml:"processed,omitempty" json:"processed,omitempty"`
	Raw       []string `yaml:"raw,omitempty" json:"raw,omitempty"`
	Metadata  []string `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

// FacetConfig defines a facet for search/filtering
type FacetConfig struct {
	Field string `yaml:"field" json:"field"`
	Label string `yaml:"label" json:"label"`
}

// Metadata represents actual metadata for a file/dataset
type Metadata struct {
	SchemaName    string                 `json:"schema_name"`
	SchemaVersion string                 `json:"schema_version"`
	Fields        map[string]interface{} `json:"fields"`
	FileInfo      FileInfo               `json:"file_info"`
	Provenance    Provenance             `json:"provenance"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// FileInfo contains file-specific metadata
type FileInfo struct {
	Filename  string    `json:"filename"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	Checksum  string    `json:"checksum"`
	Format    string    `json:"format"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
}

// Provenance tracks data lineage
type Provenance struct {
	UploadedBy       string         `json:"uploaded_by"`
	UploadedFrom     string         `json:"uploaded_from"`
	UploadedAt       time.Time      `json:"uploaded_at"`
	SourceInstrument string         `json:"source_instrument,omitempty"`
	Workflow         []WorkflowStep `json:"workflow,omitempty"`
	RelatedFiles     []string       `json:"related_files,omitempty"`
}

// WorkflowStep tracks processing steps
type WorkflowStep struct {
	Name       string                 `json:"name"`
	Tool       string                 `json:"tool"`
	Version    string                 `json:"version"`
	Timestamp  time.Time              `json:"timestamp"`
	Inputs     []string               `json:"inputs"`
	Outputs    []string               `json:"outputs"`
	Parameters map[string]interface{} `json:"parameters"`
}

// ValidationResult contains validation results
type ValidationResult struct {
	Valid    bool                `json:"valid"`
	Errors   []ValidationError   `json:"errors,omitempty"`
	Warnings []ValidationWarning `json:"warnings,omitempty"`
	Score    *QualityScore       `json:"score,omitempty"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Type    string `json:"type"` // missing, invalid_type, invalid_value, etc.
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	Field       string   `json:"field"`
	Message     string   `json:"message"`
	Suggestions []string `json:"suggestions,omitempty"`
}

// QualityScore represents metadata quality metrics
type QualityScore struct {
	Overall          int      `json:"overall"` // 0-100
	Completeness     int      `json:"completeness"`
	Consistency      int      `json:"consistency"`
	Richness         int      `json:"richness"`
	Interoperability int      `json:"interoperability"`
	Details          []string `json:"details,omitempty"`
}

// SchemaManager manages metadata schemas
type SchemaManager struct {
	schemas map[string]*Schema
	loader  SchemaLoader
}

// SchemaLoader interface for loading schemas
type SchemaLoader interface {
	Load(name string) (*Schema, error)
	LoadFromFile(path string) (*Schema, error)
	List() ([]SchemaInfo, error)
	Search(query string) ([]SchemaInfo, error)
}

// SchemaInfo contains schema metadata
type SchemaInfo struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Domain      string    `json:"domain"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	CreatedAt   time.Time `json:"created_at"`
	Downloads   int       `json:"downloads"`
	Rating      float64   `json:"rating"`
}

// NewSchemaManager creates a new schema manager
func NewSchemaManager(loader SchemaLoader) *SchemaManager {
	return &SchemaManager{
		schemas: make(map[string]*Schema),
		loader:  loader,
	}
}

// LoadSchema loads a schema by name
func (sm *SchemaManager) LoadSchema(name string) (*Schema, error) {
	// Check cache first
	if schema, ok := sm.schemas[name]; ok {
		return schema, nil
	}

	// Load from loader
	schema, err := sm.loader.Load(name)
	if err != nil {
		return nil, fmt.Errorf("failed to load schema %s: %w", name, err)
	}

	// Resolve extends
	if len(schema.Extends) > 0 {
		if err := sm.resolveExtends(schema); err != nil {
			return nil, fmt.Errorf("failed to resolve extends: %w", err)
		}
	}

	// Cache it
	sm.schemas[name] = schema
	return schema, nil
}

// resolveExtends merges parent schemas
func (sm *SchemaManager) resolveExtends(schema *Schema) error {
	for _, parent := range schema.Extends {
		parentSchema, err := sm.LoadSchema(parent)
		if err != nil {
			return err
		}

		// Merge fields
		for name, field := range parentSchema.Fields {
			if _, exists := schema.Fields[name]; !exists {
				schema.Fields[name] = field
			}
		}

		// Merge required fields
		schema.RequiredFields = append(schema.RequiredFields, parentSchema.RequiredFields...)
	}

	return nil
}

// ValidateMetadata validates metadata against a schema
func (sm *SchemaManager) ValidateMetadata(metadata *Metadata) ValidationResult {
	result := ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}

	// Load schema
	schema, err := sm.LoadSchema(metadata.SchemaName)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "schema",
			Message: fmt.Sprintf("Schema not found: %s", metadata.SchemaName),
			Type:    "schema_not_found",
		})
		return result
	}

	// Validate required fields
	for _, required := range schema.RequiredFields {
		if _, ok := metadata.Fields[required]; !ok {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:   required,
				Message: fmt.Sprintf("Required field '%s' is missing", required),
				Type:    "missing_field",
			})
		}
	}

	// Validate field types and constraints
	for fieldName, value := range metadata.Fields {
		fieldSchema, ok := schema.Fields[fieldName]
		if !ok {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Field:   fieldName,
				Message: fmt.Sprintf("Field '%s' not defined in schema", fieldName),
			})
			continue
		}

		// Type checking
		if !isValidType(value, fieldSchema.Type) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("Invalid type for '%s': expected %s", fieldName, fieldSchema.Type),
				Type:    "invalid_type",
			})
		}

		// Vocabulary checking
		if len(fieldSchema.Vocabulary) > 0 {
			strValue, ok := value.(string)
			if ok && !contains(fieldSchema.Vocabulary, strValue) {
				result.Warnings = append(result.Warnings, ValidationWarning{
					Field:       fieldName,
					Message:     fmt.Sprintf("Value '%s' not in controlled vocabulary", strValue),
					Suggestions: findSimilar(strValue, fieldSchema.Vocabulary),
				})
			}
		}

		// Pattern matching
		if fieldSchema.Pattern != "" {
			strValue, ok := value.(string)
			if ok {
				matched, _ := regexp.MatchString(fieldSchema.Pattern, strValue)
				if !matched {
					result.Valid = false
					result.Errors = append(result.Errors, ValidationError{
						Field:   fieldName,
						Message: fmt.Sprintf("Value does not match pattern: %s", fieldSchema.Pattern),
						Type:    "pattern_mismatch",
					})
				}
			}
		}

		// Range checking
		if fieldSchema.Range != nil {
			if !checkRange(value, fieldSchema.Range) {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:   fieldName,
					Message: fmt.Sprintf("Value out of range: %v", fieldSchema.Range),
					Type:    "out_of_range",
				})
			}
		}
	}

	// Custom validation rules
	for range schema.ValidationRules {
		// TODO: Implement rule evaluation
		// This would involve parsing and evaluating the rule expression
	}

	// Calculate quality score if valid
	if result.Valid {
		result.Score = calculateQualityScore(schema, metadata)
	}

	return result
}

// Helper functions

func isValidType(value interface{}, expectedType string) bool {
	switch expectedType {
	case "string":
		_, ok := value.(string)
		return ok
	case "number":
		_, ok := value.(float64)
		if ok {
			return true
		}
		_, ok = value.(int)
		return ok
	case "integer":
		_, ok := value.(int)
		if ok {
			return true
		}
		f, ok := value.(float64)
		return ok && f == float64(int(f))
	case "boolean":
		_, ok := value.(bool)
		return ok
	case "array":
		_, ok := value.([]interface{})
		return ok
	case "object":
		_, ok := value.(map[string]interface{})
		return ok
	case "date", "datetime":
		_, ok := value.(time.Time)
		if ok {
			return true
		}
		_, ok = value.(string)
		return ok // TODO: Validate date format
	default:
		return true
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func findSimilar(value string, vocabulary []string) []string {
	// TODO: Implement fuzzy matching (Levenshtein distance)
	return []string{}
}

func checkRange(value interface{}, r *Range) bool {
	// TODO: Implement range checking for different types
	return true
}

func calculateQualityScore(schema *Schema, metadata *Metadata) *QualityScore {
	score := &QualityScore{}

	// Completeness: percentage of fields filled
	totalFields := len(schema.Fields)
	filledFields := len(metadata.Fields)
	score.Completeness = (filledFields * 100) / totalFields

	// Consistency: use of controlled vocabularies
	// TODO: Implement
	score.Consistency = 85

	// Richness: optional fields filled, notes, keywords
	// TODO: Implement
	score.Richness = 70

	// Interoperability: ontology terms, standard formats
	// TODO: Implement
	score.Interoperability = 80

	// Overall: weighted average
	score.Overall = (score.Completeness + score.Consistency + score.Richness + score.Interoperability) / 4

	return score
}
