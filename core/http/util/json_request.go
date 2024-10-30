package util

import (
	"encoding/json"
	"net/http"
)

func JSONRequest(request *http.Request, value interface{}) error {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(value)
	if err != nil {
		return err
	}
	return nil
}
