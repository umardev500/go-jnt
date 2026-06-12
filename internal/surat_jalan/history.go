package suratjalan

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/umardev500/jnt-report/internal/detail"
	"github.com/umardev500/jnt-report/internal/excel"
	"github.com/umardev500/jnt-report/internal/helper"
	"github.com/umardev500/jnt-report/internal/reporting"
	"github.com/umardev500/jnt-report/internal/route"
	"github.com/umardev500/jnt-report/internal/types"
	"github.com/xuri/excelize/v2"
)

func GetLatestSequenceByDestination(f *excel.ExcelFile, targetDest, service string) (int, error) {
	now := time.Now()
	rows := f.GetRows()

	count := 0

	for _, row := range rows {
		kode := row.Values["No.Surat Jalan"]
		code := strings.TrimSpace(kode)
		if code == "" {
			log.Info().Msg("Skipping empty KODE JMS")
			continue
		}

		jadwalKeberangkatan := row.Values["Jadwal Keberangkatan"]
		fmt.Println(jadwalKeberangkatan)

		t, err := helper.ParseExcelDateTime(jadwalKeberangkatan)
		if err != nil {
			fmt.Println("err history 1")
			log.Fatal().Err(err).Msg("invalid date")
		}

		if t.Year() == now.Year() &&
			t.Month() == now.Month() &&
			t.Day() == now.Day() {
		} else {
			log.Info().Msgf("Skipping date %s", jadwalKeberangkatan)
			continue
		}

		namaRute := row.Values["Nama Rute"]
		if namaRute != targetDest {
			log.Info().Msgf("Skipping destination %s", targetDest)
			continue
		}

		layanan := row.Values["Layanan"]
		if layanan != service {
			log.Info().Msgf("Skipping service %s", service)
			continue
		}

		log.Info().Msgf("Processing... %s", row.Values["Kode Tugas"])
		count++

	}

	log.Info().Msgf("Count: %d", count+1)

	return count + 1, nil
}

func CreateHistory(ex *excel.ExcelFile, routes types.GatewayRoutes, dt *detail.ShipmentDetailResponse, report *reporting.ReportSummary, dest, service, admin string) (int, error) {
	now := time.Now()
	layout := "2006-01-02 15:04:05"
	shipment := dt.Data.ShipmentDetail
	origin := shipment.TmsShipmentStopVOList[0]
	routeCode := detail.GetRouteCode(dt)

	seq, _ := GetLatestSequenceByDestination(ex, dest, service)
	noSurat := fmt.Sprintf("BTN%s%02d/%s%s", routeCode, seq, now.Format("20060102"), service)

	f := ex.File
	rows := ex.GetRows()
	lastRow := len(rows)
	newRow := lastRow + 2
	planOperasi := shipment.PlannedDepartureTime
	planDate, err := time.Parse(layout, planOperasi)
	if err != nil {
		panic(err)
	}

	lastValues := rows[lastRow-1].Values

	green := "\033[32m"
	cyan := "\033[36m"
	yellow := "\033[33m"
	reset := "\033[0m"

	fmt.Println("=== Last Row Values ===")
	fmt.Printf("%sCW%s      : %v\n", green, reset, lastValues["CW"])
	fmt.Printf("%sColly%s   : %v\n", cyan, reset, lastValues["Colly"])
	fmt.Printf("%sCW Luar%s : %v\n", yellow, reset, lastValues["CW luar"])

	// New data
	fmt.Println("=== New Data ===")
	fmt.Printf("%sCW%s      : %v\n", green, reset, report.TotalWaybillCount)
	fmt.Printf("%sColly%s   : %v\n", cyan, reset, report.TotalPivotRows)
	fmt.Printf("%sCW Luar%s : %v\n", yellow, reset, report.BlankBaggingCount)

	cw, err := strconv.Atoi(lastValues["CW"])
	if err != nil {
		return 0, fmt.Errorf("invalid CW value in lastValues: %v", err)
	}

	sameCW := cw == report.TotalWaybillCount

	if sameCW {
		log.Error().Msg("CW is the same, skipping")
		return 0, fmt.Errorf("CW is the same, skipping")
	}

	// Make log data for ok tasks
	log.Info().Msg("=== OK ===")

	departureTime, err := time.Parse("2006-01-02 15:04:05", planOperasi)
	if err != nil {
		fmt.Println("err history")
		log.Err(err).Msg("invalid date")
	}
	rs := route.GetRitase(routes, shipment.EndName, departureTime.Format("15:04"))
	jenisPaket := route.GetJenisPaket(rs.StatusRute)
	fmt.Println("Jenis: ", jenisPaket, rs.StatusRute)

	var appTrackDepartureTime time.Time
	if origin.ActualDepartureTime != nil {
		appTrackDepartureTime, err = time.Parse(layout, *origin.ActualDepartureTime)
		if err != nil {
			log.Error().Err(err).Msg("Error parsing app track departure time")
		}
	}

	f.InsertRows("Sheet1", newRow, 1)
	ex.SetValue(newRow, "No.Surat Jalan", noSurat)
	ex.SetValue(newRow, "Tanggal Operasi", planDate)
	ex.SetValue(newRow, "Nama Rute", routeCode)
	ex.SetValue(newRow, "Plat Nomor", shipment.PlateNumber)
	ex.SetValue(newRow, "Jenis Mobil", shipment.VehicleTypeName)
	ex.SetValue(newRow, "Jenis Trip", strconv.Itoa(seq))
	// ex.SetValue(newRow, "Jadwal Keberangkatan", planDate)
	ex.SetValue(newRow, "Waktu Berangkat Mobil", appTrackDepartureTime.Format("15:04")+" / "+jenisPaket)
	ex.SetValue(newRow, "Colly", strconv.Itoa(report.TotalPivotRows))
	ex.SetValue(newRow, "CW", strconv.Itoa(report.TotalWaybillCount))
	ex.SetValue(newRow, "CW luar", strconv.Itoa(report.BlankBaggingCount))
	// ex.SetValue(newRow, "Waktu Rilis", now)
	// ex.SetValue(newRow, "Waktu Rilis", now.Format("02/01/2006 15:04:05"))

	fmtStr := "dd/mm/yyyy hh:mm:ss"
	styleID, err := ex.File.NewStyle(&excelize.Style{
		CustomNumFmt: &fmtStr,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error creating style")
		return 0, err
	}

	cell := fmt.Sprintf("L%d", newRow)

	ex.File.SetCellValue("Sheet1", cell, now)
	ex.File.SetCellStyle("Sheet1", cell, cell, styleID)

	// Set jadwal keberangkatan to use date style

	fmtStr = "dd/mm/yyyy"
	styleID, err = ex.File.NewStyle(&excelize.Style{
		CustomNumFmt: &fmtStr,
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center", // optional
		},
	})
	cell = fmt.Sprintf("G%d", newRow)
	ex.File.SetCellValue("Sheet1", cell, appTrackDepartureTime.Format("02/01/2006"))
	ex.File.SetCellStyle("Sheet1", cell, cell, styleID)

	// ex.File.SetCellValue("Sheet1", "L:"+strconv.Itoa(newRow), now)
	ex.SetValue(newRow, "Layanan", service)
	ex.SetValue(newRow, "Admin / SPV", admin)
	ex.SetValue(newRow, "Vendor", helper.MapVendor(shipment.CarrierName))
	ex.SetValue(newRow, "Kode Tugas", shipment.ShipmentNo)

	return seq, nil
}
