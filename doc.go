// doc.go - package documentation.
// SPDX-License-Identifier: GPL-3.0-or-later

/*
Package clip provides a framework for building command-line
applications with subcommands and flags parsing.

To use this package proceed as follows:

 1. Create a [RootCommand] containing either a [DispatcherCommand] or a
    [LeafCommand], depending on whether you need subcommands.

 2. If using a [DispatcherCommand], initialize its Commands map with
    either [LeafCommand] or [DispatcherCommand] instances.

 3. Inside a [LeafCommand] Run function, implement the flags parsing
    logic using a [*FlagSet] created with [clip.NewFlagSet].

See the package examples for more information.

# RootCommand

The [RootCommand] optionally allows to react to signals and otherwise
dispatches the work to the [DispatcherCommand] or [LeafCommand]
contained within.

# DispatcherCommand

A [DispatcherCommand] dispatches the execution to subcommands by
inspecting the command-line arguments and finding a match to a known
subcommand. Typically, the next command to execute is the first
entry of the arguments slice. However, [DispatcherCommand] can also
transparently handle the case where the command appears later on
in the slice. This works as long as there is a single
token in the arguments slice mapping to known subcommands. With more
than one match we return [ErrAmbiguousCommandLine]. With no match,
we return [ErrNoSuchCommand]. The [DispatcherCommand] automatically
handles `--help`, `-h`, `help` and `help COMMAND` requests by
printing its own help or dispatching to subcommands:

 1. `foo --help` and `foo -h` print the help message for `foo`.

 2. `foo help` also prints the help message for `foo`.

 3. `foo help bar` prints the help message for the `bar` subcommand.

There is no flag associated with the [DispatcherCommand]. All the
flags are only associated with the subcommands. However, given that
the command-line can be out of order, as mentioned above, the
overall usage experience is smooth and forgiving. Yet, you need
by convention to ensure command line consistency among tools. But, on
the flip side, the resulting code is straightforward.

# LeafCommand and FlagSet

To parse command line flags in a [LeafCommand], you create a
[*FlagSet] using [clip.NewFlagSet]. The [*FlagSet] uses an API
very similar to the Go standard library `flag` package. However,
by default, [*FlagSet] uses the GNU command line conventions:

 1. Short options start with `-` (e.g., `-v`).

 2. Long options start with `--` (e.g., `--verbose`).

 3. Short options may be combined (e.g., `-vf`).

 4. Short options arguments can be specified in separate arguments
    (e.g., `-vf FILE`) or inline (e.g., `-vfFILE`).

 5. Long options arguments can be specified in separate arguments
    (e.g., `--file FILE`) or inline (e.g., `--file=FILE`).

 6. The `--` separator terminates option processing and all
    subsequent arguments are treated as positional.

Yet, as documented in [pkg/flags] and [pkg/parser], it is possible to
configure an individual [*FlagSet] differently (e.g., to behave
like the Go [flag] package or like `dig`). This functionality is
actually one of the strongest selling point of this package. It
allows one to write tools combining different command-line interfaces
styles. While in principle this is not desirable for a new tool, it
simplifies the job of providing a consistent command-line interface
across different emulated tools. The poster child for this functionality
is a tool named `rbmk` where `rbmk dig` works like `dig` and
`rbmk curl` works like `curl`.

A [*FlagSet] automatically discovers the usage of help flags (by
default `-h` and `--help`) and, overall, this package ensures that
appending `--help` to an otherwise completely broken command line
always displays the help message.

# Testability

All top-level types depend on an abstract T type, bounded by the
[ExecEnv] interface. The default implementation, using the standard
library, but highly customizable is [*StdlibExecEnv]. By using
such an interface, it is possible to write highly testable code
where most of the environment dependencies can be mocked.
*/
package clip
