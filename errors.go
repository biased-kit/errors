package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// E represents an error
type E interface {
	error
	GetStack() *runtime.Frames
	GetKeyval() []interface{}
}

// Error type implements the E interface.
// It is public for serialization purpose only, in most cases you need to use E.
type Error struct {
	Err    error // when wrap standard error, preserve it for future use. Currently it's not used.
	Msg    string
	Keyval []interface{}
	Stack  *runtime.Frames
}

// Error produce error messafe for user
func (e *Error) Error() string {
	return e.Msg
}

// GetStack returns stack trace in order to log it.
func (e *Error) GetStack() *runtime.Frames {
	return e.Stack
}

// GetKeyval return key/value pairs describing error context. valuable for debuging purpose.
func (e *Error) GetKeyval() []interface{} {
	return e.Keyval
}

// New creates a new error
// values for pattern is obtained from key/values pairs, at that the keys are ignored.
func New(pattern string, keyval ...interface{}) E {
	msg := foramtMsg(pattern, keyval)
	st := frames()
	return &Error{
		Msg:    msg,
		Keyval: keyval,
		Stack:  st,
	}
}

// Wrap converts standard error to errors.E.
// In case if err is already type of errors.E the original stack trace and keyval is preserved.
// If you wrap nil you'll get nil.
func Wrap(err error, keyval ...interface{}) E {
	if err == nil {
		return nil
	}

	st := frames()
	msg := err.Error()
	if er, ok := err.(E); ok {
		for _, v := range er.GetKeyval() {
			keyval = append(keyval, v)
		}
		st = er.GetStack()
	}

	return &Error{
		Err:    err,
		Msg:    msg,
		Keyval: keyval,
		Stack:  st,
	}
}

func foramtMsg(msg string, keyval []interface{}) string {
	// todo: bench strings.Count and replace with custom implementation if it faster
	n := strings.Count(msg, "%")
	nn := strings.Count(msg, "%%")
	n -= nn * 2
	if n == 0 {
		return msg
	}
	vals := make([]interface{}, 0, n)
	for i := 1; i < len(keyval) && i < n*2; i += 2 {
		vals = append(vals, keyval[i])
	}
	return fmt.Sprintf(msg, vals...)
}

func frames() *runtime.Frames {
	var rpc [32]uintptr
	runtime.Callers(3, rpc[:])
	frames := runtime.CallersFrames(rpc[:])
	return frames
}
