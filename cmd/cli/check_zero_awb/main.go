package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// =========================
// Structs
// =========================

type APIResponse struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data []PackageData `json:"data"`
	Fail bool          `json:"fail"`
	Succ bool          `json:"succ"`
}

type PackageData struct {
	PackageNumber string   `json:"packageNumber"`
	WaybillIDs    []string `json:"waybillIds"`

	PackageNumbers any `json:"packageNumbers"`
	PackageName    any `json:"packageName"`
}

// =========================
// Fetch Function
// =========================

func FetchWaybillByPackageNumbers(
	authToken string,
	packageNumbers []string,
) (*APIResponse, error) {

	url := "https://jmsgw.jntexpress.id/operatingplatform/packScanList/waybillIdsByPackageNumber"

	jsonBody, err := json.Marshal(packageNumbers)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}

	// Headers
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("authToken", authToken)
	req.Header.Set("lang", "ID")
	req.Header.Set("langType", "ID")
	req.Header.Set("timezone", "GMT+0700")
	req.Header.Set("routeName", "trackingExpress")
	req.Header.Set("Cache-Control", "max-age=2, must-revalidate")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result APIResponse

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// =========================
// Main
// =========================

func main() {
	authToken := "fac06c03110946719816e7900b5c16a3"

	packages := []string{
		"F07536329",
		"F07536858",
		"F10315750",
	}

	result, err := FetchWaybillByPackageNumbers(
		authToken,
		packages,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Code: %d\n", result.Code)
	fmt.Printf("Message: %s\n", result.Msg)

	for _, pkg := range result.Data {
		count := len(pkg.WaybillIDs)

		if count <= 0 {
			// red highlight
			fmt.Printf("\033[31m\nPackage: %s\n", pkg.PackageNumber)
			fmt.Printf("Waybill Count: %d\033[0m\n", count)
		} else {
			fmt.Printf("\nPackage: %s\n", pkg.PackageNumber)
			fmt.Printf("Waybill Count: %d\n", count)
		}
	}
}
