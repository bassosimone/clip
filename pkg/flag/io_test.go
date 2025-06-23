// flagset_test.go - Tests for FlagSet implementation
// SPDX-License-Identifier: GPL-3.0-or-later

package flag

import (
	"io"
	"testing"
)

func TestSetStderr(t *testing.T) {
	fx := NewFlagSet("test", ContinueOnError)
	fx.SetStderr(io.Discard)
	if fx.stderr != io.Discard {
		t.Errorf("Expected stderr to be io.Discard, got %v", fx.stderr)
	}
}

func TestSetStdout(t *testing.T) {
	fx := NewFlagSet("test", ContinueOnError)
	fx.SetStdout(io.Discard)
	if fx.stdout != io.Discard {
		t.Errorf("Expected stdout to be io.Discard, got %v", fx.stdout)
	}
}
