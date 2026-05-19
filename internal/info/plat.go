package info

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "https://jmsgw.jntexpress.id/transportation/tmsNewVehicle/page"

	Green = "\033[32m"
	Red   = "\033[31m"
	Reset = "\033[0m"
)

// ================= CLIENT =================

type Client struct {
	http  *http.Client
	token string
}

func NewClient(token string) *Client {
	return &Client{
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
		token: token,
	}
}

// ================= API RESPONSE =================

type VehicleResponse struct {
	Data struct {
		Records []Vehicle `json:"records"`
	} `json:"data"`
}

// ================= API MODEL =================

type Vehicle struct {
	ID              int     `json:"id"`
	PlateNumber     string  `json:"plateNumber"`
	CarrierName     string  `json:"carrierName"`
	VehicleTypeName string  `json:"vehicleTypeName"`
	InsideLength    float64 `json:"insideLength"`
	InsideWidth     float64 `json:"insideWidth"`
	InsideHeight    float64 `json:"insideHeight"`
	LoadWeight      float64 `json:"loadWeight"`
	VehicleVolume   float64 `json:"vehicleVolume"`
	AxleNumber      int     `json:"axleNumber"`
	LastUseTime     string  `json:"lastUseTime"`
}

// ================= FETCH =================

func (c *Client) GetPlateInfo(plates []string) ([]Vehicle, error) {
	var results []Vehicle

	for _, plate := range plates {
		url := fmt.Sprintf(
			"%s?current=1&size=20&plateNumber=%s&vehicleClassId=1&isOutage=1",
			baseURL,
			plate,
		)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("authToken", c.token)
		req.Header.Set("lang", "ID")
		req.Header.Set("timezone", "GMT+0700")

		resp, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
		}

		var apiResp VehicleResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, err
		}

		results = append(results, apiResp.Data.Records...)
	}

	return results, nil
}

func (v Vehicle) Pretty(status string) string {
	color := Green
	reset := Reset

	if status != "OK" {
		color = Red
	}

	return fmt.Sprintf(
		`Status       : %s%s%s
Plate        : %s%s%s
Carrier      : %s%s%s
Type         : %s%s%s
Dimension    : %.2f x %.2f x %.2f
Load         : %.2f ton
Volume       : %.2f
Axle         : %d
Last Used    : %s
-----------------------------------`,

		color, status, reset,

		color, v.PlateNumber, reset,
		color, v.CarrierName, reset,
		color, v.VehicleTypeName, reset,

		v.InsideLength,
		v.InsideWidth,
		v.InsideHeight,

		v.LoadWeight,
		v.VehicleVolume,
		v.AxleNumber,
		v.LastUseTime,
	)
}
