// version.go - automatic handling of --version and version.
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import (
	"context"
	"fmt"

	"github.com/bassosimone/clip/pkg/nflag"
)

// VersionCommand implements the version command.
//
// The zero value is ready to use.
type VersionCommand[T ExecEnv] struct {
	// BriefDescriptionText is the optional brief description text.
	//
	// When unset, we use a reasonable default value.
	BriefDescriptionText string

	// ErrorHandling is the optional error handling strategy.
	//
	// When unset, we use [ContinueOnError].
	ErrorHandling nflag.ErrorHandling

	// HelpFlagValue is the optional help flag. When unset, we use "--help".
	HelpFlagValue string

	// Version is the optional version. When unsed, we use "dev".
	Version string
}

var _ Command[*StdlibExecEnv] = &VersionCommand[*StdlibExecEnv]{}

// BriefDescription implements [Command].
func (c *VersionCommand[T]) BriefDescription() string {
	output := "Print the program version and exit."
	if c.BriefDescriptionText != "" {
		output = c.BriefDescriptionText
	}
	return output
}

// HelpFlag implements [Command].
func (c *VersionCommand[T]) HelpFlag() string {
	output := "--help"
	if c.HelpFlagValue != "" {
		output = c.HelpFlagValue
	}
	return output
}

// PrintVersion prints the version to the stdout.
func (c *VersionCommand[T]) PrintVersion(env T) error {
	version := "dev"
	if c.Version != "" {
		version = c.Version
	}
	_, err := fmt.Fprintf(env.Stdout(), "%s\n", version)
	return err
}

// Run implements [Command].
func (c *VersionCommand[T]) Run(ctx context.Context, args *CommandArgs[T]) error {
	// Create empty command line parser.
	clp := nflag.NewFlagSet(args.CommandName, c.ErrorHandling)
	clp.Description = args.Command.BriefDescription()
	clp.PositionalArgumentsUsage = "" // do not print a name for positional arguments

	// Explicitly set positional arguments to 0 for clarity
	// even though this is already the default.
	clp.MinPositionalArgs = 0
	clp.MaxPositionalArgs = 0

	// Add the `-h, --help` flag.
	clp.AutoHelp("help", 'h', "Show this help message and exit.")

	// Parse the command line arguments.
	if err := clp.Parse(args.Args); err != nil {
		return err
	}

	// Print the version to the standard output.
	return c.PrintVersion(args.Env)
}

// SupportsSubcommands implements [Command].
func (c *VersionCommand[T]) SupportsSubcommands() bool {
	return false
}
