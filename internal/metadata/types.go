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

// This file defines instrument-specific metadata structures for different
// scientific domains. These structs extend the base Metadata type with
// domain-specific fields for microscopy, sequencing, mass spectrometry,
// flow cytometry, and other research instruments.

import "time"

// MicroscopyMetadata contains metadata specific to microscopy instruments
type MicroscopyMetadata struct {
	// Instrument information
	Manufacturer  string `json:"manufacturer" yaml:"manufacturer"`                       // e.g., "Zeiss", "Nikon", "Leica"
	Model         string `json:"model" yaml:"model"`                                     // e.g., "LSM 880", "Ti2-E"
	SerialNumber  string `json:"serial_number,omitempty" yaml:"serial_number,omitempty"` // Instrument serial number
	SoftwareVersion string `json:"software_version,omitempty" yaml:"software_version,omitempty"` // Acquisition software version

	// Imaging parameters
	Modality       string  `json:"modality" yaml:"modality"`                             // e.g., "confocal", "widefield", "TIRF", "light sheet"
	Magnification  float64 `json:"magnification,omitempty" yaml:"magnification,omitempty"` // Total magnification (e.g., 40.0)
	NumericalAperture float64 `json:"numerical_aperture,omitempty" yaml:"numerical_aperture,omitempty"` // Objective NA (e.g., 1.4)
	Objective      string  `json:"objective,omitempty" yaml:"objective,omitempty"`       // Objective description

	// Image dimensions
	Width        int     `json:"width" yaml:"width"`                                       // Image width in pixels
	Height       int     `json:"height" yaml:"height"`                                     // Image height in pixels
	Depth        int     `json:"depth,omitempty" yaml:"depth,omitempty"`                  // Z-stack depth
	Channels     int     `json:"channels" yaml:"channels"`                                 // Number of channels
	Timepoints   int     `json:"timepoints,omitempty" yaml:"timepoints,omitempty"`        // Number of timepoints
	PixelSizeX   float64 `json:"pixel_size_x,omitempty" yaml:"pixel_size_x,omitempty"`    // µm per pixel (X)
	PixelSizeY   float64 `json:"pixel_size_y,omitempty" yaml:"pixel_size_y,omitempty"`    // µm per pixel (Y)
	PixelSizeZ   float64 `json:"pixel_size_z,omitempty" yaml:"pixel_size_z,omitempty"`    // µm per slice (Z)
	VoxelSizeUnit string `json:"voxel_size_unit,omitempty" yaml:"voxel_size_unit,omitempty"` // e.g., "micrometers"

	// Channel information
	ChannelInfo []MicroscopyChannel `json:"channel_info,omitempty" yaml:"channel_info,omitempty"`

	// Acquisition settings
	ExposureTime    float64 `json:"exposure_time,omitempty" yaml:"exposure_time,omitempty"`       // Exposure time in ms
	FrameRate       float64 `json:"frame_rate,omitempty" yaml:"frame_rate,omitempty"`             // Frames per second
	BinningX        int     `json:"binning_x,omitempty" yaml:"binning_x,omitempty"`               // Binning factor X
	BinningY        int     `json:"binning_y,omitempty" yaml:"binning_y,omitempty"`               // Binning factor Y
	DetectorGain    float64 `json:"detector_gain,omitempty" yaml:"detector_gain,omitempty"`       // Detector gain
	DetectorModel   string  `json:"detector_model,omitempty" yaml:"detector_model,omitempty"`     // e.g., "Hamamatsu ORCA-Flash4.0"

	// Experiment details
	ExperimentName  string    `json:"experiment_name,omitempty" yaml:"experiment_name,omitempty"`   // User-defined experiment name
	SampleID        string    `json:"sample_id,omitempty" yaml:"sample_id,omitempty"`               // Sample identifier
	Organism        string    `json:"organism,omitempty" yaml:"organism,omitempty"`                 // e.g., "Mus musculus"
	Tissue          string    `json:"tissue,omitempty" yaml:"tissue,omitempty"`                     // e.g., "brain", "liver"
	CellLine        string    `json:"cell_line,omitempty" yaml:"cell_line,omitempty"`               // e.g., "HeLa", "CHO"
	Treatment       string    `json:"treatment,omitempty" yaml:"treatment,omitempty"`               // Experimental treatment
	AcquisitionDate time.Time `json:"acquisition_date,omitempty" yaml:"acquisition_date,omitempty"` // When image was acquired
	Operator        string    `json:"operator,omitempty" yaml:"operator,omitempty"`                 // Person who acquired image

	// Environmental conditions
	Temperature float64 `json:"temperature,omitempty" yaml:"temperature,omitempty"` // °C
	CO2Level    float64 `json:"co2_level,omitempty" yaml:"co2_level,omitempty"`    // % CO2
	Humidity    float64 `json:"humidity,omitempty" yaml:"humidity,omitempty"`      // % relative humidity
}

