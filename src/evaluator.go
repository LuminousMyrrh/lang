package main

import (
	"errors"
	"fmt"
)

type Evaluator struct {
	Entry *ProgramNode
	Environment *Env
	currentEnv *Env
	Errors []error
}

func (e *Evaluator) Eval(env *Env, entry *ProgramNode) {
	e.Environment = env
	e.Entry = entry
	e.currentEnv = env

	for _, stmt := range entry.Nodes{
		e.eval(stmt)
	}
}

func (e *Evaluator) eval(stmt Node) any {
	fmt.Printf("Evaluating node: %T\n", stmt)

	switch s := stmt.(type) {
	case VarDefNode: {
		return e.evalVarDef(s)
	}
	case FunctionDefNode: {
		return e.evalFunctionDef(s)
	}
	case BinaryOpNode: {
		return e.evalBinary(s)
	}
	case IdentifierNode: {
		return e.evalIdentifier(s)
	}
	case FunctionCallNode: {
		return e.evalFunctionCall(s)
	}
	case LiteralNode: {
		return e.evalLiteral(s)
	}
	case ReturnNode:
		return e.evalReturn(s)
	default: {
		return nil
	}
	}
}

func (e *Evaluator) evalReturn(ret ReturnNode) any {
	return e.eval(ret.Value)
}

func (e *Evaluator) evalVarDef(stmt VarDefNode) any {
	value := e.eval(stmt.Value)
	e.currentEnv.AddVarSymbol(stmt.Name, "", value)
	return value
}

func (e *Evaluator) evalLiteral(lit LiteralNode) any {
	return lit.Value
}

func (e *Evaluator) evalIdentifier(id IdentifierNode) any {
	return e.currentEnv.FindSymbol(id.Name)
}

func (e *Evaluator) evalBinary(expr BinaryOpNode) any {
	left := e.eval(expr.Left)
	right := e.eval(expr.Right)
	if _, ok := left.(int); !ok {
		e.genError(fmt.Sprintf("Failed to convert to int: %v", left))
		e.genError("Failed to convert to int")
		return nil
	}

	if _, ok := right.(int); !ok {
		e.genError(fmt.Sprintf("Failed to convert to int: %v", right))
		return nil
	}

	switch expr.Op {
	case "+": {
		return left.(int) + right.(int)
	}
	case "-": {
		return left.(int) - right.(int)
	}
	case "/": {
		return left.(int) / right.(int)
	}
	case "*": {
		return left.(int) * right.(int)
	}
	}

	return nil
}

func (e *Evaluator) genError(msg string) {
	e.Errors = append(e.Errors, errors.New(msg))
}

func (e *Evaluator) evalFunctionDef(funcDef FunctionDefNode) any {
	e.currentEnv.AddFuncSymbol(
		funcDef.Name,
		funcDef.Parameters,
		funcDef.Body,
		NewEnv(e.currentEnv),
		)
	return nil
}

func (e *Evaluator) evalFunctionCall(call FunctionCallNode) any {
	if call.Name == "print" {
		for _, arg := range call.Args {
			val := e.eval(arg)
			fmt.Print(val)
		}
		fmt.Println()
		return nil
	}

	fmt.Println(call.String())
	function := e.currentEnv.FindSymbol(call.Name)
	if function == nil {
		e.genError(fmt.Sprintf("Function '%s' not found", call.Name))
		return nil
	}

	switch f := function.(type) {
	case FuncSymbol: {
		params := f.Params
		if len(call.Args) != len(params) {
			e.genError(fmt.Sprintf(
				"Function '%s' accepts %d, but passed only %d",
					call.Name,
					len(params),
					len(call.Args)))

			return nil
		}

		prevEnv := e.currentEnv
		e.currentEnv = f.Enf

		// evaling args and adding them to function env
		for i, arg := range call.Args {
			value := e.eval(arg)
			e.currentEnv.AddVarSymbol(params[i], "", value)
		}

		for _, stmt := range f.Body.Statements {
			e.eval(stmt)
		}

		e.currentEnv = prevEnv
	}
	default:
		e.genError("Something very very wrong here")
		return nil
	}

	return nil
}
