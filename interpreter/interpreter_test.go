package interpreter

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ArthurBryan166/mango/ast"
	"github.com/ArthurBryan166/mango/lexer"
	"github.com/ArthurBryan166/mango/parser"
	"github.com/ArthurBryan166/mango/token"
)

// ============================================================
// HELPERS
// ============================================================

// runSource executa um programa Mango completo e retorna sua saída.
func runSource(t *testing.T, source string) string {
	t.Helper()

	// Lexer
	l := lexer.New(source)

	tokens, err := l.ScanTokens()
	if err != nil {
		t.Fatalf("erro no lexer: %v", err)
	}

	// Parser
	p := parser.New(tokens)

	program, err := p.Parse()
	if err != nil {
		t.Fatalf("erro no parser: %v", err)
	}

	// Interpreter
	interpreter := New()

	var output bytes.Buffer

	interpreter.writer = bufio.NewWriter(&output)

	if err := interpreter.Run(program); err != nil {
		t.Fatalf("erro no interpretador: %v", err)
	}

	if err := interpreter.writer.Flush(); err != nil {
		t.Fatalf("erro ao finalizar saída: %v", err)
	}

	return output.String()
}

func runSourceWithInput(
	t *testing.T,
	source string,
	input string,
) string {
	t.Helper()

	l := lexer.New(source)

	tokens, err := l.ScanTokens()
	if err != nil {
		t.Fatalf("erro no lexer: %v", err)
	}

	p := parser.New(tokens)

	program, err := p.Parse()
	if err != nil {
		t.Fatalf("erro no parser: %v", err)
	}

	interpreter := New()

	interpreter.reader = bufio.NewReader(
		bytes.NewBufferString(input),
	)

	var output bytes.Buffer

	interpreter.writer = bufio.NewWriter(&output)

	if err := interpreter.Run(program); err != nil {
		t.Fatalf("erro no interpretador: %v", err)
	}

	if err := interpreter.writer.Flush(); err != nil {
		t.Fatalf("erro ao finalizar saída: %v", err)
	}

	return output.String()
}

// runSourceError executa um programa que deve produzir erro.
func runSourceError(t *testing.T, source string) error {
	t.Helper()

	l := lexer.New(source)

	tokens, err := l.ScanTokens()
	if err != nil {
		return err
	}

	p := parser.New(tokens)

	program, err := p.Parse()
	if err != nil {
		return err
	}

	interpreter := New()

	var output bytes.Buffer
	interpreter.writer = bufio.NewWriter(&output)

	return interpreter.Run(program)
}

// assertOutput verifica exatamente a saída produzida.
func assertOutput(t *testing.T, source string, expected string) {
	t.Helper()

	output := runSource(t, source)

	if output != expected {
		t.Fatalf(
			"saída esperada %q, recebida %q",
			expected,
			output,
		)
	}
}

// assertError verifica apenas se ocorreu um erro.
func assertError(t *testing.T, source string) {
	t.Helper()

	err := runSourceError(t, source)

	if err == nil {
		t.Fatal("era esperado um erro, mas nenhum ocorreu")
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

// ============================================================
// LEXER + PARSER
// ============================================================

func TestFunctionCall(t *testing.T) {
	source := `
		mango dobro(x number) number {
			return x * 2
		}

		mango main() void {
			mangout(dobro(21))
		}
	`

	// -------------------------
	// Lexer
	// -------------------------

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

	// -------------------------
	// Parser
	// -------------------------

	p := parser.New(tokens)

	program, err := p.Parse()
	if err != nil {
		t.Fatalf("erro no parser: %v", err)
	}

	if len(program) == 0 {
		t.Fatal("parser não produziu nenhum nó")
	}

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

	// -------------------------
	// Parameter
	// -------------------------

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

	// -------------------------
	// Function body
	// -------------------------

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

	// -------------------------
	// Return expression
	// -------------------------

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

	_, ok = binary.Right.(ast.NumberLiteral)
	if !ok {
		t.Fatalf(
			"lado direito deveria ser NumberLiteral, recebido %T",
			binary.Right,
		)
	}
}

// ============================================================
// MANGOUT
// ============================================================

func TestMangout(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(42)
		}
	`, "42\n")
}

func TestMangoutExpression(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(10 + 5)
		}
	`, "15\n")
}

