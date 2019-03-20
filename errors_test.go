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

	if User(err) != ErrUser {
		log.Fatalf("expected err to be: %v, got: %v", ErrUser, User(err))
	}

	if Cause(err) != ErrSecret {
		log.Fatalf("expected cause to be: %v, got: %v", ErrSecret, Cause(err))
	}
}

func TestUserCauseWithThirdPartyError(t *testing.T) {
	err := errors.New("generic error")

	if User(err) != err {
		log.Fatalf("expected err to be: %v, got: %v", err, User(err))
	}

	if Cause(err) != err {
		log.Fatalf("expected cause to be: %v, got: %v", err, User(err))
	}
}
