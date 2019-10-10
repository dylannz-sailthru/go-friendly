package friendly

import (
	"errors"
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

// Friendly takes any error and will return the first user-friendly error it finds
// as it traverses up through the linked list. If there are no user-friendly
// causes found, nil is returned.
func Friendly(err error) error {
	type errCauser interface {
		Cause() error
	}
	type errFriendlyer interface {
		Friendly() error
	}

	for err != nil {
		user, ok := err.(errFriendlyer)
		if ok {
			return user.Friendly()
		}

		cause, ok := err.(errCauser)
		if !ok {
			break
		}
		err = cause.Cause()
	}

	return nil
}
