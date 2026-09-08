package interpreter

import "github.com/ArthurBryan166/mango/ast"

type ValueType int

const (
    NUMBER_VALUE ValueType = iota
    STRING_VALUE
    BOOLEAN_VALUE
    NIL_VALUE
	FUNCTION_VALUE
)

type Value struct {
    Type  ValueType
    Value any
}

func (t ValueType) String() string {
	switch t {
	case NUMBER_VALUE:
		return "number"
	case STRING_VALUE:
		return "string"
	case BOOLEAN_VALUE:
		return "boolean"
	case NIL_VALUE:
		return "nil"
	case FUNCTION_VALUE:
		return "function"
	default:
		return "unknown"
	}
}

type Function struct {
    Declaration ast.FunctionDeclaration
    Closure     *Environment
}