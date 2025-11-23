// internal/metadata/extractor.go
package metadata

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"
)

// Extractor interface for extracting metadata from files
type Extractor interface {
	// CanHandle returns true if this extractor can handle the given file
	CanHandle(filename string) bool

	// Extract extracts metadata from a file
	Extract(filepath string) (map[string]interface{}, error)

	// ExtractFromReader extracts metadata from a reader
	ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error)

	// Name returns the extractor name
	Name() string

	// SupportedFormats returns list of supported file extensions
	SupportedFormats() []string
}

// ExtractorRegistry manages metadata extractors
type ExtractorRegistry struct {
	extractors []Extractor
}

// NewExtractorRegistry creates a new extractor registry
func NewExtractorRegistry() *ExtractorRegistry {
	return &ExtractorRegistry{
		extractors: []Extractor{},
	}
}

// Register registers a new extractor
func (r *ExtractorRegistry) Register(extractor Extractor) {
	r.extractors = append(r.extractors, extractor)
}

// RegisterDefaults registers default extractors
func (r *ExtractorRegistry) RegisterDefaults() {
	// Image formats
	r.Register(&TIFFExtractor{})
	r.Register(&OMETIFFExtractor{})
	r.Register(&ZeissExtractor{}) // .czi
	r.Register(&NikonExtractor{}) // .nd2
	r.Register(&LeicaExtractor{}) // .lif

	// Sequencing formats
	r.Register(&FASTQExtractor{})
	r.Register(&BAMExtractor{})

	// Mass spec formats
	r.Register(&MzMLExtractor{})
	r.Register(&MGFExtractor{})

	// Other formats
	r.Register(&HDF5Extractor{})
	r.Register(&ZarrExtractor{})
	r.Register(&DICOMExtractor{})
	r.Register(&FCSExtractor{}) // Flow cytometry

	// Generic fallback
	r.Register(&GenericExtractor{})
}

// FindExtractor finds an extractor for the given filename
func (r *ExtractorRegistry) FindExtractor(filename string) Extractor {
	for _, extractor := range r.extractors {
		if extractor.CanHandle(filename) {
			return extractor
		}
	}
	return nil
}

// Extract extracts metadata using the appropriate extractor
func (r *ExtractorRegistry) Extract(filepath string) (map[string]interface{}, error) {
	filename := filepath
	extractor := r.FindExtractor(filename)
	if extractor == nil {
		return nil, fmt.Errorf("no extractor found for file: %s", filename)
	}

	return extractor.Extract(filepath)
}

// ListExtractors returns all registered extractors
func (r *ExtractorRegistry) ListExtractors() []ExtractorInfo {
	var info []ExtractorInfo
	for _, ext := range r.extractors {
		info = append(info, ExtractorInfo{
			Name:    ext.Name(),
			Formats: ext.SupportedFormats(),
		})
	}
	return info
}

// ExtractorInfo contains extractor information
type ExtractorInfo struct {
	Name    string   `json:"name"`
	Formats []string `json:"formats"`
}

// --- TIFF Extractor ---

type TIFFExtractor struct{}

func (e *TIFFExtractor) Name() string {
	return "TIFF"
}

func (e *TIFFExtractor) SupportedFormats() []string {
	return []string{".tif", ".tiff"}
}

func (e *TIFFExtractor) CanHandle(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, format := range e.SupportedFormats() {
		if ext == format {
			return true
		}
	}
	return false
}

func (e *TIFFExtractor) Extract(filepath string) (map[string]interface{}, error) {
	// TODO: Implement TIFF metadata extraction using tiff package
	// Extract:
	// - Image dimensions
	// - Bit depth
	// - Creation timestamp
	// - EXIF data if present
	// - Software used

	metadata := map[string]interface{}{
		"format": "TIFF",
		// Add extracted fields here
	}

	return metadata, nil
}

func (e *TIFFExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
	// TODO: Implement
	return nil, fmt.Errorf("not implemented")
}

// --- OME-TIFF Extractor ---

type OMETIFFExtractor struct{}

func (e *OMETIFFExtractor) Name() string {
	return "OME-TIFF"
}

func (e *OMETIFFExtractor) SupportedFormats() []string {
	return []string{".ome.tif", ".ome.tiff"}
}

func (e *OMETIFFExtractor) CanHandle(filename string) bool {
	lower := strings.ToLower(filename)
	return strings.HasSuffix(lower, ".ome.tif") || strings.HasSuffix(lower, ".ome.tiff")
}

func (e *OMETIFFExtractor) Extract(filepath string) (map[string]interface{}, error) {
	// TODO: Implement OME-TIFF metadata extraction
	// Parse OME-XML embedded in TIFF
	// Extract:
	// - Instrument information
	// - Objective details
	// - Channel information
	// - Dimensions (X, Y, Z, T, C)
	// - Pixel size
	// - Acquisition parameters

	metadata := map[string]interface{}{
		"format": "OME-TIFF",
	}

	return metadata, nil
}

