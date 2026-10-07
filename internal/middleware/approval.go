package middleware

import (
	"net/http"

	"github.com/umardev500/jnt-report/internal/approval"
)

func approvalService(approvalService approval.Service, next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if !approvalService.IsValid() {

			http.Error(
				w,
				"approval required",
				http.StatusForbidden,
			)

			return
		}

		next(w, r)
	}
}
