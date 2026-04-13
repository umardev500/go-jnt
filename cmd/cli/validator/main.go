package main

import (
	"github.com/umardev500/jnt-report/internal/config"
	"github.com/umardev500/jnt-report/internal/excel"
	"github.com/umardev500/jnt-report/internal/validator"
)

func main() {
	ex := excel.Open(config.GetHistoryFilePath())
	validator.ValidateSuratJalan(ex, config.AuthToken)
}
