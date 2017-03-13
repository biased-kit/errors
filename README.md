# Errors
IMO errors should serve three purposes:
 - show message for user, describing what's wrong in order to he'll be able to fix it.
 - show machines the kind of error, in order to be handled automatically.
 - show sufficient information to developers for easy debug including place (stack trace) and context (variables), this info is placed to log.

this package provides interface and implementation to solve each goal.

## What means error handling ?
* Most common case is just to transmit an error to upper caller.
* Though sometimes you want to add more context for user (change message), probably you could add context for debug as well.
* In rare case you could try to perform some action. For this purpose you may rely on error type since different errors could be handled differently. 
* Finally, you could want to raise *your own* type of error, in order to caller be able to handle it automatically.

## Error type
Error types should be described as a part of method interface. This way a caller could chose a right action for error handling automaticaly.
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
If you develop a program just use errors.E for every func declaration, like:
```go
func DoSomthing() errors.E
```
But if you develop a package with public API, I suggest your methods return standard error interface, and document that it could be safely casted to errors.E.
This allows you to safely vendor the errors package.
Although you API could expose errors.E but don't vendor the errors package in that case!
