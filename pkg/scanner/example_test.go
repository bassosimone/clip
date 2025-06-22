// example_test.go - Scanner example tests
// SPDX-License-Identifier: GPL-3.0-or-later

package scanner_test

import (
	"fmt"

	"github.com/bassosimone/clip/pkg/scanner"
)

// ExampleScanner_dig demonstrates dig command-line parsing style.
//
// Dig style:
//
//   - Traditional short options with single dash: -v, -f
//
//   - Long options with double dash: --verbose, --file
//
//   - Plus options for dig-specific features: +trace, +short, +noall
//
//   - Plus options are treated as long options (no bundling)
//
//   - Options with arguments: -f file, +timeout=5, --port=53
//
//   - Double dash separator: -- stops option parsing
//
//   - Mixed prefix styles for different option categories
func ExampleScanner_dig() {
	// Configure scanner for dig style
	s := &scanner.Scanner{
		Prefixes:   []string{"-", "--", "+"}, // Single dash, double dash, and plus prefixes
		Separators: []string{"--"},           // Only double dash separator supported
	}

	args := []string{"program", "-v", "+trace", "--verbose", "+short=yes", "-f", "config", "--", "remaining", "-args"}

	tokens, err := s.Scan(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, token := range tokens {
		fmt.Printf("%#v\n", token)
	}

	// Output:
	// scanner.ProgramNameToken{Index:0, Name:"program"}
	// scanner.OptionToken{Index:1, Prefix:"-", Name:"v"}
	// scanner.OptionToken{Index:2, Prefix:"+", Name:"trace"}
	// scanner.OptionToken{Index:3, Prefix:"--", Name:"verbose"}
	// scanner.OptionToken{Index:4, Prefix:"+", Name:"short=yes"}
	// scanner.OptionToken{Index:5, Prefix:"-", Name:"f"}
	// scanner.ArgumentToken{Index:6, Value:"config"}
	// scanner.SeparatorToken{Index:7, Separator:"--"}
	// scanner.ArgumentToken{Index:8, Value:"remaining"}
	// scanner.OptionToken{Index:9, Prefix:"-", Name:"args"}
}

// ExampleScanner_gnu demonstrates GNU command-line parsing.
//
// GNU style:
//
//   - Short options with single dash: -v, -f
//
//   - Long options with double dash: --verbose, --file
//
//   - Short options can be bundled: -vf equivalent to -v -f
//
//   - Options with arguments: -f file, -ffile, --file=name, --file name
//
//   - Double dash separator: -- stops option parsing
//
//   - Argument permutation (reordering) typically supported at parser level
func ExampleScanner_gnu() {
	// Configure scanner for GNU style
	s := &scanner.Scanner{
		Prefixes:   []string{"-", "--"}, // Single and double dash prefixes
		Separators: []string{"--"},
	}

	// Example command line: program -v --file=config.txt -abc -- --an-option input.txt
	args := []string{"program", "-v", "--file=config.txt", "-abc", "--", "--an-option", "input.txt"}

	tokens, err := s.Scan(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, token := range tokens {
		fmt.Printf("%#v\n", token)
	}

	// Output:
	// scanner.ProgramNameToken{Index:0, Name:"program"}
	// scanner.OptionToken{Index:1, Prefix:"-", Name:"v"}
	// scanner.OptionToken{Index:2, Prefix:"--", Name:"file=config.txt"}
	// scanner.OptionToken{Index:3, Prefix:"-", Name:"abc"}
	// scanner.SeparatorToken{Index:4, Separator:"--"}
	// scanner.OptionToken{Index:5, Prefix:"--", Name:"an-option"}
	// scanner.ArgumentToken{Index:6, Value:"input.txt"}
}

// ExampleScanner_go demonstrates Go command-line parsing style.
//
// Go style:
//
//   - Short and long options with single dash: -v, -f, -verbose, -file
//
//   - No short option bundling: -vf is treated as single option "vf"
//
//   - Options with arguments: -file=name, -file name
//
//   - Double dash separator: -- stops option parsing
//
//   - Simple and consistent: all options use single dash prefix
func ExampleScanner_go() {
	// Configure scanner for Go style
	s := &scanner.Scanner{
		Prefixes:   []string{"-"}, // Go uses single dash for all options
		Separators: []string{"--"},
	}

	// Example command line: program -v -file=config.txt -verbose -debug input.txt -- extra
	args := []string{"program", "-v", "-file=config.txt", "-verbose", "-debug", "input.txt", "--", "extra"}

	tokens, err := s.Scan(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, token := range tokens {
		fmt.Printf("%#v\n", token)
	}

	// Output:
	// scanner.ProgramNameToken{Index:0, Name:"program"}
	// scanner.OptionToken{Index:1, Prefix:"-", Name:"v"}
	// scanner.OptionToken{Index:2, Prefix:"-", Name:"file=config.txt"}
	// scanner.OptionToken{Index:3, Prefix:"-", Name:"verbose"}
	// scanner.OptionToken{Index:4, Prefix:"-", Name:"debug"}
	// scanner.ArgumentToken{Index:5, Value:"input.txt"}
	// scanner.SeparatorToken{Index:6, Separator:"--"}
	// scanner.ArgumentToken{Index:7, Value:"extra"}
}

// ExampleScanner_unix demonstrates traditional UNIX command-line parsing.
//
// Traditional UNIX style:
//
//   - Only single-dash short options: -v, -f
//
//   - No long options (--verbose not supported)
//
//   - Short options can be bundled: -vf equivalent to -v -f
//
//   - Options with arguments: -f file or -ffile
//
//   - No special separators (-- not recognized)
func ExampleScanner_unix() {
	// Configure scanner for traditional UNIX style
	s := &scanner.Scanner{
		Prefixes:   []string{"-"}, // Only single-dash options in traditional UNIX
		Separators: []string{},    // No separators in traditional UNIX
	}

	// Example command line: program -v -f file.txt -abc input.txt
	args := []string{"program", "-v", "-f", "file.txt", "-abc", "input.txt"}

	tokens, err := s.Scan(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, token := range tokens {
		fmt.Printf("%#v\n", token)
	}

	// Output:
	// scanner.ProgramNameToken{Index:0, Name:"program"}
	// scanner.OptionToken{Index:1, Prefix:"-", Name:"v"}
	// scanner.OptionToken{Index:2, Prefix:"-", Name:"f"}
	// scanner.ArgumentToken{Index:3, Value:"file.txt"}
	// scanner.OptionToken{Index:4, Prefix:"-", Name:"abc"}
	// scanner.ArgumentToken{Index:5, Value:"input.txt"}
}
