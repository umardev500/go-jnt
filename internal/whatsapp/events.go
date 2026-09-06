package whatsapp

import (
	"fmt"

	"go.mau.fi/whatsmeow/types/events"
)

func (w *Client) handleEvent(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if v.Message == nil {
			return
		}

		message := v.Message.GetConversation()

		fmt.Printf(
			"\nWhatsApp message from %s: %s\n",
			v.Info.Sender.User,
			message,
		)
	}
}
