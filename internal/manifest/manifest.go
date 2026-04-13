package manifest

import (
	"fmt"
	"strings"
	"time"

	"github.com/umardev500/jnt-report/internal/detail"
	"github.com/umardev500/jnt-report/internal/excel"
	"github.com/xuri/excelize/v2"
)

type PivotRow struct {
	Bagging string
	Count   int
}

func GenerateManifest(ex *excel.ExcelFile, dt *detail.ShipmentDetailResponse, pivotRow []PivotRow, sheetName string) (*excelize.File, error) {
	shipment := dt.Data.ShipmentDetail

	buf, err := ex.File.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	f, err := excelize.OpenReader(buf)
	if err != nil {
		return nil, err
	}

	sheet := sheetName

	// ✅ unified style (bold + border + center)
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})
	if err != nil {
		return nil, err
	}

	seq := 1
	insertRow := 11
	blank := PivotRow{"", 0}
	count := 0

	for _, r := range pivotRow {
		count += r.Count

		// skip blank bagging
		if r.Bagging == "" {
			blank.Count = r.Count
			continue
		}

		if err := insertPivotRow(f, sheet, insertRow, seq, r, style); err != nil {
			fmt.Println(err)
			return nil, err
		}

		insertRow++
		seq++
	}

	if insertRow > 11 {
		endRow := insertRow - 1

		if err := f.MergeCell(sheet, "B11", fmt.Sprintf("B%d", endRow)); err != nil {
			return nil, err
		}
	}

	styleID, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   80, // 🔥 font size here
			Family: "Roboto Condensed",
		},
		Alignment: &excelize.Alignment{
			Horizontal:   "center",
			Vertical:     "center",
			TextRotation: 90, // 🔥 vertical text
		},
	})
	if err != nil {
		return nil, err
	}

	dest := shipment.TmsShipmentStopVOList[1]
	parts := strings.Split(dest.NetworkName, "_")
	destName := fmt.Sprintf("%s %s", parts[0], parts[1])

	if err := f.SetCellValue(sheet, "B11", destName); err != nil {
		return nil, err
	}

	if err := f.SetCellStyle(sheet, "B11", "B11", styleID); err != nil {
		return nil, err
	}

	SetManifestHeader(dt, blank, count, f, sheetName)
	SetManifestFooter(f, sheet, count)

	return f, nil
}

func insertPivotRow(
	f *excelize.File,
	sheet string,
	rowIndex int,
	seq int,
	r PivotRow,
	style int,
) error {

	// insert row
	if err := f.InsertRows(sheet, rowIndex, 1); err != nil {
		return err
	}

	// A = sequence
	_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIndex), seq)

	// D = bagging
	_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIndex), r.Bagging)

	// E = count
	_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", rowIndex), r.Count)

	// F = optional (if needed, keep empty or fixed)
	_ = f.SetCellValue(sheet, fmt.Sprintf("F%d", rowIndex), "EZ")

	// apply style A-F
	start := fmt.Sprintf("A%d", rowIndex)
	end := fmt.Sprintf("F%d", rowIndex)

	return f.SetCellStyle(sheet, start, end, style)
}

func BuildPivot(rows [][]string) map[string]int {

	result := make(map[string]int)

	for i, row := range rows {

		if i == 0 {
			continue // skip header
		}

		waybill := ""
		bagging := ""

		if len(row) > 0 {
			waybill = row[0]
		}
		if len(row) > 1 {
			bagging = row[1]
		}

		// only count valid waybill
		if waybill == "" {
			continue
		}

		result[bagging]++
	}

	return result
}

func ToPivotRows(pivot map[string]int) []PivotRow {

	var result []PivotRow

	for bagging, count := range pivot {

		result = append(result, PivotRow{
			Bagging: bagging,
			Count:   count,
		})
	}

	return result
}

