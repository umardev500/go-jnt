package suratjalan

import (
	"fmt"
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
)

func GenateSuratJalan(ex *excel.ExcelFile, routes types.GatewayRoutes, code, token, admin, service string) (*detail.ShipmentDetailResponse, *reporting.ReportSummary) {
	dt, err := detail.GetShipmentDetail(code, token)
	if err != nil {
		log.Error().Err(err).Str("Kode Tugas", code).Msg("Error fetching shipment detail")
	}
	shipment := dt.Data.ShipmentDetail

	now := time.Now()
	today := time.Now().Format("2006-01-02")
	layout := "2006-01-02 15:04:05"
	plannedDepartureTime, err := time.Parse(layout, shipment.PlannedDepartureTime)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing planned departure time")
	}

	departureTime, err := time.Parse("2006-01-02 15:04:05", shipment.PlannedDepartureTime)
	if err != nil {
		fmt.Println("err sj")
		log.Error().Err(err).Msg("invalid date")
	}
	rs := route.GetRitase(routes, shipment.EndName, departureTime.Format("15:04"))
	jenisPaket := route.GetJenisPaket(rs.StatusRute)
	fmt.Println("Jenis: ", jenisPaket, rs.StatusRute)

	ex.SetValue(32, "C", admin)
	ex.SetValue(9, "E", today)
	ex.SetValue(10, "E", plannedDepartureTime.Format("15:04")+" / "+jenisPaket)
	ex.SetValue(11, "E", shipment.PlateNumber)
	ex.SetValue(12, "E", strings.ToUpper(shipment.DriverName))
	ex.SetValue(13, "E", shipment.VehicleTypeName)
	ex.SetValue(14, "E", shipment.DriverContact)
	ex.SetValue(15, "E", now.Format("15:04:05"))

	ex.SetValue(18, "E", shipment.CarrierName)

	ex.SetValue(17, "I", "Kode Tugas :"+code)

	ex.SetValue(16, "J", service)

	// Service
	ex.SetValue(3, "T", service)

	// PREPARE THE EXCEL FILE
	// Get pivot info
	// Step 1: Download shipment file
	log.Info().Msgf("Downloading shipment file... %s", code)
	if err := downloader.DownloadShipmentFile(code, token); err != nil {
		fmt.Println("Error downloading shipment file:", err)
	}
	log.Info().Msg("Shipment file downloaded successfully!")

	// Step 2: Generate pivot report
	report, err := reporting.GeneratePivotReport(config.GetExportedFilePath(), "Memuat dan membongkar ekspor in", true)
	if err != nil {
		fmt.Println("Error generating report:", err)
		return nil, nil
	}

	routeCode := detail.GetRouteCode(dt)

	// INFO KOLI
	ex.SetValue(21, "G", strconv.Itoa(report.TotalPivotRows))
	ex.SetValue(25, "G", strconv.Itoa(report.TotalPivotRows))
	ex.SetValue(32, "G", strings.ToUpper(shipment.DriverName))
	ex.SetValue(20, "H", strconv.Itoa(report.BlankBaggingCount))
	ex.SetValue(21, "H", strconv.Itoa(report.TotalWaybillCount))
	ex.SetValue(25, "H", strconv.Itoa(report.TotalWaybillCount+report.BlankBaggingCount))

	dest := shipment.TmsShipmentStopVOList[1]

	ex.SetValue(10, "I", dest.NetworkName)
	address := data.Routes[routeCode]
	ex.SetValue(11, "I", address)

	return dt, report
}
