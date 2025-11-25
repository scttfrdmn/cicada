package metadata

import (
	"bytes"
	"testing"
)

func TestZeissExtractor_CanHandle(t *testing.T) {
	extractor := &ZeissExtractor{}

	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{"CZI file lowercase", "sample.czi", true},
		{"CZI file uppercase", "SAMPLE.CZI", true},
		{"CZI file mixed case", "Sample.CzI", true},
		{"TIFF file", "sample.tif", false},
		{"No extension", "sample", false},
		{"Wrong extension", "sample.nd2", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractor.CanHandle(tt.filename); got != tt.want {
				t.Errorf("CanHandle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZeissExtractor_ExtractFromReader(t *testing.T) {
	extractor := &ZeissExtractor{}

	t.Run("Valid CZI header", func(t *testing.T) {
		// Create mock CZI file with valid header
		data := []byte("ZISRAWFILE\x00\x00\x00\x00\x00\x00")
		reader := bytes.NewReader(data)

		metadata, err := extractor.ExtractFromReader(reader, "test.czi")
		if err != nil {
			t.Fatalf("ExtractFromReader() error = %v", err)
		}

		// Verify basic metadata fields
		if metadata["format"] != "CZI" {
			t.Errorf("Expected format=CZI, got %v", metadata["format"])
		}

		if metadata["manufacturer"] != "Zeiss" {
			t.Errorf("Expected manufacturer=Zeiss, got %v", metadata["manufacturer"])
		}

		if metadata["file_name"] != "test.czi" {
			t.Errorf("Expected file_name=test.czi, got %v", metadata["file_name"])
		}
	})

	t.Run("Invalid CZI header", func(t *testing.T) {
		// Invalid magic bytes
		data := []byte("NOT_A_CZI_FILE\x00\x00")
		reader := bytes.NewReader(data)

		_, err := extractor.ExtractFromReader(reader, "invalid.czi")
		if err == nil {
			t.Error("Expected error for invalid CZI file, got nil")
		}
	})

	t.Run("Too short file", func(t *testing.T) {
		// File too short to contain header
		data := []byte("SHORT")
		reader := bytes.NewReader(data)

		_, err := extractor.ExtractFromReader(reader, "short.czi")
		if err == nil {
			t.Error("Expected error for too-short file, got nil")
		}
	})
}

func TestZeissExtractor_Name(t *testing.T) {
	extractor := &ZeissExtractor{}
	if got := extractor.Name(); got != "Zeiss CZI" {
		t.Errorf("Name() = %v, want 'Zeiss CZI'", got)
	}
}

func TestZeissExtractor_SupportedFormats(t *testing.T) {
	extractor := &ZeissExtractor{}
	formats := extractor.SupportedFormats()

	if len(formats) != 1 {
		t.Errorf("Expected 1 format, got %d", len(formats))
	}

	if formats[0] != ".czi" {
		t.Errorf("Expected .czi format, got %v", formats[0])
	}
}

func TestExtractorRegistry_FindExtractor(t *testing.T) {
	registry := NewExtractorRegistry()
	registry.Register(&ZeissExtractor{})
	registry.Register(&TIFFExtractor{})

	tests := []struct {
		name         string
		filename     string
		wantExtractor string
	}{
		{"CZI file", "sample.czi", "Zeiss CZI"},
		{"TIFF file", "sample.tif", "TIFF"},
		{"Unknown file", "sample.xyz", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor := registry.FindExtractor(tt.filename)
			if tt.wantExtractor == "" {
				if extractor != nil {
					t.Errorf("Expected no extractor, got %v", extractor.Name())
				}
			} else {
				if extractor == nil {
					t.Error("Expected extractor, got nil")
				} else if extractor.Name() != tt.wantExtractor {
					t.Errorf("Expected %v, got %v", tt.wantExtractor, extractor.Name())
				}
			}
		})
	}
}

func TestExtractorRegistry_Register(t *testing.T) {
	registry := NewExtractorRegistry()

	// Initially empty
	if len(registry.extractors) != 0 {
		t.Errorf("New registry should be empty, got %d extractors", len(registry.extractors))
	}

	// Register first extractor
	registry.Register(&TIFFExtractor{})
	if len(registry.extractors) != 1 {
		t.Errorf("Expected 1 extractor after registration, got %d", len(registry.extractors))
	}

	// Register second extractor
	registry.Register(&ZeissExtractor{})
	if len(registry.extractors) != 2 {
		t.Errorf("Expected 2 extractors after second registration, got %d", len(registry.extractors))
	}

	// Verify we can find both
	tiffExt := registry.FindExtractor("test.tif")
	if tiffExt == nil || tiffExt.Name() != "TIFF" {
		t.Error("Failed to find TIFF extractor after registration")
	}

	cziExt := registry.FindExtractor("test.czi")
	if cziExt == nil || cziExt.Name() != "Zeiss CZI" {
		t.Error("Failed to find CZI extractor after registration")
	}
}

func TestExtractorRegistry_ListExtractors(t *testing.T) {
	registry := NewExtractorRegistry()
	registry.Register(&TIFFExtractor{})
	registry.Register(&ZeissExtractor{})

	extractors := registry.ListExtractors()

	if len(extractors) != 2 {
		t.Errorf("Expected 2 extractors, got %d", len(extractors))
	}

	// Verify extractor info contains expected data
	foundTIFF := false
	foundCZI := false
	for _, info := range extractors {
		if info.Name == "TIFF" {
			foundTIFF = true
			if len(info.Formats) == 0 {
				t.Error("TIFF extractor should have formats")
			}
		}
		if info.Name == "Zeiss CZI" {
			foundCZI = true
			if len(info.Formats) == 0 {
				t.Error("CZI extractor should have formats")
			}
		}
	}

	if !foundTIFF {
		t.Error("TIFF extractor not found in list")
	}
	if !foundCZI {
		t.Error("CZI extractor not found in list")
	}
}

func TestExtractorRegistry_RegisterDefaults(t *testing.T) {
	registry := NewExtractorRegistry()
	registry.RegisterDefaults()

	// Should have registered multiple extractors
	if len(registry.extractors) < 5 {
		t.Errorf("Expected at least 5 default extractors, got %d", len(registry.extractors))
	}

	// Verify specific extractors are registered
	tests := []struct {
		filename string
		wantName string
	}{
		{"sample.tif", "TIFF"},
		{"sample.czi", "Zeiss CZI"},
		{"sample.fastq", "FASTQ"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			extractor := registry.FindExtractor(tt.filename)
			if extractor == nil {
				t.Errorf("Expected to find extractor for %s", tt.filename)
			} else if extractor.Name() != tt.wantName {
				t.Errorf("Expected %s extractor, got %s", tt.wantName, extractor.Name())
			}
		})
	}

	// Verify GenericExtractor is last (fallback)
	genericExt := registry.FindExtractor("unknown.xyz")
	if genericExt == nil {
		t.Error("GenericExtractor should be registered as fallback")
	} else if genericExt.Name() != "Generic" {
		t.Errorf("Expected Generic extractor for unknown file, got %s", genericExt.Name())
	}
}

func TestGenericExtractor_CanHandle(t *testing.T) {
	extractor := &GenericExtractor{}

	// GenericExtractor should handle any file
	tests := []string{
		"file.xyz",
		"file.unknown",
		"noextension",
		"file.txt",
		"file.jpg",
	}

	for _, filename := range tests {
		t.Run(filename, func(t *testing.T) {
			if !extractor.CanHandle(filename) {
				t.Errorf("GenericExtractor should handle %s", filename)
			}
		})
	}
}

func TestGenericExtractor_ExtractFromReader(t *testing.T) {
	extractor := &GenericExtractor{}

	data := []byte("test file content")
	reader := bytes.NewReader(data)

	metadata, err := extractor.ExtractFromReader(reader, "test.xyz")
	if err != nil {
		t.Fatalf("ExtractFromReader() error = %v", err)
	}

	// Should extract basic format
	if metadata["format"] == nil {
		t.Error("Expected format field in metadata")
	}

	// Format should be detected from filename
	if format, ok := metadata["format"].(string); ok {
		if format != "XYZ" && format != "Unknown" {
			t.Errorf("Unexpected format: %s", format)
		}
	}
}

func TestGenericExtractor_Name(t *testing.T) {
	extractor := &GenericExtractor{}
	if got := extractor.Name(); got != "Generic" {
		t.Errorf("Name() = %v, want 'Generic'", got)
	}
}

func TestGenericExtractor_SupportedFormats(t *testing.T) {
	extractor := &GenericExtractor{}
	formats := extractor.SupportedFormats()

	if len(formats) != 1 {
		t.Errorf("Expected 1 format, got %d", len(formats))
	}

	if formats[0] != "*" {
		t.Errorf("Expected * format (wildcard), got %v", formats[0])
	}
}

func TestNewExtractorRegistry(t *testing.T) {
	registry := NewExtractorRegistry()

	if registry == nil {
		t.Fatal("NewExtractorRegistry() returned nil")
	}

	if registry.extractors == nil {
		t.Error("Registry extractors slice should be initialized")
	}

	if len(registry.extractors) != 0 {
		t.Errorf("New registry should be empty, got %d extractors", len(registry.extractors))
	}
}
