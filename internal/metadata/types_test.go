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
	"encoding/json"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestMicroscopyMetadata_JSONSerialization(t *testing.T) {
	original := &MicroscopyMetadata{
		Manufacturer:    "Zeiss",
		Model:           "LSM 880",
		SerialNumber:    "Z880-12345",
		SoftwareVersion: "ZEN 3.2",
		Modality:        "confocal",
		Magnification:   63.0,
		NumericalAperture: 1.4,
		Objective:       "Plan-Apochromat 63x/1.4 Oil",
		Width:           1024,
		Height:          1024,
		Depth:           50,
		Channels:        3,
		Timepoints:      10,
		PixelSizeX:      0.13,
		PixelSizeY:      0.13,
		PixelSizeZ:      0.5,
		VoxelSizeUnit:   "micrometers",
		ChannelInfo: []MicroscopyChannel{
			{
				Name:                 "DAPI",
				Index:                0,
				Fluorophore:          "DAPI",
				ExcitationWavelength: 405,
				EmissionWavelength:   461,
				LaserPower:           2.5,
				Color:                "cyan",
			},
			{
				Name:                 "GFP",
				Index:                1,
				Fluorophore:          "EGFP",
				ExcitationWavelength: 488,
				EmissionWavelength:   509,
				LaserPower:           5.0,
				Color:                "green",
			},
		},
		ExposureTime:    100.0,
		FrameRate:       10.0,
		DetectorGain:    800.0,
		DetectorModel:   "Hamamatsu ORCA-Flash4.0",
		ExperimentName:  "Mitosis Time-Lapse",
		SampleID:        "S042",
		Organism:        "Homo sapiens",
		Tissue:          "HeLa cells",
		CellLine:        "HeLa",
		Treatment:       "nocodazole 100nM",
		AcquisitionDate: time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC),
		Operator:        "jsmith",
		Temperature:     37.0,
		CO2Level:        5.0,
		Humidity:        95.0,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	// Test JSON unmarshaling
	var decoded MicroscopyMetadata
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// Verify key fields
	if decoded.Manufacturer != original.Manufacturer {
		t.Errorf("Manufacturer mismatch: got %s, want %s", decoded.Manufacturer, original.Manufacturer)
	}
	if decoded.Model != original.Model {
		t.Errorf("Model mismatch: got %s, want %s", decoded.Model, original.Model)
	}
	if decoded.Width != original.Width {
		t.Errorf("Width mismatch: got %d, want %d", decoded.Width, original.Width)
	}
	if decoded.Channels != original.Channels {
		t.Errorf("Channels mismatch: got %d, want %d", decoded.Channels, original.Channels)
	}
	if len(decoded.ChannelInfo) != len(original.ChannelInfo) {
		t.Errorf("ChannelInfo length mismatch: got %d, want %d", len(decoded.ChannelInfo), len(original.ChannelInfo))
	}
	if decoded.PixelSizeX != original.PixelSizeX {
		t.Errorf("PixelSizeX mismatch: got %f, want %f", decoded.PixelSizeX, original.PixelSizeX)
	}
}

func TestMicroscopyMetadata_YAMLSerialization(t *testing.T) {
	original := &MicroscopyMetadata{
		Manufacturer:  "Nikon",
		Model:         "Ti2-E",
		Modality:      "widefield",
		Magnification: 40.0,
		Width:         2048,
		Height:        2048,
		Channels:      2,
		PixelSizeX:    0.16,
		PixelSizeY:    0.16,
	}

	// Test YAML marshaling
	yamlData, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("YAML marshal failed: %v", err)
	}

	// Test YAML unmarshaling
	var decoded MicroscopyMetadata
	if err := yaml.Unmarshal(yamlData, &decoded); err != nil {
		t.Fatalf("YAML unmarshal failed: %v", err)
	}

	// Verify key fields
	if decoded.Manufacturer != original.Manufacturer {
		t.Errorf("Manufacturer mismatch: got %s, want %s", decoded.Manufacturer, original.Manufacturer)
	}
	if decoded.Width != original.Width {
		t.Errorf("Width mismatch: got %d, want %d", decoded.Width, original.Width)
	}
}

