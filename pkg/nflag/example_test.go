// example_test.go - Examples
// SPDX-License-Identifier: GPL-3.0-or-later

package nflag_test

import (
	"fmt"
	"os"

	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/nflag"
	"github.com/bassosimone/clip/pkg/nparser"
)

// This example shows the behavior when no flags are defined.
func ExampleFlagSet_noFlags() {
	// Create an empty flag set
	fset := nflag.NewFlagSet("curl", nflag.ExitOnError)

	// Make output more pretty by editing default values
	fset.Description = "curl is an utility to transfer URLs.\n"
	fset.Examples = "Examples:\n  curl -fsSL -o index.html https://example.com/\n"
	fset.PositionalArgumentsUsage = "URL ..."
	fset.MinPositionalArgs = 1

	// Note: no flags have been configured

	// Override Exit to transform it into a panic
	fset.Exit = func(status int) {
		panic("mocked exit invocation")
	}

	// Override Stderr to be the Stdout otherwise the testable example fails
	fset.Stderr = os.Stdout

	// Handle the panic by caused by Exit by simply ignoring it
	defer func() { recover() }()

	// Invoke with `--verbose` which yields an error because `verbose` is not defined
	fset.Parse([]string{"--verbose"})

	// Output:
	// curl: unknown option: --verbose
}

// This example shows how we print the usage for a curl-like command line.
func ExampleFlagSet_curlHelp() {
	// Create an empty flag set
	fset := nflag.NewFlagSet("curl", nflag.ExitOnError)

	// Make output more pretty by editing default values
	fset.Description = "curl is an utility to transfer URLs.\n"
	fset.Examples = "Examples:\n  curl -fsSL -o index.html https://example.com/\n"
	fset.PositionalArgumentsUsage = "URL ..."
	fset.MinPositionalArgs = 1

	// Add the supported flags
	fset.BoolFlag("fail", 'f', "Fail fast with no output at all on server errors.")
	fset.BoolFlag("location", 'L', "Follow HTTP redirections.")
	fset.AutoHelp("help", 'h', "Show this help message and exit.")
	fset.StringFlag("output", 'o', "Write output to the file indicated by VALUE.")
	fset.BoolFlag("show-error", 'S', "Show an error message, even when silent, on failure.")
	fset.BoolFlag("silent", 's', "Silent or quiet mode.")

	// Override Exit to transform it into a panic
	fset.Exit = func(status int) {
		panic("mocked exit invocation")
	}

	// Handle the panic by caused by Exit by simply ignoring it
	defer func() { recover() }()

	// Invoke with `--help`
	fset.Parse([]string{"--help"})

	// Output:
	// Usage: curl [options] URL ...
	//
	// curl is an utility to transfer URLs.
	//
	// Options:
	//   -f, --fail
	//     Fail fast with no output at all on server errors.
	//
	//   -L, --location
	//     Follow HTTP redirections.
	//
	//   -h, --help
	//     Show this help message and exit.
	//
	//   -o, --output=VALUE
	//     Write output to the file indicated by VALUE.
	//
	//   -S, --show-error
	//     Show an error message, even when silent, on failure.
	//
	//   -s, --silent
	//     Silent or quiet mode.
	//
	// Examples:
	//   curl -fsSL -o index.html https://example.com/
}

