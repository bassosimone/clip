// example_test.go - getopt examples.
// SPDX-License-Identifier: GPL-3.0-or-later

package getopt_test

import (
	"fmt"

	"github.com/bassosimone/clip/pkg/getopt"
)

// ExampleShort demonstrates how to use Short to parse traditional UNIX getopt
// style command line arguments. This example shows:
//
//   - Traditional short options (-v, -f, etc.)
//
//   - Short options with arguments (-f filename)
//
//   - Bundled options (-vhq)
//
//   - Argument after bundled options (-vhqfile.txt)
//
//   - Options after separator are treated as arguments
func ExampleShort() {
	// Example command line for a hypothetical file archiver
	argv := []string{
		"program",
		"-vf",         // verbose mode and specify archive
		"archive.tar", // archive name
		"-C",          // change directory
		"/tmp/dir",    // directory to change to
		"-z",          // enable compression
		"file1.txt",   // files to archive
		"file2.txt",
	}

	// Parse using traditional getopt style with options that take arguments
	items, err := getopt.Short(argv, "vf:C:z")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Convert back to argv format for easy comparison
	result := getopt.Serialize(items)
	for _, arg := range result {
		fmt.Printf("%s\n", arg)
	}

	// Output:
	// program
	// -v
	// -f
	// archive.tar
	// -C
	// /tmp/dir
	// -z
	// file1.txt
	// file2.txt
}

// ExampleLong demonstrates how to use Long to parse GNU getopt_long style
// command line arguments. This example shows:
//
//   - Traditional short options (-v, -f)
//
//   - Long options (--verbose, --file)
//
//   - Options with required arguments (--file=name, -f name)
//
//   - Mixed short and long options
//
//   - Options with and without arguments
//
//   - Separator handling
func ExampleLong() {
	// Example command line for a hypothetical network diagnostic tool
	argv := []string{
		"program",
		"--verbose",    // verbose output
		"-p",           // port number
		"8080",         // port value
		"--timeout=30", // timeout with value
		"-r3",          // retry count
		"--",           // separator
		"ping",         // command to run
		"--host",       // host to connect to
		"example.com",  // hostname
	}

	// Define long options
	options := []getopt.Option{
		{Name: "verbose", HasArg: false},
		{Name: "timeout", HasArg: true},
	}

	// Parse using GNU getopt_long style
	items, err := getopt.Long(argv, "p:r:", options)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Convert back to argv format for easy comparison
	result := getopt.Serialize(items)
	for _, arg := range result {
		fmt.Printf("%s\n", arg)
	}

	// Output:
	// program
	// --verbose
	// -p
	// 8080
	// --timeout
	// 30
	// -r
	// 3
	// --
	// ping
	// --host
	// example.com
}

// ExampleMain demonstrates how to use Main to implement the getopt(1)
// tool. This example shows:
//
//   - Defining short options (-o optstring)
//
//   - Defining long options (--longoptions)
//
//   - Multiple --longoptions flags
//
//   - Complex option string parsing
//
//   - Separator handling
//
//   - Reordering of arguments
func ExampleMain() {
	// Example showing how to process a complex command line with getopt
	argv := []string{
		"getopt",
		"-o", "vf:C:h",
		"--longoptions", "verbose,file:,chdir:,help",
		"--longoptions", "debug,output:",
		"--",
		"-vf",
		"test.txt",
		"-C",
		"/tmp",
		"--debug",
		"extra",
		"args",
	}

	// Process using getopt main function
	result, err := getopt.Main(argv)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print the reordered arguments
	for _, arg := range result {
		fmt.Printf("%s\n", arg)
	}

	// Output:
	// -v
	// -f
	// test.txt
	// -C
	// /tmp
	// --debug
	// extra
	// args
}
