// int64_test.go - Unit tests for int64 flags
// SPDX-License-Identifier: GPL-3.0-or-later

package nflag

import "testing"

// TestInt64FlagLongShortShareSameValue verifies that when we create a flag
// with both long and short names, both flags share the same underlying Value.
func TestInt64FlagLongShortShareSameValue(t *testing.T) {
	fset := NewFlagSet("test", ContinueOnError)

	// Add a flag with both long and short names
	fset.Int64Flag("count", 'c', "count value")

	// Lookup both the long and short flags
	longFlag, ok := fset.LookupFlagLong("count")
	if !ok {
		t.Fatal("long flag not found")
	}

	shortFlag, ok := fset.LookupFlagShort('c')
	if !ok {
		t.Fatal("short flag not found")
	}

	// Verify they share the same Value (same pointer)
	if longFlag.Value != shortFlag.Value {
		t.Fatalf("long and short flags should share the same Value, but they don't: %p vs %p",
			longFlag.Value, shortFlag.Value)
	}
}

// TestInt64FlagModified tests the Modified() method behavior in various scenarios.
func TestInt64FlagModified(t *testing.T) {
	t.Run("ModifiedWithShortFlag", func(t *testing.T) {
		fset := NewFlagSet("test", ContinueOnError)

		// Add a flag with both long and short names
		value := fset.Int64Flag("count", 'c', "count value")

		// Before parsing, Modified should be false
		longFlag, _ := fset.LookupFlagLong("count")
		if longFlag.Value.Modified() {
			t.Fatal("flag should not be modified before parsing")
		}

		// Parse with the flag set
		err := fset.Parse([]string{"-c", "42"})
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		// After parsing, Modified should be true
		if !longFlag.Value.Modified() {
			t.Fatal("flag should be modified after parsing")
		}

		// The value should be set correctly
		if *value != 42 {
			t.Fatalf("expected value 42, got %d", *value)
		}
	})

	t.Run("ModifiedWithLongFlag", func(t *testing.T) {
		fset := NewFlagSet("test", ContinueOnError)

		// Add a flag with both long and short names
		fset.Int64Flag("count", 'c', "count value")

		// Parse with the long flag
		err := fset.Parse([]string{"--count", "100"})
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		// Check Modified on both long and short flags
		longFlag, _ := fset.LookupFlagLong("count")
		shortFlag, _ := fset.LookupFlagShort('c')

		if !longFlag.Value.Modified() {
			t.Fatal("long flag should be modified after parsing --count")
		}

		if !shortFlag.Value.Modified() {
			t.Fatal("short flag should be modified after parsing --count (they share the same Value)")
		}
	})

	t.Run("NotModifiedWhenNotParsed", func(t *testing.T) {
		fset := NewFlagSet("test", ContinueOnError)

		// Add two flags
		fset.Int64Flag("count", 'c', "count value")
		fset.Int64Flag("max", 'm', "max value")

		// Parse with only one flag
		err := fset.Parse([]string{"-c", "42"})
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		// The parsed flag should be modified
		countFlag, _ := fset.LookupFlagLong("count")
		if !countFlag.Value.Modified() {
			t.Fatal("count flag should be modified")
		}

		// The unparsed flag should NOT be modified
		maxFlag, _ := fset.LookupFlagLong("max")
		if maxFlag.Value.Modified() {
			t.Fatal("max flag should not be modified (was not in command line)")
		}
	})
}
