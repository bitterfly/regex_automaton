# Brief

This project is for the "Applications of Finite Automata" course at FMI (Faculty of Mathematics and Informatics) written in Go.

It builds the minimal automaton for a given dictionary. The program then can answer queries in the form of regular expressions. It returns the words from the dictionary that match the regular expression.
It recognises UTF-8 symbols.

The automaton from the dictinary is built once. For every given expression a new automaton is built using the [Thompson's construction](https://en.wikipedia.org/wiki/Thompson%27s_construction), the epsilon transitions are eliminated and the resulting automaton is then intersected with the dictionary trie automaton.

# Usage:

## Requirements

An installed version of Go.
Mind that large queries take a lot of RAM (a.k.a. 10GB for 72MB dictionary ~ 1,000,000 words)

The dictionary should be sorted and should not contain spaces or dashes (the dashes I think depend on the language and sorting).

## Installation

    > go get github.com/bitterfly/regex_automaton@latest


## Running

    > regex_automaton [OPTIONS] dictionary_file < regular_expression

There are several options:

* **-output** -  specifies file into which the matched words will be written instead of stdin.
* **-infix** - true if the regular expression is in infix notation or false for reverse polish notation. It is true by default.
* **-verbose** - if true debug information will be printed.

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

* Non-latin (matches мама and баба)

`(м|б).а.(м|б).а`

### RPN

* Matches `color` or `colour`

`colou?|r.....`

* Non-latin (matches бира and гира)

`бг|ира...`

### Examples
Using a bulgarian dictionary from [this repository](https://github.com/miglen/bulgarian-wordlists/blob/master/wordlists/bg-words-cyrillic.txt).

* Match мечка/печка
```
> echo "(п|м).е.ч.к.а" | regex_automaton <(grep -v "[- ]" bg-words-cyrillic.txt | sort)

мечка
печка

> echo ((а|б|в|г|д|е|ж|з|и|й|к|л|м|н|о|п|р|с|т|у|ф|х|ц|ч|ш|щ|ъ|ь|ю|я)*).в.а.н.е | regex_automaton <(grep -v "[- ]" bg-words-cyrillic.txt | sort)

адвокатстване
активоване
акушерстване
акцентуване
анатемосване
...
```

Using this [english dictionary](https://github.com/dwyl/english-words/blob/master/words.txt):

* Match words containing "cool".

```
> any="((a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)*)"
> echo "${any}.c.o.o.l.${any}" | go run project.go <(grep -v "[-/ ,.&']" /tmp/foo | awk '{print tolower($0)}' | sort -u)

acool
aftercooler
bellacoola
cool
coolabah
```

# TODO
Думите, които още не са обработени (т.е. reduce-нати), може да не се слагат направо в автомата, а да се пазят в статична структура и чак след обработка да се добавят.

Важно. В Автомата за регулярните изрази, от всяко състояние има най-много 2 прехода, като или има един с буква, или най-много два с епсилон. Това означава, че структурата за преходи може да се подобри **многократно**.
