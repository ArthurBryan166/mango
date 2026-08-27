package token

type TokenTipe int

const(
	MNG TokenTipe = iota
	IDENTIFIER
	NUMBER
	STRING

	PLUS
	MULTIPLY
	MINUS
    DIVIDE

	ASSIGN // <-
    EQUAL
    NOT_EQUAL
    GREATER
    GREATER_EQUAL
    LESS
    LESS_EQUAL

    LEFT_PAREN
    RIGHT_PAREN
    LEFT_BRACE
    RIGHT_BRACE

    IF
    ELSE
    FOR
    MANGO // function
    RETURN

    TRUE
    FALSE
    NIL
    EOF
)

type Token struct{
	Type TokenTipe
	Lexeme string
	Line int
}