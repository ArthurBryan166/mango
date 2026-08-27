package token

type TokenTipe int

const(
	MNG TokenTipe = iota
	IDENTIFIER
	ASSIGN // <-
	NUMBER
	PLUS
	MULTIPLY
)