func (e *OMETIFFExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

// --- Zeiss CZI Extractor ---

type ZeissExtractor struct{}

func (e *ZeissExtractor) Name() string {
	return "Zeiss CZI"
}

func (e *ZeissExtractor) SupportedFormats() []string {
	return []string{".czi"}
}

func (e *ZeissExtractor) CanHandle(filename string) bool {
	return strings.ToLower(filepath.Ext(filename)) == ".czi"
}

func (e *ZeissExtractor) Extract(filepath string) (map[string]interface{}, error) {
	// TODO: Implement CZI metadata extraction
	// Use bioformats or custom parser
	// Extract:
	// - Microscope model (e.g., "Zeiss LSM 980")
	// - Objective (magnification, NA, immersion)
	// - Channels (names, wavelengths, exposure)
	// - Dimensions
	// - Pixel size
	// - Acquisition date/time
	// - Operator (if recorded)

	metadata := map[string]interface{}{
		"format":                  "CZI",
		"microscope_manufacturer": "Zeiss",
		// Add extracted fields
	}

	return metadata, nil
}

func (e *ZeissExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

// --- Nikon ND2 Extractor ---

type NikonExtractor struct{}

func (e *NikonExtractor) Name() string {
	return "Nikon ND2"
}

func (e *NikonExtractor) SupportedFormats() []string {
	return []string{".nd2"}
}

func (e *NikonExtractor) CanHandle(filename string) bool {
	return strings.ToLower(filepath.Ext(filename)) == ".nd2"
}

func (e *NikonExtractor) Extract(filepath string) (map[string]interface{}, error) {
	// TODO: Implement ND2 metadata extraction
	metadata := map[string]interface{}{
		"format":                  "ND2",
		"microscope_manufacturer": "Nikon",
	}
	return metadata, nil
}

func (e *NikonExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

// --- Leica LIF Extractor ---

type LeicaExtractor struct{}

func (e *LeicaExtractor) Name() string {
	return "Leica LIF"
}

func (e *LeicaExtractor) SupportedFormats() []string {
	return []string{".lif"}
}

func (e *LeicaExtractor) CanHandle(filename string) bool {
	return strings.ToLower(filepath.Ext(filename)) == ".lif"
}

func (e *LeicaExtractor) Extract(filepath string) (map[string]interface{}, error) {
	// TODO: Implement LIF metadata extraction
	metadata := map[string]interface{}{
		"format":                  "LIF",
		"microscope_manufacturer": "Leica",
	}
	return metadata, nil
}

func (e *LeicaExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

// --- FASTQ Extractor ---

type FASTQExtractor struct{}

func (e *FASTQExtractor) Name() string {
	return "FASTQ"
}

func (e *FASTQExtractor) SupportedFormats() []string {
	return []string{".fastq", ".fq", ".fastq.gz", ".fq.gz"}
}

func (e *FASTQExtractor) CanHandle(filename string) bool {
	lower := strings.ToLower(filename)
	for _, format := range e.SupportedFormats() {
		if strings.HasSuffix(lower, format) {
			return true
		}
	}
	return false
}

func (e *FASTQExtractor) Extract(filepath string) (map[string]interface{}, error) {
	// TODO: Implement FASTQ metadata extraction
	// Extract:
	// - Total reads (count sequences)
	// - Read length distribution
	// - Quality score distribution
	// - Detect paired-end from filename pattern

	metadata := map[string]interface{}{
		"format": "FASTQ",
	}

	return metadata, nil
}

func (e *FASTQExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

// --- BAM Extractor ---

type BAMExtractor struct{}

func (e *BAMExtractor) Name() string {
	return "BAM"
}

func (e *BAMExtractor) SupportedFormats() []string {
	return []string{".bam"}
}

func (e *BAMExtractor) CanHandle(filename string) bool {
	return strings.ToLower(filepath.Ext(filename)) == ".bam"
}

func (e *BAMExtractor) Extract(filepath string) (map[string]interface{}, error) {
	// TODO: Implement BAM metadata extraction using htslib bindings
	// Extract header information:
	// - Reference genome
	// - Aligner used
	// - Total reads
	// - Mapped reads

	metadata := map[string]interface{}{
		"format": "BAM",
	}

	return metadata, nil
}

func (e *BAMExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

// --- mzML Extractor (Mass Spec) ---

type MzMLExtractor struct{}

func (e *MzMLExtractor) Name() string {
	return "mzML"
}

func (e *MzMLExtractor) SupportedFormats() []string {
	return []string{".mzml"}
}

func (e *MzMLExtractor) CanHandle(filename string) bool {
	return strings.ToLower(filepath.Ext(filename)) == ".mzml"
}

func (e *MzMLExtractor) Extract(filepath string) (map[string]interface{}, error) {
	// TODO: Parse mzML XML
	// Extract:
	// - Instrument info
	// - Acquisition parameters
	// - Number of spectra

	metadata := map[string]interface{}{
		"format": "mzML",
	}

	return metadata, nil
}

func (e *MzMLExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

// --- MGF Extractor (Mass Spec) ---

type MGFExtractor struct{}

func (e *MGFExtractor) Name() string {
	return "MGF"
}

func (e *MGFExtractor) SupportedFormats() []string {
	return []string{".mgf"}
}

func (e *MGFExtractor) CanHandle(filename string) bool {
	return strings.ToLower(filepath.Ext(filename)) == ".mgf"
}

func (e *MGFExtractor) Extract(filepath string) (map[string]interface{}, error) {
	metadata := map[string]interface{}{
		"format": "MGF",
	}
	return metadata, nil
}

func (e *MGFExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

// --- HDF5 Extractor ---

type HDF5Extractor struct{}

func (e *HDF5Extractor) Name() string {
	return "HDF5"
}

func (e *HDF5Extractor) SupportedFormats() []string {
	return []string{".h5", ".hdf5"}
}

func (e *HDF5Extractor) CanHandle(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".h5" || ext == ".hdf5"
}

func (e *HDF5Extractor) Extract(filepath string) (map[string]interface{}, error) {
	// TODO: Use HDF5 Go library to extract embedded attributes
	metadata := map[string]interface{}{
		"format": "HDF5",
	}
	return metadata, nil
}

func (e *HDF5Extractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

// --- Zarr Extractor ---

type ZarrExtractor struct{}

func (e *ZarrExtractor) Name() string {
	return "Zarr"
}

func (e *ZarrExtractor) SupportedFormats() []string {
	return []string{".zarr"}
}

func (e *ZarrExtractor) CanHandle(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".zarr")
}

func (e *ZarrExtractor) Extract(filepath string) (map[string]interface{}, error) {
	// TODO: Read .zattrs JSON file
	metadata := map[string]interface{}{
		"format": "Zarr",
	}
	return metadata, nil
}

func (e *ZarrExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

// --- DICOM Extractor ---

type DICOMExtractor struct{}

func (e *DICOMExtractor) Name() string {
	return "DICOM"
}

func (e *DICOMExtractor) SupportedFormats() []string {
	return []string{".dcm", ".dicom"}
}

func (e *DICOMExtractor) CanHandle(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".dcm" || ext == ".dicom"
}

func (e *DICOMExtractor) Extract(filepath string) (map[string]interface{}, error) {
	// TODO: Parse DICOM tags
	// IMPORTANT: Redact PHI/PII before storing
	metadata := map[string]interface{}{
		"format": "DICOM",
	}
	return metadata, nil
}

func (e *DICOMExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

// --- FCS Extractor (Flow Cytometry) ---

type FCSExtractor struct{}

func (e *FCSExtractor) Name() string {
	return "FCS"
}

func (e *FCSExtractor) SupportedFormats() []string {
	return []string{".fcs"}
}

func (e *FCSExtractor) CanHandle(filename string) bool {
	return strings.ToLower(filepath.Ext(filename)) == ".fcs"
}

func (e *FCSExtractor) Extract(filepath string) (map[string]interface{}, error) {
	// TODO: Parse FCS header and TEXT segment
	// Extract instrument, channels, events, etc.
	metadata := map[string]interface{}{
		"format": "FCS",
	}
	return metadata, nil
}

func (e *FCSExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

// --- Generic Extractor (Fallback) ---

type GenericExtractor struct{}

func (e *GenericExtractor) Name() string {
	return "Generic"
}

func (e *GenericExtractor) SupportedFormats() []string {
	return []string{"*"}
}

func (e *GenericExtractor) CanHandle(filename string) bool {
	return true // Always matches as fallback
}

func (e *GenericExtractor) Extract(filepath string) (map[string]interface{}, error) {
	// Extract basic file info only
	info, err := getFileInfo(filepath)
	if err != nil {
		return nil, err
	}

	metadata := map[string]interface{}{
		"format":   detectFormat(filepath),
		"filesize": info.Size,
		"modified": info.ModTime,
	}

	return metadata, nil
}

func (e *GenericExtractor) ExtractFromReader(r io.Reader, filename string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"format": detectFormat(filename),
	}, nil
}

// Helper functions

func getFileInfo(filepath string) (FileInfo, error) {
	// TODO: Implement os.Stat wrapper
	return FileInfo{}, nil
}

func detectFormat(filename string) string {
	ext := filepath.Ext(filename)
	if ext != "" {
		return strings.ToUpper(strings.TrimPrefix(ext, "."))
	}
	return "Unknown"
}

// AutoExtractMetadata is a convenience function that extracts metadata
// and wraps it in a Metadata struct
func AutoExtractMetadata(registry *ExtractorRegistry, filepath, schemaName string, user string) (*Metadata, error) {
	extracted, err := registry.Extract(filepath)
	if err != nil {
		return nil, err
	}

	fileInfo, _ := getFileInfo(filepath)

	metadata := &Metadata{
		SchemaName:    schemaName,
		SchemaVersion: "1.0",
		Fields:        extracted,
		FileInfo:      fileInfo,
		Provenance: Provenance{
			UploadedBy: user,
			UploadedAt: time.Now(),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return metadata, nil
}
