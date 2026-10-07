package app

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/umardev500/jnt-report/internal/vendor"
)

func (app *App) VendorHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getVendors(w, r)

	case http.MethodPost:
		app.createVendor(w, r)

	case http.MethodPut:
		app.updateVendor(w, r)

	case http.MethodDelete:
		app.deleteVendor(w, r)

	default:
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

func (app *App) getVendors(
	w http.ResponseWriter,
	r *http.Request,
) {
	idString := r.URL.Query().Get("id")

	if idString != "" {
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(
				w,
				"invalid id",
				http.StatusBadRequest,
			)
			return
		}

		vendor, err := app.vendorRepo.GetByID(id)
		if err != nil {
			http.Error(
				w,
				"vendor not found",
				http.StatusNotFound,
			)
			return
		}

		writeJSON(w, http.StatusOK, vendor)
		return
	}

	vendors, err := app.vendorRepo.GetAll()
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, vendors)
}

func (app *App) createVendor(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input vendor.Vendor

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if input.Name == "" {
		http.Error(
			w,
			"name is required",
			http.StatusBadRequest,
		)
		return
	}

	if err := app.vendorRepo.Create(input); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusCreated, input)
}

func (app *App) updateVendor(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input vendor.Vendor

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if input.ID <= 0 || input.Name == "" {
		http.Error(
			w,
			"id and name are required",
			http.StatusBadRequest,
		)
		return
	}

	if err := app.vendorRepo.Update(input); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, input)
}

func (app *App) deleteVendor(
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

	if err := app.vendorRepo.Delete(id); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
