package parser

import (
	"fmt"

	"github.com/ArthurBryan166/mango/ast"
	"github.com/ArthurBryan166/mango/token"
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

    if p.check(token.LEFT_PAREN) {
        p.advance()

        expr, err := p.expression()
        if err != nil {
            return nil, err
        }

        _, err = p.consume(
            token.RIGHT_PAREN,
            "esperado ')' depois da expressão",
        )

        if err != nil {
            return nil, err
        }

        return expr, nil
    }

    return nil, fmt.Errorf(
        "expressão esperada na linha %d",
        p.peek().Line,
    )
}

func (p *Parser) expression() (ast.Node, error) {
	return p.equality()
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

    if p.check(token.IF) {
        return p.ifStatement()
    }

    return nil, fmt.Errorf(
        "comando inesperado '%s' na linha %d",
        p.peek().Lexeme,
        p.peek().Line,
    )
}

func (p *Parser) unary() (ast.Node, error) {
    if p.check(token.MINUS) || p.check(token.PLUS) {
        operator := p.advance()

        right, err := p.unary()

        if err != nil {
            return nil, err
        }

        return ast.UnaryExpression{
            Operator: operator,
            Right:    right,
        }, nil
    }

    return p.primary()
}

func (p *Parser) factor() (ast.Node, error) {
	left, err := p.unary()

	if err != nil {
		return nil, err
	}

	for p.check(token.MULTIPLY) || p.check(token.DIVIDE) {
		operator := p.advance()

		right, err := p.unary()

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

func (p *Parser) term() (ast.Node, error) {
	left, err := p.factor()

	if err != nil {
		return nil, err
	}

	for p.check(token.PLUS) || p.check(token.MINUS) {
		operator := p.advance()

		right, err := p.factor()

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

func (p *Parser) equality() (ast.Node, error) {
	left, err := p.comparison()

	if err != nil {
		return nil, err
	}

	for p.check(token.EQUAL) {
		operator := p.advance()

		right, err := p.comparison()

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

func (p *Parser) comparison() (ast.Node, error) {
	left, err := p.term()

	if err != nil {
		return nil, err
	}

	for p.check(token.GREATER) ||
		p.check(token.GREATER_EQUAL) ||
		p.check(token.LESS) ||
		p.check(token.LESS_EQUAL) {

		operator := p.advance()

		right, err := p.term()

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

func (p *Parser) ifStatement() (ast.Node, error) {
    p.advance() // consome 'if'

    condition, err := p.expression()
    if err != nil {
        return nil, err
    }

    _, err = p.consume(token.LEFT_BRACE, "esperado '{' depois da condição")

    if err != nil {
        return nil, err
    }

    var body []ast.Node

    for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
        statement, err := p.declaration()

        if err != nil {
            return nil, err
        }

        body = append(body, statement)
    }

    _, err = p.consume(token.RIGHT_BRACE, "esperado '}' depois do bloco")

    if err != nil {
        return nil, err
    }

    var elseBody []ast.Node

    if p.check(token.ELSE) {
        p.advance()

        // else if
        if p.check(token.IF) {
            nestedIf, err := p.ifStatement()

            if err != nil {
                return nil, err
            }

            elseBody = append(elseBody, nestedIf)
        } else {
            // else normal
            _, err = p.consume(
                token.LEFT_BRACE,
                "esperado '{' depois de 'else'",
            )

            if err != nil {
                return nil, err
            }

            for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
                statement, err := p.declaration()

                if err != nil {
                    return nil, err
                }

                elseBody = append(elseBody, statement)
            }

            _, err = p.consume(
                token.RIGHT_BRACE,
                "esperado '}' depois do bloco 'else'",
            )

            if err != nil {
                return nil, err
            }
        }
    }

    return ast.IfStatement{
        Condition: condition,
        Body:      body,
        ElseBody:  elseBody,
    }, nil
}