# Errors
IMO error should serve three purposes:
 - shows a message for an user, describing what's wrong in order to the user be able to fix it.
 - provides the kind of error for machines, in order to be handled automatically.
 - shows sufficient information to developers for easy debug. The info includes place (stack trace) and context (variables) where the error happened.

this package provides interface and implementation to solve each goal.

## What means error handling ?
* Most common case is just to transmit an error to upper caller.
* Though sometimes you want to add more context for debug.
* In rare case you could try to perform some action. For this purpose you may rely on error type since different errors could be handled differently. 
* Finally, you could want to raise *your own* type of error, in order to caller be able to handle it automatically.

## Error type
Error types should be described as a part of method interface. This way a caller could chose a right action for error handling automatically.
I suggest to use golang struct with embedded errors.E interface as a error type, like:
```go
type ErrNotFound struct {
    errors.E
    some_field some_type // may be used for handling action.
}
```
This way the caller could use type switching to detect the right error type.

## Error interface type
If you define a public API as an interface, the errors it could produce should be defined as well.
for example:
```go
type ErrBadArgument interface {
    errors.E
    // add interface marker
    IsErrBadArgument()
    // or in specific case you could add useful method:
    GetBadArgName() string
}
```

## Error propagation
Just use errors.E for every func declaration, like:
```go
func DoSomthing() errors.E
```

inside you code use *err* for standard error interface, and *e* for the errors.E, this way you'll avoid misunderstanding:
```go
value, err := standardFunc()
if err != nil {
    errors.Wrap(err)
}

result, e := myFunc(value)
if e != nil {
    return e
}
```

## Inspired by github.com/pkg/errors