// git.go - git subcommand
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"fmt"

	"github.com/bassosimone/clip"
)

// gitInitMain is the main entry point for the 'git init' leaf command.
func gitInitMain(ctx context.Context, args *clip.CommandArgs[*clip.StdlibExecEnv]) error {
	// Create flag set
	fset := clip.NewFlagSet(args.CommandName, clip.ExitOnError)
	fset.SetDescription(args.Command.BriefDescription())
	fset.SetArgsDocs("[directory]")

	// Not strictly needed in production but necessary for testing
	fset.SetExitFunc(args.Env.Exit)
	fset.SetStderr(args.Env.Stderr())
	fset.SetStdout(args.Env.Stdout())

	// Add the --branch, -b flag
	branchFlag := fset.String("branch", 'b', "", "Branch name")

	// Add the -q, --quiet flag
	quietFlag := fset.Bool("quiet", 'q', "Run in quiet mode.")

	// Parse the flags; note that ExitOnError is set, so it will exit on error
	_ = fset.Parse(args.Args)

	// Parse the positional arguments; note that ExitOnError is set, so it will exit on error
	_ = fset.PositionalArgsRangeCheck(0, 1)

	// Print the parsed flags
	fmt.Fprintf(args.Env.Stdout(), "branch: %s\n", *branchFlag)
	fmt.Fprintf(args.Env.Stdout(), "quiet: %v\n", *quietFlag)

	// Print the positional arguments
	fmt.Fprintf(args.Env.Stdout(), "%v\n", fset.Args())
	return nil
}

// gitCloneMain is the main entry point for the 'git clone' leaf command.
func gitCloneMain(ctx context.Context, args *clip.CommandArgs[*clip.StdlibExecEnv]) error {
	// Create flag set
	fset := clip.NewFlagSet(args.CommandName, clip.ExitOnError)
	fset.SetDescription(args.Command.BriefDescription())
	fset.SetArgsDocs("<repository> [directory]")

	// Not strictly needed in production but necessary for testing
	fset.SetExitFunc(args.Env.Exit)
	fset.SetStderr(args.Env.Stderr())
	fset.SetStdout(args.Env.Stdout())

	// Add the -b flag
	branchFlag := fset.String("branch", 'b', "", "Branch name")

	// Add the -q, --quiet flag
	quietFlag := fset.Bool("quiet", 'q', "Run in quiet mode.")

	// Parse the flags; note that ExitOnError is set, so it will exit on error
	_ = fset.Parse(args.Args)

	// Parse the positional arguments; note that ExitOnError is set, so it will exit on error
	_ = fset.PositionalArgsRangeCheck(1, 2)

	// Print the parsed flags
	fmt.Fprintf(args.Env.Stdout(), "branch: %s\n", *branchFlag)
	fmt.Fprintf(args.Env.Stdout(), "quiet: %v\n", *quietFlag)

	// Print the positional arguments
	fmt.Fprintf(args.Env.Stdout(), "%v\n", fset.Args())
	return nil
}
