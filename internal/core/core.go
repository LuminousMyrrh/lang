package core

import "lang/internal/parser"

type EvaluatorInterface interface {
    eval(stmt parser.Node) any
    genError(msg string, pos parser.Position)
}
