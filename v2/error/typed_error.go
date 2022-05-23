package gl_error

type TypedErrorCode int32

const (
	CodeUndefinedErr TypedErrorCode = 0
	CodeInternalErr  TypedErrorCode = 1
	CodeAuthErr      TypedErrorCode = 2
	CodeBadInput     TypedErrorCode = 3
)

// TypedError wraps a go error and an additional TypedErrorCode value for handling errors with higher precision.
//
// TypedErrorCode can have a user defined value or any of the constants provided in the package.
//
// E.g.: CodeInternalErr, CodeAuthErr etc...
type TypedError struct {
	err  error
	code TypedErrorCode
}

// NewTypedError creates a new error instance from given go error and TypedErrorCode value.
//
// TypedErrorCode can have a user defined value or any of the constants provided in the package.
//
// E.g.: CodeInternalErr, CodeAuthErr etc...
func NewTypedError(err error, code TypedErrorCode) TypedError {
	return TypedError{
		err:  err,
		code: code,
	}
}

// Code returns TypedErrorCode defined during initialization.
func (e TypedError) Code() TypedErrorCode {
	return e.code
}

// Error returns wrapped go error.
func (e TypedError) Error() error {
	return e.err
}
