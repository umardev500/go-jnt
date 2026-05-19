package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/umardev500/jnt-report/internal/config"
	"github.com/umardev500/jnt-report/internal/manifest"
)

func main() {
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		panic(err)
	}

	// 🔥 LIST OF ROUTES
	routes := []string{"BDO", "SRG", "SOC", "SUB", "BGR", "SEG", "BKI", "JKT", "JAT", "DPK"}
	date := "2026-04-15"
	filename := "manifest.yml"

	if err := manifest.SyncManifest(filename, routes, date, cfg.Token); err != nil {
		log.Error().Err(err).Msg("failed to save routes")
		return
	}

	fmt.Println("✅ All routes saved successfully")
}
