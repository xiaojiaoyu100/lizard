package apiaccessor

import (
	"errors"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func TestDefTimestampChecker(t *testing.T) {
	err := defTimestampChecker(1602146001) // 2020-10-08 16:33:21
	assert.Equal(t, errors.Is(err, ErrTimestampTimeout), true)
	t.Log(err)
	err = defTimestampChecker(time.Now().Unix())
	assert.Equal(t, err, nil)
}
