// flagset_test.go - Tests for FlagSet implementation
// SPDX-License-Identifier: GPL-3.0-or-later

package flag

import (
	"errors"
	"testing"

	"github.com/bassosimone/clip/pkg/parser"
)

// mockValue is a test implementation of the Value interface
type mockValue struct {
	value      string
	optionType parser.OptionType
	setError   error
}

func (m *mockValue) OptionType() parser.OptionType {
	return m.optionType
}

func (m *mockValue) String() string {
	return m.value
}

func (m *mockValue) Set(value string) error {
	if m.setError != nil {
		return m.setError
	}
	m.value = value
	return nil
}

func TestFlagSet_Parse_Errors(t *testing.T) {
	tests := []struct {
		name          string
		setupFlagSet  func() *FlagSet
		args          []string
		expectedError error
		errorContains string
	}{
		{
			name: "unknown short option",
			setupFlagSet: func() *FlagSet {
				fs := NewFlagSet("test", ContinueOnError)
				return fs
			},
			args:          []string{"-x"},
			expectedError: parser.ErrUnknownOption,
			errorContains: "unknown option: x",
		},

		{
			name: "unknown long option",
			setupFlagSet: func() *FlagSet {
				fs := NewFlagSet("test", ContinueOnError)
				return fs
			},
			args:          []string{"--unknown"},
			expectedError: parser.ErrUnknownOption,
			errorContains: "unknown option: unknown",
		},

		{
			name: "option value set error",
			setupFlagSet: func() *FlagSet {
				fs := NewFlagSet("test", ContinueOnError)
				mockVal := &mockValue{
					optionType: parser.OptionTypeString,
					setError:   errors.New("invalid value"),
				}
				opt := &Option{
					LongName: "test",
					Value:    mockVal,
				}
				fs.AddOption(opt)
				return fs
			},
			args:          []string{"--test=value"},
			expectedError: nil, // mocked error, so we check contains instead
			errorContains: "when setting value \"value\" for option \"test\": invalid value",
		},

		{
			name: "option value set error with short option",
			setupFlagSet: func() *FlagSet {
				fs := NewFlagSet("test", ContinueOnError)
				mockVal := &mockValue{
					optionType: parser.OptionTypeString,
					setError:   errors.New("conversion failed"),
				}
				opt := &Option{
					ShortName: 't',
					Value:     mockVal,
				}
				fs.AddOption(opt)
				return fs
			},
			args:          []string{"-t", "badvalue"},
			expectedError: nil, // mocked error, so we check contains instead
			errorContains: "when setting value \"badvalue\" for option \"t\": conversion failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.setupFlagSet()
			err := fs.Parse(tt.args)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("FlagSet.Parse() expected error %v, got nil", tt.expectedError)
				} else if !errors.Is(err, tt.expectedError) {
					t.Errorf("FlagSet.Parse() error = %v, want %v", err, tt.expectedError)
				}
			}

			if tt.errorContains != "" {
				if err == nil {
					t.Errorf("FlagSet.Parse() expected error containing %q, got nil", tt.errorContains)
				} else if err.Error() != tt.errorContains {
					t.Errorf("FlagSet.Parse() error = %q, want error containing %q", err.Error(), tt.errorContains)
				}
			}
		})
	}
}

