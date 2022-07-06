package sln

// Token represents a lexical token.
type Token int

const (
	// Unknown characters that doesn't fit in any other group
	Unknown Token = iota
	// EOF - end of file
	EOF
	// WS - whitespace
	WS

	// IDENT - characters
	IDENT

	// ASTERISK *
	ASTERISK
	// COMMA ,
	COMMA
	// OPEN_PAREN (
	OPEN_PAREN
	// CLOSE_PAREN )
	CLOSE_PAREN
	// QUOTE "
	QUOTE
	// EQUAL =
	EQUAL

	// PROJECT - project begin
	PROJECT
	// ENDPROJECT - project end
	ENDPROJECT
)
