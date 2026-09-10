package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func SendNotification(phone, message string) error {
	apiPort := os.Getenv("API_PORT")

	if apiPort == "" {
		return fmt.Errorf("API_PORT is not set")
	}

	payload := map[string]string{
		"phone":   phone,
		"message": message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"http://localhost:%s/whatsapp/send",
		apiPort,
	)

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(
			"WhatsApp API returned status %d",
			resp.StatusCode,
		)
	}

	return nil
}
