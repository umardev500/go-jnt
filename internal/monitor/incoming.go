package monitor

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/xuri/excelize/v2"
)

func UpdateIncomingByLokasiSebelumnya(
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

	lokasiCol := -1
	tujuanCol := -1
	klasifikasiCol := -1
	provinsiCol := -1

	// Find columns
	for i, header := range rows[0] {

		switch strings.TrimSpace(header) {

		case "Lokasi Sebelumnya":
			lokasiCol = i

		case "Tujuan":
			tujuanCol = i

		case "KLASIFIKASI":
			klasifikasiCol = i

		case "PROVINSI":
			provinsiCol = i
		}
	}

	if lokasiCol == -1 {
		return fmt.Errorf("Lokasi Sebelumnya column not found")
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

		lokasi := ""

		if lokasiCol < len(rows[row-1]) {
			lokasi = strings.ToUpper(
				strings.TrimSpace(
					rows[row-1][lokasiCol],
				),
			)
		}

		// blank -> do nothing
		if lokasi == "" {
			continue
		}

		// Check gateway / transit
		if !strings.Contains(lokasi, "GATEWAY") &&
			!strings.Contains(lokasi, "JKS_TRANSIT_CENTER") {
			continue
		}

		// Update C, D, E
		if err := f.SetCellValue(
			sheet,
			fmt.Sprintf("%s%d", ExcelColumnName(tujuanCol+1), row),
			"INCOMING",
		); err != nil {
			return err
		}

		if err := f.SetCellValue(
			sheet,
			fmt.Sprintf("%s%d", ExcelColumnName(klasifikasiCol+1), row),
			"INCOMING",
		); err != nil {
			return err
		}

		if err := f.SetCellValue(
			sheet,
			fmt.Sprintf("%s%d", ExcelColumnName(provinsiCol+1), row),
			"INCOMING",
		); err != nil {
			return err
		}

		log.Info().
			Int("row", row).
			Str("lokasi_sebelumnya", lokasi).
			Msg("updated incoming")
	}

	return nil
}
