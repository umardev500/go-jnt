package monitor

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/xuri/excelize/v2"
)

func Process(
	input string,
	output string,
	sheet string,
	lookupFile string,
	log zerolog.Logger,
) error {
	inputFile, err := excelize.OpenFile(input)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	lookupFileXlsx, err := excelize.OpenFile(lookupFile)
	if err != nil {
		return err
	}
	defer lookupFileXlsx.Close()

	// ----------------------------------------
	// Lookup: BTN KE TUJUAN
	// A = key
	// B = klasifikasi
	// ----------------------------------------
	btnLookup, err := buildLookup(
		lookupFileXlsx,
		"BTN KE TUJUAN",
	)
	if err != nil {
		return err
	}

	// ----------------------------------------
	// Lookup: PROVINSI ESTIMASI BARANG
	// A = klasifikasi
	// B = provinsi
	// ----------------------------------------
	provinsiLookup, err := buildLookup(
		lookupFileXlsx,
		"PROVINSI ESTIMASI BARANG",
	)
	if err != nil {
		return err
	}

	sheets := inputFile.GetSheetList()
	if len(sheets) == 0 {
		return fmt.Errorf("excel file has no sheets")
	}

	// ----------------------------------------
	// Read input
	// ----------------------------------------
	rows, err := inputFile.GetRows(sheets[0])
	if err != nil {
		return err
	}

	// ----------------------------------------
	// Create output
	// ----------------------------------------
	out := excelize.NewFile()
	defer out.Close()

	destSheet := out.GetSheetName(0)

	for rowIndex, row := range rows {
		destRow := rowIndex + 1

		// ------------------------------------
		// Header
		// ------------------------------------
		if rowIndex == 0 {
			// Original A -> Output A
			if len(row) > 0 {
				if err := out.SetCellValue(
					destSheet,
					fmt.Sprintf("A%d", destRow),
					row[0],
				); err != nil {
					return err
				}
			}

			// Original C -> Output B
			if len(row) > 2 {
				if err := out.SetCellValue(
					destSheet,
					fmt.Sprintf("B%d", destRow),
					row[2],
				); err != nil {
					return err
				}
			}

			// New columns
			if err := out.SetCellValue(
				destSheet,
				"C1",
				"klasifikas",
			); err != nil {
				return err
			}

			if err := out.SetCellValue(
				destSheet,
				"D1",
				"provinse",
			); err != nil {
				return err
			}

			// Original K -> Output E
			if len(row) > 10 {
				if err := out.SetCellValue(
					destSheet,
					"E1",
					row[10],
				); err != nil {
					return err
				}
			}

			continue
		}

		// ------------------------------------
		// Original A -> Output A
		// ------------------------------------
		if len(row) > 0 {
			if err := out.SetCellValue(
				destSheet,
				fmt.Sprintf("A%d", destRow),
				row[0],
			); err != nil {
				return err
			}
		}

		// ------------------------------------
		// Original C -> Output B
		//
		// This is the key for klasifikasi.
		// ------------------------------------
		klasifikasiKey := ""

		if len(row) > 2 {
			klasifikasiKey = strings.TrimSpace(row[2])

			if err := out.SetCellValue(
				destSheet,
				fmt.Sprintf("B%d", destRow),
				row[2],
			); err != nil {
				return err
			}
		}

		// ------------------------------------
		// Klasifikasi lookup
		//
		// Original C
		//   -> BTN KE TUJUAN A
		//   -> BTN KE TUJUAN B
		//   -> Output C
		// ------------------------------------
		klasifikasi := btnLookup[klasifikasiKey]

		// ------------------------------------
		// Replace klasifikasi
		//
		// MES_GATEWAY -> JKT_GATEWAY
		// ------------------------------------
		if klasifikasi == "MES_GATEWAY" {
			klasifikasi = "JKT_GATEWAY"
		}

		// Write klasifikasi
		if klasifikasi != "" {
			if err := out.SetCellValue(
				destSheet,
				fmt.Sprintf("C%d", destRow),
				klasifikasi,
			); err != nil {
				return err
			}
		}

		// ------------------------------------
		// Provinsi lookup
		//
		// Uses the FINAL klasifikasi value.
		//
		// Klasifikasi
		//   -> PROVINSI ESTIMASI BARANG A
		//   -> PROVINSI ESTIMASI BARANG B
		//   -> Output D
		// ------------------------------------
		provinsi := provinsiLookup[strings.TrimSpace(klasifikasi)]

		if provinsi != "" {
			if err := out.SetCellValue(
				destSheet,
				fmt.Sprintf("D%d", destRow),
				provinsi,
			); err != nil {
				return err
			}
		}

		// ------------------------------------
		// Original K -> Output E
		// EZ / NDD remain unchanged.
		// Everything else becomes REGULER.
		// ------------------------------------
		if len(row) > 10 {
			value := strings.TrimSpace(row[10])

			if value != "EZ" && value != "NDD" {
				value = "REGULER"
			}

			if err := out.SetCellValue(
				destSheet,
				fmt.Sprintf("E%d", destRow),
				value,
			); err != nil {
				return err
			}
		}

	}

	return out.SaveAs(output)
}

func buildLookup(
	f *excelize.File,
	sheet string,
) (map[string]string, error) {
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	lookup := make(map[string]string)

	// Start from row 2
	for i := 1; i < len(rows); i++ {
		if len(rows[i]) < 2 {
			continue
		}

		key := strings.TrimSpace(rows[i][0])
		value := strings.TrimSpace(rows[i][1])

		if key == "" {
			continue
		}

		lookup[key] = value
	}

	return lookup, nil
}
