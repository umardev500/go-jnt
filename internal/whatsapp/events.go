package whatsapp

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/umardev500/jnt-report/internal/helper"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

const groupAdmin = "120363413605685992@g.us"

func (w *Client) handleEvent(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if v.Message == nil {
			log.Debug().Msg("Ignored message")
			return
		}

		// Accept messages only from your specified group.
		if v.Info.Chat.String() != groupAdmin {
			return
		}

		message := strings.TrimSpace(v.Message.GetConversation())
		if message == "" {
			message = strings.TrimSpace(
				v.Message.GetExtendedTextMessage().GetText(),
			)
		}
		if message == "" {
			return
		}

		fmt.Printf(
			"\nWhatsApp message from %s: %s\n",
			v.Info.Sender.User,
			message,
		)

		// Only respond to group commands.
		if v.Info.Chat.Server != types.GroupServer {
			return
		}

		var response string

		switch strings.ToLower(message) {
		case "/menu":
			response = "🤖 SERVER MENU\n\n" +
				"/menu - Show commands\n" +
				"/acc - Grant approval\n" +
				"/revoke - Revoke approval"

		case "/acc":
			response = "⏳ Processing approval..."

			// TODO: Trigger your approval service here.
			if err := helper.CallCurrentTunnel("/granted"); err != nil {
				log.Error().Err(err).Msg("Grant failed")
				// Send an error reply to the group.

				return
			}

		case "/revoke":
			response = "⏳ Processing revoke..."

			// TODO: Trigger your revoke service here.
			if err := helper.CallCurrentTunnel("/revoke"); err != nil {
				log.Error().Err(err).Msg("Revoke failed")
				// Send an error reply to the group.
				return
			}
		default:
			return
		}

		// IMPORTANT:
		// Use the correct destination expected by your API.
		// A group JID is not the same as a phone number.
		if err := SendNotification(
			groupAdmin,
			response,
		); err != nil {
			fmt.Printf("Failed to send command response: %v\n", err)
		}
	}
}
