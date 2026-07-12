package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const AuthToken = "0fce2491fd23443bbc0a774956b93e62"
const AssetDir = `C:\Users\User\Projects\go-report\assets`

func GetUnitsFilePath() string {
	baseDir := `C:\Users\User\Projects\go-report` // absolute path
	return filepath.Join(baseDir, "units.xlsx")
}

func GetVehicleFilePath() string {
	baseDir := `C:\Users\User\Projects\go-report` // absolute path
	return filepath.Join(baseDir, "ACTUAL_VEHICLE.xlsx")
}

func GetExportedFilePath() string {
	baseDir := `C:\Users\User\Projects\go-report\public` // absolute path
	return filepath.Join(baseDir, "exported_file.xlsx")
}

func GetReportFilePath() string {
	baseDir := `C:\Users\User\Projects\go-report\public` // absolute path
	return filepath.Join(baseDir, "report.xlsx")
}

func GetManifestFilePath() string {
	baseDir := `C:\Users\User\Projects\go-report\public` // absolute path
	return filepath.Join(baseDir, "manifest_template.xlsx")
}

func GetSKOFilePath() string {
	baseDir := `C:\Users\User\Projects\go-report\assets` // absolute path
	return filepath.Join(baseDir, "sko.json")
}

func GetSJValidatorFilePath() string {
	baseDir := `C:\Users\User\Projects\go-report\public` // absolute path
	return filepath.Join(baseDir, "sj.xlsx")
}

func GetSuratJalanFilePath() string {
	baseDir := `C:\Users\User\Projects\go-report\public` // absolute path
	return filepath.Join(baseDir, "surat_jalan.xlsm")
}

func GetSuratJalanPrintFilePath() string {
	baseDir := `C:\Users\User\Projects\go-report\public` // absolute path
	return filepath.Join(baseDir, "surat_jalan_updated.xlsm")
}

func GetHistoryFilePath() string {
	baseDir := `C:\Users\User\Projects\go-report\assets` // absolute path
	return filepath.Join(baseDir, "sj_history.xlsx")
}

// ================= YAML STRUCT =================

type Config struct {
	Token   string               `yaml:"token"`
	Vendors map[string][]Vehicle `yaml:"vendors"`
}

type Vehicle struct {
	Plate string `yaml:"plate"`
	Type  string `yaml:"type"`
}

// ================= LOAD CONFIG =================

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// ================= EXTRACT PLATES (FLATTEN) =================

func ExtractPlates(cfg *Config) []string {
	var plates []string

	for _, vehicles := range cfg.Vendors {
		for _, v := range vehicles {
			plates = append(plates, v.Plate)
		}
	}

	return plates
}

// ================= OPTIONAL: FLAT VIEW (FOR REPORT) =================

type FlatVehicle struct {
	Vendor string
	Plate  string
	Type   string
}

func Flatten(cfg *Config) []FlatVehicle {
	var out []FlatVehicle

	for vendor, vehicles := range cfg.Vendors {
		for _, v := range vehicles {
			out = append(out, FlatVehicle{
				Vendor: vendor,
				Plate:  v.Plate,
				Type:   v.Type,
			})
		}
	}

	return out
}
