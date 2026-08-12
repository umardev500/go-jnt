package monitor

import (
	"strings"

	"github.com/xuri/excelize/v2"
)

func BuildLookup(
	file string,
	sheet string,
	keyHeader string,
	valueHeader string,
) (map[string]string, error) {

	f, err := excelize.OpenFile(file)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	rows, err := f.GetRows(sheet)

	if err != nil {
		return nil, err
	}

	keyCol := -1
	valueCol := -1

	for i, h := range rows[0] {

		switch strings.TrimSpace(h) {

		case keyHeader:
			keyCol = i

		case valueHeader:
			valueCol = i
		}
	}

	result := map[string]string{}

	for _, row := range rows[1:] {

		if keyCol >= len(row) ||
			valueCol >= len(row) {
			continue
		}

		result[strings.ToUpper(
			strings.TrimSpace(row[keyCol]),
		)] = strings.TrimSpace(
			row[valueCol],
		)
	}

	return result, nil
}
