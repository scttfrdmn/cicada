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
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// CZI file format constants
const (
	cziMagicBytes    = "ZISRAWFILE"
	cziSegmentHeader = 16 // Size of segment header
)

// CZI segment types
const (
	segmentMetadata = 0x00000001 // XML metadata segment
	segmentSubBlock = 0x00000002 // Image data subblock
	segmentAttach   = 0x00000003 // Attachment segment
)

// ZeissCZIExtractor extracts metadata from Zeiss CZI microscopy files.
//
// CZI files are complex binary containers with:
// - ZIP-based compression
// - Segmented structure (header, metadata, subblocks, attachments)
// - XML metadata stored as UTF-8
//
// This implementation focuses on extracting the XML metadata segment
// and parsing key fields for instrument awareness and FAIR compliance.
type ZeissCZIExtractor struct{}

// cziMetadataXML represents the root structure of CZI XML metadata.
// CZI metadata follows a hierarchical XML schema with instrument,
// image, and acquisition information.
type cziMetadataXML struct {
	XMLName      xml.Name             `xml:"ImageDocument"`
	Metadata     cziMetadata          `xml:"Metadata"`
	Information  cziInformation       `xml:"Information"`
}

type cziMetadata struct {
	Scaling      cziScaling           `xml:"Scaling"`
	Information  cziInstrumentInfo    `xml:"Information"`
}

type cziScaling struct {
	Items        []cziScalingItem     `xml:"Items>Distance"`
}

type cziScalingItem struct {
	ID           string               `xml:"Id,attr"`
	Value        float64              `xml:"Value"`
}

type cziInstrumentInfo struct {
	Instrument   cziInstrument        `xml:"Instrument"`
	Image        cziImage             `xml:"Image"`
	User         cziUser              `xml:"User"`
}

type cziInstrument struct {
	Microscopes  []cziMicroscope      `xml:"Microscopes>Microscope"`
	Objectives   []cziObjective       `xml:"Objectives>Objective"`
}

type cziMicroscope struct {
	Name         string               `xml:"Name,attr"`
	System       string               `xml:"System"`
}

type cziImage struct {
	SizeX        int                  `xml:"SizeX"`
	SizeY        int                  `xml:"SizeY"`
	SizeZ        int                  `xml:"SizeZ"`
	SizeC        int                  `xml:"SizeC"` // Number of channels
	SizeT        int                  `xml:"SizeT"` // Number of time points
	ComponentBitCount int             `xml:"ComponentBitCount"`
}

type cziUser struct {
	Name         string               `xml:"Name"`
}

type cziObjective struct {
	Name         string               `xml:"Name,attr"`
	Id           string               `xml:"Id,attr"`
	Magnification float64             `xml:"NominalMagnification"`
	NA           float64              `xml:"LensNA"` // Numerical Aperture
	Immersion    string               `xml:"Immersion"`
	WorkingDistance float64           `xml:"WorkingDistance"`
}

type cziInformation struct {
	Application  cziApplication       `xml:"Application"`
	Image        cziImageInfo         `xml:"Image"`
	Document     cziDocument          `xml:"Document"`
}

type cziApplication struct {
	Name         string               `xml:"Name"`
	Version      string               `xml:"Version"`
}

type cziImageInfo struct {
	AcquisitionDateAndTime string     `xml:"AcquisitionDateAndTime"`
	Dimensions   cziDimensions        `xml:"Dimensions"`
}

type cziDimensions struct {
	Channels     []cziChannel         `xml:"Channels>Channel"`
}

type cziChannel struct {
	Id           string               `xml:"Id,attr"`
	Name         string               `xml:"Name"`
	EmissionWavelength float64       `xml:"EmissionWavelength"`
	ExcitationWavelength float64     `xml:"ExcitationWavelength"`
	DyeName      string               `xml:"DyeName"`
}

type cziDocument struct {
	Name         string               `xml:"Name"`
	CreationDate string               `xml:"CreationDate"`
}

// Name returns the extractor name.
func (e *ZeissCZIExtractor) Name() string {
	return "Zeiss CZI"
}

// SupportedFormats returns the file extensions this extractor handles.
func (e *ZeissCZIExtractor) SupportedFormats() []string {
	return []string{".czi"}
}

