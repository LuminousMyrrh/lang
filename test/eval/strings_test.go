package eval

import (
	"lang/internal/core"
	"lang/internal/eval"
	"lang/internal/parser"
	"reflect"
	"testing"
)

func ident(name string) *parser.IdentifierNode {
	return &parser.IdentifierNode{Name: name}
}

func varDef(name string, value parser.Node) parser.VarDefNode {
	return parser.VarDefNode{
		Name:     name,
		Value:    value,
		Position: parser.Position{1, 1},
	}
}

func checkSymbolType(symbol core.Symbol, expectedType any) bool {
	actualType := reflect.TypeOf(symbol.Value())
	return actualType == reflect.TypeOf(expectedType)
}

func TestCapitalize(t *testing.T) {
	varDecl := varDef(
		"testVar",
		&parser.LiteralNode{Value: "test"},
	)

	input := parser.ProgramNode{
		Nodes: []parser.Node{
			&varDecl,
			&parser.StructMethodCall{
				Caller:     ident("testVar"),
				MethodName: "capitalize",
				Position:   parser.Position{1, 1},
			},
		},
	}

	expectErrs := 0
	expect := "This test"
	evaluator := eval.NewEvaluatorAutoEnv(&input)
	evaluator.Eval()
	if expectErrs != len(evaluator.Errors) {
		t.Errorf("Expected %v errors, got %v", expectErrs, evaluator.Errors)
	}

	varSym := evaluator.Environment.FindSymbol("testVar")

	if varSym.Value() == expect {
		t.Errorf("Expect '%s', got '%s'", expect, varSym.Value())
	}
}

func TestStringDoesContains(t *testing.T) {
	varDecl := varDef(
		"testVar",
		&parser.LiteralNode{Value: "it does contain test"},
	)

	input := parser.ProgramNode{
		Nodes: []parser.Node{
			&varDecl,
			&parser.StructMethodCall{
				Caller:     ident("testVar"),
				MethodName: "contains",
				Args: []parser.Node{
					&parser.LiteralNode{Value: "test"},
				},
				Position: parser.Position{1, 1},
			},
		},
	}

	expectErrs := 0
	expect := "This test"
	evaluator := eval.NewEvaluatorAutoEnv(&input)
	evaluator.Eval()
	if expectErrs != len(evaluator.Errors) {
		t.Errorf("Expected %v errors, got %v", expectErrs, evaluator.Errors)
	}

	varSym := evaluator.Environment.FindSymbol("testVar")

	if varSym.Value() == expect {
		t.Errorf("Expect '%s', got '%s'", expect, varSym.Value())
	}
}
