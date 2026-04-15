package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/umardev500/jnt-report/internal/manifest"
)

func main() {

	// 🔥 LIST OF ROUTES
	routes := []string{"BDO", "SRG", "SOC", "SUB", "BGR", "SEG", "BKI", "JKT", "JAT", "DPK"}
	date := "2026-04-10"
	filename := "manifest.yml"

	if err := manifest.SyncManifest(filename, routes, date); err != nil {
		log.Error().Err(err).Msg("failed to save routes")
		return
	}

	fmt.Println("✅ All routes saved successfully")
}