// This example shows how we print errors when there are too few arguments.
func ExampleFlagSet_curlTooFewArguments() {
	// Create an empty flag set
	fset := nflag.NewFlagSet("curl", nflag.ExitOnError)

	// Make output more pretty by editing default values
	fset.Description = "curl is an utility to transfer URLs.\n"
	fset.Examples = "Examples:\n  curl -fsSL -o index.html https://example.com/\n"
	fset.PositionalArgumentsUsage = "URL ..."
	fset.MinPositionalArgs = 1

	// Add the supported flags
	fset.BoolFlag("fail", 'f', "Fail fast with no output at all on server errors.")
	fset.BoolFlag("location", 'L', "Follow HTTP redirections.")
	fset.AutoHelp("help", 'h', "Show this help message and exit.")
	fset.StringFlag("output", 'o', "Write output to the file indicated by VALUE.")
	fset.BoolFlag("show-error", 'S', "Show an error message, even when silent, on failure.")
	fset.BoolFlag("silent", 's', "Silent or quiet mode.")

	// Override Exit to transform it into a panic
	fset.Exit = func(status int) {
		panic("mocked exit invocation")
	}

	// Override Stderr to be the Stdout otherwise the testable example fails
	fset.Stderr = os.Stdout

	// Handle the panic by caused by Exit by simply ignoring it
	defer func() { recover() }()

	// Invoked with not arguments at all
	fset.Parse([]string{})

	// Output:
	// curl: too few positional arguments: expected at least 1, got 0
	// Try 'curl --help' for more help.
}

// This example shows a successful invocation of a curl-like tool.
func ExampleFlagSet_curlSuccess() {
	// Create an empty flag set
	fset := nflag.NewFlagSet("curl", nflag.ExitOnError)

	// Make output more pretty by editing default values
	fset.Description = "curl is an utility to transfer URLs.\n"
	fset.Examples = "Examples:\n  curl -fsSL -o index.html https://example.com/\n"
	fset.PositionalArgumentsUsage = "URL ..."
	fset.MinPositionalArgs = 1

	// Add the supported flags
	ffail := fset.BoolFlag("fail", 'f', "Fail fast with no output at all on server errors.")
	flocation := fset.BoolFlag("location", 'L', "Follow HTTP redirections.")
	fmaxFilesize := fset.Int64Flag("max-filesize", 'm', "Fail if the file is larger than VALUE bytes.")
	fset.AutoHelp("help", 'h', "Show this help message and exit.")
	foutput := fset.StringFlag("output", 'o', "Write output to the file indicated by VALUE.")
	fshowError := fset.BoolFlag("show-error", 'S', "Show an error message, even when silent, on failure.")
	fsilent := fset.BoolFlag("silent", 's', "Silent or quiet mode.")

	// Invoke with command line arguments
	//
	// Note that `-fsSL` is equivalent to [`-f`, `-s`, `-S`, `-L`].
	//
	// Note that `-m1024` is equivalent to [`-m`, `1024`].
	fset.Parse([]string{"-fsSL", "-m1024", "-o", "index.html", "https://example.com/"})

	// Print the parsed flags
	fmt.Println("---")
	fmt.Printf("fail: %v\n", *ffail)
	fmt.Printf("location: %v\n", *flocation)
	fmt.Printf("max-filesize: %v\n", *fmaxFilesize)
	fmt.Printf("output: %s\n", *foutput)
	fmt.Printf("show-error: %v\n", *fshowError)
	fmt.Printf("silent: %v\n", *fsilent)

	// Same as above but by accessing the flags directly
	fmt.Println("---")
	for _, pair := range fset.Flags() {
		assert.True(pair.LongFlag != nil && pair.ShortFlag != nil, "expected both to be not nil")
		fmt.Printf("%s %s: %s\n", pair.LongFlag.Option.Name,
			pair.ShortFlag.Option.Name, pair.Value.String())
	}
	fmt.Println("---")

	// Print the positional arguments
	fmt.Printf("positional arguments: %v\n", fset.Args())

	// Output:
	// ---
	// fail: true
	// location: true
	// max-filesize: 1024
	// output: index.html
	// show-error: true
	// silent: true
	// ---
	// fail f: true
	// location L: true
	// max-filesize m: 1024
	// help h: false
	// output o: index.html
	// show-error S: true
	// silent s: true
	// ---
	// positional arguments: [https://example.com/]
}

