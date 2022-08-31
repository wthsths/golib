package gl_routing

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Router struct {
	// Contains static path definitions in mapping as follows: QueryStrippedPath -> Method -> *RouteRule
	staticPaths map[string]map[string]*RouteRule
	// Contains dynamic path definitions which have named route parameters in them.
	dynamicPaths []*RouteRule
	// Contains all route rules as key: path, value: method pairs.
	// Meant to be used for checking duplicates during initialization.
	allPaths map[string]string
}

// NewRouter creates http router from input routeRules.
//
// It will return error upon invalid data.
func NewRouter(routeRules []*RouteRule) (*Router, error) {
	router := &Router{
		staticPaths:  make(map[string]map[string]*RouteRule, len(routeRules)),
		dynamicPaths: make([]*RouteRule, 0, len(routeRules)),
		allPaths:     make(map[string]string, len(routeRules)),
	}

	for _, r := range routeRules {
		path, ok := router.allPaths[r.Path]
		if ok && path == r.Method {
			return nil, fmt.Errorf("path: '%s' is registered multiple times to method: '%s'", r.Path, r.Method)
		}
		router.allPaths[r.Path] = r.Method

		if !r.DynamicPath {
			if router.staticPaths[r.Path] == nil {
				router.staticPaths[r.Path] = make(map[string]*RouteRule)
			}

			router.staticPaths[r.Path][r.Method] = r
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
func (sr *Router) FindMatch(r *http.Request) *RouteRule {
	queryStrippedPath := strings.Split(r.URL.RequestURI(), "?")[0]
	staticPathRecord := sr.staticPaths[queryStrippedPath]
	if staticPathRecord != nil {
		staticRouteRule, ok := staticPathRecord[r.Method]
		if ok {
			return staticRouteRule
		}
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

			v.routeParams = result
			return v
		}
	}

	return nil
}

// HasMatch returns true if input request matches with any of the registered routed rules.
func (sr *Router) HasMatch(r *http.Request) bool {
	queryStrippedPath := strings.Split(r.URL.RequestURI(), "?")[0]
	staticPathRecord := sr.staticPaths[queryStrippedPath]
	if staticPathRecord != nil {
		_, ok := staticPathRecord[r.Method]
		if ok {
			return true
		}
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

	routeParams map[string]string
}

func (r *RouteRule) GetRouteParams() map[string]string {
	return r.routeParams
}
