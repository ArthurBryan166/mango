package parser

import (
	"github.com/ArthurBryan166/mango/token"
	"github.com/ArthurBryan166/mango/ast"
	"fmt"
)

type Parser struct {
    tokens  []token.Token
    current int
}

func (p *Parser) isAtEnd() bool {
    return p.peek().Type == token.EOF
}

func (p *Parser) peek() token.Token {
    return p.tokens[p.current]
}

func (p *Parser) previous() token.Token {
    return p.tokens[p.current-1]
}

func (p *Parser) advance() token.Token {
    if !p.isAtEnd() {
        p.current++
    }

    return p.previous()
}

func New(tokens []token.Token) *Parser {
    return &Parser{
        tokens:  tokens,
        current: 0,
    }
}

func (p *Parser) check(t token.TokenType) bool {
    if p.isAtEnd() {
        return false
    }

    return p.peek().Type == t
}

func (p *Parser) consume(t token.TokenType, message string) (token.Token, error) {
    if p.check(t) {
        return p.advance(), nil
    }

    return token.Token{}, fmt.Errorf("%s na linha %d", message, p.peek().Line)
}

func (p *Parser) primary() (ast.Node, error) {
    if p.check(token.NUMBER) {
        tok := p.advance()

        return ast.NumberLiteral{
            Value: tok.Lexeme,
        }, nil
    }

    if p.check(token.IDENTIFIER) {
        tok := p.advance()

        return ast.VariableExpression{
            Name: tok.Lexeme,
        }, nil
    }

    return nil, fmt.Errorf(
        "expressão esperada na linha %d",
        p.peek().Line,
    )
}

func (p *Parser) expression() (ast.Node, error) {
    return p.term()
}

func (p *Parser) variableDeclaration() (ast.Node, error) {
    _, err := p.consume(
        token.MNG,
        "esperado 'mng'",
    )

    if err != nil {
        return nil, err
    }

    name, err := p.consume(
        token.IDENTIFIER,
        "esperado um identificador depois de 'mng'",
    )

    if err != nil {
        return nil, err
    }

    _, err = p.consume(
        token.ASSIGN,
        "esperado '<-' depois do identificador",
    )

    if err != nil {
        return nil, err
    }

    initializer, err := p.expression()

    if err != nil {
        return nil, err
    }

    return ast.VariableDeclaration{
        Name:        name.Lexeme,
        Initializer: initializer,
    }, nil
}

func (p *Parser) Parse() ([]ast.Node, error) {
    var nodes []ast.Node

    for !p.isAtEnd() {
        node, err := p.declaration()

        if err != nil {
            return nil, err
        }

        nodes = append(nodes, node)
    }

    return nodes, nil
}

func (p *Parser) declaration() (ast.Node, error) {
    if p.check(token.MNG) {
        return p.variableDeclaration()
    }

    return nil, fmt.Errorf("comando inesperado '%s' na linha %d", p.peek().Lexeme, p.peek().Line)
}

func (p *Parser) term() (ast.Node, error) {
    left, err := p.primary()

    if err != nil {
        return nil, err
    }

    for p.check(token.PLUS) || p.check(token.MINUS) {
        operator := p.advance()

        right, err := p.primary()

        if err != nil {
            return nil, err
        }

        left = ast.BinaryExpression{
            Left:     left,
            Operator: operator,
            Right:    right,
        }
    }

    return left, nil
}