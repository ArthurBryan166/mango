package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ArthurBryan166/mango/ast"
	"github.com/ArthurBryan166/mango/token"
)

type Interpreter struct {
    environment     *Environment
    reader          *bufio.Reader
	writer 		    *bufio.Writer
    insideFunction 	bool
}

func New() *Interpreter {
    return &Interpreter{
        environment:    NewEnvironment(),
        reader:         bufio.NewReader(os.Stdin),
        writer:         bufio.NewWriter(os.Stdout),
        insideFunction: false,
    }
}

func (i *Interpreter) Run(program []ast.Node) error {
    for _, node := range program {
        function, ok := node.(ast.FunctionDeclaration)

        if !ok {
            return fmt.Errorf(
                "código executável não pode existir fora de uma função",
            )
        }

        _, err := i.evaluate(function)
        if err != nil {
            return err
        }
    }

    mainValue, exists := i.environment.Get("main")

    if !exists {
        return fmt.Errorf("função 'main' não foi declarada")
    }

    if mainValue.Type != FUNCTION_VALUE {
        return fmt.Errorf("'main' não é uma função")
    }

    mainFunction := mainValue.Value.(Function)

    if mainFunction.Declaration.ReturnType != "void" {
        return fmt.Errorf("função 'main' deve retornar void")
    }

    if len(mainFunction.Declaration.Parameters) != 0 {
        return fmt.Errorf("função 'main' não pode receber parâmetros")
    }

    _, err := i.callFunction(mainFunction, nil)

    return err
}

