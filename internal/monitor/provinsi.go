package monitor

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/xuri/excelize/v2"
)

func FillProvinsi(
	f *excelize.File,
	sheet string,
	file string,
	lookupSheet string,
	log zerolog.Logger,
) error {

	lookup, err := BuildLookup(
		file,
		lookupSheet,
		"SCAN KIRIM TUJUAN",
		"KLASIFIKASI PROVINSI",
	)

	if err != nil {
		return err
	}

	rows, err := f.GetRows(sheet)

	if err != nil {
		return err
	}

	for row := 2; row <= len(rows); row++ {

		klasifikasi, err := f.GetCellValue(
			sheet,
			fmt.Sprintf("D%d", row),
		)

		if err != nil {
			return err
		}

		value, ok := lookup[strings.ToUpper(
			strings.TrimSpace(klasifikasi),
		)]

		if !ok {
			value = "#N/A"
		}

		f.SetCellValue(
			sheet,
			fmt.Sprintf("E%d", row),
			value,
		)
	}

	log.Info().
		Msg("filled provinsi")

	return nil
}
