// optiontype_test.go - option type tests.
// SPDX-License-Identifier: GPL-3.0-or-later

package nparser

import "testing"

func TestOptionType(t *testing.T) {
	type testcase struct {
		name         string
		input        OptionType
		isEarly      bool
		isStandalone bool
		isGroupable  bool
	}

	cases := []testcase{
		{
			name:    "OptionTypeEarlyArgumentNone",
			input:   OptionTypeEarlyArgumentNone,
			isEarly: true,
		},

		{
			name:         "OptionTypeStandaloneArgumentNone",
			input:        OptionTypeStandaloneArgumentNone,
			isStandalone: true,
		},

		{
			name:         "OptionTypeStandaloneArgumentRequired",
			input:        OptionTypeStandaloneArgumentRequired,
			isStandalone: true,
		},

		{
			name:         "OptionTypeStandaloneArgumentOptional",
			input:        OptionTypeStandaloneArgumentOptional,
			isStandalone: true,
		},

		{
			name:        "OptionTypeGroupableArgumentNone",
			input:       OptionTypeGroupableArgumentNone,
			isGroupable: true,
		},

		{
			name:        "OptionTypeGroupableArgumentRequired",
			input:       OptionTypeGroupableArgumentRequired,
			isGroupable: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.input.isEarly(); got != tc.isEarly {
				t.Fatal("isEarly: expected", tc.isEarly, "got", got)
			}
			if got := tc.input.isStandalone(); got != tc.isStandalone {
				t.Fatal("isStandalone: expected", tc.isStandalone, "got", got)
			}
			if got := tc.input.isGroupable(); got != tc.isGroupable {
				t.Fatal("isGroupable: expected", tc.isGroupable, "got", got)
			}
		})
	}
}
