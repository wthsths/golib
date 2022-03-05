package http

import (
	"fmt"
	go_http "net/http"
	"regexp"
	"strings"
)

type router struct {
	staticPaths  map[string]*RouteRule
	dynamicPaths []*RouteRule
	allPaths     map[string]bool
}

// NewRouter creates http router from input routeRules.
// It will return error upon invalid data.
func NewRouter(routeRules []*RouteRule) (*router, error) {
	router := &router{
		staticPaths:  make(map[string]*RouteRule, len(routeRules)),
		dynamicPaths: make([]*RouteRule, 0, len(routeRules)),
		allPaths:     make(map[string]bool, len(routeRules)),
	}

	for _, r := range routeRules {
		_, ok := router.allPaths[r.Path]
		if ok {
			return nil, fmt.Errorf("path: '%s' is registered multiple times", r.Path)
		}
		router.allPaths[r.Path] = true

		if !r.DynamicPath {
			router.staticPaths[r.Path] = r
		} else {
			compiled, err := regexp.Compile(r.Path)
			if err != nil {
				return nil, fmt.Errorf("unable parse dynamic path: '%s' error: %w", r.Path, err)
			}
			r.regex = compiled
			router.dynamicPaths = append(router.dynamicPaths, r)
		}

	}
	return router, nil
}

// FindMatch can be used inside a http.Handle() to check if incoming request complies with any of the routing rules.
// It returns routeTo func of the match.
// It also extracts and returns route parameters from named regex groups.
//
// E.g: Input path: `/Transfer/(?P<guid>\S+)`
// Request: `/Transfer/abcdef` will register "guid"="abcdef" to routeParams.
func (rt *router) FindMatch(r *go_http.Request) (routeTo func(w go_http.ResponseWriter, r *go_http.Request, routeParams map[string]string), requiresAuth bool, routeParams map[string]string) {
	queryStrippedPath := strings.Split(r.URL.RequestURI(), "?")[0]
	staticRoute := rt.staticPaths[queryStrippedPath]
	if staticRoute != nil && staticRoute.Method == r.Method {
		return staticRoute.RouteTo, staticRoute.RequiresAuth, nil
	}

	for _, v := range rt.dynamicPaths {
		matchDynamicPath := v.regex.MatchString(queryStrippedPath)
		if matchDynamicPath && v.Method == r.Method {
			result := make(map[string]string)
			match := v.regex.FindStringSubmatch(queryStrippedPath)
			for i, name := range v.regex.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}

			return v.RouteTo, v.RequiresAuth, result
		}
	}

	return nil, false, nil
}

// RouteRule is used for registering rules to Router.
// Any request path with route parameters in it should be registered as named regex group.
// Also they should be registered as DynamicPath=true.
// Example path:  `/Transfer/(?P<guid>\S+)`
//
// Query parameters in a url are ignored during checking.
// Therefore, request paths that have only query parameters in it should be registered as DynamicPath=false.
type RouteRule struct {
	Method string
	Path   string
	// DynamicPath should be set to true if endpoint has route parameters in it.
	// Query parameters however are not considered as a part of dynamic path.
	DynamicPath  bool
	RequiresAuth bool
	RouteTo      func(w go_http.ResponseWriter, r *go_http.Request, routeParams map[string]string)
	regex        *regexp.Regexp
}
