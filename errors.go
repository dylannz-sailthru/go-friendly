package errors

import "errors"

type Error struct {
	error     error
	userError error
}

func (e Error) Error() string {
	return e.error.Error()
}

func (e Error) UserError() string {
	return e.userError.Error()
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

func NewError(user, err error) error {
	if user == nil {
		return nil
	}

	if err == nil {
		err = user
	}

	return &Error{
		error:     err,
		userError: user,
	}
}

func NewErrorString(user, err string) error {
	if user == "" {
		return nil
	}

	if err == "" {
		err = user
	}

	return &Error{
		error:     errors.New(err),
		userError: errors.New(user),
	}
}
