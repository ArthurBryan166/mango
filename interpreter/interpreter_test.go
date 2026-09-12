package interpreter

import "testing"

func TestFunctionCall(t *testing.T) {
    source := `
        mango dobro(x number) number {
            return x * 2
        }

        mango main() void {
            mangout(dobro(21))
        }
    `

    // Aqui vamos executar o código Mango.
}