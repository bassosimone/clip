// autohelp_test.go - Unit tests for help flags
// SPDX-License-Identifier: GPL-3.0-or-later

package nflag

import "testing"

// TestAutoHelpLongShortShareSameValue verifies that when we create a help flag
// with both long and short names, both flags share the same underlying Value.
func TestAutoHelpLongShortShareSameValue(t *testing.T) {
	fset := NewFlagSet("test", ContinueOnError)

	// Add a help flag with both long and short names
	fset.AutoHelp("help", 'h', "show help")

	// Lookup both the long and short flags
	longFlag, ok := fset.LookupFlagLong("help")
	if !ok {
		t.Fatal("long flag not found")
	}

	shortFlag, ok := fset.LookupFlagShort('h')
	if !ok {
		t.Fatal("short flag not found")
	}

	// Verify they share the same Value (same pointer)
	if longFlag.Value != shortFlag.Value {
		t.Fatalf("long and short flags should share the same Value, but they don't: %p vs %p",
			longFlag.Value, shortFlag.Value)
	}
}

// TestAutoHelpModified tests the Modified() method behavior in various scenarios.
func TestAutoHelpModified(t *testing.T) {
	t.Run("ModifiedWithShortFlag", func(t *testing.T) {
		fset := NewFlagSet("test", ContinueOnError)

		// Add a help flag with both long and short names
		fset.AutoHelp("help", 'h', "show help")

		// Before parsing, Modified should be false
		longFlag, _ := fset.LookupFlagLong("help")
		if longFlag.Value.Modified() {
			t.Fatal("flag should not be modified before parsing")
		}

		// Parse with the short flag (should return ErrHelp)
		err := fset.Parse([]string{"-h"})
		if err != ErrHelp {
			t.Fatalf("expected ErrHelp, got %v", err)
		}

		// After parsing, Modified should be true
		if !longFlag.Value.Modified() {
			t.Fatal("flag should be modified after parsing")
		}
	})

	t.Run("ModifiedWithLongFlag", func(t *testing.T) {
		fset := NewFlagSet("test", ContinueOnError)

		// Add a help flag with both long and short names
		fset.AutoHelp("help", 'h', "show help")

		// Parse with the long flag (should return ErrHelp)
		err := fset.Parse([]string{"--help"})
		if err != ErrHelp {
			t.Fatalf("expected ErrHelp, got %v", err)
		}

		// Check Modified on both long and short flags
		longFlag, _ := fset.LookupFlagLong("help")
		shortFlag, _ := fset.LookupFlagShort('h')

		if !longFlag.Value.Modified() {
			t.Fatal("long flag should be modified after parsing --help")
		}

		if !shortFlag.Value.Modified() {
			t.Fatal("short flag should be modified after parsing --help (they share the same Value)")
		}
	})

	t.Run("NotModifiedWhenNotParsed", func(t *testing.T) {
		fset := NewFlagSet("test", ContinueOnError)

		// Add help flag and another flag
		fset.AutoHelp("help", 'h', "show help")
		fset.BoolFlag("verbose", 'v', "verbose mode")

		// Parse with only the verbose flag
		err := fset.Parse([]string{"-v"})
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		// The verbose flag should be modified
		verboseFlag, _ := fset.LookupFlagLong("verbose")
		if !verboseFlag.Value.Modified() {
			t.Fatal("verbose flag should be modified")
		}

		// The help flag should NOT be modified
		helpFlag, _ := fset.LookupFlagLong("help")
		if helpFlag.Value.Modified() {
			t.Fatal("help flag should not be modified (was not in command line)")
		}
	})
}
