package repl

import (
	"bufio"
	"fmt"
	"lang/internal/env"
	"lang/internal/eval"
	"lang/internal/lexer"
	"lang/internal/parser"
	"os"
	"strings"
)

type Repl struct {
	Eval *eval.Evaluator
	Environment *env.Env
}

func NewRepl() *Repl {
	return &Repl {
		Eval: &eval.Evaluator{},
		Environment: env.NewEnv(nil, "global"),
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
	lexer := lexer.Lexer{}
	toks, err := lexer.Read(source)
	if err != nil {
		return err
	}
	parser := parser.NewParser(toks)
	mnode, errs := parser.Parse()

	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
	} else {
		// mnode.Print()
		evaluator := eval.NewEvaluator(r.Environment, mnode)
		evaluator.Eval()
		errs = evaluator.Errors
		if len(errs) != 0 {
			for _, err := range errs {
				fmt.Println("Runtime error: ", err)
			}
		} else {
			r.Environment = evaluator.Environment
		}
	}

	return nil
}
