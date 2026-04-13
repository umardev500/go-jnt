package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/umardev500/jnt-report/internal/config"
	"github.com/umardev500/jnt-report/internal/detail"
	"github.com/umardev500/jnt-report/internal/downloader"
	"github.com/umardev500/jnt-report/internal/excel"
	"github.com/umardev500/jnt-report/internal/manifest"
	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
)

type CodeMap map[string][]string

func main() {
	codesBySheet, err := LoadCodes("manifest.yml")
	if err != nil {
		fmt.Println("failed to load yaml:", err)
		return
	}

	fmt.Println(codesBySheet)

	if len(os.Args) < 3 {
		fmt.Println("Usage: app <sheet> <code>")
		return
	}

	sheet := os.Args[1]
	code := os.Args[2]

	fmt.Println("Sheet:", sheet)
	fmt.Println("Code:", code)

	// code := "DBGX26041301196" // JAT
	// sheet := "JKT"

	// Step 1: Download shipment file
	log.Info().Msgf("Downloading shipment file... %s", code)
	if err := downloader.DownloadShipmentFile(code); err != nil {
		fmt.Println("Error downloading shipment file:", err)
		return
	}
	fmt.Println("Shipment file downloaded successfully!")

	pivotExcel := excel.Open("public/exported_file.xlsx")
	pivotSheet := "Memuat dan membongkar ekspor in"

	rows, err := pivotExcel.File.GetRows(pivotSheet)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get pivot data")
	}
	pivotMap := manifest.BuildPivot(rows)
	result := manifest.ToPivotRows(pivotMap)

	blankBagging := make([]manifest.PivotRow, 1)

	i := 0
	total := 0
	for _, r := range result {
		total += r.Count

		// ❌ skip blank bagging in loop
		if r.Bagging == "" {
			blankBagging[0] = r
			fmt.Println(r.Count)
			continue
		}

		// fmt.Printf("Row %d - Bagging: %s | Count: %d\n", i+1, r.Bagging, r.Count)
		i++
	}

	fmt.Println(total)

	dt, err := detail.GetShipmentDetail(code, config.AuthToken)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get shipment detail")
	}

	ex := excel.Open(config.GetManifestFilePath())
	f, err := manifest.GenerateManifest(ex, dt, result, sheet)

	resultFile, err := excelize.OpenFile("public/manifest_out.xlsx")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open result file")
	}
	manifestResultFile, err := manifest.AppendGeneratedToResult(f, resultFile, sheet)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to append generated to result")
	}

	if err := manifestResultFile.SaveAs("public/manifest_out.xlsx"); err != nil {
		log.Fatal().Err(err).Msg("failed to save result file")
	}

	if err := f.SaveAs("public/manifest_generated.xlsx"); err != nil {
		log.Fatal().Err(err).Msg("failed to save manifest")
	}

}

func LoadCodes(path string) (CodeMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m CodeMap
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return m, nil
}
