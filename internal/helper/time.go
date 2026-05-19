package helper

import (
	"fmt"
	"strings"
	"time"
)

func ParseExcelDateTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)

	// toggle this if you want debug logs
	debug := false

	layouts := []string{
		// 🔥 MOST COMMON FIRST (fix your error case)
		"1/2/06 15:04",   // 4/20/26 08:00
		"1/2/2006 15:04", // 4/20/2026 08:00
		"01/02/2006 15:04",

		// without time (US)
		"1/2/06",
		"01/02/06",
		"1/2/2006",
		"01/02/2006",

		// 🌍 International (DD/MM)
		"2/1/2006 15:04",
		"02/01/2006 15:04",
		"2/1/2006",
		"02/01/2006",

		// dash variants
		"1-2-06 15:04",
		"01-02-06 15:04",
		"1-2-2006 15:04",
		"01-02-2006 15:04",

		"1-2-06",
		"01-02-06",
		"1-2-2006",
		"01-02-2006",

		// ISO (safe fallback)
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		} else if debug {
			fmt.Println("failed layout:", layout, "error:", err)
		}
	}

	return time.Time{}, fmt.Errorf("unsupported date format: %s", s)
}
