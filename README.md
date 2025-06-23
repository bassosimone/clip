# clip: Command Line Parser

[![GoDoc](https://pkg.go.dev/badge/github.com/bassosimone/clip)](https://pkg.go.dev/github.com/bassosimone/clip) [![Build Status](https://github.com/bassosimone/clip/actions/workflows/go.yml/badge.svg)](https://github.com/bassosimone/clip/actions) [![codecov](https://codecov.io/gh/bassosimone/clip/branch/main/graph/badge.svg)](https://codecov.io/gh/bassosimone/clip)

This repository implements a very flexible command line parser
written in Go. It provides an intuitive API modeled after the
standard library's [flag](https://pkg.go.dev/flag) package.

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
command line tools behavior in their subcommands. For example, the
[rbmk](https://github.com/rbmk-project/rbmk) network measurement tool.

## Architecture

```
    +--------+-------+
    | getopt | flag  |
    +--------+-------+
    |     parser     |
    +----------------+
    |     scanner    |
    +----------------+
```

- [./pkg/flag](pkg/flag): [flag](https://pkg.go.dev/flag)
inspired implementation (uses the parser).

- [./pkg/getopt](pkg/getopt): getopt compatible implementation (uses the parser).

- [./pkg/parser](pkg/parser): parser for CLI options (uses the scanner).

- [./pkg/scanner](pkg/scanner): scanner for CLI options.

## Examples

The following example shows how to use the [./pkg/flag](pkg/flag) package:

```Go
package main

import (
	"fmt"

	"github.com/bassosimone/clip/pkg/flag"
)

func main() {
	// Create FlagSet
	fset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Define flags with both short and long options
	help := fset.Bool("help", 'h', false, "display help information")
	verbose := fset.Bool("verbose", 'v', false, "enable verbose output")
	output := fset.String("output", 'o', "output.txt", "output file path")
	format := fset.String("format", 'f', "text", "output format")

	// Parse command line arguments (simulated here)
	if err := fset.Parse(os.Args[1:]); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		return
	}

	// Print the values of the flags and arguments
	fmt.Printf("help = %v\n", *help)
	fmt.Printf("verbose = %v\n", *verbose)
	fmt.Printf("output = %v\n", *output)
	fmt.Printf("format = %v\n", *format)
	fmt.Printf("args = %v\n", fset.Args())
}
```

See also the testable examples:

- [./pkg/flag/example_test.go](pkg/flag/example_test.go).

- [./pkg/getopt/example_test.go](pkg/getopt/example_test.go).

- [./pkg/parser/example_test.go](pkg/parser/example_test.go).

- [./pkg/scanner/example_test.go](pkg/scanner/example_test.go).

## Documentation

- [pkg.go.dev/github.com/bassosimone/clip](https://pkg.go.dev/github.com/bassosimone/clip)

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
