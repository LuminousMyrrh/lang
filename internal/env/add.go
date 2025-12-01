package env

import (
	"lang/internal/core"
	"lang/internal/parser"
)

func (e *Env) AddVarSymbol(name, varType string, val any) {
	e.Symbols[name] = &VarSymbol{value: val, typeName: varType}
}

func (e *Env) AddFuncSymbol(name string, params []string, body *parser.BlockNode, newEnv *Env) {
	e.Symbols[name] = &FuncSymbol{
		Body:     body,
		Params:   params,
		TypeName: "nan",
		Env:      newEnv,
	}
}

func (e *Env) AddStructSymbol(name string, structEnv *Env) {
	e.Symbols[name] = &StructSymbol{
		TypeName:    name,
		Environment: structEnv,
	}
}

func (e *Env) AddStructMethod(
	structName,
	methodName string,
	params []string,
	body []parser.Node,
	nativeFn any) int {
	for env := e; env != nil; env = env.Parent {
		sym, exists := env.Symbols[structName]
		if !exists {
			continue
		}

		structSym, ok := sym.(*StructSymbol)
		if !ok {
			return -1
		}

		if structSym.Environment.Symbols == nil {
			structSym.Environment.Symbols = make(core.SymbolStore)
		}

		if _, exists := structSym.Environment.Symbols[methodName]; exists {
			return -1
		}

		funcSym := &FuncSymbol{
			Body:     &parser.BlockNode{Statements: body},
			Params:   params,
			TypeName: "nan",
			Env:      NewEnv(structSym.Environment, "method"),
		}

		structSym.Environment.Symbols[methodName] = funcSym
		env.Symbols[structName] = structSym
		return 0
	}
	return 1
}
