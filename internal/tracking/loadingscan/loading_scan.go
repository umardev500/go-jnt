package loadingscan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const endpoint = "https://jmsgw.jntexpress.id/transportation/trackingDeatil/loading/scan/list"

type Scan struct {
	JobCode              string `json:"jobCode"`
	LoadingTypeName      string `json:"loadingTypeName"`
	ScanStartTime        string `json:"scanStartTime"`
	ScanEndTime          string `json:"scanEndTime"`
	ScanTotalTime        string `json:"scanTotalTime"`
	ScanNetworkName      string `json:"scanNetworkName"`
	ScanNetworkCode      string `json:"scanNetworkCode"`
	ScanWaybillSingleNum int    `json:"scanWaybillSingleNum"`
	ScanPackageNum       int    `json:"scanPackageNum"`
	ScanWaybillNum       int    `json:"scanWaybillNum"`
	OnlyLoadCount        *int   `json:"onlyLoadCount"`
	OnlyUnloadCount      *int   `json:"onlyUnloadCount"`
}

type Response struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    []Scan `json:"data"`
	Details any    `json:"details"`
	Fail    bool   `json:"fail"`
	Succ    bool   `json:"succ"`
}

func Get(shipmentNo, authToken string) (*Scan, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s?shipmentNo=%s", endpoint, url.QueryEscape(shipmentNo)),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("authToken", authToken)
	req.Header.Set("lang", "ID")
	req.Header.Set("langType", "ID")
	req.Header.Set("timezone", "GMT+0700")
	req.Header.Set("routeName", "monitoringSearchView")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Succ || result.Code != 1 {
		return nil, fmt.Errorf("%s", result.Msg)
	}

	switch len(result.Data) {
	case 0:
		return nil, fmt.Errorf("loading scan not found")
	case 1:
		return &result.Data[0], nil
	default:
		return nil, fmt.Errorf("expected 1 loading scan, got %d", len(result.Data))
	}
}
