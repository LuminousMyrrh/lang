package main

import (
	"errors"
	"fmt"
)

func (e *Evaluator) genError(msg string, pos Position) {
	e.Errors = append(
		e.Errors, errors.New(fmt.Sprintf(
			"%d, %d: %s",
			pos.Row,
			pos.Column,
			msg,
			)),
		)
}

func (e *Evaluator) resolveType(value any, pos Position) string {
	switch v := value.(type) {
	case string:
		return "string"
	case int:
		return "int"
	case bool:
		return "bool"
	case *Env:
		return v.Type
	case []any:
		return "[]"
	case nilValue:
		return "nil"
	default:
		e.genError(fmt.Sprintf("Unknown type '%v'", v), pos)
		return ""
	}
}

