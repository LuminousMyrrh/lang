package parser

import (
	"reflect"
	"testing"
)

func TestParser_ParseVarDef(t *testing.T) {
    tokens := []*Token{
        {TType: Var, Lexeme: "var", Line: 1, Column: 1},
        {TType: Identifier, Lexeme: "x", Line: 1, Column: 5},
        {TType: Assign, Lexeme: "=", Line: 1, Column: 7},
        {TType: Digit, Lexeme: "42", Line: 1, Column: 9},
        {TType: Semicolon, Lexeme: ";", Line: 1, Column: 11},
    }

    parser := NewParser(NewEnv(nil), tokens)
    program, errs := parser.Parse()

    if len(errs) > 0 {
        t.Fatalf("Unexpected parse errors: %v", errs)
    }

    // Expected AST: ProgramNode with one VarDefNode
    want := &ProgramNode{
        Nodes: []Node{
            VarDefNode{
                Name:  "x",
                Value: LiteralNode{Value: 42},
            },
        },
    }

    if !reflect.DeepEqual(program, want) {
        t.Errorf("Parse() = %#v; want %#v", program, want)
    }
}

func TestParser_ParseVarDefWithFuncCallAndBinaryOp(t *testing.T) {
    tokens := []*Token{
        {TType: Var, Lexeme: "var", Line: 1, Column: 1},
        {TType: Identifier, Lexeme: "x", Line: 1, Column: 5},
        {TType: Assign, Lexeme: "=", Line: 1, Column: 7},
        {TType: Identifier, Lexeme: "y", Line: 1, Column: 9},
        {TType: LParen, Lexeme: "(", Line: 1, Column: 10},
        {TType: RParen, Lexeme: ")", Line: 1, Column: 11},
        {TType: Plus, Lexeme: "+", Line: 1, Column: 13},
        {TType: Digit, Lexeme: "10", Line: 1, Column: 15},
        {TType: Semicolon, Lexeme: ";", Line: 1, Column: 17},
    }

    parser := NewParser(NewEnv(nil), tokens)
    program, errs := parser.Parse()

    if len(errs) > 0 {
        t.Fatalf("Unexpected parse errors: %v", errs)
    }

    // Expected AST: var x = y() + 10;
    want := &ProgramNode{
        Nodes: []Node{
            VarDefNode{
                Name: "x",
                Value: BinaryOpNode{
                    Op: "+",
                    Left: FunctionCallNode{
                        Name: "y",
                        Args: []Node{},
                    },
                    Right: LiteralNode{Value: 10},
                },
            },
        },
    }

    if !reflect.DeepEqual(program, want) {
        t.Errorf("Parse() = %#v; want %#v", program, want)
    }
}

func TestParser_ParseSimpleIf(t *testing.T) {
    tokens := []*Token{
        {TType: If, Lexeme: "if", Line: 1, Column: 1},
        {TType: LParen, Lexeme: "(", Line: 1, Column: 3},
        {TType: Identifier, Lexeme: "x", Line: 1, Column: 4},
        {TType: Equals, Lexeme: "==", Line: 1, Column: 6},
        {TType: Digit, Lexeme: "1", Line: 1, Column: 9},
        {TType: RParen, Lexeme: ")", Line: 1, Column: 10},
        {TType: LCurly, Lexeme: "{", Line: 1, Column: 12},
        {TType: Identifier, Lexeme: "print", Line: 2, Column: 5},
        {TType: LParen, Lexeme: "(", Line: 2, Column: 10},
        {TType: String, Lexeme: "ok", Line: 2, Column: 11},
        {TType: RParen, Lexeme: ")", Line: 2, Column: 15},
        {TType: Semicolon, Lexeme: ";", Line: 2, Column: 16},
        {TType: RCurly, Lexeme: "}", Line: 3, Column: 1},
    }

    parser := NewParser(NewEnv(nil), tokens)
    program, errs := parser.Parse()

    if len(errs) > 0 {
        t.Fatalf("Unexpected parse errors: %v", errs)
    }

    want := &ProgramNode{
        Nodes: []Node{
            IfNode{
                Condition: BinaryOpNode{
                    Op: "==",
                    Left: IdentifierNode{Name: "x"},
                    Right: LiteralNode{Value: 1},
                },
                ThenBranch: &BlockNode{
                    Statements: []Node{
                        ExpressionStatementNode{
                            Expr: FunctionCallNode{
                                Name: "print",
                                Args: []Node{
                                    LiteralNode{Value: "ok"},
                                },
                            },
                        },
                    },
                },
                ElseBranch: nil,
            },
        },
    }

    if !reflect.DeepEqual(program, want) {
        t.Errorf("Parse() = %#v; want %#v", program, want)
    }
}

func TestParser_ParseIfElse(t *testing.T) {
    tokens := []*Token{
        {TType: If, Lexeme: "if", Line: 1, Column: 1},
        {TType: LParen, Lexeme: "(", Line: 1, Column: 3},
        {TType: Identifier, Lexeme: "flag", Line: 1, Column: 4},
        {TType: RParen, Lexeme: ")", Line: 1, Column: 8},
        {TType: LCurly, Lexeme: "{", Line: 1, Column: 10},
        {TType: Identifier, Lexeme: "print", Line: 2, Column: 5},
        {TType: LParen, Lexeme: "(", Line: 2, Column: 10},
        {TType: String, Lexeme: "yes", Line: 2, Column: 11},
        {TType: RParen, Lexeme: ")", Line: 2, Column: 16},
        {TType: Semicolon, Lexeme: ";", Line: 2, Column: 17},
        {TType: RCurly, Lexeme: "}", Line: 3, Column: 1},
        {TType: Else, Lexeme: "else", Line: 3, Column: 3},
        {TType: LCurly, Lexeme: "{", Line: 3, Column: 8},
        {TType: Identifier, Lexeme: "print", Line: 4, Column: 5},
        {TType: LParen, Lexeme: "(", Line: 4, Column: 10},
        {TType: String, Lexeme: "no", Line: 4, Column: 11},
        {TType: RParen, Lexeme: ")", Line: 4, Column: 15},
        {TType: Semicolon, Lexeme: ";", Line: 4, Column: 16},
        {TType: RCurly, Lexeme: "}", Line: 5, Column: 1},
    }

    parser := NewParser(NewEnv(nil), tokens)
    program, errs := parser.Parse()

    if len(errs) > 0 {
        t.Fatalf("Unexpected parse errors: %v", errs)
    }

    want := &ProgramNode{
        Nodes: []Node{
            IfNode{
                Condition: IdentifierNode{Name: "flag"},
                ThenBranch: &BlockNode{
                    Statements: []Node{
                        ExpressionStatementNode{
                            Expr: FunctionCallNode{
                                Name: "print",
                                Args: []Node{
                                    LiteralNode{Value: "yes"},
                                },
                            },
                        },
                    },
                },
                ElseBranch: &BlockNode{
                    Statements: []Node{
                        ExpressionStatementNode{
                            Expr: FunctionCallNode{
                                Name: "print",
                                Args: []Node{
                                    LiteralNode{Value: "no"},
                                },
                            },
                        },
                    },
                },
            },
        },
    }

    if !reflect.DeepEqual(program, want) {
        t.Errorf("Parse() = %#v; want %#v", program, want)
    }
}
