// optiontype.go - option type definition.
// SPDX-License-Identifier: GPL-3.0-or-later

package nparser

// OptionType is the type of an [Option].
type OptionType int64

const (
	optionKindEarly = OptionType(1 << (iota + 4))
	optionKindStandalone
	optionKindGroupable
)

const (
	optionArgumentNone = OptionType(1 << iota)
	optionArgumentRequired
	optionArgumentOptional
)

func (ot OptionType) isEarly() bool {
	return (ot & optionKindEarly) != 0
}

func (ot OptionType) isStandalone() bool {
	return (ot & optionKindStandalone) != 0
}

func (ot OptionType) isGroupable() bool {
	return (ot & optionKindGroupable) != 0
}

// These constants define the allowed [OptionType] values.
const (
	// OptionTypeEarlyArgumentNone indicates an early option requiring no arguments.
	OptionTypeEarlyArgumentNone = optionKindEarly | optionArgumentNone

	// OptionTypeStandaloneArgumentNone indicates a standalone option requiring no arguments.
	OptionTypeStandaloneArgumentNone = optionKindStandalone | optionArgumentNone

	// OptionTypeStandaloneArgumentRequired indicates a standalone option requiring an argument.
	OptionTypeStandaloneArgumentRequired = optionKindStandalone | optionArgumentRequired

	// OptionTypeStandaloneArgumentOptional indicates a standalone option with an optional argument.
	OptionTypeStandaloneArgumentOptional = optionKindStandalone | optionArgumentOptional

	// OptionTypeGroupableArgumentNone indicates a groupable option requiring no arguments.
	OptionTypeGroupableArgumentNone = optionKindGroupable | optionArgumentNone

	// OptionTypeGroupableArgumentRequired indicates groupable option requiring an argument.
	OptionTypeGroupableArgumentRequired = optionKindGroupable | optionArgumentRequired
)
