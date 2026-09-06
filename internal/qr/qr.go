package qr

import (
	"os"

	"github.com/mdp/qrterminal/v3"
)

func Print(data string) {
	config := qrterminal.Config{
		Level:      qrterminal.M,
		Writer:     os.Stdout,
		HalfBlocks: true,
		QuietZone:  1,
	}

	qrterminal.GenerateWithConfig(data, config)
}
