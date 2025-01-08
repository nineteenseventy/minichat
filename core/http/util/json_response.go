package util

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
