package gl_error

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Pass_By_Value(t *testing.T) {
	someErr := NewTypedError(fmt.Errorf("this is an error"), 1)
	modErr := someFunc(someErr)

	assert.Equal(t, TypedErrorCode(9), modErr.code)
	assert.Equal(t, "modified error", modErr.err.Error())
	assert.Equal(t, TypedErrorCode(1), someErr.code)
	assert.Equal(t, "this is an error", someErr.err.Error())
}

func someFunc(e TypedError) TypedError {
	e.code = 9
	e.err = fmt.Errorf("modified error")
	return e
}
