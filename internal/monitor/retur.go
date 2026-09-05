package monitor

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/xuri/excelize/v2"
)

func UpdateTujuanFromRetur(
	f *excelize.File,
	sheet string,
	log zerolog.Logger,
) error {

	rows, err := f.GetRows(sheet)

	if err != nil {
		return err
	}

	statusCol := -1

	for i, h := range rows[0] {
		h = normalizeHeader(h)
		h = strings.ToUpper(h)
		fmt.Println(h)

		if h == "STATUS RETUR" {
			statusCol = i
			break
		}
	}

	if statusCol == -1 {
		return fmt.Errorf("status retur missing")
	}

	for i := 1; i < len(rows); i++ {

		if statusCol >= len(rows[i]) {
			continue
		}

		if strings.TrimSpace(rows[i][statusCol]) == "Y" {

			row := i + 1

			value, err := f.GetCellValue(
				sheet,
				fmt.Sprintf("B%d", row),
			)

			if err != nil {
				return err
			}

			f.SetCellValue(
				sheet,
				fmt.Sprintf("C%d", row),
				value,
			)
		}
	}

	log.Info().
		Msg("updated retur")

	return nil
}
