package app

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/umardev500/jnt-report/internal/unit"
)

func (app *App) UnitHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getUnits(w, r)

	case http.MethodPost:
		app.createUnit(w, r)

	case http.MethodPut:
		app.updateUnit(w, r)

	case http.MethodDelete:
		app.deleteUnit(w, r)

	default:
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

func (app *App) getUnits(
	w http.ResponseWriter,
	r *http.Request,
) {
	idString := r.URL.Query().Get("id")

	if idString != "" {
		id, err := strconv.Atoi(idString)
		if err != nil || id <= 0 {
			http.Error(
				w,
				"invalid id",
				http.StatusBadRequest,
			)
			return
		}

		unit, err := app.unitRepo.GetByID(id)
		if err != nil {
			http.Error(
				w,
				"unit not found",
				http.StatusNotFound,
			)
			return
		}

		writeJSON(w, http.StatusOK, unit)
		return
	}

	units, err := app.unitRepo.GetAll()
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, units)
}

func (app *App) createUnit(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input unit.Unit

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if input.VendorID <= 0 {
		http.Error(
			w,
			"vendor_id is required",
			http.StatusBadRequest,
		)
		return
	}

	if input.PlateNumber == "" {
		http.Error(
			w,
			"plate_number is required",
			http.StatusBadRequest,
		)
		return
	}

	if input.UnitTypeID <= 0 {
		http.Error(
			w,
			"unit_type_id is required",
			http.StatusBadRequest,
		)
		return
	}

	if input.CreatedBy == "" {
		http.Error(
			w,
			"created_by is required",
			http.StatusBadRequest,
		)
		return
	}

	if err := app.unitRepo.Create(input); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusCreated, input)
}

func (app *App) updateUnit(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input unit.Unit

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if input.ID <= 0 {
		http.Error(
			w,
			"id is required",
			http.StatusBadRequest,
		)
		return
	}

	if input.VendorID <= 0 {
		http.Error(
			w,
			"vendor_id is required",
			http.StatusBadRequest,
		)
		return
	}

	if input.PlateNumber == "" {
		http.Error(
			w,
			"plate_number is required",
			http.StatusBadRequest,
		)
		return
	}

	if input.UnitTypeID <= 0 {
		http.Error(
			w,
			"unit_type_id is required",
			http.StatusBadRequest,
		)
		return
	}

	if input.UpdatedBy == "" {
		http.Error(
			w,
			"updated_by is required",
			http.StatusBadRequest,
		)
		return
	}

	if err := app.unitRepo.Update(input); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, input)
}

func (app *App) deleteUnit(
	w http.ResponseWriter,
	r *http.Request,
) {
	idString := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idString)
	if err != nil || id <= 0 {
		http.Error(
			w,
			"invalid id",
			http.StatusBadRequest,
		)
		return
	}

	if err := app.unitRepo.Delete(id); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