func (i *Interpreter) evaluate(node ast.Node) (Value, error) {
	switch n := node.(type) {
	case ast.NumberLiteral:
		value, err := strconv.ParseFloat(n.Value, 64)

		if err != nil {
			return Value{}, fmt.Errorf("número inválido: %s", n.Value)
		}

		return Value{
			Type:  NUMBER_VALUE,
			Value: value,
		}, nil

	case ast.StringLiteral:
		return Value{
			Type:  STRING_VALUE,
			Value: n.Value,
		}, nil

	case ast.BooleanLiteral:
		return Value{
			Type:  BOOLEAN_VALUE,
			Value: n.Value,
		}, nil

	case ast.VariableExpression:
		value, exists := i.environment.Get(n.Name)

		if !exists {
			return Value{}, fmt.Errorf("variável '%s' não foi declarada", n.Name)
		}

		return value, nil

	case ast.VariableDeclaration:
		value, err := i.evaluate(n.Initializer)

		if err != nil {
			return Value{}, err
		}

		err = i.environment.Define(n.Name, value)

		if err != nil {
			return Value{}, err
		}

		return value, nil

	case ast.Assignment:
		value, err := i.evaluate(n.Value)

		if err != nil {
			return Value{}, err
		}

		err = i.environment.Assign(n.Name, value)

		if err != nil {
			return Value{}, err
		}

		return value, nil

	case ast.BinaryExpression:
		left, err := i.evaluate(n.Left)
		if err != nil {
			return Value{}, err
		}

		switch n.Operator.Type {

		case token.AND:
			if left.Type != BOOLEAN_VALUE {
				return Value{}, fmt.Errorf(
					"operador 'and' requer dois booleanos",
				)
			}

			if !left.Value.(bool) {
				return Value{
					Type:  BOOLEAN_VALUE,
					Value: false,
				}, nil
			}

			right, err := i.evaluate(n.Right)
			if err != nil {
				return Value{}, err
			}

			if right.Type != BOOLEAN_VALUE {
				return Value{}, fmt.Errorf(
					"operador 'and' requer dois booleanos",
				)
			}

			return Value{
				Type:  BOOLEAN_VALUE,
				Value: right.Value.(bool),
			}, nil

		case token.OR:
			if left.Type != BOOLEAN_VALUE {
				return Value{}, fmt.Errorf(
					"operador 'or' requer dois booleanos",
				)
			}

			if left.Value.(bool) {
				return Value{
					Type:  BOOLEAN_VALUE,
					Value: true,
				}, nil
			}

			right, err := i.evaluate(n.Right)
			if err != nil {
				return Value{}, err
			}

			if right.Type != BOOLEAN_VALUE {
				return Value{}, fmt.Errorf(
					"operador 'or' requer dois booleanos",
				)
			}

			return Value{
				Type:  BOOLEAN_VALUE,
				Value: right.Value.(bool),
			}, nil

		default:
			right, err := i.evaluate(n.Right)
			if err != nil {
				return Value{}, err
			}

			switch n.Operator.Type {

			case token.PLUS:
				if left.Type == NUMBER_VALUE && right.Type == NUMBER_VALUE {
					return Value{
						Type:  NUMBER_VALUE,
						Value: left.Value.(float64) + right.Value.(float64),
					}, nil
				}

				if left.Type == STRING_VALUE && right.Type == STRING_VALUE {
					return Value{
						Type:  STRING_VALUE,
						Value: left.Value.(string) + right.Value.(string),
					}, nil
				}

				return Value{}, fmt.Errorf(
					"operador '+' requer dois números ou duas strings",
				)

			case token.MINUS, token.MULTIPLY, token.DIVIDE:
				if left.Type != NUMBER_VALUE || right.Type != NUMBER_VALUE {
					return Value{}, fmt.Errorf(
						"operador '%s' requer dois números",
						n.Operator.Lexeme,
					)
				}

				leftValue := left.Value.(float64)
				rightValue := right.Value.(float64)

				switch n.Operator.Type {
				case token.MINUS:
					return Value{
						Type:  NUMBER_VALUE,
						Value: leftValue - rightValue,
					}, nil

				case token.MULTIPLY:
					return Value{
						Type:  NUMBER_VALUE,
						Value: leftValue * rightValue,
					}, nil

				case token.DIVIDE:
					if rightValue == 0 {
						return Value{}, fmt.Errorf("divisão por zero")
					}

					return Value{
						Type:  NUMBER_VALUE,
						Value: leftValue / rightValue,
					}, nil
				}

			case token.LESS, token.LESS_EQUAL,
				token.GREATER, token.GREATER_EQUAL:

				if left.Type != NUMBER_VALUE || right.Type != NUMBER_VALUE {
					return Value{}, fmt.Errorf(
						"operador '%s' requer dois números",
						n.Operator.Lexeme,
					)
				}

				leftValue := left.Value.(float64)
				rightValue := right.Value.(float64)

				switch n.Operator.Type {
				case token.LESS:
					return Value{
						Type:  BOOLEAN_VALUE,
						Value: leftValue < rightValue,
					}, nil

				case token.LESS_EQUAL:
					return Value{
						Type:  BOOLEAN_VALUE,
						Value: leftValue <= rightValue,
					}, nil

				case token.GREATER:
					return Value{
						Type:  BOOLEAN_VALUE,
						Value: leftValue > rightValue,
					}, nil

				case token.GREATER_EQUAL:
					return Value{
						Type:  BOOLEAN_VALUE,
						Value: leftValue >= rightValue,
					}, nil
				}

			case token.EQUAL:
				if left.Type != right.Type {
					return Value{
						Type:  BOOLEAN_VALUE,
						Value: false,
					}, nil
				}

				switch left.Type {
				case NUMBER_VALUE:
					return Value{
						Type:  BOOLEAN_VALUE,
						Value: left.Value.(float64) == right.Value.(float64),
					}, nil

				case STRING_VALUE:
					return Value{
						Type:  BOOLEAN_VALUE,
						Value: left.Value.(string) == right.Value.(string),
					}, nil

				case BOOLEAN_VALUE:
					return Value{
						Type:  BOOLEAN_VALUE,
						Value: left.Value.(bool) == right.Value.(bool),
					}, nil

				case NIL_VALUE:
					return Value{
						Type:  BOOLEAN_VALUE,
						Value: true,
					}, nil
				}

			default:
				return Value{}, fmt.Errorf(
					"operador '%s' não suportado",
					n.Operator.Lexeme,
				)
			}

			return Value{}, fmt.Errorf(
				"operação '%s' inválida para os tipos %s e %s",
				n.Operator.Lexeme,
				left.Type,
				right.Type,
			)
		}

	case ast.UnaryExpression:
		right, err := i.evaluate(n.Right)

		if err != nil {
			return Value{}, err
		}

		switch n.Operator.Type {
		case token.MINUS:
			if right.Type != NUMBER_VALUE {
				return Value{}, fmt.Errorf(
					"operador '-' requer um número",
				)
			}

			return Value{
				Type:  NUMBER_VALUE,
				Value: -right.Value.(float64),
			}, nil

		case token.PLUS:
			if right.Type != NUMBER_VALUE {
				return Value{}, fmt.Errorf(
					"operador '+' requer um número",
				)
			}

			return Value{
				Type:  NUMBER_VALUE,
				Value: right.Value.(float64),
			}, nil

		case token.NOT:
			if right.Type != BOOLEAN_VALUE {
				return Value{}, fmt.Errorf(
					"operador 'not' requer um booleano",
				)
			}

			return Value{
				Type:  BOOLEAN_VALUE,
				Value: !right.Value.(bool),
			}, nil

		default:
			return Value{}, fmt.Errorf(
				"operador unário '%s' não suportado",
				n.Operator.Lexeme,
			)
		}

	case ast.IfStatement:
		condition, err := i.evaluate(n.Condition)
		if err != nil {
			return Value{}, err
		}

		if condition.Type != BOOLEAN_VALUE {
			return Value{}, fmt.Errorf(
				"condição do 'if' deve ser um booleano",
			)
		}

		if condition.Value.(bool) {
			environment := NewEnclosedEnvironment(i.environment)

			err := i.executeBlock(n.Body, environment)
			if err != nil {
				return Value{}, err
			}
		} else {
			environment := NewEnclosedEnvironment(i.environment)

			err := i.executeBlock(n.ElseBody, environment)
			if err != nil {
				return Value{}, err
			}
		}

		return Value{
			Type:  NIL_VALUE,
			Value: nil,
		}, nil

	case ast.PrintStatement:
		for _, expression := range n.Expressions {
			value, err := i.evaluate(expression)

			if err != nil {
				return Value{}, err
			}

			fmt.Fprint(i.writer, value.Value)
		}

		fmt.Fprintln(i.writer)

		return Value{
			Type:  NIL_VALUE,
			Value: nil,
		}, nil

	case ast.InputStatement:
		input, err := i.readInput()

		if err != nil {
			return Value{}, fmt.Errorf(
				"erro ao ler entrada: %v",
				err,
			)
		}

		currentValue, exists := i.environment.Get(n.Variable.Name)

		if !exists {
			return Value{}, fmt.Errorf(
				"variável '%s' não foi declarada",
				n.Variable.Name,
			)
		}

		var value Value

		switch currentValue.Type {
		case NUMBER_VALUE:
			number, err := strconv.ParseFloat(input, 64)

			if err != nil {
				return Value{}, fmt.Errorf(
					"entrada inválida: '%s' não é um número",
					input,
				)
			}

			value = Value{
				Type:  NUMBER_VALUE,
				Value: number,
			}

		case STRING_VALUE:
			value = Value{
				Type:  STRING_VALUE,
				Value: input,
			}

		case BOOLEAN_VALUE:
			switch strings.ToLower(input) {
			case "true":
				value = Value{
					Type:  BOOLEAN_VALUE,
					Value: true,
				}

			case "false":
				value = Value{
					Type:  BOOLEAN_VALUE,
					Value: false,
				}

			default:
				return Value{}, fmt.Errorf(
					"entrada inválida: '%s' não é um booleano",
					input,
				)
			}

		default:
			return Value{}, fmt.Errorf(
				"tipo de variável não pode ser usado com 'mangoin'",
			)
		}

		if err := i.environment.Assign(n.Variable.Name, value); err != nil {
			return Value{}, err
		}

		return Value{
			Type:  NIL_VALUE,
			Value: nil,
		}, nil

	case ast.ForStatement:
		previous := i.environment

		loopEnvironment := NewEnclosedEnvironment(previous)
		i.environment = loopEnvironment

		defer func() {
			i.environment = previous
		}()

		if n.Initializer != nil {
			_, err := i.evaluate(n.Initializer)
			if err != nil {
				return Value{}, err
			}
		}

		for {
			if n.Condition != nil {
				condition, err := i.evaluate(n.Condition)
				if err != nil {
					return Value{}, err
				}

				if condition.Type != BOOLEAN_VALUE {
					return Value{}, fmt.Errorf(
						"condição do 'for' deve ser um booleano",
					)
				}

				if !condition.Value.(bool) {
					break
				}
			}

			bodyEnvironment := NewEnclosedEnvironment(loopEnvironment)

			if err := i.executeBlock(
				n.Body,
				bodyEnvironment,
			); err != nil {
				return Value{}, err
			}

			if n.Increment != nil {
				_, err := i.evaluate(n.Increment)
				if err != nil {
					return Value{}, err
				}
			}
		}

		return Value{Type: NIL_VALUE, Value: nil}, nil
		
	case ast.ReturnStatement:
		if !i.insideFunction {
			return Value{}, fmt.Errorf(
				"'return' só pode ser usado dentro de uma função",
			)
		}

		value, err := i.evaluate(n.Value)

		if err != nil {
			return Value{}, err
		}

		return Value{}, ReturnValue{
			Value: value,
		}

	case ast.FunctionDeclaration:
		err := i.environment.Define(
			n.Name,
			Value{
				Type: FUNCTION_VALUE,
				Value: Function{
					Declaration: n,
					Closure:     i.environment,
				},
			},
		)

		if err != nil {
			return Value{}, err
		}

		return Value{}, nil

	case ast.CallExpression:
		callee, err := i.evaluate(n.Callee)

		if err != nil {
			return Value{}, err
		}

		if callee.Type != FUNCTION_VALUE {
			return Value{}, fmt.Errorf(
				"não é possível chamar %s",
				callee.Type,
			)
		}

		function := callee.Value.(Function)

		arguments := make([]Value, len(n.Arguments))

		for j, argument := range n.Arguments {
			value, err := i.evaluate(argument)

			if err != nil {
				return Value{}, err
			}

			arguments[j] = value
		}

		return i.callFunction(function, arguments)

	case ast.ExpressionStatement:
		_, err := i.evaluate(n.Expression)

		if err != nil {
			return Value{}, err
		}

		return Value{
			Type:  NIL_VALUE,
			Value: nil,
		}, nil

	default:
		return Value{}, fmt.Errorf("nó não suportado: %T", node)
	}
}

