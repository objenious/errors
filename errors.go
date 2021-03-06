// Package errors provides simple error handling primitives.
//
// The traditional error handling idiom in Go is roughly akin to
//
//     if err != nil {
//             return err
//     }
//
// which when applied recursively up the call stack results in error reports
// without context or debugging information. The errors package allows
// programmers to add context to the failure path in their code in a way
// that does not destroy the original value of the error.
//
// Adding context to an error
//
// The errors.Wrap function returns a new error that adds context to the
// original error by recording a stack trace at the point Wrap is called,
// together with the supplied message. For example
//
//     _, err := ioutil.ReadAll(r)
//     if err != nil {
//             return errors.Wrap(err, "read failed")
//     }
//
// If additional control is required, the errors.WithStack and
// errors.WithMessage functions destructure errors.Wrap into its component
// operations: annotating an error with a stack trace and with a message,
// respectively.
//
// Retrieving the cause of an error
//
// Using errors.Wrap constructs a stack of errors compatible with Go 1.13 errors wrapping
//
// It can be inspected by errors.Cause. errors.Cause will recursively retrieve
// the topmost error that does not implement causer, which is assumed to be
// the original cause. For example:
//
//     switch err := errors.Cause(err).(type) {
//     case *MyError:
//             // handle specifically
//     default:
//             // unknown error
//     }
//
// Formatted printing of errors
//
// All error values returned from this package implement fmt.Formatter and can
// be formatted by the fmt package. The following verbs are supported:
//
//     %s    print the error. If the error has a Cause it will be
//           printed recursively.
//     %v    see %s
//     %+v   extended format. Each Frame of the error's StackTrace will
//           be printed in detail.
//
// Retrieving the stack trace of an error or wrapper
//
// New, Errorf, Wrap, and Wrapf record a stack trace at the point they are
// invoked. This information can be retrieved with the following interface:
//
//     type stackTracer interface {
//             StackTrace() errors.StackTrace
//     }
//
// The returned errors.StackTrace type is defined as
//
//     type StackTrace []Frame
//
// The Frame type represents a call site in the stack trace. Frame supports
// the fmt.Formatter interface that can be used for printing information about
// the stack trace of this error. For example:
//
//     if err, ok := err.(stackTracer); ok {
//             for _, f := range err.StackTrace() {
//                     fmt.Printf("%+s:%d\n", f, f)
//             }
//     }
//
// Although the stackTracer interface is not exported by this package, it is
// considered a part of its stable public interface.
//
// See the documentation for Frame.Format for more details.
package errors

import (
	goerrors "errors"
	"fmt"
	"io"
)

// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
func New(message string) error {
	return &withStack{
		goerrors.New(message),
		callers(),
		"",
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
func Errorf(format string, args ...interface{}) error {
	return &withStack{
		fmt.Errorf(format, args...),
		callers(),
		"",
	}
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{
		err,
		callers(),
		"",
	}
}

type withStack struct {
	error
	*stack
	msg string
}

// Unwrap unwraps one level of this error
func (w *withStack) Unwrap() error {
	return goerrors.Unwrap(w.error)
}

// Cause is the same as Unwrap, returns the cause of this error
func (w *withStack) Cause() error {
	return w.Unwrap()
}

// Format formats the error with stack trace if available
func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if w.msg == "" {
				br := ""
				if w, ok := s.Width(); ok && w > 0 {
					br = "\n"
				}
				_, _ = fmt.Fprintf(s, "%s%+v", br, w.error) // recursive : go to bottom
			} else {
				if cause := w.Cause(); cause != nil {
					_, _ = fmt.Fprintf(s, "%+v", cause) // recursive : go to bottom

					if causeWithStack, ok := cause.(*withStack); ok && w.msg != "" && causeWithStack.msg != w.msg || w.msg != "" && cause.Error() != w.msg {
						_, _ = fmt.Fprintf(s, "\n%+v", w.msg)
					}
				} else {
					// root format
					_, _ = fmt.Fprintf(s, "%+v", w.error)
				}
			}
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, w.error.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	err = fmt.Errorf("%s: %w", message, err)
	return &withStack{
		err,
		callers(),
		message,
	}
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	msg := fmt.Sprintf(format, args...)
	args = append(args, err)
	err = fmt.Errorf(format+": %w", args...)
	return &withStack{
		err,
		callers(),
		msg,
	}
}

// Cause returns the underlying cause of the error, if possible.
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
func Cause(err error) error {
	for {
		unwrap := goerrors.Unwrap(err)
		if unwrap == nil {
			if wrap, ok := err.(*withStack); ok {
				return wrap.error
			}
			return err
		}
		err = unwrap
	}
}
