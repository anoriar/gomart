package errors

import "errors"

var ErrConflict = errors.New("conflict error")
var ErrNotFound = errors.New("not found error")
var ErrTooManyRequests = errors.New("too many requests")
var ErrInternalError = errors.New("internal error")
var ErrDependencyFailure = errors.New("dependency failure")
