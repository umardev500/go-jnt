package manifest

import (
	"fmt"
	"strings"
	"time"

	"github.com/umardev500/jnt-report/internal/detail"
	"github.com/umardev500/jnt-report/internal/excel"
	"github.com/umardev500/jnt-report/internal/pivot"
	"github.com/xuri/excelize/v2"
)

type PivotRow struct {
	Bagging string
	Count   int
}

func GenerateManifest_old(ex *excel.ExcelFile, dt *detail.ShipmentDetailResponse, pivotRow []pivot.PivotRow, sheetName string) (*excelize.File, error) {
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
	blank := pivot.PivotRow{Bagging: "", Count: 0}
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

	// SetManifestHeader(dt, blank, count, f, sheetName)
	// SetManifestFooter(f, sheet, count)

	return f, nil
}

func insertPivotRow(
	f *excelize.File,
	sheet string,
	rowIndex int,
	seq int,
	r pivot.PivotRow,
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

	// Log if count is zero
	if r.Count == 0 {
		fmt.Printf(
			"\033[1;33m[WARN]\033[0m Zero count for Bagging: \033[1;36m%s\033[0m | Count: \033[1;31m%d\033[0m\n",
			r.Bagging,
			r.Count,
		)
	}
	// F = optional (if needed, keep empty or fixed)
	_ = f.SetCellValue(sheet, fmt.Sprintf("F%d", rowIndex), "EZ")

	// apply style A-F
	start := fmt.Sprintf("A%d", rowIndex)
	end := fmt.Sprintf("F%d", rowIndex)

	return f.SetCellStyle(sheet, start, end, style)
}

