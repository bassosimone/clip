// permute.go - argv permutation.
// SPDX-License-Identifier: GPL-3.0-or-later

package nparser

func permute(cfg *config, programName Value, options, positionals []Value) []Value {
	// Determine what to do depending on the configuration
	switch {

	// When permutation is disabled, restore the original token order
	case cfg.disablePermute():
		output := make([]Value, 0, 1+len(options)+len(positionals))
		output = append(output, programName)
		output = append(output, options...)
		output = append(output, positionals...)
		sortValues(output)
		return output

	// Otherwise, merge options together and sort options and arguments independently
	default:
		sortValues(options)
		sortValues(positionals)
		output := make([]Value, 0, 1+len(options)+len(positionals))
		output = append(output, programName)
		output = append(output, options...)
		output = append(output, positionals...)
		return output
	}
}
