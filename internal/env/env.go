package env

import (
	"fmt"
	"lang/internal/core"
	"strings"
)

type Env struct {
	Type    string
	Symbols map[string]core.Symbol
	Parent  *Env
}

func NewEnv(parent *Env, envType string) *Env {
	return &Env{
		Type:    envType,
		Symbols: make(map[string]core.Symbol),
		Parent:  parent,
	}
}

func (e *Env) ListSymbols() string {
	var lines []string
	for name, val := range e.Symbols {
		lines = append(lines, fmt.Sprintf("%s(%T) -> %v\n", name, val.Type(), val))
	}
	return strings.Join(lines, "\n")
}

func UnwrapBuiltinValue(v any) any {
	if instEnv, ok := v.(*Env); ok {
		if instEnv.Parent != nil {
			pName := instEnv.Parent.Type
			if pName == "string" || pName == "int" || pName == "float" {
				if valueSym, ok := instEnv.Symbols["value"]; ok {
					return valueSym.Value()
				}
			}
		}
	}
	return v
}
