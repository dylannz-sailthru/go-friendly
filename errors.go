package errors

import (
	"errors"
)

var DefaultUserError = errors.New("Something went wrong.")

type Error struct {
	error     error
	userError error
}

func (e Error) Error() string {
	return e.error.Error()
}

func (e Error) UserError() error {
	return e.userError
}

// Cause takes any error and will return the first user-friendly error it finds
// as it traverses up through the linked list. If there are no user-friendly
// causes found, nil is returned.
func User(err error) error {
	type errCauser interface {
		Cause() error
	}
	type errUser interface {
		UserError() error
	}

	for err != nil {
		user, ok := err.(errUser)
		if ok {
			return user.UserError()
		}

		cause, ok := err.(errCauser)
		if !ok {
			break
		}
		err = cause.Cause()
	}

	return nil
}

func NewError(err, user error) error {
	if err == nil {
		return nil
	}

	if user == nil {
		user = DefaultUserError
	}

	return &Error{
		error:     err,
		userError: user,
	}
}

func NewErrorString(err, user string) error {
	if err == "" {
		return nil
	}

	return &Error{
		error:     errors.New(err),
		userError: errors.New(user),
	}
}
