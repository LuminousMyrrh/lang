package eval

import (
	"fmt"
	"lang/internal/core"
	"lang/internal/env"
	"lang/internal/parser"
)

func (e *Evaluator) evalFunctionDef(funcDef *parser.FunctionDefNode) any {
	e.currentEnv.AddFuncSymbol(
		funcDef.Name,
		funcDef.Parameters,
		funcDef.Body,
		env.NewEnv(e.currentEnv, "function"),
	)
	return 1
}

func (e *Evaluator) evalFunctionCall(call *parser.FunctionCallNode) any {
	ident, ok := call.Name.(*parser.IdentifierNode)
	if !ok {
		e.GenError("Function call must be on an identifier",
			call.Position)
		return nil
	}

	// builtins
	if builtin, ok := e.Builtins[ident.Name]; ok {
		return builtin(e, call.Args, call.Position)
	}

	// user defined
	function := e.currentEnv.FindSymbol(ident.Name)
	if function == nil {
		e.GenError(fmt.Sprintf("Function '%s' not found", call.Name),
			call.Position)
		return nil
	}

	f, ok := function.(*env.FuncSymbol)
	if !ok {
		e.GenError("Something very very wrong here", call.Position)
		return nil
	}

	params := f.Params
	if len(call.Args) != len(params) {
		e.GenError(fmt.Sprintf(
			"Function '%s' accepts %d, but passed only %d",
			call.Name,
			len(params),
			len(call.Args)),
			call.Position)

		return nil
	}

	// 1. Eval args
	argValues := make([]any, len(call.Args))
	for i, arg := range call.Args {
		argValues[i] = e.EvalNode(arg)
	}

	// 2. Switch to the function's environment
	prevEnv := e.currentEnv
	callEnv := env.NewEnv(f.Env, ident.Name)
	e.currentEnv = callEnv

	// 3. Add parameters to the new environment
	for i, val := range argValues {
		e.currentEnv.AddVarSymbol(f.Params[i],
			e.resolveType(val, call.Position), val)
	}

	// 4. Evaluate the function body
	var result any
	for _, stmt := range f.Body.Statements {
		result = e.EvalNode(stmt)
		if ret, ok := result.(core.ReturnValue); ok {
			e.currentEnv = prevEnv
			return ret.Value
		}
	}

	e.currentEnv = prevEnv
	return result
}

func (e *Evaluator) evalFuncBlock(block *parser.BlockNode) any {
	var result any
	for _, stmt := range block.Statements {
		result = e.EvalNode(stmt)
		if ret, ok := result.(core.ReturnValue); ok {
			return ret
		}
	}
	return result
}
