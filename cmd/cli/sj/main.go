package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/jnt-report/internal/approval"
	"github.com/umardev500/jnt-report/internal/config"
	"github.com/umardev500/jnt-report/internal/detail"
	"github.com/umardev500/jnt-report/internal/excel"
	"github.com/umardev500/jnt-report/internal/payment"
	suratjalan "github.com/umardev500/jnt-report/internal/surat_jalan"
	"github.com/umardev500/jnt-report/internal/tracking/loadingscan"
	"github.com/umardev500/jnt-report/internal/types"
	"github.com/xuri/excelize/v2"
)

func printExcel(filePath string, prod bool) error {
	var err error
	path := filepath.Join(
		"C:\\", "Users", "User", "Projects", "go-report", filePath,
	)

	if prod {
		path, err = filepath.Abs(filePath)
		if err != nil {
			return err
		}

	}

	cmd := exec.Command(
		"C:\\Program Files\\Microsoft Office\\root\\Office16\\EXCEL.EXE",
		path,
	)

	return cmd.Run()
}

type Admin struct {
	FullName string
	Partner  string
}

var adminData = map[string]Admin{
	"RAHMAWAN": {
		FullName: "RAHMAWAN RAMADHAN PRIATNA",
		Partner:  "STEVE NICHOLAS FERNANDO LORENTE",
	},
	"STEVE NICHOLAS FERNANDO LORENTE": {
		FullName: "INDRA GUNAWAN",
		Partner:  "RAHMAWAN RAMADHAN PRIATNA",
	},
	"UMAR": {
		FullName: "UMAR",
		Partner:  "WISNU JATI KUSUMA",
	},
	"WISNU JATI KUSUMA": {
		FullName: "WISNU JATI KUSUMA",
		Partner:  "UMAR",
	},
	"IMAM": {
		FullName: "IMAM MANSUR SUJANA",
		Partner:  "ODEH",
	},
	"ODEH": {
		FullName: "ODEH",
		Partner:  "IMAM MANSUR SUJANA",
	},
}

func getFullName(name string) string {
	name = strings.ToUpper(strings.TrimSpace(name))
	if val, ok := adminData[name]; ok {
		return val.FullName
	}
	return ""
}

func getPartner(name string) string {
	name = strings.ToUpper(strings.TrimSpace(name))
	if val, ok := adminData[name]; ok {
		return val.Partner
	}
	return ""
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err)
	}

	fmt.Println(os.Getenv("AUTH_TOKEN"))
	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = zerolog.New(
		zerolog.ConsoleWriter{
			Out: os.Stdout,
		},
	).With().Timestamp().Logger()
}

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

type Scanner struct {
	Service string
	Code    string
}

func (s *Scanner) SetServiceFromFlag(eco bool) {
	if eco {
		s.Service = "ECO"
	} else {
		s.Service = "EZ"
	}
}

func checkDuplicateInColumnP(value string, prod bool) error {
	filePath := config.GetHistoryFilePath(prod)

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open history file: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return fmt.Errorf("no sheet found")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to read rows: %w", err)
	}

	for i, row := range rows {
		// Column P = index 15 (0-based)
		if len(row) <= 15 {
			continue
		}

		if row[15] == value {
			return fmt.Errorf("duplicate value found in column P at row %d: %s", i+1, value)
		}
	}

	return nil
}

