package routing

import (
	"fmt"
	"regexp"
)

type RouteTable struct {
	routeRules []*RouteRule
}

// NewRouteTable checks validity of input routeRules.
// Route rules must contain have regex compatible path values.
func NewRouteTable(routeRules []*RouteRule) (*RouteTable, error) {
	table := &RouteTable{
		routeRules: routeRules,
	}
	var err error
	for _, e := range table.routeRules {
		e.regexp, err = regexp.Compile(e.path)
		if err != nil {
			return nil, fmt.Errorf("can not compile: '%s': %w", e.path, err)
		}
	}
	return table, nil
}

type RouteRule struct {
	method string
	path   string
	regexp *regexp.Regexp
}

// NewRouteRule creates a single entry for RouteTable.
func NewRouteRule(method, path string) *RouteRule {
	return &RouteRule{
		method: method,
		path:   path,
	}
}

func (rr *RouteRule) Method() string {
	return rr.method
}

func (rr *RouteRule) Path() string {
	return rr.path
}

func (rr *RouteRule) Regexp() regexp.Regexp {
	return *rr.regexp
}