func SetManifestHeader(dt *detail.ShipmentDetailResponse, blank pivot.PivotRow, total int, f *excelize.File, sheet string) error {
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

func BuildManifest(
	templateFile *excelize.File,
	outFile *excelize.File,
	sheet string,
	data *detail.ShipmentDetailResponse,
	pivotRow []pivot.PivotRow,
) error {
	shipment := data.Data.ShipmentDetail

	// ✅ unified style (bold + border + center)
	style, err := outFile.NewStyle(&excelize.Style{
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
		return err
	}

	seq := 1
	insertRow := 11
	blank := pivot.PivotRow{Bagging: "", Count: 0}
	count := 0

	fmt.Println("--------------------------------------------------")
	fmt.Println("\033[1;34m[INFO]\033[0m Inserting pivot rows... start")
	for _, r := range pivotRow {
		count += r.Count

		// skip blank bagging
		if r.Bagging == "" {
			blank.Bagging = r.Bagging
			blank.Count = r.Count
			continue
		}

		err := insertPivotRow(outFile, sheet, insertRow, seq, r, style)
		if err != nil {
			fmt.Println("\033[1;31m[ERROR]\033[0m Insert pivot row failed")
			return err
		}

		insertRow++
		seq++
	}
	fmt.Printf(
		"\033[1;32m[DONE]\033[0m Inserted pivot rows successfully | total=%d\n",
		insertRow,
	)
	fmt.Println("--------------------------------------------------")

	if insertRow > 11 {
		endRow := insertRow - 1

		if err := outFile.MergeCell(sheet, "B11", fmt.Sprintf("B%d", endRow)); err != nil {
			return err
		}
	}

	styleID, err := outFile.NewStyle(&excelize.Style{
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
		return err
	}

	dest := shipment.TmsShipmentStopVOList[1]
	parts := strings.Split(dest.NetworkName, "_")
	destName := fmt.Sprintf("%s %s", parts[0], parts[1])

	if err := outFile.SetCellValue(sheet, "B11", destName); err != nil {
		return err
	}

	if err := outFile.SetCellStyle(sheet, "B11", "B11", styleID); err != nil {
		return err
	}

	SetManifestHeader(data, blank, count, outFile, sheet)
	SetManifestFooter(outFile, sheet, count)

	return nil
}

func InsertManifest(
	sourceFile *excelize.File,
	destFile *excelize.File,
	sheetDest, templateSheet string,
) error {

	// ---------------------------
	// 1. Ensure destination sheet exists
	// ---------------------------
	idx, err := destFile.GetSheetIndex(sheetDest)
	if err != nil {
		return err
	}
	if idx == -1 {
		destFile.NewSheet(sheetDest)
	}

	// ---------------------------
	// 2. Find start row (append position)
	// ---------------------------
	dstRows, err := destFile.GetRows(sheetDest)
	if err != nil {
		return err
	}

	margin := 4 // change to 2–3 if you want spacing
	startRow := len(dstRows) + 1 + margin

	// ---------------------------
	// 3. Read source (template) sheet
	// ---------------------------
	srcRows, err := sourceFile.GetRows(templateSheet)
	if err != nil {
		return err
	}
	if len(srcRows) == 0 {
		return fmt.Errorf("source sheet is empty")
	}

	cols := 6 // adjust if needed

	// ---------------------------
	// 4. Copy values + styles (FIXED)
	// ---------------------------
	styleCache := map[int]int{} // sourceStyleID → destStyleID

	for r := 1; r <= len(srcRows); r++ {
		for c := 1; c <= cols; c++ {

			col, _ := excelize.ColumnNumberToName(c)

			srcCell := fmt.Sprintf("%s%d", col, r)
			dstCell := fmt.Sprintf("%s%d", col, startRow+(r-1))

			// value
			val, _ := sourceFile.GetCellValue(templateSheet, srcCell)
			if err := destFile.SetCellValue(sheetDest, dstCell, val); err != nil {
				return fmt.Errorf("set value: %w", err)
			}

			// ---------------------------
			// ✅ FIX: correct sheet + portable style
			// ---------------------------
			styleID, err := sourceFile.GetCellStyle(templateSheet, srcCell)
			if err == nil && styleID != 0 {

				var newStyleID int

				// reuse style if already created
				if cached, ok := styleCache[styleID]; ok {
					newStyleID = cached
				} else {
					styleDef, err := sourceFile.GetStyle(styleID)
					if err == nil {
						newStyleID, err = destFile.NewStyle(styleDef)
						if err == nil {
							styleCache[styleID] = newStyleID
						}
					}
				}

				if newStyleID != 0 {
					_ = destFile.SetCellStyle(sheetDest, dstCell, dstCell, newStyleID)
				}
			}
		}
	}

	// ---------------------------
	// 5. Copy merged cells (FIXED sheet + offset)
	// ---------------------------
	merges, err := sourceFile.GetMergeCells(templateSheet)
	if err == nil {
		for _, m := range merges {

			start := m.GetStartAxis()
			end := m.GetEndAxis()

			sCol, sRow, _ := excelize.CellNameToCoordinates(start)
			eCol, eRow, _ := excelize.CellNameToCoordinates(end)

			startColName, _ := excelize.ColumnNumberToName(sCol)
			endColName, _ := excelize.ColumnNumberToName(eCol)

			newStart := fmt.Sprintf("%s%d", startColName, startRow+(sRow-1))
			newEnd := fmt.Sprintf("%s%d", endColName, startRow+(eRow-1))

			_ = destFile.MergeCell(sheetDest, newStart, newEnd)
		}
	}

	return nil
}

func CloneExcelFile(src *excelize.File) (*excelize.File, error) {
	buf, err := src.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return excelize.OpenReader(buf)
}

func colName(n int) string {
	col, _ := excelize.ColumnNumberToName(n)
	return col
}

func ensureSheetExists(f *excelize.File, sheet string) error {
	idx, err := f.GetSheetIndex(sheet)
	if err != nil {
		return err
	}

	if idx == -1 {
		f.NewSheet(sheet)
	}

	return nil
}
