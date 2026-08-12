package monitor

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/xuri/excelize/v2"
)

func UpdateBelumPickup(
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

	tujuanCol := -1
	klasifikasiCol := -1
	provinsiCol := -1

	// Find columns
	for i, header := range rows[0] {

		switch strings.TrimSpace(header) {

		case "Tujuan":
			tujuanCol = i

		case "KLASIFIKASI":
			klasifikasiCol = i

		case "PROVINSI":
			provinsiCol = i
		}
	}

	if tujuanCol == -1 {
		return fmt.Errorf("Tujuan column not found")
	}

	if klasifikasiCol == -1 {
		return fmt.Errorf("KLASIFIKASI column not found")
	}

	if provinsiCol == -1 {
		return fmt.Errorf("PROVINSI column not found")
	}

	// Process rows
	for row := 2; row <= len(rows); row++ {

		tujuan := ""

		if tujuanCol < len(rows[row-1]) {
			tujuan = strings.TrimSpace(
				rows[row-1][tujuanCol],
			)
		}

		// Only process blank Tujuan
		if tujuan != "" {
			continue
		}

		// Update Tujuan C
		if err := f.SetCellValue(
			sheet,
			fmt.Sprintf("%s%d", ExcelColumnName(tujuanCol+1), row),
			"BLM PU",
		); err != nil {
			return err
		}

		// Update KLASIFIKASI D
		if err := f.SetCellValue(
			sheet,
			fmt.Sprintf("%s%d", ExcelColumnName(klasifikasiCol+1), row),
			"BELUM PICKUP",
		); err != nil {
			return err
		}

		// Update PROVINSI E
		if err := f.SetCellValue(
			sheet,
			fmt.Sprintf("%s%d", ExcelColumnName(provinsiCol+1), row),
			"BELUM PICKUP",
		); err != nil {
			return err
		}

		log.Info().
			Int("row", row).
			Msg("updated blank tujuan as belum pickup")
	}

	return nil
}