func TestSequencingMetadata_JSONSerialization(t *testing.T) {
	original := &SequencingMetadata{
		Platform:        "Illumina",
		Model:           "NovaSeq 6000",
		SerialNumber:    "NS6000-789",
		SoftwareVersion: "RTA 3.4.4",
		RunID:           "210115_NS6000_0123",
		FlowcellID:      "H7GJKDSXY",
		Lane:            1,
		RunDate:         time.Date(2025, 1, 15, 8, 0, 0, 0, time.UTC),
		Operator:        "jdoe",
		LibraryID:       "LIB042",
		LibraryKit:      "TruSeq Stranded mRNA",
		InsertSize:      350,
		IndexSequence:   "AATTCCGG",
		ReadType:        "paired-end",
		ReadLength:      150,
		Read1Length:     150,
		Read2Length:     150,
		IndexLength:     8,
		TotalReads:      500000000,
		PassFilterReads: 480000000,
		QualityScore:    36.5,
		PercentQ30:      92.3,
		DuplicationRate: 15.2,
		SampleID:        "S042",
		SampleName:      "Brain_Tumor_Patient1",
		Organism:        "Homo sapiens",
		Tissue:          "glioblastoma",
		ReferenceGenome: "GRCh38",
		AssayType:       "RNA-Seq",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	// Test JSON unmarshaling
	var decoded SequencingMetadata
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// Verify key fields
	if decoded.Platform != original.Platform {
		t.Errorf("Platform mismatch: got %s, want %s", decoded.Platform, original.Platform)
	}
	if decoded.Model != original.Model {
		t.Errorf("Model mismatch: got %s, want %s", decoded.Model, original.Model)
	}
	if decoded.TotalReads != original.TotalReads {
		t.Errorf("TotalReads mismatch: got %d, want %d", decoded.TotalReads, original.TotalReads)
	}
	if decoded.ReadLength != original.ReadLength {
		t.Errorf("ReadLength mismatch: got %d, want %d", decoded.ReadLength, original.ReadLength)
	}
	if decoded.PercentQ30 != original.PercentQ30 {
		t.Errorf("PercentQ30 mismatch: got %f, want %f", decoded.PercentQ30, original.PercentQ30)
	}
}

func TestSequencingMetadata_YAMLSerialization(t *testing.T) {
	original := &SequencingMetadata{
		Platform:   "PacBio",
		Model:      "Sequel IIe",
		ReadType:   "single-molecule",
		TotalReads: 1000000,
		Organism:   "Escherichia coli",
		AssayType:  "WGS",
	}

	// Test YAML marshaling
	yamlData, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("YAML marshal failed: %v", err)
	}

	// Test YAML unmarshaling
	var decoded SequencingMetadata
	if err := yaml.Unmarshal(yamlData, &decoded); err != nil {
		t.Fatalf("YAML unmarshal failed: %v", err)
	}

	// Verify key fields
	if decoded.Platform != original.Platform {
		t.Errorf("Platform mismatch: got %s, want %s", decoded.Platform, original.Platform)
	}
	if decoded.TotalReads != original.TotalReads {
		t.Errorf("TotalReads mismatch: got %d, want %d", decoded.TotalReads, original.TotalReads)
	}
}

func TestMassSpecMetadata_JSONSerialization(t *testing.T) {
	original := &MassSpecMetadata{
		Manufacturer:    "Thermo Fisher Scientific",
		Model:           "Q Exactive HF",
		InstrumentType:  "Orbitrap",
		IonizationMode:  "ESI",
		Polarity:        "positive",
		MassAnalyzer:    "Orbitrap",
		ScanRange:       "300-2000",
		Resolution:      70000,
		TotalSpectra:    50000,
		MS1Spectra:      10000,
		MS2Spectra:      40000,
		ChromatographyType: "UPLC",
		ColumnType:      "C18",
		ColumnLength:    150.0,
		FlowRate:        300.0,
		RunTime:         60.0,
		SampleID:        "P042",
		SampleType:      "protein digest",
		Organism:        "Homo sapiens",
		AcquisitionDate: time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC),
		Operator:        "mchen",
		ExperimentType:  "proteomics",
		AcquisitionMode: "DDA",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	// Test JSON unmarshaling
	var decoded MassSpecMetadata
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// Verify key fields
	if decoded.Manufacturer != original.Manufacturer {
		t.Errorf("Manufacturer mismatch: got %s, want %s", decoded.Manufacturer, original.Manufacturer)
	}
	if decoded.Resolution != original.Resolution {
		t.Errorf("Resolution mismatch: got %d, want %d", decoded.Resolution, original.Resolution)
	}
	if decoded.TotalSpectra != original.TotalSpectra {
		t.Errorf("TotalSpectra mismatch: got %d, want %d", decoded.TotalSpectra, original.TotalSpectra)
	}
}

