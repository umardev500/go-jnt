package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	app1 "github.com/umardev500/jnt-report/internal/app"
	"github.com/umardev500/jnt-report/internal/approval"
	"github.com/umardev500/jnt-report/internal/config"
	"github.com/umardev500/jnt-report/internal/db"
	"github.com/umardev500/jnt-report/internal/whatsapp"
	"github.com/xuri/excelize/v2"
)

const reportURL = "https://jmsgw.jntexpress.id/transportation/tmsShipmentEvent/report"

// ================== MODELS ==================

type Response struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    Data   `json:"data"`
	Details any    `json:"details"`
	Fail    bool   `json:"fail"`
	Succ    bool   `json:"succ"`
}

type Data struct {
	Records     []Record `json:"records"`
	Total       int      `json:"total"`
	Size        int      `json:"size"`
	Current     int      `json:"current"`
	SearchCount bool     `json:"searchCount"`
	Pages       int      `json:"pages"`
}

// keep your full Record struct here
type Record struct {
	TmsType                  int      `json:"tmsType"`
	ActualArrivalTime        *string  `json:"actualArrivalTime"`
	CarrierCheckoutAgentCode string   `json:"carrierCheckoutAgentCode"`
	CarrierCheckoutAgentName string   `json:"carrierCheckoutAgentName"`
	VehicleBelongCode        int      `json:"vehicleBelongCode"`
	VehicleBelongName        string   `json:"vehicleBelongName"`
	Total                    int      `json:"total"`
	ActualBatchTime          *string  `json:"actualBatchTime"`
	ActualDepartureTime      *string  `json:"actualDepartureTime"`
	ActualStopTime           *string  `json:"actualStopTime"`
	ActualUseTime            *string  `json:"actualUseTime"`
	AgingType                *string  `json:"agingType"`
	ArriveNetworkCode        string   `json:"arriveNetworkCode"`
	ArriveNetworkName        string   `json:"arriveNetworkName"`
	ArriveProvince           *string  `json:"arriveProvince"`
	BusinessAttribute        int      `json:"businessAttribute"`
	CarrierName              string   `json:"carrierName"`
	CarrierType              int      `json:"carrierType"`
	DelayStopTime            *string  `json:"delayStopTime"`
	DelayTime                *string  `json:"delayTime"`
	DriverContact            *string  `json:"driverContact"`
	DriverName               string   `json:"driverName"`
	EstimateBatchTime        *string  `json:"estimateBatchTime"`
	InBoundWeighTime         *string  `json:"inBoundWeighTime"`
	LoadCapacity             *float64 `json:"loadCapacity"`
	LoadingScanEndTime       string   `json:"loadingScanEndTime"`
	LoadingScanStartTime     string   `json:"loadingScanStartTime"`
	LoadingScanTotalTime     int      `json:"loadingScanTotalTime"`
	LoadWeight               float64  `json:"loadWeight"`
	OutBoundGrossWeight      *float64 `json:"outBoundGrossWeight"`
	OutBoundSuttleWeight     *float64 `json:"outBoundSuttleWeight"`
	OutBoundTareWeight       *float64 `json:"outBoundTareWeight"`
	OutBoundWeighTime        *string  `json:"outBoundWeighTime"`
	PlannedArrivalTime       string   `json:"plannedArrivalTime"`
	PlannedDepartureTime     string   `json:"plannedDepartureTime"`
	PlateNumber              string   `json:"plateNumber"`
	PlateColor               int      `json:"plateColor"`
	PredictArriveTime        *string  `json:"predictArriveTime"`
	QuotationModel           int      `json:"quotationModel"`
	ScanTime                 *string  `json:"scanTime"`
	SendNetworkCode          string   `json:"sendNetworkCode"`
	SendNetworkName          string   `json:"sendNetworkName"`
	ShipmentName             string   `json:"shipmentName"` // NEW
	ShipmentNo               string   `json:"shipmentNo"`
	ShipmentState            int      `json:"shipmentState"`
	ShipmentType             int      `json:"shipmentType"`
	StopTime                 *string  `json:"stopTime"`
	TardyTime                *string  `json:"tardyTime"`
	TrackInTime              *string  `json:"trackInTime"`
	TrackOutTime             *string  `json:"trackOutTime"`
	TrailerNumber            *string  `json:"trailerNumber"`
	UnLoadLineTime           *string  `json:"unLoadLineTime"`
	UnLoadingScanEndTime     *string  `json:"unLoadingScanEndTime"`
	UnLoadingScanStartTime   *string  `json:"unLoadingScanStartTime"`
	UnLoadingScanTotalTime   *int     `json:"unLoadingScanTotalTime"`
	UnScanTime               *string  `json:"unScanTime"`
	UseTime                  int      `json:"useTime"`
	UseWayTime               *string  `json:"useWayTime"`
	VehicleDoorCnt           *int     `json:"vehicleDoorCnt"`
	VehiclelineCode          string   `json:"vehiclelineCode"`
	VehiclelineName          string   `json:"vehiclelineName"`
	VehicletypeName          string   `json:"vehicletypeName"`
	StationWaitingTime       *string  `json:"stationWaitingTime"`
	CubeNumber               *float64 `json:"cubeNumber"`
	Promotion                int      `json:"promotion"`
	Shifts                   int      `json:"shifts"`
	OperationModel           int      `json:"operationModel"`
	Mileage                  float64  `json:"mileage"`
	CarrierShortName         string   `json:"carrierShortName"`
	ApplyReasonItem          *int     `json:"applyReasonItem"`
	ApplyReason              *string  `json:"applyReason"`
	AuditStatus              int      `json:"auditStatus"`
	AuditRemark              *string  `json:"auditRemark"`
	Auditer                  *string  `json:"auditer"`
	VehicleTypegroup         string   `json:"vehicleTypegroup"`
	AxleNumber               int      `json:"axleNumber"`
	VehicleOrigin            string   `json:"vehicleOrigin"`
	OvertimeType             *int     `json:"overtimeType"`
	OvertimeReasons          *string  `json:"overtimeReasons"`
	OriRegShiftCarrierId     *string  `json:"oriRegShiftCarrierId"`
	OriRegShiftCarrierName   *string  `json:"oriRegShiftCarrierName"`
	IsAssistLine             int      `json:"isAssistLine"`
	BillingWay               int      `json:"billingWay"`
	FreightCode              *string  `json:"freightCode"`
	IsStop                   int      `json:"isStop"`
	AppTrackDepartureTime    *string  `json:"appTrackDepartureTime"`
	AppTrackArrivalTime      *string  `json:"appTrackArrivalTime"`
	CarrierId                string   `json:"carrierId"`
	CarrierCode              *string  `json:"carrierCode"`
	BandLineId               int      `json:"bandLineId"`
	BandLineName             string   `json:"bandLineName"`
	LoadingStatus            int      `json:"loadingStatus"`
	StartHandlingType        string   `json:"startHandlingType"`
	EndHandlingType          string   `json:"endHandlingType"`
	ScanPackageNum           int      `json:"scanPackageNum"`
	ScanWaybillNum           int      `json:"scanWaybillNum"`
	DriverActualUseTime      *int     `json:"driverActualUseTime"`
	DriverUseWayTime         *int     `json:"driverUseWayTime"`
	IsGuaranteeLine          int      `json:"isGuaranteeLine"`
	HasReturnCar             int      `json:"hasReturnCar"`
	ActualHasReturnCar       int      `json:"actualHasReturnCar"`
	UnScanPackageNum         *int     `json:"unScanPackageNum"`
	UnScanWaybillNum         *int     `json:"unScanWaybillNum"`
	WinCarrierStatus         *string  `json:"winCarrierStatus"`
	WinCarrierRanking        *int     `json:"winCarrierRanking"`
	WinCarrierId             *string  `json:"winCarrierId"`
	WinCarrierName           *string  `json:"winCarrierName"`
	IsOneHourInterval        int      `json:"isOneHourInterval"`
	LinePartVolumeRate       *string  `json:"linepartvolumerate"`
}

