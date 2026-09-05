package suratjalan

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/umardev500/jnt-report/internal/config"
	"github.com/umardev500/jnt-report/internal/data"
	"github.com/umardev500/jnt-report/internal/detail"
	"github.com/umardev500/jnt-report/internal/downloader"
	"github.com/umardev500/jnt-report/internal/excel"
	"github.com/umardev500/jnt-report/internal/reporting"
	"github.com/umardev500/jnt-report/internal/route"
	"github.com/umardev500/jnt-report/internal/types"
	"github.com/xuri/excelize/v2"
	"golang.org/x/term"
)

type Vehicle struct {
	Vendor string
	Plat   string
	Type   string
}

const (
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

func prettyJSONLog(level string, fields map[string]any, message string) {
	data := map[string]any{
		"level":   level,
		"message": message,
	}

	for key, value := range fields {
		data[key] = value
	}

	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("failed to marshal pretty JSON: %v\n", err)
		return
	}

	var pretty map[string]any
	if err := json.Unmarshal(output, &pretty); err != nil {
		fmt.Printf("failed to parse JSON: %v\n", err)
		return
	}

	fmt.Println("{")

	i := 0
	for key, value := range pretty {
		comma := ","
		if i == len(pretty)-1 {
			comma = ""
		}

		valueJSON, _ := json.MarshalIndent(value, "  ", "  ")

		fmt.Printf(
			"  %q: %s%s%s%s\n",
			key,
			colorGreen,
			valueJSON,
			colorReset,
			comma,
		)

		i++
	}

	fmt.Println("}")
}

func logSection(title string) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width < 40 {
		width = 80
	}

	line := strings.Repeat("━", width)

	fmt.Println(line)
	fmt.Println(title)
	fmt.Println(line)
}

func confirmContinue(prompt string) bool {
	var answer string

	fmt.Print(prompt)
	fmt.Scanln(&answer)

	return strings.EqualFold(strings.TrimSpace(answer), "Y")
}

func FindVehicleByPlat(filename, sheetName, plat string) (*Vehicle, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	for i, row := range rows {
		// Skip header
		if i == 0 {
			continue
		}

		// Need at least columns A, B, C
		if len(row) < 3 {
			continue
		}

		if strings.EqualFold(strings.TrimSpace(row[1]), strings.TrimSpace(plat)) {
			return &Vehicle{
				Vendor: row[0], // Column A
				Plat:   row[1], // Column B
				Type:   row[2], // Column C
			}, nil
		}
	}

	return nil, fmt.Errorf("plat %q not found", plat)
}

