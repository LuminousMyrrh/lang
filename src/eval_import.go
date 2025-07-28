package main

import (
	"fmt"
	"os"
)

func (e *Evaluator) evalImport(stmt *ImportNode) any {
	fileName := stmt.File + ".lang"
	data, err := os.ReadFile(fileName)
	if err != nil {
		e.genError(fmt.Sprintf(
			"Failed to read imported file: %s", err),
			stmt.Position,
			)
		return nil
	}

	content := string(data)
	e.lexer = &Lexer{};
	toks, err := e.lexer.Read(content)
	if err != nil {
		e.genError(fmt.Sprintf("Failed to read file: %s", err),
			stmt.Position,
			)
		return nil
	}
	e.parser = NewParser(toks)
	mnode, errs := e.parser.Parse()
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		return nil
	}

	// -----------------

	if len(stmt.Symbols) == 0 { 

		// Utils
		structName := capitalizeFirstLetter(stmt.File)

		// utils
		varName := stmt.File

		importEnv := NewEnv(e.currentEnv, structName)
		e.currentEnv.AddStructSymbol(structName, importEnv)

		for _, stmt := range mnode.Nodes {
			if def, ok := stmt.(*FunctionDefNode); ok {
				e.currentEnv.AddStructMethod(
					structName,
					def.Name,
					def.Parameters,
					def.Body.Statements,
					)
			}
		}

		e.currentEnv.AddVarSymbol(
			varName, structName, importEnv)

		fmt.Println(importEnv)

		return 1
	} else {
		for _, symbol := range stmt.Symbols {
			node, err := mnode.Find(symbol)
			if err != nil {
				e.genError(err.Error(), stmt.Position)
				return nil
			}
			e.eval(node)
		}
		return 1

	}
}
