// option_test.go - Tests for Option
// SPDX-License-Identifier: GPL-3.0-or-later

package flag

import (
	"testing"

	"github.com/bassosimone/clip/pkg/parser"
)

func TestOption_FormatParamName(t *testing.T) {
	type testcase struct {
		name    string
		opttype parser.OptionType
		pname   string
		expect  string
	}

	cases := []testcase{
		{
			name:    "with a boolean option and no param name",
			opttype: parser.OptionTypeBool,
			pname:   "",
			expect:  "",
		},

		{
			name:    "with a boolean option and a given param name",
			opttype: parser.OptionTypeBool,
			pname:   "ANTANI",
			expect:  "",
		},

		{
			name:    "with a string option and no param name",
			opttype: parser.OptionTypeString,
			pname:   "",
			expect:  "VALUE",
		},

		{
			name:    "with a string option and a given param name",
			opttype: parser.OptionTypeString,
			pname:   "ANTANI",
			expect:  "ANTANI",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opt := &Option{
				LongName:  "verbose",
				Modified:  false,
				ParamName: tc.pname,
				ShortName: 'v',
				Usage:     "",
				Value: &mockValue{
					value:      "",
					optionType: tc.opttype,
					setError:   nil,
				},
			}
			paramName := opt.FormatParamName()
			if paramName != tc.expect {
				t.Errorf("expected %q, got %q", tc.expect, paramName)
			}
		})
	}
}
