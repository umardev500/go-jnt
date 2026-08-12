package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/jnt-report/internal/approval"
	"github.com/umardev500/jnt-report/internal/config"
	"github.com/umardev500/jnt-report/internal/detail"
	"github.com/umardev500/jnt-report/internal/downloader"
	"github.com/umardev500/jnt-report/internal/excel"
	"github.com/umardev500/jnt-report/internal/manifest"
	"github.com/umardev500/jnt-report/internal/pivot"
	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
)

type CodeMap map[string][]string

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		panic(err)
	}

	approval.InitDB()
	if !approval.IsValid() {
		log.Fatal().Msg("Approval not valid")
		return
	}

	if len(os.Args) < 2 {
		panic("usage: go run main.go 2026-04-10")
	}

	dateStr := os.Args[1] // 2026-04-10

	// ✅ validate YYYY-MM-DD
	_, err = time.Parse("2006-01-02", dateStr)
	if err != nil {
		panic("invalid date format, expected YYYY-MM-DD")
	}

	// 🔥 LIST OF ROUTES
	routes := []string{"BDO", "SRG", "SOC", "SUB", "BGR", "SEG", "BKI", "JKT", "JAT", "DPK", "TGL", "TSK", "PTI", "PRO", "MDN", "CRN", "CIM", "JBR", "JOG", "CKP"}
	// routes := []string{"JKT"}
	filename := "manifest.yml"

	if err := manifest.SyncManifest(filename, routes, dateStr, cfg.Token); err != nil {
		log.Error().Err(err).Msg("failed to save routes")
		return
	}

	// Step 1: Load manifest config
	codesBySheet, err := loadCodes("manifest.yml")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load manifest file")
	}

	templatePath := fmt.Sprintf("%s\\template\\manifest_blank.xlsm", config.AssetDir)
	finalPath := fmt.Sprintf("%s\\generated\\manifest\\manifest_%s.xlsm", config.AssetDir, dateStr)
	manifestPath := fmt.Sprintf("%s\\manifest.xlsm", config.AssetDir)

	// 6. Recreate fresh working file (important)
	if err := copyFile(templatePath, manifestPath); err != nil {
		fmt.Printf("reset template: %v", err)
	}

	log.Debug().
		Interface("codes_by_sheet", codesBySheet).
		Msg("loaded codes")

	sheetNum := 0

	for sheet, codes := range codesBySheet {
		sheetNum++
		fmt.Printf("%d. 🔎 Processing sheet: %s\n", sheetNum, sheet)

		count := 0
		for _, code := range codes {
			count++
			fmt.Printf("%d. 🔎 Processing code: %s\n", count, code)

			if err := runJob(sheet, code, codesBySheet, cfg.Token); err != nil {
				log.Error().
					Err(err).
					Str("sheet", sheet).
					Str("code", code).
					Msg("job failed")

				continue // don't stop entire batch
			}
		}
	}

	if err := os.Rename(manifestPath, finalPath); err != nil {
		fmt.Printf("move file: %v", err)
	}
}

func generateManifest(
	data *detail.ShipmentDetailResponse,
	templatePath string,
	outPath string,
	sheet string,
	pivotRows []pivot.PivotRow,
) (*excelize.File, error) {

	log.Info().Msg("loading template...")
	templateFile, err := excelize.OpenFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("open template: %w", err)
	}
	defer templateFile.Close()

	log.Info().Msg("template loaded successfully")

	log.Info().Msg("cloning template...")
	outFile, err := manifest.CloneExcelFile(templateFile)
	if err != nil {
		return nil, fmt.Errorf("clone template: %w", err)
	}
	defer outFile.Close()

	log.Info().Msg("template cloned successfully")

	log.Info().Msg("building manifest...")
	if err := manifest.BuildManifest(templateFile, outFile, sheet, data, pivotRows); err != nil {
		return nil, fmt.Errorf("build manifest: %w", err)
	}
	log.Info().Msg("manifest built successfully")

	log.Info().Msg("saving manifest...")
	if err := outFile.SaveAs(outPath); err != nil {
		return nil, fmt.Errorf("save manifest: %w", err)
	}

	log.Info().Msg("manifest saved successfully")

	return outFile, nil
}

// -------------------- Helpers --------------------

// loadCodes reads and parses the YAML manifest file
func loadCodes(path string) (CodeMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var codes CodeMap
	if err := yaml.Unmarshal(data, &codes); err != nil {
		return nil, err
	}

	return codes, nil
}

// parseArgs validates and returns CLI arguments
func parseArgs() (sheet, code string) {
	if len(os.Args) < 3 {
		log.Fatal().Msg("Usage: app <sheet> <code>")
	}

	return os.Args[1], os.Args[2]
}

// downloadShipment handles shipment file download
func downloadShipment(code string, token string) error {
	log.Info().Str("code", code).Msg("downloading shipment file")
	return downloader.DownloadShipmentFile(code, token)
}

// processPivot reads and transforms pivot Excel data
func processPivot(path, sheet string) ([]pivot.PivotRow, error) {
	file := excel.Open(path)

	rows, err := file.File.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	pivotMap := pivot.BuildPivot(rows)
	return pivot.ToPivotRows(pivotMap), nil
}

// fetchShipmentDetail gets shipment detail from API
func fetchShipmentDetail(code string, token string) (*detail.ShipmentDetailResponse, error) {
	return detail.GetShipmentDetail(code, token)
}

func runJob(sheet, code string, codesBySheet CodeMap, token string) error {
	templateSheet := "TEMPLATE"

	// Step 1: Download shipment file
	if err := downloadShipment(code, token); err != nil {
		return fmt.Errorf("download shipment: %w", err)
	}

	// Step 2: Process pivot Excel
	pivotRows, err := processPivot(
		"public/exported_file.xlsx",
		"Memuat dan membongkar ekspor in",
	)
	if err != nil {
		return fmt.Errorf("process pivot: %w", err)
	}

	// Step 3: Fetch shipment detail
	data, err := fetchShipmentDetail(code, token)
	if err != nil {
		return fmt.Errorf("fetch shipment detail: %w", err)
	}

	log.Info().
		Int("pivot_rows", len(pivotRows)).
		Msg("pivot processed successfully")

	// Step 4: Build manifest
	outFile, err := generateManifest(
		data,
		fmt.Sprintf("%s\\template\\manifest_template.xlsm", config.AssetDir),
		fmt.Sprintf("%s\\generated\\temp\\manifest_%s.xlsx", config.AssetDir, code),
		templateSheet,
		pivotRows,
	)
	if err != nil {
		return fmt.Errorf("generate manifest: %w", err)
	}

	// Step 5: Insert into master manifest
	manifestPath := fmt.Sprintf("%s\\manifest.xlsm", config.AssetDir)

	manifestFile, err := excelize.OpenFile(manifestPath)
	if err != nil {
		return fmt.Errorf("open manifest file: %w", err)
	}

	if err := manifest.InsertManifest(outFile, manifestFile, sheet, templateSheet); err != nil {
		return fmt.Errorf("insert manifest: %w", err)
	}

	if err := manifestFile.Save(); err != nil {
		return fmt.Errorf("save manifest file: %w", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
