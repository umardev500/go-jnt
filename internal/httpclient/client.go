package httpclient

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func GetJSON(url string, headers map[string]string, result interface{}) error {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(result)
}

func PostJSON(url string, headers map[string]string, payload any, result any) error {
	body, _ := json.Marshal(payload)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(result)
}
