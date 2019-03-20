package errors

type Error struct {
	User  error
	Cause error
}

type ErrCauser interface {
	Error() string
	Cause() string
}

// Error calls the user .Error method. The design is to return the user-safe
// error by default.
func (s Error) Error() string {
	return s.User.Error()
}

// Cause takes any error and will return the underlying cause (if it's an
// instance of the Error struct in this package).
//
// Could be an improvement to cast to an interface instead of a concrete type,
// would allow for more flexibility.
func Cause(err error) error {
	e, ok := err.(*Error)
	if ok {
		return e.Cause
	}

	return err
}

func NewError(user, cause error) error {
	if user == nil {
		return nil
	}

	if cause == nil {
		cause = user
	}

	return &Error{
		User:  user,
		Cause: cause,
	}
}
