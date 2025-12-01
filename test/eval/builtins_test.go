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

func TestBuiltinInput(t *testing.T) {
    // Provide the input string to simulate user typing
    inputString := "test input\n"

    // Backup original stdin
    oldStdin := os.Stdin

    // Create a reader with the inputString and set to Stdin
    r, w, err := os.Pipe()
    if err != nil {
        t.Fatalf("Failed to create pipe: %v", err)
    }
    _, err = w.Write([]byte(inputString))
    if err != nil {
        t.Fatalf("Failed to write to pipe: %v", err)
    }
    w.Close()
    os.Stdin = r

    // Create a FunctionCallNode representing input()
    inputNode := parser.FunctionCallNode{
        Name: &parser.IdentifierNode{Name: "input"},
        Args: []parser.Node{},
    }

    // Setup ProgramNode with this call
    program := parser.ProgramNode{
        Nodes: []parser.Node{&inputNode},
    }

    evaluator := eval.NewEvaluatorAutoEnv(&program)
    // Evaluate the input Node via evaluator to trigger builtinInput
    result := evaluator.EvalNode(&inputNode)

    // Restore original stdin
    os.Stdin = oldStdin

    // Check for errors
    if len(evaluator.Errors) > 0 {
        t.Fatalf("Unexpected errors: %v", evaluator.Errors)
    }

    // Assert that result matches inputString without newline
    expected := "test input"
    if result != expected {
        t.Errorf("Expected input result %q, got %q", expected, result)
    }
}
