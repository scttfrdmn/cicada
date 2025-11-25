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

package metadata

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestZeissCZIExtractor_Name(t *testing.T) {
	extractor := &ZeissCZIExtractor{}
	if extractor.Name() != "Zeiss CZI" {
		t.Errorf("Name() = %v, want %v", extractor.Name(), "Zeiss CZI")
	}
}

func TestZeissCZIExtractor_SupportedFormats(t *testing.T) {
	extractor := &ZeissCZIExtractor{}
	formats := extractor.SupportedFormats()
	if len(formats) != 1 || formats[0] != ".czi" {
		t.Errorf("SupportedFormats() = %v, want [.czi]", formats)
	}
}

func TestZeissCZIExtractor_CanHandle(t *testing.T) {
	extractor := &ZeissCZIExtractor{}
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{"CZI file lowercase", "test.czi", true},
		{"CZI file uppercase", "test.CZI", true},
		{"CZI file mixed case", "test.Czi", true},
		{"TIFF file", "test.tiff", false},
		{"No extension", "test", false},
		{"Wrong extension", "test.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractor.CanHandle(tt.filename); got != tt.want {
				t.Errorf("CanHandle(%v) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

func TestZeissCZIExtractor_ExtractFromReader_InvalidMagicBytes(t *testing.T) {
	extractor := &ZeissCZIExtractor{}

	// Create invalid CZI file (wrong magic bytes)
	data := []byte("NOTAVALIDCZI123456789")

	_, err := extractor.ExtractFromReader(bytes.NewReader(data), "test.czi")
	if err == nil {
		t.Error("ExtractFromReader() should return error for invalid magic bytes")
	}
	if !strings.Contains(err.Error(), "not a valid CZI file") {
		t.Errorf("ExtractFromReader() error = %v, want error containing 'not a valid CZI file'", err)
	}
}

func TestZeissCZIExtractor_ExtractFromReader_TooShort(t *testing.T) {
	extractor := &ZeissCZIExtractor{}

	// Create file that's too short
	data := []byte("ZISRAW")

	_, err := extractor.ExtractFromReader(bytes.NewReader(data), "test.czi")
	if err == nil {
		t.Error("ExtractFromReader() should return error for too short file")
	}
}

func TestZeissCZIExtractor_ExtractFromReader_ValidHeaderNoMetadata(t *testing.T) {
	extractor := &ZeissCZIExtractor{}

	// Create minimal valid CZI file with just header (no metadata segment)
	data := createMinimalCZIFile(nil)

	metadata, err := extractor.ExtractFromReader(bytes.NewReader(data), "test.czi")
	if err != nil {
		t.Fatalf("ExtractFromReader() error = %v", err)
	}

	// Check basic fields
	if metadata["format"] != "CZI" {
		t.Errorf("format = %v, want CZI", metadata["format"])
	}
	if metadata["manufacturer"] != "Zeiss" {
		t.Errorf("manufacturer = %v, want Zeiss", metadata["manufacturer"])
	}
	if metadata["file_name"] != "test.czi" {
		t.Errorf("file_name = %v, want test.czi", metadata["file_name"])
	}
}

func TestZeissCZIExtractor_ExtractFromReader_WithXMLMetadata(t *testing.T) {
	extractor := &ZeissCZIExtractor{}

	// Create CZI file with XML metadata
	xmlData := []byte(`<?xml version="1.0" encoding="utf-8"?>
<ImageDocument>
  <Metadata>
    <Information>
      <Instrument>
        <Microscopes>
          <Microscope Name="LSM980">
            <System>LSM 980</System>
          </Microscope>
        </Microscopes>
        <Objectives>
          <Objective Name="Plan-Apochromat" Id="1">
            <NominalMagnification>40</NominalMagnification>
            <LensNA>1.3</LensNA>
            <Immersion>Oil</Immersion>
          </Objective>
        </Objectives>
      </Instrument>
      <Image>
        <SizeX>1024</SizeX>
        <SizeY>1024</SizeY>
        <SizeZ>10</SizeZ>
        <SizeC>3</SizeC>
        <SizeT>1</SizeT>
        <ComponentBitCount>16</ComponentBitCount>
      </Image>
      <User>
        <Name>John Doe</Name>
      </User>
    </Information>
    <Scaling>
      <Items>
        <Distance Id="X">
          <Value>0.0001</Value>
        </Distance>
        <Distance Id="Y">
          <Value>0.0001</Value>
        </Distance>
        <Distance Id="Z">
          <Value>0.0003</Value>
        </Distance>
      </Items>
    </Scaling>
  </Metadata>
  <Information>
    <Application>
      <Name>ZEN</Name>
      <Version>3.5</Version>
    </Application>
    <Image>
      <AcquisitionDateAndTime>2025-11-23T10:30:00Z</AcquisitionDateAndTime>
      <Dimensions>
        <Channels>
          <Channel Id="Ch1">
            <Name>DAPI</Name>
            <EmissionWavelength>461</EmissionWavelength>
            <ExcitationWavelength>405</ExcitationWavelength>
            <DyeName>DAPI</DyeName>
          </Channel>
          <Channel Id="Ch2">
            <Name>GFP</Name>
            <EmissionWavelength>509</EmissionWavelength>
            <ExcitationWavelength>488</ExcitationWavelength>
            <DyeName>GFP</DyeName>
          </Channel>
        </Channels>
      </Dimensions>
    </Image>
    <Document>
      <Name>test_image</Name>
      <CreationDate>2025-11-23</CreationDate>
    </Document>
  </Information>
</ImageDocument>`)

	data := createMinimalCZIFile(xmlData)

	metadata, err := extractor.ExtractFromReader(bytes.NewReader(data), "test.czi")
	if err != nil {
		t.Fatalf("ExtractFromReader() error = %v", err)
	}

	// Check basic fields
	if metadata["format"] != "CZI" {
		t.Errorf("format = %v, want CZI", metadata["format"])
	}
	if metadata["manufacturer"] != "Zeiss" {
		t.Errorf("manufacturer = %v, want Zeiss", metadata["manufacturer"])
	}

	// Check instrument fields
	if metadata["instrument_model"] != "LSM 980" {
		t.Errorf("instrument_model = %v, want LSM 980", metadata["instrument_model"])
	}
	if metadata["instrument_type"] != "microscopy" {
		t.Errorf("instrument_type = %v, want microscopy", metadata["instrument_type"])
	}

	// Check objective fields
	if metadata["objective_magnification"] != float64(40) {
		t.Errorf("objective_magnification = %v, want 40", metadata["objective_magnification"])
	}
	if metadata["objective_na"] != 1.3 {
		t.Errorf("objective_na = %v, want 1.3", metadata["objective_na"])
	}
	if metadata["objective_immersion"] != "Oil" {
		t.Errorf("objective_immersion = %v, want Oil", metadata["objective_immersion"])
	}

	// Check image dimensions
	if metadata["image_width"] != 1024 {
		t.Errorf("image_width = %v, want 1024", metadata["image_width"])
	}
	if metadata["image_height"] != 1024 {
		t.Errorf("image_height = %v, want 1024", metadata["image_height"])
	}
	if metadata["image_depth"] != 10 {
		t.Errorf("image_depth = %v, want 10", metadata["image_depth"])
	}
	if metadata["num_channels"] != 3 {
		t.Errorf("num_channels = %v, want 3", metadata["num_channels"])
	}
	if metadata["bit_depth"] != 16 {
		t.Errorf("bit_depth = %v, want 16", metadata["bit_depth"])
	}

	// Check pixel size (converted to micrometers)
	if metadata["pixel_size_x_um"] != 100.0 {
		t.Errorf("pixel_size_x_um = %v, want 100.0", metadata["pixel_size_x_um"])
	}
	if metadata["pixel_size_y_um"] != 100.0 {
		t.Errorf("pixel_size_y_um = %v, want 100.0", metadata["pixel_size_y_um"])
	}
	if metadata["pixel_size_z_um"] != 300.0 {
		t.Errorf("pixel_size_z_um = %v, want 300.0", metadata["pixel_size_z_um"])
	}

	// Check acquisition date
	if metadata["acquisition_date"] != "2025-11-23T10:30:00Z" {
		t.Errorf("acquisition_date = %v, want 2025-11-23T10:30:00Z", metadata["acquisition_date"])
	}

	// Check operator
	if metadata["operator"] != "John Doe" {
		t.Errorf("operator = %v, want John Doe", metadata["operator"])
	}

	// Check software
	if metadata["software_name"] != "ZEN" {
		t.Errorf("software_name = %v, want ZEN", metadata["software_name"])
	}
	if metadata["software_version"] != "3.5" {
		t.Errorf("software_version = %v, want 3.5", metadata["software_version"])
	}

	// Check channels
	channels, ok := metadata["channels"].([]map[string]interface{})
	if !ok {
		t.Fatalf("channels is not []map[string]interface{}")
	}
	if len(channels) != 2 {
		t.Errorf("len(channels) = %v, want 2", len(channels))
	}

	// Check first channel (DAPI)
	if channels[0]["name"] != "DAPI" {
		t.Errorf("channel[0].name = %v, want DAPI", channels[0]["name"])
	}
	if channels[0]["emission_wavelength_nm"] != float64(461) {
		t.Errorf("channel[0].emission_wavelength_nm = %v, want 461", channels[0]["emission_wavelength_nm"])
	}

	// Check extractor metadata
	if metadata["extractor_name"] != "zeiss_czi" {
		t.Errorf("extractor_name = %v, want zeiss_czi", metadata["extractor_name"])
	}
	if metadata["schema_name"] != "zeiss_czi_v1" {
		t.Errorf("schema_name = %v, want zeiss_czi_v1", metadata["schema_name"])
	}
}

func TestZeissCZIExtractor_ExtractFromReader_PartialMetadata(t *testing.T) {
	extractor := &ZeissCZIExtractor{}

	// Create CZI file with minimal XML metadata (some fields missing)
	xmlData := []byte(`<?xml version="1.0" encoding="utf-8"?>
<ImageDocument>
  <Metadata>
    <Information>
      <Instrument>
        <Microscopes>
          <Microscope Name="LSM880">
            <System>LSM 880</System>
          </Microscope>
        </Microscopes>
      </Instrument>
    </Information>
  </Metadata>
  <Information>
    <Application>
      <Name>ZEN</Name>
    </Application>
  </Information>
</ImageDocument>`)

	data := createMinimalCZIFile(xmlData)

	metadata, err := extractor.ExtractFromReader(bytes.NewReader(data), "test.czi")
	if err != nil {
		t.Fatalf("ExtractFromReader() error = %v", err)
	}

	// Should have basic fields
	if metadata["format"] != "CZI" {
		t.Errorf("format = %v, want CZI", metadata["format"])
	}
	if metadata["instrument_model"] != "LSM 880" {
		t.Errorf("instrument_model = %v, want LSM 880", metadata["instrument_model"])
	}
	if metadata["software_name"] != "ZEN" {
		t.Errorf("software_name = %v, want ZEN", metadata["software_name"])
	}

	// Should not have fields that weren't in XML
	if _, exists := metadata["objective_magnification"]; exists {
		t.Error("objective_magnification should not exist when not in XML")
	}
}

func TestZeissCZIExtractor_ExtractXMLMetadata(t *testing.T) {
	extractor := &ZeissCZIExtractor{}

	// Create CZI file with metadata segment
	xmlData := []byte(`<test>data</test>`)
	data := createMinimalCZIFile(xmlData)

	// Extract XML
	extractedXML, err := extractor.extractXMLMetadata(data)
	if err != nil {
		t.Fatalf("extractXMLMetadata() error = %v", err)
	}

	if extractedXML == nil {
		t.Fatal("extractXMLMetadata() returned nil")
	}

	if !bytes.Contains(extractedXML, []byte("<test>data</test>")) {
		t.Errorf("extractXMLMetadata() = %s, want XML containing <test>data</test>", extractedXML)
	}
}

func TestZeissCZIExtractor_Extract_FromFile(t *testing.T) {
	extractor := &ZeissCZIExtractor{}

	// Create test CZI file with XML metadata
	xmlData := []byte(`<?xml version="1.0" encoding="utf-8"?>
<ImageDocument>
  <Metadata>
    <Information>
      <Instrument>
        <Microscopes>
          <Microscope Name="LSM880">
            <System>LSM 880</System>
          </Microscope>
        </Microscopes>
      </Instrument>
      <Image>
        <SizeX>512</SizeX>
        <SizeY>512</SizeY>
      </Image>
    </Information>
  </Metadata>
  <Information>
    <Application>
      <Name>ZEN</Name>
      <Version>3.0</Version>
    </Application>
  </Information>
</ImageDocument>`)

	data := createMinimalCZIFile(xmlData)

	// Write to temporary file
	tmpFile := t.TempDir() + "/test.czi"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	// Extract from file
	metadata, err := extractor.Extract(tmpFile)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	// Verify metadata
	if metadata["format"] != "CZI" {
		t.Errorf("format = %v, want CZI", metadata["format"])
	}
	if metadata["instrument_model"] != "LSM 880" {
		t.Errorf("instrument_model = %v, want LSM 880", metadata["instrument_model"])
	}
	if metadata["image_width"] != 512 {
		t.Errorf("image_width = %v, want 512", metadata["image_width"])
	}
}

func TestZeissCZIExtractor_Extract_NonExistentFile(t *testing.T) {
	extractor := &ZeissCZIExtractor{}

	_, err := extractor.Extract("/nonexistent/file.czi")
	if err == nil {
		t.Error("Extract() should return error for non-existent file")
	}
	if !strings.Contains(err.Error(), "failed to open file") {
		t.Errorf("Extract() error = %v, want error containing 'failed to open file'", err)
	}
}

// createMinimalCZIFile creates a minimal valid CZI file structure for testing.
// CZI file structure:
//   - 16-byte file header with magic bytes "ZISRAWFILE"
//   - Optional metadata segment with XML data
func createMinimalCZIFile(xmlMetadata []byte) []byte {
	buf := new(bytes.Buffer)

	// Write file header (16 bytes)
	// Magic bytes: "ZISRAWFILE" (10 bytes) + padding (6 bytes)
	buf.WriteString("ZISRAWFILE")
	buf.Write(make([]byte, 6)) // Padding

	// If XML metadata provided, add metadata segment
	if xmlMetadata != nil {
		// Segment header (32 bytes):
		// - Segment ID (16 bytes): "ZISRAWMETADATA\0\0"
		// - Allocated size (8 bytes, int64)
		// - Used size (8 bytes, int64)

		segmentID := make([]byte, 16)
		copy(segmentID, "ZISRAWMETADATA")
		buf.Write(segmentID)

		// Allocated size = used size for simplicity
		allocatedSize := int64(len(xmlMetadata))
		if err := binary.Write(buf, binary.LittleEndian, allocatedSize); err != nil {
			panic(fmt.Sprintf("failed to write allocated size: %v", err))
		}

		// Used size
		usedSize := int64(len(xmlMetadata))
		if err := binary.Write(buf, binary.LittleEndian, usedSize); err != nil {
			panic(fmt.Sprintf("failed to write used size: %v", err))
		}

		// XML data
		buf.Write(xmlMetadata)
	}

	return buf.Bytes()
}