func (i *Interpreter) Evaluate(node ast.Node) (Value, error) {
	return i.evaluate(node)
}

func (i *Interpreter) Define(name string, value Value) error {
    return i.environment.Define(name, value)
}

func (i *Interpreter) readInput() (string, error) {
	input, err := i.reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	return strings.TrimRight(input, "\r\n"), nil
}

func (i *Interpreter) executeBlock(
    statements []ast.Node,
    environment *Environment,
) error {
    previous := i.environment
    i.environment = environment

    defer func() {
        i.environment = previous
    }()

    for _, statement := range statements {
        _, err := i.evaluate(statement)

        if err != nil {
            return err
        }
    }

    return nil
}

func (i *Interpreter) callFunction(
	function Function,
	arguments []Value,
) (Value, error) {

	if len(arguments) != len(function.Declaration.Parameters) {
		return Value{}, fmt.Errorf(
			"função %s esperava %d argumentos, recebeu %d",
			function.Declaration.Name,
			len(function.Declaration.Parameters),
			len(arguments),
		)
	}

	for j, argument := range arguments {
		expectedType := function.Declaration.Parameters[j].Type

		if !matchesType(argument, expectedType) {
			return Value{}, fmt.Errorf(
				"argumento %d da função '%s' deveria ser %s, mas recebeu %s",
				j+1,
				function.Declaration.Name,
				expectedType,
				argument.Type,
			)
		}
	}

	functionEnvironment := NewEnclosedEnvironment(function.Closure)

	for j, parameter := range function.Declaration.Parameters {
		err := functionEnvironment.Define(
			parameter.Name,
			arguments[j],
		)

		if err != nil {
			return Value{}, err
		}
	}

	previousFunctionContext := i.insideFunction
	i.insideFunction = true

	err := i.executeBlock(
		function.Declaration.Body,
		functionEnvironment,
	)

	i.insideFunction = previousFunctionContext

	if err != nil {
		if returnValue, ok := err.(ReturnValue); ok {

			if function.Declaration.ReturnType == "void" {
				return Value{}, fmt.Errorf(
					"função '%s' é void e não pode usar 'return'",
					function.Declaration.Name,
				)
			}

			if !matchesType(
				returnValue.Value,
				function.Declaration.ReturnType,
			) {
				return Value{}, fmt.Errorf(
					"função '%s' deveria retornar %s, mas retornou %s",
					function.Declaration.Name,
					function.Declaration.ReturnType,
					returnValue.Value.Type,
				)
			}

			return returnValue.Value, nil
		}

		return Value{}, err
	}

	if function.Declaration.ReturnType != "void" {
		return Value{}, fmt.Errorf(
			"função '%s' deveria retornar %s, mas não retornou nenhum valor",
			function.Declaration.Name,
			function.Declaration.ReturnType,
		)
	}

	return Value{
		Type:  NIL_VALUE,
		Value: nil,
	}, nil
}

func matchesType(value Value, expected string) bool {
	switch expected {
	case "number":
		return value.Type == NUMBER_VALUE

	case "string":
		return value.Type == STRING_VALUE

	case "bool":
		return value.Type == BOOLEAN_VALUE

	default:
		return false
	}
}