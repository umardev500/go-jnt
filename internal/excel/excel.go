package excel

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

type RowData struct {
	RowIndex int
	Values   map[string]string
}

type ExcelFile struct {
	File      *excelize.File
	SheetName string
	HeaderMap map[string]int // "KODE JMS" -> column index
}

// Open and map headers
func Open(filePath string) *ExcelFile {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	sheet := f.GetSheetList()[0]

	rows, err := f.GetRows(sheet)
	if err != nil {
		log.Fatal(err)
	}

	headerMap := make(map[string]int)
	if len(rows) == 0 {
		log.Println("⚠️ Sheet has no rows:", sheet)
		return &ExcelFile{
			File:      f,
			SheetName: sheet,
			HeaderMap: headerMap,
		}
	}

	for i, col := range rows[0] {
		headerMap[col] = i
	}

	return &ExcelFile{
		File:      f,
		SheetName: sheet,
		HeaderMap: headerMap,
	}
}

// Generic row reader
func (e *ExcelFile) GetRows() []RowData {
	rows, err := e.File.GetRows(e.SheetName)
	if err != nil {
		log.Fatal(err)
	}

	var result []RowData

	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}

		values := make(map[string]string)

		for header, idx := range e.HeaderMap {
			if idx < len(row) {
				values[header] = row[idx]
			} else {
				values[header] = ""
			}
		}

		result = append(result, RowData{
			RowIndex: i + 1,
			Values:   values,
		})
	}

	return result
}

func (e *ExcelFile) SetValue(row int, header string, value interface{}) {
	colIdx, ok := e.HeaderMap[header]
	if !ok {
		log.Fatalf("Column %s not found", header)
	}

	colName, err := excelize.ColumnNumberToName(colIdx + 1)
	if err != nil {
		log.Fatal(err)
	}

	cell := colName + fmt.Sprint(row)

	e.File.SetCellValue(e.SheetName, cell, value)
}

// Save file
func (e *ExcelFile) Save(path string) {
	if err := e.File.SaveAs(path); err != nil {
		log.Fatal(err)
	}
}

func (e *ExcelFile) SetFontColor(row int, header string, color string) {
	colIdx, ok := e.HeaderMap[header]
	if !ok {
		log.Fatalf("Column %s not found", header)
	}

	colName, err := excelize.ColumnNumberToName(colIdx + 1)
	if err != nil {
		log.Fatal(err)
	}

	cell := colName + fmt.Sprint(row)

	styleID, err := e.File.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color: color,
			Bold:  true,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	err = e.File.SetCellStyle(e.SheetName, cell, cell, styleID)
	if err != nil {
		log.Fatal(err)
	}
}
