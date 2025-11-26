package env

import (
	"fmt"
	"strings"
)

type Env struct {
	Type    string
	Symbols map[string]Symbol
	Parent  *Env
}

func NewEnv(parent *Env, envType string) *Env {
	return &Env{
		Type:    envType,
		Symbols: make(map[string]Symbol),
		Parent:  parent,
	}
}

func (e *Env) String() string {
	var lines []string
	for name, val := range e.Symbols {
		lines = append(lines, fmt.Sprintf("%s(%T) -> %v\n", name, val, val.Name()))
	}
	return strings.Join(lines, "\n")
}

func (e *Env) InitStringBuiltin() {
	stringSymbol := e.FindStructSymbol("string")
	if stringSymbol != nil {
		stringSymbol.Symbols["substring"] = &symbol.FuncSymbol{
			NativeFunc: stringSubstring,
			TypeName:   "string",
		}
		stringSymbol.Symbols["capitalize"] = &env.FuncSymbol{
			NativeFunc: stringCapitalize,
			TypeName:   "string",
		}
		stringSymbol.Symbols["contains"] = &env.FuncSymbol{
			NativeFunc: stringContains,
			TypeName:   "string",
		}
		stringSymbol.Symbols["empty"] = &env.FuncSymbol{
			NativeFunc: stringEmpty,
			TypeName:   "string",
		}
		stringSymbol.Symbols["isDigit"] = &env.FuncSymbol{
			NativeFunc: stringIsDigit,
			TypeName:   "string",
		}
		stringSymbol.Symbols["isAlph"] = &env.FuncSymbol{
			NativeFunc: stringIsAlph,
			TypeName:   "string",
		}
	} else {
		return -1
	}
}
