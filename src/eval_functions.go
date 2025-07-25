package main

import (
	"fmt"
)

func (e *Evaluator) evalFunctionDef(funcDef *FunctionDefNode) any {
	e.currentEnv.AddFuncSymbol(
		funcDef.Name,
		funcDef.Parameters,
		funcDef.Body,
		NewEnv(e.currentEnv, "function"),
		)
	return 1
}

func (e *Evaluator) evalFunctionCall(call *FunctionCallNode) any {
	ident, ok := call.Name.(*IdentifierNode)
	if !ok {
		e.genError("Function call must be on an identifier",
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
		e.genError(fmt.Sprintf("Function '%s' not found", call.Name),
			call.Position)
		return nil
	}

	if f, ok := function.(*FuncSymbol); ok {
		params := f.Params
		if len(call.Args) != len(params) {
			e.genError(fmt.Sprintf(
					"Function '%s' accepts %d, but passed only %d",
					call.Name,
					len(params),
					len(call.Args)),
				call.Position)

			return nil
		}

		argValues := make([]any, len(call.Args))
		for i, arg := range call.Args {
			argValues[i] = e.eval(arg)
		}

		// 2. Switch to the function's environment
		prevEnv := e.currentEnv
		callEnv := NewEnv(f.Env, "function")
		e.currentEnv = callEnv

		// 3. Add parameters to the new environment
		for i, val := range argValues {
			e.currentEnv.AddVarSymbol(f.Params[i], "", val)
		}

		// 4. Evaluate the function body
		result := e.evalFuncBlock(f.Body)
		if ret, ok := result.(returnValue); ok {
			return ret.value
		}

		e.currentEnv = prevEnv

		return result

	}
	e.genError("Something very very wrong here", call.Position)
	return nil
}

func (e *Evaluator) evalFuncBlock(block *BlockNode) any {
	for _, stmt := range block.Statements {
		result := e.eval(stmt)
		if ret, ok := result.(returnValue); ok {
			return ret
		}
	}
	return nil
}

