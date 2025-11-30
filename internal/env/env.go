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

