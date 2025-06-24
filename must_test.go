// must_test.go - the must function tests
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import (
	"errors"
	"fmt"
	"testing"
)

func TestMust(t *testing.T) {
	t.Run("Must with nil error", func(t *testing.T) {
		env := NewStdlibExecEnv()
		Must(env, nil) // Should not exit
	})

	t.Run("Must with non-nil error", func(t *testing.T) {
		errMockExit := errors.New("mock exit")

		defer func() {
			err := recover().(error)
			if !errors.Is(err, errMockExit) {
				t.Errorf("expected exit error, got %v", err)
			}
		}()

		env := NewStdlibExecEnv()
		var exitcode int
		env.OSExit = func(code int) {
			exitcode = code
			panic(errMockExit)
		}
		Must(env, fmt.Errorf("test error"))

		if exitcode != 1 {
			t.Errorf("expected exit code 1, got %d", exitcode)
		}
	})
}
