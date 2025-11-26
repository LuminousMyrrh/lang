package eval

import (
	"errors"
	"fmt"
	"lang/internal/env"
	"lang/internal/parser"
	"strings"
	"unicode"
)

func getValue(self *env.Env) (string, error) {
	valSym, ok := self.Symbols["value"]
	if !ok {
		return "", errors.New("Struct instance does not have a 'value' field")
	}

	varSym, ok := valSym.(*env.VarSymbol)
	if !ok {
		return "", errors.New("'value' field is not a variable symbol")
	}

	s, ok := varSym.Value().(string)
	if !ok {
		return "", errors.New("'value' field is not a string")
	}

	return s, nil
}

func stringSubstring(e *Evaluator, self *env.Env, args []any, pos parser.Position) any {
	s, err := getValue(self)
	if err != nil {
		e.genError(err.Error(), pos)
		return nil
	}

	if len(args) != 2 {
		e.genError("substring: expects exactly 2 arguments", pos)
		return nil
	}
	from, ok1 := args[0].(int)
	to, ok2 := args[1].(int)
	if !ok1 || !ok2 {
		e.genError("substring: arguments must be integers", pos)
		return nil
	}
	if from < 0 || to >= len(s) || from > to {
		e.genError(fmt.Sprintf(
			"substring: invalid indices '%d and %d' with length %d",
			from, to, len(s)),
			pos)
		return nil
	}
	return s[from : to+1]
}

func stringCapitalize(e *Evaluator, self *env.Env, args []any, pos parser.Position) any {
	if len(args) != 0 {
		e.genError("'capitalize' doesn't accept any arguments", pos)
		return nil
	}
	s, err := getValue(self)
	if err != nil {
		e.genError(err.Error(), pos)
		return nil
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	s = string(runes)
	self.UpdateSymbol("value", s, "string")

	return nilValue{}
}

// func stringFind(e *Evaluator, self *env.Env, args []any, pos parser.Position) any {
// 	if len(args) < 2 || len(args) > 3 {
// 		e.genError("'find' should have at least two arguments", pos)
// 		return nil
// 	}
// 
// 	subs, ok := args[0].(string)
// 	if !ok {
// 		e.genError("First argument in 'find' should be string", pos)
// 		return nil
// 	}
// 	start, ok := args[1].(int)
// 	if !ok {
// 		e.genError("Second argument in 'find' should be int", pos)
// 		return nil
// 	}
// 
// 	s, err := getValue(self)
// 	if err != nil {
// 		e.genError(err.Error(), pos)
// 		return nil
// 	}
// 
// 	end := len(s)
// 	if len(args) == 3 {
// 		end, ok = args[2].(int)
// 		if !ok {
// 			e.genError("Third argument in 'find' should be int", pos)
// 			return nil
// 		}
// 	}
// 
// 	var index int = -1
// 
// 
// 	return index
// }

func stringContains(e *Evaluator, self *env.Env, args []any, pos parser.Position) any {
	if len(args) != 1 {
		e.genError("'contains' accept exactly one argument", pos)
		return nil
	}

	subs, ok := args[0].(string)
	if !ok {
		e.genError("Argument in 'contains' should be a string", pos)
		return nil
	}

	s, err := getValue(self)
	if err != nil {
		e.genError(err.Error(), pos)
		return nil
	}

	return strings.Contains(s, subs)
}

func stringEmpty(e *Evaluator, self *env.Env, args []any, pos parser.Position) any {
	if len(args) != 0 {
		e.genError("'empty' doesn't accent any arguments", pos)
		return nil
	}

	s, err := getValue(self)
	if err != nil {
		e.genError(err.Error(), pos)
		return nil
	}

	return len(s) == 0
}

func stringIsDigit(e *Evaluator, self *env.Env, args []any, pos parser.Position) any {
	if len(args) != 0 {
		e.genError("'empty' doesn't accent any arguments", pos)
		return nil
	}

	s, err := getValue(self)
	if err != nil {
		e.genError(err.Error(), pos)
		return nil
	}

	for _, ch := range s {
		if !unicode.IsDigit(ch) {
			return false
		}
	}

	return true
}

func stringIsAlph(e *Evaluator, self *env.Env, args []any, pos parser.Position) any {
	if len(args) != 0 {
		e.genError("'empty' doesn't accent any arguments", pos)
		return nil
	}

	s, err := getValue(self)
	if err != nil {
		e.genError(err.Error(), pos)
		return nil
	}

	for _, ch := range s {
		if !unicode.IsLetter(ch) {
			return false
		}
	}

	return true
}
