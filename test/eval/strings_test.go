package eval

import (
	"lang/internal/eval"
	"lang/internal/parser"
	"testing"
)

func ident(name string) parser.IdentifierNode {
	return parser.IdentifierNode{Name: name}
}

func varDef(name string, value parser.Node) parser.VarDefNode {
	return parser.VarDefNode{
		Name:     "testVar",
		Value:    value,
		Position: parser.Position{1, 1},
	}
}

func TestCapitalize(t *testing.T) {
	varDecl := parser.VarDefNode{
		Name:     "testVar",
		Value:    &parser.LiteralNode{Value: "test"},
		Position: parser.Position{1, 1},
	}

	input := parser.ProgramNode{
		Nodes: []parser.Node{
			&varDecl,
			&parser.FunctionCallNode{
				Name: &parser.IdentifierNode{Name: "print"},
				Args: []parser.Node{
					&parser.StructMethodCall{
						Caller:
					},
				},
			},
		},
	}

	expectErrs := 0
	expectOutput := "this test"
	evaluator := eval.NewEvaluatorAutoEnv(&input)
	evaluator.Eval()
}
