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
