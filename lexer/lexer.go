package lexer

import "github.com/ArthurBryan166/mango.git/token"

type Lexer struct{
	source string
	start int
	current int
	line int
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
func (l *Lexer) makeToken(tokenType token.TokenType) token.Token {
    return token.Token{
        Type:   tokenType,
        Lexeme: l.source[l.start:l.current],
        Line:   l.line,
    }
}

// lê um token e o classifica
func (l *Lexer) scanToken() token.Token {
	l.skipWhitespace()

	l.start = l.current

    if l.isAtEnd() {
        return token.Token{
            Type: token.EOF,
            Line: l.line,
        }
    }

    char := l.advance()

	// tokens de um caractere
    switch char {
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
    }

	// token numeral
	if isDigit(char) {
        return l.scanNumber()
    }

    return token.Token{}
}

// lê um número e o classifica
func (l *Lexer) scanNumber() token.Token {
    for isDigit(l.peek()) {
        l.advance()
    }

    return l.makeToken(token.NUMBER)
}
