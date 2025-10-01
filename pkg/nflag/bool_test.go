// bool_test.go - Unit tests for bool flags
// SPDX-License-Identifier: GPL-3.0-or-later

package nflag

import "testing"

// TestBoolFlagLongShortShareSameValue verifies that when we create a flag
// with both long and short names, both flags share the same underlying Value.
func TestBoolFlagLongShortShareSameValue(t *testing.T) {
	fset := NewFlagSet("test", ContinueOnError)

	// Add a flag with both long and short names
	fset.BoolFlag("verbose", 'v', "verbose mode")

	// Lookup both the long and short flags
	longFlag, ok := fset.LookupFlagLong("verbose")
	if !ok {
		t.Fatal("long flag not found")
	}

	shortFlag, ok := fset.LookupFlagShort('v')
	if !ok {
		t.Fatal("short flag not found")
	}

	// Verify they share the same Value (same pointer)
	if longFlag.Value != shortFlag.Value {
		t.Fatalf("long and short flags should share the same Value, but they don't: %p vs %p",
			longFlag.Value, shortFlag.Value)
	}
}

// TestBoolFlagModified tests the Modified() method behavior in various scenarios.
func TestBoolFlagModified(t *testing.T) {
	t.Run("ModifiedWithShortFlag", func(t *testing.T) {
		fset := NewFlagSet("test", ContinueOnError)

		// Add a flag with both long and short names
		value := fset.BoolFlag("verbose", 'v', "verbose mode")

		// Before parsing, Modified should be false
		longFlag, _ := fset.LookupFlagLong("verbose")
		if longFlag.Value.Modified() {
			t.Fatal("flag should not be modified before parsing")
		}

		// Parse with the flag set
		err := fset.Parse([]string{"-v"})
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		// After parsing, Modified should be true
		if !longFlag.Value.Modified() {
			t.Fatal("flag should be modified after parsing")
		}

		// The value should be set correctly
		if !*value {
			t.Fatal("expected value true, got false")
		}
	})

	t.Run("ModifiedWithLongFlag", func(t *testing.T) {
		fset := NewFlagSet("test", ContinueOnError)

		// Add a flag with both long and short names
		fset.BoolFlag("verbose", 'v', "verbose mode")

		// Parse with the long flag
		err := fset.Parse([]string{"--verbose"})
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		// Check Modified on both long and short flags
		longFlag, _ := fset.LookupFlagLong("verbose")
		shortFlag, _ := fset.LookupFlagShort('v')

		if !longFlag.Value.Modified() {
			t.Fatal("long flag should be modified after parsing --verbose")
		}

		if !shortFlag.Value.Modified() {
			t.Fatal("short flag should be modified after parsing --verbose (they share the same Value)")
		}
	})

	t.Run("NotModifiedWhenNotParsed", func(t *testing.T) {
		fset := NewFlagSet("test", ContinueOnError)

		// Add two flags
		fset.BoolFlag("verbose", 'v', "verbose mode")
		fset.BoolFlag("quiet", 'q', "quiet mode")

		// Parse with only one flag
		err := fset.Parse([]string{"-v"})
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		// The parsed flag should be modified
		verboseFlag, _ := fset.LookupFlagLong("verbose")
		if !verboseFlag.Value.Modified() {
			t.Fatal("verbose flag should be modified")
		}

		// The unparsed flag should NOT be modified
		quietFlag, _ := fset.LookupFlagLong("quiet")
		if quietFlag.Value.Modified() {
			t.Fatal("quiet flag should not be modified (was not in command line)")
		}
	})
}
