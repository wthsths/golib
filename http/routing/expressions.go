package routing

import (
	"regexp"
	"strings"
)

// RegToRouteExp converts regular expression to route expression format.
//
// E.g.: '(?P<guid>\S+)' will convert to '{guid}'.
func RegToRouteExp(regex string) (string, error) {
	compiledReg, err := regexp.Compile(`((?:\(\?P<\w+\>\\S\+\))+),?`)
	if err != nil {
		return "", err
	}

	matches := compiledReg.FindAllString(regex, -1)
	routeex := regex

	for _, m := range matches {
		modified := strings.Replace(m, `(?P<`, `{`, 1)
		modified = strings.Replace(modified, `>\S+)`, `}`, 1)

		routeex = strings.Replace(routeex, m, modified, 1)
	}

	return routeex, nil
}

// RouteToRegExp converts route expression to regular expression format.
//
// E.g.: '{guid}' will convert to '(?P<guid>\S+)'.
func RouteToRegExp(routeex string) (string, error) {
	compiledReg, err := regexp.Compile(`((?:\{\w+\})+),?`)
	if err != nil {
		return "", err
	}

	matches := compiledReg.FindAllString(routeex, -1)
	regex := routeex

	for _, m := range matches {
		modified := strings.Replace(m, `{`, `(?P<`, 1)
		modified = strings.Replace(modified, `}`, `>\S+)`, 1)

		regex = strings.Replace(regex, m, modified, 1)
	}

	return regex, nil
}
