package monitor

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/xuri/excelize/v2"
)

func normalizeHeader(s string) string {
	s = strings.ReplaceAll(s, "\u00A0", " ") // non-breaking space
	s = strings.ReplaceAll(s, "\u200B", "")  // zero-width space
	s = strings.ReplaceAll(s, "\uFEFF", "")  // BOM
	return strings.TrimSpace(s)
}

func UpdateKlasifikasiByJenisLayanan(
	f *excelize.File,
	sheet string,
	log zerolog.Logger,
) error {

	rows, err := f.GetRows(sheet)
	if err != nil {
		return err
	}

	if len(rows) <= 1 {
		return fmt.Errorf("no data rows found")
	}

	jenisLayananCol := -1
	klasifikasiCol := -1

	// Find columns
	for i, header := range rows[0] {

		header = normalizeHeader(header)
		header = strings.ToUpper(header)
		fmt.Println(header)

		switch header {
		case "JENIS LAYANAN":
			jenisLayananCol = i

		case "KLASIFIKASI":
			klasifikasiCol = i
		}
	}

	if jenisLayananCol == -1 {
		return fmt.Errorf("Jenis Layanan column not found")
	}

	if klasifikasiCol == -1 {
		return fmt.Errorf("KLASIFIKASI column not found")
	}

	// Process rows
	for row := 2; row <= len(rows); row++ {

		jenisLayanan := ""

		if jenisLayananCol < len(rows[row-1]) {

			jenisLayanan = strings.ToUpper(
				strings.TrimSpace(
					rows[row-1][jenisLayananCol],
				),
			)
		}

		klasifikasi, err := f.GetCellValue(
			sheet,
			fmt.Sprintf(
				"%s%d",
				ExcelColumnName(klasifikasiCol+1),
				row,
			),
		)

		if err != nil {
			return err
		}

		klasifikasi = strings.ToUpper(
			strings.TrimSpace(
				klasifikasi,
			),
		)

		// Ignore blank values
		if klasifikasi == "" || jenisLayanan == "" {
			continue
		}

		// Only modify JKT_GATEWAY
		if klasifikasi != "JKT_GATEWAY" {
			continue
		}

		newValue := ""

		if jenisLayanan == "ECO" {

			newValue = "JKT_GATEWAY ECO"

		} else {

			newValue = "JKT_GATEWAY AF"
		}

		err = f.SetCellValue(
			sheet,
			fmt.Sprintf(
				"%s%d",
				ExcelColumnName(klasifikasiCol+1),
				row,
			),
			newValue,
		)

		if err != nil {
			return err
		}

		log.Info().
			Int("row", row).
			Str("jenis_layanan", jenisLayanan).
			Str("old_klasifikasi", klasifikasi).
			Str("new_klasifikasi", newValue).
			Msg("updated klasifikasi by jenis layanan")
	}

	return nil
}
