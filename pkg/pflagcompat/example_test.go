// example_test.go - Examples
// SPDX-License-Identifier: GPL-3.0-or-later

package pflagcompat_test

import (
	"fmt"
	"os"

	"github.com/bassosimone/clip/pkg/nflag"
	"github.com/bassosimone/clip/pkg/pflagcompat"
)

// This example shows how we print the usage for a curl-like command line.
func ExampleFlagSet_curlHelp() {
	// Create an empty flag set
	fset := pflagcompat.NewFlagSet("curl", nflag.ExitOnError)

	// Make output more pretty by editing default values
	fset.Set.Description = "curl is an utility to transfer URLs.\n"
	fset.Set.Examples = "Examples:\n  curl -fsSL -o index.html https://example.com/\n"
	fset.Set.PositionalArgumentsUsage = "URL ..."
	fset.Set.MinPositionalArgs = 1

	// Add the supported flags using pflag-compatible API
	fset.BoolP("fail", 'f', false, "Fail fast with no output at all on server errors.")
	fset.BoolP("location", 'L', false, "Follow HTTP redirections.")
	fset.Set.AutoHelp("help", 'h', "Show this help message and exit.")
	fset.StringP("output", 'o', "", "Write output to the file indicated by VALUE.")
	fset.BoolP("show-error", 'S', false, "Show an error message, even when silent, on failure.")
	fset.BoolP("silent", 's', false, "Silent or quiet mode.")

	// Override Exit to transform it into a panic
	fset.Set.Exit = func(status int) {
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
	fset := pflagcompat.NewFlagSet("curl", nflag.ExitOnError)

	// Make output more pretty by editing default values
	fset.Set.Description = "curl is an utility to transfer URLs.\n"
	fset.Set.Examples = "Examples:\n  curl -fsSL -o index.html https://example.com/\n"
	fset.Set.PositionalArgumentsUsage = "URL ..."
	fset.Set.MinPositionalArgs = 1

	// Add the supported flags using pflag-compatible API
	fset.BoolP("fail", 'f', false, "Fail fast with no output at all on server errors.")
	fset.BoolP("location", 'L', false, "Follow HTTP redirections.")
	fset.Set.AutoHelp("help", 'h', "Show this help message and exit.")
	fset.StringP("output", 'o', "", "Write output to the file indicated by VALUE.")
	fset.BoolP("show-error", 'S', false, "Show an error message, even when silent, on failure.")
	fset.BoolP("silent", 's', false, "Silent or quiet mode.")

	// Override Exit to transform it into a panic
	fset.Set.Exit = func(status int) {
		panic("mocked exit invocation")
	}

	// Override Stderr to be the Stdout otherwise the testable example fails
	fset.Set.Stderr = os.Stdout

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
	fset := pflagcompat.NewFlagSet("curl", nflag.ExitOnError)

	// Make output more pretty by editing default values
	fset.Set.Description = "curl is an utility to transfer URLs.\n"
	fset.Set.Examples = "Examples:\n  curl -fsSL -o index.html https://example.com/\n"
	fset.Set.PositionalArgumentsUsage = "URL ..."
	fset.Set.MinPositionalArgs = 1

	// Add the supported flags using pflag-compatible API
	ffail := fset.BoolP("fail", 'f', false, "Fail fast with no output at all on server errors.")
	flocation := fset.BoolP("location", 'L', false, "Follow HTTP redirections.")
	fmaxFilesize := fset.Int64P("max-filesize", 'm', 0, "Fail if the file is larger than VALUE bytes.")
	fset.Set.AutoHelp("help", 'h', "Show this help message and exit.")
	foutput := fset.StringP("output", 'o', "", "Write output to the file indicated by VALUE.")
	fshowError := fset.BoolP("show-error", 'S', false, "Show an error message, even when silent, on failure.")
	fsilent := fset.BoolP("silent", 's', false, "Silent or quiet mode.")

	// Invoke with command line arguments
	//
	// Note that `-fsSL` is equivalent to [`-f`, `-s`, `-S`, `-L`].
	//
	// Note that `-m1024` is equivalent to [`-m`, `1024`].
	fset.Parse([]string{"-fsSL", "-m1024", "-o", "index.html", "https://example.com/"})

	// Print the parsed flags
	fmt.Printf("fail: %v\n", *ffail)
	fmt.Printf("location: %v\n", *flocation)
	fmt.Printf("max-filesize: %v\n", *fmaxFilesize)
	fmt.Printf("output: %s\n", *foutput)
	fmt.Printf("show-error: %v\n", *fshowError)
	fmt.Printf("silent: %v\n", *fsilent)
	fmt.Printf("args: %v\n", fset.Set.Args())

	// Output:
	// fail: true
	// location: true
	// max-filesize: 1024
	// output: index.html
	// show-error: true
	// silent: true
	// args: [https://example.com/]
}

// This example demonstrates using the Var variants with default values.
func ExampleFlagSet_curlWithVar() {
	// Create an empty flag set
	fset := pflagcompat.NewFlagSet("curl", nflag.ExitOnError)

	// Make output more pretty by editing default values
	fset.Set.Description = "curl is an utility to transfer URLs.\n"
	fset.Set.PositionalArgumentsUsage = "URL ..."
	fset.Set.MinPositionalArgs = 1

	// Define variables and add flags using Var variants
	var (
		fail      bool
		location  bool
		output    string
		silent    bool
		maxSize   int64
	)

	fset.BoolVarP(&fail, "fail", 'f', false, "Fail fast with no output at all on server errors.")
	fset.BoolVarP(&location, "location", 'L', false, "Follow HTTP redirections.")
	fset.StringVarP(&output, "output", 'o', "", "Write output to the file indicated by VALUE.")
	fset.BoolVarP(&silent, "silent", 's', false, "Silent or quiet mode.")
	fset.Int64VarP(&maxSize, "max-filesize", 'm', 2048, "Fail if the file is larger than VALUE bytes.")

	// Invoke with command line arguments
	fset.Parse([]string{"-fsL", "-o", "page.html", "https://example.com/"})

	// Print the parsed flags
	fmt.Printf("fail: %v\n", fail)
	fmt.Printf("location: %v\n", location)
	fmt.Printf("output: %s\n", output)
	fmt.Printf("silent: %v\n", silent)
	fmt.Printf("max-filesize: %v\n", maxSize)
	fmt.Printf("args: %v\n", fset.Set.Args())

	// Output:
	// fail: true
	// location: true
	// output: page.html
	// silent: true
	// max-filesize: 2048
	// args: [https://example.com/]
}

// This example demonstrates using the non-P variants (long-only flags).
func ExampleFlagSet_longOnlyFlags() {
	// Create an empty flag set
	fset := pflagcompat.NewFlagSet("tool", nflag.ExitOnError)

	// Make output more pretty by editing default values
	fset.Set.Description = "A tool with long-only flags.\n"
	fset.Set.PositionalArgumentsUsage = "FILE ..."
	fset.Set.MinPositionalArgs = 0

	// Add long-only flags
	verbose := fset.Bool("verbose", false, "Enable verbose output.")
	configFile := fset.String("config", "default.conf", "Configuration file path.")
	maxRetries := fset.Int64("max-retries", 3, "Maximum number of retries.")

	// Invoke with command line arguments
	fset.Parse([]string{"--verbose", "--config", "custom.conf", "file1.txt", "file2.txt"})

	// Print the parsed flags
	fmt.Printf("verbose: %v\n", *verbose)
	fmt.Printf("config: %s\n", *configFile)
	fmt.Printf("max-retries: %v\n", *maxRetries)
	fmt.Printf("args: %v\n", fset.Set.Args())

	// Output:
	// verbose: true
	// config: custom.conf
	// max-retries: 3
	// args: [file1.txt file2.txt]
}

// This example demonstrates using the Var variants with long-only flags.
func ExampleFlagSet_longOnlyVarFlags() {
	// Create an empty flag set
	fset := pflagcompat.NewFlagSet("tool", nflag.ExitOnError)

	// Make output more pretty by editing default values
	fset.Set.Description = "A tool with long-only Var flags.\n"
	fset.Set.PositionalArgumentsUsage = "FILE ..."
	fset.Set.MinPositionalArgs = 0

	// Define variables and add long-only flags using Var variants
	var (
		verbose    bool
		configFile string
		maxRetries int64
	)

	fset.BoolVar(&verbose, "verbose", false, "Enable verbose output.")
	fset.StringVar(&configFile, "config", "default.conf", "Configuration file path.")
	fset.Int64Var(&maxRetries, "max-retries", 3, "Maximum number of retries.")

	// Invoke with command line arguments
	fset.Parse([]string{"--verbose", "--config", "custom.conf", "file1.txt", "file2.txt"})

	// Print the parsed flags
	fmt.Printf("verbose: %v\n", verbose)
	fmt.Printf("config: %s\n", configFile)
	fmt.Printf("max-retries: %v\n", maxRetries)
	fmt.Printf("args: %v\n", fset.Set.Args())

	// Output:
	// verbose: true
	// config: custom.conf
	// max-retries: 3
	// args: [file1.txt file2.txt]
}
