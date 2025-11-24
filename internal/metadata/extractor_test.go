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
