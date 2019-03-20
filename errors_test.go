package errors_test

import (
	"errors"
	"log"
	"testing"

	. "github.com/dylannz-sailthru/errors"
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
		return NewError(ErrUser, err)
	}

	return nil
}

func TestConstructor(t *testing.T) {
	err := implementation()

	e, ok := err.(*Error)
	if !ok {
		log.Fatalf("expected err to be of type *Error, was: %T", err)
	}
	if e.User != ErrUser {
		log.Fatalf("expected err to be: %v, got: %v", ErrUser, e.User)
	}

	if Cause(err) != ErrSecret {
		log.Fatalf("expected cause to be: %v, got: %v", ErrSecret, Cause(err))
	}
}

func TestUserCauseWithThirdPartyError(t *testing.T) {
	c := errors.New("generic error")
	err := NewError(c, nil)

	e, ok := err.(*Error)
	if !ok {
		log.Fatalf("expected err to be of type *Error, was: %T", err)
	}
	if e.User != c {
		log.Fatalf("expected err to be: %v, got: %v", c, e.User)
	}

	if Cause(err) != c {
		log.Fatalf("expected cause to be: %v, got: %v", c, e.Cause)
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
