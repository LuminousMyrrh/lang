package eval

import (
	"bytes"
	"lang/internal/eval"
	"lang/internal/parser"
	"os"
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

func TestVarDeclAndUse(t *testing.T) {
    // Program:
    // var x = 42
    // print(x)

    varDecl := &parser.VarDefNode{
        Name: "x",
        Value: &parser.LiteralNode{Value: 42},
    }

    printCall := &parser.FunctionCallNode{
        Name: &parser.IdentifierNode{Name: "print"},
        Args: []parser.Node{
            &parser.IdentifierNode{Name: "x"},
        },
    }

    program := parser.ProgramNode{
        Nodes: []parser.Node{varDecl, printCall},
    }

    var buf bytes.Buffer
    oldStdout := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    evaluator := eval.NewEvaluatorAutoEnv(&program)
    evaluator.Eval()

    w.Close()
    os.Stdout = oldStdout
    buf.ReadFrom(r)
    output := buf.String()

    expectedOutput := "42"

    if len(evaluator.Errors) != 0 {
        t.Fatalf("Unexpected errors: %v", evaluator.Errors)
    }
    if output != expectedOutput {
        t.Errorf("Expected output %q, got %q", expectedOutput, output)
    }
}
