// curl.go - curl subcommand
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"fmt"
	"math"

	"github.com/bassosimone/clip"
	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/nflag"
)

// curlMain is the main entry point for the curl leaf command.
func curlMain(ctx context.Context, args *clip.CommandArgs[*clip.StdlibExecEnv]) error {
	// Create flag set
	fset := nflag.NewFlagSet(args.CommandName, nflag.ExitOnError)
	fset.Description = args.Command.BriefDescription()
	fset.PositionalArgumentsUsage = "URL ..."
	fset.MinPositionalArgs = 1
	fset.MaxPositionalArgs = math.MaxInt

	// Not strictly needed in production but necessary for testing
	fset.Exit = args.Env.Exit
	fset.Stderr = args.Env.Stderr()
	fset.Stdout = args.Env.Stdout()

	// Add the --cacert flag
	cacertFlag := fset.StringFlag("cacert", 0, "Add part to the CA certificate file.")

	// Add the -c, --cookiejar flag
	cookieJarFlag := fset.StringFlag("cookiejar", 'c', "Path of the file containing cookies data")

	// Add the --help flag
	fset.AutoHelp("help", 'h', "Print this help message and exit.")

	// Add the -v, --verbose flag
	verboseFlag := fset.BoolFlag("verbose", 'v', "Run in verbose mode.")

	// Parse the flags
	assert.NotError(fset.Parse(args.Args))

	// Print the parsed flags
	fmt.Fprintf(args.Env.Stdout(), "cacert: %s\n", *cacertFlag)
	fmt.Fprintf(args.Env.Stdout(), "cookiejar: %s\n", *cookieJarFlag)
	fmt.Fprintf(args.Env.Stdout(), "verbose: %v\n", *verboseFlag)

	// Print the positional arguments
	fmt.Fprintf(args.Env.Stdout(), "%v\n", fset.Args())
	return nil
}
