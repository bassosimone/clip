// command.go - tests for Command
// SPDX-License-Identifier: GPL-3.0-or-later

package clip_test

import (
	"context"
	"fmt"
	"math"

	"github.com/bassosimone/clip"
	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/nflag"
)

var gzipSubcommand = &clip.LeafCommand[*clip.StdlibExecEnv]{
	BriefDescriptionText: "Compress or expand files.",
	HelpFlagValue:        "--help",
	RunFunc: func(ctx context.Context, args *clip.CommandArgs[*clip.StdlibExecEnv]) error {
		// Create command line parser
		fset := nflag.NewFlagSet(args.CommandName, nflag.ExitOnError)
		fset.Description = args.Command.BriefDescription()
		fset.PositionalArgumentsUsage = "file ..."
		fset.MinPositionalArgs = 1
		fset.MaxPositionalArgs = math.MaxInt

		// Add the --help flag
		fset.AutoHelp("help", 'h', "Print this help message and exit.")

		// Add the `-v` option
		vflag := fset.BoolFlag("verbose", 'v', "verbose mode")

		// Parse command line arguments
		assert.NotError(fset.Parse(args.Args))

		// Print the received flags and arguments
		fmt.Fprintf(args.Env.Stdout(), "Flags: -v=%v\n", *vflag)
		fmt.Fprintf(args.Env.Stdout(), "Arguments: %v\n", fset.Args())
		return nil
	},
}

var tarSubcommand = &clip.LeafCommand[*clip.StdlibExecEnv]{
	BriefDescriptionText: "Archiving utility.",
	HelpFlagValue:        "--help",
	RunFunc: func(ctx context.Context, args *clip.CommandArgs[*clip.StdlibExecEnv]) error {
		// Create command line parser
		fset := nflag.NewFlagSet(args.CommandName, nflag.ExitOnError)
		fset.Description = args.Command.BriefDescription()
		fset.PositionalArgumentsUsage = "file ..."
		fset.MinPositionalArgs = 1
		fset.MaxPositionalArgs = math.MaxInt

		// Add the `-c` option
		cflag := fset.BoolFlag("create", 'c', "create a new archive")

		// Add the `-f` option
		fflag := fset.StringFlag("file", 'f', "archive file name")

		// Add the --help flag
		fset.AutoHelp("help", 'h', "Print this help message and exit.")

		// Add the `-v` option
		vflag := fset.BoolFlag("verbose", 'v', "verbose mode")

		// Add the `-z` option
		zflag := fset.BoolFlag("gzip", 'z', "gzip compression")

		// Parse command line arguments
		assert.NotError(fset.Parse(args.Args))

		// Print the received flags and arguments
		fmt.Fprintf(args.Env.Stdout(), "Flags: -c=%v, -f=%v, -v=%v, -z=%v\n", *cflag, *fflag, *vflag, *zflag)
		fmt.Fprintf(args.Env.Stdout(), "Arguments: %v\n", fset.Args())
		return nil
	},
}

const toolVersion = "0.1.0"

var toolsDispatcher = &clip.DispatcherCommand[*clip.StdlibExecEnv]{
	BriefDescriptionText: "UNIX command-line tools.",
	Commands: map[string]clip.Command[*clip.StdlibExecEnv]{
		"gzip": gzipSubcommand,
		"tar":  tarSubcommand,
	},
	ErrorHandling:             nflag.ExitOnError,
	Version:                   toolVersion,
	OptionPrefixes:            []string{"-", "--"},
	OptionsArgumentsSeparator: "--",
}

var toplevelDispatcher = &clip.DispatcherCommand[*clip.StdlibExecEnv]{
	BriefDescriptionText: "Swiss Army Knife command-line tools.",
	Commands: map[string]clip.Command[*clip.StdlibExecEnv]{
		"tools": toolsDispatcher,
	},
	ErrorHandling:             nflag.ExitOnError,
	Version:                   toolVersion,
	OptionPrefixes:            []string{"-", "--"},
	OptionsArgumentsSeparator: "--",
}

// rootCommand is the root command of the application.
var rootCommand = &clip.RootCommand[*clip.StdlibExecEnv]{
	Command: toplevelDispatcher,
}

// This example shows how to construct a complex command line
// interface using the subcommand package.
func Example() {
	// Create environment using the standard library I/O
	env := clip.NewStdlibExecEnv()

	// execute the root command
	rootCommand.Main(env)
}
