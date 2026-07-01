package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func toFloat(s string) float64 {
	s = strings.TrimSpace(strings.ReplaceAll(s, ",", ""))
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func writeShiftedRow(dst *excelize.File, sheet string, row []string, outRow int) {

	// =========================
	// A - L (0 - 11)
	// =========================
	for col := 0; col < 12 && col < len(row); col++ {
		cell, _ := excelize.CoordinatesToCellName(col+1, outRow)
		dst.SetCellValue(sheet, cell, row[col])
	}

	// =========================
	// TOTAL (M)
	// =========================
	k := ""
	l := ""

	if len(row) > 10 {
		k = row[10]
	}
	if len(row) > 11 {
		l = row[11]
	}

	total := toFloat(k) + toFloat(l)
	dst.SetCellValue(sheet, fmt.Sprintf("M%d", outRow), total)

	// =========================
	// SHIFT REST (old M → N etc)
	// =========================
	for col := 12; col < len(row); col++ {
		cell, _ := excelize.CoordinatesToCellName(col+2, outRow)
		dst.SetCellValue(sheet, cell, row[col])
	}
}

func writeShiftedHeader(dst *excelize.File, sheet string, header []string) {

	// A - L
	for col := 0; col < 12 && col < len(header); col++ {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		dst.SetCellValue(sheet, cell, header[col])
	}

	// M = TOTAL
	dst.SetCellValue(sheet, "M1", "TOTAL")

	// shift rest header (M → N etc)
	for col := 12; col < len(header); col++ {
		cell, _ := excelize.CoordinatesToCellName(col+2, 1)
		dst.SetCellValue(sheet, cell, header[col])
	}
}

func main() {

	// =========================
	// LOAD CODES (ORDERED)
	// =========================
	data, err := os.ReadFile("codes.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(data), "\n")

	var codes []string
	for _, line := range lines {
		line = strings.TrimSpace(strings.ReplaceAll(line, "\r", ""))
		if line != "" {
			codes = append(codes, line)
		}
	}

	// =========================
	// OPEN EXCEL
	// =========================
	src, err := excelize.OpenFile("SURAT JALAN.xlsm")
	if err != nil {
		log.Fatal(err)
	}
	defer src.Close()

	sheet := "HistoryInput"

	rows, err := src.GetRows(sheet)
	if err != nil {
		log.Fatal(err)
	}

	// =========================
	// INDEX BY COLUMN Q
	// =========================
	index := make(map[string][]string)

	for i := 2; i <= len(rows); i++ {
		code, _ := src.GetCellValue(sheet, fmt.Sprintf("Q%d", i))
		code = strings.TrimSpace(code)

		if code == "" {
			continue
		}

		index[code] = rows[i-1]
	}

	// =========================
	// CREATE OUTPUT
	// =========================
	dst := excelize.NewFile()
	dstSheet := dst.GetSheetName(0)

	outRow := 1

	// =========================
	// SHIFTED HEADER
	// =========================
	writeShiftedHeader(dst, dstSheet, rows[0])

	outRow++

	// =========================
	// SHIFTED DATA (ORDERED)
	// =========================
	for i, code := range codes {

		row, ok := index[code]

		if !ok {
			log.Printf("[NOT FOUND] %s", code)
			continue
		}

		log.Printf("[MATCH] %d -> %s", i+1, code)

		writeShiftedRow(dst, dstSheet, row, outRow)
		outRow++
	}

	// =========================
	// SAVE
	// =========================
	if err := dst.SaveAs("filtered_surat_jalan.xlsx"); err != nil {
		log.Fatal(err)
	}

	log.Println("DONE")
}
