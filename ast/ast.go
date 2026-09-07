package ast

import (
	"github.com/ArthurBryan166/mango/token"
)

type Node interface {
    String() string
}

type NumberLiteral struct {
    Value string
}

func (n NumberLiteral) String() string {
    return n.Value
}

type StringLiteral struct {
    Value string
}

func (s StringLiteral) String() string {
    return s.Value
}

type BooleanLiteral struct {
    Value bool
}

func (b BooleanLiteral) String() string {
    if b.Value {
        return "true"
    }

    return "false"
}

type VariableExpression struct {
    Name string
}

func (v VariableExpression) String() string {
    return v.Name
}

type BinaryExpression struct {
    Left     Node
    Operator token.Token
    Right    Node
}

func (b BinaryExpression) String() string {
    return "(" + b.Left.String() + " " +
        b.Operator.Lexeme + " " +
        b.Right.String() + ")"
}

type UnaryExpression struct {
    Operator token.Token
    Right    Node
}

func (u UnaryExpression) String() string {
    return "(" + u.Operator.Lexeme + u.Right.String() + ")"
}

type VariableDeclaration struct {
    Name        string
    Initializer Node
}

func (v VariableDeclaration) String() string {
    return v.Name + " <- " + v.Initializer.String()
}

type Assignment struct {
	Name  string
	Value Node
}

func (a Assignment) String() string {
	return a.Name + " <- " + a.Value.String()
}

type IfStatement struct {
    Condition Node
    Body      []Node
    ElseBody []Node
}

func (i IfStatement) String() string {
    result := "if " + i.Condition.String() + " {\n"

    for _, statement := range i.Body {
        result += "    " + statement.String() + "\n"
    }

    result += "}"

    if len(i.ElseBody) > 0 {
        result += " else {\n"

        for _, statement := range i.ElseBody {
            result += "    " + statement.String() + "\n"
        }

        result += "}"
    }

    return result
}

type ForStatement struct {
    Initializer Node
    Condition   Node
    Increment   Node
    Body        []Node
}

func (f ForStatement) String() string {
    result := "for "

    if f.Initializer != nil {
        result += f.Initializer.String()
    }

    if f.Condition != nil {
        if f.Initializer != nil {
            result += "; "
        }

        result += f.Condition.String()
    }

    if f.Increment != nil {
        result += "; "
        result += f.Increment.String()
    }

    result += " {\n"

    for _, statement := range f.Body {
        result += "    " + statement.String() + "\n"
    }

    result += "}"

    return result
}

type PrintStatement struct {
    Expressions []Node
}

func (p PrintStatement) String() string {
    result := "mangout("

    for i, expression := range p.Expressions {
        if i > 0 {
            result += ", "
        }

        result += expression.String()
    }

    result += ")"

    return result
}

type InputStatement struct {
    Variable VariableExpression
}

func (i InputStatement) String() string {
    return "mangoin(" + i.Variable.String() + ")"
}

type Parameter struct {
	Name string
	Type string
}

type FunctionDeclaration struct {
	Name       string
	Parameters []Parameter
	ReturnType string
	Body       []Node
}

func (f FunctionDeclaration) String() string {
	result := "mango " + f.Name + "("

	for i, parameter := range f.Parameters {
		if i > 0 {
			result += ", "
		}

		result += parameter.Name + " " + parameter.Type
	}

	result += ") " + f.ReturnType + " {\n"

	for _, statement := range f.Body {
		result += "    " + statement.String() + "\n"
	}

	result += "}"

	return result
}

type ReturnStatement struct {
	Value Node
}

func (r ReturnStatement) String() string {
	if r.Value == nil {
		return "return"
	}

	return "return " + r.Value.String()
}

type CallExpression struct {
	Callee    Node
	Arguments []Node
}

func (c CallExpression) String() string {
	result := c.Callee.String() + "("

	for i, argument := range c.Arguments {
		if i > 0 {
			result += ", "
		}

		result += argument.String()
	}

	result += ")"

	return result
}

type ExpressionStatement struct {
    Expression Node
}

func (e ExpressionStatement) String() string {
    return e.Expression.String()
}