package multierror

import (
	"errors"
	"fmt"
	"io"

	perr "github.com/pingcap/errors"
)

const (
	multiErrorStartString    = "==========multi error: print number %d started==========\n"
	multiErrorCompleteString = "\n==========multi error: print number %d completed==========\n"
)

// Error is an error type to track multiple errors. This is used to
// accumulate errors in cases and return them as a single "error".
type Error struct {
	Errs        []error
	ErrorFormat ErrorFormatFunc
}

func (e *Error) Error() string {
	fn := e.ErrorFormat
	if fn == nil {
		fn = ListFormatFunc
	}

	return fn(e.Errs)
}

// ErrorOrNil returns an error interface if this Error represents
// a list of errors, or returns nil if the list of errors is empty. This
// function is useful at the end of accumulation to make sure that the value
// returned represents the existence of errors.
func (e *Error) ErrorOrNil() error {
	if e == nil {
		return nil
	}
	if len(e.Errs) == 0 {
		return nil
	}

	return e
}

func (e *Error) GoString() string {
	return fmt.Sprintf("*%#v", *e)
}

// WrappedErrors returns the list of errors that this Error is wrapping. It is
// an implementation of the errwrap.Wrapper interface so that multierror.Error
// can be used with that library.
//
// This method is not safe to be called concurrently. Unlike accessing the
// Errs field directly, this function also checks if the multierror is nil to
// prevent a null-pointer panic. It satisfies the errwrap.Wrapper interface.
func (e *Error) WrappedErrors() []error {
	if e == nil {
		return nil
	}
	return e.Errs
}

// Unwrap returns an error from Error (or nil if there are no errors).
// This error returned will further support Unwrap to get the next error,
// etc. The order will match the order of Errs in the multierror.Error
// at the time of calling.
//
// The resulting error supports errors.As/Is/Unwrap so you can continue
// to use the stdlib errors package to introspect further.
//
// This will perform a shallow copy of the errors slice. Any errors appended
// to this error after calling Unwrap will not be available until a new
// Unwrap is called on the multierror.Error.
func (e *Error) Unwrap() error {
	// If we have no errors then we do nothing
	if e == nil || len(e.Errs) == 0 {
		return nil
	}

	// If we have exactly one error, we can just return that directly.
	if len(e.Errs) == 1 {
		return e.Errs[0]
	}

	// Shallow copy the slice
	errs := make([]error, len(e.Errs))
	copy(errs, e.Errs)
	return chain(errs)
}

// Errors returns the error list, it implements Errors interface of github.com/pingcap/errors
func (e *Error) Errors() []error {
	return e.Errs
}

// Format implements fmt.Formatter interface
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, fmt.Sprintf("multi error: %d errors occurred.\n", len(e.Errs)))

			for i, err := range e.WrappedErrors() {
				_, _ = io.WriteString(s, fmt.Sprintf(multiErrorStartString, i))

				_, _ = io.WriteString(s, err.Error())
				if perr.HasStack(err) {
					perr.GetStackTracer(err).StackTrace().Format(s, verb)
				}
				_, _ = io.WriteString(s, fmt.Sprintf(multiErrorCompleteString, i))
			}

			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	}
}

// chain implements the interfaces necessary for errors.Is/As/Unwrap to
// work in a deterministic way with multierror. A chain tracks a list of
// errors while accounting for the current represented error. This lets
// Is/As be meaningful.
//
// Unwrap returns the next error. In the cleanest form, Unwrap would return
// the wrapped error here but we can't do that if we want to properly
// get access to all the errors. Instead, users are recommended to use
// Is/As to get the correct error type out.
//
// Precondition: []error is non-empty (len > 0)
type chain []error

// Error implements the error interface
func (e chain) Error() string {
	return e[0].Error()
}

// Unwrap implements errors.Unwrap by returning the next error in the
// chain or nil if there are no more errors.
func (e chain) Unwrap() error {
	if len(e) == 1 {
		return nil
	}

	return e[1:]
}

// As implements errors.As by attempting to map to the current value.
func (e chain) As(target interface{}) bool {
	return errors.As(e[0], target)
}

// Is implements errors.Is by comparing the current value directly.
func (e chain) Is(target error) bool {
	return errors.Is(e[0], target)
}
