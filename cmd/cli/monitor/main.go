package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/jnt-report/internal/monitor"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	log := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()

	start := time.Now()

	err := monitor.Process(
		"monitor.xlsx",
		"monitor_result.xlsx",
		"sheet1",
		"KLASIFIKASI.xlsx",
		log,
	)

	elapsed := time.Since(start).Seconds()

	if err != nil {
		log.Fatal().
			Err(err).
			Float64("elapsed_seconds", elapsed).
			Msg("failed")
	}

	log.Info().
		Float64("elapsed_seconds", elapsed).
		Msg("completed")
}
