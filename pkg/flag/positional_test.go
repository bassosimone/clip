// positional_test.go - Tests for positional arguments checking
// SPDX-License-Identifier: GPL-3.0-or-later

package flag

import (
	"errors"
	"testing"
)

func TestFlagSet_NArg(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected int
	}{
		{
			name:     "no args",
			args:     []string{},
			expected: 0,
		},

		{
			name:     "one arg",
			args:     []string{"foo"},
			expected: 1,
		},

		{
			name:     "multiple args",
			args:     []string{"foo", "bar", "baz"},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := NewFlagSet("test", ContinueOnError)
			fs.args = tt.args
			got := fs.NArg()
			if got != tt.expected {
				t.Errorf("NArg() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestFlagSet_PositionalArgsRangeCheck(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		min       int
		max       int
		wantErr   error
		wantIsErr bool
	}{
		{
			name:      "within range",
			args:      []string{"foo", "bar"},
			min:       2,
			max:       3,
			wantErr:   nil,
			wantIsErr: false,
		},

		{
			name:      "exact min and max, correct count",
			args:      []string{"foo", "bar"},
			min:       2,
			max:       2,
			wantErr:   nil,
			wantIsErr: false,
		},

		{
			name:      "too few args",
			args:      []string{"foo"},
			min:       2,
			max:       3,
			wantErr:   ErrTooFewPositionalArgs,
			wantIsErr: true,
		},

		{
			name:      "too many args",
			args:      []string{"foo", "bar", "baz", "qux"},
			min:       1,
			max:       3,
			wantErr:   ErrTooManyPositionalArgs,
			wantIsErr: true,
		},

		{
			name:      "exact match required, wrong count",
			args:      []string{"foo"},
			min:       2,
			max:       2,
			wantErr:   ErrUnexpectedNumberOfPositionalArgs,
			wantIsErr: true,
		},

		{
			name:      "exact match required, too many",
			args:      []string{"foo", "bar", "baz"},
			min:       2,
			max:       2,
			wantErr:   ErrUnexpectedNumberOfPositionalArgs,
			wantIsErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := NewFlagSet("test", ContinueOnError)
			fs.args = tt.args
			err := fs.PositionalArgsRangeCheck(tt.min, tt.max)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("PositionalArgsRangeCheck() unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("PositionalArgsRangeCheck() expected error, got nil")
				} else if !errors.Is(err, tt.wantErr) {
					t.Errorf("PositionalArgsRangeCheck() error = %v, want error wrapping %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestFlagSet_PositionalArgsEqualCheck(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		n         int
		wantErr   error
		wantIsErr bool
	}{
		{
			name:      "equal args",
			args:      []string{"foo", "bar"},
			n:         2,
			wantErr:   nil,
			wantIsErr: false,
		},

		{
			name:      "too few args",
			args:      []string{"foo"},
			n:         2,
			wantErr:   ErrUnexpectedNumberOfPositionalArgs,
			wantIsErr: true,
		},

		{
			name:      "too many args",
			args:      []string{"foo", "bar", "baz"},
			n:         2,
			wantErr:   ErrUnexpectedNumberOfPositionalArgs,
			wantIsErr: true,
		},

		{
			name:      "zero args expected, zero given",
			args:      []string{},
			n:         0,
			wantErr:   nil,
			wantIsErr: false,
		},

		{
			name:      "zero args expected, some given",
			args:      []string{"foo"},
			n:         0,
			wantErr:   ErrUnexpectedNumberOfPositionalArgs,
			wantIsErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := NewFlagSet("test", ContinueOnError)
			fs.args = tt.args
			err := fs.PositionalArgsEqualCheck(tt.n)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("PositionalArgsEqualCheck() unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("PositionalArgsEqualCheck() expected error, got nil")
				} else if !errors.Is(err, tt.wantErr) {
					t.Errorf("PositionalArgsEqualCheck() error = %v, want error wrapping %v", err, tt.wantErr)
				}
			}
		})
	}
}
