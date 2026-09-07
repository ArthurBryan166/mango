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

func (p *Parser) peekNext() token.Token {
	if p.current+1 >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}

	return p.tokens[p.current+1]
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

    if p.check(token.STRING) {
        tok := p.advance()

        return ast.StringLiteral{
            Value: tok.Lexeme[1 : len(tok.Lexeme)-1],
        }, nil
    }

    if p.check(token.TRUE) {
        p.advance()

        return ast.BooleanLiteral{
            Value: true,
        }, nil
    }

    if p.check(token.FALSE) {
        p.advance()

        return ast.BooleanLiteral{
            Value: false,
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
	return p.or()
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

    if p.check(token.MANGOUT) {
        return p.printStatement()
    }

    if p.check(token.MANGOIN) {
        return p.inputStatement()
    }

    if p.check(token.IDENTIFIER) && p.peekNext().Type == token.ASSIGN {
        return p.assignmentStatement()
    }

    if p.check(token.IDENTIFIER) {
        expression, err := p.expression()

        if err != nil {
            return nil, err
        }

        return ast.ExpressionStatement{
            Expression: expression,
        }, nil
    }

    if p.check(token.FOR) {
        return p.forStatement()
    }

    if p.check(token.MANGO) {
        return p.functionDeclaration()
    }

    if p.check(token.RETURN) {
        return p.returnStatement()
    }

    return nil, fmt.Errorf(
        "comando inesperado '%s' na linha %d",
        p.peek().Lexeme,
        p.peek().Line,
    )
}

func (p *Parser) unary() (ast.Node, error) {
    if p.check(token.MINUS) || p.check(token.PLUS) || p.check(token.NOT) {
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

    return p.callExpression()
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

func (p *Parser) and() (ast.Node, error) {
    left, err := p.equality()
    if err != nil {
        return nil, err
    }

    for p.check(token.AND) {
        operator := p.advance()

        right, err := p.equality()
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

func (p *Parser) or() (ast.Node, error) {
    left, err := p.and()
    if err != nil {
        return nil, err
    }

    for p.check(token.OR) {
        operator := p.advance()

        right, err := p.and()
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

func (p *Parser) forStatement() (ast.Node, error) {
	p.advance()

	var initializer ast.Node
	var condition ast.Node
	var increment ast.Node
	var err error

	// for { ... }
	if p.check(token.LEFT_BRACE) {
		p.advance()

		body, err := p.parseBlock()
		if err != nil {
			return nil, err
		}

		return ast.ForStatement{
			Initializer: nil,
			Condition:   nil,
			Increment:   nil,
			Body:        body,
		}, nil
	}

	// for mng i <- 0; ...
	if p.check(token.MNG) {
		initializer, err = p.variableDeclaration()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(
			token.SEMICOLON,
			"esperado ';' depois da inicialização do 'for'",
		)
		if err != nil {
			return nil, err
		}
	}

	// for mng i <- 0; condição; incremento { ... }
	if initializer != nil {
		if !p.check(token.SEMICOLON) {
			condition, err = p.expression()
			if err != nil {
				return nil, err
			}
		}

		_, err = p.consume(
			token.SEMICOLON,
			"esperado ';' antes do incremento",
		)
		if err != nil {
			return nil, err
		}

		if !p.check(token.LEFT_BRACE) {
			if !p.check(token.IDENTIFIER) ||
				p.peekNext().Type != token.ASSIGN {
				return nil, fmt.Errorf(
					"esperada uma atribuição como incremento na linha %d",
					p.peek().Line,
				)
			}

			increment, err = p.assignmentStatement()
			if err != nil {
				return nil, err
			}
		}
	} else {
		// for condição { ... }
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(
		token.LEFT_BRACE,
		"esperado '{' depois do 'for'",
	)
	if err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return ast.ForStatement{
		Initializer: initializer,
		Condition:   condition,
		Increment:   increment,
		Body:        body,
	}, nil
}

func (p *Parser) parseBlock() ([]ast.Node, error) {
    var body []ast.Node

    for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
        statement, err := p.declaration()

        if err != nil {
            return nil, err
        }

        body = append(body, statement)
    }

    _, err := p.consume(
        token.RIGHT_BRACE,
        "esperado '}' depois do bloco",
    )

    if err != nil {
        return nil, err
    }

    return body, nil
}

func (p *Parser) printStatement() (ast.Node, error) {
    p.advance()

    _, err := p.consume(
        token.LEFT_PAREN,
        "esperado '(' depois de 'mangout'",
    )

    if err != nil {
        return nil, err
    }

    expressions := []ast.Node{}

    expression, err := p.expression()

    if err != nil {
        return nil, err
    }

    expressions = append(expressions, expression)

    for p.check(token.COMMA) {
        p.advance()

        expression, err := p.expression()

        if err != nil {
            return nil, err
        }

        expressions = append(expressions, expression)
    }

    _, err = p.consume(
        token.RIGHT_PAREN,
        "esperado ')' depois das expressões",
    )

    if err != nil {
        return nil, err
    }

    return ast.PrintStatement{
        Expressions: expressions,
    }, nil
}

func (p *Parser) inputStatement() (ast.Node, error) {
    p.advance()

    _, err := p.consume(
        token.LEFT_PAREN,
        "esperado '(' depois de 'mangoin'",
    )
    if err != nil {
        return nil, err
    }

    name, err := p.consume(
        token.IDENTIFIER,
        "esperado uma variável dentro de 'mangoin'",
    )
    if err != nil {
        return nil, err
    }

    _, err = p.consume(
        token.RIGHT_PAREN,
        "esperado ')' depois da variável",
    )
    if err != nil {
        return nil, err
    }

    return ast.InputStatement{
        Variable: ast.VariableExpression{
            Name: name.Lexeme,
        },
    }, nil
}

func (p *Parser) assignmentStatement() (ast.Node, error) {
	name := p.advance()

	_, err := p.consume(
		token.ASSIGN,
		"esperado '<-' depois do identificador",
	)
	if err != nil {
		return nil, err
	}

	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	return ast.Assignment{
		Name:  name.Lexeme,
		Value: value,
	}, nil
}

func (p *Parser) consumeParameterType() (string, error) {
	switch p.peek().Type {
	case token.NUMBER_TYPE:
		p.advance()
		return "number", nil

	case token.STRING_TYPE:
		p.advance()
		return "string", nil

	case token.BOOL_TYPE:
		p.advance()
		return "bool", nil

	default:
		return "", fmt.Errorf(
			"esperado tipo 'number', 'string' ou 'bool' na linha %d",
			p.peek().Line,
		)
	}
}

func (p *Parser) consumeReturnType() (string, error) {
	switch p.peek().Type {
	case token.NUMBER_TYPE:
		p.advance()
		return "number", nil

	case token.STRING_TYPE:
		p.advance()
		return "string", nil

	case token.BOOL_TYPE:
		p.advance()
		return "bool", nil

	case token.VOID_TYPE:
		p.advance()
		return "void", nil

	default:
		return "", fmt.Errorf("tipo de retorno esperado")
	}
}

func (p *Parser) functionDeclaration() (ast.Node, error) {
	_, err := p.consume(
		token.MANGO,
		"esperado 'mango'",
	)
	if err != nil {
		return nil, err
	}

	name, err := p.consume(
		token.IDENTIFIER,
		"esperado nome da função depois de 'mango'",
	)
	if err != nil {
		return nil, err
	}

	_, err = p.consume(
		token.LEFT_PAREN,
		"esperado '(' depois do nome da função",
	)
	if err != nil {
		return nil, err
	}

	var parameters []ast.Parameter

	if !p.check(token.RIGHT_PAREN) {
		for {
			parameterName, err := p.consume(
				token.IDENTIFIER,
				"esperado nome do parâmetro",
			)
			if err != nil {
				return nil, err
			}

			parameterType, err := p.consumeParameterType()
			if err != nil {
				return nil, err
			}

			parameters = append(parameters, ast.Parameter{
				Name: parameterName.Lexeme,
				Type: parameterType,
			})

			if !p.check(token.COMMA) {
				break
			}

			p.advance()
		}
	}

	_, err = p.consume(
		token.RIGHT_PAREN,
		"esperado ')' depois dos parâmetros",
	)
	if err != nil {
		return nil, err
	}

	returnType, err := p.consumeReturnType()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(
		token.LEFT_BRACE,
		"esperado '{' depois do tipo de retorno",
	)
	if err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return ast.FunctionDeclaration{
		Name:       name.Lexeme,
		Parameters: parameters,
		ReturnType: returnType,
		Body:       body,
	}, nil
}

func (p *Parser) returnStatement() (ast.Node, error) {
	p.advance()

	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	return ast.ReturnStatement{
		Value: value,
	}, nil
}

func (p *Parser) callExpression() (ast.Node, error) {
	callee, err := p.primary()
	if err != nil {
		return nil, err
	}

	for p.check(token.LEFT_PAREN) {
		p.advance()

		var arguments []ast.Node

		if !p.check(token.RIGHT_PAREN) {
			for {
				argument, err := p.expression()
				if err != nil {
					return nil, err
				}

				arguments = append(arguments, argument)

				if !p.check(token.COMMA) {
					break
				}

				p.advance()
			}
		}

		_, err = p.consume(
			token.RIGHT_PAREN,
			"esperado ')' depois dos argumentos",
		)
		if err != nil {
			return nil, err
		}

		callee = ast.CallExpression{
			Callee:    callee,
			Arguments: arguments,
		}
	}

	return callee, nil
}