func TestMangoutString(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout("Olá")
		}
	`, "Olá\n")
}

func TestMangoutBoolean(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(true)
		}
	`, "true\n")
}

func TestMangoutMultipleArguments(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout("Valor: ", 42)
		}
	`, "Valor: 42\n")
}

// ============================================================
// VARIÁVEIS
// ============================================================

func TestVariableDeclaration(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mng x <- 10
			mangout(x)
		}
	`, "10\n")
}

func TestStringVariable(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mng nome <- "Arthur"
			mangout(nome)
		}
	`, "Arthur\n")
}

func TestBooleanVariable(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mng ativo <- true
			mangout(ativo)
		}
	`, "true\n")
}

func TestVariableAssignment(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mng x <- 10
			x <- 15
			mangout(x)
		}
	`, "15\n")
}

func TestAssignmentExpression(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mng x <- 10
			x <- x + 5
			mangout(x)
		}
	`, "15\n")
}

// ============================================================
// OPERAÇÕES ARITMÉTICAS
// ============================================================

func TestAddition(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(10 + 5)
		}
	`, "15\n")
}

func TestSubtraction(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(10 - 5)
		}
	`, "5\n")
}

func TestMultiplication(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(10 * 5)
		}
	`, "50\n")
}

func TestDivision(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(10 / 2)
		}
	`, "5\n")
}

func TestUnaryMinus(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(-10)
		}
	`, "-10\n")
}

func TestOperatorPrecedence(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(2 + 3 * 4)
		}
	`, "14\n")
}

func TestParentheses(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout((2 + 3) * 4)
		}
	`, "20\n")
}

// ============================================================
// COMPARAÇÕES
// ============================================================

func TestEqualNumbers(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(10 = 10)
		}
	`, "true\n")
}

func TestDifferentNumbers(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(10 = 5)
		}
	`, "false\n")
}

func TestLessThan(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(5 < 10)
		}
	`, "true\n")
}

func TestLessEqual(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(10 <= 10)
		}
	`, "true\n")
}

func TestGreaterThan(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(10 > 5)
		}
	`, "true\n")
}

func TestGreaterEqual(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(10 >= 10)
		}
	`, "true\n")
}

// ============================================================
// OPERADORES LÓGICOS
// ============================================================

func TestLogicalAnd(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(true and true)
		}
	`, "true\n")
}

func TestLogicalAndFalse(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(true and false)
		}
	`, "false\n")
}

func TestLogicalOr(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(false or true)
		}
	`, "true\n")
}

func TestLogicalOrFalse(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(false or false)
		}
	`, "false\n")
}

func TestLogicalNot(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mangout(not true)
		}
	`, "false\n")
}

// ============================================================
// IF / ELSE
// ============================================================

func TestIfTrue(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			if true {
				mangout("sim")
			}
		}
	`, "sim\n")
}

func TestIfFalse(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			if false {
				mangout("sim")
			}
		}
	`, "")
}

func TestIfElse(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			if false {
				mangout("sim")
			} else {
				mangout("não")
			}
		}
	`, "não\n")
}

func TestElseIf(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mng x <- 2

			if x = 1 {
				mangout("um")
			} else if x = 2 {
				mangout("dois")
			} else {
				mangout("outro")
			}
		}
	`, "dois\n")
}

func TestIfWithVariable(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mng idade <- 18

			if idade >= 18 {
				mangout("maior")
			}
		}
	`, "maior\n")
}

// ============================================================
// FOR
// ============================================================

func TestForClassic(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			for mng i <- 0; i < 3; i <- i + 1 {
				mangout(i)
			}
		}
	`, "0\n1\n2\n")
}

func TestForCondition(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mng i <- 0

			for i < 3 {
				mangout(i)
				i <- i + 1
			}
		}
	`, "0\n1\n2\n")
}

func TestForInfiniteStyleWithCondition(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mng i <- 3

			for i > 0 {
				mangout(i)
				i <- i - 1
			}
		}
	`, "3\n2\n1\n")
}