// MicroscopyChannel describes a single channel/fluorophore in microscopy
type MicroscopyChannel struct {
	Name            string  `json:"name" yaml:"name"`                                           // Channel name (e.g., "DAPI", "GFP")
	Index           int     `json:"index" yaml:"index"`                                         // Channel index (0-based)
	Fluorophore     string  `json:"fluorophore,omitempty" yaml:"fluorophore,omitempty"`         // e.g., "Alexa488", "mCherry"
	ExcitationWavelength int `json:"excitation_wavelength,omitempty" yaml:"excitation_wavelength,omitempty"` // nm
	EmissionWavelength   int `json:"emission_wavelength,omitempty" yaml:"emission_wavelength,omitempty"`     // nm
	LaserPower      float64 `json:"laser_power,omitempty" yaml:"laser_power,omitempty"`         // % or mW
	Color           string  `json:"color,omitempty" yaml:"color,omitempty"`                     // Display color (e.g., "cyan", "magenta")
	ContrastMin     float64 `json:"contrast_min,omitempty" yaml:"contrast_min,omitempty"`       // Display range min
	ContrastMax     float64 `json:"contrast_max,omitempty" yaml:"contrast_max,omitempty"`       // Display range max
}

// SequencingMetadata contains metadata specific to DNA/RNA sequencing
type SequencingMetadata struct {
	// Instrument information
	Platform        string `json:"platform" yaml:"platform"`                                   // e.g., "Illumina", "PacBio", "ONT"
	Model           string `json:"model" yaml:"model"`                                         // e.g., "NovaSeq 6000", "Sequel IIe"
	SerialNumber    string `json:"serial_number,omitempty" yaml:"serial_number,omitempty"`     // Instrument serial number
	SoftwareVersion string `json:"software_version,omitempty" yaml:"software_version,omitempty"` // e.g., "RTA 3.4.4"

	// Run information
	RunID           string    `json:"run_id" yaml:"run_id"`                                       // Unique run identifier
	FlowcellID      string    `json:"flowcell_id,omitempty" yaml:"flowcell_id,omitempty"`         // Flowcell barcode
	Lane            int       `json:"lane,omitempty" yaml:"lane,omitempty"`                       // Lane number (1-8 for NovaSeq)
	RunDate         time.Time `json:"run_date,omitempty" yaml:"run_date,omitempty"`               // When run was performed
	Operator        string    `json:"operator,omitempty" yaml:"operator,omitempty"`               // Person who ran the sequencer

	// Library information
	LibraryID       string `json:"library_id,omitempty" yaml:"library_id,omitempty"`               // Library identifier
	LibraryKit      string `json:"library_kit,omitempty" yaml:"library_kit,omitempty"`             // Kit used for library prep
	LibraryProtocol string `json:"library_protocol,omitempty" yaml:"library_protocol,omitempty"`   // Protocol name/version
	InsertSize      int    `json:"insert_size,omitempty" yaml:"insert_size,omitempty"`             // Average insert size (bp)
	IndexSequence   string `json:"index_sequence,omitempty" yaml:"index_sequence,omitempty"`       // Sample barcode/index

	// Read configuration
	ReadType        string `json:"read_type,omitempty" yaml:"read_type,omitempty"`                 // e.g., "paired-end", "single-end"
	ReadLength      int    `json:"read_length,omitempty" yaml:"read_length,omitempty"`             // Read length in bp
	Read1Length     int    `json:"read1_length,omitempty" yaml:"read1_length,omitempty"`           // R1 length (paired-end)
	Read2Length     int    `json:"read2_length,omitempty" yaml:"read2_length,omitempty"`           // R2 length (paired-end)
	IndexLength     int    `json:"index_length,omitempty" yaml:"index_length,omitempty"`           // Index read length

	// Quality metrics
	TotalReads      int64   `json:"total_reads,omitempty" yaml:"total_reads,omitempty"`             // Total number of reads
	PassFilterReads int64   `json:"pass_filter_reads,omitempty" yaml:"pass_filter_reads,omitempty"` // Reads passing filter
	QualityScore    float64 `json:"quality_score,omitempty" yaml:"quality_score,omitempty"`         // Average Phred score
	PercentQ30      float64 `json:"percent_q30,omitempty" yaml:"percent_q30,omitempty"`             // % bases ≥Q30
	DuplicationRate float64 `json:"duplication_rate,omitempty" yaml:"duplication_rate,omitempty"`   // % duplicate reads

	// Sample information
	SampleID        string `json:"sample_id,omitempty" yaml:"sample_id,omitempty"`                 // Sample identifier
	SampleName      string `json:"sample_name,omitempty" yaml:"sample_name,omitempty"`             // Human-readable name
	Organism        string `json:"organism,omitempty" yaml:"organism,omitempty"`                   // e.g., "Homo sapiens"
	Tissue          string `json:"tissue,omitempty" yaml:"tissue,omitempty"`                       // Tissue type
	CellType        string `json:"cell_type,omitempty" yaml:"cell_type,omitempty"`                 // Cell type
	Treatment       string `json:"treatment,omitempty" yaml:"treatment,omitempty"`                 // Experimental treatment
	ReferenceGenome string `json:"reference_genome,omitempty" yaml:"reference_genome,omitempty"`   // e.g., "GRCh38", "mm10"

	// Experiment type
	AssayType       string `json:"assay_type,omitempty" yaml:"assay_type,omitempty"`               // e.g., "RNA-Seq", "ChIP-Seq", "WGS"
	TargetRegion    string `json:"target_region,omitempty" yaml:"target_region,omitempty"`         // e.g., "exome", "16S rRNA"
	EnrichmentKit   string `json:"enrichment_kit,omitempty" yaml:"enrichment_kit,omitempty"`       // Target enrichment kit
}

