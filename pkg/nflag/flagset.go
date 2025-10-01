// flagset.go - FlagSet implementation for command-line flag parsing
// SPDX-License-Identifier: GPL-3.0-or-later

package nflag

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/nparser"
)

// --- types ---

// ErrorHandling controls [*FlagSet.Parse] error handling.
type ErrorHandling int

// These constants define the allowed [ErrorHandling] values.
const (
	// ContinueOnError causes [*FlagSet] to return the parse error.
	ContinueOnError = ErrorHandling(iota)

	// ExitOnError causes [*Flagset] to call Exit with code 2 on error.
	ExitOnError

	// PanicOnError causes [*FlagSet] to panic on error.
	PanicOnError
)

// Flag is a long or short flag managed by the [*FlagSet].
//
// You can add a [*Flag] to a [*FlagSet] using, e.g.:
//
//  1. [*FlagSet.BoolVar] to add a boolean flag.
//
//  2. [*FlagSet.StringVar] to add a string flag.
//
// These methods usually add two flags per invocation: a long
// flag and a short flag. See also their documentation.
type Flag struct {
	// Modified indicates whether the flag was modified as
	// a side effect of parsing the command line.
	Modified bool

	// Option is the related parser option.
	Option *nparser.Option

	// TakesArg is true if this flag takes an argument.
	TakesArg bool

	// Usage contains the usage message.
	Usage string

	// Value is the value assigned-to when parsing.
	Value Value
}

// LongShortFlag contains a long and a short flag that are logically
// bound to setting the same [Value]. For example, when you call
// [*FlagSet.BoolVar] with a nonempty long name and a nonempty short
// name, this method creates a [LongShortFlag] pointing to a valid
// long flag and to a valid short flag.
type LongShortFlag struct {
	// LongFlag is the long flag.
	//
	// Warning: it may be nil.
	LongFlag *Flag

	// ShortFlag is the short flag.
	//
	// Warning: it may be nil.
	ShortFlag *Flag

	// TakesArg is true if the flags takes an argument.
	TakesArg bool

	// Usage is the usage string.
	Usage string

	// Value is the value assigned-to when parsing.
	Value Value
}

