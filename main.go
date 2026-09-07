package main

import (
	"log"

	// "github.com/ArthurBryan166/mango/ast"
	"github.com/ArthurBryan166/mango/interpreter"
	"github.com/ArthurBryan166/mango/lexer"
	"github.com/ArthurBryan166/mango/parser"
)

func main() {
	source := `
	mango main() void {
		mangout(not 10)
	}
`

	l := lexer.New(source)

	tokens, err := l.ScanTokens()
	if err != nil {
		log.Fatal(err)
	}

	p := parser.New(tokens)

	nodes, err := p.Parse()
	if err != nil {
		log.Fatal(err)
	}

	i := interpreter.New()

	if err := i.Run(nodes); err != nil {
		log.Fatal(err)
	}
}
