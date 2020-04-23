package lockguard

import "errors"

// Error error
type Error string

const (
	errLockNotObtained = Error("lock not obtained")
)

// Error reports an error.
func (e Error) Error() string {
	return string(e)
}

// IsLockNotObtained reports a lock which is not obtained.
func IsLockNotObtained(err error) bool {
	return errors.Is(err, errLockNotObtained)
}
