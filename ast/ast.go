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