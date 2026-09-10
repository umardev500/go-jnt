package payment

import (
	"bufio"
	"fmt"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/umardev500/jnt-report/internal/config"
	"github.com/umardev500/jnt-report/internal/whatsapp"
)

const adminPhone = "6283142765573"

func paymentSubmittedMessage(name string) string {
	return fmt.Sprintf(`💳 *Payment Received*

Hi %s!

Thank you! Your payment has been submitted successfully.

⏳ We are now verifying your payment. Please keep this application open while we complete the confirmation.

You will receive another notification once your payment has been confirmed.`, name)
}

func paymentConfirmedMessage(name string) string {
	return fmt.Sprintf(`✅ *Payment Confirmed*

Hi %s!

Your payment has been successfully confirmed.

🎉 Your application is now *activated* and ready to use.

Thank you for your payment!`, name)
}

type User struct {
	Name    string
	StaffNo string
	Phone   string
}

func WaitForActivation(
	cfg *config.Config,
	isApproved func() bool,
	scanner *bufio.Scanner,
	qrisPath string,
) error {
	users := map[string]User{
		"01830317": {
			StaffNo: "01830317",
			Name:    "RAHMAWAN",
			Phone:   "6281389741086",
		},
		"01829395": {
			StaffNo: "01829395",
			Name:    "UMAR",
			Phone:   "6283142765573",
		},
	}

	user, ok := users[cfg.StaffNo]
	if !ok {
		return fmt.Errorf("user not found for staff no: %s", cfg.StaffNo)
	}

	// Already activated.
	if isApproved() {
		return nil
	}

	log.Info().
		Str("staffNo", cfg.StaffNo).
		Str("user", user.Name).
		Msg("Payment required to activate this app.")

	log.Info().
		Msg("Scan the QR code below to complete payment and activate your app.")

	log.Info().
		Msg("QR code opened in a new window.")

	if err := openQRIS(qrisPath); err != nil {
		log.Warn().
			Err(err).
			Msg("Failed to open QR code")
	}

	log.Info().Msg("Press ENTER after you have completed the payment...")

	if !scanner.Scan() {
		return fmt.Errorf("failed to read user input")
	}

	if err := notifyAdminPaymentSubmitted(user); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to notify admin about payment submission")
	}

	// Payment has been submitted.
	if err := whatsapp.SendNotification(
		user.Phone,
		paymentSubmittedMessage(user.Name),
	); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to send payment waiting notification")
	}

	log.Info().
		Str("user", user.Name).
		Msg("Waiting for payment confirmation...")

	// Keep checking until payment is confirmed.
	for {
		if isApproved() {
			log.Info().
				Str("user", user.Name).
				Msg("Payment confirmed. App activated.")

			// Notify user.
			if err := whatsapp.SendNotification(
				user.Phone,
				paymentConfirmedMessage(user.Name),
			); err != nil {
				log.Error().
					Err(err).
					Msg("Failed to send payment confirmation notification")
			}

			// Notify admin.
			if err := notifyAdminPaymentConfirmed(user); err != nil {
				log.Error().
					Err(err).
					Msg("Failed to notify admin")
			}

			return nil
		}

		time.Sleep(3 * time.Second)
	}
}

func notifyAdminPaymentSubmitted(user User) error {
	message := fmt.Sprintf(
		`💳 *Payment Submitted*

A user has submitted a payment for activation.

👤 *Name:* %s
🆔 *Staff No:* %s
📱 *WhatsApp:* %s
⏳ *Payment:* Waiting for confirmation
🚀 *Application:* Pending activation

Please verify the payment and approve the user.`,
		user.Name,
		user.StaffNo,
		user.Phone,
	)

	return whatsapp.SendNotification(
		adminPhone,
		message,
	)
}

func notifyAdminPaymentConfirmed(user User) error {
	message := fmt.Sprintf(
		`💰 *Payment Confirmed*

A user's payment has been successfully confirmed.

👤 *Name:* %s
🆔 *Staff No:* %s
📱 *WhatsApp:* %s
✅ *Payment:* Confirmed
🚀 *Application:* Activated

The user has been notified successfully.`,
		user.Name,
		user.StaffNo,
		user.Phone,
	)

	return whatsapp.SendNotification(
		adminPhone,
		message,
	)
}

func openQRIS(path string) error {
	if path == "" {
		return fmt.Errorf("QRIS path is empty")
	}

	return exec.Command(
		"powershell",
		"-Command",
		"Start-Process",
		path,
	).Run()
}