// This example shows how we print the usage for a dig-like tool.
func ExampleFlagSet_digHelp() {
	// Create an empty flag set
	fset := nflag.NewFlagSet("dig", nflag.ExitOnError)

	// Make output more pretty by editing default values
	fset.Description = "dig is an utility to query the domain name system.\n"
	fset.Examples = "Examples:\n  dig +short IN A -46 example.com +https\n"
	fset.PositionalArgumentsUsage = "[@server] [name] [type] [class]"
	fset.MaxPositionalArgs = 4

	// Modify the long prefix to use dig conventions
	fset.LongFlagPrefix = "+"

	// Add the supported flags
	fset.BoolFlag("", '4', "Enable using IPv4.")
	fset.BoolFlag("", '6', "Enable using IPv6.")
	fset.AutoHelp("", 'h', "Show this help message and exit.")
	fset.StringFlag("https", 0, "Use DNS-over-HTTPS optionally setting URL path to VALUE.")
	fset.BoolFlag("short", 0, "Write terse output.")

	// Modify the "https" flag to *optionally* accept a value
	assert.True1(fset.LookupFlagLong("https")).Option.Type = nparser.OptionTypeStandaloneArgumentOptional

	// Override Exit to transform it into a panic
	fset.Exit = func(status int) {
		panic("mocked exit invocation")
	}

	// Handle the panic by caused by Exit by simply ignoring it
	defer func() { recover() }()

	// Invoke with `-h`
	fset.Parse([]string{"-h"})

	// Output:
	// Usage: dig [options] [@server] [name] [type] [class]
	//
	// dig is an utility to query the domain name system.
	//
	// Options:
	//   -4
	//     Enable using IPv4.
	//
	//   -6
	//     Enable using IPv6.
	//
	//   -h
	//     Show this help message and exit.
	//
	//   +https=VALUE
	//     Use DNS-over-HTTPS optionally setting URL path to VALUE.
	//
	//   +short
	//     Write terse output.
	//
	// Examples:
	//   dig +short IN A -46 example.com +https
}

// This example shows how we print errors caused by invalid flags.
func ExampleFlagSet_digInvalidFlag() {
	// Create an empty flag set
	fset := nflag.NewFlagSet("dig", nflag.ExitOnError)

	// Make output more pretty by editing default values
	fset.Description = "dig is an utility to query the domain name system.\n"
	fset.Examples = "Examples:\n  dig +short IN A -46 example.com +https\n"
	fset.PositionalArgumentsUsage = "[@server] [name] [type] [class]"
	fset.MaxPositionalArgs = 4

	// Modify the long prefix to use dig conventions
	fset.LongFlagPrefix = "+"

	// Add the supported flags
	fset.BoolFlag("", '4', "Enable using IPv4.")
	fset.BoolFlag("", '6', "Enable using IPv6.")
	fset.AutoHelp("", 'h', "Show this help message and exit.")
	fset.StringFlag("https", 0, "Use DNS-over-HTTPS optionally setting URL path to VALUE.")
	fset.BoolFlag("short", 0, "Write terse output.")

	// Modify the "https" flag to *optionally* accept a value
	assert.True1(fset.LookupFlagLong("https")).Option.Type = nparser.OptionTypeStandaloneArgumentOptional

	// Override Exit to transform it into a panic
	fset.Exit = func(status int) {
		panic("mocked exit invocation")
	}

	// Override Stderr to be the Stdout otherwise the testable example fails
	fset.Stderr = os.Stdout

	// Handle the panic by caused by Exit by simply ignoring it
	defer func() { recover() }()

	// Invoke with a flag that has not been defined
	fset.Parse([]string{"+tls"})

	// Output:
	// dig: unknown option: +tls
	// Try 'dig -h' for more help.
}

