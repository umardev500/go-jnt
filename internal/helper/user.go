package helper

import (
	"fmt"

	"github.com/umardev500/jnt-report/internal/models"
)

var users = map[string]models.User{
	"01830317": {
		StaffNo: "01830317",
		Name:    "RAHMAWAN",
		Phone:   "6281389741086",
	},
	"01829395": {
		StaffNo: "01829395",
		Name:    "UMAR",
		Phone:   "6283142765573",
	},
	"01753499": {
		StaffNo: "01753499",
		Name:    "IMAM",
		Phone:   "6285199679372",
	},
}

func GetUser(staffNo string) (models.User, error) {
	user, ok := users[staffNo]
	if !ok {
		return models.User{}, fmt.Errorf(
			"user not found for staff no: %s", staffNo,
		)
	}

	return user, nil
}
