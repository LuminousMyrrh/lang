package main

import (
	"fmt"
)

type returnValue struct {
    value any
}

type breakSignal struct {}

type nilValue struct {}

type Evaluator struct {
	Entry *ProgramNode
	Environment *Env
	currentEnv *Env
	Errors []error
	parser *Parser
	lexer *Lexer
	Builtins map[string]BuiltinFunction
}

func (e *Evaluator) Eval(env *Env, entry *ProgramNode) {
	e.Environment = env
	e.Entry = entry
	e.currentEnv = env
	e.Builtins = map[string]BuiltinFunction {
		"printf":   builtinPrintf,
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
	}
	e.initBuiltintClasses()
	e.initBuiltintMethods()

	for _, stmt := range entry.Nodes{
		if e.eval(stmt) == nil {
			return
		}
	}

	e.Environment = e.currentEnv
}

func (e *Evaluator) eval(stmt Node) any {
    // fmt.Printf("eval node type: %T\n", stmt)
    res := e.evalX(stmt)  // actual dispatch
    if _, ok := res.(returnValue); ok {
        // fmt.Printf("eval returned returnValue wrapping: %v\n", ret.value)
    } else {
        //fmt.Printf("eval returned: %v\n", res)
    }
    return res
}

func (e *Evaluator) evalX(stmt Node) any {

	switch s := stmt.(type) {
	case *VarDefNode:
		return e.evalVarDef(s)
	case *FunctionDefNode:
		return e.evalFunctionDef(s)
	case *BinaryOpNode:
		return e.evalBinary(s)
	case *IdentifierNode:
		return e.evalIdentifier(s)
	case *FunctionCallNode: 
		return e.evalFunctionCall(s)
	case *LiteralNode: 
		return e.evalLiteral(s)
	case *TrueNode:
		return e.evalTrue()
	case *FalseNode:
		return e.evalFalse()
	case *ReturnNode:
		return e.evalReturn(s)
	case *AssignmentNode:
		return e.evalAssignment(s)
	case *IfNode: 
		return e.evalIf(s)
	case *BlockNode:
		return e.evalBlock(s)
	case *UnaryOpNode:
		return e.evalUnary(s)
	case *WhileNode:
		return e.evalWhile(s)
	case *ForNode:
		return e.evalFor(s)
	case *ArrayNode:
		return e.evalArray(s)
	case *ArrayAccessNode:
		return e.evalArrayAccess(s)
	case *ArrayAssign:
		return e.evalArrayAssign(s)
	case *ImportNode: 
		return e.evalImport(s)
	case *StructDefNode:
		return e.evalStructDef(s)
	case *StructMethodDef:
		return e.evalStructMethodDef(s)
	case *StructInitNode:
		return e.evalStructInit(s)
	case *StructMethodCall:
		return e.evalStructMemberAccess(s)
	case *NilNode:
		return e.evalNil(s)
	default: {
		e.genError(fmt.Sprintf("Unknown node type: %T", s), Position{-1, -1})
		return nil
	}
	}
}

func (e *Evaluator) createString(value string) *Env {
	stringEnv := e.currentEnv.FindStructSymbol("string")
	if stringEnv == nil {
		return nil
	}
	instEnv := NewEnv(stringEnv, "string")
	instEnv.AddVarSymbol("value", "string", value)
	return instEnv
}
