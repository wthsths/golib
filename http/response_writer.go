package http

import (
	"encoding/json"
	go_http "net/http"
)

// WriteCustomJsonResponse serializes input res and creates response payload from it.
func WriteCustomJsonResponse(w go_http.ResponseWriter, statusCode int, res interface{}) (writtenRes []byte, err error) {
	resJson, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_, err = w.Write(resJson)
	if err != nil {
		return nil, err
	}
	return resJson, nil
}