// MassSpecMetadata contains metadata specific to mass spectrometry
type MassSpecMetadata struct {
	// Instrument information
	Manufacturer    string `json:"manufacturer" yaml:"manufacturer"`                               // e.g., "Thermo", "Waters", "Bruker"
	Model           string `json:"model" yaml:"model"`                                             // e.g., "Q Exactive HF"
	SerialNumber    string `json:"serial_number,omitempty" yaml:"serial_number,omitempty"`         // Instrument serial number
	SoftwareVersion string `json:"software_version,omitempty" yaml:"software_version,omitempty"`   // Acquisition software version

	// Mass spectrometer type
	InstrumentType  string `json:"instrument_type,omitempty" yaml:"instrument_type,omitempty"`     // e.g., "Orbitrap", "Q-TOF", "Triple Quad"
	IonizationMode  string `json:"ionization_mode,omitempty" yaml:"ionization_mode,omitempty"`     // e.g., "ESI", "MALDI", "APCI"
	Polarity        string `json:"polarity,omitempty" yaml:"polarity,omitempty"`                   // "positive" or "negative"
	MassAnalyzer    string `json:"mass_analyzer,omitempty" yaml:"mass_analyzer,omitempty"`         // e.g., "Orbitrap", "TOF", "Quadrupole"

	// Acquisition parameters
	ScanRange       string  `json:"scan_range,omitempty" yaml:"scan_range,omitempty"`               // e.g., "300-2000 m/z"
	Resolution      int     `json:"resolution,omitempty" yaml:"resolution,omitempty"`               // Mass resolution (e.g., 70000)
	ScanRate        float64 `json:"scan_rate,omitempty" yaml:"scan_rate,omitempty"`                 // Scans per second
	TotalSpectra    int     `json:"total_spectra,omitempty" yaml:"total_spectra,omitempty"`         // Number of spectra
	MS1Spectra      int     `json:"ms1_spectra,omitempty" yaml:"ms1_spectra,omitempty"`             // Number of MS1 scans
	MS2Spectra      int     `json:"ms2_spectra,omitempty" yaml:"ms2_spectra,omitempty"`             // Number of MS2 scans

	// Chromatography (if LC-MS)
	ChromatographyType string  `json:"chromatography_type,omitempty" yaml:"chromatography_type,omitempty"` // e.g., "HPLC", "UPLC"
	ColumnType         string  `json:"column_type,omitempty" yaml:"column_type,omitempty"`                 // e.g., "C18"
	ColumnLength       float64 `json:"column_length,omitempty" yaml:"column_length,omitempty"`             // mm
	FlowRate           float64 `json:"flow_rate,omitempty" yaml:"flow_rate,omitempty"`                     // µL/min
	RunTime            float64 `json:"run_time,omitempty" yaml:"run_time,omitempty"`                       // minutes

	// Sample information
	SampleID        string    `json:"sample_id,omitempty" yaml:"sample_id,omitempty"`                 // Sample identifier
	SampleType      string    `json:"sample_type,omitempty" yaml:"sample_type,omitempty"`             // e.g., "protein digest", "metabolite extract"
	Organism        string    `json:"organism,omitempty" yaml:"organism,omitempty"`                   // Source organism
	PrepMethod      string    `json:"prep_method,omitempty" yaml:"prep_method,omitempty"`             // Sample preparation protocol
	AcquisitionDate time.Time `json:"acquisition_date,omitempty" yaml:"acquisition_date,omitempty"`   // When data was acquired
	Operator        string    `json:"operator,omitempty" yaml:"operator,omitempty"`                   // Person who ran instrument

	// Experiment type
	ExperimentType  string `json:"experiment_type,omitempty" yaml:"experiment_type,omitempty"`       // e.g., "proteomics", "metabolomics", "lipidomics"
	AcquisitionMode string `json:"acquisition_mode,omitempty" yaml:"acquisition_mode,omitempty"`     // e.g., "DDA", "DIA", "targeted"
}

