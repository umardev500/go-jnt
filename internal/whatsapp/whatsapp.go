package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type NotificationType string

const (
	NotificationSuratJalan NotificationType = "SURAT_JALAN"
	NotificationManifest   NotificationType = "MANIFEST"
	NotificationOutgoing   NotificationType = "OUTGOING"
)

const notificationGroupJID = "120363413605685992@g.us"

func SendGroupNotification(notificationType NotificationType) error {
	var message string

	switch notificationType {
	case NotificationSuratJalan:
		message = "🚚 Someone is processing Surat Jalan."
	case NotificationManifest:
		message = "📦 Someone is processing Manifest."
	case NotificationOutgoing:
		message = "📤 Someone is reporting Outgoing."
	default:
		return fmt.Errorf("unsupported notification type: %s", notificationType)
	}

	return SendNotification(notificationGroupJID, message)
}

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
