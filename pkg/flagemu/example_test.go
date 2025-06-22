// example_test.go - Comprehensive examples demonstrating flagemu package usage
// SPDX-License-Identifier: GPL-3.0-or-later

package flagemu_test

import (
	"fmt"
	"log"

	"github.com/bassosimone/clip/pkg/flagemu"
	"github.com/bassosimone/clip/pkg/parser"
)

// Example_basic demonstrates the basic usage of the flagemu package
// with boolean and string flags.
func Example_basic() {
	// Create a new flag set
	fs := flagemu.NewFlagSet("myapp", flagemu.ContinueOnError)

	// Define flags with both long and short names
	verbose := fs.Bool("verbose", 'v', false, "enable verbose output")
	output := fs.String("output", 'o', "stdout", "output destination")
	debug := fs.Bool("debug", 'd', false, "enable debug mode")

	// Example command line: myapp -v --output=file.txt --debug input.txt extra.txt
	args := []string{"-v", "--output=file.txt", "--debug", "input.txt", "extra.txt"}

	// Parse the arguments
	err := fs.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	// Use the parsed flags
	fmt.Printf("Verbose: %t\n", *verbose)
	fmt.Printf("Output: %s\n", *output)
	fmt.Printf("Debug: %t\n", *debug)
	fmt.Printf("Args: %v\n", fs.Args())

	// Output:
	// Verbose: true
	// Output: file.txt
	// Debug: true
	// Args: [input.txt extra.txt]
}

