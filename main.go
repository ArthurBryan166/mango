package main

import (
    "fmt"
    "github.com/ArthurBryan166/mango.git/lexer"
)

func main() {
    source := `
	mng idade <- 18 // idade
	idade = 20
	`

    l := lexer.New(source)

    tokens, err := l.ScanTokens()

    if err != nil {
        fmt.Println("Erro:", err)
        return
    }

    for _, tok := range tokens {
        fmt.Printf("%+v\n", tok)
    }
}