// ================== REQUEST ==================

type RequestPayload struct {
	Current         int    `json:"current"`
	Size            int    `json:"size"`
	ShipmentState   int    `json:"shipmentState"`
	TmsType         int    `json:"tmsType"`
	SendNetworkCode string `json:"sendNetworkCode"`
	TimeType        int    `json:"timeType"`
	StartTime       string `json:"startTime"`
	EndTime         string `json:"endTime"`
	CountryID       string `json:"countryId"`
	ShipmentName    string `json:"shipmentName,omitempty"`
}

type VehicleCheckPayload struct {
	ShipmentNo  string `json:"shipmentNo"`
	PlateNumber string `json:"plateNumber"`
	VehicleType string `json:"vehicleType"`
}

type ExcelRow struct {
	VehicleType string
	Status      string
}

type VehicleCheckResult struct {
	ShipmentNo  string `json:"shipmentNo"`
	PlateNumber string `json:"plateNumber"`
	VehicleType string `json:"vehicleType"`
	Status      string `json:"status,omitempty"`
	Found       bool   `json:"found"`
}

// ================== APP ==================

type App struct {
	configStore *config.Store
	client      *http.Client
	whatsapp    *whatsapp.Client
}

func NewApp(
	store *config.Store,
	whatsappClient *whatsapp.Client,
) *App {
	return &App{
		configStore: store,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		whatsapp: whatsappClient,
	}
}

