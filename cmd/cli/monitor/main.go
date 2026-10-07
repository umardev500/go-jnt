package main

import (
	"bufio"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/jnt-report/internal/approval"
	"github.com/umardev500/jnt-report/internal/config"
	"github.com/umardev500/jnt-report/internal/db"
	"github.com/umardev500/jnt-report/internal/monitor"
	"github.com/umardev500/jnt-report/internal/payment"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	log := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()
	log = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	start := time.Now()

	_ = godotenv.Load()

	prod := os.Getenv("APP_ENV") == "prod"
	log.Info().Msgf("Running in %v mode", prod)

	database, err := db.Init()
	if err != nil {
		log.Fatal().Err(err)
	}
	defer database.Close()

	approvalService := approval.New(database)

	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	if err := payment.WaitForActivation(
		cfg,
		approvalService.IsValid,
		scanner,
		"qris.jpg",
	); err != nil {
		log.Fatal().
			Err(err).
			Msg("Payment activation failed")
	}

	log.Info().Msg("Starting monitor...")

	err = monitor.Process(
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
