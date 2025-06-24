// stdlib_test.go - standard library execution environment tests
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStdlibExecEnv(t *testing.T) {

	t.Run("Args", func(t *testing.T) {
		env := NewStdlibExecEnv()
		if diff := cmp.Diff(env.Args(), os.Args); diff != "" {
			t.Error(diff)
		}
	})

	t.Run("Exit", func(t *testing.T) {
		env := NewStdlibExecEnv()
		var code int
		env.OSExit = func(exitcode int) {
			code = exitcode
		}
		env.Exit(117)
		if code != 117 {
			t.Errorf("exit(%d)", code)
		}
	})

	t.Run("LookupEnv", func(t *testing.T) {
		env := NewStdlibExecEnv()
		env.OSLookupEnv = func(key string) (string, bool) {
			return "FOOBAR", true
		}
		got, _ := env.LookupEnv("FOO")
		if got != "FOOBAR" {
			t.Errorf("LookupEnv(%q) = %q, want %q", "FOO", got, "FOOBAR")
		}
	})

	t.Run("SignalNotify", func(t *testing.T) {
		env := NewStdlibExecEnv()
		var got chan<- os.Signal
		env.SignalNotifyFunc = func(c chan<- os.Signal, sig ...os.Signal) {
			got = c
		}
		ch := make(chan os.Signal, 1)
		env.SignalNotify(ch, os.Interrupt)
		if got != ch {
			t.Errorf("SignalNotify() = %v, want %v", got, ch)
		}
	})

	t.Run("Stderr", func(t *testing.T) {
		env := NewStdlibExecEnv()
		if env.Stderr() != os.Stderr {
			t.Errorf("Stderr() = %v, want %v", env.Stderr(), os.Stderr)
		}
	})

	t.Run("Stdout", func(t *testing.T) {
		env := NewStdlibExecEnv()
		if env.Stdout() != os.Stdout {
			t.Errorf("Stdout() = %v, want %v", env.Stdout(), os.Stdout)
		}
	})

	t.Run("Stdin", func(t *testing.T) {
		env := NewStdlibExecEnv()
		if env.Stdin() != os.Stdin {
			t.Errorf("Stdin() = %v, want %v", env.Stdin(), os.Stdin)
		}
	})
}
