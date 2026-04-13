package reporting

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

type ReportSummary struct {
	FilePath          string `json:"file_path"`
	SheetName         string `json:"sheet_name"`
	TotalRows         int    `json:"total_rows"`
	TotalPivotRows    int    `json:"total_pivot_rows"`
	TotalWaybillCount int    `json:"total_waybill_count"`
	BlankBaggingCount int    `json:"blank_bagging_count"`
	OutputFile        string `json:"output_file"`
}

func GeneratePivotReport(filePath, sheetName string, createCopy bool) (*ReportSummary, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	rowCount := len(rows)

	dataRange := fmt.Sprintf("%s!A1:B%d", sheetName, rowCount)
	pivotRange := fmt.Sprintf("%s!J1:L%d", sheetName, rowCount+10)

	counts := make(map[string]int)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		bagging := ""
		if len(row) > 1 {
			bagging = row[1]
		}
		waybill := ""
		if len(row) > 0 {
			waybill = row[0]
		}
		if waybill != "" {
			counts[bagging]++
		}
	}

	totalRowsInPivot := 0
	totalSumCountWaybill := 0
	blankCount := 0
	for bagging, c := range counts {
		if bagging == "" {
			blankCount = c
			continue
		}
		totalRowsInPivot++
		totalSumCountWaybill += c
	}

	report := &ReportSummary{
		FilePath:          filePath,
		SheetName:         sheetName,
		TotalRows:         len(rows),
		TotalPivotRows:    totalRowsInPivot,
		TotalWaybillCount: totalSumCountWaybill,
		BlankBaggingCount: blankCount,
		OutputFile:        filePath,
	}

	// --------------------------
	// Create pivot table in Excel
	// --------------------------
	err = f.AddPivotTable(&excelize.PivotTableOptions{
		DataRange:       dataRange,
		PivotTableRange: pivotRange,
		Rows: []excelize.PivotTableField{
			{Data: "Kepemilikan No. Bagging"},
		},
		Data: []excelize.PivotTableField{
			{
				Data:     "No. Waybill",
				Name:     "Count of Waybill",
				Subtotal: "Count",
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if createCopy {

		// Remove original columns A-G to clean sheet
		if err := f.SetColVisible(sheetName, "A:G", false); err != nil {
			log.Fatal(err)
		}

		// Save pivot table
		if err := f.SaveAs(report.OutputFile); err != nil {
			return nil, err
		}

		log.Println("Pivot report generated:", report.OutputFile)
	}

	return report, nil
}
