package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/umardev500/jnt-report/internal/whatsapp"
)

func main() {
	client, err := whatsapp.New()
	if err != nil {
		fmt.Printf("Failed to initialize WhatsApp: %v\n", err)
		os.Exit(1)
	}

	defer client.Disconnect()

	if err := client.Connect(); err != nil {
		fmt.Printf("Failed to connect WhatsApp: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("WhatsApp is ready.")

	// Get WhatsApp groups.
	groups, err := client.GetGroups(context.Background())
	if err != nil {
		fmt.Printf("Failed to get groups: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("WhatsApp Groups:")
	fmt.Println("----------------")

	for _, group := range groups {
		fmt.Printf("Name: %s\n", group.Name)
		fmt.Printf("JID:  %s\n", group.JID)
		fmt.Println()
	}

	fmt.Println("Press Ctrl+C to exit.")

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	<-ctx.Done()

	fmt.Println()
	fmt.Println("Disconnecting WhatsApp...")
}
