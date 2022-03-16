package routing

import (
	"fmt"
	"regexp"
)

type RouteTable struct {
	routeRules []*RouteRule
}

// NewRouteTable checks validity of input routeRules.
//
// Note that rules for paths with route parameters must be defined with curly brackets.
//
// E.g: /Transfer/{guid}
func NewRouteTable(routeRules []*RouteRule) (*RouteTable, error) {
	table := &RouteTable{
		routeRules: routeRules,
	}

	for _, e := range table.routeRules {
		regexConv, err := RouteToRegExp(e.path)
		if err != nil {
			return nil, fmt.Errorf("invalid path definition: '%s'", e.path)
		}

		e.regexp, err = regexp.Compile(regexConv)
		if err != nil {
			return nil, fmt.Errorf("can not compile: '%s': %s", e.path, err.Error())
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
//
// Note that rules for paths with route parameters must be defined with curly brackets.
//
// E.g: /Transfer/{guid}
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
