// command.go - tests for Command
// SPDX-License-Identifier: GPL-3.0-or-later

package clip_test

import (
	"context"
	"fmt"
	"math"

	"github.com/bassosimone/clip"
)

var gzipSubcommand = &clip.LeafCommand[*clip.StdlibExecEnv]{
	BriefDescriptionText: "Compress or expand files.",
	HelpFlagValue:        "--help",
	RunFunc: func(ctx context.Context, args *clip.CommandArgs[*clip.StdlibExecEnv]) error {
		// Create command line parser
		fset := clip.NewFlagSet(args.CommandName, clip.ExitOnError)
		fset.SetDescription(args.Command.BriefDescription())
		fset.SetArgsDocs("file ...")

		// Add the `-v` option
		vflag := fset.Bool("verbose", 'v', false, "verbose mode")

		// Parse command line arguments
		clip.Must(args.Env, fset.Parse(args.Args))

		// Validate number of positional arguments
		clip.Must(args.Env, fset.PositionalArgsRangeCheck(1, math.MaxInt))

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
		fset := clip.NewFlagSet(args.CommandName, clip.ExitOnError)
		fset.SetDescription(args.Command.BriefDescription())
		fset.SetArgsDocs("file ...")

		// Add the `-c` option
		cflag := fset.Bool("create", 'c', false, "create a new archive")

		// Add the `-f` option
		fflag := fset.String("file", 'f', "", "archive file name")

		// Add the `-v` option
		vflag := fset.Bool("verbose", 'v', false, "verbose mode")

		// Add the `-z` option
		zflag := fset.Bool("gzip", 'z', false, "gzip compression")

		// Parse command line arguments
		clip.Must(args.Env, fset.Parse(args.Args))

		// Validate number of positional arguments
		clip.Must(args.Env, fset.PositionalArgsRangeCheck(1, math.MaxInt))

		// Print the received flags and arguments
		fmt.Fprintf(args.Env.Stdout(), "Flags: -c=%v, -f=%v, -v=%v, -z=%v\n", *cflag, *fflag, *vflag, *zflag)
		fmt.Fprintf(args.Env.Stdout(), "Arguments: %v\n", fset.Args())
		return nil
	},
}

var versionSubcommand = &clip.LeafCommand[*clip.StdlibExecEnv]{
	BriefDescriptionText: "Print the version of the program.",
	HelpFlagValue:        "--help",
	RunFunc: func(ctx context.Context, args *clip.CommandArgs[*clip.StdlibExecEnv]) error {
		// Create command line parser
		fset := clip.NewFlagSet(args.CommandName, clip.ExitOnError)
		fset.SetDescription("Print the version of the program.")
		fset.SetArgsDocs("")

		// Parse command line arguments
		clip.Must(args.Env, fset.Parse(args.Args))

		// Validate number of positional arguments
		clip.Must(args.Env, fset.PositionalArgsEqualCheck(0))

		// Print the received flags and arguments
		fmt.Fprintf(args.Env.Stdout(), "0.1.0\n")
		return nil
	},
}

var toolsDispatcher = &clip.DispatcherCommand[*clip.StdlibExecEnv]{
	BriefDescriptionText: "UNIX command-line tools.",
	Commands: map[string]clip.Command[*clip.StdlibExecEnv]{
		"gzip": gzipSubcommand,
		"tar":  tarSubcommand,
	},
	ErrorHandling: clip.ExitOnError,
}

var toplevelDispatcher = &clip.DispatcherCommand[*clip.StdlibExecEnv]{
	BriefDescriptionText: "Swiss Army Knife command-line tools.",
	Commands: map[string]clip.Command[*clip.StdlibExecEnv]{
		"tools":   toolsDispatcher,
		"version": versionSubcommand,
	},
	ErrorHandling: clip.ExitOnError,
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
