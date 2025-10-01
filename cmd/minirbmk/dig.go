// dig.go - dig subcommand
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"fmt"

	"github.com/bassosimone/clip"
	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/nflag"
)

// digMain is the main entry point for the dig leaf command.
func digMain(ctx context.Context, args *clip.CommandArgs[*clip.StdlibExecEnv]) error {
	// Create flag set
	fset := nflag.NewFlagSet(args.CommandName, nflag.ExitOnError)
	fset.Description = args.Command.BriefDescription()
	fset.PositionalArgumentsUsage = "[@server] name [type] [class]"
	fset.MinPositionalArgs = 1
	fset.MaxPositionalArgs = 4
	fset.LongFlagPrefix = "+"
	fset.ShortFlagPrefix = "-" // already the default, but set explicitly for clarity

	// Not strictly needed in production but necessary for testing
	fset.Exit = args.Env.Exit
	fset.Stderr = args.Env.Stderr()
	fset.Stdout = args.Env.Stdout()

	// Add the -4 flag
	fourFlag := fset.BoolFlag("", '4', "Only use IPv4")

	// Add the -h flag
	fset.AutoHelp("", 'h', "Print this help message and exit.")

	// Add the +short flag
	shortFlag := fset.BoolFlag("short", 0, "Print a terse query representation.")

	// Parse the flags
	assert.NotError(fset.Parse(args.Args))

	// Print the parsed flags
	fmt.Fprintf(args.Env.Stdout(), "-4: %v\n", *fourFlag)
	fmt.Fprintf(args.Env.Stdout(), "+short: %v\n", *shortFlag)

	// Print the positional arguments
	fmt.Fprintf(args.Env.Stdout(), "%v\n", fset.Args())
	return nil
}
