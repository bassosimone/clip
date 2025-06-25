// curl.go - curl subcommand
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"fmt"
	"math"

	"github.com/bassosimone/clip"
)

// curlMain is the main entry point for the curl leaf command.
func curlMain(ctx context.Context, args *clip.CommandArgs[*clip.StdlibExecEnv]) error {
	// Create flag set
	fset := clip.NewFlagSet(args.CommandName, clip.ExitOnError)
	fset.SetDescription(args.Command.BriefDescription())
	fset.SetArgsDocs("URL ...")

	// Not strictly needed in production but necessary for testing
	fset.SetExitFunc(args.Env.Exit)
	fset.SetStderr(args.Env.Stderr())
	fset.SetStdout(args.Env.Stdout())

	// Add the --cacert flag
	cacertFlag := fset.StringLong("cacert", "", "Add part to the CA certificate file.")

	// Add the -c, --cookiejar flag
	cookieJarFlag := fset.String("cookiejar", 'c', "", "Path of the file containing cookies data")

	// Add the -v, --verbose flag
	verboseFlag := fset.Bool("verbose", 'v', "Run in verbose mode.")

	// Parse the flags; note that ExitOnError is set, so it will exit on error
	_ = fset.Parse(args.Args)

	// Parse the positional arguments; note that ExitOnError is set, so it will exit on error
	_ = fset.PositionalArgsRangeCheck(1, math.MaxInt)

	// Print the parsed flags
	fmt.Fprintf(args.Env.Stdout(), "cacert: %s\n", *cacertFlag)
	fmt.Fprintf(args.Env.Stdout(), "cookiejar: %s\n", *cookieJarFlag)
	fmt.Fprintf(args.Env.Stdout(), "verbose: %v\n", *verboseFlag)

	// Print the positional arguments
	fmt.Fprintf(args.Env.Stdout(), "%v\n", fset.Args())
	return nil
}