func main() {
	_ = godotenv.Load()

	prod := os.Getenv("APP_ENV") == "prod"
	log.Info().Msgf("Running in %v mode", prod)

	now := time.Now()
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		panic(err)
	}

	approval.InitDB()
	scanner1 := bufio.NewScanner(os.Stdin)

	if err := payment.WaitForActivation(
		cfg,
		approval.IsValid,
		scanner1,
		"qris.jpg",
	); err != nil {
		log.Fatal().
			Err(err).
			Msg("Payment activation failed")
	}

	// Flags
	flagEco := flag.Bool("eco", false, "use ECO service")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("usage: go run main.go [--eco] <code>")
		return
	}

	code := args[0]
	admin := os.Getenv("ADMIN")
	if admin == "" {
		log.Fatal().Msg("ADMIN env variable is not set")
		return
	}

	// Load loading scan
	log.Info().Msg("Loading loading scan...")

	_, err = loadingscan.Get(code, cfg.Token)
	if err != nil {
		fmt.Printf("Error fetching loading scan: %v\n", err)

		var answer string
		fmt.Print("Continue anyway? (Y/N): ")
		fmt.Scanln(&answer)

		switch strings.ToUpper(answer) {
		case "Y":
			log.Warn().Msg("Continuing despite the error...")
		case "N":
			log.Info().Msg("Operation cancelled.")
			return
		default:
			log.Info().Msg("Invalid input. Operation cancelled.")
			return
		}
	}

	log.Info().Msg("Loading scan loaded successfully")

	if err := checkDuplicateInColumnP(code, prod); err != nil {
		log.Fatal().Err(err).Msg("Duplicate error")
		return
	}

	// Load routes
	routes, err := loadGatewayRoutes(config.GetSKOFilePath(prod))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load gateway routes")
	}

	// Scanner
	scanner := &Scanner{Code: code}
	scanner.SetServiceFromFlag(*flagEco)

	log.Info().
		Str("service", scanner.Service).
		Str("code", scanner.Code).
		Msg("scanner initialized")

	// Open Excel files
	ex := excel.Open(config.GetSuratJalanFilePath(prod))
	exHist := excel.Open(config.GetHistoryFilePath(prod))

	// Generate Surat Jalan
	dt, report := suratjalan.GenateSuratJalan(
		ex,
		routes,
		scanner.Code,
		cfg.Token,
		getFullName(admin),
		scanner.Service,
		prod,
	)
	if dt == nil && report == nil {
		log.Error().Msg("failed to generate surat jalan")
		return
	}

	// Extract route code safely
	routeCode := detail.GetRouteCode(dt)
	t, err := time.Parse("2006-01-02 15:04:05", dt.Data.ShipmentDetail.PlannedDepartureTime)
	if err != nil {
		panic(err)
	}

	filenameTime := t.Format("15_04")

	// Create history + sequence
	seq, err := suratjalan.CreateHistory(
		exHist,
		routes,
		dt,
		report,
		routeCode,
		scanner.Service,
		fmt.Sprintf("%s / %s", getFullName(admin), getPartner(admin)),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create history")
		return
	}

	// Generate document number
	noSurat := generateNoSurat(routeCode, seq, now, scanner.Service)
	ex.SetValue(8, "G", noSurat)

	// Set print area
	// setPrintArea(ex)

	// Build output p
	sjOut := buildOutputPath(now, fmt.Sprintf("%s_%s", routeCode, filenameTime))

	// Save files
	ex.Save(sjOut)

	// Save history
	if prod {
		exHist.Save("sj_history.xlsx")
	} else {
		exHist.Save("assets/sj_history.xlsx")
	}

	// exHist.Save("assets/sj_history.xlsx")

	// Print
	fmt.Println("printing....")
	if err := printExcel(sjOut, prod); err != nil {
		log.Error().Err(err).Msg("print failed")
	}

	log.Info().Msg("done")
}

//
// ======================= HELPERS =======================
//

func getRouteCode(endName string) string {
	parts := strings.Split(endName, "_")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func generateNoSurat(routeCode string, seq int, now time.Time, service string) string {
	return fmt.Sprintf(
		"BTN%s%02d/%s%s",
		routeCode,
		seq,
		now.Format("20060102"),
		service,
	)
}

func setPrintArea(ex *excel.ExcelFile) {
	err := ex.File.SetDefinedName(&excelize.DefinedName{
		Name:     "_xlnm.Print_Area",
		RefersTo: "BTN!$B$2:$L$41",
		Scope:    "BTN",
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to set print area")
	}
}

func buildOutputPath(now time.Time, routeCode string) string {

	folder := filepath.Join(
		"generated", "excel", "sj",
		now.Format("2006_01_02"),
	)

	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		log.Fatal().Err(err).Msg("failed to create folder")
	}

	filename := fmt.Sprintf(
		"sj_print_%s_%s.xlsm",
		routeCode,
		now.Format("20060102_150405"),
	)

	return filepath.Join(folder, filename)
}
