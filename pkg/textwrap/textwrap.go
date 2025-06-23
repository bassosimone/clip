// textwrap.go - text wrapping utils.
// SPDX-License-Identifier: GPL-3.0-or-later

package textwrap

import "strings"

// Do wraps text to the given width with the specified indentation
func Do(text string, width int, indent string) string {
	words := strings.Fields(text)
	if len(words) <= 0 {
		return ""
	}

	var lines []string
	current := indent + words[0]

	for _, word := range words[1:] {
		if len(current)+1+len(word) <= width {
			current += " " + word
			continue
		}
		lines = append(lines, current)
		current = indent + word
	}
	lines = append(lines, current)

	return strings.Join(lines, "\n")
}
