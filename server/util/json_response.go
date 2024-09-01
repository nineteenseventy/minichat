package util

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(v)
	if err != nil {
		return err
	}
	return nil
}
