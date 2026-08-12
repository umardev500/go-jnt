package monitor

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/xuri/excelize/v2"
)

func RemoveUnusedColumns(
	f *excelize.File,
	sheet string,
	columns []string,
	log zerolog.Logger,
) error {

	for _, col := range columns {

		if err := f.RemoveCol(
			sheet,
			col,
		); err != nil {
			return err
		}

		log.Info().
			Str("column", col).
			Msg("removed")
	}

	return nil
}

func InsertStyledColumnsAfterC(
	f *excelize.File,
	sheet string,
	log zerolog.Logger,
) error {

	if err := f.InsertCols(
		sheet,
		"D",
		2,
	); err != nil {
		return err
	}

	rows, err := f.GetRows(sheet)

	if err != nil {
		return err
	}

	for i := 1; i <= len(rows); i++ {

		style, err := f.GetCellStyle(
			sheet,
			fmt.Sprintf("F%d", i),
		)

		if err != nil {
			return err
		}

		for _, col := range []string{"D", "E"} {

			cell := fmt.Sprintf(
				"%s%d",
				col,
				i,
			)

			err = f.SetCellStyle(
				sheet,
				cell,
				cell,
				style,
			)

			if err != nil {
				return err
			}
		}
	}

	f.SetCellValue(
		sheet,
		"D1",
		"KLASIFIKASI",
	)

	f.SetCellValue(
		sheet,
		"E1",
		"PROVINSI",
	)

	log.Info().
		Msg("created columns")

	return nil
}
