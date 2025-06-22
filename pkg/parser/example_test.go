// example_test.go - Parse example tests.
// SPDX-License-Identifier: GPL-3.0-or-later

package parser_test

import (
	"fmt"

	"github.com/bassosimone/clip/pkg/parser"
)

// ExampleParser_dig demonstrates dig command-line parsing style.
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
func ExampleParser_dig() {
	// Configure parser for dig style
	p := &parser.Parser{
		ShortOptions: map[string]parser.OptionType{
			"v": parser.OptionTypeBool,   // -v for verbose
			"f": parser.OptionTypeString, // -f for file
			"h": parser.OptionTypeBool,   // -h for help
		},
		LongOptions: map[string]parser.OptionType{
			"verbose": parser.OptionTypeBool,   // --verbose
			"file":    parser.OptionTypeString, // --file
			"port":    parser.OptionTypeString, // --port
			"trace":   parser.OptionTypeBool,   // +trace (using long options map for + prefix)
			"short":   parser.OptionTypeBool,   // +short
			"timeout": parser.OptionTypeString, // +timeout
		},
		ShortOptionPrefixes: []string{"-"},
		LongOptionPrefixes:  []string{"--", "+"}, // Both -- and + prefixes for long options
		Separators:          []string{"--"},
	}

	// Example dig command line: program -v IN A +trace --port=53 +timeout=5 -f config -- remaining args
	args := []string{"program", "-v", "IN", "A", "+trace", "--port=53", "+timeout=5", "-f", "config", "--", "remaining", "args"}

	items, err := p.Parse(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, item := range items {
		fmt.Printf("%#v\n", item)
	}

	// Output:
	// parser.ProgramNameItem{Name:"program", Token:scanner.ProgramNameToken{Index:0, Name:"program"}}
	// parser.OptionItem{Name:"v", Token:scanner.OptionToken{Index:1, Prefix:"-", Name:"v"}, Value:"true", IsShort:true, Type:2, Prefix:"-"}
	// parser.ArgumentItem{Token:scanner.ArgumentToken{Index:2, Value:"IN"}, Value:"IN"}
	// parser.ArgumentItem{Token:scanner.ArgumentToken{Index:3, Value:"A"}, Value:"A"}
	// parser.OptionItem{Name:"trace", Token:scanner.OptionToken{Index:4, Prefix:"+", Name:"trace"}, Value:"true", IsShort:false, Type:2, Prefix:"+"}
	// parser.OptionItem{Name:"port", Token:scanner.OptionToken{Index:5, Prefix:"--", Name:"port=53"}, Value:"53", IsShort:false, Type:1, Prefix:"--"}
	// parser.OptionItem{Name:"timeout", Token:scanner.OptionToken{Index:6, Prefix:"+", Name:"timeout=5"}, Value:"5", IsShort:false, Type:1, Prefix:"+"}
	// parser.OptionItem{Name:"f", Token:scanner.OptionToken{Index:7, Prefix:"-", Name:"f"}, Value:"config", IsShort:true, Type:1, Prefix:"-"}
	// parser.SeparatorItem{Token:scanner.SeparatorToken{Index:9, Separator:"--"}, Separator:"--"}
	// parser.ArgumentItem{Token:scanner.ArgumentToken{Index:10, Value:"remaining"}, Value:"remaining"}
	// parser.ArgumentItem{Token:scanner.ArgumentToken{Index:11, Value:"args"}, Value:"args"}
}

// ExampleParser_gnu demonstrates GNU command-line parsing.
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
//   - Argument permutation (reordering) typically supported at parser levelfunc ExampleParser_gnu() {
func ExampleParser_gnu() {
	// Configure parser for GNU style
	p := &parser.Parser{
		ShortOptions: map[string]parser.OptionType{
			"v": parser.OptionTypeBool,   // -v for verbose
			"f": parser.OptionTypeString, // -f for file
			"h": parser.OptionTypeBool,   // -h for help
		},
		LongOptions: map[string]parser.OptionType{
			"verbose": parser.OptionTypeBool,   // --verbose
			"file":    parser.OptionTypeString, // --file
			"help":    parser.OptionTypeBool,   // --help
		},
		ShortOptionPrefixes: []string{"-"},
		LongOptionPrefixes:  []string{"--"},
		Separators:          []string{"--"},
	}

	// Example GNU command line: program -vh target --file=config.txt input.txt --verbose -- --not-an-option
	args := []string{"program", "-vh", "target", "--file=config.txt", "input.txt", "--verbose", "--", "--not-an-option", "output.txt"}

	items, err := p.Parse(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, item := range items {
		fmt.Printf("%#v\n", item)
	}

	// Output:
	// parser.ProgramNameItem{Name:"program", Token:scanner.ProgramNameToken{Index:0, Name:"program"}}
	// parser.OptionItem{Name:"v", Token:scanner.OptionToken{Index:1, Prefix:"-", Name:"vh"}, Value:"true", IsShort:true, Type:2, Prefix:"-"}
	// parser.OptionItem{Name:"h", Token:scanner.OptionToken{Index:1, Prefix:"-", Name:"vh"}, Value:"true", IsShort:true, Type:2, Prefix:"-"}
	// parser.ArgumentItem{Token:scanner.ArgumentToken{Index:2, Value:"target"}, Value:"target"}
	// parser.OptionItem{Name:"file", Token:scanner.OptionToken{Index:3, Prefix:"--", Name:"file=config.txt"}, Value:"config.txt", IsShort:false, Type:1, Prefix:"--"}
	// parser.ArgumentItem{Token:scanner.ArgumentToken{Index:4, Value:"input.txt"}, Value:"input.txt"}
	// parser.OptionItem{Name:"verbose", Token:scanner.OptionToken{Index:5, Prefix:"--", Name:"verbose"}, Value:"true", IsShort:false, Type:2, Prefix:"--"}
	// parser.SeparatorItem{Token:scanner.SeparatorToken{Index:6, Separator:"--"}, Separator:"--"}
	// parser.ArgumentItem{Token:scanner.OptionToken{Index:7, Prefix:"--", Name:"not-an-option"}, Value:"--not-an-option"}
	// parser.ArgumentItem{Token:scanner.ArgumentToken{Index:8, Value:"output.txt"}, Value:"output.txt"}
}

// ExampleParser_go demonstrates Go command-line parsing style.
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
func ExampleParser_go() {
	// Configure parser for Go style
	p := &parser.Parser{
		ShortOptions: map[string]parser.OptionType{}, // Go doesn't distinguish between short/long
		LongOptions: map[string]parser.OptionType{
			"v":       parser.OptionTypeBool,   // -v for verbose (single character with single dash)
			"file":    parser.OptionTypeString, // -file for file (multi-character with single dash)
			"verbose": parser.OptionTypeBool,   // -verbose (multi-character with single dash)
			"help":    parser.OptionTypeBool,   // -help (multi-character with single dash)
		},
		ShortOptionPrefixes: []string{},
		LongOptionPrefixes:  []string{"-"},  // Go uses single dash for all options
		Separators:          []string{"--"}, // Recognize the GNU like separator
	}

	// Example Go command line: program -v -file=config.txt -verbose input.txt -- -extra
	args := []string{"program", "-v", "-file=config.txt", "-verbose", "input.txt", "--", "-extra"}

	items, err := p.Parse(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, item := range items {
		fmt.Printf("%#v\n", item)
	}

	// Output:
	// parser.ProgramNameItem{Name:"program", Token:scanner.ProgramNameToken{Index:0, Name:"program"}}
	// parser.OptionItem{Name:"v", Token:scanner.OptionToken{Index:1, Prefix:"-", Name:"v"}, Value:"true", IsShort:false, Type:2, Prefix:"-"}
	// parser.OptionItem{Name:"file", Token:scanner.OptionToken{Index:2, Prefix:"-", Name:"file=config.txt"}, Value:"config.txt", IsShort:false, Type:1, Prefix:"-"}
	// parser.OptionItem{Name:"verbose", Token:scanner.OptionToken{Index:3, Prefix:"-", Name:"verbose"}, Value:"true", IsShort:false, Type:2, Prefix:"-"}
	// parser.ArgumentItem{Token:scanner.ArgumentToken{Index:4, Value:"input.txt"}, Value:"input.txt"}
	// parser.SeparatorItem{Token:scanner.SeparatorToken{Index:5, Separator:"--"}, Separator:"--"}
	// parser.ArgumentItem{Token:scanner.OptionToken{Index:6, Prefix:"-", Name:"extra"}, Value:"-extra"}
}

// ExampleParser_noPermute demonstrates parsing with FlagNoPermute.
//
// This is like the POSIXLY_CORRECT behavior in GNU getopt.
func ExampleParser_noPermute() {
	// Configure parser with FlagNoPermute
	p := &parser.Parser{
		Flags: parser.FlagNoPermute, // Stop parsing options at first argument
		ShortOptions: map[string]parser.OptionType{
			"v": parser.OptionTypeBool,   // -v for verbose
			"f": parser.OptionTypeString, // -f for file
			"h": parser.OptionTypeBool,   // -h for help
		},
		LongOptions: map[string]parser.OptionType{
			"verbose": parser.OptionTypeBool,   // --verbose
			"file":    parser.OptionTypeString, // --file
			"help":    parser.OptionTypeBool,   // --help
		},
		ShortOptionPrefixes: []string{"-"},
		LongOptionPrefixes:  []string{"--"},
		Separators:          []string{"--"},
	}

	// Example with no-permute: program -v input.txt --verbose -f config.txt
	// After "input.txt", everything is treated as arguments, including --verbose and -f
	args := []string{"program", "-v", "input.txt", "--verbose", "-f", "config.txt", "output.txt"}

	items, err := p.Parse(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, item := range items {
		fmt.Printf("%#v\n", item)
	}

	// Output:
	// parser.ProgramNameItem{Name:"program", Token:scanner.ProgramNameToken{Index:0, Name:"program"}}
	// parser.OptionItem{Name:"v", Token:scanner.OptionToken{Index:1, Prefix:"-", Name:"v"}, Value:"true", IsShort:true, Type:2, Prefix:"-"}
	// parser.ArgumentItem{Token:scanner.ArgumentToken{Index:2, Value:"input.txt"}, Value:"input.txt"}
	// parser.ArgumentItem{Token:scanner.OptionToken{Index:3, Prefix:"--", Name:"verbose"}, Value:"--verbose"}
	// parser.ArgumentItem{Token:scanner.OptionToken{Index:4, Prefix:"-", Name:"f"}, Value:"-f"}
	// parser.ArgumentItem{Token:scanner.ArgumentToken{Index:5, Value:"config.txt"}, Value:"config.txt"}
	// parser.ArgumentItem{Token:scanner.ArgumentToken{Index:6, Value:"output.txt"}, Value:"output.txt"}
}

// ExampleParser_unix demonstrates traditional UNIX command-line parsing.
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
func ExampleParser_unix() {
	// Configure parser for traditional UNIX style
	p := &parser.Parser{
		ShortOptions: map[string]parser.OptionType{
			"v": parser.OptionTypeBool,   // -v for verbose
			"f": parser.OptionTypeString, // -f for file
			"h": parser.OptionTypeBool,   // -h for help
			"d": parser.OptionTypeBool,   // -d for debug
		},
		LongOptions:         map[string]parser.OptionType{}, // No long options in traditional UNIX
		ShortOptionPrefixes: []string{"-"},
		LongOptionPrefixes:  []string{}, // No long option prefixes
		Separators:          []string{}, // No separators in traditional UNIX
	}

	// Example traditional UNIX command line: program -vh -f config.txt -d input.txt
	args := []string{"program", "-vh", "-f", "config.txt", "-d", "input.txt"}

	items, err := p.Parse(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, item := range items {
		fmt.Printf("%#v\n", item)
	}

	// Output:
	// parser.ProgramNameItem{Name:"program", Token:scanner.ProgramNameToken{Index:0, Name:"program"}}
	// parser.OptionItem{Name:"v", Token:scanner.OptionToken{Index:1, Prefix:"-", Name:"vh"}, Value:"true", IsShort:true, Type:2, Prefix:"-"}
	// parser.OptionItem{Name:"h", Token:scanner.OptionToken{Index:1, Prefix:"-", Name:"vh"}, Value:"true", IsShort:true, Type:2, Prefix:"-"}
	// parser.OptionItem{Name:"f", Token:scanner.OptionToken{Index:2, Prefix:"-", Name:"f"}, Value:"config.txt", IsShort:true, Type:1, Prefix:"-"}
	// parser.OptionItem{Name:"d", Token:scanner.OptionToken{Index:4, Prefix:"-", Name:"d"}, Value:"true", IsShort:true, Type:2, Prefix:"-"}
	// parser.ArgumentItem{Token:scanner.ArgumentToken{Index:5, Value:"input.txt"}, Value:"input.txt"}
}
