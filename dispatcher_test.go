// dispatcher_test.go - subcommand dispatcher tests
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDispatch(t *testing.T) {
	t.Run("BriefDescription", func(t *testing.T) {
		dx := &DispatcherCommand[*StdlibExecEnv]{BriefDescriptionText: "Test Dispatcher"}
		if got := dx.BriefDescription(); got != "Test Dispatcher" {
			t.Errorf("BriefDescription() = %q, want %q", got, "Test Dispatcher")
		}
	})

	t.Run("HelpFlag", func(t *testing.T) {
		dx := &DispatcherCommand[*StdlibExecEnv]{}
		if got := dx.HelpFlag(); got != "--help" {
			t.Errorf("HelpFlag() = %q, want %q", got, "--help")
		}
	})

	t.Run("SupportsSubcommands", func(t *testing.T) {
		dx := &DispatcherCommand[*StdlibExecEnv]{
			Commands: map[string]Command[*StdlibExecEnv]{
				"test": &DispatcherCommand[*StdlibExecEnv]{},
			},
		}
		if !dx.SupportsSubcommands() {
			t.Error("SupportsSubcommands() = false, want true")
		}
	})

	t.Run("formatUsage", func(t *testing.T) {
		expected := "Usage: test [subcommand]\n\nAvailable subcommands:\n  dig    Dig for something\n  tools  Tools for various tasks\n"
		dx := &DispatcherCommand[*StdlibExecEnv]{
			Usage: expected,
		}
		if diff := cmp.Diff(dx.formatUsage("test"), expected); diff != "" {
			t.Errorf("formatUsage() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Main", func(t *testing.T) {

		type testcase struct {
			// name is the name of the test case.
			name string

			// skipTest indicates whether we should skip this test case.
			skipTest bool

			// args contains the command line arguments to test.
			args []string

			// cmdReturnError is the error that the command should return.
			cmdReturnError error

			// expectRunCmd is whether we expect to run the command.
			expectRunCmd bool

			// expectError is the expected error
			expectError error
		}

		// Define an error to simulate the command failing
		errMocked := errors.New("mocked error")

		cases := []testcase{
			{
				name:           "We invoke without arguments",
				skipTest:       false,
				args:           []string{},
				cmdReturnError: nil,
				expectRunCmd:   false,
				expectError:    nil,
			},

			{
				name:           "We invoke with --help",
				skipTest:       false,
				args:           []string{"--help"},
				cmdReturnError: nil,
				expectRunCmd:   false,
				expectError:    nil,
			},

			{
				name:           "We invoke with -h",
				skipTest:       false,
				args:           []string{"-h"},
				cmdReturnError: nil,
				expectRunCmd:   false,
				expectError:    nil,
			},

			{
				name:           "We invoke with help without arguments",
				skipTest:       false,
				args:           []string{"help"},
				cmdReturnError: nil,
				expectRunCmd:   false,
				expectError:    nil,
			},

			{
				name:           "We invoke with help for nonexisting subcommand",
				skipTest:       false,
				args:           []string{"help", "__antani__"},
				cmdReturnError: nil,
				expectRunCmd:   false,
				expectError:    ErrNoSuchCommand,
			},

			{
				name:           "We invoke with help for the dig subcommand",
				skipTest:       false,
				args:           []string{"help", "tools", "dig"},
				cmdReturnError: nil,
				expectRunCmd:   true,
				expectError:    nil,
			},

			{
				name:           "We invoke with help for the pipelines subcommand",
				skipTest:       false,
				args:           []string{"help", "pipelines"},
				cmdReturnError: nil,
				expectRunCmd:   false,
				expectError:    nil,
			},

			{
				name:           "We invoke a nonexisting subcommand",
				skipTest:       false,
				args:           []string{"__antani__"},
				cmdReturnError: nil,
				expectRunCmd:   false,
				expectError:    ErrNoSuchCommand,
			},

			{
				name:           "We invoke the dig command and it succeeds",
				skipTest:       false,
				args:           []string{"tools", "dig", "IN", "+short", "A", "example.com"},
				cmdReturnError: nil,
				expectRunCmd:   true,
				expectError:    nil,
			},

			{
				name:           "We invoke the dig command and it fails",
				skipTest:       false,
				args:           []string{"tools", "dig", "IN", "+short", "A", "example.com"},
				cmdReturnError: errMocked,
				expectRunCmd:   true,
				expectError:    errMocked,
			},

			{
				name:           "Repair with ambiguous command line",
				skipTest:       false,
				args:           []string{"IN", "A", "+short", "tools", "dig", "example.com"},
				cmdReturnError: errMocked,
				expectRunCmd:   true,
				expectError:    errMocked,
			},

			{
				name:           "Impossible to repair ambiguity",
				skipTest:       false,
				args:           []string{"IN", "A", "+short", "tools", "dig", "example.com", "pipelines"},
				cmdReturnError: nil,
				expectRunCmd:   false,
				expectError:    ErrAmbiguousCommandLine,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				// Create the execution environment mock
				env := NewStdlibExecEnv()

				// Create the dispatcher with a single command
				didCmdRun := false
				dx := &DispatcherCommand[*StdlibExecEnv]{
					Commands: map[string]Command[*StdlibExecEnv]{
						"tools": &DispatcherCommand[*StdlibExecEnv]{
							Commands: map[string]Command[*StdlibExecEnv]{
								"dig": &LeafCommand[*StdlibExecEnv]{
									RunFunc: func(ctx context.Context, args *CommandArgs[*StdlibExecEnv]) error {
										didCmdRun = true
										return tc.cmdReturnError
									},
								},
							},
						},
						"pipelines": &DispatcherCommand[*StdlibExecEnv]{},
					},
				}

				// Run the command
				args := &CommandArgs[*StdlibExecEnv]{
					Args:        tc.args,
					Command:     dx,
					CommandName: "test",
					Env:         env,
					Parent:      nil,
				}
				err := dx.Run(context.Background(), args)

				// Check whether the error is consistent with the expectation
				switch {
				case tc.expectError == nil && err == nil:
					// As expected
				case tc.expectError != nil && err == nil:
					t.Fatal("expected", tc.expectError, "got", err)
				case tc.expectError == nil && err != nil:
					t.Fatal("expected", tc.expectError, "got", err)
				default:
					if !errors.Is(err, tc.expectError) {
						t.Fatalf("unexpected error: want=%T, got=%T\n", tc.expectError, err)
					}
				}

				// Check whether the command did run
				if tc.expectRunCmd != didCmdRun {
					t.Fatalf("expected command to run %v, got %v", tc.expectRunCmd, didCmdRun)
				}
			})
		}
	})
}