// ================== CLIENT ==================

func (a *App) fetchReport(payload RequestPayload) (*Response, error) {
	cfg := a.configStore.Load()

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		reportURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("authToken", cfg.Token)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Response

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func getClientIP(r *http.Request) string {

	// cloudflare
	ip := r.Header.Get("CF-Connecting-IP")
	if ip != "" {
		return ip
	}

	// proxy / load balancer
	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}

	// nginx
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// fallback
	return r.RemoteAddr
}

// ================== HANDLER ==================

func (a *App) reportHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// log client info
	clientIP := getClientIP(r)

	log.Info().
		Str("ip", clientIP).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("ua", r.UserAgent()).
		Msg("incoming request")

	var req RequestPayload

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	applyDefaults(&req)

	data, err := a.fetchReport(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *App) checkVehicleHandler(w http.ResponseWriter, r *http.Request) {
	var payload []VehicleCheckPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	lookup, err := loadUnits("ACTUAL_VEHICLE.xlsx")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to load excel", http.StatusInternalServerError)
		return
	}

	results := make([]VehicleCheckResult, 0, len(payload))

	for _, p := range payload {
		row, found := lookup[p.PlateNumber]

		if !found {
			results = append(results, VehicleCheckResult{
				ShipmentNo:  p.ShipmentNo,
				PlateNumber: p.PlateNumber,
				VehicleType: p.VehicleType,
				Found:       false,
			})
			continue
		}

		// if row.VehicleType != p.VehicleType {
		// 	results = append(results, VehicleCheckResult{
		// 		ShipmentNo:  p.ShipmentNo,
		// 		PlateNumber: p.PlateNumber,
		// 		VehicleType: p.VehicleType,
		// 		Status:      "vehicle type mismatch",
		// 		Found:       false,
		// 	})
		// 	continue
		// }

		results = append(results, VehicleCheckResult{
			ShipmentNo:  p.ShipmentNo,
			PlateNumber: p.PlateNumber,
			VehicleType: p.VehicleType,
			Status:      row.Status,
			Found:       true,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    results,
	})
}

func (app *App) updateTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Token   string `json:"token"`
		Name    string `json:"name"`
		StaffNo string `json:"staffNo"`
		Email   string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Token == "" {
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}

	if err := app.configStore.UpdateUserInfo(req.Token, req.Name, req.StaffNo, req.Email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("token updated"))
}

