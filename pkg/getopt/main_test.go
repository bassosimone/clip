// main_test.go - getopt tests.
// SPDX-License-Identifier: GPL-3.0-or-later

package getopt

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMain(t *testing.T) {
	tests := []struct {
		name    string
		argv    []string
		want    []string
		wantErr error
	}{
		{
			name: "success: valid short and long options",
			argv: []string{
				"getopt",
				"-o", "a:bc",
				"--longoptions", "alpha:,beta,gamma",
				"--",
				"-a", "value", "log.txt", "-bc",
			},
			want: []string{
				"-a",
				"value",
				"-b",
				"-c",
				"log.txt",
			},
			wantErr: nil,
		},

		{
			name: "error: missing separator before argument",
			argv: []string{
				"getopt",
				"-o", "abc",
				"argument", // argument without separator
			},
			want:    nil,
			wantErr: ErrExpectedSeparator,
		},

		{
			name: "error: unknown option x before separator",
			argv: []string{
				"getopt",
				"-o", "abc",
				"-x", // invalid option
				"--",
				"-a", "-b", "-c",
			},
			want:    nil,
			wantErr: errors.New("unknown option: -x"),
		},

		{
			name:    "invocation without arguments",
			argv:    []string{"getopt"},
			want:    []string{},
			wantErr: nil,
		},

		{
			name: "no options to parse",
			argv: []string{
				"getopt",
				"-o", "abc",
				"--",
			},
			want:    []string{},
			wantErr: nil,
		},

		{
			name: "error: unknown option after separator",
			argv: []string{
				"getopt",
				"-o", "a:bc",
				"--longoptions", "alpha:,beta,gamma",
				"--",
				"-a", "value", "log.txt", "--foo", "-bc",
			},
			want:    nil,
			wantErr: errors.New("unknown option: --foo"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Main(tt.argv)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("Main() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr != nil {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr.Error()) {
					t.Errorf("Main() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Main() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
