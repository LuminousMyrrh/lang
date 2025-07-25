package main

import (
	"fmt"
	"os"
)

func eval(source string) (error) {
	globalEnv := NewEnv(nil, "global")
	lexer := Lexer{};
	toks, err := lexer.Read(source);
	//fmt.Println("Scanning done")
	parser := NewParser(toks);
	if err != nil {
		return err
	}

	if len(toks) == 0 {
		fmt.Println("No tokens found");
		return nil
	}

	mnode, errs := parser.Parse()

	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
	} else {
		// mnode.Print()
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

func enterRepl() {
	repl := NewRepl()

	repl.Start()
}

func main() {
	args := os.Args
	if len(args) == 1 {

	} else {
		if args[1] == "repl" {
			enterRepl()
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
