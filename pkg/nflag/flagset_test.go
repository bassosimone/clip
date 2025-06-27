// flagset_test.go - FlagSet tests.
// SPDX-License-Identifier: GPL-3.0-or-later

package nflag

import (
	"errors"
	"testing"
)

func TestFlagSet_LookupFlagShort(t *testing.T) {
	fset := NewFlagSet("test", ContinueOnError)
	flag, ok := fset.LookupFlagShort('c')
	if ok {
		t.Errorf("expected no flag for short name 'c', got %v", flag)
	}
	if flag != nil {
		t.Errorf("expected nil flag for short name 'c', got %v", flag)
	}
}

func TestFlagSet_maybeHandleError(t *testing.T) {
	t.Run("PanicOnError", func(t *testing.T) {
		// make sure a panic happens
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic for PanicOnError, but did not panic")
			}
		}()

		// configure to panic on error
		fset := NewFlagSet("test", PanicOnError)

		// invoke with unknown flag
		fset.Parse([]string{"--unknown"})

		// make sure we do not reach this point
		t.Fatal("should have panicked but did not")
	})

	t.Run("ContinueOnError", func(t *testing.T) {
		// configure to continue on error
		fset := NewFlagSet("test", ContinueOnError)

		// add support for --help
		fset.AutoHelp("help", 'h', "Show this help message and exit.")

		// parse with --help
		err := fset.Parse([]string{"--help"})

		// make sure the error is [ErrHelp]
		if !errors.Is(err, ErrHelp) {
			t.Errorf("expected error to be %v, got %v", ErrHelp, err)
		}
	})
}