func (a *App) whatsappSendHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	var req struct {
		Phone   string `json:"phone"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if req.Phone == "" {
		http.Error(
			w,
			"phone is required",
			http.StatusBadRequest,
		)
		return
	}

	if req.Message == "" {
		http.Error(
			w,
			"message is required",
			http.StatusBadRequest,
		)
		return
	}

	err := a.whatsapp.SendText(
		r.Context(),
		req.Phone,
		req.Message,
	)

	if err != nil {
		log.Error().
			Err(err).
			Str("phone", req.Phone).
			Msg("Failed to send WhatsApp message")

		http.Error(
			w,
			"failed to send WhatsApp message",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "WhatsApp message sent",
	})
}

// ================== HELPERS ==================
func loadUnits(filePath string) (map[string]ExcelRow, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("VEHICLES")
	if err != nil {
		return nil, err
	}

	lookup := make(map[string]ExcelRow)

	for i, row := range rows {
		if i == 0 || len(row) < 3 {
			continue
		}

		plate := row[1]
		lookup[plate] = ExcelRow{
			VehicleType: row[2],
			Status:      "none",
		}
	}

	return lookup, nil
}

func applyDefaults(req *RequestPayload) {

	if req.Current == 0 {
		req.Current = 1
	}

	if req.Size == 0 {
		req.Size = 100
	}

	if req.TmsType == 0 {
		req.TmsType = 1
	}

	if req.SendNetworkCode == "" {
		req.SendNetworkCode = "BTN777"
	}

	if req.TimeType == 0 {
		req.TimeType = 2
	}

	if req.CountryID == "" {
		req.CountryID = "1"
	}
}

func init() {
	// Pretty console logging
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		},
	)
}

func watchConfig(store *config.Store) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal().Err(err)
	}
	defer watcher.Close()

	if err := watcher.Add("."); err != nil {
		log.Fatal().Err(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if filepath.Base(event.Name) != "config.yml" {
				continue
			}

			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
				if err := store.Reload(); err != nil {
					log.Err(err).Msg("Config reload failed")
				}
			}

		case err := <-watcher.Errors:
			log.Err(err).Msg("Config watcher error")
		}
	}
}

// ================== MAIN ==================

func main() {
	database, err := db.Init()
	if err != nil {
		log.Fatal().Err(err)
	}
	defer database.Close()

	approvalService := approval.New(database)

	configStore := &config.Store{}

	if err := configStore.Reload(); err != nil {
		panic(err)
	}

	go watchConfig(configStore)

	// ================== WHATSAPP ==================

	wa, err := whatsapp.New()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to initialize WhatsApp")
	}

	go func() {
		log.Info().Msg("Starting WhatsApp connection...")

		if err := wa.Connect(); err != nil {
			log.Error().
				Err(err).
				Msg("WhatsApp connection failed")

			return
		}

		log.Info().Msg("WhatsApp connected")
	}()

	defer wa.Disconnect()

	// ================== APP ==================

	app := NewApp(configStore, wa)

	mux := http.NewServeMux()
	apps := app1.New(database)

	mux.HandleFunc("/users", apps.UserHandler)
	mux.HandleFunc("/vendors", apps.VendorHandler)
	mux.HandleFunc("/units", apps.UnitHandler)

	mux.HandleFunc("/report", app.reportHandler)
	mux.HandleFunc("/check-vehicle", app.checkVehicleHandler)
	mux.HandleFunc("/whatsapp/send", app.whatsappSendHandler)

	mux.HandleFunc("/granted", func(w http.ResponseWriter, r *http.Request) {
		err := approvalService.Grant(725)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Write([]byte("approved for 15 minutes"))
	})

	mux.HandleFunc("/revoke", func(w http.ResponseWriter, r *http.Request) {
		err := approvalService.Revoke()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("approval revoked"))
	})

	mux.HandleFunc("/token", app.updateTokenHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	apiPort := os.Getenv("API_PORT")

	if apiPort == "" {
		apiPort = "8080"
	}

	log.Info().
		Str("port", apiPort).
		Msg("API listening")

	if err := http.ListenAndServe(
		":"+apiPort,
		c.Handler(mux),
	); err != nil {
		log.Fatal().
			Err(err).
			Msg("API server stopped")
	}

}
