package handlers

import (
	"encoding/json"
	"net/http"
)

func ReadRequestBody(r *http.Request, data interface{}) error {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(data); err != nil {
		return err
	}
	return nil
}
