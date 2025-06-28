// dispatcher_test.go - subcommand dispatcher tests
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/bassosimone/clip/pkg/nflag"
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
				name:           "We invoke with --version",
				skipTest:       false,
				args:           []string{"--version"},
				cmdReturnError: nil,
				expectRunCmd:   false,
				expectError:    nil,
			},

			{
				name:           "We invoke with version",
				skipTest:       false,
				args:           []string{"version"},
				cmdReturnError: nil,
				expectRunCmd:   false,
				expectError:    nil,
			},

			{
				name:           "We obtain the help of the version subcommand using --help",
				skipTest:       false,
				args:           []string{"version", "--help"},
				cmdReturnError: nil,
				expectRunCmd:   false,
				expectError:    nflag.ErrHelp,
			},

			{
				name:           "We obtain the help of the version subcommand using help",
				skipTest:       false,
				args:           []string{"help", "version"},
				cmdReturnError: nil,
				expectRunCmd:   false,
				expectError:    nflag.ErrHelp,
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
				args:           []string{"+short", "tools", "dig", "IN", "A", "example.com"},
				cmdReturnError: errMocked,
				expectRunCmd:   true,
				expectError:    errMocked,
			},

			{
				name:           "No command name just flags before the separator",
				skipTest:       false,
				args:           []string{"+short", "--", "tools", "dig"},
				cmdReturnError: errMocked,
				expectRunCmd:   false,
				expectError:    ErrInvalidFlags,
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
							Version:                   "0.1.0",
							OptionPrefixes:            []string{"--", "-", "+"},
							OptionsArgumentsSeparator: "--",
						},
						"pipelines": &DispatcherCommand[*StdlibExecEnv]{},
					},
					Version:                   "0.1.0",
					OptionPrefixes:            []string{"--", "-", "+"},
					OptionsArgumentsSeparator: "--",
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
						t.Fatalf("unexpected error: want=%v, got=%v\n", tc.expectError, err)
					}
				}

				// Check whether the command did run
				if tc.expectRunCmd != didCmdRun {
					t.Fatalf("expected command to run %v, got %v", tc.expectRunCmd, didCmdRun)
				}
			})
		}
	})

	t.Run("ExitOnError policy", func(t *testing.T) {
		// Create a command dispatcher with a command that fails.
		errCommandFailed := errors.New("command failed")
		dx := &DispatcherCommand[*StdlibExecEnv]{
			Commands: map[string]Command[*StdlibExecEnv]{
				"dig": &LeafCommand[*StdlibExecEnv]{
					RunFunc: func(ctx context.Context, args *CommandArgs[*StdlibExecEnv]) error {
						return errCommandFailed
					},
				},
				"grep": &LeafCommand[*StdlibExecEnv]{},
			},
			ErrorHandling: 0, // we will modify this field in various test cases
		}

		// Customize a stdlib env to panic on Exit
		errExit := errors.New("exit")
		env := NewStdlibExecEnv()
		env.OSExit = func(exitcode int) {
			panic(fmt.Errorf("%w: %d", errExit, exitcode))
		}

		// Create a background context
		ctx := context.Background()

		t.Run("without error", func(t *testing.T) {
			for _, policy := range []nflag.ErrorHandling{nflag.ContinueOnError, nflag.ExitOnError, nflag.PanicOnError} {
				t.Run(fmt.Sprintf("with policy %d", policy), func(t *testing.T) {
					// Set the proper error handling strategy
					dx.ErrorHandling = policy

					// Run the command
					err := dx.Run(ctx, &CommandArgs[*StdlibExecEnv]{
						Args:        []string{},
						Command:     dx,
						CommandName: "main",
						Env:         env,
					})

					// Verify the status
					if err != nil {
						t.Fatal(err)
					}
				})
			}
		})

		t.Run("ContinueOnError", func(t *testing.T) {
			// Set the proper error handling strategy
			dx.ErrorHandling = nflag.ContinueOnError

			// Run the command
			err := dx.Run(ctx, &CommandArgs[*StdlibExecEnv]{
				Args:        []string{"dig"},
				Command:     dx,
				CommandName: "main",
				Env:         env,
			})

			// Verify the status
			if !errors.Is(err, errCommandFailed) {
				t.Fatalf("expected errCommandFailed, got %T", err)
			}
		})

		t.Run("PanicOnError", func(t *testing.T) {
			// Verify the status after the panic has occurred
			defer func() {
				err := recover().(error)
				if !errors.Is(err, errCommandFailed) {
					t.Fatalf("expected errCommandFailed, got %T", err)
				}
			}()

			// Set the proper error handling strategy
			dx.ErrorHandling = nflag.PanicOnError

			// Run the command
			_ = dx.Run(ctx, &CommandArgs[*StdlibExecEnv]{
				Args:        []string{"dig"},
				Command:     dx,
				CommandName: "main",
				Env:         env,
			})

			t.Fatal("expected panic, but did not occur")
		})

		t.Run("ExitOnError with exit code 1", func(t *testing.T) {
			// Verify the status after the panic has occurred
			defer func() {
				err := recover().(error)
				if !errors.Is(err, errExit) {
					t.Fatalf("expected errCommandFailed, got %T", err)
				}
				if reason := err.Error(); reason != "exit: 1" {
					t.Fatalf("expected 'exit: 1', got %s", reason)
				}
			}()

			// Set the proper error handling strategy
			dx.ErrorHandling = nflag.ExitOnError

			// Run the command
			_ = dx.Run(ctx, &CommandArgs[*StdlibExecEnv]{
				Args:        []string{"dig"},
				Command:     dx,
				CommandName: "main",
				Env:         env,
			})

			t.Fatal("expected panic, but did not occur")
		})

		t.Run("ExitOnError with exit code 2 for usage error", func(t *testing.T) {
			// Verify the status after the panic has occurred
			defer func() {
				err := recover().(error)
				if !errors.Is(err, errExit) {
					t.Fatalf("expected errCommandFailed, got %T", err)
				}
				if reason := err.Error(); reason != "exit: 2" {
					t.Fatalf("expected 'exit: 2', got %s", reason)
				}
			}()

			// Set the proper error handling strategy
			dx.ErrorHandling = nflag.ExitOnError

			// Run the command
			_ = dx.Run(ctx, &CommandArgs[*StdlibExecEnv]{
				Args:        []string{"__antani__"},
				Command:     dx,
				CommandName: "main",
				Env:         env,
			})

			t.Fatal("expected panic, but did not occur")
		})
	})
}
