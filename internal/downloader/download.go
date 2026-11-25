package downloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/umardev500/jnt-report/internal/config"
	"github.com/umardev500/jnt-report/internal/httpclient"
)

// Core structs
type BongkarMuatRecord struct {
	JobCode         string `json:"jobCode"`
	ScanNetworkCode string `json:"scanNetworkCode"`
	ScanWaybillNum  int    `json:"scanWaybillNum"`
}

type ExportRecord struct {
	OssUrl   string `json:"ossUrl"`
	TaskName string `json:"taskName"`
}

// Wait helper
func wait(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// -------------------- API Steps --------------------

func FetchBongkarMuat(shipmentNo string, token string) (*BongkarMuatRecord, error) {
	url := fmt.Sprintf("https://jmsgw.jntexpress.id/transportation/trackingDeatil/loading/scan/list?shipmentNo=%s", shipmentNo)
	headers := map[string]string{
		"User-Agent":   "Go-http-client/1.1",
		"Accept":       "application/json",
		"authToken":    token,
		"lang":         "ID",
		"langType":     "ID",
		"timezone":     "GMT+0700",
		"Content-Type": "application/json;charset=utf-8",
	}

	var res struct {
		Data []BongkarMuatRecord `json:"data"`
	}

	if err := httpclient.GetJSON(url, headers, &res); err != nil {
		return nil, err
	}

	for _, r := range res.Data {
		if r.ScanNetworkCode == "BTN777" {
			return &r, nil
		}
	}
	return nil, errors.New("no record found for BTN777")
}

func TriggerExport(record *BongkarMuatRecord, token string) error {
	url := "https://jmsgw.jntexpress.id/transportation/trackingDeatil/exportScanList"
	headers := map[string]string{
		"User-Agent":   "Go-http-client/1.1",
		"Accept":       "application/json",
		"authToken":    token,
		"lang":         "ID",
		"langType":     "ID",
		"timezone":     "GMT+0700",
		"Content-Type": "application/json;charset=utf-8",
	}

	payload := map[string]any{
		"current":         1,
		"size":            20,
		"shipmentNo":      record.JobCode,
		"scanNetworkCode": record.ScanNetworkCode,
		"countNum":        record.ScanWaybillNum,
		"tag":             1,
		"countryId":       "1",
	}

	var res map[string]any
	return httpclient.PostJSON(url, headers, payload, &res)
}

func GetLastExportedRecordWithRetry(startTime, endTime string, maxRetries, delayMs int, token string) (*ExportRecord, error) {
	headers := map[string]string{
		"User-Agent":   "Go-http-client/1.1",
		"Accept":       "application/json",
		"authToken":    token,
		"lang":         "ID",
		"langType":     "ID",
		"timezone":     "GMT+0700",
		"Content-Type": "application/json;charset=utf-8",
	}

	var res struct {
		Data struct {
			Records []ExportRecord `json:"records"`
		} `json:"data"`
	}

	url := fmt.Sprintf(
		"https://jmsgw.jntexpress.id/transportation/tmsExportTransportReport/task?endTime=%s&startTime=%s&current=1&size=20&moduleType=",
		endTime,
		startTime,
	)

	for i := range maxRetries {
		attempt := i + 1

		log.Info().
			Int("attempt", attempt).
			Int("max_retries", maxRetries).
			Msg("Checking exported record")

		if err := httpclient.GetJSON(url, headers, &res); err != nil {
			log.Error().
				Err(err).
				Int("attempt", attempt).
				Msg("Failed to fetch exported record")

			return nil, err
		}

		if len(res.Data.Records) > 0 && res.Data.Records[0].OssUrl != "" {
			log.Info().
				Int("attempt", attempt).
				Msg("Exported record is ready")

			return &res.Data.Records[0], nil
		}

		if attempt < maxRetries {
			log.Warn().
				Int("attempt", attempt).
				Int("max_retries", maxRetries).
				Int("delay_ms", delayMs).
				Msg("Export not ready, retrying...")

			wait(delayMs)
		}
	}

	log.Error().
		Int("max_retries", maxRetries).
		Msg("Export is not ready after maximum retries")

	return nil, errors.New("export not ready")
}

func GetDownloadUrl(record *ExportRecord, token string) (string, error) {
	url := "https://jmsgw.jntexpress.id/transportation/file/oss/getDownloadSignedUrl"
	headers := map[string]string{
		"User-Agent":   "Go-http-client/1.1",
		"Accept":       "application/json",
		"authToken":    token,
		"lang":         "ID",
		"langType":     "ID",
		"timezone":     "GMT+0700",
		"Content-Type": "application/json;charset=utf-8",
	}

	payload := map[string]string{
		"data":      record.OssUrl,
		"taskName":  record.TaskName,
		"countryId": "1",
	}

	var res struct {
		Data string `json:"data"`
	}
	if err := httpclient.PostJSON(url, headers, payload, &res); err != nil {
		return "", err
	}
	if res.Data == "" {
		return "", errors.New("failed to get signed download URL")
	}
	return res.Data, nil
}

func DownloadFile(url, outputPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func DownloadShipmentFile(shipmentNo string, token string, prod bool) error {
	record, err := FetchBongkarMuat(shipmentNo, token)
	if err != nil {
		return err
	}

	if err := TriggerExport(record, token); err != nil {
		return err
	}

	now := time.Now()
	layout := "2006-01-02+15:04:05"

	// Start of today (00:00:00)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// End of today (23:59:59)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	log.Info().Msgf("Start time: %s", startOfDay.Format(layout))
	log.Info().Msgf("End time: %s", endOfDay.Format(layout))

	// startTime := "2026-04-08+00:00:00"
	// endTime := "2026-04-08+23:59:59"

	lastRecord, err := GetLastExportedRecordWithRetry(startOfDay.Format(layout), endOfDay.Format(layout), 50, 3000, token)
	if err != nil {
		return err
	}

	downloadUrl, err := GetDownloadUrl(lastRecord, token)
	if err != nil {
		return err
	}

	return DownloadFile(downloadUrl, config.GetExportedFilePath(prod))
}
