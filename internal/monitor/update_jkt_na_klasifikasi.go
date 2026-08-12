package monitor

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/xuri/excelize/v2"
)

func UpdateJKTIfKlasifikasiNA(
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

	// Find columns by header
	for i, header := range rows[0] {
		switch strings.TrimSpace(header) {
		case "Tujuan":
			tujuanCol = i
		case "KLASIFIKASI":
			klasifikasiCol = i
		}
	}

	if tujuanCol == -1 {
		return fmt.Errorf("Tujuan column not found")
	}

	if klasifikasiCol == -1 {
		return fmt.Errorf("KLASIFIKASI column not found")
	}

	// Process rows
	for row := 2; row <= len(rows); row++ {

		klasifikasi := ""

		if klasifikasiCol < len(rows[row-1]) {
			klasifikasi = strings.TrimSpace(rows[row-1][klasifikasiCol])
		}

		// Only update when KLASIFIKASI == "#N/A"
		if klasifikasi != "#N/A" {
			continue
		}

		if err := f.SetCellValue(
			sheet,
			fmt.Sprintf("%s%d", ExcelColumnName(tujuanCol+1), row),
			"JKT",
		); err != nil {
			return err
		}

		log.Info().
			Int("row", row).
			Msg("updated tujuan to JKT because klasifikasi is #N/A")
	}

	return nil
}