// Example_shortFlags demonstrates using short flags with bundling-like syntax.
func Example_shortFlags() {
	fs := flagemu.NewFlagSet("tool", flagemu.ContinueOnError)

	// Define flags with short names
	verbose := fs.Bool("verbose", 'v', false, "verbose output")
	all := fs.Bool("all", 'a', false, "show all files")
	human := fs.Bool("human", 'h', false, "human readable")
	file := fs.String("file", 'f', "", "input file")

	// Parse short flags: tool -v -a -h -f config.txt
	args := []string{"-v", "-a", "-h", "-f", "config.txt"}

	err := fs.Parse(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Verbose: %t\n", *verbose)
	fmt.Printf("All: %t\n", *all)
	fmt.Printf("Human: %t\n", *human)
	fmt.Printf("File: %s\n", *file)

	// Output:
	// Verbose: true
	// All: true
	// Human: true
	// File: config.txt
}

// Example_longFlags demonstrates using long flags with various syntaxes.
func Example_longFlags() {
	fs := flagemu.NewFlagSet("server", flagemu.ContinueOnError)

	// Define flags with long names
	port := fs.String("port", 'p', "8080", "server port")
	host := fs.String("host", 'h', "localhost", "server host")
	ssl := fs.Bool("ssl", 's', false, "enable SSL")
	config := fs.String("config", 'c', "", "configuration file")

	// Parse long flags with different syntaxes
	args := []string{
		"--port=9000",       // equals syntax
		"--host", "0.0.0.0", // separate argument
		"--ssl",                  // boolean flag
		"--config=/etc/app.conf", // equals syntax with path
	}

	err := fs.Parse(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Port: %s\n", *port)
	fmt.Printf("Host: %s\n", *host)
	fmt.Printf("SSL: %t\n", *ssl)
	fmt.Printf("Config: %s\n", *config)

	// Output:
	// Port: 9000
	// Host: 0.0.0.0
	// SSL: true
	// Config: /etc/app.conf
}

// Example_mixedFlags demonstrates mixing short and long flags in the same command.
func Example_mixedFlags() {
	fs := flagemu.NewFlagSet("backup", flagemu.ContinueOnError)

	// Define flags that can be used with both short and long names
	verbose := fs.Bool("verbose", 'v', false, "verbose output")
	dryRun := fs.Bool("dry-run", 'n', false, "perform a trial run")
	exclude := fs.String("exclude", 'e', "", "exclude pattern")
	target := fs.String("target", 't', "", "backup target")

	// Mix of short and long flags
	args := []string{
		"-v",          // short boolean
		"--dry-run",   // long boolean
		"-e", "*.tmp", // short string with separate argument
		"--target=/backup/dir", // long string with equals
		"/home/user",           // positional argument
	}

	err := fs.Parse(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Verbose: %t\n", *verbose)
	fmt.Printf("Dry Run: %t\n", *dryRun)
	fmt.Printf("Exclude: %s\n", *exclude)
	fmt.Printf("Target: %s\n", *target)
	fmt.Printf("Source: %v\n", fs.Args())

	// Output:
	// Verbose: true
	// Dry Run: true
	// Exclude: *.tmp
	// Target: /backup/dir
	// Source: [/home/user]
}

// Example_defaultValues demonstrates how default values work.
func Example_defaultValues() {
	fs := flagemu.NewFlagSet("config", flagemu.ContinueOnError)

	// Define flags with default values
	timeout := fs.String("timeout", 't', "30s", "connection timeout")
	retries := fs.String("retries", 'r', "3", "number of retries")
	enabled := fs.Bool("enabled", 'e', true, "feature enabled")
	disabled := fs.Bool("disabled", 'd', false, "feature disabled")

	// Parse with no arguments to show defaults
	args := []string{}

	err := fs.Parse(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Timeout: %s\n", *timeout)
	fmt.Printf("Retries: %s\n", *retries)
	fmt.Printf("Enabled: %t\n", *enabled)
	fmt.Printf("Disabled: %t\n", *disabled)

	// Parse with some arguments to show overrides
	fmt.Println("\nWith arguments:")

	fs2 := flagemu.NewFlagSet("config", flagemu.ContinueOnError)
	timeout2 := fs2.String("timeout", 't', "30s", "connection timeout")
	retries2 := fs2.String("retries", 'r', "3", "number of retries")
	enabled2 := fs2.Bool("enabled", 'e', true, "feature enabled")
	disabled2 := fs2.Bool("disabled", 'd', false, "feature disabled")

	args2 := []string{"--timeout=60s", "--disabled"}
	err = fs2.Parse(args2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Timeout: %s\n", *timeout2)
	fmt.Printf("Retries: %s\n", *retries2)
	fmt.Printf("Enabled: %t\n", *enabled2)
	fmt.Printf("Disabled: %t\n", *disabled2)

	// Output:
	// Timeout: 30s
	// Retries: 3
	// Enabled: true
	// Disabled: false
	//
	// With arguments:
	// Timeout: 60s
	// Retries: 3
	// Enabled: true
	// Disabled: true
}

// Example_noPermute demonstrates the SetNoPermute option.
func Example_noPermute() {
	fmt.Println("Default behavior (with permutation):")
	fs1 := flagemu.NewFlagSet("app", flagemu.ContinueOnError)
	verbose1 := fs1.Bool("verbose", 'v', false, "verbose output")

	// Arguments mixed with flags
	args1 := []string{"file1.txt", "-v", "file2.txt"}
	err := fs1.Parse(args1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Verbose: %t\n", *verbose1)
	fmt.Printf("Args: %v\n", fs1.Args())

	fmt.Println("\nWith no permutation:")
	fs2 := flagemu.NewFlagSet("app", flagemu.ContinueOnError)
	fs2.SetNoPermute() // Disable argument permutation
	verbose2 := fs2.Bool("verbose", 'v', false, "verbose output")

	// Same arguments, but permutation disabled
	args2 := []string{"file1.txt", "-v", "file2.txt"}
	err = fs2.Parse(args2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Verbose: %t\n", *verbose2)
	fmt.Printf("Args: %v\n", fs2.Args())

	// Output:
	// Default behavior (with permutation):
	// Verbose: true
	// Args: [file1.txt file2.txt]
	//
	// With no permutation:
	// Verbose: false
	// Args: [file1.txt -v file2.txt]
}

// Example_realWorld demonstrates a real-world usage scenario
// similar to common command-line tools.
func Example_realWorld() {
	fs := flagemu.NewFlagSet("fileprocessor", flagemu.ContinueOnError)

	// Define flags similar to common tools like grep, find, etc.
	recursive := fs.Bool("recursive", 'r', false, "process directories recursively")
	ignoreCase := fs.Bool("ignore-case", 'i', false, "ignore case distinctions")
	pattern := fs.String("pattern", 'p', "", "search pattern")
	output := fs.String("output", 'o', "", "output file (default: stdout)")
	verbose := fs.Bool("verbose", 'v', false, "verbose output")
	count := fs.Bool("count", 'c', false, "only show count of matches")

	// Simulate a complex command line
	args := []string{
		"-r",               // recursive
		"--ignore-case",    // ignore case
		"-p", "TODO|FIXME", // pattern
		"--output=results.txt", // output file
		"-v",                   // verbose
		"src/",                 // directory to process
		"docs/",                // another directory
	}

	err := fs.Parse(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Display parsed configuration
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Recursive: %t\n", *recursive)
	fmt.Printf("  Ignore Case: %t\n", *ignoreCase)
	fmt.Printf("  Pattern: %s\n", *pattern)
	fmt.Printf("  Output: %s\n", *output)
	fmt.Printf("  Verbose: %t\n", *verbose)
	fmt.Printf("  Count Only: %t\n", *count)
	fmt.Printf("  Directories: %v\n", fs.Args())

	// Simulate processing
	fmt.Printf("\nProcessing...\n")
	for _, dir := range fs.Args() {
		if *verbose {
			fmt.Printf("Processing directory: %s\n", dir)
		}
		if *recursive {
			fmt.Printf("  Searching recursively in %s\n", dir)
		}
		if *ignoreCase {
			fmt.Printf("  Using case-insensitive search\n")
		}
		fmt.Printf("  Looking for pattern: %s\n", *pattern)
	}

	if *output != "" {
		fmt.Printf("Results will be saved to: %s\n", *output)
	}

	// Output:
	// Configuration:
	//   Recursive: true
	//   Ignore Case: true
	//   Pattern: TODO|FIXME
	//   Output: results.txt
	//   Verbose: true
	//   Count Only: false
	//   Directories: [src/ docs/]
	//
	// Processing...
	// Processing directory: src/
	//   Searching recursively in src/
	//   Using case-insensitive search
	//   Looking for pattern: TODO|FIXME
	// Processing directory: docs/
	//   Searching recursively in docs/
	//   Using case-insensitive search
	//   Looking for pattern: TODO|FIXME
	// Results will be saved to: results.txt
}

// Example_errorHandling demonstrates error handling in flag parsing.
func Example_errorHandling() {
	fs := flagemu.NewFlagSet("errortest", flagemu.ContinueOnError)

	// Define some flags
	verbose := fs.Bool("verbose", 'v', false, "verbose output")
	file := fs.String("file", 'f', "", "input file")

	fmt.Println("Valid arguments:")
	args1 := []string{"-v", "--file=test.txt"}
	err := fs.Parse(args1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success - Verbose: %t, File: %s\n", *verbose, *file)
	}

	fmt.Println("\nInvalid flag:")
	fs2 := flagemu.NewFlagSet("errortest", flagemu.ContinueOnError)
	fs2.Bool("verbose", 'v', false, "verbose output")

	args2 := []string{"--unknown-flag"}
	err = fs2.Parse(args2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Output:
	// Valid arguments:
	// Success - Verbose: true, File: test.txt
	//
	// Invalid flag:
	// Error: unknown option: unknown-flag
}

// Example_customParser demonstrates how to customize the parser using FlagSet.NewParser
// and FlagSet.ParseWithParser to support dig-style command-line options.
// This shows how to reuse parsing code for different command-line styles.
//
// Key Features Demonstrated:
// 1. Using FlagSet.NewParser() to get a customizable parser
// 2. Adding custom prefixes (like + for dig-style options)
// 3. Adding custom separators (like ++ or @@)
// 4. Using FlagSet.ParseWithParser() with the customized parser
//
// This approach allows you to:
// - Support multiple command-line styles (Unix, dig, custom) with the same flag definitions
// - Add domain-specific prefixes and separators
// - Maintain consistent flag behavior across different parsing styles
func Example_customParser() {
	fmt.Println("Dig-style parsing with + prefix support:")

	// Create a flag set with standard options
	fs := flagemu.NewFlagSet("dig", flagemu.ContinueOnError)

	// Define standard flags
	verbose := fs.Bool("verbose", 'v', false, "verbose output")
	file := fs.String("file", 'f', "", "config file")

	// Define dig-style plus options (using long names)
	trace := fs.Bool("trace", 0, false, "enable trace mode")
	short := fs.Bool("short", 0, false, "short output format")
	timeout := fs.String("timeout", 0, "5", "query timeout in seconds")

	// Create a custom parser using NewParser
	parser := fs.NewParser()

	// Customize the parser to support dig-style + prefixes
	// Add + as a long option prefix alongside --
	parser.LongOptionPrefixes = append(parser.LongOptionPrefixes, "+")

	// Add ++ as a separator alongside --
	parser.Separators = append(parser.Separators, "++")

	// Example dig-style command line: dig -v +trace --file=config +timeout=10 +short ++ domain.com
	args := []string{"-v", "+trace", "--file=config.txt", "+timeout=10", "+short", "++", "example.com", "extra.com"}

	// Parse using the customized parser
	err := fs.ParseWithParser(parser, args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Display results
	fmt.Printf("Standard flags:\n")
	fmt.Printf("  Verbose: %t\n", *verbose)
	fmt.Printf("  File: %s\n", *file)
	fmt.Printf("Plus options:\n")
	fmt.Printf("  Trace: %t\n", *trace)
	fmt.Printf("  Short: %t\n", *short)
	fmt.Printf("  Timeout: %s\n", *timeout)
	fmt.Printf("Arguments: %v\n", fs.Args())

	fmt.Println("\nMultiple prefix support:")

	// Create another flag set demonstrating multiple prefixes
	fs2 := flagemu.NewFlagSet("multiapp", flagemu.ContinueOnError)
	help := fs2.Bool("help", 'h', false, "show help")
	output := fs2.String("output", 'o', "", "output file")
	enable := fs2.Bool("enable", 0, false, "enable feature")

	// Create and customize parser to support both - and + prefixes
	parser2 := fs2.NewParser()

	// Add + as an additional short option prefix (alongside -)
	parser2.ShortOptionPrefixes = append(parser2.ShortOptionPrefixes, "+")
	parser2.LongOptionPrefixes = append(parser2.LongOptionPrefixes, "+")

	// Mixed prefix command: multiapp -h +enable --output=result.txt +help files.txt
	args2 := []string{"-h", "+enable", "--output=result.txt", "files.txt", "more.txt"}

	err = fs2.ParseWithParser(parser2, args2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Multiple prefix flags:\n")
	fmt.Printf("  Help: %t\n", *help)
	fmt.Printf("  Output: %s\n", *output)
	fmt.Printf("  Enable: %t\n", *enable)
	fmt.Printf("Arguments: %v\n", fs2.Args())

	fmt.Println("\nCustom separators example:")

	// Demonstrate custom separators
	fs3 := flagemu.NewFlagSet("customsep", flagemu.ContinueOnError)
	debug := fs3.Bool("debug", 'd', false, "debug mode")
	config := fs3.String("config", 'c', "", "configuration")

	parser3 := fs3.NewParser()

	// Add custom separators
	parser3.Separators = []string{"--", "@@", "STOP"}

	// Command with custom separator: app -d --config=app.conf @@ these are all arguments
	args3 := []string{"-d", "--config=app.conf", "@@", "these", "are", "all", "arguments"}

	err = fs3.ParseWithParser(parser3, args3)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Custom separator parsing:\n")
	fmt.Printf("  Debug: %t\n", *debug)
	fmt.Printf("  Config: %s\n", *config)
	fmt.Printf("  Arguments: %v\n", fs3.Args())

	// Output:
	// Dig-style parsing with + prefix support:
	// Standard flags:
	//   Verbose: true
	//   File: config.txt
	// Plus options:
	//   Trace: true
	//   Short: true
	//   Timeout: 10
	// Arguments: [example.com extra.com]
	//
	// Multiple prefix support:
	// Multiple prefix flags:
	//   Help: true
	//   Output: result.txt
	//   Enable: true
	// Arguments: [files.txt more.txt]
	//
	// Custom separators example:
	// Custom separator parsing:
	//   Debug: true
	//   Config: app.conf
	//   Arguments: [these are all arguments]
}

// Example_goStyleParsing demonstrates Go-style command-line parsing
// where all options use a single dash prefix.
func Example_goStyleParsing() {
	fmt.Println("Go-style parsing example:")

	// Create a flag set with both short and long option names
	fs := flagemu.NewFlagSet("myapp", flagemu.ContinueOnError)

	// Define flags with both long and short names
	verbose := fs.Bool("verbose", 'v', false, "enable verbose output")
	help := fs.Bool("help", 'h', false, "show help message")
	output := fs.String("output", 'o', "output.txt", "output file path")
	timeout := fs.String("timeout", 't', "30s", "operation timeout")

	// Create a Go-style parser using the dedicated method
	// This ensures both short and long options work with single dash prefix
	parser := fs.NewGoStyleParser()

	// Example Go-style command line: myapp -v -help -output=result.txt -timeout=60s input.txt
	args := []string{"-v", "-help", "-output=result.txt", "-timeout=60s", "input.txt", "extra.txt"}

	err := fs.ParseWithParser(parser, args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Display the parsed results
	fmt.Printf("Parsed flags:\n")
	fmt.Printf("  Verbose: %t\n", *verbose)
	fmt.Printf("  Help: %t\n", *help)
	fmt.Printf("  Output: %s\n", *output)
	fmt.Printf("  Timeout: %s\n", *timeout)
	fmt.Printf("Arguments: %v\n", fs.Args())

	fmt.Println("\nAlternative using short names:")

	// Create another flag set to demonstrate short names work too
	fs2 := flagemu.NewFlagSet("myapp2", flagemu.ContinueOnError)
	debug := fs2.Bool("debug", 'd', false, "enable debug mode")
	file := fs2.String("file", 'f', "", "input file")

	parser2 := fs2.NewGoStyleParser()

	// Using short names: myapp2 -d -f=input.dat data.txt
	args2 := []string{"-d", "-f=input.dat", "data.txt"}

	err = fs2.ParseWithParser(parser2, args2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Short name parsing:\n")
	fmt.Printf("  Debug: %t\n", *debug)
	fmt.Printf("  File: %s\n", *file)
	fmt.Printf("Arguments: %v\n", fs2.Args())

	// Output:
	// Go-style parsing example:
	// Parsed flags:
	//   Verbose: true
	//   Help: true
	//   Output: result.txt
	//   Timeout: 60s
	// Arguments: [input.txt extra.txt]
	//
	// Alternative using short names:
	// Short name parsing:
	//   Debug: true
	//   File: input.dat
	// Arguments: [data.txt]
}

// Example_parserReuse demonstrates reusing parsing logic for different command-line styles,
// similar to how dig command supports multiple option formats.
//
// This example shows the power of the FlagSet.NewParser/ParseWithParser pattern:
// - Define flags once with standard FlagSet methods
// - Create different parser configurations for different command-line styles
// - Parse the same logical command with different syntaxes
// - Maintain consistent behavior and validation across all styles
//
// Benefits of this approach:
// 1. Code reuse: Same flag definitions work with multiple parsing styles
// 2. Consistency: Same validation and behavior regardless of syntax
// 3. Flexibility: Easy to add new command-line styles without changing core logic
// 4. Maintainability: Changes to flag behavior apply to all parsing styles
//
// This is particularly useful for tools that need to support legacy command-line
// formats or multiple user preferences (like dig's support for both -- and + prefixes).
func Example_parserReuse() {
	fmt.Println("Reusable parsing configuration:")

	// Define a common configuration structure
	type Config struct {
		Verbose bool
		Debug   bool
		Output  string
		Timeout string
		Files   []string
	}

	// Function to parse with a given style
	parseWithStyle := func(name string, setupParser func(*flagemu.FlagSet) *parser.Parser, args []string) *Config {
		config := &Config{Timeout: "30s"} // default values
		fs := flagemu.NewFlagSet(name, flagemu.ContinueOnError)

		// Get pointers to the actual flag variables
		verbose := fs.Bool("verbose", 'v', config.Verbose, "verbose output")
		debug := fs.Bool("debug", 'd', config.Debug, "debug mode")
		output := fs.String("output", 'o', config.Output, "output file")
		timeout := fs.String("timeout", 't', config.Timeout, "timeout duration")

		// Get custom parser and parse
		parser := setupParser(fs)
		err := fs.ParseWithParser(parser, args)
		if err != nil {
			fmt.Printf("Error parsing %s style: %v\n", name, err)
			return config
		}

		// Update config with parsed values
		config.Verbose = *verbose
		config.Debug = *debug
		config.Output = *output
		config.Timeout = *timeout
		config.Files = fs.Args()

		return config
	}

	// Unix-style parser setup
	unixStyle := func(fs *flagemu.FlagSet) *parser.Parser {
		parser := fs.NewParser()
		// Unix style is the default, no changes needed
		return parser
	}

	// Dig-style parser setup
	digStyle := func(fs *flagemu.FlagSet) *parser.Parser {
		parser := fs.NewParser()
		parser.LongOptionPrefixes = append(parser.LongOptionPrefixes, "+")
		parser.Separators = append(parser.Separators, "++")
		return parser
	}

	// Multiple separator parser setup
	multiSepStyle := func(fs *flagemu.FlagSet) *parser.Parser {
		parser := fs.NewParser()
		// Add multiple separators
		parser.Separators = []string{"--", "@@", "STOP"}
		return parser
	}

	// Parse the same logical command with different styles

	// Unix style: app -v --debug --output=result.txt --timeout=60s file1.txt file2.txt
	unixArgs := []string{"-v", "--debug", "--output=result.txt", "--timeout=60s", "file1.txt", "file2.txt"}
	unixConfig := parseWithStyle("unix-app", unixStyle, unixArgs)

	fmt.Printf("Unix style result:\n")
	fmt.Printf("  Verbose: %t, Debug: %t, Output: %s, Timeout: %s\n",
		unixConfig.Verbose, unixConfig.Debug, unixConfig.Output, unixConfig.Timeout)
	fmt.Printf("  Files: %v\n", unixConfig.Files)

	// Dig style: dig -v +debug --output=result.txt +timeout=60s file1.txt file2.txt
	digArgs := []string{"-v", "+debug", "--output=result.txt", "+timeout=60s", "file1.txt", "file2.txt"}
	digConfig := parseWithStyle("dig-app", digStyle, digArgs)

	fmt.Printf("Dig style result:\n")
	fmt.Printf("  Verbose: %t, Debug: %t, Output: %s, Timeout: %s\n",
		digConfig.Verbose, digConfig.Debug, digConfig.Output, digConfig.Timeout)
	fmt.Printf("  Files: %v\n", digConfig.Files)

	// Multi-separator style: app -v --debug STOP these are all arguments
	multiSepArgs := []string{"-v", "--debug", "--output=result.txt", "STOP", "these", "are", "arguments"}
	multiSepConfig := parseWithStyle("multi-app", multiSepStyle, multiSepArgs)

	fmt.Printf("Multi-separator style result:\n")
	fmt.Printf("  Verbose: %t, Debug: %t, Output: %s, Timeout: %s\n",
		multiSepConfig.Verbose, multiSepConfig.Debug, multiSepConfig.Output, multiSepConfig.Timeout)
	fmt.Printf("  Files: %v\n", multiSepConfig.Files)

	// Output:
	// Reusable parsing configuration:
	// Unix style result:
	//   Verbose: true, Debug: true, Output: result.txt, Timeout: 60s
	//   Files: [file1.txt file2.txt]
	// Dig style result:
	//   Verbose: true, Debug: true, Output: result.txt, Timeout: 60s
	//   Files: [file1.txt file2.txt]
	// Multi-separator style result:
	//   Verbose: true, Debug: true, Output: result.txt, Timeout: 30s
	//   Files: [these are arguments]
}

// Analysis of the FlagSet.NewParser/ParseWithParser Solution:
//
// STRENGTHS:
// 1. **Code Reuse**: The same flag definitions can be used with multiple parsing styles.
//    This is exactly what dig-style tools need - the ability to support both traditional
//    Unix flags (--verbose) and dig-specific syntax (+trace).
//
// 2. **Separation of Concerns**: Flag definition is separate from parsing style.
//    You define what flags exist once, then customize how they're parsed.
//
// 3. **Flexibility**: Easy to add new command-line styles without changing existing code.
//    New styles just require a new parser configuration function.
//
// 4. **Consistency**: All parsing styles share the same validation, type conversion,
//    and error handling logic. This ensures consistent behavior regardless of syntax.
//
// 5. **Backward Compatibility**: Can support legacy command-line formats while
//    introducing new ones, which is crucial for tools with established user bases.
//
// PRACTICAL APPLICATIONS:
// - dig-style DNS tools that support both --option and +option syntax
// - Cross-platform tools that need to support both Unix (-) and Windows (/) prefixes
// - Applications migrating between different command-line conventions
// - Tools that need custom separators for domain-specific syntax
//
// COMPARISON TO ALTERNATIVES:
// - Better than multiple flag parsers: Avoids code duplication and inconsistencies
// - Better than conditional parsing: Cleaner separation and easier to extend
// - Better than string manipulation: Leverages existing parser infrastructure
//
// RECOMMENDATION:
// This solution is excellent for tools that need to support multiple command-line
// styles while maintaining code clarity and consistency. It's particularly valuable
// for dig-style applications where users expect both traditional and domain-specific
// option syntax to work seamlessly.