func TestForEmpty(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			for false {
				mangout("não")
			}
		}
	`, "")
}

// ============================================================
// FUNÇÕES
// ============================================================

func TestFunctionCallExecution(t *testing.T) {
	assertOutput(t, `
		mango dobro(x number) number {
			return x * 2
		}

		mango main() void {
			mangout(dobro(21))
		}
	`, "42\n")
}

func TestFunctionMultipleParameters(t *testing.T) {
	assertOutput(t, `
		mango soma(a number, b number) number {
			return a + b
		}

		mango main() void {
			mangout(soma(10, 20))
		}
	`, "30\n")
}

func TestFunctionStringParameter(t *testing.T) {
	assertOutput(t, `
		mango mostrar(nome string) void {
			mangout("Olá ", nome)
		}

		mango main() void {
			mostrar("Arthur")
		}
	`, "Olá Arthur\n")
}

func TestVoidFunction(t *testing.T) {
	assertOutput(t, `
		mango mostrar() void {
			mangout("Olá")
		}

		mango main() void {
			mostrar()
		}
	`, "Olá\n")
}

func TestFunctionCallingFunction(t *testing.T) {
	assertOutput(t, `
		mango dobro(x number) number {
			return x * 2
		}

		mango quadrado(x number) number {
			return x * x
		}

		mango main() void {
			mangout(quadrado(dobro(3)))
		}
	`, "36\n")
}

// ============================================================
// RECURSÃO
// ============================================================

func TestRecursion(t *testing.T) {
	assertOutput(t, `
		mango fatorial(n number) number {
			if n = 0 {
				return 1
			}

			return n * fatorial(n - 1)
		}

		mango main() void {
			mangout(fatorial(5))
		}
	`, "120\n")
}

// ============================================================
// ESCOPO
// ============================================================

func TestLocalVariable(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			if true {
				mng x <- 10
				mangout(x)
			}
		}
	`, "10\n")
}

