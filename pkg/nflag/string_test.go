// string_test.go - Unit tests for string flags
// SPDX-License-Identifier: GPL-3.0-or-later

package nflag

import "testing"

// TestStringFlagLongShortShareSameValue verifies that when we create a flag
// with both long and short names, both flags share the same underlying Value.
func TestStringFlagLongShortShareSameValue(t *testing.T) {
	fset := NewFlagSet("test", ContinueOnError)

	// Add a flag with both long and short names
	fset.StringFlag("output", 'o', "output file")

	// Lookup both the long and short flags
	longFlag, ok := fset.LookupFlagLong("output")
	if !ok {
		t.Fatal("long flag not found")
	}

	shortFlag, ok := fset.LookupFlagShort('o')
	if !ok {
		t.Fatal("short flag not found")
	}

	// Verify they share the same Value (same pointer)
	if longFlag.Value != shortFlag.Value {
		t.Fatalf("long and short flags should share the same Value, but they don't: %p vs %p",
			longFlag.Value, shortFlag.Value)
	}
}

// TestStringFlagModified tests the Modified() method behavior in various scenarios.
func TestStringFlagModified(t *testing.T) {
	t.Run("ModifiedWithShortFlag", func(t *testing.T) {
		fset := NewFlagSet("test", ContinueOnError)

		// Add a flag with both long and short names
		value := fset.StringFlag("output", 'o', "output file")

		// Before parsing, Modified should be false
		longFlag, _ := fset.LookupFlagLong("output")
		if longFlag.Value.Modified() {
			t.Fatal("flag should not be modified before parsing")
		}

		// Parse with the flag set
		err := fset.Parse([]string{"-o", "file.txt"})
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		// After parsing, Modified should be true
		if !longFlag.Value.Modified() {
			t.Fatal("flag should be modified after parsing")
		}

		// The value should be set correctly
		if *value != "file.txt" {
			t.Fatalf("expected value 'file.txt', got %s", *value)
		}
	})

	t.Run("ModifiedWithLongFlag", func(t *testing.T) {
		fset := NewFlagSet("test", ContinueOnError)

		// Add a flag with both long and short names
		fset.StringFlag("output", 'o', "output file")

		// Parse with the long flag
		err := fset.Parse([]string{"--output", "result.txt"})
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		// Check Modified on both long and short flags
		longFlag, _ := fset.LookupFlagLong("output")
		shortFlag, _ := fset.LookupFlagShort('o')

		if !longFlag.Value.Modified() {
			t.Fatal("long flag should be modified after parsing --output")
		}

		if !shortFlag.Value.Modified() {
			t.Fatal("short flag should be modified after parsing --output (they share the same Value)")
		}
	})

	t.Run("NotModifiedWhenNotParsed", func(t *testing.T) {
		fset := NewFlagSet("test", ContinueOnError)

		// Add two flags
		fset.StringFlag("output", 'o', "output file")
		fset.StringFlag("input", 'i', "input file")

		// Parse with only one flag
		err := fset.Parse([]string{"-o", "file.txt"})
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		// The parsed flag should be modified
		outputFlag, _ := fset.LookupFlagLong("output")
		if !outputFlag.Value.Modified() {
			t.Fatal("output flag should be modified")
		}

		// The unparsed flag should NOT be modified
		inputFlag, _ := fset.LookupFlagLong("input")
		if inputFlag.Value.Modified() {
			t.Fatal("input flag should not be modified (was not in command line)")
		}
	})
}
