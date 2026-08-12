package monitor

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/xuri/excelize/v2"
)

func FillKlasifikasi(
	f *excelize.File,
	sheet string,
	file string,
	lookupSheet string,
	log zerolog.Logger,
) error {

	lookup, err := BuildLookup(
		file,
		lookupSheet,
		"DP TUJUAN",
		"KLASIFIKASI SCAN KIRIM",
	)

	if err != nil {
		return err
	}

	rows, err := f.GetRows(sheet)

	if err != nil {
		return err
	}

	tujuanCol := 2 // column C

	for row := 2; row <= len(rows); row++ {

		tujuan := strings.ToUpper(
			strings.TrimSpace(
				rows[row-1][tujuanCol],
			),
		)

		value, ok := lookup[tujuan]

		if !ok {
			value = "#N/A"
		}

		f.SetCellValue(
			sheet,
			fmt.Sprintf("D%d", row),
			value,
		)
	}

	log.Info().
		Msg("filled klasifikasi")

	return nil
}
