package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ArthurBryan166/mango/interpreter"
	"github.com/ArthurBryan166/mango/lexer"
	"github.com/ArthurBryan166/mango/parser"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("uso: mango run <arquivo.mg>")
	}

	switch os.Args[1] {
	case "run":
		run(os.Args[2:])

	default:
		log.Fatalf("comando desconhecido: %s", os.Args[1])
	}
}

func run(args []string) {
	if len(args) != 1 {
		log.Fatal("uso: mango run <arquivo.mg>")
	}

	filename := args[0]

	if filepath.Ext(filename) != ".mg" {
		log.Fatal("o arquivo deve possuir a extensão .mg")
	}

	source, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("erro ao ler arquivo: %v", err)
	}

	l := lexer.New(string(source))

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