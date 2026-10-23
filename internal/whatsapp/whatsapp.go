package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/umardev500/jnt-report/internal/config"
	"github.com/umardev500/jnt-report/internal/helper"
)

type NotificationType string

const (
	NotificationSuratJalan NotificationType = "SURAT_JALAN"
	NotificationManifest   NotificationType = "MANIFEST"
	NotificationOutgoing   NotificationType = "OUTGOING"
)

const notificationGroupJID = "120363413605685992@g.us"

func SendGroupNotification(notificationType NotificationType) error {
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		return err
	}

	user, err := helper.GetUser(cfg.StaffNo)
	if err != nil {
		return err
	}

	var action string

	switch notificationType {
	case NotificationSuratJalan:
		action = "processing Surat Jalan"
	case NotificationManifest:
		action = "processing Manifest"
	case NotificationOutgoing:
		action = "reporting Outgoing"
	default:
		return fmt.Errorf("unsupported notification type: %s", notificationType)
	}

	// Message for the group
	groupMessage := fmt.Sprintf(
		"🔔 Activity Notification\n\n"+
			"👤 Name: %s\n"+
			"🆔 Staff No: %s\n"+
			"📱 Phone: %s\n"+
			"⚙️ Action: %s.",
		user.Name,
		user.StaffNo,
		user.Phone,
		action,
	)

	// Message for the user
	userMessage := fmt.Sprintf(
		"Hi %s, 👋\n\n"+
			"This is to notify you that you are currently %s.",
		user.Name,
		action,
	)

	// Send to group
	if err := SendNotification(notificationGroupJID, groupMessage); err != nil {
		return fmt.Errorf("send group notification: %w", err)
	}

	// Send directly to the user
	// userJID := fmt.Sprintf("%s@s.whatsapp.net", user.Phone)
	if err := SendNotification(user.Phone, userMessage); err != nil {
		return fmt.Errorf("send user notification: %w", err)
	}

	return nil
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
