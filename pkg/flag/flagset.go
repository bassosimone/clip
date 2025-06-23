// flagset.go - FlagSet implementation for command-line flag parsing
// SPDX-License-Identifier: GPL-3.0-or-later

package flag

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/parser"
)

// ErrorHandling defines the behavior of the flag set when an error occurs.
type ErrorHandling int

// These flags control how we handle flag parsing errors.
const (
	// ContinueOnError causes the flag set to continue parsing after an error occurs.
	ContinueOnError = ErrorHandling(iota)

	// ExitOnError causes the flag set to call Exit(2) after an error occurs.
	ExitOnError

	// PanicOnError causes the flag set to panic after an error occurs.
	PanicOnError
)

// FlagSet represents a set of command-line flags.
//
// This struct contains a [*parser.Parser] configured
// as documented by [*parser.Parser].
//
// The zero value is invalid. Construct using [NewFlagSet].
type FlagSet struct {
	// args contains the positional arguments.
	args []string

	// argsdocs contains the positional arguments documentation.
	argsdocs string

	// customusage is an optional custom usage function.
	customusage func(fx *FlagSet) string

	// description contains the optional description.
	description string

	// handling controls how we handle flag parsing errors.
	handling ErrorHandling

	// optionsShort maps short options to their values.
	optionsShort map[string]*Option

	// optionsLong maps long options to their values.
	optionsLong map[string]*Option

	// parser contains the flag parser.
	parser *parser.Parser

	// progname contains the name of the program.
	progname string

	// stderr is the output stream for error messages.
	stderr io.Writer

	// stdout is the output stream for usage messages.
	stdout io.Writer
}

func newDefaultParser() *parser.Parser {
	return &parser.Parser{
		Flags:               0, // permute by default
		LongOptions:         map[string]parser.OptionType{},
		LongOptionPrefixes:  []string{"--"},
		Separators:          []string{"--"},
		ShortOptions:        map[string]parser.OptionType{},
		ShortOptionPrefixes: []string{"-"},
	}
}

// NewFlagSet creates a new [*FlagSet] with the given program name and error handling.
func NewFlagSet(progname string, handling ErrorHandling) *FlagSet {
	return &FlagSet{
		args:         []string{},
		argsdocs:     "[arguments]",
		handling:     handling,
		optionsShort: make(map[string]*Option),
		optionsLong:  make(map[string]*Option),
		parser:       newDefaultParser(),
		progname:     progname,
		stderr:       os.Stderr,
		stdout:       os.Stdout,
	}
}

// Args returns the parsed positional arguments.
//
// This method MUST be invoked after [*FlagSet.Parse].
func (fx *FlagSet) Args() []string {
	return fx.args
}

// Parser returns the [*parser.Parser] used by the [*FlagSet].
//
// You usually want to use this method to customize the flag prefixes
// and the explicit options-arguments separator.
//
// The default [*parser.Parser] configuration is this:
//
//   - Long option prefixes: "--"
//
//   - Short option prefixes: "-"
//
//   - Options-arguments eparators: "--"
//
// You MUST use this method before calling [*FlagSet.Parse].
func (fx *FlagSet) Parser() *parser.Parser {
	return fx.parser
}

// AddOption adds an [*Option] to the [*FlagSet]. Usually, you use more
// specialized methods to add options. This method is useful when you
// need to add an option [Value] not supported by this package.
//
// You MUST use this method before calling [*FlagSet.Parse].
func (fx *FlagSet) AddOption(opt *Option) {
	if opt.LongName != "" {
		fx.parser.LongOptions[opt.LongName] = opt.Value.OptionType()
		fx.optionsLong[opt.LongName] = opt
	}
	if opt.ShortName != 0 {
		fx.parser.ShortOptions[string(opt.ShortName)] = opt.Value.OptionType()
		fx.optionsShort[string(opt.ShortName)] = opt
	}
}

// errUnexpectedItemType is returned when an unexpected item type is encountered.
var errUnexpectedItemType = errors.New("unexpected item type")

// Parse parses the command line arguments. The arguments must
// not contain the program name. Depending on how the [*FlagSet]
// is configured, this method may call [os.Exit] or panic in
// case of parsing errors (see [NewFlagSet]).
//
// When using [ContinueOnError] this function returns:
//
//  1. any error that [*parser.Parser.Parse] may return.
//
//  2. nil on success.
//
// With [ExitOnError] or [PanicOnError] this function does
// not ever return a non-nil error.
//
// This function may panic on internal errors.
func (fx *FlagSet) Parse(arguments []string) error {
	return fx.maybeHandleError(fx.parse(arguments))
}

func (fx *FlagSet) parse(arguments []string) error {
	// Generate the argv from the program name and the arguments.
	argv := append([]string{fx.progname}, arguments...)

	// Parse the argv into items
	items, err := fx.parser.Parse(argv)

	// Handle the user requesting for help
	if errors.Is(err, parser.ErrHelp) {
		fmt.Fprintf(fx.stdout, "%s\n", fx.Usage())
		// fallthrough
	}

	// Handle any other error
	if err != nil {
		return err
	}

	// Ensure the results are consistent
	assert.True(len(items) >= 0, "expected at least one item")
	_, ok := items[0].(parser.ProgramNameItem)
	assert.True(ok, "expected program name item")
	items = items[1:]

	// Process each item in the command line
	for _, item := range items {
		switch item := item.(type) {

		// Handle option
		case parser.OptionItem:
			option, found := fx.findOption(item.Name)
			assert.True(found, "option not found but was successfully parsed")
			if err := option.Value.Set(item.Value); err != nil {
				return fmt.Errorf("when setting value %q for option %q: %w", item.Value, item.Name, err)
			}
			option.Modified = true // helps with testing

		// Handle argument
		case parser.ArgumentItem:
			fx.args = append(fx.args, item.Value)

		// Ignore any other item type
		default:
			continue
		}
	}
	return nil
}

var exitfn = os.Exit // for testing

func (fx *FlagSet) maybeHandleError(err error) error {
	switch {
	case err == nil:
		return nil

	case fx.handling == ContinueOnError:
		return err

	case fx.handling == ExitOnError && errors.Is(err, parser.ErrHelp):
		exitfn(0)
		panic(err) // just in case exitfn does not exit

	case fx.handling == ExitOnError:
		fmt.Fprintf(fx.stderr, "%s: %s\n", fx.progname, err.Error())
		if lpref := fx.firstLongOptionsPrefix(); lpref != "" {
			fmt.Fprintf(fx.stderr, "Try '%s %shelp' for more help.\n", fx.progname, lpref)
		}
		exitfn(2)
		panic(err) // just in case exitfn does not exit

	default:
		panic(err) // catches also any unhandled fx.handling value
	}
}

func (fx *FlagSet) findOption(name string) (*Option, bool) {
	if option, found := fx.optionsShort[name]; found {
		return option, true
	}
	if option, found := fx.optionsLong[name]; found {
		return option, true
	}
	return nil, false
}
