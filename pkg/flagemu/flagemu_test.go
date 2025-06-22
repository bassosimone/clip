// flagemu_test.go - Comprehensive unit tests for the flagemu package
// SPDX-License-Identifier: GPL-3.0-or-later

package flagemu

import (
	"testing"

	"github.com/bassosimone/clip/pkg/parser"
)

func TestNewFlagSet(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)
	if fs == nil {
		t.Fatal("NewFlagSet returned nil")
	}
	if fs.progname != "test" {
		t.Errorf("Expected progname 'test', got '%s'", fs.progname)
	}
	if fs.valuesShort == nil {
		t.Error("valuesShort map is nil")
	}
	if fs.valuesLong == nil {
		t.Error("valuesLong map is nil")
	}
	if len(fs.args) != 0 {
		t.Errorf("Expected empty args, got %v", fs.args)
	}
}

func TestFlagSet_Bool(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)

	// Test with both long and short names
	flag := fs.Bool("verbose", 'v', false, "verbose output")
	if flag == nil {
		t.Fatal("Bool returned nil pointer")
	}
	if *flag != false {
		t.Errorf("Expected default value false, got %t", *flag)
	}

	// Check that both long and short names are registered
	if _, exists := fs.valuesLong["verbose"]; !exists {
		t.Error("Long name 'verbose' not registered")
	}
	if _, exists := fs.valuesShort["v"]; !exists {
		t.Error("Short name 'v' not registered")
	}

	// Test with only long name
	flag2 := fs.Bool("debug", 0, true, "debug mode")
	if *flag2 != true {
		t.Errorf("Expected default value true, got %t", *flag2)
	}
	if _, exists := fs.valuesLong["debug"]; !exists {
		t.Error("Long name 'debug' not registered")
	}

	// Test with only short name
	_ = fs.Bool("", 'q', false, "quiet mode")
	if _, exists := fs.valuesShort["q"]; !exists {
		t.Error("Short name 'q' not registered")
	}
}

func TestFlagSet_String(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)

	// Test with both long and short names
	flag := fs.String("output", 'o', "default.txt", "output file")
	if flag == nil {
		t.Fatal("String returned nil pointer")
	}
	if *flag != "default.txt" {
		t.Errorf("Expected default value 'default.txt', got '%s'", *flag)
	}

	// Check that both long and short names are registered
	if _, exists := fs.valuesLong["output"]; !exists {
		t.Error("Long name 'output' not registered")
	}
	if _, exists := fs.valuesShort["o"]; !exists {
		t.Error("Short name 'o' not registered")
	}
}

func TestFlagSet_SetNoPermute(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)

	// Initially flags should be 0
	if fs.flags != 0 {
		t.Errorf("Expected initial flags to be 0, got %d", fs.flags)
	}

	fs.SetNoPermute()

	// After SetNoPermute, flags should be non-zero
	if fs.flags == 0 {
		t.Error("SetNoPermute did not set flags")
	}
}

func TestFlagSet_ParseBoolFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{"short flag", []string{"-v"}, true},
		{"long flag", []string{"--verbose"}, true},
		{"no flag", []string{}, false},
		{"bool with true", []string{"--verbose=true"}, true},
		{"bool with false", []string{"--verbose=false"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := NewFlagSet("test", ContinueOnError)
			verbose := fs.Bool("verbose", 'v', false, "verbose output")

			err := fs.Parse(tt.args)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			if *verbose != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, *verbose)
			}
		})
	}
}

func TestFlagSet_ParseStringFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"short flag separate", []string{"-o", "output.txt"}, "output.txt"},
		{"long flag separate", []string{"--output", "output.txt"}, "output.txt"},
		{"long flag equals", []string{"--output=output.txt"}, "output.txt"},
		{"no flag", []string{}, "default.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := NewFlagSet("test", ContinueOnError)
			output := fs.String("output", 'o', "default.txt", "output file")

			err := fs.Parse(tt.args)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			if *output != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, *output)
			}
		})
	}
}

func TestFlagSet_ParseMixedFlags(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)
	verbose := fs.Bool("verbose", 'v', false, "verbose output")
	output := fs.String("output", 'o', "", "output file")
	debug := fs.Bool("debug", 'd', false, "debug mode")

	args := []string{"-v", "--output=file.txt", "-d", "arg1", "arg2"}

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if !*verbose {
		t.Error("Expected verbose to be true")
	}
	if *output != "file.txt" {
		t.Errorf("Expected output 'file.txt', got '%s'", *output)
	}
	if !*debug {
		t.Error("Expected debug to be true")
	}

	positionalArgs := fs.Args()
	expected := []string{"arg1", "arg2"}
	if len(positionalArgs) != len(expected) {
		t.Errorf("Expected %d args, got %d", len(expected), len(positionalArgs))
	}
	for i, arg := range positionalArgs {
		if arg != expected[i] {
			t.Errorf("Expected arg[%d] '%s', got '%s'", i, expected[i], arg)
		}
	}
}

