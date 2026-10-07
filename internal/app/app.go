package app

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/umardev500/jnt-report/internal/approval"
	"github.com/umardev500/jnt-report/internal/unit"
	"github.com/umardev500/jnt-report/internal/user"
	"github.com/umardev500/jnt-report/internal/vendor"
)

type App struct {
	db *sql.DB

	userRepo   *user.Repository
	vendorRepo *vendor.Repository
	unitRepo   *unit.Repository
	approval   *approval.Service
}

func New(db *sql.DB) *App {
	return &App{
		db: db,

		userRepo:   user.NewRepository(db),
		vendorRepo: vendor.NewRepository(db),
		unitRepo:   unit.NewRepository(db),
		approval:   approval.New(db),
	}
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}
