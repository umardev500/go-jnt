package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/cors"
	"github.com/umardev500/jnt-report/internal/config"
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
	ActualUseTime            *int     `json:"actualUseTime"`
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
	UseWayTime               *int     `json:"useWayTime"`
	VehicleDoorCnt           *int     `json:"vehicleDoorCnt"`
	VehiclelineCode          string   `json:"vehiclelineCode"`
	VehiclelineName          string   `json:"vehiclelineName"`
	VehicletypeName          string   `json:"vehicletypeName"`
	StationWaitingTime       *int     `json:"stationWaitingTime"`
	CubeNumber               *float64 `json:"cubeNumber"`
	Promotion                int      `json:"promotion"`
	Shifts                   int      `json:"shifts"`
	OperationModel           int      `json:"operationModel"`
	Mileage                  int      `json:"mileage"`
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

// ================== APP ==================

type App struct {
	cfg    *config.Config
	client *http.Client
}

func NewApp(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// ================== CLIENT ==================

func (a *App) fetchReport(payload RequestPayload) (*Response, error) {

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
	req.Header.Set("authToken", a.cfg.Token)

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

	fmt.Printf(
		"[REQUEST] ip=%s method=%s path=%s ua=%s\n",
		clientIP,
		r.Method,
		r.URL.Path,
		r.UserAgent(),
	)

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

// ================== HELPERS ==================

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

// ================== MAIN ==================

func main() {

	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		panic(err)
	}

	app := NewApp(cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/report", app.reportHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	fmt.Println("Listening on port 8080...")

	if err := http.ListenAndServe(
		":8080",
		c.Handler(mux),
	); err != nil {
		panic(err)
	}
}
