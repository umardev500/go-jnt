package whatsapp

import (
	"context"
	"fmt"
	"os"

	_ "modernc.org/sqlite"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"

	"github.com/mdp/qrterminal/v3"
)

type Client struct {
	client    *whatsmeow.Client
	container *sqlstore.Container
}

func New() (*Client, error) {
	ctx := context.Background()

	dbLog := waLog.Stdout("Database", "ERROR", true)

	container, err := sqlstore.New(
		ctx,
		"sqlite",
		"file:data/whatsapp/whatsapp.db?_pragma=foreign_keys(1)",
		dbLog,
	)
	if err != nil {
		return nil, fmt.Errorf("create WhatsApp database: %w", err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("get WhatsApp device: %w", err)
	}

	clientLog := waLog.Stdout("WhatsApp", "INFO", true)

	client := whatsmeow.NewClient(
		deviceStore,
		clientLog,
	)

	wc := &Client{
		client:    client,
		container: container,
	}

	client.AddEventHandler(wc.handleEvent)

	return wc, nil
}

func (w *Client) Connect() error {
	ctx := context.Background()

	// Already paired.
	if w.client.Store.ID != nil {
		fmt.Println("WhatsApp session found.")
		fmt.Println("Connecting to WhatsApp...")

		if err := w.client.Connect(); err != nil {
			return fmt.Errorf("connect WhatsApp: %w", err)
		}

		fmt.Println("WhatsApp connected.")

		return nil
	}

	// New device, need QR code.
	fmt.Println("No WhatsApp session found.")
	fmt.Println("Preparing QR code...")
	fmt.Println()

	qrChan, err := w.client.GetQRChannel(ctx)
	if err != nil {
		return fmt.Errorf("get QR channel: %w", err)
	}

	if err := w.client.Connect(); err != nil {
		return fmt.Errorf("connect WhatsApp: %w", err)
	}

	for evt := range qrChan {
		switch evt.Event {
		case "code":
			fmt.Println()
			fmt.Println("Scan this QR code with WhatsApp:")
			fmt.Println()

			qrterminal.GenerateHalfBlock(
				evt.Code,
				qrterminal.L,
				os.Stdout,
			)

			fmt.Println()
			fmt.Println("Waiting for WhatsApp login...")

		case "success":
			fmt.Println()
			fmt.Println("WhatsApp login successful.")

		default:
			fmt.Printf("WhatsApp login event: %s\n", evt.Event)
		}
	}

	return nil
}

func (w *Client) Disconnect() {
	if w.client != nil {
		w.client.Disconnect()
	}

	if w.container != nil {
		_ = w.container.Close()
	}
}
