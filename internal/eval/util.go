package eval

import (
	"errors"
	"fmt"
	"lang/internal/env"
	"lang/internal/parser"
	"unicode"
)

func (e *Evaluator) genError(msg string, pos parser.Position) {
	e.Errors = append(
		e.Errors, errors.New(fmt.Sprintf(
			"%d, %d: %s",
			pos.Row,
			pos.Column,
			msg,
			)),
		)
}

func (e *Evaluator) resolveType(value any, pos parser.Position) string {
	switch v := value.(type) {
	case string:
		return "string"
	case int:
		return "int"
	case float64:
		return "float"
	case bool:
		return "bool"
	case *env.Env:
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

func capitalizeFirstLetter(s string) string {
    if len(s) == 0 {
        return s
    }
    runes := []rune(s)
    runes[0] = unicode.ToUpper(runes[0])
    return string(runes)
}

func UnwrapBuiltinValue(v any) any {
    if instEnv, ok := v.(*env.Env); ok {
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
