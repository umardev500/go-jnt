package main

import (
	"fmt"

	"github.com/umardev500/jnt-report/internal/config"
	"github.com/umardev500/jnt-report/internal/info"
)

func main() {
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		panic(err)
	}

	plates := config.ExtractPlates(cfg)

	client := info.NewClient(cfg.Token)

	apiData, err := client.GetPlateInfo(plates)
	if err != nil {
		panic(err)
	}

	compare := info.CompareVehicles(cfg.Vendors, apiData)

	for _, r := range compare {
		status := "OK"
		if !r.Match {
			status = "FAIL: " + r.Reason
		}

		fmt.Println(r.API.Pretty(status))
		fmt.Println()
	}
}
