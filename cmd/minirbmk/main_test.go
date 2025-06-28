// main.go - Main for the minirbmk example
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"errors"
	"fmt"
	"testing"
)

type errExitStatus struct {
	code int
}

func (err errExitStatus) Error() string {
	return fmt.Sprintf("exit: %d", err.code)
}

func Test_main(t *testing.T) {
	type testcase struct {
		argv   []string
		expect int
	}

	cases := []testcase{

		{
			argv:   []string{"minirbmk"},
			expect: 0,
		},

		{
			argv:   []string{"minirbmk", "--help"},
			expect: 0,
		},

		{
			argv:   []string{"minirbmk", "-h"},
			expect: 0,
		},

		{
			argv:   []string{"minirbmk", "--version"},
			expect: 0,
		},

		{
			argv:   []string{"minirbmk", "dig", "-h"},
			expect: 0,
		},

		{
			argv:   []string{"minirbmk", "dig", "+short", "-4", "www.google.com"},
			expect: 0,
		},

		{
			argv:   []string{"minirbmk", "dig", "+wrong", "-4", "www.google.com"},
			expect: 2,
		},

		{
			argv:   []string{"minirbmk", "dig", "+short", "-4"},
			expect: 2,
		},

		{
			argv:   []string{"minirbmk", "curl", "--help"},
			expect: 0,
		},

		{
			argv:   []string{"minirbmk", "curl", "-c", "xo", "--cacert", "yo", "-v", "URL"},
			expect: 0,
		},

		{
			argv:   []string{"minirbmk", "curl", "-Z", "xo", "--cacert", "yo", "-v", "URL"},
			expect: 2,
		},

		{
			argv:   []string{"minirbmk", "curl", "-c", "xo", "--cacert", "yo", "-v"},
			expect: 2,
		},

		{
			argv:   []string{"minirbmk", "git", "--help"},
			expect: 0,
		},

		{
			argv:   []string{"minirbmk", "git", "init", "--help"},
			expect: 0,
		},

		{
			argv:   []string{"minirbmk", "git", "init", "-qb", "main"},
			expect: 0,
		},

		{
			argv:   []string{"minirbmk", "git", "init", "-qb", "main", "-z"},
			expect: 2,
		},

		{
			argv:   []string{"minirbmk", "git", "init", "-qb", "main", "z"},
			expect: 0,
		},

		{
			argv:   []string{"minirbmk", "git", "init", "-qb", "main", "z", "z"},
			expect: 2,
		},

		{
			argv:   []string{"minirbmk", "git", "clone", "--help"},
			expect: 0,
		},

		{
			argv:   []string{"minirbmk", "git", "clone", "-qb", "main", "REPO"},
			expect: 0,
		},

		{
			argv:   []string{"minirbmk", "git", "clone", "-qb", "main", "REPO", "DIR"},
			expect: 0,
		},

		{
			argv:   []string{"minirbmk", "git", "clone", "-qb", "main", "REPO", "-z"},
			expect: 2,
		},

		{
			argv:   []string{"minirbmk", "git", "clone", "-qb", "main"},
			expect: 2,
		},

		{
			argv:   []string{"minirbmk", "git", "__antani__", "--help"},
			expect: 2,
		},

		{
			argv:   []string{"minirbmk", "__antani__", "--help"},
			expect: 2,
		},
	}

	// modify the environment to collect the exit status
	env.OSExit = func(exitcode int) {
		panic(errExitStatus{exitcode})
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.argv), func(t *testing.T) {

			// modify the argv
			env.OSArgs = tc.argv

			// prepare to handle the panic and check the status code
			defer func() {
				r := recover()

				// handle the case where exit was called
				if r != nil {
					err := r.(error)
					var exitStatus errExitStatus
					if !errors.As(err, &exitStatus) {
						t.Fatal("unexpected error", err)
					}
					if exitStatus.code != tc.expect {
						t.Fatal("expected", tc.expect, "got", exitStatus.code)
					}
					return
				}

				// otherwise we must have 0 expected exit status
				if tc.expect != 0 {
					t.Fatal("exit not called but we expected", tc.expect, "as exit code")
				}
			}()

			// run the main function
			main()
		})
	}
}
