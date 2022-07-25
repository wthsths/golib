package gl_http

import (
	"net/url"
)

type QueryExtractor struct{}

func NewQueryExtractor() *QueryExtractor {
	return &QueryExtractor{}
}

// ReadAll returns top level query parameter pairs from input url.
func (e *QueryExtractor) ReadAll(url *url.URL) map[string]string {
	queryParams := url.Query()
	paramsToReturn := make(map[string]string, len(queryParams))
	for k, v := range queryParams {
		if len(v) > 0 {
			paramsToReturn[k] = v[0]
		}
	}

	return paramsToReturn
}
