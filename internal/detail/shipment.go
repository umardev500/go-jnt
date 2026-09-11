package detail

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Shipment represents the main shipment info
type Shipment struct {
	ShipmentName         string `json:"shipmentName"`
	ShipmentNo           string `json:"shipmentNo"`
	PlateNumber          string `json:"plateNumber"`
	DriverName           string `json:"driverName"`
	DriverContact        string `json:"driverContact"`
	CarrierName          string `json:"carrierName"`
	VehicleTypeName      string `json:"vehicletypeName"`
	PlannedDepartureTime string `json:"plannedDepartureTime"`
	EndName              string `json:"endName"`
}

// Stop represents each stop
type Stop struct {
	CreateTime           string  `json:"createTime"`
	CreateByName         string  `json:"createByName"`
	CarrierCheckoutAgent *string `json:"carrierCheckoutAgentName"`
	PlannedDepartureTime *string `json:"plannedDepartureTime"`
	ActualDepartureTime  *string `json:"actualDepartureTime"`
	PlannedDepartureDay  *string `json:"plannedDepartureDay"`
	PlannedDepartureHour *string `json:"plannedDepartureHour"`
	ScanTime             *string `json:"appDepartureTime"`
	AppDriverDeparture   *string `json:"appTrackDepartureTime"`
	ActualArrivalTime    *string `json:"actualArrivalTime"`
	NetworkCode          string  `json:"networkCode"`
	NetworkName          string  `json:"networkName"`
}

// ShipmentDetailResponse is the top-level JSON response
type ShipmentDetailResponse struct {
	Data struct {
		ShipmentDetail struct {
			StartCode             string `json:"startCode"`
			EndCode               string `json:"endCode"`
			ShipmentName          string `json:"shipmentName"`
			ShipmentNo            string `json:"shipmentNo"`
			PlateNumber           string `json:"plateNumber"`
			DriverName            string `json:"driverName"`
			DriverContact         string `json:"driverContact"`
			CarrierName           string `json:"carrierName"`
			VehicleTypeName       string `json:"vehicletypeName"`
			EndName               string `json:"endName"`
			TmsShipmentStopVOList []Stop `json:"tmsShipmentStopVOList"`
			PlannedDepartureTime  string `json:"plannedDepartureTime"`
			ActualDepartureTime   string `json:"actualDepartureTime"`
			ActualArrivalTime     string `json:"actualArrivalTime"`
			TotalRuntime          int    `json:"totalRuntime"`
		} `json:"shipmentDetail"`
	} `json:"data"`
}

// ---------------- Service Function ----------------

// GetShipmentDetail fetches shipment info by shipmentNo and returns structured data
func GetShipmentDetail(shipmentNo string, authToken string) (*ShipmentDetailResponse, error) {
	url := fmt.Sprintf("https://jmsgw.jntexpress.id/transportation/tmsShipment/traceDetail?shipmentNo=%s", shipmentNo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// headers from JS fetch
	req.Header.Set("User-Agent", "Go-http-client/1.1")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("authToken", authToken)
	req.Header.Set("lang", "ID")
	req.Header.Set("langType", "ID")
	req.Header.Set("timezone", "GMT+0700")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res ShipmentDetailResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func GetRouteCode(dt *ShipmentDetailResponse) string {
	shipment := dt.Data.ShipmentDetail

	destination := shipment.TmsShipmentStopVOList[1]
	parts := strings.Split(destination.NetworkName, "_")
	routeCode := parts[0]

	switch routeCode {
	case "PMH2":
		return "PMH"
	case "GSK2":
		return "GSK"
	default:
		return routeCode
	}
}
