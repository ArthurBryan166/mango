package lexer

import (
	"github.com/ArthurBryan166/mango.git/token"
	"fmt"
)

type Lexer struct{
	source string
	start int
	current int
	line int
}

var keywords = map[string]token.TokenType{
    "mng":		token.MNG,
    "if":		token.IF,
    "else":		token.ELSE,
    "for":		token.FOR,
    "mango":	token.MANGO,
    "return":	token.RETURN,
    "true":		token.TRUE,
    "false":	token.FALSE,
    "nil":    	token.NIL,
}

// gera um novo lexer
func New(source string) *Lexer {
    return &Lexer{
        source:  source,
        current: 0,
        line:    1,
    }
}

// lê um caractere e avança uma posicao
func (l *Lexer) advance() byte {
    char := l.source[l.current]
    l.current++

    return char
}

// verifica se a linha chegou ao fim
func (l *Lexer) isAtEnd() bool {
    return l.current >= len(l.source)
}

// lê um caractere mas não avança
func (l *Lexer) peek() byte {
    if l.isAtEnd() {
        return 0
    }

    return l.source[l.current]
}

// pula espaços, tabs e \n
func (l *Lexer) skipWhitespace() {
    for {
        switch l.peek() {
        case ' ', '\t':
            l.advance()

        case '\n':
            l.advance()
            l.line++

        default:
            return
        }
    }
}

// verifica se é um dígito
func isDigit(char byte) bool {
    return char >= '0' && char <= '9'
}

// classifica um token
func (l *Lexer) makeToken(tokenType token.TokenType) (token.Token, error) {
    return token.Token{
        Type:   tokenType,
        Lexeme: l.source[l.start:l.current],
        Line:   l.line,
    }, nil
}

// verifica se o caractere faz parte do alfabeto
func isAlpha(char byte) bool {
    return (char >= 'a' && char <= 'z') ||
        (char >= 'A' && char <= 'Z')
}

// verifica se o caractere faz parte do alfabeto, ou é um número ou é _
func isAlphaNumeric(char byte) bool {
    return isAlpha(char) || isDigit(char) || char == '_'
}

func (l *Lexer) match(expected byte) bool {
    if l.isAtEnd() {
        return false
    }

    if l.source[l.current] != expected {
        return false
    }

    l.current++
    return true
}

// lê um token e o classifica
func (l *Lexer) scanToken() (token.Token, error) {
	l.skipWhitespace()

	l.start = l.current

    if l.isAtEnd() {
        return token.Token{
            Type: token.EOF,
            Line: l.line,
        }, nil
    }

    char := l.advance()

	// token numeral
	if isDigit(char) {
        return l.scanNumber()
    }

	// palavras-chave ou identificadores
	if isAlpha(char) {
		return l.scanIdentifier()
	}

	// tokens de um caractere
    switch char {
	case '"':
    	return l.scanString()

    case '+':
        return l.makeToken(token.PLUS)

    case '-':
        return l.makeToken(token.MINUS)

	case '*':
        return l.makeToken(token.MULTIPLY)

	case '/':
        return l.makeToken(token.DIVIDE)
	
	case '{':
        return l.makeToken(token.LEFT_BRACE)

    case '}':
        return l.makeToken(token.RIGHT_BRACE)

	case '(':
        return l.makeToken(token.LEFT_PAREN)

	case ')':
        return l.makeToken(token.RIGHT_PAREN)
	
	case '=':
    	return l.makeToken(token.EQUAL)

	case '<':
		if l.match('-') {
			return l.makeToken(token.ASSIGN)
		}

		if l.match('=') {
			return l.makeToken(token.LESS_EQUAL)
		}

		return l.makeToken(token.LESS)

	case '>':
		if l.match('=') {
			return l.makeToken(token.GREATER_EQUAL)
		}

		return l.makeToken(token.GREATER)
	}

    return token.Token{}, nil
}

// lê um número e o classifica
func (l *Lexer) scanNumber() (token.Token, error) {
    for isDigit(l.peek()) {
        l.advance()
    }

    return l.makeToken(token.NUMBER)
}

func (l *Lexer) scanIdentifier() (token.Token, error) {
    for isAlphaNumeric(l.peek()) {
        l.advance()
    }

    text := l.source[l.start:l.current]

    if tokenType, ok := keywords[text]; ok {
        return l.makeToken(tokenType)
    }

    return l.makeToken(token.IDENTIFIER)
}

func (l *Lexer) scanString() (token.Token, error) {
    for l.peek() != '"' && !l.isAtEnd() {
        l.advance()
    }

    if l.isAtEnd() {
        return token.Token{}, fmt.Errorf("string não fechada na linha %d", l.line,)
    }

    l.advance()

    return l.makeToken(token.STRING)
}

func (l *Lexer) ScanTokens() ([]token.Token, error) {
    var tokens []token.Token

    for !l.isAtEnd() {
        tok, err := l.scanToken()

        if err != nil {
            return nil, err
        }

        tokens = append(tokens, tok)
    }

    tokens = append(tokens, token.Token{
        Type: token.EOF,
        Line: l.line,
    })

    return tokens, nil
}