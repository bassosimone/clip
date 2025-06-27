// early.go - early option parsing.
// SPDX-License-Identifier: GPL-3.0-or-later

package nparser

import "github.com/bassosimone/clip/pkg/scanner"

func searchEarly(px *Parser, argv []string) (Value, bool) {
	for idx := 1; idx < len(argv); idx++ {
		for _, option := range px.Options {
			if argv[idx] == option.Prefix+option.Name && option.Type.isEarly() {
				// We have found an early option, return it
				ovalue := ValueOption{
					Option: option,
					Tok: scanner.OptionToken{
						Idx:    idx,
						Prefix: option.Prefix,
						Name:   option.Name,
					},
					Value: "",
				}
				return ovalue, true
			}
		}
	}
	return nil, false
}
