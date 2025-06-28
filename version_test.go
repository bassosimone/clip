// version_test.go - tests for version handling
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import (
	"context"
	"errors"
	"testing"

	"github.com/bassosimone/clip/pkg/nparser"
)

func TestVersionCommand(t *testing.T) {
	t.Run("BriefDescription", func(t *testing.T) {
		expect := "foo"
		vc := &VersionCommand[*StdlibExecEnv]{BriefDescriptionText: expect}
		got := vc.BriefDescription()
		if got != expect {
			t.Fatalf("BriefDescription() = %q, want %q", got, expect)
		}
	})

	t.Run("HelpFlag when explicitly set", func(t *testing.T) {
		expect := "/help"
		vc := &VersionCommand[*StdlibExecEnv]{HelpFlagValue: expect}
		got := vc.HelpFlag()
		if got != expect {
			t.Fatalf("HelpFlag() = %q, want %q", got, expect)
		}
	})

	t.Run("HelpFlag default value", func(t *testing.T) {
		expect := "--help"
		vc := &VersionCommand[*StdlibExecEnv]{}
		got := vc.HelpFlag()
		if got != expect {
			t.Fatalf("HelpFlag() = %q, want %q", got, expect)
		}
	})

	t.Run("SupportsSubcommands", func(t *testing.T) {
		expect := false
		vc := &VersionCommand[*StdlibExecEnv]{}
		got := vc.SupportsSubcommands()
		if got != expect {
			t.Fatalf("SupportsSubcommands() = %v, want %v", got, expect)
		}
	})

	t.Run("Run with unexpected flags", func(t *testing.T) {
		vc := &VersionCommand[*StdlibExecEnv]{}
		args := &CommandArgs[*StdlibExecEnv]{
			Args:        []string{"--unexpected"},
			Command:     vc,
			CommandName: "version",
			Env:         NewStdlibExecEnv(),
		}
		err := vc.Run(context.Background(), args)
		var unknownOption nparser.ErrUnknownOption
		if !errors.As(err, &unknownOption) {
			t.Fatalf("Run() = %v, want %v", err, unknownOption)
		}
	})

	t.Run("Run with unexpected positional arguments", func(t *testing.T) {
		vc := &VersionCommand[*StdlibExecEnv]{}
		args := &CommandArgs[*StdlibExecEnv]{
			Args:        []string{"unexpected"},
			Command:     vc,
			CommandName: "version",
			Env:         NewStdlibExecEnv(),
		}
		err := vc.Run(context.Background(), args)
		var unknownPositional nparser.ErrTooManyPositionalArguments
		if !errors.As(err, &unknownPositional) {
			t.Fatalf("Run() = %v, want %v", err, unknownPositional)
		}
	})

	t.Run("Run with no arguments and no flags", func(t *testing.T) {
		vc := &VersionCommand[*StdlibExecEnv]{
			Version: "0.1.0",
		}
		args := &CommandArgs[*StdlibExecEnv]{
			Args:        []string{},
			Command:     vc,
			CommandName: "version",
			Env:         NewStdlibExecEnv(),
		}
		err := vc.Run(context.Background(), args)
		if err != nil {
			t.Fatal(err)
		}
	})
}
