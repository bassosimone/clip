# clip: Command Line Parser

[![Package-level Go docs](https://pkg.go.dev/badge/github.com/bassosimone/clip)](https://pkg.go.dev/github.com/bassosimone/clip) [![Build Status](https://github.com/bassosimone/clip/actions/workflows/go.yml/badge.svg)](https://github.com/bassosimone/clip/actions) [![codecov](https://codecov.io/gh/bassosimone/clip/branch/main/graph/badge.svg)](https://codecov.io/gh/bassosimone/clip)

This repository implements a very flexible command line parser
written in Go. It provides an intuitive flag parsing API modeled
after the standard library's [flag](https://pkg.go.dev/flag)
package. It also provides means to create subcommands.

By default, [clip](https://github.com/bassosimone/clip) implements
[GNU getopt](https://linux.die.net/man/3/getopt) compatible
command line parsing:

1. Short options introduced by `-`, long options introduced by `--`.

2. The ability to mix short and long options.

3. Automatic (but configurable) options and arguments permutation.

4. The `--` separator to terminate options processing.

However, what sets [clip](https://github.com/bassosimone/clip) apart is
the possibility of customizing the prefixes used for short and long
options. For example, it is possible to customize it to:

1. Only use long options introduced by `-` and `--`, like in the
Go [flag](https://pkg.go.dev/flag) package.

2. Allow long options to also be introduced by `+` like in the
[dig](https://linux.die.net/man/1/dig) command line tool.

3. Have options be introduced by `/`, like Windows tools do.

Therefore, [clip](https://github.com/bassosimone/clip) is suitable
for writing complex command line tools that require to emulate other
command line tools behavior in their subcommands, such as, for
example, the [rbmk](https://github.com/rbmk-project/rbmk)
network measurement tool.

## Examples

This section shows how to use the packages in this repository.

### pkg/flag

The following example shows how to use the [clip](.) package:

```Go
package main

import (
	"context"
	"fmt"
	"math"

	"github.com/bassosimone/clip"
)

// Create a subcommand working a bit like gzip
var gzipSubcommand = &clip.LeafCommand[*clip.StdlibExecEnv]{
	BriefDescriptionText: "Compress or expand files.",
	HelpFlagValue:        "--help",
	RunFunc: func(ctx context.Context, args *clip.CommandArgs[*clip.StdlibExecEnv]) error {
		// Create command line parser
		fset := clip.NewFlagSet(args.CommandName, clip.ExitOnError)
		fset.SetDescription(args.Command.BriefDescription())
		fset.SetArgsDocs("file ...")

		// Add the options
		vflag := fset.Bool("verbose", 'v', false, "verbose mode")

		// Parse command line arguments
		clip.Must(args.Env, fset.Parse(args.Args))

		// Validate number of positional arguments
		clip.Must(args.Env, fset.PositionalArgsRangeCheck(1, math.MaxInt))

		// ...
	},
}

// Create a subcommand working a bit like tar
var tarSubcommand = &clip.LeafCommand[*clip.StdlibExecEnv]{
	BriefDescriptionText: "Archiving utility.",
	HelpFlagValue:        "--help",
	RunFunc: func(ctx context.Context, args *clip.CommandArgs[*clip.StdlibExecEnv]) error {
		// Create command line parser
		fset := clip.NewFlagSet(args.CommandName, clip.ExitOnError)
		fset.SetDescription(args.Command.BriefDescription())
		fset.SetArgsDocs("file ...")

		// Add the options
		cflag := fset.Bool("create", 'c', false, "create a new archive")
		fflag := fset.String("file", 'f', "", "archive file name")
		vflag := fset.Bool("verbose", 'v', false, "verbose mode")
		zflag := fset.Bool("gzip", 'z', false, "gzip compression")

		// Parse command line arguments
		clip.Must(args.Env, fset.Parse(args.Args))

		// Validate number of positional arguments
		clip.Must(args.Env, fset.PositionalArgsRangeCheck(1, math.MaxInt))

		// ...
	},
}

// Create a dispatcher handling control over either tar or gzip
var toolsDispatcher = &clip.DispatcherCommand[*clip.StdlibExecEnv]{
	BriefDescriptionText: "UNIX command-line tools.",
	Commands: map[string]clip.Command[*clip.StdlibExecEnv]{
		"gzip": gzipSubcommand,
		"tar":  tarSubcommand,
	},
}

// Create a root command to wrap it all
var rootCommand = &clip.RootCommand[*clip.StdlibExecEnv]{
	Command: toplevelDispatcher,
}

func main() {
	// Create environment using the standard library I/O
	env := clip.NewStdlibExecEnv()

	// execute the root command
	rootCommand.Main(env)
}
```

See also the testable examples at [example_test.go](example_test.go).

### pkg/flag

See the testable examples at [./pkg/flag/example_test.go](pkg/flag/example_test.go).

### pkg/getopt

See the testable examples at [./pkg/getopt/example_test.go](pkg/getopt/example_test.go).

### pkg/parser

See the testable examples at [./pkg/parser/example_test.go](pkg/parser/example_test.go).

### pkg/scanner

See the testable examples at [./pkg/scanner/example_test.go](pkg/scanner/example_test.go).

## Architecture

The following diagram illustrates the package architecture:

```mermaid
flowchart TD
    assert[pkg/assert]
    clip
    getopt[pkg/getopt]
    flag[pkg/flag]
    parser[pkg/parser]
    scanner[pkg/scanner]
    textwrap[pkg/textwrap]

    clip --> flag
    clip --> assert
    clip --> textwrap
    getopt --> parser
    flag --> parser
    parser --> scanner
    flag --> assert
    getopt --> assert
    parser --> assert
    flag --> textwrap
```

The following subsections illustrate each package.

### clip

[![clip docs](https://pkg.go.dev/badge/github.com/bassosimone/clip)](
https://pkg.go.dev/github.com/bassosimone/clip) [![clip code](
https://img.shields.io/badge/GitHub-pkg/clip-blue?logo=github)](
https://github.com/bassosimone/clip)

Top-level API integrating [flag](./pkg/flag) with subcommands.

### pkg/flag

[![pkg/flag docs](https://pkg.go.dev/badge/github.com/bassosimone/clip/pkg/flag)](
https://pkg.go.dev/github.com/bassosimone/clip/pkg/flag) [![pkg/flag code](
https://img.shields.io/badge/GitHub-pkg/flag-blue?logo=github)](
https://github.com/bassosimone/clip/tree/main/pkg/flag)

[Flag](https://pkg.go.dev/flag) inspired implementation (uses the parser).

### pkg/getopt

[![pkg/getopt docs](https://pkg.go.dev/badge/github.com/bassosimone/clip/pkg/getopt)](
https://pkg.go.dev/github.com/bassosimone/clip/pkg/getopt) [![pkg/getopt code](
https://img.shields.io/badge/GitHub-pkg/getopt-blue?logo=github)](
https://github.com/bassosimone/clip/tree/main/pkg/getopt)

[GNU getopt](https://linux.die.net/man/3/getopt) compatible implementation (uses the parser).

### pkg/parser

[![pkg/parser docs](https://pkg.go.dev/badge/github.com/bassosimone/clip/pkg/parser)](
https://pkg.go.dev/github.com/bassosimone/clip/pkg/parser) [![pkg/parser code](
https://img.shields.io/badge/GitHub-pkg/parser-blue?logo=github)](
https://github.com/bassosimone/clip/tree/main/pkg/parser)

Parser for CLI options (uses the scanner).

### pkg/scanner

[![pkg/scanner docs](https://pkg.go.dev/badge/github.com/bassosimone/clip/pkg/scanner)](
https://pkg.go.dev/github.com/bassosimone/clip/pkg/scanner) [![pkg/scanner code](
https://img.shields.io/badge/GitHub-pkg/scanner-blue?logo=github)](
https://github.com/bassosimone/clip/tree/main/pkg/scanner)

Scanner for CLI options.

### pkg/textwrap

[![pkg/textwrap docs](https://pkg.go.dev/badge/github.com/bassosimone/clip/pkg/textwrap)](
https://pkg.go.dev/github.com/bassosimone/clip/pkg/textwrap) [![pkg/textwrap code](
https://img.shields.io/badge/GitHub-pkg/textwrap-blue?logo=github)](
https://github.com/bassosimone/clip/tree/main/pkg/textwrap)

Utility code to wrap and indent text.

### pkg/assert

[![pkg/assert docs](https://pkg.go.dev/badge/github.com/bassosimone/clip/pkg/assert)](
https://pkg.go.dev/github.com/bassosimone/clip/pkg/assert) [![pkg/assert code](
https://img.shields.io/badge/GitHub-pkg/assert-blue?logo=github)](
https://github.com/bassosimone/clip/tree/main/pkg/assert)

Code to write runtime assertions that panic in case of failure.

## Documentation

Read the package documentation at [pkg.go.dev/github.com/bassosimone/clip](https://pkg.go.dev/github.com/bassosimone/clip)

## Minimum Supported Go Version

Go 1.24

## Installation

```bash
go get -u -v github.com/bassosimone/clip
```

## API Stability Guarantees

This package is experimental and the API may change in the future. Yet,
we will not anticiapte break the existing API without a compelling reason
to do so (e.g., bugs or significant design flaws).

## Running Tests

```
go test -race -count 1 -cover ./...
```

## Dependencies

- [github.com/google/go-cmp](https://pkg.go.dev/github.com/google/go-cmp)
for improving the comparison of structs in unit tests.

## License

```
SPDX-License-Identifier: GPL-3.0-or-later
```
