package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/jnt-report/internal/detail"
	"github.com/umardev500/jnt-report/internal/excel"
)

var (
	greenBold = color.New(color.FgGreen, color.Bold)
	redBold   = color.New(color.FgRed, color.Bold)
	spaceRe   = regexp.MustCompile(`\s+`)
)

// ---------- GENERIC HELPERS ----------

func statusColor(matched bool) string {
	if matched {
		return greenBold.Sprint("Matched")
	}
	return redBold.Sprint("Not matched")
}

func normalize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "shuttle", "")
	s = strings.ReplaceAll(s, "pt", "")
	s = strings.TrimSpace(s)
	return spaceRe.ReplaceAllString(s, " ")
}

func compareStrings(a, b string) bool {
	aNorm := normalize(a)
	bNorm := normalize(b)
	return strings.Contains(bNorm, aNorm) || strings.Contains(aNorm, bNorm)
}

// ---------- CAR TYPE ----------

func matchCarType(s1, s2 string) bool {
	s1 = strings.TrimSpace(s1)
	s2 = strings.TrimSpace(s2)

	// Manual override
	if strings.EqualFold(s1, "CDDL") && strings.Contains(strings.ToUpper(s2), "CDD LONG") {
		return true
	}

	// fallback
	s1Norm := strings.ToLower(s1)
	s2Norm := strings.ToLower(s2)
	return strings.Contains(s2Norm, s1Norm) || strings.Contains(s1Norm, s2Norm)
}

// ---------- PLATE ----------

func normalizePlate(s string) string {
	s = strings.ToUpper(s)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "-", "")
	return strings.TrimSpace(s)
}

func matchPlate(p1, p2 string) bool {
	return normalizePlate(p1) == normalizePlate(p2)
}

// ---------- MAIN ----------

func ValidateSuratJalan(ex *excel.ExcelFile, token string) {
	rows := ex.GetRows()

	// unified formats
	labelFmt := "%-14s: %s\n"
	rowFmt := "%-14s: %-30s | %-40s | %s\n"
	separator := strings.Repeat("-", 100)

	for _, r := range rows {
		kode := strings.TrimSpace(r.Values["Kode Tugas"])
		if kode == "" {
			log.Info().Msg("Skipping empty KODE JMS")
			continue
		}

		dt, err := detail.GetShipmentDetail(kode, token)
		if err != nil {
			log.Error().Err(err).Str("Kode Tugas", kode).Msg("Error fetching shipment detail")
			continue
		}

		// --- DATA ---
		shipment := dt.Data.ShipmentDetail.CarrierName
		vendor := r.Values["Vendor"]
		car1 := dt.Data.ShipmentDetail.VehicleTypeName
		car2 := r.Values["Jenis Mobil"]
		plate1 := dt.Data.ShipmentDetail.PlateNumber
		plate2 := r.Values["Plat Nomor"]

		// --- HEADER ---
		fmt.Printf(labelFmt, "Kode Tugas", kode)

		// --- VALIDATIONS (ALL SAME STRUCTURE) ---
		fmt.Printf(rowFmt, "Vendor",
			shipment,
			vendor,
			statusColor(compareStrings(shipment, vendor)),
		)

		fmt.Printf(rowFmt, "Car Type",
			car1,
			car2,
			statusColor(matchCarType(car1, car2)),
		)

		fmt.Printf(rowFmt, "Plate",
			plate1,
			plate2,
			statusColor(matchPlate(plate1, plate2)),
		)

		fmt.Println(separator)
	}
}
