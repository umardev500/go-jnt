package manifest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

// ================= REQUEST =================

type shipmentRequest struct {
	Current       int    `json:"current"`
	Size          int    `json:"size"`
	ShipmentName  string `json:"shipmentName"`
	TmsType       int    `json:"tmsType"`
	StartCode     string `json:"startCode"`
	StartDateTime string `json:"startDateTime"`
	EndDateTime   string `json:"endDateTime"`
	SearchType    string `json:"searchType"`
	CountryId     string `json:"countryId"`
}

// ================= RESPONSE =================

type shipmentResponse struct {
	Data struct {
		Records []struct {
			ShipmentNo    string `json:"shipmentNo"`
			ShipmentState int    `json:"shipmentState"`
			PlateNumber   string `json:"plateNumber"`
		} `json:"records"`
	} `json:"data"`
}

// ================= YAML STRUCT =================

type CodesYML map[string][]string

// ================= API FUNCTION =================

func GetShipmentCodes(
	authToken string,
	current int,
	size int,
	shipmentName string,
	startDateTime string,
	endDateTime string,
) ([]string, error) {

	url := "https://jmsgw.jntexpress.id/transportation/tmsShipment/post/page"

	reqBody := shipmentRequest{
		Current:       current,
		Size:          size,
		ShipmentName:  shipmentName,
		TmsType:       1,
		StartCode:     "BTN777",
		StartDateTime: startDateTime,
		EndDateTime:   endDateTime,
		SearchType:    "manage",
		CountryId:     "1",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("authToken", authToken)
	req.Header.Set("lang", "ID")
	req.Header.Set("langType", "ID")
	req.Header.Set("timezone", "GMT+0700")

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

	var result shipmentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var shipmentCodes []string

	for _, r := range result.Data.Records {

		// skip deleted
		if r.ShipmentState == 5 {
			msg := fmt.Sprintf("ℹ️  Skipping deleted shipment %s", r.ShipmentNo)
			fmt.Println(msg)
			continue
		}

		// Log if plate number is dummy which is like B1234HQ
		// Never skipped just log
		if r.PlateNumber == "B1234HQ" {
			msg := fmt.Sprintf("⚠️ Dummy plate number %s for shipment %s", r.PlateNumber, r.ShipmentNo)
			fmt.Println(msg)
		}

		shipmentCodes = append(shipmentCodes, r.ShipmentNo)
	}

	return shipmentCodes, nil
}

//
// ================= YAML LOAD / SAVE =================
//

func LoadCodes(filename string) (CodesYML, error) {
	data := make(CodesYML)

	file, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return nil, err
	}

	err = json.Unmarshal(file, &data) // NOTE: YAML removed? (fix below)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func SaveCodes(filename string, data CodesYML) error {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, out, 0644)
}

//
// ================= UTIL =================
//

func Unique(list []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(list))

	for _, v := range list {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

//
// ================= MAIN WORKFLOW =================
//

func SyncManifest(
	filename string,
	routes []string,
	date string,
	token string,
) error {

	// reset file each run
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		return err
	}

	data, err := LoadCodes(filename)
	if err != nil {
		return err
	}

	from := date + " 00:00:00"
	to := date + " 23:59:59"

	for _, route := range routes {

		fmt.Println()
		fmt.Printf("🔎 Fetching codes for route %s\n", route)

		codes, err := GetShipmentCodes(
			token,
			1,
			100,
			route,
			from,
			to,
		)
		if err != nil {
			log.Error().Err(err).Msgf("Route failed: %s", route)
			continue
		}

		data[route] = Unique(append(data[route], codes...))

		fmt.Println("✔ Saved:", route, "count:", len(data[route]))
	}

	if err := SaveCodes(filename, data); err != nil {
		return err
	}

	log.Info().Msg("All routes saved successfully")
	return nil
}
