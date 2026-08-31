package main

import (
	"fmt"
	"log"

	"github.com/ArthurBryan166/mango/lexer"
	"github.com/ArthurBryan166/mango/parser"
)

func main() {
	source := `
    mangout("Olá ", nome, ", tudo bem?")
    `

	l := lexer.New(source)

	tokens, err := l.ScanTokens()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("TOKENS:")
	for _, tok := range tokens {
		fmt.Printf("%s(%s)\n", tok.Type, tok.Lexeme)
	}

	p := parser.New(tokens)

	nodes, err := p.Parse()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nAST:")
	for _, node := range nodes {
		fmt.Println(node.String())
	}
}
