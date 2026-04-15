package pivot

type PivotRow struct {
	Bagging string
	Count   int
}

// countPivotRows calculates total count and extracts blank bagging row (if any)
func CountPivotRows(rows []PivotRow) (total int, blank *PivotRow) {
	for _, r := range rows {
		total += r.Count

		if r.Bagging == "" {
			// keep the first blank bagging row
			if blank == nil {
				tmp := r
				blank = &tmp
			}
			continue
		}
	}

	return total, blank
}

// BuildPivot reads and transforms pivot Excel data
func BuildPivot(rows [][]string) map[string]int {

	result := make(map[string]int)

	for i, row := range rows {

		if i == 0 {
			continue // skip header
		}

		waybill := ""
		bagging := ""

		if len(row) > 0 {
			waybill = row[0]
		}
		if len(row) > 1 {
			bagging = row[1]
		}

		// only count valid waybill
		if waybill == "" {
			continue
		}

		result[bagging]++
	}

	return result
}

// ToPivotRows converts map[string]int to []PivotRow
func ToPivotRows(pivot map[string]int) []PivotRow {

	var result []PivotRow

	for bagging, count := range pivot {

		result = append(result, PivotRow{
			Bagging: bagging,
			Count:   count,
		})
	}

	return result
}
