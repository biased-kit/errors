package errors

import (
	"fmt"
	"runtime"
)

// E represents an error
type E interface {
	error
	// StackTrace returns runtime.Frames , that represent execution stack
	StackTrace() *runtime.Frames
	// KeyValues return key/value pairs that represent user-added params.
	KeyValues() []interface{}
	// With allows to add key/values pairs it return E in order to support "fluent interface"
	// if count is not even the nil value is added to the end.
	With(keyvals ...interface{}) E
}

// error type implements the E interface.
type erro struct {
	Err     error // when wrap standard error, preserve it for future use. Currently it's not used.
	Msg     string
	Keyvals []interface{}
	Stack   *runtime.Frames
}

// Error produce error messafe for user
func (e *erro) Error() string {
	return e.Msg
}

// GetStack returns stack trace in order to log it.
func (e *erro) StackTrace() *runtime.Frames {
	return e.Stack
}

// GetKeyvals return key/value pairs describing error context. valuable for debuging purpose.
func (e *erro) KeyValues() []interface{} {
	return e.Keyvals
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
		Msg:   msg,
		Stack: frames(lvl),
	}
}

// New creates a new error
func New(msg string) E {
	return create(msg, 4)
}

// New creates a new error from format string and arguments
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
	if er, ok := err.(E); ok {
		return er
	}

	return create(err.Error(), 4)
}

// WrapWith is similar to Wrap except it also could add key/values pairs.
func WrapWith(err error, keyvals ...interface{}) E {
	er := Wrap(err)
	if err != nil {
		er.With(keyvals...)
	}
	return er
}

func frames(lvl int) *runtime.Frames {
	var rpc [32]uintptr
	runtime.Callers(lvl, rpc[:])
	frames := runtime.CallersFrames(rpc[:])
	return frames
}
