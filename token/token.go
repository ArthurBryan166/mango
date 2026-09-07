package token

import (
	"fmt"
)

type TokenType int

const(
	MNG TokenType = iota
    MANGOUT
    MANGOIN
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
    AND
    OR
    NOT

    LEFT_PAREN
    RIGHT_PAREN
    LEFT_BRACE
    RIGHT_BRACE
    COMMA
    SEMICOLON

    IF
    ELSE
    FOR
    MANGO // function
    RETURN

    TRUE
    FALSE
    NIL
    EOF

    NUMBER_TYPE
    STRING_TYPE
    BOOL_TYPE
    VOID_TYPE
)

type Token struct{
	Type TokenType
	Lexeme string
	Line int
}

func (t TokenType) String() string {
    switch t {
    case MNG:
        return "MNG"
    case MANGOUT:
        return "MANGOUT"
    case MANGOIN:
        return "MANGOIN"
    case IDENTIFIER:
        return "IDENTIFIER"
    case NUMBER:
        return "NUMBER"
    case STRING:
        return "STRING"
    case PLUS:
        return "PLUS"
    case MINUS:
        return "MINUS"
    case MULTIPLY:
        return "MULTIPLY"
    case DIVIDE:
        return "DIVIDE"
    case ASSIGN:
        return "ASSIGN"
    case EQUAL:
        return "EQUAL"
    case NOT_EQUAL:
        return "NOT_EQUAL"
    case GREATER:
        return "GREATER"
    case GREATER_EQUAL:
        return "GREATER_EQUAL"
    case LESS:
        return "LESS"
    case LESS_EQUAL:
        return "LESS_EQUAL"
    case AND:
        return "AND"
    case OR:
        return "OR"
    case NOT:
        return "NOT"
    case LEFT_PAREN:
        return "LEFT_PAREN"
    case RIGHT_PAREN:
        return "RIGHT_PAREN"
    case LEFT_BRACE:
        return "LEFT_BRACE"
    case RIGHT_BRACE:
        return "RIGHT_BRACE"
    case COMMA:
        return "COMMA"
    case SEMICOLON:
        return "SEMICOLON"
    case IF:
        return "IF"
    case ELSE:
        return "ELSE"
    case FOR:
        return "FOR"
    case MANGO:
        return "MANGO"
    case RETURN:
        return "RETURN"
    case TRUE:
        return "TRUE"
    case FALSE:
        return "FALSE"
    case NIL:
        return "NIL"
    case EOF:
        return "EOF"
    case NUMBER_TYPE:
        return "NUMBER_TYPE"
    case STRING_TYPE:
        return "STRING_TYPE"
    case BOOL_TYPE:
        return "BOOL_TYPE"
    case VOID_TYPE:
        return "VOID_TYPE"
    default:
        return "UNKNOWN"
    }
}

func (t Token) String() string {
    return fmt.Sprintf("%s(%q)", t.Type, t.Lexeme)
}