func TestFlagSet_ParseErrors(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{"unknown long flag", []string{"--unknown"}, true},
		{"unknown short flag", []string{"-x"}, true},
		{"valid flags", []string{"-v"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := NewFlagSet("test", ContinueOnError)
			fs.Bool("verbose", 'v', false, "verbose output")

			err := fs.Parse(tt.args)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestFlagSet_Args(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)
	fs.Bool("verbose", 'v', false, "verbose output")

	args := []string{"-v", "file1.txt", "file2.txt"}

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	positionalArgs := fs.Args()
	expected := []string{"file1.txt", "file2.txt"}

	if len(positionalArgs) != len(expected) {
		t.Errorf("Expected %d args, got %d", len(expected), len(positionalArgs))
	}

	for i, arg := range positionalArgs {
		if arg != expected[i] {
			t.Errorf("Expected arg[%d] '%s', got '%s'", i, expected[i], arg)
		}
	}
}

func TestBoolValue(t *testing.T) {
	var flag bool
	bv := newBoolValue(&flag)

	// Test OptionType
	if bv.OptionType() != parser.OptionTypeBool {
		t.Errorf("Expected OptionType to be OptionTypeBool")
	}

	// Test String representation
	if bv.String() != "false" {
		t.Errorf("Expected String 'false', got '%s'", bv.String())
	}

	// Test Set with "true"
	err := bv.Set("true")
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	if !flag {
		t.Error("Expected flag to be true after Set('true')")
	}
	if bv.String() != "true" {
		t.Errorf("Expected String 'true', got '%s'", bv.String())
	}

	// Test Set with "false"
	err = bv.Set("false")
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	if flag {
		t.Error("Expected flag to be false after Set('false')")
	}

	// Test Set with other values (should be false for unknown values)
	err = bv.Set("anything")
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	if flag {
		t.Error("Expected flag to be false after Set('anything') (unknown values default to false)")
	}
}

func TestBoolValue_SimpleParsing(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		// Only "true" is true
		{"true", "true", true},
		{"false", "false", false},

		// All other values are false
		{"True", "True", false},
		{"TRUE", "TRUE", false},
		{"1", "1", false},
		{"yes", "yes", false},
		{"on", "on", false},
		{"0", "0", false},
		{"no", "no", false},
		{"off", "off", false},
		{"unknown", "unknown", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var flag bool
			bv := newBoolValue(&flag)

			err := bv.Set(tt.value)
			if err != nil {
				t.Errorf("Set failed: %v", err)
			}

			if flag != tt.expected {
				t.Errorf("Expected %t for value '%s', got %t", tt.expected, tt.value, flag)
			}
		})
	}
}

func TestStringValue(t *testing.T) {
	var str string = "initial"
	sv := newStringValue(&str)

	// Test OptionType
	if sv.OptionType() != parser.OptionTypeString {
		t.Errorf("Expected OptionType to be OptionTypeString")
	}

	// Test String representation
	if sv.String() != "initial" {
		t.Errorf("Expected String 'initial', got '%s'", sv.String())
	}

	// Test Set
	err := sv.Set("new value")
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	if str != "new value" {
		t.Errorf("Expected string to be 'new value', got '%s'", str)
	}
	if sv.String() != "new value" {
		t.Errorf("Expected String 'new value', got '%s'", sv.String())
	}

	// Test Set with empty string
	err = sv.Set("")
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	if str != "" {
		t.Errorf("Expected string to be empty, got '%s'", str)
	}
}

func TestFlagSet_SeparatorHandling(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)
	verbose := fs.Bool("verbose", 'v', false, "verbose output")

	// Test with -- separator
	args := []string{"-v", "--", "--not-a-flag", "regular-arg"}

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if !*verbose {
		t.Error("Expected verbose to be true")
	}

	positionalArgs := fs.Args()
	expected := []string{"--not-a-flag", "regular-arg"}

	if len(positionalArgs) != len(expected) {
		t.Errorf("Expected %d args, got %d", len(expected), len(positionalArgs))
	}

	for i, arg := range positionalArgs {
		if arg != expected[i] {
			t.Errorf("Expected arg[%d] '%s', got '%s'", i, expected[i], arg)
		}
	}
}

