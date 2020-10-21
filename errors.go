package friendly

import (
	"errors"
	"fmt"

	"github.com/hashicorp/errwrap"
)

var DefaultUserError = errors.New("Something went wrong")

type Error struct {
	cause    error
	friendly error
}

func New() Error {
	return Error{
		cause:    DefaultUserError,
		friendly: DefaultUserError,
	}
}

// WithCause sets the internal non-user-safe cause of the error.
func (e Error) WithCause(err error) Error {
	e.cause = err
	return e
}

// WithFriendly sets the user-safe cause of the error.
func (e Error) WithFriendly(err error) Error {
	e.friendly = err
	return e
}

// WithCauseString sets the internal non-user-safe cause of the error.
func (e Error) WithCauseString(err string) Error {
	e.cause = errors.New(err)
	return e
}

// WithFriendlyString sets the user-safe cause of the error.
func (e Error) WithFriendlyString(err string) Error {
	e.friendly = errors.New(err)
	return e
}

// Err returns the friendly error as an 'error', if an underlying cause is
// present.
func (e Error) Err() error {
	if e.cause == nil {
		return nil
	}
	return e
}

// Cause returns the underlying cause error.
func (e Error) Cause() error {
	return e.cause
}

// Friendly returns the underlying friendly error.
func (e Error) Friendly() error {
	return e.friendly
}

func (e Error) Error() string {
	return e.cause.Error()
}

func unwrap(err error) []error {
	// support standard library Unwrap()
	if e, ok := err.(interface {
		Unwrap() error
	}); ok {
		if unwrapped := e.Unwrap(); e != nil {
			return []error{unwrapped}
		}
		return nil
	}

	// support legacy pkg/errors Cause()
	if e, ok := err.(interface {
		Cause() error
	}); ok {
		if unwrapped := e.Cause(); e != nil {
			return []error{unwrapped}
		}
		return nil
	}

	// support legacy hashicorp errwrap (before they added support for Unwrap())
	if e, ok := err.(errwrap.Wrapper); ok {
		return e.WrappedErrors()
	}

	return nil
}

func friendly(errs []error) error {
	type errFriendlyer interface {
		Friendly() error
	}

	nextErrs := []error{}
	for _, err := range errs {
		if err == nil {
			continue
		}

		if e, ok := err.(errFriendlyer); ok {
			return e.Friendly()
		}

		nextErrs = append(nextErrs, unwrap(err)...)
	}

	if len(nextErrs) > 0 {
		return friendly(nextErrs)
	}

	return nil
}

// Friendly takes any error and will return the first user-friendly error it finds
// as it traverses up through the linked list. If there are no user-friendly
// causes found, nil is returned.
func Friendly(err error) error {
	return friendly([]error{err})
}

// Wrap is a convience method to easily add a friendly message for an existing
// error.
func Wrap(cause error, friendly string) error {
	return New().WithCause(cause).WithFriendlyString(friendly).Err()
}

// Wrapf is a convience method to easily add a friendly message for an existing
// error.
func Wrapf(cause error, friendly string, a ...interface{}) error {
	return New().WithCause(cause).WithFriendlyString(fmt.Sprintf(friendly, a...)).Err()
}
