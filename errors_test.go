package friendly_test

import (
	"log"
	"testing"

	. "github.com/dylannz-sailthru/go-friendly"
	"github.com/pkg/errors"
)

var (
	ErrCause    = errors.New("something secret went wrong")
	ErrFriendly = errors.New("something public went wrong")
)

func someErr() error {
	return ErrCause
}

func implementation() error {
	err := someErr()
	if err != nil {
		return New().WithCause(err).WithFriendly(ErrFriendly)
	}

	return nil
}

func TestConstructor(t *testing.T) {
	err := implementation()

	e, ok := err.(Error)
	if !ok {
		log.Fatalf("expected err to be of type *Error, was: %T", err)
	}
	if e.Friendly() != ErrFriendly {
		log.Fatalf("expected user error to be: %v, got: %v", ErrFriendly.Error(), e.Friendly())
	}

	if err.Error() != ErrCause.Error() {
		log.Fatalf("expected cause to be: %v, got: %v", ErrCause.Error(), err.Error())
	}
}

func TestUserCauseWithThirdPartyError(t *testing.T) {
	c := errors.New("generic error")
	err := New().WithCause(c).Err()

	e, ok := err.(Error)
	if !ok {
		log.Fatalf("expected err to be of type *Error, was: %T", err)
	}
	if e.Friendly() != DefaultUserError {
		log.Fatalf("expected user error to be: %v, got: %v", DefaultUserError, e.Friendly())
	}

	if err.Error() != c.Error() {
		log.Fatalf("expected cause to be: %v, got: %v", c.Error(), err.Error())
	}
}

func TestConstructorWithNilCauseAndNilFriendly(t *testing.T) {
	err := New().WithCause(nil).WithFriendly(nil).Err()
	if err != nil {
		log.Fatalf("expected err to be nil, got: %v", err)
	}
}

func TestConstructorWithNilCause(t *testing.T) {
	err := New().WithCause(nil).WithFriendly(errors.New("some error")).Err()
	if err != nil {
		log.Fatalf("expected err to be nil, got: %v", err)
	}
}

func TestUser(t *testing.T) {
	err := New().WithCause(ErrCause).WithFriendly(ErrFriendly).Err()
	wrapped := errors.Wrap(err, "some wrapper")

	if ErrFriendly != Friendly(wrapped) {
		log.Fatalf("expected unwrapped user error (%v) to equal original user err (%v)", Friendly(wrapped), ErrFriendly)
	}
}

func TestUserWithNonUserError(t *testing.T) {
	err := errors.New("some non-user error")
	wrapped := errors.Wrap(err, "some wrapper")

	if Friendly(wrapped) != nil {
		log.Fatalf("expected unwrapped user error (%v) to be nil", Friendly(wrapped))
	}
}