func TestFlagSet_MultipleCallsToSameFlag(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)
	output := fs.String("output", 'o', "default", "output file")

	// Parse multiple times with same flag - last one should win
	args := []string{"-o", "first.txt", "--output", "second.txt", "-o", "third.txt"}

	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// The last value should be set
	if *output != "third.txt" {
		t.Errorf("Expected output 'third.txt', got '%s'", *output)
	}
}

func TestErrorHandling(t *testing.T) {
	// Test that ErrorHandling constants are defined
	if ContinueOnError != 0 {
		t.Errorf("Expected ContinueOnError to be 0, got %d", ContinueOnError)
	}
}

func TestFlagSet_EmptyArguments(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)
	verbose := fs.Bool("verbose", 'v', false, "verbose output")
	output := fs.String("output", 'o', "default.txt", "output file")

	// Parse empty arguments
	err := fs.Parse([]string{})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Should have default values
	if *verbose != false {
		t.Errorf("Expected verbose false, got %t", *verbose)
	}
	if *output != "default.txt" {
		t.Errorf("Expected output 'default.txt', got '%s'", *output)
	}

	// Should have no positional arguments
	if len(fs.Args()) != 0 {
		t.Errorf("Expected no args, got %v", fs.Args())
	}
}

func TestFlagSet_GoStyleParsing(t *testing.T) {
	// This test verifies that Go-style parsing works correctly
	// where both short and long options use the same prefix
	fs := NewFlagSet("test", ContinueOnError)
	verbose := fs.Bool("verbose", 'v', false, "verbose output")
	help := fs.Bool("help", 'h', false, "show help")

	// Create a Go-style parser using the dedicated method
	px := fs.NewGoStyleParser()

	// Try to parse Go-style arguments
	args := []string{"-v", "-help"}

	err := fs.ParseWithParser(px, args)

	// This should now work with the fix:
	// All options (both short and long) are put in LongOptions for Go-style parsing
	// since they all use the same "-" prefix

	if err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}

	// Check the values
	if !*verbose {
		t.Error("Expected verbose to be true")
	}
	if !*help {
		t.Error("Expected help to be true")
	}
}

func TestFlagSet_GNUStyleParsing(t *testing.T) {
	// This test verifies that traditional GNU-style parsing still works correctly
	// after our refactoring to separate short and long option maps
	fs := NewFlagSet("test", ContinueOnError)
	verbose := fs.Bool("verbose", 'v', false, "verbose output")
	output := fs.String("output", 'o', "", "output file")
	debug := fs.Bool("debug", 'd', false, "debug mode")

	// Use the default parser (GNU-style)
	px := fs.NewParser()

	// Test various GNU-style argument combinations
	args := []string{"-v", "--output=file.txt", "-d", "arg1", "arg2"}

	err := fs.ParseWithParser(px, args)

	if err != nil {
		t.Fatalf("GNU-style parsing failed: %v", err)
	}

	// Check short flags work
	if !*verbose {
		t.Error("Expected verbose (-v) to be true")
	}
	if !*debug {
		t.Error("Expected debug (-d) to be true")
	}

	// Check long flag with value works
	if *output != "file.txt" {
		t.Errorf("Expected output to be 'file.txt', got '%s'", *output)
	}

	// Check positional arguments
	positionalArgs := fs.Args()
	expected := []string{"arg1", "arg2"}
	if len(positionalArgs) != len(expected) {
		t.Errorf("Expected %d args, got %d", len(expected), len(positionalArgs))
	}
	for i, arg := range positionalArgs {
		if arg != expected[i] {
			t.Errorf("Expected arg[%d] '%s', got '%s'", i, expected[i], arg)
		}
	}

	// Test that short and long versions of the same flag work
	fs2 := NewFlagSet("test2", ContinueOnError)
	quiet := fs2.Bool("quiet", 'q', false, "quiet mode")

	px2 := fs2.NewParser()

	// Test short version
	err = fs2.ParseWithParser(px2, []string{"-q"})
	if err != nil {
		t.Fatalf("Short flag parsing failed: %v", err)
	}
	if !*quiet {
		t.Error("Expected quiet (-q) to be true")
	}

	// Reset and test long version
	*quiet = false
	err = fs2.ParseWithParser(px2, []string{"--quiet"})
	if err != nil {
		t.Fatalf("Long flag parsing failed: %v", err)
	}
	if !*quiet {
		t.Error("Expected quiet (--quiet) to be true")
	}
}
