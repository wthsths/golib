package gl_routing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegExRouteExConversions(t *testing.T) {
	type testStruct struct {
		regex   string
		routeex string
	}

	testData := []testStruct{
		{
			regex:   `(?P<guid>\S+)`,
			routeex: `{guid}`,
		},
		{
			regex:   `abc/(?P<guid>\S+)/abc`,
			routeex: `abc/{guid}/abc`,
		},
		{
			regex:   `abc/(?P<param1>\S+)/abcd/(?P<param2>\S+)`,
			routeex: `abc/{param1}/abcd/{param2}`,
		},
		{
			regex:   `/api/transfers/(?P<id>\S+)/something/(?P<ref>\S+)`,
			routeex: `/api/transfers/{id}/something/{ref}`,
		},
	}

	// Route expression to regular expression conversions.
	for _, td := range testData {
		regex, err := RouteToRegExp(td.routeex)

		assert.NoError(t, err)
		assert.Equal(t, td.regex, regex)
	}

	// Regular expression to route expression conversions.
	for _, td := range testData {
		routeEx, err := RegToRouteExp(td.regex)

		assert.NoError(t, err)
		assert.Equal(t, td.routeex, routeEx)
	}
}
