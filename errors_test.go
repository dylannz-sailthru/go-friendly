package errors_test

import (
	"log"
	"testing"

	. "github.com/dylannz-sailthru/errors"
	"github.com/pkg/errors"
)

var (
	ErrSecret = errors.New("something secret went wrong")
	ErrUser   = errors.New("something public went wrong")
)

func someErr() error {
	return ErrSecret
}

func implementation() error {
	err := someErr()
	if err != nil {
		return NewError(err, ErrUser)
	}

	return nil
}

func TestConstructor(t *testing.T) {
	err := implementation()

	e, ok := err.(*Error)
	if !ok {
		log.Fatalf("expected err to be of type *Error, was: %T", err)
	}
	if e.UserError() != ErrUser {
		log.Fatalf("expected user error to be: %v, got: %v", ErrUser.Error(), e.UserError())
	}

	if err.Error() != ErrSecret.Error() {
		log.Fatalf("expected cause to be: %v, got: %v", ErrSecret.Error(), err.Error())
	}
}

func TestUserCauseWithThirdPartyError(t *testing.T) {
	c := errors.New("generic error")
	err := NewError(c, nil)

	e, ok := err.(*Error)
	if !ok {
		log.Fatalf("expected err to be of type *Error, was: %T", err)
	}
	if e.UserError() != DefaultUserError {
		log.Fatalf("expected user error to be: %v, got: %v", DefaultUserError, e.UserError())
	}

	if err.Error() != c.Error() {
		log.Fatalf("expected cause to be: %v, got: %v", c.Error(), err.Error())
	}
}

func TestConstructorWithNilUserError(t *testing.T) {
	err := NewError(nil, nil)
	if err != nil {
		log.Fatalf("expected err to be nil, got: %v", err)
	}
	err = NewError(nil, errors.New("some error"))
	if err != nil {
		log.Fatalf("expected err to be nil, got: %v", err)
	}
}

func TestUser(t *testing.T) {
	err := NewError(ErrSecret, ErrUser)
	wrapped := errors.Wrap(err, "some wrapper")

	if ErrUser != User(wrapped) {
		log.Fatalf("expected unwrapped user error (%v) to equal original user err (%v)", User(wrapped), ErrUser)
	}
}

func TestUserWithNonUserError(t *testing.T) {
	err := errors.New("some non-user error")
	wrapped := errors.Wrap(err, "some wrapper")

	if User(wrapped) != nil {
		log.Fatalf("expected unwrapped user error (%v) to be nil", User(wrapped))
	}
}