func TestFlagSet_ErrorHandling(t *testing.T) {
	t.Run("ContinueOnError returns error", func(t *testing.T) {
		fs := NewFlagSet("test", ContinueOnError)
		err := fs.Parse([]string{"--unknown"})

		if err == nil {
			t.Error("Expected error but got nil")
		} else if !errors.Is(err, parser.ErrUnknownOption) {
			t.Errorf("FlagSet.Parse() error = %v, want %v", err, parser.ErrUnknownOption)
		}
	})

	t.Run("ExitOnError calls exit with parser.ErrHelp", func(t *testing.T) {
		var exitCalled bool
		var exitCode int

		fs := NewFlagSet("test", ExitOnError)

		// Mock the exit function to simulate os.Exit behavior
		fs.SetExitFunc(func(code int) {
			exitCalled = true
			exitCode = code
			// Use panic to simulate exit stopping execution
			// In real usage, `os.Exit` would stop the program
			panic("simulated exit")
		})

		// Expect the simulated exit panic
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected simulated exit panic but none occurred")
			} else if r != "simulated exit" {
				t.Errorf("Expected simulated exit panic, got %v", r)
			} else {
				// Check that exit was called with correct code
				if !exitCalled {
					t.Error("Expected exit to be called but it wasn't")
				}
				if exitCode != 0 {
					t.Errorf("Expected exit code 0, got %d", exitCode)
				}
			}
		}()

		// Parse with an unknown option to trigger an error
		fs.Parse([]string{"--help"})

		// Should not reach here due to panic
		t.Error("Should not reach here - exit should have been called")
	})

	t.Run("ExitOnError calls exit with other errors", func(t *testing.T) {
		var exitCalled bool
		var exitCode int

		fs := NewFlagSet("test", ExitOnError)

		// Mock the exit function to simulate os.Exit behavior
		fs.SetExitFunc(func(code int) {
			exitCalled = true
			exitCode = code
			// Use panic to simulate exit stopping execution
			// In real usage, `os.Exit` would stop the program
			panic("simulated exit")
		})

		// Expect the simulated exit panic
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected simulated exit panic but none occurred")
			} else if r != "simulated exit" {
				t.Errorf("Expected simulated exit panic, got %v", r)
			} else {
				// Check that exit was called with correct code
				if !exitCalled {
					t.Error("Expected exit to be called but it wasn't")
				}
				if exitCode != 2 {
					t.Errorf("Expected exit code 2, got %d", exitCode)
				}
			}
		}()

		// Parse with an unknown option to trigger an error
		fs.Parse([]string{"--unknown"})

		// Should not reach here due to panic
		t.Error("Should not reach here - exit should have been called")
	})

	t.Run("PanicOnError panics", func(t *testing.T) {
		fs := NewFlagSet("test", PanicOnError)

		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic but none occurred")
			}
		}()

		// Parse with an unknown option to trigger an error
		fs.Parse([]string{"--unknown"})

		// Should not reach here due to panic
		t.Error("Should not reach here - panic should have occurred")
	})

	t.Run("invalid ErrorHandling value defaults to panic", func(t *testing.T) {
		fs := NewFlagSet("test", ErrorHandling(999))

		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic but none occurred")
			}
		}()

		// Parse with an unknown option to trigger an error
		fs.Parse([]string{"--unknown"})

		// Should not reach here due to panic
		t.Error("Should not reach here - panic should have occurred")
	})
}

func TestFlagSet_Parse_Success(t *testing.T) {
	tests := []struct {
		name          string
		setupFlagSet  func() *FlagSet
		args          []string
		expectedArgs  []string
		expectedFlags []string
	}{
		{
			name: "no arguments",
			setupFlagSet: func() *FlagSet {
				return NewFlagSet("test", ContinueOnError)
			},
			args:          []string{},
			expectedArgs:  []string{},
			expectedFlags: []string{},
		},

		{
			name: "only positional arguments",
			setupFlagSet: func() *FlagSet {
				return NewFlagSet("test", ContinueOnError)
			},
			args:          []string{"arg1", "arg2", "arg3"},
			expectedArgs:  []string{"arg1", "arg2", "arg3"},
			expectedFlags: []string{},
		},

		{
			name: "mixed options and arguments",
			setupFlagSet: func() *FlagSet {
				fs := NewFlagSet("test", ContinueOnError)
				mockVal := &mockValue{optionType: parser.OptionTypeString}
				opt := &Option{
					LongName: "file",
					Value:    mockVal,
				}
				fs.AddOption(opt)
				return fs
			},
			args:          []string{"--file=test.txt", "arg1", "arg2"},
			expectedArgs:  []string{"arg1", "arg2"},
			expectedFlags: []string{"file"},
		},

		{
			name: "separator handling",
			setupFlagSet: func() *FlagSet {
				fs := NewFlagSet("test", ContinueOnError)
				mockVal := &mockValue{optionType: parser.OptionTypeString}
				opt := &Option{
					LongName: "file",
					Value:    mockVal,
				}
				fs.AddOption(opt)
				return fs
			},
			args:          []string{"--file=test.txt", "--", "--not-an-option"},
			expectedArgs:  []string{"--not-an-option"},
			expectedFlags: []string{"file"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.setupFlagSet()
			err := fs.Parse(tt.args)

			if err != nil {
				t.Fatalf("FlagSet.Parse() unexpected error = %v", err)
			}

			// Make sure all expected args are present
			args := fs.Args()
			if len(args) != len(tt.expectedArgs) {
				t.Fatalf("FlagSet.Args() length = %d, want %d", len(args), len(tt.expectedArgs))
			}
			for i, arg := range args {
				if i >= len(tt.expectedArgs) || arg != tt.expectedArgs[i] {
					t.Fatalf("FlagSet.Args()[%d] = %q, want %q", i, arg, tt.expectedArgs[i])
				}
			}

			// Make sure all expected flags are present
			const (
				expectflag = 1 << iota
				gotflag
			)
			account := map[string]int{}
			for _, opt := range fs.Options() {
				if opt.Modified {
					if opt.LongName != "" {
						account[opt.LongName] |= gotflag
					}
					if opt.ShortName != 0 {
						account[string(opt.ShortName)] |= gotflag
					}
				}
			}
			for _, opt := range tt.expectedFlags {
				account[opt] |= expectflag
			}
			for name, value := range account {
				if value != expectflag|gotflag {
					t.Errorf("FlagSet.Option(%q) = %d, want %d", name, value, expectflag|gotflag)
				}
			}
		})
	}
}

