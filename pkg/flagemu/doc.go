// doc.go - Package documentation for flagemu
// SPDX-License-Identifier: GPL-3.0-or-later

/*
Package flagemu provides a flag package emulation library that mimics the
standard library's flag package interface while using the generic parser
from the clip project internally.

This package is designed to provide a familiar flag parsing experience
for applications that need to parse command-line arguments in a traditional
GNU-style format while leveraging the flexible parsing capabilities of
the underlying clip parser.

# Basic Usage

	fs := flagemu.NewFlagSet("myprogram", flagemu.ContinueOnError)
	verbose := fs.Bool("verbose", 'v', false, "enable verbose output")
	output := fs.String("output", 'o', "", "output file name")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Verbose: %t\n", *verbose)
	fmt.Printf("Output: %s\n", *output)
	fmt.Printf("Args: %v\n", fs.Args())

# Supported Flag Types

- Bool: Boolean flags that can be specified with or without values

- String: String flags that require a value

The package supports both short (-v) and long (--verbose) option formats,
following GNU getopt conventions. Options can be specified using the equals
syntax (--output=file.txt) or as separate arguments (--output file.txt).

# Error Handling

The package provides configurable error handling through the [ErrorHandling]
type. Currently, only [ContinueOnError] is implemented, which allows parsing to
continue after encountering errors.

# Parser Configuration

By default, the package uses GNU-style parsing with separate short (-v) and
long (--verbose) option prefixes. However, it also supports Go-style parsing
where all options use a single dash prefix:

	// GNU-style parsing (default)
	fs := flagemu.NewFlagSet("myapp", flagemu.ContinueOnError)
	verbose := fs.Bool("verbose", 'v', false, "verbose mode")
	parser := fs.NewParser()  // Creates GNU-style parser

	// Go-style parsing
	fs := flagemu.NewFlagSet("myapp", flagemu.ContinueOnError)
	verbose := fs.Bool("verbose", 'v', false, "verbose mode")
	parser := fs.NewGoStyleParser()  // Creates Go-style parser

The Go-style parser is useful for creating command-line interfaces that
follow Go's flag conventions, where both -v and -verbose would use the
same single-dash prefix.

# Comparison with Standard Flag Package

This package provides a similar API to Go's standard flag package but with
the following differences:

- Supports both short and long option names in a single flag definition

- Uses the clip parser internally for more flexible parsing options

- Supports both GNU-style (default) and Go-style parsing conventions

- Provides more detailed error information through the underlying parser

# Thread Safety

[FlagSet] instances are not safe for concurrent use. Each goroutine should use
its own [FlagSet] instance or external synchronization should be provided.
*/
package flagemu