// FlowCytometryMetadata contains metadata specific to flow cytometry
type FlowCytometryMetadata struct {
	// Instrument information
	Manufacturer    string `json:"manufacturer" yaml:"manufacturer"`                               // e.g., "BD", "Beckman Coulter"
	Model           string `json:"model" yaml:"model"`                                             // e.g., "FACSAria III"
	SerialNumber    string `json:"serial_number,omitempty" yaml:"serial_number,omitempty"`         // Instrument serial number
	SoftwareVersion string `json:"software_version,omitempty" yaml:"software_version,omitempty"`   // e.g., "FACSDiva 8.0"

	// Acquisition parameters
	TotalEvents     int       `json:"total_events" yaml:"total_events"`                             // Total events acquired
	EventRate       float64   `json:"event_rate,omitempty" yaml:"event_rate,omitempty"`             // Events per second
	AbortedEvents   int       `json:"aborted_events,omitempty" yaml:"aborted_events,omitempty"`     // Events aborted
	AcquisitionTime float64   `json:"acquisition_time,omitempty" yaml:"acquisition_time,omitempty"` // seconds
	AcquisitionDate time.Time `json:"acquisition_date,omitempty" yaml:"acquisition_date,omitempty"` // When acquired

	// Channels/Parameters
	Parameters      []FlowCytometryParameter `json:"parameters,omitempty" yaml:"parameters,omitempty"` // Measured parameters
	CompensationMatrix [][]float64 `json:"compensation_matrix,omitempty" yaml:"compensation_matrix,omitempty"` // Compensation matrix

	// Sample information
	SampleID        string `json:"sample_id,omitempty" yaml:"sample_id,omitempty"`                 // Sample identifier
	TubeID          string `json:"tube_id,omitempty" yaml:"tube_id,omitempty"`                     // Tube identifier
	Organism        string `json:"organism,omitempty" yaml:"organism,omitempty"`                   // Source organism
	CellType        string `json:"cell_type,omitempty" yaml:"cell_type,omitempty"`                 // Cell type analyzed
	Treatment       string `json:"treatment,omitempty" yaml:"treatment,omitempty"`                 // Experimental treatment
	Operator        string `json:"operator,omitempty" yaml:"operator,omitempty"`                   // Person who ran instrument

	// Gating/Analysis
	Populations     []string `json:"populations,omitempty" yaml:"populations,omitempty"`             // Identified populations
	GatingStrategy  string   `json:"gating_strategy,omitempty" yaml:"gating_strategy,omitempty"`     // Description of gating
}