func SetManifestHeader(dt *detail.ShipmentDetailResponse, blank PivotRow, total int, f *excelize.File, sheet string) error {
	shipment := dt.Data.ShipmentDetail

	// Origin to Destination A2
	if err := f.SetCellValue(sheet, "A2", fmt.Sprintf("DATA OUTGOING %s TO %s", shipment.StartCode, shipment.EndCode)); err != nil {
		return err
	}

	// Plate Number → F4
	if err := f.SetCellValue(sheet, "F4", shipment.PlateNumber); err != nil {
		return err
	}

	// Driver Name → F5
	if err := f.SetCellValue(sheet, "F5", shipment.DriverName); err != nil {
		return err
	}

	// Driver Contact → F6
	if err := f.SetCellValue(sheet, "F6", shipment.DriverContact); err != nil {
		return err
	}

	// Shipment No → E8
	if err := f.SetCellValue(sheet, "E8", shipment.ShipmentNo); err != nil {
		return err
	}

	// Carrier → A3
	if err := f.SetCellValue(sheet, "A3", shipment.CarrierName); err != nil {
		return err
	}

	// Vehicle Type → A6
	if err := f.SetCellValue(sheet, "B6", shipment.VehicleTypeName); err != nil {
		return err
	}

	// Clock out
	actualDepartureTime := shipment.ActualDepartureTime

	actualDepartureTimeParsed, err := time.Parse("2006-01-02 15:04:05", actualDepartureTime)
	if err != nil {
		return err
	}
	fmt.Println(actualDepartureTime)

	if err := f.SetCellValue(sheet, "B5", actualDepartureTimeParsed.Format("15:04")); err != nil {
		return err
	}

	// Date
	if err := f.SetCellValue(sheet, "B4", actualDepartureTimeParsed.Format("2006-01-02")); err != nil {
		return err
	}

	// Blank E10
	if err := f.SetCellValue(sheet, "E10", blank.Count); err != nil {
		return err
	}

	return nil
}

func SetManifestFooter(
	f *excelize.File,
	sheet string,
	total int,
) error {

	rows, err := f.GetRows(sheet)
	if err != nil {
		return err
	}

	lastRow := len(rows)
	footerRow := lastRow

	// Optional label in D
	if err := f.SetCellValue(sheet, fmt.Sprintf("D%d", footerRow), "TOTAL"); err != nil {
		return err
	}

	// Total in E
	if err := f.SetCellValue(sheet, fmt.Sprintf("E%d", footerRow), total); err != nil {
		return err
	}

	return nil
}

func AppendGeneratedToResult(
	resultFile *excelize.File,
	generated *excelize.File,
	sheet string,
) (*excelize.File, error) {

	// ---------------------------
	// 1. Ensure sheet exists
	// ---------------------------
	idx, err := resultFile.GetSheetIndex(sheet)
	if err != nil {
		return nil, err
	}

	if idx == -1 {
		resultFile.NewSheet(sheet)
	}

	// ---------------------------
	// 2. Find start row (append position)
	// ---------------------------
	dstRows, err := resultFile.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	margin := 3
	startRow := len(dstRows) + 1 + margin

	// ---------------------------
	// 3. Read generated sheet
	// ---------------------------
	srcRows, err := generated.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	if len(srcRows) == 0 {
		return nil, fmt.Errorf("generated sheet is empty")
	}

	cols := 6 // A–F fixed layout in your manifest

	// ---------------------------
	// 4. Copy values + styles
	// ---------------------------
	for r := 1; r <= len(srcRows); r++ {
		for c := 1; c <= cols; c++ {

			col, _ := excelize.ColumnNumberToName(c)

			srcCell := fmt.Sprintf("%s%d", col, r)
			dstCell := fmt.Sprintf("%s%d", col, startRow+(r-1))

			// value
			val, _ := generated.GetCellValue(sheet, srcCell)

			// style
			styleID, _ := generated.GetCellStyle(sheet, srcCell)

			_ = resultFile.SetCellValue(sheet, dstCell, val)

			if styleID != 0 {
				_ = resultFile.SetCellStyle(sheet, dstCell, dstCell, styleID)
			}
		}
	}

	// ---------------------------
	// 5. Copy merged cells (IMPORTANT for header/footer)
	// ---------------------------
	merges, _ := generated.GetMergeCells(sheet)

	for _, m := range merges {

		start := m.GetStartAxis()
		end := m.GetEndAxis()

		sCol, sRow, _ := excelize.CellNameToCoordinates(start)
		eCol, eRow, _ := excelize.CellNameToCoordinates(end)

		shift := startRow - 1

		newStart := fmt.Sprintf("%s%d",
			colName(sCol),
			sRow+shift,
		)

		newEnd := fmt.Sprintf("%s%d",
			colName(eCol),
			eRow+shift,
		)

		_ = resultFile.MergeCell(sheet, newStart, newEnd)
	}

	return resultFile, nil
}

func colName(n int) string {
	col, _ := excelize.ColumnNumberToName(n)
	return col
}
