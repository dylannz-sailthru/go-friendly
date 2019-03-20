package errors

type Error struct {
	User  error
	Cause error
}

type ErrCauser interface {
	Error() string
	Cause() string
}

func (s Error) Error() string {
	return s.User.Error()
}

func Cause(err error) error {
	cause, ok := err.(*Error)
	if ok {
		return cause.Cause
	}

	return err
}

func User(err error) error {
	cause, ok := err.(*Error)
	if ok {
		return cause.User
	}

	return err
}

func NewError(user, cause error) error {
	return &Error{
		User:  user,
		Cause: cause,
	}
}
