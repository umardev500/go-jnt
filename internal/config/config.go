package config

import "path/filepath"

const AuthToken = "9f2eef81b7d5491ab9d9eb12af1a0a39"
const AssetDir = `C:\Users\User\Projects\go-report\assets`

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