// FlowCytometryParameter describes a single measured parameter
type FlowCytometryParameter struct {
	Name        string  `json:"name" yaml:"name"`                                       // e.g., "FSC-A", "PE-A"
	Description string  `json:"description,omitempty" yaml:"description,omitempty"`     // e.g., "Forward Scatter Area"
	Range       int     `json:"range,omitempty" yaml:"range,omitempty"`                 // Data range (e.g., 262144)
	Bits        int     `json:"bits,omitempty" yaml:"bits,omitempty"`                   // Bit depth (e.g., 18)
	Gain        float64 `json:"gain,omitempty" yaml:"gain,omitempty"`                   // Detector gain
	Voltage     float64 `json:"voltage,omitempty" yaml:"voltage,omitempty"`             // PMT voltage
	Filter      string  `json:"filter,omitempty" yaml:"filter,omitempty"`               // Optical filter (e.g., "530/30")
	Fluorochrome string `json:"fluorochrome,omitempty" yaml:"fluorochrome,omitempty"`   // Fluorochrome (e.g., "FITC", "PE")
}

// CryoEMMetadata contains metadata specific to cryo-electron microscopy
type CryoEMMetadata struct {
	// Instrument information
	Manufacturer    string `json:"manufacturer" yaml:"manufacturer"`                               // e.g., "Thermo Fisher", "JEOL"
	Model           string `json:"model" yaml:"model"`                                             // e.g., "Titan Krios"
	Voltage         int    `json:"voltage,omitempty" yaml:"voltage,omitempty"`                     // kV (e.g., 300)
	SerialNumber    string `json:"serial_number,omitempty" yaml:"serial_number,omitempty"`         // Instrument serial number

	// Detector
	DetectorModel   string `json:"detector_model,omitempty" yaml:"detector_model,omitempty"`       // e.g., "Gatan K3"
	DetectorMode    string `json:"detector_mode,omitempty" yaml:"detector_mode,omitempty"`         // e.g., "counting", "super-resolution"
	PixelSize       float64 `json:"pixel_size,omitempty" yaml:"pixel_size,omitempty"`              // Å per pixel

	// Acquisition parameters
	Magnification   float64 `json:"magnification,omitempty" yaml:"magnification,omitempty"`         // Nominal magnification
	Defocus         float64 `json:"defocus,omitempty" yaml:"defocus,omitempty"`                     // µm
	ExposureTime    float64 `json:"exposure_time,omitempty" yaml:"exposure_time,omitempty"`         // seconds
	TotalDose       float64 `json:"total_dose,omitempty" yaml:"total_dose,omitempty"`               // e⁻/Å²
	FramesPerMovie  int     `json:"frames_per_movie,omitempty" yaml:"frames_per_movie,omitempty"`   // Number of frames
	TotalMovies     int     `json:"total_movies,omitempty" yaml:"total_movies,omitempty"`           // Number of movies acquired

	// Sample information
	SampleID        string    `json:"sample_id,omitempty" yaml:"sample_id,omitempty"`               // Sample identifier
	Protein         string    `json:"protein,omitempty" yaml:"protein,omitempty"`                   // Protein/complex name
	Organism        string    `json:"organism,omitempty" yaml:"organism,omitempty"`                 // Source organism
	GridType        string    `json:"grid_type,omitempty" yaml:"grid_type,omitempty"`               // e.g., "Quantifoil R1.2/1.3"
	FreezingMethod  string    `json:"freezing_method,omitempty" yaml:"freezing_method,omitempty"`   // e.g., "plunge freezing"
	AcquisitionDate time.Time `json:"acquisition_date,omitempty" yaml:"acquisition_date,omitempty"` // When data was acquired
	Operator        string    `json:"operator,omitempty" yaml:"operator,omitempty"`                 // Person who ran instrument
}

