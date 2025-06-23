// example_test.go - Flag package example tests.
// SPDX-License-Identifier: GPL-3.0-or-later

package flag_test

import (
	"fmt"

	"github.com/bassosimone/clip/pkg/flag"
)

// ExampleFlagSet demonstrates how to use various flag types and parse command line
// arguments. This example showcases both boolean and string flags with their
// short and long forms.
func ExampleFlagSet() {
	// Create FlagSet
	fset := flag.NewFlagSet("program", flag.ContinueOnError)

	// Define flags with both short and long options
	help := fset.Bool("help", 'h', false, "display help information")
	verbose := fset.Bool("verbose", 'v', false, "enable verbose output")
	output := fset.String("output", 'o', "output.txt", "output file path")
	format := fset.String("format", 'f', "text", "output format")

	// Parse command line arguments (simulated here)
	args := []string{"-v", "--output=result.txt", "-f", "json", "input1.txt", "input2.txt"}
	if err := fset.Parse(args); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		return
	}

	// Print the values of the flags and arguments
	fmt.Printf("help = %v\n", *help)
	fmt.Printf("verbose = %v\n", *verbose)
	fmt.Printf("output = %v\n", *output)
	fmt.Printf("format = %v\n", *format)
	fmt.Printf("args = %v\n", fset.Args())

	// Output:
	// help = false
	// verbose = true
	// output = result.txt
	// format = json
	// args = [input1.txt input2.txt]
}

// ExampleFlagSet_usage demonstrates how the usage information is formatted
// and displayed for a set of flags.
func ExampleFlagSet_usage() {
	// Create FlagSet for showing usage
	fset := flag.NewFlagSet("myprogram", flag.ContinueOnError)
	fset.SetDescription("Illustrates how to print the usage information.")

	// Define various flags to demonstrate usage formatting
	fset.Bool("help", 'h', false, "display this help message")
	fset.Bool("verbose", 'v', false, "enable verbose output mode")
	fset.String("output", 'o', "output.txt", "path to the output file")
	fset.String("format", 'f', "text", "output format (text or json)")
	fset.Bool("debug", 'd', false, "enable debug logging")
	fset.String("config", 'c', "", "path to configuration file")

	// Get the usage string and print it
	usage := fset.Usage()
	fmt.Print(usage)

	// Output:
	// Usage: myprogram [options] [--] [arguments]
	//
	// Illustrates how to print the usage information.
	//
	// Options:
	//   -c, --config VALUE
	//     path to configuration file
	//
	//   -d, --debug
	//     enable debug logging
	//
	//   -f, --format VALUE
	//     output format (text or json)
	//
	//   -h, --help
	//     display this help message
	//
	//   -o, --output VALUE
	//     path to the output file
	//
	//   -v, --verbose
	//     enable verbose output mode
	//
	// Use 'myprogram --help' to show this help screen.
}

// ExampleFlagSet_allMethods demonstrates all the different ways to define flags
// using various methods like Type, TypeVar, TypeLong, TypeShort, etc.
func ExampleFlagSet_allMethods() {
	// Create FlagSet
	fset := flag.NewFlagSet("allmethods", flag.ContinueOnError)

	// 1. Regular methods (return pointer)
	// Bool variations
	debug := fset.Bool("debug", 'd', false, "enable debug mode") // both long and short
	trace := fset.BoolLong("trace", false, "enable trace mode")  // long only
	quiet := fset.BoolShort('q', "suppress output")              // short only

	// String variations
	config := fset.String("config", 'c', "config.json", "config file path") // both long and short
	logfile := fset.StringLong("logfile", "app.log", "log file path")       // long only
	mode := fset.StringShort('m', "operation mode")                         // short only

	// 2. Var methods (use existing pointer)
	var (
		verbose bool   = false
		force   bool   = false
		dryrun  bool   = false
		output  string = "output.txt"
		format  string = "text"
		level   string = "info"
	)

	// BoolVar variations
	fset.BoolVar(&verbose, "verbose", 'v', "verbose output") // both long and short
	fset.BoolLongVar(&force, "force", "force operation")     // long only
	fset.BoolShortVar(&dryrun, 'n', "dry run mode")          // short only

	// StringVar variations
	fset.StringVar(&output, "output", 'o', "output file")  // both long and short
	fset.StringLongVar(&format, "format", "output format") // long only
	fset.StringShortVar(&level, 'l', "log level")          // short only

	// Parse example arguments using various flag styles
	args := []string{
		"-d",                // short bool
		"--trace",           // long bool
		"-q",                // short-only bool
		"-v",                // var short bool
		"--force",           // var long bool
		"-n",                // var short-only bool
		"--config=test.cfg", // long with value
		"-c", "other.cfg",   // short with value
		"--logfile=app.log", // long-only with value
		"-m", "fast",        // short-only with value
		"-o", "result.txt", // var short with value
		"--format=json", // var long with value
		"-l", "debug",   // var short-only with value
	}

	if err := fset.Parse(args); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		return
	}

	// Print values to demonstrate all flags were properly set
	fmt.Printf("Regular methods (return pointer):\n")
	fmt.Printf("  debug = %v\n", *debug)
	fmt.Printf("  trace = %v\n", *trace)
	fmt.Printf("  quiet = %v\n", *quiet)
	fmt.Printf("  config = %v\n", *config)
	fmt.Printf("  logfile = %v\n", *logfile)
	fmt.Printf("  mode = %v\n", *mode)

	fmt.Printf("\nVar methods (existing pointer):\n")
	fmt.Printf("  verbose = %v\n", verbose)
	fmt.Printf("  force = %v\n", force)
	fmt.Printf("  dryrun = %v\n", dryrun)
	fmt.Printf("  output = %v\n", output)
	fmt.Printf("  format = %v\n", format)
	fmt.Printf("  level = %v\n", level)

	// Output:
	// Regular methods (return pointer):
	//   debug = true
	//   trace = true
	//   quiet = true
	//   config = other.cfg
	//   logfile = app.log
	//   mode = fast
	//
	// Var methods (existing pointer):
	//   verbose = true
	//   force = true
	//   dryrun = true
	//   output = result.txt
	//   format = json
	//   level = debug
}

// ExampleFlagSet_customParser demonstrates how to customize the parser's behavior
// by modifying option prefixes and separators.
func ExampleFlagSet_customParser() {
	// Create FlagSet
	fset := flag.NewFlagSet("customprogram", flag.ContinueOnError)

	// Get the parser and customize it
	parser := fset.Parser()
	parser.LongOptionPrefixes = []string{"--", "+"} // Allow both -- and + for long options
	parser.ShortOptionPrefixes = []string{"-", "/"} // Allow both - and / for short options
	parser.Separators = []string{"--", "++"}        // Allow both -- and ++ as separators

	// Define some flags
	verbose := fset.Bool("verbose", 'v', false, "enable verbose mode")
	debug := fset.Bool("debug", 'd', false, "enable debug mode")
	output := fset.String("output", 'o', "out.txt", "output file")

	// Parse with custom prefixes
	args := []string{"/v", "+debug", "--output=test.txt", "++", "--notaflag"}
	if err := fset.Parse(args); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		return
	}

	// Print results
	fmt.Printf("verbose = %v\n", *verbose)
	fmt.Printf("debug = %v\n", *debug)
	fmt.Printf("output = %v\n", *output)
	fmt.Printf("args = %v\n", fset.Args())

	// Output:
	// verbose = true
	// debug = true
	// output = test.txt
	// args = [--notaflag]
}
