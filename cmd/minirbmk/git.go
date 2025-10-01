// git.go - git subcommand
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"fmt"

	"github.com/bassosimone/clip"
	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/nflag"
)

// gitInitMain is the main entry point for the 'git init' leaf command.
func gitInitMain(ctx context.Context, args *clip.CommandArgs[*clip.StdlibExecEnv]) error {
	// Create flag set
	fset := nflag.NewFlagSet(args.CommandName, nflag.ExitOnError)
	fset.Description = args.Command.BriefDescription()
	fset.PositionalArgumentsUsage = "[directory]"
	fset.MinPositionalArgs = 0
	fset.MaxPositionalArgs = 1

	// Not strictly needed in production but necessary for testing
	fset.Exit = args.Env.Exit
	fset.Stderr = args.Env.Stderr()
	fset.Stdout = args.Env.Stdout()

	// Add the --branch, -b flag
	branchFlag := fset.StringFlag("branch", 'b', "Branch name")

	// Add the --help flag
	fset.AutoHelp("help", 'h', "Print this help message and exit.")

	// Add the -q, --quiet flag
	quietFlag := fset.BoolFlag("quiet", 'q', "Run in quiet mode.")

	// Parse the flags
	assert.NotError(fset.Parse(args.Args))

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
	fset := nflag.NewFlagSet(args.CommandName, nflag.ExitOnError)
	fset.Description = args.Command.BriefDescription()
	fset.PositionalArgumentsUsage = "<repository> [directory]"
	fset.MinPositionalArgs = 1
	fset.MaxPositionalArgs = 2

	// Not strictly needed in production but necessary for testing
	fset.Exit = args.Env.Exit
	fset.Stderr = args.Env.Stderr()
	fset.Stdout = args.Env.Stdout()

	// Add the -b flag
	branchFlag := fset.StringFlag("branch", 'b', "Branch name")

	// Add the --help flag
	fset.AutoHelp("help", 'h', "Print this help message and exit.")

	// Add the -q, --quiet flag
	quietFlag := fset.BoolFlag("quiet", 'q', "Run in quiet mode.")

	// Parse the flags
	assert.NotError(fset.Parse(args.Args))

	// Print the parsed flags
	fmt.Fprintf(args.Env.Stdout(), "branch: %s\n", *branchFlag)
	fmt.Fprintf(args.Env.Stdout(), "quiet: %v\n", *quietFlag)

	// Print the positional arguments
	fmt.Fprintf(args.Env.Stdout(), "%v\n", fset.Args())
	return nil
}
