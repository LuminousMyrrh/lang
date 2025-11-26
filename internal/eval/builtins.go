package eval

import (
	"fmt"
	"io"
	"lang/internal/env"
	"lang/internal/parser"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type BuiltinFunction func(e *Evaluator, args []parser.Node, pos parser.Position) any

func decodeEscapeSequences(input string) (string, error) {
	// Wrap in quotes to make it a valid Go string literal
	quoted := "\"" + input + "\""
	return strconv.Unquote(quoted)
}

func builtinPrint(e *Evaluator, args []parser.Node, pos parser.Position) any {
	for _, arg := range args {
		val := UnwrapBuiltinValue(e.eval(arg))
		fmt.Print(val)
	}
	fmt.Println()
	return nilValue{}
}

func builtinPrintf(e *Evaluator, args []parser.Node, pos parser.Position) any {
	for _, arg := range args {
		val := UnwrapBuiltinValue(e.eval(arg))
		if s, ok := val.(string); ok {
			decoded, err := decodeEscapeSequences(s)
			if err != nil {
				e.genError(
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
	return nilValue{}
}

func builtinPrintln(e *Evaluator, args []parser.Node, pos parser.Position) any {
	for _, arg := range args {
		val := UnwrapBuiltinValue(e.eval(arg))
		if s, ok := val.(string); ok {
			decoded, err := decodeEscapeSequences(s)
			if err != nil {
				e.genError(
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
	return nilValue{}
}

func builtinType(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) != 1 {
		e.genError("type: expects one argument", pos)
		return nil
	}
	val := UnwrapBuiltinValue(e.eval(args[0]))
	return e.resolveType(val, pos)
}

func builtinInput(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) > 1 {
		e.genError("input: expects zero or one argument", pos)
		return nil
	}
	if len(args) == 1 {
		fmt.Print(UnwrapBuiltinValue(e.eval(args[0])))
	}
	var input string
	fmt.Scanln(&input)
	return input
}

func builtinInt(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) != 1 {
		e.genError("atoi: expects one argument", pos)
		return nil
	}
	val := UnwrapBuiltinValue(e.eval(args[0]))
	switch s := val.(type) {
	case string:
		result, err := strconv.Atoi(s)
		if err != nil {
			e.genError("int: invalid string format", pos)
			return nil
		}
		return result
	case float64:
		return int(s)
	case int:
		return s
	default:
		e.genError(fmt.Sprintf(
			"Unsupported type: %T", s), pos)
		return nil
	}
}

func builtinFloat(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) != 1 {
		e.genError("float: expects one argument", pos)
		return nil
	}
	val := UnwrapBuiltinValue(e.eval(args[0]))
	switch s := val.(type) {
	case string:
		result, err := strconv.ParseFloat(s, 64)
		if err != nil {
			e.genError("int: invalid string format", pos)
			return nil
		}
		return result
	case int:
		return float64(s)
	case float64:
		return s
	default:
		e.genError(fmt.Sprintf(
			"Unsupported type: %T", s), pos)
		return nil
	}
}

func builtinString(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) != 1 {
		e.genError("itoa: expects one argument", pos)
		return nil
	}
	val := UnwrapBuiltinValue(e.eval(args[0]))
	switch v := val.(type) {
	case int:
		return e.createString(strconv.Itoa(v))
	case float64:
		return e.createString(strconv.FormatFloat(v, 'g', -1, 64))
	case string:
		return e.createString(v)
	default:
		e.genError(fmt.Sprintf("Unsupported type: %T", v), pos)
		return nil
	}
}

func builtinLen(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) != 1 {
		e.genError("len: expects one argument", pos)
		return nil
	}
	arr := UnwrapBuiltinValue(e.eval(args[0]))
	switch a := arr.(type) {
	case []any:
		return len(a)
	case string:
		return len(a)
	}
	e.genError("len: argument must be array/string", pos)
	return nil
}

// builtinReadAll implements the 'readAll' builtin function.
// Usage: readAll("filename") -> file contents as bytes
func builtinReadAll(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) == 0 {
		e.genError("Function 'readAll' expects at least one argument, but 0 were provided", pos)
		return nil
	}
	evaledFileName := UnwrapBuiltinValue(e.eval(args[0]))
	fileName, ok := evaledFileName.(string)
	if !ok {
		e.genError(fmt.Sprintf(
			"Argument to 'readAll' must be a file name as string, got: %v", evaledFileName), pos)
		return nil
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		e.genError(err.Error(), pos)
		return nil
	}
	return e.createString(string(data))
}

func builtinFetch(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) != 1 {
		e.genError(
			"Function 'fetch' expect only one argument", pos)
		return nil
	}

	evaledName := UnwrapBuiltinValue(e.eval(args[0]))
	name, ok := evaledName.(string)
	if !ok {
		e.genError(
			"Argument should be string", pos)
		return nil
	}
	if !strings.HasPrefix("https://", name) {
		name = "https://" + name
	}

	resp, err := http.Get(name)
	if err != nil {
		e.genError(fmt.Sprintf("Failed to fetch: %s", err), pos)
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		e.genError(fmt.Sprintf("Failed to read body: %s", err), pos)
		return nil
	}
	return body
}

func builtinWrite(e *Evaluator, args []parser.Node, pos parser.Position) any {
	if len(args) != 2 {
		e.genError(
			"Function 'write' expect two arguments", pos)
		return nil
	}
	evaledFileName := e.eval(args[0])
	evaledValue := e.eval(args[1])
	if fileName, ok := evaledFileName.(string); ok {
		var value []byte
		switch val := evaledValue.(type) {
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
				e.genError(fmt.Sprintf(
					"Usupported type for 'write': %T", val),
					pos)
				return nil
			}
		}
		data := os.WriteFile(fileName, value, 0644)
		return data
	} else {
		e.genError(fmt.Sprintf(
			"Argument in 'readAll' should be file name as string, but got: %v",
			evaledFileName),
			pos)
		return nil
	}
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

func (e *Evaluator) initBuiltintMethods() int {
	stringEnv := e.currentEnv.FindStructSymbol("string")
	if stringEnv != nil {
		stringEnv.Symbols["substring"] = &env.FuncSymbol{
			NativeFunc: stringSubstring,
			TypeName:   "string",
		}
		stringEnv.Symbols["capitalize"] = &env.FuncSymbol{
			NativeFunc: stringCapitalize,
			TypeName:   "string",
		}
		stringEnv.Symbols["contains"] = &env.FuncSymbol{
			NativeFunc: stringContains,
			TypeName:   "string",
		}
		stringEnv.Symbols["empty"] = &env.FuncSymbol{
			NativeFunc: stringEmpty,
			TypeName:   "string",
		}
		stringEnv.Symbols["isDigit"] = &env.FuncSymbol{
			NativeFunc: stringIsDigit,
			TypeName:   "string",
		}
		stringEnv.Symbols["isAlph"] = &env.FuncSymbol{
			NativeFunc: stringIsAlph,
			TypeName:   "string",
		}
	} else {
		return -1
	}
	return 0
}
