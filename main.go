package main

import (
	"fmt"
	"log"

	"github.com/ArthurBryan166/mango/lexer"
	"github.com/ArthurBryan166/mango/parser"
)

func main() {
	source := `
    if idade >= 18 {
		mng resultado <- 1
	} else if idade >= 13 {
		mng resultado <- 2
	} else {
		mng resultado <- 3
	}
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
