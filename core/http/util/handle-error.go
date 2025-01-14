package util

import (
	"net/http"

	"github.com/jackc/pgx/v5"
)

func HandleError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	if err == pgx.ErrNoRows {
		http.Error(w, err.Error(), http.StatusNotFound)
		return true
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
	return true
}
