package gl_routing

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Router struct {
	staticPaths  map[string]*RouteRule
	dynamicPaths []*RouteRule
	allPaths     map[string]bool
}

// NewRouter creates http router from input routeRules.
//
// It will return error upon invalid data.
func NewRouter(routeRules []*RouteRule) (*Router, error) {
	router := &Router{
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
			regexConv, err := RouteToRegExp(r.Path)
			if err != nil {
				return nil, fmt.Errorf("invalid path definition: '%s'", r.Path)
			}

			compiled, err := regexp.Compile(regexConv)
			if err != nil {
				return nil, fmt.Errorf("unable parse dynamic path: '%s' error: %s", r.Path, err.Error())
			}
			r.regex = compiled
			router.dynamicPaths = append(router.dynamicPaths, r)
		}

	}
	return router, nil
}

// FindMatch can be used inside a http.Handle() to check if incoming request matches with any of the routing rules.
// It returns routeTo func of the match.
// It also extracts and returns route parameters from curly bracket definitions.
//
// E.g: Input path: `/Transfer/{guid}`
//
// Request: `/Transfer/abcdef` will register as "guid"="abcdef" to routeParams.
func (sr *Router) FindMatch(r *http.Request) (authWith func(sessionID string, w http.ResponseWriter, r *http.Request) error, routeTo func(w http.ResponseWriter, r *http.Request, sessionID string, routeParams map[string]string), routeParams map[string]string) {
	queryStrippedPath := strings.Split(r.URL.RequestURI(), "?")[0]
	staticRoute := sr.staticPaths[queryStrippedPath]
	if staticRoute != nil && staticRoute.Method == r.Method {
		return staticRoute.AuthWith, staticRoute.RouteTo, nil
	}

	for _, v := range sr.dynamicPaths {
		matchDynamicPath := v.regex.MatchString(queryStrippedPath)
		if matchDynamicPath && v.Method == r.Method {
			result := make(map[string]string)
			match := v.regex.FindStringSubmatch(queryStrippedPath)
			for i, name := range v.regex.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}

			return v.AuthWith, v.RouteTo, result
		}
	}

	return nil, nil, nil
}

// HasMatch returns true if input request matches with any of the registered routed rules.
func (sr *Router) HasMatch(r *http.Request) bool {
	queryStrippedPath := strings.Split(r.URL.RequestURI(), "?")[0]
	staticRoute := sr.staticPaths[queryStrippedPath]
	if staticRoute != nil && staticRoute.Method == r.Method {
		return true
	}

	for _, v := range sr.dynamicPaths {
		matchDynamicPath := v.regex.MatchString(queryStrippedPath)
		if matchDynamicPath && v.Method == r.Method {
			result := make(map[string]string)
			match := v.regex.FindStringSubmatch(queryStrippedPath)
			for i, name := range v.regex.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}

			return true
		}
	}

	return false
}

// RouteRule is used for registering rules to Router.
// Any request path with route parameters in it should be registered with within curly brackets.
// They should also be registered as DynamicPath=true.
//
// Example path:  `/Transfer/{guid}`
//
// Query parameters in a url are ignored during checking.
// Therefore, request paths that have query parameters in it (but have no route parameters) should be registered as DynamicPath=false.
type RouteRule struct {
	Method string
	Path   string
	// DynamicPath should be set to true if endpoint has route parameters in it.
	// Query parameters however are NOT considered as a part of dynamic path.
	DynamicPath bool
	AuthWith    func(sessionID string, w http.ResponseWriter, r *http.Request) error
	RouteTo     func(w http.ResponseWriter, r *http.Request, sessionID string, routeParams map[string]string)

	regex *regexp.Regexp
}
