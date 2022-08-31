package gl_routing

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Route_Rules(t *testing.T) {
	routeRules := []*RouteRule{
		// Tests comparisons are made according to indices.
		// Therefore changing order of route rule definitions can break test results.
		{Method: `GET`, Path: `/api/accounts`, DynamicPath: false},
		{Method: `GET`, Path: `/api/transfers/{uniqueID}`, DynamicPath: true},
		{Method: `GET`, Path: `/api/entity/{id}/reference/{ref}`, DynamicPath: true},
		{Method: `GET`, Path: `/api/transfers`, DynamicPath: false},
		{Method: `POST`, Path: `/api/transfers`, DynamicPath: false},
	}

	router, err := NewRouter(routeRules)
	assert.NoError(t, err)

	type testData struct {
		fullPath       string
		method         string
		routeRuleIndex int
		match          bool
		routeParams    map[string]string
	}

	data := []testData{
		{fullPath: `/api/accounts`, method: `GET`, routeRuleIndex: 0, match: true, routeParams: nil},
		{fullPath: `/api/transfers/12345`, method: `GET`, routeRuleIndex: 1, match: true, routeParams: map[string]string{"uniqueID": "12345"}},
		{fullPath: `/api/entity/45/reference/xyz`, method: `GET`, routeRuleIndex: 2, match: true, routeParams: map[string]string{"id": "45", "ref": "xyz"}},
		{fullPath: `/api/transfers`, method: `GET`, routeRuleIndex: 3, match: true, routeParams: nil},
		{fullPath: `/api/transfers`, method: `POST`, routeRuleIndex: 3, match: false, routeParams: nil},
		{fullPath: `/api/transfers`, method: `POST`, routeRuleIndex: 4, match: true, routeParams: nil},
	}

	for _, td := range data {
		req := toHttpRequest(td.method, td.fullPath)

		foundRouteRule := router.FindMatch(req)

		if td.match {
			assert.Equal(t, routeRules[td.routeRuleIndex], foundRouteRule)
		} else {
			assert.NotEqual(t, routeRules[td.routeRuleIndex], foundRouteRule)
		}

		foundRuleRouteParams := foundRouteRule.GetRouteParams()

		for k, v := range foundRuleRouteParams {
			assert.Equal(t, v, td.routeParams[k])
		}
	}

}

func toHttpRequest(method, path string) *http.Request {
	split := strings.Split(path, "?")
	queryStrippedPath := split[0]
	query := ""

	if len(split) > 1 {
		query = split[1]
	}

	return &http.Request{
		Method: method,
		URL: &url.URL{
			Scheme:      "",
			Opaque:      "",
			User:        &url.Userinfo{},
			Host:        "",
			Path:        queryStrippedPath,
			RawPath:     "",
			ForceQuery:  false,
			RawQuery:    query,
			Fragment:    "",
			RawFragment: "",
		},
		RequestURI: path,
	}
}
