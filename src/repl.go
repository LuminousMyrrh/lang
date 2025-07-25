package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Repl struct {
	Eval *Evaluator
	Environment *Env
}

func NewRepl() *Repl {
	return &Repl {
		Eval: &Evaluator{},
		Environment: NewEnv(nil, "global"),
	}
}

func (r *Repl) Start() {
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
		if err := r.eval(input); err != nil {
			fmt.Println("Error: ", err)
		}
	}
}

func (r *Repl) eval(source string) error {
	lexer := Lexer{}
	toks, err := lexer.Read(source)
	if err != nil {
		return err
	}
	parser := NewParser(toks)
	mnode, errs := parser.Parse()

	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
	} else {
		// mnode.Print()
		evaler := Evaluator{}
		evaler.Eval(r.Environment, mnode)
		errs = evaler.Errors
		if len(errs) != 0 {
			for _, err := range errs {
				fmt.Println("Runtime error: ", err)
			}
		} else {
			r.Environment = evaler.Environment
		}
	}

	return nil
}
