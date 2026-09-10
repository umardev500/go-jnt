package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/umardev500/jnt-report/internal/approval"
	"github.com/umardev500/jnt-report/internal/config"
	"github.com/umardev500/jnt-report/internal/detail"
	"github.com/umardev500/jnt-report/internal/downloader"
	"github.com/umardev500/jnt-report/internal/payment"
	"github.com/umardev500/jnt-report/internal/reporting"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	_ = godotenv.Load()

	prod := os.Getenv("APP_ENV") == "prod"
	log.Info().Msgf("Running in %v mode", prod)

	approval.InitDB()

	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	if err := payment.WaitForActivation(
		cfg,
		approval.IsValid,
		scanner,
		"qris.jpg",
	); err != nil {
		log.Fatal().
			Err(err).
			Msg("Payment activation failed")
	}

	log.Info().Msg("Starting JNT Report...")

	if !approval.IsValid() {
		log.Fatal().Msg("Approval not valid")
		return
	}

	// Flags
	withReport := flag.Bool("with-report", false, "Download shipment file, generate pivot report, and print summary")
	driverTimes := flag.Bool("driver-times", false, "Print only Waktu Berangkat Driver and Waktu Sampai Driver")
	flag.Parse()

	// Check shipment number
	if flag.NArg() < 1 {
		fmt.Println("Usage: cli [options] <shipmentNo>")
		flag.PrintDefaults()
		return
	}
	shipmentNo := flag.Arg(0)

	// Steps 1-3: Optional group
	if *withReport {
		// Step 1: Download shipment file
		if err := downloader.DownloadShipmentFile(shipmentNo, cfg.Token, prod); err != nil {
			fmt.Println("Error downloading shipment file:", err)
			return
		}
		fmt.Println("Shipment file downloaded successfully!")

		// Step 2: Generate pivot report
		report, err := reporting.GeneratePivotReport(config.GetExportedFilePath(prod), "Memuat dan membongkar ekspor in", true)
		if err != nil {
			fmt.Println("Error generating report:", err)
			return
		}

		// Step 3: Print report summary
		printReportSummary(report)
	}

	// Step 4: Always fetch shipment detail
	res, err := detail.GetShipmentDetail(shipmentNo, cfg.Token)
	if err != nil {
		fmt.Println("Error fetching shipment detail:", err)
		return
	}

	// Conditional output
	if *driverTimes {
		printDriverTimes(res) // New method to show only driver times
	} else {
		printShipmentDetail(res) // Existing full detail
	}
}

func sendWhatsAppNotification(phone, message string) error {
	payload := map[string]string{
		"phone":   phone,
		"message": message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8081/whatsapp/send",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("WhatsApp API returned status %d", resp.StatusCode)
	}

	return nil
}

func printReportSummary(report *reporting.ReportSummary) {
	// ANSI colors
	cyan := "\033[36m"
	magenta := "\033[35m"
	yellow := "\033[33m"
	reset := "\033[0m"

	fmt.Println()
	fmt.Println("================================================")
	fmt.Println("=============== PIVOT REPORT ==================")
	fmt.Println()

	fmt.Printf("%-25s: %s%s%s\n", "File Path", magenta, report.FilePath, reset)
	fmt.Printf("%-25s: %s%s%s\n", "Sheet Name", magenta, report.SheetName, reset)
	fmt.Printf("%-25s: %s%d%s\n", "Total Rows", cyan, report.TotalRows, reset)
	fmt.Printf("%-25s: %s%d%s\n", "Pivot Rows (non-blank)", cyan, report.TotalPivotRows, reset)
	fmt.Printf("%-25s: %s%d%s\n", "Waybill Total", yellow, report.TotalWaybillCount, reset)
	fmt.Printf("%-25s: %s%d%s\n", "Blank Bagging Count", yellow, report.BlankBaggingCount, reset)
	fmt.Printf("%-25s: %s%s%s\n", "Output File", magenta, report.OutputFile, reset)

	fmt.Println()
	fmt.Println("================================================")
	fmt.Println()
}

// helper to handle nil pointers with color
func coloredValue(p *string) string {
	if p == nil {
		return color.RedString("<nil>")
	}
	return color.GreenString(*p)
}

// PrintShipmentDetail prints a ShipmentResponse nicely
func printShipmentDetail(res *detail.ShipmentDetailResponse) {
	s := res.Data.ShipmentDetail

	fmt.Println("============== SHIPMENT ==============")
	fmt.Printf("Shipment Name     : %s\n", color.GreenString(s.ShipmentName))
	fmt.Printf("Shipment No       : %s\n", color.GreenString(s.ShipmentNo))
	fmt.Printf("Plate Number      : %s\n", color.GreenString(s.PlateNumber))
	fmt.Printf("Planned Departure : %s\n", coloredValue(&s.PlannedDepartureTime))
	fmt.Printf("Driver Name       : %s\n", color.GreenString(s.DriverName))
	fmt.Printf("Driver Contact    : %s\n", color.GreenString(s.DriverContact))
	fmt.Printf("Carrier Name      : %s\n", color.GreenString(s.CarrierName))
	// fmt.Printf("Vehicle Type      : %s\n", color.GreenString(s.ActualDepartureTime))
	fmt.Printf("Vehicle Type      : %s\n", color.GreenString(s.VehicleTypeName))
	fmt.Printf("Percentage        : %s\n", "")

	layout := "2006-01-02 15:04:05"
	departure, _ := time.Parse(layout, s.ActualDepartureTime)
	arrival, _ := time.Parse(layout, s.ActualArrivalTime)

	duration := arrival.Sub(departure)
	minutes := int(math.Round(duration.Minutes()))
	fmt.Println(duration, " ", minutes)

	fmt.Println("\n=============== STOPS =================")
	for i, stop := range s.TmsShipmentStopVOList {
		fmt.Printf("%d:\n", i+1)
		fmt.Printf("  Create Time             : %s\n", color.GreenString(stop.CreateTime))
		fmt.Printf("  Create By               : %s\n", color.GreenString(stop.CreateByName))
		fmt.Printf("  Network Code            : %s\n", color.GreenString(stop.NetworkCode))
		fmt.Printf("  Planned Departure       : %s\n", coloredValue(stop.PlannedDepartureTime))
		fmt.Printf("  Waktu Berangkat Driver  : %s\n", coloredValue(stop.AppDriverDeparture))
		fmt.Printf("  Planned Departure Day   : %s\n", coloredValue(stop.PlannedDepartureDay))
		fmt.Printf("  Planned Departure Hour  : %s\n", coloredValue(stop.PlannedDepartureHour))
		fmt.Printf("  Scan Kirim Mobil        : %s\n", coloredValue(stop.ScanTime))
		fmt.Printf("  Waktu Sampai Driver     : %s\n\n", coloredValue(stop.ActualArrivalTime))
	}
}

// PrintDriverTimes prints only Waktu Berangkat Driver and Waktu Sampai Driver for all stops
func printDriverTimes(res *detail.ShipmentDetailResponse) {
	s := res.Data.ShipmentDetail
	keyWidth := 25 // consistent alignment

	fmt.Println("=========== DRIVER TIMES ===========")
	for i, stop := range s.TmsShipmentStopVOList {
		fmt.Printf("Stop %d:\n", i+1)
		fmt.Printf("  %-*s: %s\n", keyWidth-2, "Waktu Berangkat Driver", coloredValue(stop.AppDriverDeparture))
		fmt.Printf("  %-*s: %s\n\n", keyWidth-2, "Waktu Sampai Driver", coloredValue(stop.ActualArrivalTime))
	}
}
