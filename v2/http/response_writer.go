package gl_http

import (
	"encoding/json"
	"net/http"

	gl_constants "github.com/payports/golib/v2/constants"
)

type responseWriter struct {
}

// NewResponseWriter creates an utility instance for handling of http responses.
func NewResponseWriter() *responseWriter {
	return &responseWriter{}
}

// WriteCustomJsonResponse serializes input res and creates response payload from it.
func (hu *responseWriter) WriteCustomJsonResponse(w http.ResponseWriter, statusCode int, res interface{}) (writtenRes []byte, err error) {
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

	resPrettyJson, _ := json.MarshalIndent(res, "", gl_constants.JsonIndentDefault)
	return resPrettyJson, nil
}
