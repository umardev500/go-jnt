package monitor

func ExcelColumnName(n int) string {
	name := ""

	for n > 0 {
		n--
		name = string(rune('A'+n%26)) + name
		n /= 26
	}

	return name
}
