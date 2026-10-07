package app

import (
	"encoding/json"
	"net/http"

	"github.com/umardev500/jnt-report/internal/user"
)

func (app *App) UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getUsers(w, r)

	case http.MethodPost:
		app.createUser(w, r)

	case http.MethodPut:
		app.updateUser(w, r)

	case http.MethodDelete:
		app.deleteUser(w, r)

	default:
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

func (app *App) getUsers(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.URL.Query().Get("id")

	if id != "" {
		user, err := app.userRepo.GetByID(id)
		if err != nil {
			http.Error(
				w,
				"user not found",
				http.StatusNotFound,
			)
			return
		}

		writeJSON(w, http.StatusOK, user)
		return
	}

	users, err := app.userRepo.GetAll()
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, users)
}

func (app *App) createUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input user.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if input.ID == "" {
		http.Error(
			w,
			"id is required",
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

	if err := app.userRepo.Create(input); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusCreated, input)
}

func (app *App) updateUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input user.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if input.ID == "" || input.Name == "" {
		http.Error(
			w,
			"id and name are required",
			http.StatusBadRequest,
		)
		return
	}

	if err := app.userRepo.Update(input); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, input)
}

func (app *App) deleteUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(
			w,
			"id is required",
			http.StatusBadRequest,
		)
		return
	}

	if err := app.userRepo.Delete(id); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
