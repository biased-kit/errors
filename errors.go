package errors

import (
	stderrors "errors"
	"fmt"
	"runtime"
)

// E represents an error
type E interface {
	error
	// StackTrace returns program counters of function invocations on the calling goroutine's stack.
	// See runtime.Callers()
	StackTrace() []uintptr
	// KeyValues return key/value pairs that represent user-added params.
	KeyValues() []interface{}
	// With allows to add key/values pairs it return E in order to support "fluent interface"
	// if count is not even the nil value is added to the end.
	With(keyvals ...interface{}) E
	// Cause returns underlaying err
	Cause() error
}

// error type implements the E interface.
type erro struct {
	error
	Keyvals []interface{}
	Stack   []uintptr
}

// GetStack returns stack trace.
func (e *erro) StackTrace() []uintptr {
	return e.Stack
}

// GetKeyvals return key/value pairs describing error context. valuable for debuging purpose.
func (e *erro) KeyValues() []interface{} {
	return e.Keyvals
}

// Cause returns underlaying err
func (e *erro) Cause() error {
	return e.error
}

func (e *erro) With(keyvals ...interface{}) E {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, nil)
	}
	e.Keyvals = append(e.Keyvals, keyvals...)
	return e
}

// New creates a new error
func create(msg string, lvl int) E {
	return &erro{
		error: stderrors.New(msg),
		Stack: frames(lvl),
	}
}

// New creates a new error
func New(msg string) E {
	return create(msg, 4)
}

// Newf creates a new error from format string and arguments
func Newf(format string, args ...interface{}) E {
	msg := fmt.Sprintf(format, args...)
	return create(msg, 4)
}

// Wrap converts standard error to errors.E adding stack trace.
// In case if err is already type of errors.E it won't be changed.
// If you wrap nil you'll get nil.
func Wrap(err error) E {
	if err == nil {
		return nil
	}
	if e, ok := err.(E); ok {
		return e
	}

	return create(err.Error(), 4)
}

// WrapWith is similar to Wrap except it also could add key/values pairs.
func WrapWith(err error, keyvals ...interface{}) E {
	if err == nil {
		return nil
	}

	e, ok := err.(E)
	if !ok {
		e = create(err.Error(), 4)
	}

	e.With(keyvals...)
	return e
}

// Unwrap returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Unwrap(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

func frames(lvl int) []uintptr {
	rpc := make([]uintptr, 32)
	runtime.Callers(lvl, rpc)
	return rpc
}
