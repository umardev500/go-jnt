package helper

import (
	"fmt"
	"strings"
	"time"
)

func ParseExcelDateTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)

	layouts := []string{
		"1-2-06 15:04",
		"01-02-06 15:04",
		"1/2/06 15:04",
		"01/02/06 15:04",

		"1-2-06",
		"01-02-06",
		"1/2/06",
		"01/02/06",

		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported date format: %s", s)
}
