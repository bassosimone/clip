// dig.go - dig subcommand
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"fmt"

	"github.com/bassosimone/clip"
)

// digMain is the main entry point for the dig leaf command.
func digMain(ctx context.Context, args *clip.CommandArgs[*clip.StdlibExecEnv]) error {
	// Create flag set
	fset := clip.NewFlagSet(args.CommandName, clip.ExitOnError)
	fset.SetDescription(args.Command.BriefDescription())
	fset.SetArgsDocs("[@server] name [type] [class]")

	// Not strictly needed in production but necessary for testing
	fset.SetExitFunc(args.Env.Exit)
	fset.SetStderr(args.Env.Stderr())
	fset.SetStdout(args.Env.Stdout())

	// Customize the parser
	px := fset.Parser()
	px.LongOptionPrefixes = []string{"+", "--"}

	// Add the -4 flag
	fourFlag := fset.BoolShort('4', "Only use IPv4")

	// Add the +short flag
	shortFlag := fset.BoolLong("short", "Print a terse query representation.")

	// Parse the flags; note that ExitOnError is set, so it will exit on error
	_ = fset.Parse(args.Args)

	// Parse the positional arguments; note that ExitOnError is set, so it will exit on error
	_ = fset.PositionalArgsRangeCheck(1, 4)

	// Print the parsed flags
	fmt.Fprintf(args.Env.Stdout(), "-4: %v\n", *fourFlag)
	fmt.Fprintf(args.Env.Stdout(), "+short: %v\n", *shortFlag)

	// Print the positional arguments
	fmt.Fprintf(args.Env.Stdout(), "%v\n", fset.Args())
	return nil
}
