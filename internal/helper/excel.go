package helper

import (
	"github.com/xuri/excelize/v2"
)

func CreateBorderStyle(f *excelize.File) (int, error) {
	return f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})
}