func TestFlagSet_AddOption(t *testing.T) {
	tests := []struct {
		name      string
		option    *Option
		testLong  bool
		testShort bool
		longName  string
		shortName string
	}{
		{
			name: "long option only",
			option: &Option{
				LongName: "verbose",
				Value:    &mockValue{optionType: parser.OptionTypeBool},
			},
			testLong:  true,
			testShort: false,
			longName:  "verbose",
		},

		{
			name: "short option only",
			option: &Option{
				ShortName: 'v',
				Value:     &mockValue{optionType: parser.OptionTypeBool},
			},
			testLong:  false,
			testShort: true,
			shortName: "v",
		},

		{
			name: "both long and short",
			option: &Option{
				LongName:  "verbose",
				ShortName: 'v',
				Value:     &mockValue{optionType: parser.OptionTypeBool},
			},
			testLong:  true,
			testShort: true,
			longName:  "verbose",
			shortName: "v",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := NewFlagSet("test", ContinueOnError)
			fs.AddOption(tt.option)

			if tt.testLong {
				if _, found := fs.optionsLong[tt.longName]; !found {
					t.Errorf("Long option %q was not added", tt.longName)
				}
				if _, found := fs.parser.LongOptions[tt.longName]; !found {
					t.Errorf("Long option %q was not added to parser", tt.longName)
				}
			}

			if tt.testShort {
				if _, found := fs.optionsShort[tt.shortName]; !found {
					t.Errorf("Short option %q was not added", tt.shortName)
				}
				if _, found := fs.parser.ShortOptions[tt.shortName]; !found {
					t.Errorf("Short option %q was not added to parser", tt.shortName)
				}
			}
		})
	}
}

func TestFlagSet_findOption(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)

	longOpt := &Option{
		LongName: "verbose",
		Value:    &mockValue{optionType: parser.OptionTypeBool},
	}
	shortOpt := &Option{
		ShortName: 'v',
		Value:     &mockValue{optionType: parser.OptionTypeBool},
	}

	fs.AddOption(longOpt)
	fs.AddOption(shortOpt)

	// Test finding long option
	opt, found := fs.findOption("verbose")
	if !found {
		t.Error("Expected to find long option 'verbose'")
	}
	if opt != longOpt {
		t.Error("Found option doesn't match expected long option")
	}

	// Test finding short option
	opt, found = fs.findOption("v")
	if !found {
		t.Error("Expected to find short option 'v'")
	}
	if opt != shortOpt {
		t.Error("Found option doesn't match expected short option")
	}

	// Test not finding option
	_, found = fs.findOption("nonexistent")
	if found {
		t.Error("Expected not to find option 'nonexistent'")
	}
}

func TestNewFlagSet(t *testing.T) {
	tests := []struct {
		name     string
		progname string
		handling ErrorHandling
	}{
		{
			name:     "ContinueOnError",
			progname: "test",
			handling: ContinueOnError,
		},

		{
			name:     "ExitOnError",
			progname: "myprogram",
			handling: ExitOnError,
		},

		{
			name:     "PanicOnError",
			progname: "panictest",
			handling: PanicOnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := NewFlagSet(tt.progname, tt.handling)

			if fs.progname != tt.progname {
				t.Errorf("NewFlagSet() progname = %q, want %q", fs.progname, tt.progname)
			}

			if fs.handling != tt.handling {
				t.Errorf("NewFlagSet() handling = %v, want %v", fs.handling, tt.handling)
			}

			if fs.args == nil {
				t.Error("NewFlagSet() args slice is nil")
			}

			if fs.optionsShort == nil {
				t.Error("NewFlagSet() optionsShort map is nil")
			}

			if fs.optionsLong == nil {
				t.Error("NewFlagSet() optionsLong map is nil")
			}

			if fs.parser == nil {
				t.Error("NewFlagSet() parser is nil")
			}
		})
	}
}

func TestFlagSet_Parser(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)
	parser := fs.Parser()

	if parser == nil {
		t.Error("Parser() returned nil")
	}

	if parser != fs.parser {
		t.Error("Parser() returned different parser than internal one")
	}

	// Test that we can modify the parser
	parser.LongOptionPrefixes = append(parser.LongOptionPrefixes, "---")
	if len(fs.parser.LongOptionPrefixes) != 2 {
		t.Error("Parser modification was not reflected in flagset")
	}
}
