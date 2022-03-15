package error

type TypedErrorCode int32

type typedError struct {
	err  error
	code TypedErrorCode
}

func NewTypedError(err error, code TypedErrorCode) *typedError {
	return &typedError{
		err:  err,
		code: code,
	}
}

func (e *typedError) Code() TypedErrorCode {
	return e.code
}

func (e *typedError) Error() error {
	return e.err
}
