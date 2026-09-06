package whatsapp

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

func (w *Client) SendText(
	ctx context.Context,
	phone string,
	text string,
) error {
	if w.client == nil {
		return fmt.Errorf("WhatsApp client is not initialized")
	}

	if !w.client.IsConnected() {
		return fmt.Errorf("WhatsApp client is not connected")
	}

	jid := types.NewJID(phone, types.DefaultUserServer)

	msg := &waE2E.Message{
		Conversation: &text,
	}

	_, err := w.client.SendMessage(ctx, jid, msg)
	if err != nil {
		return fmt.Errorf("send WhatsApp message: %w", err)
	}

	return nil
}
