package main

import (
	"fmt"
	"strconv"
)

type BuiltinFunction func(e *Evaluator, args []Node, pos Position) any

func decodeEscapeSequences(input string) (string, error) {
    // Wrap in quotes to make it a valid Go string literal
    quoted := "\"" + input + "\""
    return strconv.Unquote(quoted)
}

func builtinPrint(e *Evaluator, args []Node, pos Position) any {
    for _, arg := range args {
        val := e.eval(arg)
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
    return 1
}

func builtinPrintln(e *Evaluator, args []Node, pos Position) any {
    for _, arg := range args {
        val := e.eval(arg)
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
    fmt.Println()
    return 1
}

func builtinType(e *Evaluator, args []Node, pos Position) any {
    if len(args) != 1 {
        e.genError("type: expects one argument", pos)
        return nil
    }
    val := e.eval(args[0])
    return e.resolveType(val, pos)
}

func builtinInput(e *Evaluator, args []Node, pos Position) any {
    if len(args) > 1 {
        e.genError("input: expects zero or one argument", pos)
        return nil
    }
    if len(args) == 1 {
        fmt.Print(e.eval(args[0]))
    }
    var input string
    fmt.Scanln(&input)
    return input
}

func builtinAtoi(e *Evaluator, args []Node, pos Position) any {
    if len(args) != 1 {
        e.genError("atoi: expects one argument", pos)
        return nil
    }
    val := e.eval(args[0])
    str, ok := val.(string)
    if !ok {
        e.genError("atoi: argument must be string", pos)
        return nil
    }
    result, err := strconv.Atoi(str)
    if err != nil {
        e.genError("atoi: invalid string format", pos)
        return nil
    }
    return result
}

func builtinItoa(e *Evaluator, args []Node, pos Position) any {
    if len(args) != 1 {
        e.genError("itoa: expects one argument", pos)
        return nil
    }
    val := e.eval(args[0])
    i, ok := val.(int)
    if !ok {
        e.genError("itoa: argument must be int", pos)
        return nil
    }
    return strconv.Itoa(i)
}

func builtinLen(e *Evaluator, args []Node, pos Position) any {
    if len(args) != 1 {
        e.genError("len: expects one argument", pos)
        return nil
    }
    arr := e.eval(args[0])
    if a, ok := arr.([]any); ok {
        return len(a)
    }
    e.genError("len: argument must be array", pos)
    return nil
}