// This example shows how we print the usage for a tar-like tool.
func ExampleFlagSet_tarHelp() {
	// Create an empty flag set
	fset := nflag.NewFlagSet("tar", nflag.ExitOnError)

	// Make output more pretty by editing default values
	fset.Description = "tar is an utility to manage possibly-compressed archives.\n"
	fset.Examples = "Examples:\n  tar -cvzf archive.tar.gz file1.txt file2.txt file3.txt\n"
	fset.PositionalArgumentsUsage = "[FILE ...]"

	// Add the supported flags
	fset.BoolFlag("", 'c', "Create a new archive.")
	fset.StringFlag("", 'f', "Specify the output file path VALUE.")
	fset.AutoHelp("help", 'h', "Show this help message and exit.")
	fset.BoolFlag("", 'v', "Print files added to the archive to the stdout.")
	fset.BoolFlag("", 'z', "Compress using gzip.")

	// Override Exit to transform it into a panic
	fset.Exit = func(status int) {
		panic("mocked exit invocation")
	}

	// Handle the panic by caused by Exit by simply ignoring it
	defer func() { recover() }()

	// Invoke with `--help`
	fset.Parse([]string{"--help"})

	// Output:
	// Usage: tar [options] [FILE ...]
	//
	// tar is an utility to manage possibly-compressed archives.
	//
	// Options:
	//   -c
	//     Create a new archive.
	//
	//   -f VALUE
	//     Specify the output file path VALUE.
	//
	//   -h, --help
	//     Show this help message and exit.
	//
	//   -v
	//     Print files added to the archive to the stdout.
	//
	//   -z
	//     Compress using gzip.
	//
	// Examples:
	//   tar -cvzf archive.tar.gz file1.txt file2.txt file3.txt
}

// This example shows how we print errors caused by a missing mandatory argument.
func ExampleFlagSet_tarMissingOptionArgument() {
	// Create an empty flag set
	fset := nflag.NewFlagSet("tar", nflag.ExitOnError)

	// Make output more pretty by editing default values
	fset.Description = "tar is an utility to manage possibly-compressed archives.\n"
	fset.Examples = "Examples:\n  tar -cvzf archive.tar.gz file1.txt file2.txt file3.txt\n"
	fset.PositionalArgumentsUsage = "[FILE ...]"

	// Add the supported flags
	fset.BoolFlag("", 'c', "Create a new archive.")
	fset.StringFlag("", 'f', "Specify the output file path VALUE.")
	fset.AutoHelp("help", 'h', "Show this help message and exit.")
	fset.BoolFlag("", 'v', "Print files added to the archive to the stdout.")
	fset.BoolFlag("", 'z', "Compress using gzip.")

	// Override Exit to transform it into a panic
	fset.Exit = func(status int) {
		panic("mocked exit invocation")
	}

	// Override Stderr to be the Stdout otherwise the testable example fails
	fset.Stderr = os.Stdout

	// Handle the panic by caused by Exit by simply ignoring it
	defer func() { recover() }()

	// Invoke missing argument for the `-f` option
	fset.Parse([]string{"-f"})

	// Output:
	// tar: option requires an argument: -f
	// Try 'tar --help' for more help.
}

// This example shows how we print the usage for a go-like tool.
func ExampleFlagSet_goHelp() {
	// Create an empty flag set
	fset := nflag.NewFlagSet("go test", nflag.ExitOnError)

	// Make output more pretty by editing default values
	fset.Description = "go test runs package tests.\n"
	fset.Examples = "Examples:\n  go test -race -count=1 -v ./...\n"
	fset.PositionalArgumentsUsage = "[package ...]"

	// Disable using short flags and use `-` for long flags
	fset.ShortFlagPrefix = ""
	fset.LongFlagPrefix = "-"

	// Add the supported flags
	fset.Int64Flag("count", 0, "Set VALUE to 1 to avoid using the test cache.")
	fset.AutoHelp("h", 0, "Show this help message and exit.")
	fset.AutoHelp("help", 0, "Alias for -h.")
	fset.BoolFlag("race", 0, "Run tests using the race detector.")
	fset.BoolFlag("v", 0, "Print details about the tests progress and results.")

	// Override Exit to transform it into a panic
	fset.Exit = func(status int) {
		panic("mocked exit invocation")
	}

	// Handle the panic by caused by Exit by simply ignoring it
	defer func() { recover() }()

	// Invoke with `-help`
	fset.Parse([]string{"-help"})

	// Output:
	// Usage: go test [options] [package ...]
	//
	// go test runs package tests.
	//
	// Options:
	//   -count=VALUE
	//     Set VALUE to 1 to avoid using the test cache.
	//
	//   -h
	//     Show this help message and exit.
	//
	//   -help
	//     Alias for -h.
	//
	//   -race
	//     Run tests using the race detector.
	//
	//   -v
	//     Print details about the tests progress and results.
	//
	// Examples:
	//   go test -race -count=1 -v ./...
}

