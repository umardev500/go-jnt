package whatsapp

import (
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

func (w *Client) SendText(
	ctx context.Context,
	target string,
	text string,
) error {
	if w.client == nil {
		return fmt.Errorf("WhatsApp client is not initialized")
	}

	if !w.client.IsConnected() {
		return fmt.Errorf("WhatsApp client is not connected")
	}

	var jid types.JID

	if strings.HasSuffix(target, "@g.us") {
		var err error
		jid, err = types.ParseJID(target)
		if err != nil {
			return fmt.Errorf("invalid group JID: %w", err)
		}
	} else {
		jid = types.NewJID(target, types.DefaultUserServer)
	}

	msg := &waE2E.Message{
		Conversation: &text,
	}

	_, err := w.client.SendMessage(ctx, jid, msg)
	if err != nil {
		return fmt.Errorf("send WhatsApp message: %w", err)
	}

	return nil
}
