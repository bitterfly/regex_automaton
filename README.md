# Brief

This project is for applications of automata  course at FMI(Faculty of Mathematics and Informatics) written in GO

It builds minimal automaton from given dictionary. The program then can answers queries in the form of regular expressions. It returns the words from the dictionary that match the regular expression.
It recognises UTF8 symbols.

# Usage:

## Requirements

An installed version of Go.
Mind that large queries take a lot of RAM (a.k.a. 10GB for 72MB dictionary ~ 1,000,000 words)

## Instalation
$ go install github.com/bitterfly/pka

## Running

pka dictionary_file < regular_expression
//or simply from STDIN

There are several options:

* **-output** -  specifies file into which the matched words will be written
* **-infix** - true if the regular expression is in infix notation or false for reverse polish notation. It is true by default.

pka -output="matched.txt" fictionary_file < regular_expression

# Regular expression format

The recognised special symbols are

* **(**

* **)**

* **.** - concatenation

* **|** - or

* **\*** - Kleene's star

* **?** - empty-word

## Examples:

### Infix:

* Matches every english word:
`(?|a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)*`

* Matches `color` or `colour`
`c.o.l.o.(u|?).r`

* Non-latin (matches мама and baba)
`(м|б).а.(м|б).а`

### RPN

* Matches `color` or `colour`
`colou?|r.....`

* Non-latin (matches бира and гира)
`бг|ира...`


# TODO
Думите, които още не са обработени (т.е. reduce-нати), може да не се слагат направо в автомата, а да се пазят в статична структура и чак след обработка да се добавят.

Важно. В Автомата за регулярните изрази, от всяко състояние има най-много 2 прехода, като или има един с буква, или най-много два с епсилон. Това означава, че структурата за преходи може да се подобри **многократно**.