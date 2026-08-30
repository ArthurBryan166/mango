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

type VariableDeclaration struct {
    Name        string
    Initializer Node
}

func (v VariableDeclaration) String() string {
    return v.Name + " <- " + v.Initializer.String()
}

