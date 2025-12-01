package core

import (
	"lang/internal/parser"
)

type Evaluator interface {
    EvalNode(stmt parser.Node) any
    GenError(msg string, pos parser.Position)
}

type ReturnValue struct {
    Value any
}

type BreakSignal struct {}

type NilValue struct {}
