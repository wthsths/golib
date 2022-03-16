package routing

import (
	"fmt"
	"regexp"
)

type RouteTable struct {
	routeRules []*ProxyRouteRule
}

// NewProxyRouteTable checks validity of input routeRules.
//
// Note that rules for paths with route parameters must be defined with curly brackets.
//
// E.g: /Transfer/{guid}
func NewProxyRouteTable(routeRules []*ProxyRouteRule) (*RouteTable, error) {
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

type ProxyRouteRule struct {
	method string
	path   string
	regexp *regexp.Regexp
}

// NewProxyRouteRule creates a single entry for RouteTable.
//
// Note that rules for paths with route parameters must be defined with curly brackets.
//
// E.g: /Transfer/{guid}
func NewProxyRouteRule(method, path string) *ProxyRouteRule {
	return &ProxyRouteRule{
		method: method,
		path:   path,
	}
}

func (rr *ProxyRouteRule) Method() string {
	return rr.method
}

func (rr *ProxyRouteRule) Path() string {
	return rr.path
}

func (rr *ProxyRouteRule) Regexp() regexp.Regexp {
	return *rr.regexp
}
