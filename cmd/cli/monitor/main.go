package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/xuri/excelize/v2"
)

func removeUnusedColumns(f *excelize.File, sheet string, columns []string, log zerolog.Logger) error {
	// Remove columns from right to left
	for _, col := range columns {
		log.Info().
			Str("column", col).
			Msg("removing column")

		if err := f.RemoveCol(sheet, col); err != nil {
			return err
		}
	}

	return nil
}

func insertStyledColumnsAfterC(f *excelize.File, sheet string, log zerolog.Logger) error {
	// Insert two columns before D
	if err := f.InsertCols(sheet, "D", 2); err != nil {
		return err
	}

	// Copy style from original D (now shifted to F) to new D and E
	rows, err := f.GetRows(sheet)
	if err != nil {
		return err
	}

	for i := 1; i <= len(rows); i++ {
		sourceCell := fmt.Sprintf("F%d", i) // original D after insert

		targetCells := []string{
			fmt.Sprintf("D%d", i),
			fmt.Sprintf("E%d", i),
		}

		styleID, err := f.GetCellStyle(sheet, sourceCell)
		if err != nil {
			return err
		}

		for _, target := range targetCells {
			if err := f.SetCellStyle(sheet, target, target, styleID); err != nil {
				return err
			}
		}
	}

	// Add headers for the new columns
	if err := f.SetCellValue(sheet, "D1", "KLASIFIKASI"); err != nil {
		return err
	}

	if err := f.SetCellValue(sheet, "E1", "PROVINSI"); err != nil {
		return err
	}

	log.Info().
		Str("sheet", sheet).
		Msg("inserted styled columns with headers")

	return nil
}

func main() {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	fileName := "monitor.xlsx"

	f, err := excelize.OpenFile(fileName)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open Excel file")
	}
	defer f.Close()

	sheet := "sheet1"

	unusedColumns := []string{"M", "L", "J", "I", "G", "E", "D"}

	if err := removeUnusedColumns(f, sheet, unusedColumns, log); err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to remove unused columns")
	}

	if err := insertStyledColumnsAfterC(f, "sheet1", log); err != nil {
		log.Fatal().Err(err).Msg("column insert failed")
	}

	if err := f.SaveAs("monitor_modified.xlsx"); err != nil {
		log.Fatal().Err(err).Msg("failed to save Excel file")
	}

	log.Info().Msg("columns removed successfully")
}
