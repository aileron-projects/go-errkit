package errkit_test

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/aileron-projects/go-errkit"
	"github.com/aileron-projects/go-tester"
)

func TestUnwrapErr(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		err := errkit.UnwrapErr(nil)
		tester.AssertEqual(t, nil, err)
	})
	t.Run("errorf can be unwrapped", func(t *testing.T) {
		err := fmt.Errorf("outer error [%w]", io.EOF) // Implements `Unwrap() error`
		err = errkit.UnwrapErr(err)
		tester.AssertEqual(t, io.EOF, err)
	})
	t.Run("join cannet be unwrapped", func(t *testing.T) {
		err := errkit.UnwrapErr(errors.Join(io.EOF, io.EOF)) // Implements `Unwrap() []error`
		tester.AssertEqual(t, nil, err)
	})
}

func TestUnwrapErrs(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		errs := errkit.UnwrapErrs(nil)
		tester.AssertDeepEqual(t, nil, errs)
	})
	t.Run("errorf cannot be unwrapped", func(t *testing.T) {
		err := fmt.Errorf("outer error [%w]", io.EOF) // Implements `Unwrap() error`
		errs := errkit.UnwrapErrs(err)
		tester.AssertDeepEqual(t, nil, errs)
	})
	t.Run("join can be unwrapped", func(t *testing.T) {
		errs := errkit.UnwrapErrs(errors.Join(io.EOF, io.ErrUnexpectedEOF)) // Implements `Unwrap() []error`
		tester.AssertEqual(t, 2, len(errs))
		tester.AssertDeepEqual(t, []error{io.EOF, io.ErrUnexpectedEOF}, errs)
	})
}