func GenateSuratJalan(ex *excel.ExcelFile, routes types.GatewayRoutes, code, token, admin, service string, prod bool) (*detail.ShipmentDetailResponse, *reporting.ReportSummary) {
	logSection("🚚 STARTING SURAT JALAN GENERATION")

	dt, err := detail.GetShipmentDetail(code, token)
	if err != nil {
		log.Error().Err(err).Str("Kode Tugas", code).Msg("Error fetching shipment detail")
	}
	shipment := dt.Data.ShipmentDetail
	dest := shipment.TmsShipmentStopVOList[1]
	origin := shipment.TmsShipmentStopVOList[0]

	plateNumber := dt.Data.ShipmentDetail.PlateNumber
	vtn := dt.Data.ShipmentDetail.VehicleTypeName

	log.Info().Str("Plat", plateNumber).Str("Type", dt.Data.ShipmentDetail.VehicleTypeName).Msg("find vehicle by plat")
	vehicle, err := FindVehicleByPlat(config.GetVehicleFilePath(prod), "VEHICLES", plateNumber)
	if err != nil {
		log.Error().Err(err).Msg("failed to find vehicle by plat, continue...")
	}

	// ============================================================
	// VENDOR VALIDATION
	// ============================================================

	log.Warn().Msg("⚠️ Vendor confirmation required")

	prettyJSONLog(
		"warn",
		map[string]any{
			"Plat":   plateNumber,
			"Vendor": shipment.CarrierName,
		},
		"Vendor confirmation required",
	)

	if !confirmContinue("Continue anyway? (Y/N): ") {
		log.Warn().
			Str("Plat", vehicle.Plat).
			Str("Vendor", vehicle.Vendor).
			Msg("❌ Skipping...")

		return nil, nil
	}

	log.Info().
		Str("Plat", vtn).
		Str("Vendor", shipment.CarrierName).
		Msg("▶️ Continuing...")

	// ============================================================
	// VEHICLE TYPE VALIDATION
	// ============================================================

	if vehicle != nil && vehicle.Type == shipment.VehicleTypeName {
		log.Info().Msg("✅ Vehicle type match")

		prettyJSONLog(
			"info",
			map[string]any{
				"Carrier":      shipment.CarrierName,
				"Plat":         plateNumber,
				"ExpectedType": vtn,
				"ActualType":   vehicle.Type,
			},
			"Vehicle type match",
		)
	} else {
		log.Warn().Msg("⚠️ Vehicle type mismatch")

		prettyJSONLog(
			"warn",
			map[string]any{
				"Carrier":      shipment.CarrierName,
				"Plat":         plateNumber,
				"ExpectedType": vtn,
				"ActualType":   "-",
			},
			"Vehicle type mismatch",
		)

		if !confirmContinue("🚨 Vehicle type mismatch. Continue anyway? (Y/N): ") {
			log.Warn().
				Str("Plat", plateNumber).
				Str("ExpectedType", vtn).
				Str("ActualType", "-").
				Msg("🛑 Exiting...")

			return nil, nil
		}

		log.Info().
			Str("Plat", plateNumber).
			Str("ExpectedType", vtn).
			Str("ActualType", "-").
			Msg("▶️ Continuing...")
	}

	// ============================================================
	// APPLY VEHICLE TYPE
	// ============================================================

	if vehicle != nil {
		shipment.VehicleTypeName = vehicle.Type
		dt.Data.ShipmentDetail.VehicleTypeName = vehicle.Type

		log.Info().
			Str("Plat", vehicle.Plat).
			Str("VehicleType", vehicle.Type).
			Msg("🚚 Vehicle type applied")
	} else {
		log.Warn().
			Str("Plat", plateNumber).
			Msg("🚨 Vehicle plate not found in vehicle list")
	}

	now := time.Now()
	// today := time.Now().Format("2006-01-02")
	layout := "2006-01-02 15:04:05"
	// plannedDepartureTime, err := time.Parse(layout, shipment.PlannedDepartureTime)
	// if err != nil {
	// 	log.Error().Err(err).Msg("Error parsing planned departure time")
	// }

	var appTrackDepartureTime time.Time
	if origin.ActualDepartureTime != nil {
		appTrackDepartureTime, err = time.Parse(layout, *origin.ActualDepartureTime)
		if err != nil {
			log.Error().Err(err).Msg("Error parsing app track departure time")
		}
	}

	departureTime, err := time.Parse("2006-01-02 15:04:05", shipment.PlannedDepartureTime)
	if err != nil {
		fmt.Println("err sj")
		log.Error().Err(err).Msg("invalid date")
	}
	rs := route.GetRitase(routes, shipment.EndName, departureTime.Format("15:04"))
	jenisPaket := route.GetJenisPaket(rs.StatusRute)

	ex.SetValue(6, "B", shipment.PlannedDepartureTime)
	ex.SetValue(33, "C", admin)
	ex.SetValue(9, "E", appTrackDepartureTime.Format("02-01-2006"))
	ex.SetValue(10, "E", appTrackDepartureTime.Format("15:04:05")+" / "+jenisPaket)
	ex.SetValue(11, "E", shipment.PlateNumber)
	ex.SetValue(12, "E", strings.ToUpper(shipment.DriverName))
	ex.SetValue(13, "E", shipment.VehicleTypeName)
	ex.SetValue(14, "E", shipment.DriverContact)
	// ex.SetValue(15, "E", now.Format("15:04:05"))
	ex.SetValue(15, "E", now.Format("02/01/2006 15:04:05"))

	ex.SetValue(19, "E", normalizeCompanyName(shipment.CarrierName))

	ex.SetValue(17, "I", "Kode Tugas :"+code)
	ex.SetValue(18, "I", "Kode Tugas :"+code)

	ex.SetValue(16, "J", service)

	// Service
	ex.SetValue(3, "T", service)

	// PREPARE THE EXCEL FILE
	// Get pivot info
	// Step 1: Download shipment file
	log.Info().Msgf("Downloading shipment file... %s", code)
	if err := downloader.DownloadShipmentFile(code, token, prod); err != nil {
		fmt.Println("Error downloading shipment file:", err)
	}
	log.Info().Msg("Shipment file downloaded successfully!")

	// Step 2: Generate pivot report
	report, err := reporting.GeneratePivotReport(config.GetExportedFilePath(prod), "Memuat dan membongkar ekspor in", true)
	if err != nil {
		fmt.Println("Error generating report:", err)
		return nil, nil
	}

	routeCode := detail.GetRouteCode(dt)
	fmt.Println(report)
	fmt.Println(report.TotalWaybillCount)

	// INFO KOLI
	ex.SetValue(22, "G", strconv.Itoa(report.TotalPivotRows))
	ex.SetValue(26, "G", strconv.Itoa(report.TotalPivotRows))
	ex.SetValue(33, "G", strings.ToUpper(shipment.DriverName))
	ex.SetValue(21, "H", strconv.Itoa(report.BlankBaggingCount))
	ex.SetValue(22, "H", strconv.Itoa(report.TotalWaybillCount))
	ex.SetValue(26, "H", strconv.Itoa(report.TotalWaybillCount+report.BlankBaggingCount))

	ex.SetValue(10, "I", dest.NetworkName)
	address := data.Routes[routeCode]
	ex.SetValue(11, "I", address)

	return dt, report
}

func normalizeCompanyName(name string) string {
	if strings.HasPrefix(name, "PT ") {
		return strings.Replace(name, "PT ", "PT. ", 1)
	}
	return name
}
