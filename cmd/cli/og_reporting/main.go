package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/jnt-report/internal/config"
	"github.com/umardev500/jnt-report/internal/detail"
	"github.com/umardev500/jnt-report/internal/downloader"
	"github.com/umardev500/jnt-report/internal/excel"
	"github.com/umardev500/jnt-report/internal/reporting"
	"github.com/umardev500/jnt-report/internal/route"
	"github.com/umardev500/jnt-report/internal/types"
)

// LoadGatewayRoutes reads a JSON file and unmarshals its content into GatewayRoutes.
func loadGatewayRoutes(filename string) (types.GatewayRoutes, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var data types.GatewayRoutes
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return data, nil
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

}

func main() {
	_ = godotenv.Load()

	prod := os.Getenv("APP_ENV") == "prod"
	log.Info().Msgf("Running in %v mode", prod)

	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		panic(err)
	}
	ex := excel.Open(config.GetReportFilePath(prod))

	routes, err := loadGatewayRoutes(config.GetSKOFilePath(prod))
	if err != nil {
		fmt.Println("Error loading gateway routes:", err)
		return
	}

	rows := ex.GetRows()

	for _, r := range rows {
		kode := r.Values["KODE JMS"]
		kode = strings.TrimSpace(kode)
		if kode == "" {
			log.Info().Msg("Skipping empty KODE JMS")
			continue
		}
		log.Info().Msgf("Processing %s", kode)

		// call API using kode...
		res, err := detail.GetShipmentDetail(kode, cfg.Token)
		if err != nil {
			fmt.Println("Error fetching shipment detail:", err)
			return
		}

		shipment := res.Data.ShipmentDetail

		var stop = new(detail.Stop)
		for _, r := range shipment.TmsShipmentStopVOList {
			if r.NetworkCode == "BTN777" {
				stop = &r
				break
			}
		}
		if stop == nil {
			fmt.Println("No record found for BTN777")
			return
		}

		// Destination
		dest := shipment.TmsShipmentStopVOList[1]
		sampaiDriverPtr := dest.ActualArrivalTime

		var sampaiDriver string
		if sampaiDriverPtr != nil {
			sampaiDriver = (*sampaiDriverPtr)[11:]
		}

		re := regexp.MustCompile(`^([A-Z]{3}\d{3}(?:-[A-Z]{3}\d{3})+)`)
		shipmentName := re.FindString(shipment.ShipmentName)

		// find route
		log.Info().Msgf("Finding route for %s : %s", res.Data.ShipmentDetail.EndName, *stop.PlannedDepartureHour)
		rs := route.GetRitase(routes, res.Data.ShipmentDetail.EndName, *stop.PlannedDepartureHour)
		fmt.Println(rs.StatusRute)

		layout := "2006-01-02 15:04:05"
		// departure, _ := time.Parse(layout, "2026-05-18 21:48:18")
		// arrival, _ := time.Parse(layout, "2026-05-19 00:59:02")
		departureStr := res.Data.ShipmentDetail.ActualDepartureTime
		arrivalStr := res.Data.ShipmentDetail.ActualArrivalTime

		departure, depErr := time.Parse(layout, departureStr)
		arrival, arrErr := time.Parse(layout, arrivalStr)

		fmt.Println("departure:", departure, "arrival:", arrival)

		hasValidTime := depErr == nil && arrErr == nil &&
			departureStr != "" && arrivalStr != "" &&
			!departure.IsZero() && !arrival.IsZero()

		ex.SetValue(r.RowIndex, "PLANNED DEPARTURE", res.Data.ShipmentDetail.PlannedDepartureTime)
		ex.SetValue(r.RowIndex, "SLA", res.Data.ShipmentDetail.TotalRuntime)

		if hasValidTime {
			duration := arrival.Sub(departure)
			minutes := int(math.Ceil(duration.Minutes()))
			selisih := minutes - res.Data.ShipmentDetail.TotalRuntime

			ex.SetValue(r.RowIndex, "AKTUAL", minutes)
			ex.SetValue(r.RowIndex, "SELISIH", selisih)
			if selisih > 0 {
				ex.SetFontColor(r.RowIndex, "SELISIH", "FF0000")
			}
		}
		ex.SetValue(r.RowIndex, "RUTE", shipmentName)
		ex.SetValue(r.RowIndex, "RITASE", rs.StatusRute)
		ex.SetValue(r.RowIndex, "NAMA", strings.ToUpper(shipment.DriverName))
		ex.SetValue(r.RowIndex, "NO HP", shipment.DriverContact)
		ex.SetValue(r.RowIndex, "VENDOR", shipment.CarrierName)
		ex.SetValue(r.RowIndex, "NOPOL", shipment.PlateNumber)
		if stop.AppDriverDeparture != nil {
			// ex.SetValue(r.RowIndex, "WAKTU KEBERANGKATAN APP DRIVER", (*stop.AppDriverDeparture)[11:])
			ex.SetValue(r.RowIndex, "WAKTU KEBERANGKATAN APP DRIVER", *stop.AppDriverDeparture)
		}
		ex.SetValue(r.RowIndex, "WAKTU SAMPAI APP DRIVER", sampaiDriver)
		ex.SetValue(r.RowIndex, "JENIS MOBIL", shipment.VehicleTypeName)

		// Set scan kirim mobil
		if stop.ScanTime != nil {
			// ex.SetValue(r.RowIndex, "WAKTU SCAN KIRIM MOBIL", (*stop.ScanTime)[11:])
			ex.SetValue(r.RowIndex, "WAKTU SCAN KIRIM MOBIL", *stop.ScanTime)
		}

		// Get pivot info
		// Step 1: Download shipment file
		skipDownload := false
		log.Info().Msgf("Downloading shipment file... %s", kode)
		if err := downloader.DownloadShipmentFile(kode, cfg.Token, prod); err != nil {
			fmt.Println("Error downloading shipment file:", err)
			skipDownload = true
		}
		log.Info().Msg("Shipment file downloaded successfully!")

		if skipDownload {
			ex.SetValue(r.RowIndex, "KOLI", "0")
			ex.SetValue(r.RowIndex, "ISI MUATAN", "0")
			continue
		}

		// Step 2: Generate pivot report
		report, err := reporting.GeneratePivotReport(config.GetExportedFilePath(prod), "Memuat dan membongkar ekspor in", true)
		if err != nil {
			fmt.Println("Error generating report:", err)
			return
		}

		ex.SetValue(r.RowIndex, "KOLI", strconv.Itoa(report.TotalPivotRows))
		totalMuatan := report.BlankBaggingCount + report.TotalWaybillCount
		ex.SetValue(r.RowIndex, "ISI MUATAN", strconv.Itoa(totalMuatan))
	}

	if prod {
		ex.Save("report_result.xlsx")
		return
	}
	ex.Save("public/report_result.xlsx")
}
