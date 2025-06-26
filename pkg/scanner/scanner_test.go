// scanner_test.go - Tests for command line scanner.
// SPDX-License-Identifier: GPL-3.0-or-later

package scanner

import "testing"

func TestTokenIndex(t *testing.T) {
	tests := []struct {
		name     string
		token    Token
		expected int
	}{
		{
			name:     "ProgramNameToken",
			token:    ProgramNameToken{Idx: 1},
			expected: 1,
		},
		{
			name:     "OptionToken",
			token:    OptionToken{Idx: 1},
			expected: 1,
		},
		{
			name:     "ArgumentToken",
			token:    PositionalArgumentToken{Idx: 1},
			expected: 1,
		},
		{
			name:     "SeparatorToken",
			token:    OptionsArgumentsSeparatorToken{Idx: 1},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.token.Index()
			if got != tt.expected {
				t.Errorf("Token.Index() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestTokenString(t *testing.T) {
	tests := []struct {
		name     string
		token    Token
		expected string
	}{
		{
			name:     "ProgramNameToken",
			token:    ProgramNameToken{Name: "test"},
			expected: "test",
		},
		{
			name:     "OptionToken with single dash",
			token:    OptionToken{Prefix: "-", Name: "v"},
			expected: "-v",
		},
		{
			name:     "OptionToken with double dash",
			token:    OptionToken{Prefix: "--", Name: "verbose"},
			expected: "--verbose",
		},
		{
			name:     "ArgumentToken",
			token:    PositionalArgumentToken{Value: "file.txt"},
			expected: "file.txt",
		},
		{
			name:     "SeparatorToken",
			token:    OptionsArgumentsSeparatorToken{Separator: "--"},
			expected: "--",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.token.String()
			if got != tt.expected {
				t.Errorf("Token.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestScannerMissingProgramName(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "nil args",
			args: nil,
		},
		{
			name: "empty args",
			args: []string{},
		},
	}

	scanner := &Scanner{
		Prefixes:   []string{"-", "--"},
		Separators: []string{"--"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := scanner.Scan(tt.args)
			if err != ErrMissingProgramName {
				t.Errorf("Scanner.Scan() error = %v, want %v", err, ErrMissingProgramName)
			}
		})
	}
}

func TestScannerZeroLengthOption(t *testing.T) {
	scanner := &Scanner{
		Prefixes:   []string{"-"},
		Separators: []string{},
	}

	args := []string{"prog", "-"}
	tokens, err := scanner.Scan(args)
	if err != nil {
		t.Fatalf("Scanner.Scan() error = %v", err)
	}

	if len(tokens) != 2 {
		t.Errorf("Expected 1 token, got %d", len(tokens))
	}

	if _, ok := tokens[0].(ProgramNameToken); !ok {
		t.Errorf("Expected ProgramNameToken, got %T", tokens[0])
	}

	if _, ok := tokens[1].(PositionalArgumentToken); !ok {
		t.Errorf("Expected PositionalArgumentToken, got %T", tokens[1])
	}
}