// CanHandle returns true if this extractor can handle the given filename.
func (e *ZeissCZIExtractor) CanHandle(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".czi")
}

// Extract extracts metadata from a CZI file.
func (e *ZeissCZIExtractor) Extract(filepath string) (map[string]interface{}, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	// Get file info
	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	return e.extractFromFile(f, info.Size(), filepath)
}

// ExtractFromReader extracts metadata from a reader.
func (e *ZeissCZIExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
	// Read entire file into memory for ZIP processing
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return e.extractFromBytes(data, filename)
}

// extractFromFile extracts metadata from an open file.
func (e *ZeissCZIExtractor) extractFromFile(f *os.File, size int64, filepath string) (map[string]interface{}, error) {
	// Verify CZI magic bytes
	header := make([]byte, 16)
	if _, err := io.ReadFull(f, header); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	if string(header[0:10]) != cziMagicBytes {
		return nil, fmt.Errorf("not a valid CZI file: invalid magic bytes")
	}

	// Reset to beginning for full read
	if _, err := f.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek: %w", err)
	}

	// Read full file for processing
	data := make([]byte, size)
	if _, err := io.ReadFull(f, data); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return e.extractFromBytes(data, filepath)
}

// extractFromBytes extracts metadata from raw file bytes.
func (e *ZeissCZIExtractor) extractFromBytes(data []byte, filename string) (map[string]interface{}, error) {
	// Verify magic bytes
	if len(data) < 16 || string(data[0:10]) != cziMagicBytes {
		return nil, fmt.Errorf("not a valid CZI file")
	}

	// Initialize metadata
	metadata := map[string]interface{}{
		"format":       "CZI",
		"manufacturer": "Zeiss",
		"file_name":    filename,
		"file_size":    len(data),
	}

	// Parse CZI segments to find XML metadata
	xmlData, err := e.extractXMLMetadata(data)
	if err != nil {
		return nil, fmt.Errorf("failed to extract XML metadata: %w", err)
	}

	if xmlData != nil {
		// Parse XML and extract fields
		if err := e.parseXMLMetadata(xmlData, metadata); err != nil {
			return nil, fmt.Errorf("failed to parse XML metadata: %w", err)
		}
	} else {
		// No XML found, return basic metadata
		metadata["extraction_note"] = "No XML metadata segment found in CZI file"
	}

	return metadata, nil
}

// extractXMLMetadata finds and extracts the XML metadata segment from CZI file.
// CZI files have a segmented structure with a 16-byte header followed by segments.
func (e *ZeissCZIExtractor) extractXMLMetadata(data []byte) ([]byte, error) {
	offset := 16 // Skip file header

	// Scan through segments looking for metadata segment
	for offset < len(data)-16 {
		// Read segment header (16 bytes)
		if offset+16 > len(data) {
			break
		}

		// Segment structure:
		// 0-15: Segment ID (16 bytes, "ZISRAWMETADATA\0\0" for metadata)
		// 16-23: Allocated size (int64)
		// 24-31: Used size (int64)

		segmentID := string(bytes.TrimRight(data[offset:offset+16], "\x00"))

		if offset+32 > len(data) {
			break
		}

		allocatedSize := int64(binary.LittleEndian.Uint64(data[offset+16 : offset+24]))
		usedSize := int64(binary.LittleEndian.Uint64(data[offset+24 : offset+32]))

		// Check if this is the metadata segment
		if strings.HasPrefix(segmentID, "ZISRAWMETADATA") {
			// XML data starts after 32-byte segment header
			xmlStart := offset + 32
			xmlEnd := xmlStart + int(usedSize)

			if xmlEnd <= len(data) {
				return data[xmlStart:xmlEnd], nil
			}
		}

		// Move to next segment
		offset += 32 + int(allocatedSize)
	}

	return nil, nil // No metadata segment found
}

