package eval

import (
	"fmt"
	"lang/internal/env"
	"lang/internal/lexer"
	"lang/internal/parser"
)

type returnValue struct {
    value any
}

type breakSignal struct {}

type nilValue struct {}

type Evaluator struct {
	Entry *parser.ProgramNode
	Environment *env.Env
	currentEnv *env.Env
	Errors []error
	parser *parser.Parser
	lexer *lexer.Lexer
	Builtins map[string]BuiltinFunction
}

func NewEvaluator(env *env.Env, entry *parser.ProgramNode) *Evaluator {
	builtins := map[string]BuiltinFunction {
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
		"fetch":   builtinFetch,
	}

	return &Evaluator{
		Environment: env,
		Entry: entry,
		currentEnv: env,
		Builtins: builtins,
	}
}

func (e *Evaluator) Eval() {
	e.initBuiltintClasses()
	e.initBuiltintMethods()

	for _, stmt := range e.Entry.Nodes{
		if e.eval(stmt) == nil {
			return
		}
	}

	e.Environment = e.currentEnv
}

func (e *Evaluator) eval(stmt parser.Node) any {
    // fmt.Printf("eval node type: %T\n", stmt)
    res := e.evalX(stmt)  // actual dispatch
    if _, ok := res.(returnValue); ok {
        // fmt.Printf("eval returned returnValue wrapping: %v\n", ret.value)
    } else {
        //fmt.Printf("eval returned: %v\n", res)
    }
    return res
}

func (e *Evaluator) evalX(stmt parser.Node) any {

	switch s := stmt.(type) {
	case *parser.VarDefNode:
		return e.evalVarDef(s)
	case *parser.FunctionDefNode:
		return e.evalFunctionDef(s)
	case *parser.BinaryOpNode:
		return e.evalBinary(s)
	case *parser.IdentifierNode:
		return e.evalIdentifier(s)
	case *parser.FunctionCallNode: 
		return e.evalFunctionCall(s)
	case *parser.LiteralNode: 
		return e.evalLiteral(s)
	case *parser.TrueNode:
		return e.evalTrue()
	case *parser.FalseNode:
		return e.evalFalse()
	case *parser.ReturnNode:
		return e.evalReturn(s)
	case *parser.AssignmentNode:
		return e.evalAssignment(s)
	case *parser.IfNode: 
		return e.evalIf(s)
	case *parser.BlockNode:
		return e.evalBlock(s)
	case *parser.UnaryOpNode:
		return e.evalUnary(s)
	case *parser.WhileNode:
		return e.evalWhile(s)
	case *parser.ForNode:
		return e.evalFor(s)
	case *parser.ArrayNode:
		return e.evalArray(s)
	case *parser.ArrayAccessNode:
		return e.evalArrayAccess(s)
	case *parser.ArrayAssign:
		return e.evalArrayAssign(s)
	case *parser.ImportNode: 
		return e.evalImport(s)
	case *parser.StructDefNode:
		return e.evalStructDef(s)
	case *parser.StructMethodDef:
		return e.evalStructMethodDef(s)
	case *parser.StructInitNode:
		return e.evalStructInit(s)
	case *parser.StructMethodCall:
		return e.evalStructMemberAccess(s)
	case *parser.NilNode:
		return e.evalNil(s)
	default: {
		e.genError(fmt.Sprintf(
			"Unknown node type: %T", s), parser.Position{Row: -1, Column: -1})
		return nil
	}
	}
}

func (e *Evaluator) createString(value string) *env.Env {
	stringEnv := e.currentEnv.FindStructSymbol("string")
	if stringEnv == nil {
		return nil
	}
	instEnv := env.NewEnv(stringEnv, "string")
	instEnv.AddVarSymbol("value", "string", value)
	return instEnv
}