func TestNestedScope(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mng x <- 10

			if true {
				mng y <- 20

				if true {
					mangout(x)
					mangout(y)
				}
			}
		}
	`, "10\n20\n")
}

func TestFunctionScope(t *testing.T) {
	assertOutput(t, `
		mango mostrar() void {
			mng x <- 10
			mangout(x)
		}

		mango main() void {
			mostrar()
		}
	`, "10\n")
}

func TestVariableShadowing(t *testing.T) {
	assertOutput(t, `
		mango main() void {
			mng x <- 10

			if true {
				mng x <- 20
				mangout(x)
			}

			mangout(x)
		}
	`, "20\n10\n")
}

// ============================================================
// ERROS DE VARIÁVEIS
// ============================================================

func TestRedeclarationError(t *testing.T) {
	assertError(t, `
		mango main() void {
			mng x <- 10
			mng x <- 20
		}
	`)
}

func TestUndefinedVariableError(t *testing.T) {
	assertError(t, `
		mango main() void {
			mangout(x)
		}
	`)
}

func TestAssignmentUndefinedVariableError(t *testing.T) {
	assertError(t, `
		mango main() void {
			x <- 10
		}
	`)
}

// ============================================================
// ERROS DE TIPOS
// ============================================================

func TestInvalidAdditionTypes(t *testing.T) {
	assertError(t, `
		mango main() void {
			mangout(10 + "Olá")
		}
	`)
}

func TestInvalidSubtractionTypes(t *testing.T) {
	assertError(t, `
		mango main() void {
			mangout("Olá" - "Mundo")
		}
	`)
}

func TestInvalidComparisonTypes(t *testing.T) {
	assertError(t, `
		mango main() void {
			mangout(10 < "Olá")
		}
	`)
}

func TestInvalidNotType(t *testing.T) {
	assertError(t, `
		mango main() void {
			mangout(not 10)
		}
	`)
}

func TestDivisionByZero(t *testing.T) {
	assertError(t, `
		mango main() void {
			mangout(10 / 0)
		}
	`)
}

// ============================================================
// ERROS DE FUNÇÕES
// ============================================================

func TestWrongArgumentCount(t *testing.T) {
	assertError(t, `
		mango soma(a number, b number) number {
			return a + b
		}

		mango main() void {
			mangout(soma(10))
		}
	`)
}

func TestWrongArgumentType(t *testing.T) {
	assertError(t, `
		mango dobro(x number) number {
			return x * 2
		}

		mango main() void {
			mangout(dobro("Olá"))
		}
	`)
}

func TestWrongReturnType(t *testing.T) {
	assertError(t, `
		mango teste() number {
			return "Olá"
		}

		mango main() void {
			teste()
		}
	`)
}

func TestMissingReturn(t *testing.T) {
	assertError(t, `
		mango teste() number {
			mng x <- 10
		}

		mango main() void {
			teste()
		}
	`)
}

func TestVoidReturnError(t *testing.T) {
	assertError(t, `
		mango teste() void {
			return 10
		}

		mango main() void {
			teste()
		}
	`)
}

// ============================================================
// MAIN
// ============================================================

func TestMainMustExist(t *testing.T) {
	assertError(t, `
		mango teste() void {
			mangout("Olá")
		}
	`)
}

func TestMainMustBeVoid(t *testing.T) {
	assertError(t, `
		mango main() number {
			return 10
		}
	`)
}

func TestMainMustHaveNoParameters(t *testing.T) {
	assertError(t, `
		mango main(x number) void {
			mangout(x)
		}
	`)
}

// ============================================================
// INTEGRAÇÃO
// ============================================================

func TestCompleteProgram(t *testing.T) {
	assertOutput(t, `
		mango dobro(x number) number {
			return x * 2
		}

		mango fatorial(n number) number {
			if n = 0 {
				return 1
			}

			return n * fatorial(n - 1)
		}

		mango mostrarNumero(x number) void {
			mangout("Número: ", x)
		}

		mango main() void {
			mng x <- 10
			mng ativo <- true

			if ativo and x > 5 {
				mostrarNumero(dobro(x))
			}

			mangout(fatorial(5))

			for mng i <- 0; i < 3; i <- i + 1 {
				mangout(i)
			}
		}
	`, "Número: 20\n120\n0\n1\n2\n")
}

func TestMangoinNumber(t *testing.T) {
	output := runSourceWithInput(t, `
		mango main() void {
			mng idade <- 0

			mangoin(idade)

			mangout("Idade: ", idade)
		}
	`, "16\n")

	expected := "Idade: 16\n"

	if output != expected {
		t.Fatalf(
			"saída esperada %q, recebida %q",
			expected,
			output,
		)
	}
}

func TestNil(t *testing.T) {
    source := `
        mango main() void {
            mangout(nil)
            mangout(nil = nil)
            mangout(nil = 10)
            mangout(nil = "mango")
            mangout(nil != nil)
        }
    `

    output := runSourceWithInput(t, source, "")

    expected := "<nil>\ntrue\nfalse\nfalse\nfalse\n"

    if output != expected {
        t.Fatalf(
            "saída inesperada:\nobtida:\n%q\nesperada:\n%q",
            output,
            expected,
        )
    }
}

func TestReturnWithoutValue(t *testing.T) {
	source := `
		mango mostrar() void {
			mangout("antes")
			return
			mangout("depois")
		}

		mango main() void {
			mostrar()
		}
	`

	output := runSourceWithInput(t, source, "")

	expected := "antes\n"

	if output != expected {
		t.Fatalf(
			"saída inesperada:\nobtida:\n%q\nesperada:\n%q",
			output,
			expected,
		)
	}
}

func TestReturnValueInVoidFunction(t *testing.T) {
	source := `
		mango teste() void {
			return 10
		}

		mango main() void {
			teste()
		}
	`

	l := lexer.New(source)
	tokens, err := l.ScanTokens()
	if err != nil {
		t.Fatalf("erro no lexer: %v", err)
	}

	p := parser.New(tokens)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("erro no parser: %v", err)
	}

	interpreter := New()

	err = interpreter.Run(program)

	if err == nil {
		t.Fatal("esperado erro ao retornar valor em função void")
	}
}

func TestEmptyReturnInNonVoidFunction(t *testing.T) {
	source := `
		mango teste() number {
			return
		}

		mango main() void {
			mangout(teste())
		}
	`

	l := lexer.New(source)
	tokens, err := l.ScanTokens()
	if err != nil {
		t.Fatalf("erro no lexer: %v", err)
	}

	p := parser.New(tokens)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("erro no parser: %v", err)
	}

	interpreter := New()

	err = interpreter.Run(program)

	if err == nil {
		t.Fatal("esperado erro ao usar return sem valor em função number")
	}
}