// parseXMLMetadata parses the CZI XML metadata and extracts key fields.
func (e *ZeissCZIExtractor) parseXMLMetadata(xmlData []byte, metadata map[string]interface{}) error {
	var doc cziMetadataXML
	if err := xml.Unmarshal(xmlData, &doc); err != nil {
		return fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	// Extract instrument information
	if len(doc.Metadata.Information.Instrument.Microscopes) > 0 {
		microscope := doc.Metadata.Information.Instrument.Microscopes[0]
		if microscope.System != "" {
			metadata["instrument_model"] = microscope.System
			metadata["instrument_type"] = "microscopy"
		}
		if microscope.Name != "" {
			metadata["microscope_name"] = microscope.Name
		}
	}

	// Extract objective information
	if len(doc.Metadata.Information.Instrument.Objectives) > 0 {
		objective := doc.Metadata.Information.Instrument.Objectives[0]
		if objective.Magnification > 0 {
			metadata["objective_magnification"] = objective.Magnification
		}
		if objective.NA > 0 {
			metadata["objective_na"] = objective.NA
		}
		if objective.Immersion != "" {
			metadata["objective_immersion"] = objective.Immersion
		}
		if objective.Name != "" {
			metadata["objective_name"] = objective.Name
		}
	}

	// Extract image dimensions
	if doc.Metadata.Information.Image.SizeX > 0 {
		metadata["image_width"] = doc.Metadata.Information.Image.SizeX
	}
	if doc.Metadata.Information.Image.SizeY > 0 {
		metadata["image_height"] = doc.Metadata.Information.Image.SizeY
	}
	if doc.Metadata.Information.Image.SizeZ > 0 {
		metadata["image_depth"] = doc.Metadata.Information.Image.SizeZ
	}
	if doc.Metadata.Information.Image.SizeC > 0 {
		metadata["num_channels"] = doc.Metadata.Information.Image.SizeC
	}
	if doc.Metadata.Information.Image.SizeT > 0 {
		metadata["num_timepoints"] = doc.Metadata.Information.Image.SizeT
	}
	if doc.Metadata.Information.Image.ComponentBitCount > 0 {
		metadata["bit_depth"] = doc.Metadata.Information.Image.ComponentBitCount
	}

	// Extract pixel scaling
	for _, item := range doc.Metadata.Scaling.Items {
		switch item.ID {
		case "X":
			metadata["pixel_size_x_um"] = item.Value * 1e6 // Convert to micrometers
		case "Y":
			metadata["pixel_size_y_um"] = item.Value * 1e6
		case "Z":
			metadata["pixel_size_z_um"] = item.Value * 1e6
		}
	}

	// Extract acquisition date
	if doc.Information.Image.AcquisitionDateAndTime != "" {
		// Try to parse as RFC3339
		if t, err := time.Parse(time.RFC3339, doc.Information.Image.AcquisitionDateAndTime); err == nil {
			metadata["acquisition_date"] = t.Format(time.RFC3339)
		} else {
			metadata["acquisition_date"] = doc.Information.Image.AcquisitionDateAndTime
		}
	} else if doc.Information.Document.CreationDate != "" {
		metadata["acquisition_date"] = doc.Information.Document.CreationDate
	}

	// Extract user/operator
	if doc.Metadata.Information.User.Name != "" {
		metadata["operator"] = doc.Metadata.Information.User.Name
	}

	// Extract application info
	if doc.Information.Application.Name != "" {
		metadata["software_name"] = doc.Information.Application.Name
	}
	if doc.Information.Application.Version != "" {
		metadata["software_version"] = doc.Information.Application.Version
	}

	// Extract channel information
	if len(doc.Information.Image.Dimensions.Channels) > 0 {
		channels := make([]map[string]interface{}, 0)
		for _, ch := range doc.Information.Image.Dimensions.Channels {
			channelInfo := map[string]interface{}{
				"id":   ch.Id,
				"name": ch.Name,
			}
			if ch.EmissionWavelength > 0 {
				channelInfo["emission_wavelength_nm"] = ch.EmissionWavelength
			}
			if ch.ExcitationWavelength > 0 {
				channelInfo["excitation_wavelength_nm"] = ch.ExcitationWavelength
			}
			if ch.DyeName != "" {
				channelInfo["dye_name"] = ch.DyeName
			}
			channels = append(channels, channelInfo)
		}
		metadata["channels"] = channels
	}

	// Extract document name
	if doc.Information.Document.Name != "" {
		metadata["document_name"] = doc.Information.Document.Name
	}

	// Add extractor info
	metadata["extractor_name"] = "zeiss_czi"
	metadata["schema_name"] = "zeiss_czi_v1"

	return nil
}