// This example shows a successful invocation of a go-like tool.
func ExampleFlagSet_goSuccess() {
	// Create an empty flag set
	fset := nflag.NewFlagSet("go test", nflag.ExitOnError)

	// Make output more pretty by editing default values
	fset.Description = "go test runs package tests.\n"
	fset.Examples = "Examples:\n  go test -race -count=1 -v ./...\n"
	fset.PositionalArgumentsUsage = "[package ...]"

	// Disable using short flags and use `-` for long flags
	fset.ShortFlagPrefix = ""
	fset.LongFlagPrefix = "-"

	// Add the supported flags
	fcount := fset.Int64Flag("count", 0, "Set VALUE to 1 to avoid using the test cache.")
	fset.AutoHelp("h", 0, "Show this help message and exit.")
	fset.AutoHelp("help", 0, "Alias for -h.")
	frace := fset.BoolFlag("race", 0, "Run tests using the race detector.")
	fv := fset.BoolFlag("v", 0, "Print details about the tests progress and results.")

	// Invoke with command line arguments.
	//
	// Note that `-count=1` is equivalent to [`-count`, `1`].
	fset.Parse([]string{"-race", "-count", "1", "-v", "./..."})

	// Print the parsed flags
	fmt.Println("---")
	fmt.Printf("fcount: %v\n", *fcount)
	fmt.Printf("race: %v\n", *frace)
	fmt.Printf("v: %v\n", *fv)

	// Same as above but by accessing the flags directly
	fmt.Println("---")
	for _, pair := range fset.Flags() {
		assert.True(
			pair.LongFlag != nil && pair.ShortFlag == nil,
			"expected LongFlag to be not nil and ShortFlag to be nil",
		)
		fmt.Printf("%s: %s\n", pair.LongFlag.Option.Name, pair.Value.String())
	}
	fmt.Println("---")

	// Print the positional arguments
	fmt.Printf("positional arguments: %v\n", fset.Args())

	// Output:
	// ---
	// fcount: 1
	// race: true
	// v: true
	// ---
	// count: 1
	// h: false
	// help: false
	// race: true
	// v: true
	// ---
	// positional arguments: [./...]
}

// This example shows how we print errors caused by strconv failures.
func ExampleFlagSet_goStrconvError() {
	// Create an empty flag set
	fset := nflag.NewFlagSet("go test", nflag.ExitOnError)

	// Make output more pretty by editing default values
	fset.Description = "go test runs package tests.\n"
	fset.Examples = "Examples:\n  go test -race -count=1 -v ./...\n"
	fset.PositionalArgumentsUsage = "[package ...]"

	// Disable using short flags and use `-` for long flags
	fset.ShortFlagPrefix = ""
	fset.LongFlagPrefix = "-"

	// Add the supported flags
	fset.Int64Flag("count", 0, "Set VALUE to 1 to avoid using the test cache.")
	fset.AutoHelp("h", 0, "Show this help message and exit.")
	fset.AutoHelp("help", 0, "Alias for -h.")
	fset.BoolFlag("race", 0, "Run tests using the race detector.")
	fset.BoolFlag("v", 0, "Print details about the tests progress and results.")

	// Override Exit to transform it into a panic
	fset.Exit = func(status int) {
		panic("mocked exit invocation")
	}

	// Override Stderr to be the Stdout otherwise the testable example fails
	fset.Stderr = os.Stdout

	// Handle the panic by caused by Exit by simply ignoring it
	defer func() { recover() }()

	// Invoke with `-help`
	fset.Parse([]string{"-count", "not-a-number"})

	// Output:
	// go test: strconv.ParseInt: parsing "not-a-number": invalid syntax
	// Try 'go test -h' for more help.
}
