package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func eval(source string) (error) {
	globalEnv := NewEnv(nil)
	lexer := Lexer{};
	toks, err := lexer.Read(source);
	fmt.Println("Scanning done")
	parser := NewParser(globalEnv, toks);
	if err != nil {
		return err
	}

	if len(toks) == 0 {
		fmt.Println("No tokens found");
		return nil
	}
	for _, tok := range toks {
		fmt.Printf("Tok: %v Type: %v\n", tok.Lexeme, tok.TType)
	}

	mnode, errs := parser.Parse(toks)
	fmt.Println("Parsing done")

	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
	} else {
		//mnode.Print()
		evaluator := Evaluator{}
		evaluator.Eval(globalEnv, mnode)
		if len(evaluator.Errors) > 0 {
			for _, err := range evaluator.Errors {
				fmt.Println("Runtime error: ", err)
			}
		}
	}

	return nil
}

func repl() {
	reader := bufio.NewReader(os.Stdin)
	var history []string

	for {
		fmt.Print("&> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		history = append(history, input)
		if input == "exit" || input == ":q" {
			break
		}
		if err := eval(input); err != nil {
			fmt.Println("Error: ", err)
		}
	}
}

func main() {
	args := os.Args
	if len(args) == 1 {

	} else {
		if args[1] == "repl" {
			repl()
		} else {
			fileName := args[1]
			data, err := os.ReadFile(fileName)
			if err != nil {
				fmt.Println("Failed to read file: ", err)
				return
			}
			content := string(data)
			eval(content)
		}
	}
}
