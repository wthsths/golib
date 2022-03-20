package gl_http

import "net/url"

type queryExtractor struct{}

func NewQueryExtractor() *queryExtractor {
	return &queryExtractor{}
}

// ReadAll returns top level query parameter pairs from input url.
func (e *queryExtractor) ReadAll(url *url.URL) map[string]interface{} {
	returnMap := make(map[string]interface{}, len(url.Query()))
	queryParams := url.Query()

	for k, v := range queryParams {
		if len(v) > 0 {
			returnMap[k] = v[0]
		}
	}

	return returnMap
}
