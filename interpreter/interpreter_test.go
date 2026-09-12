package interpreter

import (
	"testing"

	"github.com/ArthurBryan166/mango/ast"
	"github.com/ArthurBryan166/mango/lexer"
	"github.com/ArthurBryan166/mango/parser"
	"github.com/ArthurBryan166/mango/token"
)

func TestFunctionCall(t *testing.T) {
	source := `
        mango dobro(x number) number {
            return x * 2
        }

        mango main() void {
            mangout(dobro(21))
        }
    `

	// =========================
	// LEXER
	// =========================

	l := lexer.New(source)

	tokens, err := l.ScanTokens()

	if err != nil {
		t.Fatalf("erro no lexer: %v", err)
	}

	if len(tokens) == 0 {
		t.Fatal("lexer não produziu nenhum token")
	}

	expected := []token.TokenType{
		token.MANGO,
		token.IDENTIFIER,
		token.LEFT_PAREN,
		token.IDENTIFIER,
		token.NUMBER_TYPE,
		token.RIGHT_PAREN,
		token.NUMBER_TYPE,
		token.LEFT_BRACE,
	}

	if len(tokens) < len(expected) {
		t.Fatalf(
			"esperados pelo menos %d tokens, recebidos %d",
			len(expected),
			len(tokens),
		)
	}

	for i, expectedType := range expected {
		if tokens[i].Type != expectedType {
			t.Fatalf(
				"token %d: esperado %v, recebido %v",
				i,
				expectedType,
				tokens[i].Type,
			)
		}
	}

	assertToken(t, tokens[0], token.MANGO, "mango")
	assertToken(t, tokens[1], token.IDENTIFIER, "dobro")
	assertToken(t, tokens[2], token.LEFT_PAREN, "(")
	assertToken(t, tokens[3], token.IDENTIFIER, "x")
	assertToken(t, tokens[4], token.NUMBER_TYPE, "number")

	// =========================
	// PARSER
	// =========================

	p := parser.New(tokens)

	program, err := p.Parse()

	if err != nil {
		t.Fatalf("erro no parser: %v", err)
	}

	if len(program) == 0 {
		t.Fatal("parser não produziu nenhum nó")
	}

	// =========================
	// FUNCTION DECLARATION
	// =========================

	function, ok := program[0].(ast.FunctionDeclaration)

	if !ok {
		t.Fatalf(
			"esperada FunctionDeclaration, recebido %T",
			program[0],
		)
	}

	if function.Name != "dobro" {
		t.Fatalf(
			"nome esperado %q, recebido %q",
			"dobro",
			function.Name,
		)
	}

	if function.ReturnType != "number" {
		t.Fatalf(
			"tipo de retorno esperado %q, recebido %q",
			"number",
			function.ReturnType,
		)
	}

	// =========================
	// PARAMETER
	// =========================

	if len(function.Parameters) != 1 {
		t.Fatalf(
			"esperado 1 parâmetro, recebido %d",
			len(function.Parameters),
		)
	}

	parameter := function.Parameters[0]

	if parameter.Name != "x" {
		t.Fatalf(
			"nome do parâmetro esperado %q, recebido %q",
			"x",
			parameter.Name,
		)
	}

	if parameter.Type != "number" {
		t.Fatalf(
			"tipo do parâmetro esperado %q, recebido %q",
			"number",
			parameter.Type,
		)
	}

	// =========================
	// FUNCTION BODY
	// =========================

	if len(function.Body) != 1 {
		t.Fatalf(
			"esperado 1 declaração no corpo, recebido %d",
			len(function.Body),
		)
	}

	returnStatement, ok := function.Body[0].(ast.ReturnStatement)

	if !ok {
		t.Fatalf(
			"esperado ReturnStatement, recebido %T",
			function.Body[0],
		)
	}

	// =========================
	// RETURN EXPRESSION
	// =========================

	binary, ok := returnStatement.Value.(ast.BinaryExpression)

	if !ok {
		t.Fatalf(
			"esperado BinaryExpression, recebido %T",
			returnStatement.Value,
		)
	}

	if binary.Operator.Type != token.MULTIPLY {
		t.Fatalf(
			"operador esperado *, recebido %v",
			binary.Operator,
		)
	}

	// =========================
	// LEFT SIDE: x
	// =========================

	left, ok := binary.Left.(ast.VariableExpression)

	if !ok {
		t.Fatalf(
			"lado esquerdo deveria ser VariableExpression, recebido %T",
			binary.Left,
		)
	}

	if left.Name != "x" {
		t.Fatalf(
			"variável esperada %q, recebida %q",
			"x",
			left.Name,
		)
	}

	// =========================
	// RIGHT SIDE: 2
	// =========================

	_, ok = binary.Right.(ast.NumberLiteral)

	if !ok {
		t.Fatalf(
			"lado direito deveria ser NumberLiteral, recebido %T",
			binary.Right,
		)
	}
}

// assertToken verifica o tipo e o lexema de um token.
func assertToken(
	t *testing.T,
	actual token.Token,
	expectedType token.TokenType,
	expectedLexeme string,
) {
	t.Helper()

	if actual.Type != expectedType {
		t.Fatalf(
			"tipo incorreto: esperado %v, recebido %v",
			expectedType,
			actual.Type,
		)
	}

	if actual.Lexeme != expectedLexeme {
		t.Fatalf(
			"lexema incorreto: esperado %q, recebido %q",
			expectedLexeme,
			actual.Lexeme,
		)
	}
}