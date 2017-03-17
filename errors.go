package errors

import (
	"fmt"
	"runtime"
)

// E represents an error
type E interface {
	error
	StackTrace() *runtime.Frames
	KeyvalPairs() []interface{}
	With(keyvals ...interface{}) E
}

// Error type implements the E interface.
// It is public for serialization purpose only, in most cases you need to use E.
type Error struct {
	Err     error // when wrap standard error, preserve it for future use. Currently it's not used.
	Msg     string
	Keyvals []interface{}
	Stack   *runtime.Frames
}

// Error produce error messafe for user
func (e *Error) Error() string {
	return e.Msg
}

// GetStack returns stack trace in order to log it.
func (e *Error) StackTrace() *runtime.Frames {
	return e.Stack
}

// GetKeyvals return key/value pairs describing error context. valuable for debuging purpose.
func (e *Error) KeyvalPairs() []interface{} {
	return e.Keyvals
}

func (e *Error) With(keyvals ...interface{}) E {
	e.Keyvals = append(e.Keyvals, keyvals...)
	return e
}

// New creates a new error
// pattern and args means the same as for fmt.Printf()
func New(pattern string, args ...interface{}) E {
	msg := fmt.Sprintf(pattern, args...)
	st := frames()
	return &Error{
		Msg:   msg,
		Stack: st,
	}
}

// Wrap converts standard error to errors.E.
// In case if err is already type of errors.E the original stack trace and keyvals is preserved.
// If you wrap nil you'll get nil.
func Wrap(err error, keyvals ...interface{}) E {
	if err == nil {
		return nil
	}

	st := frames()
	msg := err.Error()
	if er, ok := err.(E); ok {
		for _, v := range er.KeyvalPairs() {
			keyvals = append(keyvals, v)
		}
		st = er.StackTrace()
	}

	return &Error{
		Err:     err,
		Msg:     msg,
		Keyvals: keyvals,
		Stack:   st,
	}
}

func frames() *runtime.Frames {
	var rpc [32]uintptr
	runtime.Callers(3, rpc[:])
	frames := runtime.CallersFrames(rpc[:])
	return frames
}
