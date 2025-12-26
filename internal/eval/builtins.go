package eval

import (
	"fmt"
	"io"
	"lang/internal/core"
	"lang/internal/env"
	"lang/internal/parser"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type BuiltinFunction func(e *Evaluator,
	args []parser.Node, pos parser.Position) any

func decodeEscapeSequences(input string) (string, error) {
	// Wrap in quotes to make it a valid Go string literal
	quoted := "\"" + input + "\""
	return strconv.Unquote(quoted)
}

func builtinPrint(e *Evaluator, args []parser.Node, pos parser.Position) any {
	for _, arg := range args {
		val := env.UnwrapBuiltinValue(e.EvalNode(arg))
		fmt.Print(val)
	}
	return core.NilValue{}
}

func builtinPrintf(e *Evaluator, args []parser.Node, pos parser.Position) any {
	for _, arg := range args {
		val := env.UnwrapBuiltinValue(e.EvalNode(arg))
		if s, ok := val.(string); ok {
			decoded, err := decodeEscapeSequences(s)
			if err != nil {
				e.GenError(
					"print: invalid escape sequence",
					pos,
				)
				return nil
			}
			fmt.Print(decoded)
		} else {
			fmt.Print(val)
		}
	}
	fmt.Println()
	return core.NilValue{}
}

func builtinPrintln(e *Evaluator, args []parser.Node, pos parser.Position) any {
	for _, arg := range args {
		val := env.UnwrapBuiltinValue(e.EvalNode(arg))
		if s, ok := val.(string); ok {
			decoded, err := decodeEscapeSequences(s)
			if err != nil {
				e.GenError(
					"println: invalid escape sequence",
					pos,
				)
				return nil
			}
			fmt.Print(decoded)
		} else {
			fmt.Print(val)
		}
	}
	fmt.Println()
	return core.NilValue{}
}

func builtinType(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) != 1 {
		e.GenError("type: expects one argument", pos)
		return nil
	}
	val := env.UnwrapBuiltinValue(e.EvalNode(args[0]))
	return e.resolveType(val, pos)
}

func builtinInput(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) > 1 {
		e.GenError("input: expects zero or one argument", pos)
		return nil
	}
	if len(args) == 1 {
		fmt.Print(env.UnwrapBuiltinValue(e.EvalNode(args[0])))
	}
	var input string
	fmt.Scanln(&input)
	return input
}

func builtinInt(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) != 1 {
		e.GenError("atoi: expects one argument", pos)
		return nil
	}
	val := env.UnwrapBuiltinValue(e.EvalNode(args[0]))
	switch s := val.(type) {
	case string:
		result, err := strconv.Atoi(s)
		if err != nil {
			e.GenError("int: invalid string format", pos)
			return nil
		}
		return result
	case float64:
		return int(s)
	case int:
		return s
	default:
		e.GenError(fmt.Sprintf(
			"Unsupported type: %T", s), pos)
		return nil
	}
}

func builtinFloat(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) != 1 {
		e.GenError("float: expects one argument", pos)
		return nil
	}
	val := env.UnwrapBuiltinValue(e.EvalNode(args[0]))
	switch s := val.(type) {
	case string:
		result, err := strconv.ParseFloat(s, 64)
		if err != nil {
			e.GenError("int: invalid string format", pos)
			return nil
		}
		return result
	case int:
		return float64(s)
	case float64:
		return s
	default:
		e.GenError(fmt.Sprintf(
			"Unsupported type: %T", s), pos)
		return nil
	}
}

func builtinString(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) != 1 {
		e.GenError("itoa: expects one argument", pos)
		return nil
	}
	val := env.UnwrapBuiltinValue(e.EvalNode(args[0]))
	switch v := val.(type) {
	case int:
		return e.createString(strconv.Itoa(v))
	case float64:
		return e.createString(strconv.FormatFloat(v, 'g', -1, 64))
	case string:
		return e.createString(v)
	default:
		e.GenError(fmt.Sprintf("Unsupported type: %T", v), pos)
		return nil
	}
}

func builtinLen(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) != 1 {
		e.GenError("len: expects one argument", pos)
		return nil
	}
	arr := env.UnwrapBuiltinValue(e.EvalNode(args[0]))
	switch a := arr.(type) {
	case []any:
		return len(a)
	case string:
		return len(a)
	}
	e.GenError("len: argument must be array/string", pos)
	return nil
}

