# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code)
when working with code in this repository.

## Overview

Clip is a flexible command line parser written in Go that supports
nested subcommands and configurable flag parsing styles. It enables
implementing CLI tools with different parsing conventions (GNU
getopt, dig-style, etc.) per subcommand.

## Build and Test Commands

### Running Tests

```bash
go test -race -count 1 -cover ./...
```

### Running a Single Test

```bash
# For a specific test function
go test -race -run TestFunctionName ./path/to/package

# For a specific package
go test -race -count 1 -cover ./pkg/nflag
```

### Building Example Commands

```bash
# Build the minirbmk example
go build -o minirbmk ./cmd/minirbmk

# Run minirbmk subcommands
./minirbmk curl --help
./minirbmk dig +short example.com
./minirbmk git status
```

## Architecture

### Package Hierarchy

The codebase follows a layered architecture from
low-level to high-level:

1. **pkg/scanner**: Tokenizes command-line arguments into
option tokens, positional arguments, and separators

2. **pkg/nparser**: Parses tokenized options using configuration
from higher-level packages

3. **pkg/nflag**: Provides stdlib-like flag parsing
API (uses nparser)

4. **pkg/getopt**: GNU getopt-compatible
implementation (uses nparser)

5. **pkg/pflagcompat**: Provides spf13/pflag-compatible
API subset wrapper around nflag (uses nflag)

6. **clip** (root): High-level command dispatcher
integrating nflag with subcommand support

Supporting packages:

1. **pkg/assert**: Runtime assertions that panic on failure

### Core Types

The clip package uses a generic command system parameterized
by `ExecEnv` (execution environment) for testability:

- **Command[T ExecEnv]**: Generic interface for all commands

- **RootCommand[T]**: Entry point wrapping the top-level command

- **DispatcherCommand[T]**: Dispatches to subcommands, handles `help` and `version`

- **LeafCommand[T]**: Terminal command that parses flags using FlagSet

Key points:

- Commands are generic over `ExecEnv` which
mocks I/O functions (stdin/stdout/stderr/exit/args)

- `StdlibExecEnv` is the production implementation
using standard library

- All commands implement the `Command[T]` interface

### Argument Reordering

DispatcherCommand can reorder arguments to handle cases like:

```bash
tool -czf archive.tar.gz tar file1.txt file2.txt
```
becoming:

```bash
tool tar -czf archive.tar.gz file1.txt file2.txt
```

This requires configuring `OptionPrefixes` and
`OptionsArgumentsSeparator` fields on the dispatcher.

### Flag Parsing Customization

Each LeafCommand creates its own FlagSet with customizable:

- Option prefixes (e.g., `-`, `--`, `+`, `/`)

- Positional argument constraints (min/max counts)

- Error handling policy (ContinueOnError, ExitOnError,
PanicOnError)

Example: The minirbmk/dig.go subcommand uses `+`
prefix to emulate dig's flag style.

## Development Patterns

### Creating New Subcommands

1. Define a `LeafCommand[*clip.StdlibExecEnv]` with:

- BriefDescriptionText (required)

- HelpFlagValue (e.g., "--help")

- RunFunc that creates a FlagSet, defines flags, parses, and executes

2. Register in parent DispatcherCommand's Commands map

3. Inside RunFunc:

- Create FlagSet with `nflag.NewFlagSet(args.CommandName, errorHandling)`

- Set Description, PositionalArgumentsUsage, Min/MaxPositionalArgs

- Define flags with BoolFlag(), StringFlag(), Int64Flag()

- Call AutoHelp() for automatic help handling

- Parse with `fset.Parse(args.Args)`

- Access parsed values and positional args

### Migrating from spf13/pflag

If migrating from spf13/pflag, use pkg/pflagcompat for mechanical migration (only
works for a subset of the spf13/pflag package API):

1. Replace import: `"github.com/spf13/pflag"` → `"github.com/bassosimone/clip/pkg/pflagcompat"`

2. Replace constructor: `pflag.NewFlagSet()` → `pflagcompat.NewFlagSet()`

3. Access underlying nflag.FlagSet via `fset.Set` field for nflag-specific features

4. The API is identical: `Bool()`, `BoolP()`, `BoolVar()`, `BoolVarP()`, `String()`, `StringP()`, etc.

See pkg/pflagcompat/example_test.go for comprehensive examples.

### Testing Commands

Use the generic ExecEnv pattern:

- Create test environment implementing ExecEnv interface

- Mock stdin/stdout/stderr with buffers

- Test command execution by calling Run() with test args

- Verify output and exit codes

See `dispatcher_test.go`, `leaf_test.go` for examples.

## Minimum Go Version

Go 1.25

## License

GPL-3.0-or-later
