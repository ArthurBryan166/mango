# Mango

Mango é uma linguagem de programação interpretada criada em Go como projeto de estudos.

A linguagem possui uma sintaxe própria e atualmente conta com variáveis, funções, estruturas condicionais, laços de repetição, entrada e saída de dados, operadores lógicos e tratamento de valores.

## Exemplo

```mango
mango main() void {
    mng nome <- ""

    mangoin(nome)

    mangout("Olá, ", nome)

    if nome = "Arthur" {
        mangout("Bem-vindo!")
    } else {
        mangout("Olá!")
    }
}
```

## Funcionalidades

Algumas das funcionalidades presentes no Mango:

- Variáveis com inferência de tipo
- Tipos `number`, `string` e `bool`
- Valor `nil`
- Operadores matemáticos
- Operadores de comparação
- Operadores lógicos `and`, `or` e `not`
- `if`, `else if` e `else`
- Laços `for`
- Funções
- Recursão
- `return`
- Entrada de dados com `mangoin`
- Saída de dados com `mangout`
- Diferentes escopos para variáveis
- Interface de linha de comando

## Sintaxe

### Variáveis

Para criar uma variável, é utilizado `mng`.

```mango
mng idade <- 16
mng nome <- "Arthur"
mng aprovado <- true
mng sobrenome <- ""
```

O tipo da variável é definido automaticamente de acordo com o valor utilizado na sua criação.

É obrigatório conceder um valor inicial à variável.

Para alterar o valor de uma variável, é utilizado `<-`.

```mango
idade <- 17
```

### Operadores

O Mango possui operadores matemáticos como:

```text
+
-
*
/
%
```

Também possui operadores de comparação:

```text
>
>=
<
<=
=
!=
```

E operadores lógicos:

```text
and
or
not
```

Exemplo:

```mango
if idade >= 16 and aprovado = true {
    mangout("Pode continuar.")
}
```

### Condicionais

O Mango possui `if`, `else if` e `else`.

```mango
if idade >= 18 {
    mangout("Maior de idade")
} else if idade >= 16 {
    mangout("Adolescente")
} else {
    mangout("Menor de 16")
}
```

### Laços

O `for` é utilizado para repetição.

Uma das formas possíveis é utilizar apenas uma condição:

```mango
mng i <- 0

for i < 5 {
    mangout(i)
    i <- i + 1
}
```

Também existe a forma com inicialização, condição e incremento:

```mango
for mng i <- 0; i < 5; i <- i + 1 {
    mangout(i)
}
```

### Funções

Funções são criadas utilizando `mango`.

```mango
mango soma(a number, b number) number {
    return a + b
}
```

Os parâmetros possuem tipos definidos.

Para utilizar a função:

```mango
mng resultado <- soma(10, 20)

mangout(resultado)
```

Funções também podem chamar outras funções, incluindo a própria função, permitindo fazer recursão.

### Entrada de dados

A entrada de dados é feita utilizando `mangoin`.

A variável que receberá o valor é passada para a função:

```mango
mng nome <- ""

mangoin(nome)
```

### Saída de dados

A saída de dados é feita utilizando `mangout`.

```mango
mangout("Olá, mundo!")
```

Também é possível passar vários valores separados por vírgulas:

```mango
mng nome <- "Arthur"

mangout("Olá, ", nome, "!")
```

### Retorno

Funções podem retornar valores utilizando `return`.

```mango
mango dobro(numero number) number {
    return numero * 2
}
```

Funções `void` também podem utilizar `return` sem um valor:

```mango
mango exemplo() void {
    mangout("Antes")

    return

    mangout("Depois")
}
```

Nesse caso, o código depois do `return` não é executado.

### Nil

O Mango possui o valor `nil`.

Ele pode ser utilizado e comparado como um valor:

```mango
mng valor <- nil

mangout(valor = nil)
```

## Função main

Todo programa Mango precisa possuir uma função `main`.

Ela não recebe parâmetros e deve possuir retorno `void`.

```mango
mango main() void {
    mangout("Olá, Mango!")
}
```

A execução do programa começa pela função `main`.

## Executando um programa

Os arquivos Mango utilizam a extensão `.mg`.

Com o Mango instalado, um programa pode ser executado utilizando:

```bash
mango run programa.mg
```

Por exemplo:

```bash
mango run exemplo.mg
```

## Versão

### Mango 1.0.0

Esta é a primeira versão estável do Mango.

O conjunto de funcionalidades desta versão já está definido. Novas funcionalidades e mudanças maiores poderão ser feitas em versões futuras.

## Sobre o projeto

O Mango começou como um projeto para aprender mais sobre como linguagens de programação funcionam por dentro e para diversão própria.

A ideia foi construir a linguagem aos poucos, entendo como funciona cada etapa de funcionamento da linguagem, desde o código fonte até no resultado no terminal.

Este projeto também serve como uma forma de praticar Go, envolvendo organização de código, estruturas de dados, tratamento de erros e interpretação de programas.

## Instalação

Os executáveis do Mango estão disponíveis na página de releases.

Baixe a versão correspondente ao seu sistema operacional e arquitetura e coloque o executável no `PATH`.

Depois, execute um programa com:

```bash
mango run programa.mg