// builtinReadAll implements the 'readAll' builtin function.
// Usage: readAll("filename") -> file contents as bytes
func builtinReadAll(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) == 0 {
		e.GenError("Function 'readAll' expects at least one argument, but 0 were provided", pos)
		return nil
	}
	EvalNodeedFileName := env.UnwrapBuiltinValue(e.EvalNode(args[0]))
	fileName, ok := EvalNodeedFileName.(string)
	if !ok {
		e.GenError(fmt.Sprintf(
			"Argument to 'readAll' must be a file name as string, got: %v", EvalNodeedFileName), pos)
		return nil
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		e.GenError(err.Error(), pos)
		return nil
	}
	return e.createString(string(data))
}

func builtinFetch(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) != 1 {
		e.GenError(
			"Function 'fetch' expect only one argument", pos)
		return nil
	}

	EvalNodeedName := env.UnwrapBuiltinValue(e.EvalNode(args[0]))
	name, ok := EvalNodeedName.(string)
	if !ok {
		e.GenError(
			"Argument should be string", pos)
		return nil
	}
	if !strings.HasPrefix("https://", name) {
		name = "https://" + name
	}

	resp, err := http.Get(name)
	if err != nil {
		e.GenError(fmt.Sprintf("Failed to fetch: %s", err), pos)
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		e.GenError(fmt.Sprintf("Failed to read body: %s", err), pos)
		return nil
	}
	return body
}

func builtinWrite(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) != 2 {
		e.GenError(
			"Function 'write' expect two arguments", pos)
		return nil
	}
	EvalNodeedFileName := e.EvalNode(args[0])
	EvalNodeedValue := e.EvalNode(args[1])
	if fileName, ok := EvalNodeedFileName.(string); ok {
		var value []byte
		switch val := EvalNodeedValue.(type) {
		case string:
			{
				value = []byte(val)
			}
		case int:
			{
				value = []byte(strconv.Itoa(val))
			}
		default:
			{
				e.GenError(fmt.Sprintf(
					"Usupported type for 'write': %T", val),
					pos)
				return nil
			}
		}
		data := os.WriteFile(fileName, value, 0644)
		return data
	} else {
		e.GenError(fmt.Sprintf(
			"Argument in 'readAll' should be file name as string, but got: %v",
			EvalNodeedFileName),
			pos)
		return nil
	}
}

func builtinMod(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) != 2 {
		e.GenError("float: expects two arguments", pos)
		return nil
	}
	first := env.UnwrapBuiltinValue(e.EvalNode(args[0]))
	second := env.UnwrapBuiltinValue(e.EvalNode(args[1]))

	if v1, ok := first.(int); ok {
		if v2, ok := second.(int); ok {
			return v1 % v2
		}
	}
	e.GenError(fmt.Sprintf("Types aren't the same: %T and %T", first, second), pos)
	return nil
}

func builtinOrd(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) != 1 {
		e.GenError("float: expects one argument", pos)
		return nil
	}
	symbol := env.UnwrapBuiltinValue(e.EvalNode(args[0]))

	if s, ok := symbol.(string); ok {
		rs := []rune(s);
		if len(rs) != 1 {
			return nil
		}
		
		return int(rs[0])
	}

	e.GenError(fmt.Sprintf("%T", symbol), pos)
	return nil
}

func (e *Evaluator) initBuiltintClasses() {
	stringEnv := env.NewEnv(nil, "string")
	stringEnv.AddVarSymbol(
		"value",
		"string",
		nil)
	e.currentEnv.AddStructSymbol("string", stringEnv)

	intEnv := env.NewEnv(nil, "int")
	intEnv.AddVarSymbol(
		"value",
		"int",
		nil)
	e.currentEnv.AddStructSymbol("int", intEnv)

	floatEnv := env.NewEnv(nil, "float")
	floatEnv.AddVarSymbol(
		"value",
		"float",
		nil)
	e.currentEnv.AddStructSymbol("float", floatEnv)
}

func (e *Evaluator) initBuiltinMethods() {
	builtins := map[string]BuiltinFunction{
		"printf":  builtinPrintf,
		"print":   builtinPrint,
		"println": builtinPrintln,
		"type":    builtinType,
		"input":   builtinInput,
		"int":     builtinInt,
		"float":   builtinFloat,
		"string":  builtinString,
		"len":     builtinLen,
		"readAll": builtinReadAll,
		"write":   builtinWrite,
		"fetch":   builtinFetch,
		"mod":     builtinMod,
		"ord":     builtinOrd,
	}

	e.Builtins = builtins
}
