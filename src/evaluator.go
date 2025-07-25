package main

import (
	"fmt"
	"os"
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
	e.Builtins = map[string]BuiltinFunction{
		"print":   builtinPrint,
		"println": builtinPrintln,
		"type":    builtinType,
		"input":   builtinInput,
		"atoi":    builtinAtoi,
		"itoa":    builtinItoa,
		"len":     builtinLen,
	}


	for _, stmt := range entry.Nodes{
		if val := e.eval(stmt); val == nil {
			return
		}
	}

	e.Environment = e.currentEnv
}

func (e *Evaluator) eval(stmt Node) any {
	// fmt.Printf("Evaluating node: %T\n", stmt)

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

func (e *Evaluator) evalImport(stmt *ImportNode) any {
	fileName := stmt.File + ".lang"
	data, err := os.ReadFile(fileName)
	if err != nil {
		e.genError(fmt.Sprintf(
			"Failed to read imported file: %s", err),
			stmt.Position,
			)
		return nil
	}

	content := string(data)
	e.lexer = &Lexer{};
	toks, err := e.lexer.Read(content)
	if err != nil {
		e.genError(fmt.Sprintf("Failed to read file: %s", err),
			stmt.Position,
			)
		return nil
	}
	e.parser = NewParser(toks)
	mnode, errs := e.parser.Parse()
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		return nil
	}

	// -----------------

	if len(stmt.Symbols) == 0 { 
		// mnode.Print() 
		evaluator := Evaluator{}
		evaluator.Eval(NewEnv(nil, "global"), mnode)
		if len(evaluator.Errors) > 0 {
			for _, err := range evaluator.Errors {
				fmt.Println("Runtime error: ", err)
			}
			return nil
		}

		return 1

	} else {
		for _, symbol := range stmt.Symbols {
			node, err := mnode.Find(symbol)
			if err != nil {
				e.genError(err.Error(), stmt.Position)
				return nil
			}
			e.eval(node)
		}
		return 1

	}
}