// XRayMetadata contains metadata specific to X-ray diffraction/crystallography
type XRayMetadata struct {
	// Instrument information
	Facility        string `json:"facility,omitempty" yaml:"facility,omitempty"`                   // e.g., "APS", "ESRF"
	Beamline        string `json:"beamline,omitempty" yaml:"beamline,omitempty"`                   // e.g., "24-ID-C"
	DetectorModel   string `json:"detector_model,omitempty" yaml:"detector_model,omitempty"`       // e.g., "Pilatus3 6M"

	// Acquisition parameters
	Wavelength      float64 `json:"wavelength,omitempty" yaml:"wavelength,omitempty"`               // Å
	Energy          float64 `json:"energy,omitempty" yaml:"energy,omitempty"`                       // keV
	Distance        float64 `json:"distance,omitempty" yaml:"distance,omitempty"`                   // Detector distance (mm)
	ExposureTime    float64 `json:"exposure_time,omitempty" yaml:"exposure_time,omitempty"`         // seconds per frame
	Oscillation     float64 `json:"oscillation,omitempty" yaml:"oscillation,omitempty"`             // degrees per frame
	TotalFrames     int     `json:"total_frames,omitempty" yaml:"total_frames,omitempty"`           // Number of frames

	// Crystal information
	SpaceGroup      string  `json:"space_group,omitempty" yaml:"space_group,omitempty"`             // e.g., "P212121"
	UnitCellA       float64 `json:"unit_cell_a,omitempty" yaml:"unit_cell_a,omitempty"`             // Å
	UnitCellB       float64 `json:"unit_cell_b,omitempty" yaml:"unit_cell_b,omitempty"`             // Å
	UnitCellC       float64 `json:"unit_cell_c,omitempty" yaml:"unit_cell_c,omitempty"`             // Å
	UnitCellAlpha   float64 `json:"unit_cell_alpha,omitempty" yaml:"unit_cell_alpha,omitempty"`     // degrees
	UnitCellBeta    float64 `json:"unit_cell_beta,omitempty" yaml:"unit_cell_beta,omitempty"`       // degrees
	UnitCellGamma   float64 `json:"unit_cell_gamma,omitempty" yaml:"unit_cell_gamma,omitempty"`     // degrees

	// Sample information
	SampleID        string    `json:"sample_id,omitempty" yaml:"sample_id,omitempty"`               // Sample identifier
	Protein         string    `json:"protein,omitempty" yaml:"protein,omitempty"`                   // Protein name
	Organism        string    `json:"organism,omitempty" yaml:"organism,omitempty"`                 // Source organism
	AcquisitionDate time.Time `json:"acquisition_date,omitempty" yaml:"acquisition_date,omitempty"` // When data was acquired
	Operator        string    `json:"operator,omitempty" yaml:"operator,omitempty"`                 // Person who collected data

	// Processing information
	Resolution      float64 `json:"resolution,omitempty" yaml:"resolution,omitempty"`               // Å
	Completeness    float64 `json:"completeness,omitempty" yaml:"completeness,omitempty"`           // %
	Rmerge          float64 `json:"rmerge,omitempty" yaml:"rmerge,omitempty"`                       // Merging R-factor
	IOverSigma      float64 `json:"i_over_sigma,omitempty" yaml:"i_over_sigma,omitempty"`           // <I/σ(I)>
}
