package main

import (
	"github.com/umardev500/jnt-report/internal/config"
	"github.com/umardev500/jnt-report/internal/excel"
	"github.com/umardev500/jnt-report/internal/validator"
)

func main() {
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		panic(err)
	}

	ex := excel.Open(config.GetHistoryFilePath())
	validator.ValidateSuratJalan(ex, cfg.Token)
}