// FlagSet allows to parse flags from the command line. The zero value is not
// ready to use. Construct using the [NewFlagSet] constructor.
//
// Note: a [*FlagSet] where you have not added any flags through methods
// like [*FlagSet.BoolVar] or [*FlagSet.StringVar] defaults to parsing options
// using GNU conventions. That is, if you run the related program with:
//
//	./program --verbose
//
// The [*FlagSet] will recognize `--verbose` as a syntactically valid flag
// that has not been configured and print an "unknown flag" error.
type FlagSet struct {
	// Description is the program description used when printing the usage.
	//
	// [NewFlagSet] initializes this field to "".
	Description string

	// DisablePermute disable the permutation of options and arguments.
	//
	// [NewFlagSet] initializes this field to false.
	//
	// Consider the following command line:
	//
	// 	curl https://www.google.com/ -H 'Host: google.com'
	//
	// The default behavior is to permute this to:
	//
	// 	curl -H 'Host: google.com' https://www.google.com/
	//
	// However, when DisablePermute is true, we keep the command
	// line unmodified. While permuting is a nice-to-have property
	// in general, consider instead the following case:
	//
	// 	multirepo foreach -kx git status -v
	//
	// With permutation, this command line would become:
	//
	// 	multirepo foreach git status -kv -v
	//
	// This is not the desired behavior if the foreach command
	// takes another command and its options as arguments.
	//
	// To make the above command line work with permutation, a
	// user would instead need to write this:
	//
	// 	multirepo foreach -kv -- git status -v
	//
	// By setting DisablePermute to true, the `--` separator
	// becomes unnecessary and the UX is improved.
	DisablePermute bool

	// Examples contains examples when printing the usage.
	//
	// [NewFlagSet] initializes this field to "".
	Examples string

	// Exit is the function to call with the [ExitOnError] policy.
	//
	// [NewFlagSet] initializes this field to [os.Exit].
	Exit func(status int)

	// ErrorHandling is the [ErrorHandling] policy.
	//
	// [NewFlagSet] initializes this field to [ContinueOnError].
	ErrorHandling ErrorHandling

	// LongFlagPrefix is the prefix for parsing long flags.
	//
	// [NewFlagSet] initializes this field to "--".
	//
	// The default configuration is compatible with the GNU standards
	// where long flags are like `--verbose`, `--output <file>`.
	//
	// Because each [*Flag] stores its own prefix, modifying this
	// field after creating flags do not retroactively change their
	// prefix definition. This may be useful to incrementally make
	// a [*FlagSet] handling mixed prefixes.
	LongFlagPrefix string

	// MaxPositionalArgs is the maximum number of positional arguments.
	//
	// [NewFlagSet] initializes this field to [math.MaxInt].
	//
	// The default configuration, thus, allows for any number of
	// positional arguments to be on the command line.
	MaxPositionalArgs int

	// MinPositionalArgs is the minimum number of positional arguments.
	//
	// [NewFlagSet] initializes this field to 0.
	//
	// The default configuration, thus, allows for any number of
	// positional arguments to be on the command line.
	MinPositionalArgs int

	// OptionsArgumentsSeparator separates options and arguments.
	//
	// [NewFlagSet] initializes this field to "--".
	//
	// The default configuration is compatible with the GNU standards
	// where "--" instructs getopt to stop processing flags and to treat
	// all the remaining entries as positional arguments.
	OptionsArgumentsSeparator string

	// PositionalArgumentsUsage is the usage string for postional arguments.
	//
	// [NewFlagSet] initializes this field to "arg ..."
	PositionalArgumentsUsage string

	// ProgramName is the program name.
	//
	// [NewFlagSet] initializes this field to the given program name.
	ProgramName string

	// ShortFlagPrefix is the prefix for parsing short flags.
	//
	// [NewFlagSet] initializes this field to "-".
	//
	// The default configuration is compatible with the GNU standards
	// where long flags are like `-v`, `-o <file>`, `-vo <file>`.
	//
	// Because each [*Flag] stores its own prefix, modifying this
	// field after creating flags do not retroactively change their
	// prefix definition. This may be useful to incrementally make
	// a [*FlagSet] handling mixed prefixes.
	ShortFlagPrefix string

	// Stderr is the [io.Writer] to use as the stderr.
	//
	// [NewFlagSet] initializes this field to [os.Stderr].
	//
	// We use this field with [ExitOnError] policy.
	Stderr io.Writer

	// Stdout is the [io.Writer] to use as the stdout.
	//
	// [NewFlagSet] initializes this field to [os.Stdout].
	//
	// We use this field with [ExitOnError] policy.
	Stdout io.Writer

	// parserView organizes flags for parsing.
	parserView map[string]*Flag

	// positional buffers the positional arguments.
	positionals []string

	// usaveView organizes flags for printing the usage message.
	usageView []LongShortFlag
}

// --- constructor ---

// NewFlagSet returns a new [*FlagSet] instance. We use the given progname as
// the ProgramName field and the given handling as the ErrorHandling field. We
// initialize all the other fields using sensible defaults. We document these
// defaults in the [*FlagSet] documentation.
//
// This function panics if the program name is empty. Using an unknown value
// for the handling parameter is equivalent to using [PanicOnError].
func NewFlagSet(progname string, handling ErrorHandling) *FlagSet {
	// make sure the program name is not empty
	assert.True(progname != "", "program name must not be empty")

	// create with default settings
	return &FlagSet{
		Description:               "",
		DisablePermute:            false,
		ErrorHandling:             handling,
		Examples:                  "",
		Exit:                      os.Exit,
		LongFlagPrefix:            "--",
		MaxPositionalArgs:         math.MaxInt,
		MinPositionalArgs:         0,
		ProgramName:               progname,
		OptionsArgumentsSeparator: "--",
		PositionalArgumentsUsage:  "arg ...",
		ShortFlagPrefix:           "-",
		Stderr:                    os.Stderr,
		Stdout:                    os.Stdout,
		parserView:                map[string]*Flag{},
		positionals:               []string{},
		usageView:                 []LongShortFlag{},
	}
}

// --- getters ---

// Args returns the positional arguments collected by [*FlagSet.Parse].
func (fx *FlagSet) Args() []string {
	return fx.positionals
}

// LookupFlagLong returns the long [*Flag] associated with the given name.
func (fx *FlagSet) LookupFlagLong(name string) (*Flag, bool) {
	flag, ok := fx.parserView[name]
	return flag, ok
}

// LookupFlagShort returns the short [*Flag] associated with the given name.
func (fx *FlagSet) LookupFlagShort(name byte) (*Flag, bool) {
	flag, ok := fx.parserView[string(name)]
	return flag, ok
}

