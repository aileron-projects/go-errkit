package errkit

// UnwrapErr returns the result of calling the Unwrap method on err
// if the given err implements Unwrap() that returns an error.
// Otherwise, UnwrapErr returns nil.
//
// UnwrapErr only calls a method of the form "Unwrap() error".
// In particular UnwrapErr does not unwrap errors returned by [errors.Join].
// See also [UnwrapErrs] and [errors.Unwrap].
func UnwrapErr(err error) error {
	u, ok := err.(interface{ Unwrap() error })
	if !ok {
		return nil
	}
	return u.Unwrap()
}

// UnwrapErrs returns the result of calling the Unwrap method on err
// if the given err implements Unwrap() that returns a []error.
// Otherwise, UnwrapErrs returns nil slice.
//
// UnwrapErrs only calls a method of the form "Unwrap() []error".
// UnwrapErrs can unwrap errors returned by [errors.Join].
// See also [UnwrapErr] and [errors.Unwrap].
func UnwrapErrs(err error) []error {
	u, ok := err.(interface{ Unwrap() []error })
	if !ok {
		return nil
	}
	return u.Unwrap()
}
