package core

import (
	"lang/internal/parser"
)

type Evaluator interface {
    EvalNode(stmt parser.Node) any
    GenError(msg string, pos parser.Position)
}

type Symbol interface {
	Value() any
	Type() string
}

type SymbolStore map[string]Symbol

type Env interface {
	Type() string
	Parent() Env
	Symbols() SymbolStore
}

type ReturnValue struct {
    Value any
}

type BreakSignal struct {}

type NilValue struct {}
