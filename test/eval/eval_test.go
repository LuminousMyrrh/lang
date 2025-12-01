package eval

import (
	"lang/internal/eval"
	"lang/internal/parser"
	"testing"
)

func TestEmptyNodes(t *testing.T) {
	input := parser.ProgramNode{
		Nodes: []parser.Node{},
	}
	expectErrs := 0
	evaluator := eval.NewEvaluatorAutoEnv(&input)
	evaluator.Eval()
	if len(evaluator.Errors) != expectErrs {
		t.Errorf("Expected %v errors, got %v", expectErrs, evaluator.Errors)
	}
}

func TestBasic(t *testing.T) {
	input := parser.ProgramNode{
		Nodes: []parser.Node{
			&parser.VarDefNode{
				Position: parser.Position{Row: 0, Column: 0},
				Name:     "test",
				Value: &parser.LiteralNode{
					Value: 10,
				},
			},
		},
	}
	expectErrs := 0
	evaluator := eval.NewEvaluatorAutoEnv(&input)
	evaluator.Eval()
	if len(evaluator.Errors) != expectErrs {
		t.Errorf("Expected %v errors, got %v", expectErrs, evaluator.Errors)
	}
}

