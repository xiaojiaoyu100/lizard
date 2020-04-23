package errorkit

import "github.com/pkg/errors"

// WithStack wrap an error with stack.
func WithStack(err error) error {
	return errors.WithStack(err)
}
