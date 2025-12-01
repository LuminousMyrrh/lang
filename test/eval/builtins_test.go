package eval

import (
	"bytes"
	"lang/internal/eval"
	"lang/internal/parser"
	"os"
	"testing"
)

func TestBuiltinPrint(t *testing.T) {
	var buf bytes.Buffer

	// Save original stdout
	old := os.Stdout

	// Create a pipe
	r, w, _ := os.Pipe()
	os.Stdout = w

	input := parser.ProgramNode{
		Nodes: []parser.Node{
			&parser.FunctionCallNode{
				Name: &parser.IdentifierNode{Name: "print"},
				Args: []parser.Node{
					&parser.LiteralNode{Value: "this test"},
				},
			},
		},
	}
	expectErrs := 0
	expectOutput := "this test"
	evaluator := eval.NewEvaluatorAutoEnv(&input)
	evaluator.Eval()
	w.Close()
	if len(evaluator.Errors) != expectErrs {
		t.Errorf("Expected %v errors, got %v", expectErrs, evaluator.Errors)
	}

	os.Stdout = old
	buf.ReadFrom(r)
	output := buf.String()

	if expectOutput != output {
		t.Errorf("Expected '%s', got '%s'", expectOutput, output)
	}
}

func TestBuiltinLen(t *testing.T) {
	input := parser.ProgramNode{
		Nodes: []parser.Node{
			&parser.FunctionCallNode{
				Name: &parser.IdentifierNode{Name: "len"},
				Args: []parser.Node{
					&parser.ArrayNode{
						Elements: []parser.Node{
							&parser.LiteralNode{Value: 1},
							&parser.LiteralNode{Value: 2},
							&parser.LiteralNode{Value: 3},
						},
					},
				},
			},
		},
	}

	evaluator := eval.NewEvaluatorAutoEnv(&input)
	expectErrs := 0
	// expectValue := 3

	evaluator.Eval()
	if len(evaluator.Errors) != expectErrs {
		t.Errorf("Expected %v errors, got %v", expectErrs, evaluator.Errors)
	}

	// if result != expectValue {
	//     t.Errorf("Expected result %v, got %v", expectValue, result)
	// }
}

func _TestBuiltinInput(t *testing.T) {
    inputString := "test input\n"

    // Backup original stdin and stdout
    oldStdin := os.Stdin
    oldStdout := os.Stdout

    // Setup stdin pipe with inputString
    rIn, wIn, err := os.Pipe()
    if err != nil {
        t.Fatalf("Failed to create stdin pipe: %v", err)
    }
    _, err = wIn.Write([]byte(inputString))
    if err != nil {
        t.Fatalf("Failed to write to stdin pipe: %v", err)
    }
    wIn.Close()
    os.Stdin = rIn

    // Setup stdout pipe to capture output
    rOut, wOut, err := os.Pipe()
    if err != nil {
        t.Fatalf("Failed to create stdout pipe: %v", err)
    }
    os.Stdout = wOut

    // Build AST: println(input())
    inputCall := &parser.FunctionCallNode{
        Name: &parser.IdentifierNode{Name: "input"},
        Args: []parser.Node{},
    }

    printlnCall := &parser.FunctionCallNode{
        Name: &parser.IdentifierNode{Name: "println"},
        Args: []parser.Node{inputCall},
    }

    program := parser.ProgramNode{
        Nodes: []parser.Node{printlnCall},
    }

    evaluator := eval.NewEvaluatorAutoEnv(&program)
    evaluator.Eval()

    // Close and restore stdout & stdin
    wOut.Close()
    os.Stdout = oldStdout
    os.Stdin = oldStdin

    // Read and check captured stdout
    var outputBuf bytes.Buffer
    outputBuf.ReadFrom(rOut)
    output := outputBuf.String()

    expectedOutput := "test input\n"
    if len(evaluator.Errors) > 0 {
        t.Fatalf("Unexpected errors: %v", evaluator.Errors)
    }
    if output != expectedOutput {
        t.Errorf("Expected output %q, got %q", expectedOutput, output)
    }
}