// Flags returns the [LongShortFlag] set of defined flag. Beware that
// for some flags either the long or short pointer may be nil.
func (fx *FlagSet) Flags() []LongShortFlag {
	return fx.usageView
}

// --- parsing code ---

// Parse parses the given command line arguments, It assigns positional arguments
// and each flag [Value] as a side effect of parsing.
//
// The args MUST NOT contain the program name. That is, if there are no command
// line arguments beyond the program name, arguments must be empty.
//
// Depending on the [ErrorHandling] policy, on failure, this method may return the
// error invoke [os.Exit], or call panic with the error that occurred.
func (fx *FlagSet) Parse(args []string) error {
	return fx.maybeHandleError(fx.parse(args))
}

// ErrHelp is the error returned in case the user requested for `help`.
//
// Use [*FlagSet.AutoHelp] to enable recognizing help flags.
var ErrHelp = errors.New("help requested")

func (fx *FlagSet) parse(args []string) error {
	// create an argument vector that includes the program name
	argv := make([]string, 0, 1+len(args))
	argv = append(argv, fx.ProgramName)
	argv = append(argv, args...)

	// configure the command line parser
	px := &nparser.Parser{
		DisablePermute:            fx.DisablePermute,
		MaxPositionalArguments:    fx.MaxPositionalArgs,
		MinPositionalArguments:    fx.MinPositionalArgs,
		OptionsArgumentsSeparator: fx.OptionsArgumentsSeparator,
		Options:                   []*nparser.Option{},
	}
	for _, view := range fx.parserView {
		px.Options = append(px.Options, view.Option)
	}

	// parse the command line
	values, err := px.Parse(argv)
	if err != nil {
		return err
	}

	// map the parsed values back to options and positionals
	for _, value := range values {
		switch value := value.(type) {

		// positional argument: just add to the internal slice of positionals
		case nparser.ValuePositionalArgument:
			fx.positionals = append(fx.positionals, value.Value)

		// option: find the corresponding value and attempt to set it
		case nparser.ValueOption:
			// attempt to get the right parser view
			optname := value.Option.Name
			flag, found := fx.parserView[optname]
			assert.True(found, fmt.Sprintf("expected to find flag %q", optname))

			// assign a value to the flag
			if err := flag.Value.Set(value.Value); err != nil {
				return err
			}

			// mark the flag as modified
			flag.Modified = true

			// detect [helpValue] and transform it to [ErrHelp]
			if _, ok := flag.Value.(helpValue); ok {
				return ErrHelp
			}
		}
	}
	return nil
}

func (fx *FlagSet) maybeHandleError(err error) error {
	switch {
	case err == nil:
		return nil

	case fx.ErrorHandling == ContinueOnError:
		return err

	case fx.ErrorHandling == ExitOnError && errors.Is(err, ErrHelp):
		var sb strings.Builder
		fx.PrintUsage(&sb)
		fmt.Fprint(fx.Stdout, sb.String())
		fx.Exit(0)

	case fx.ErrorHandling == ExitOnError:
		fmt.Fprintf(fx.Stderr, "%s: %s\n", fx.ProgramName, err.Error())
		var sb strings.Builder
		fx.PrintHelpHint(&sb)
		if sb.Len() > 0 {
			fmt.Fprint(fx.Stderr, sb.String())
		}
		fx.Exit(2)
	}

	// catch all the remaining cases with a panic
	panic(err)
}

// --- code to register flags ---

func (fx *FlagSet) mustAddLongAndShortFlag(long, short *Flag) {
	// define utility function for adding a single flag
	var (
		takesArg bool
		usage    string
		value    Value
	)
	addx := func(fpv *Flag) {
		if fpv != nil {
			fname := fpv.Option.Name
			_, found := fx.parserView[fname]
			assert.True(!found, fmt.Sprintf("flag %q already defined", fname))
			fx.parserView[fname] = fpv
			takesArg = fpv.TakesArg
			usage = fpv.Usage
			value = fpv.Value
		}
	}

	// ensure at least one of them is defined
	assert.True(long != nil || short != nil, "mustAddLongAndShortFlag passed two nil pointers")

	// add both of them
	addx(long)
	addx(short)

	// update the corresponding usahe view
	fx.usageView = append(fx.usageView, LongShortFlag{
		ShortFlag: short,
		LongFlag:  long,
		TakesArg:  takesArg,
		Usage:     usage,
		Value:     value,
	})
}