func TestFlowCytometryMetadata_JSONSerialization(t *testing.T) {
	original := &FlowCytometryMetadata{
		Manufacturer:    "BD Biosciences",
		Model:           "FACSAria III",
		SoftwareVersion: "FACSDiva 8.0",
		TotalEvents:     100000,
		EventRate:       5000.0,
		AcquisitionTime: 20.0,
		AcquisitionDate: time.Date(2025, 1, 15, 11, 0, 0, 0, time.UTC),
		Parameters: []FlowCytometryParameter{
			{
				Name:         "FSC-A",
				Description:  "Forward Scatter Area",
				Range:        262144,
				Bits:         18,
				Gain:         1.0,
			},
			{
				Name:         "PE-A",
				Description:  "PE Area",
				Range:        262144,
				Bits:         18,
				Voltage:      450.0,
				Filter:       "575/26",
				Fluorochrome: "PE",
			},
		},
		SampleID:   "FC042",
		Organism:   "Mus musculus",
		CellType:   "splenocytes",
		Treatment:  "LPS 100ng/mL",
		Operator:   "kwong",
		Populations: []string{"Lymphocytes", "Monocytes", "Granulocytes"},
		GatingStrategy: "FSC/SSC → Singlets → Live/Dead → CD4+/CD8+",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	// Test JSON unmarshaling
	var decoded FlowCytometryMetadata
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// Verify key fields
	if decoded.Manufacturer != original.Manufacturer {
		t.Errorf("Manufacturer mismatch: got %s, want %s", decoded.Manufacturer, original.Manufacturer)
	}
	if decoded.TotalEvents != original.TotalEvents {
		t.Errorf("TotalEvents mismatch: got %d, want %d", decoded.TotalEvents, original.TotalEvents)
	}
	if len(decoded.Parameters) != len(original.Parameters) {
		t.Errorf("Parameters length mismatch: got %d, want %d", len(decoded.Parameters), len(original.Parameters))
	}
	if len(decoded.Populations) != len(original.Populations) {
		t.Errorf("Populations length mismatch: got %d, want %d", len(decoded.Populations), len(original.Populations))
	}
}

func TestCryoEMMetadata_JSONSerialization(t *testing.T) {
	original := &CryoEMMetadata{
		Manufacturer:    "Thermo Fisher Scientific",
		Model:           "Titan Krios",
		Voltage:         300,
		DetectorModel:   "Gatan K3",
		DetectorMode:    "super-resolution",
		PixelSize:       0.825,
		Magnification:   130000.0,
		Defocus:         -1.5,
		ExposureTime:    2.5,
		TotalDose:       50.0,
		FramesPerMovie:  40,
		TotalMovies:     5000,
		SampleID:        "EM042",
		Protein:         "20S proteasome",
		Organism:        "Thermoplasma acidophilum",
		GridType:        "Quantifoil R1.2/1.3",
		FreezingMethod:  "plunge freezing",
		AcquisitionDate: time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC),
		Operator:        "alee",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	// Test JSON unmarshaling
	var decoded CryoEMMetadata
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// Verify key fields
	if decoded.Model != original.Model {
		t.Errorf("Model mismatch: got %s, want %s", decoded.Model, original.Model)
	}
	if decoded.Voltage != original.Voltage {
		t.Errorf("Voltage mismatch: got %d, want %d", decoded.Voltage, original.Voltage)
	}
	if decoded.TotalMovies != original.TotalMovies {
		t.Errorf("TotalMovies mismatch: got %d, want %d", decoded.TotalMovies, original.TotalMovies)
	}
	if decoded.TotalDose != original.TotalDose {
		t.Errorf("TotalDose mismatch: got %f, want %f", decoded.TotalDose, original.TotalDose)
	}
}

func TestXRayMetadata_JSONSerialization(t *testing.T) {
	original := &XRayMetadata{
		Facility:        "APS",
		Beamline:        "24-ID-C",
		DetectorModel:   "Pilatus3 6M",
		Wavelength:      0.9792,
		Energy:          12.658,
		Distance:        250.0,
		ExposureTime:    0.1,
		Oscillation:     0.1,
		TotalFrames:     3600,
		SpaceGroup:      "P212121",
		UnitCellA:       78.9,
		UnitCellB:       95.2,
		UnitCellC:       104.3,
		UnitCellAlpha:   90.0,
		UnitCellBeta:    90.0,
		UnitCellGamma:   90.0,
		SampleID:        "XRD042",
		Protein:         "lysozyme",
		Organism:        "Gallus gallus",
		AcquisitionDate: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
		Operator:        "rgarcia",
		Resolution:      1.8,
		Completeness:    99.5,
		Rmerge:          0.065,
		IOverSigma:      25.3,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	// Test JSON unmarshaling
	var decoded XRayMetadata
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// Verify key fields
	if decoded.Facility != original.Facility {
		t.Errorf("Facility mismatch: got %s, want %s", decoded.Facility, original.Facility)
	}
	if decoded.SpaceGroup != original.SpaceGroup {
		t.Errorf("SpaceGroup mismatch: got %s, want %s", decoded.SpaceGroup, original.SpaceGroup)
	}
	if decoded.Resolution != original.Resolution {
		t.Errorf("Resolution mismatch: got %f, want %f", decoded.Resolution, original.Resolution)
	}
	if decoded.TotalFrames != original.TotalFrames {
		t.Errorf("TotalFrames mismatch: got %d, want %d", decoded.TotalFrames, original.TotalFrames)
	}
}

func TestMicroscopyChannel_JSONRoundTrip(t *testing.T) {
	original := &MicroscopyChannel{
		Name:                 "Cy5",
		Index:                2,
		Fluorophore:          "Cy5",
		ExcitationWavelength: 649,
		EmissionWavelength:   670,
		LaserPower:           10.0,
		Color:                "magenta",
		ContrastMin:          100.0,
		ContrastMax:          4095.0,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	// Test JSON unmarshaling
	var decoded MicroscopyChannel
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// Verify all fields
	if decoded.Name != original.Name {
		t.Errorf("Name mismatch: got %s, want %s", decoded.Name, original.Name)
	}
	if decoded.Index != original.Index {
		t.Errorf("Index mismatch: got %d, want %d", decoded.Index, original.Index)
	}
	if decoded.ExcitationWavelength != original.ExcitationWavelength {
		t.Errorf("ExcitationWavelength mismatch: got %d, want %d", decoded.ExcitationWavelength, original.ExcitationWavelength)
	}
	if decoded.LaserPower != original.LaserPower {
		t.Errorf("LaserPower mismatch: got %f, want %f", decoded.LaserPower, original.LaserPower)
	}
}

func TestFlowCytometryParameter_JSONRoundTrip(t *testing.T) {
	original := &FlowCytometryParameter{
		Name:         "APC-A",
		Description:  "APC Area",
		Range:        262144,
		Bits:         18,
		Gain:         1.5,
		Voltage:      500.0,
		Filter:       "660/20",
		Fluorochrome: "APC",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	// Test JSON unmarshaling
	var decoded FlowCytometryParameter
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// Verify all fields
	if decoded.Name != original.Name {
		t.Errorf("Name mismatch: got %s, want %s", decoded.Name, original.Name)
	}
	if decoded.Range != original.Range {
		t.Errorf("Range mismatch: got %d, want %d", decoded.Range, original.Range)
	}
	if decoded.Voltage != original.Voltage {
		t.Errorf("Voltage mismatch: got %f, want %f", decoded.Voltage, original.Voltage)
	}
}

// Benchmark tests for performance
func BenchmarkMicroscopyMetadata_JSONMarshal(b *testing.B) {
	metadata := &MicroscopyMetadata{
		Manufacturer: "Zeiss",
		Model:        "LSM 880",
		Width:        1024,
		Height:       1024,
		Channels:     3,
		ChannelInfo: []MicroscopyChannel{
			{Name: "DAPI", Index: 0, ExcitationWavelength: 405, EmissionWavelength: 461},
			{Name: "GFP", Index: 1, ExcitationWavelength: 488, EmissionWavelength: 509},
			{Name: "RFP", Index: 2, ExcitationWavelength: 561, EmissionWavelength: 582},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(metadata)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSequencingMetadata_JSONMarshal(b *testing.B) {
	metadata := &SequencingMetadata{
		Platform:        "Illumina",
		Model:           "NovaSeq 6000",
		TotalReads:      500000000,
		PassFilterReads: 480000000,
		ReadLength:      150,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(metadata)
		if err != nil {
			b.Fatal(err)
		}
	}
}
