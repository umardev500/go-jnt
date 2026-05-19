package info

import (
	"fmt"

	"github.com/umardev500/jnt-report/internal/config"
)

type CompareResult struct {
	Vendor string
	Plate  string
	Match  bool
	Reason string
	API    Vehicle
	YAML   config.Vehicle
}

// ================= MAIN COMPARE =================

func CompareVehicles(yamlData map[string][]config.Vehicle, apiData []Vehicle) []CompareResult {
	apiMap := make(map[string]Vehicle)

	for _, v := range apiData {
		apiMap[v.PlateNumber] = v
	}

	var result []CompareResult

	for vendor, vehicles := range yamlData {
		for _, y := range vehicles {

			api, ok := apiMap[y.Plate]

			if !ok {
				result = append(result, CompareResult{
					Vendor: vendor,
					Plate:  y.Plate,
					Match:  false,
					Reason: "NOT FOUND IN API",
					YAML:   y,
				})
				continue
			}

			match := api.CarrierName == vendor && api.VehicleTypeName == y.Type

			reason := "OK"
			if !match {
				reason = fmt.Sprintf(
					"MISMATCH (API vendor=%s type=%s)",
					api.CarrierName,
					api.VehicleTypeName,
				)
			}

			result = append(result, CompareResult{
				Vendor: vendor,
				Plate:  y.Plate,
				Match:  match,
				Reason: reason,
				API:    api,
				YAML:   y,
			})
		}
	}

	return result
}
